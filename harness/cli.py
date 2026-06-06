#!/usr/bin/env python3
"""Build and deploy game-server-manager services to a local Kubernetes cluster."""

from __future__ import annotations

import enum
import os
import shlex
import subprocess
import sys
import tempfile
from pathlib import Path
from typing import Optional

import typer

ROOT = Path(__file__).resolve().parent.parent
SRC = ROOT / "src"
INFRA = ROOT / "infra" / "k8s"

SERVICES: dict[str, Path] = {
    "api": SRC / "services" / "control-plane-api",
    "web": SRC / "services" / "control-plane-web",
    "bot": SRC / "discord" / "control-bot",
}

SERVICE_PORTS: dict[str, int] = {
    "api": 8080,
    "web": 3000,
}

IMAGES: dict[str, str] = {
    "api": "control-plane-api",
    "web": "control-plane-web",
    "bot": "control-bot",
}

OVERLAYS = ["local", "prod"]

class Cluster(enum.Enum):
    KIND = "kind"
    MINIKUBE = "minikube"
    RANCHER_DESKTOP = "rancher-desktop"
    OTHER = "other"  # kubectl available but context doesn't match known clusters

# Maps kubectl context prefixes to Cluster members.
# Kind uses "kind-<name>" contexts; minikube and Rancher Desktop match exactly.
CLUSTER_PREFIXES: dict[str, Cluster] = {
    "kind-": Cluster.KIND,
    "minikube": Cluster.MINIKUBE,
    "rancher-desktop": Cluster.RANCHER_DESKTOP,
}

app = typer.Typer(help="Build and deploy game-server-manager services to a local Kubernetes cluster.")


# ---------------------------------------------------------------------------
# helpers
# ---------------------------------------------------------------------------


def run(cmd: list[str], **kwargs) -> subprocess.CompletedProcess:
    print(f"  \x1b[36m>>> {' '.join(shlex.quote(a) for a in cmd)}\x1b[0m")
    return subprocess.run(cmd, check=True, **kwargs)


def run_quiet(cmd: list[str], **kwargs) -> str:
    """Run a command and return its stdout (stripped), swallows stderr."""
    return subprocess.run(
        cmd, capture_output=True, text=True, **kwargs
    ).stdout.strip()


def container_cli() -> str:
    """Return 'podman' or 'docker', preferring whichever is installed."""
    for candidate in ("podman", "docker"):
        try:
            subprocess.run(
                [candidate, "version"], capture_output=True, check=True
            )
            return candidate
        except Exception:
            continue
    sys.exit(
        "Neither podman nor docker is installed/available on PATH"
    )


def detect_cluster() -> Cluster | None:
    """Return the local cluster type or None if kubectl is unavailable."""
    try:
        ctx = run_quiet(["kubectl", "config", "current-context"])
    except Exception:
        return None

    for prefix, member in CLUSTER_PREFIXES.items():
        if ctx.startswith(prefix):
            return member

    return Cluster.OTHER


def load_images(images: list[str]) -> None:
    """Push images into the local cluster cache if needed."""
    cluster = detect_cluster()
    cli = container_cli()

    if cluster is Cluster.MINIKUBE:
        for img in images:
            if cli == "podman":
                fd, path = tempfile.mkstemp(suffix=".tar")
                os.close(fd)
                try:
                    run(["podman", "save", img, "-o", path])
                    run(["minikube", "image", "load", path])
                finally:
                    Path(path).unlink(missing_ok=True)
            else:
                run(["minikube", "image", "load", img])

    elif cluster is Cluster.KIND:
        for img in images:
            run(["kind", "load", "docker-image", img])

    # Docker Desktop and Rancher Desktop share the host Docker daemon —
    # images are already available.
    elif cluster in (None, Cluster.RANCHER_DESKTOP, Cluster.OTHER):
        pass

    else:
        print(f"  (skipped image load — unknown cluster {cluster!r})")


# ---------------------------------------------------------------------------
# commands
# ---------------------------------------------------------------------------


@app.command()
def build(
    services: Optional[list[str]] = typer.Argument(None, help="Services to build (default: all). Choices: api, web, bot"),
    tag: str = typer.Option("latest", "--tag", help="Image tag"),
) -> None:
    if services:
        invalid = set(services) - set(SERVICES)
        if invalid:
            raise typer.BadParameter(
                f"Invalid service(s): {', '.join(invalid)}. Choices: {', '.join(SERVICES)}"
            )
    chosen = [s for s in SERVICES if not services or s in services]
    if not chosen:
        raise typer.Exit(code=1)

    images: list[str] = []
    for name in chosen:
        img = f"localhost/{IMAGES[name]}:{tag}"
        print(f"\n--- Building {name} ({img}) ---")
        run([container_cli(), "build", "-t", img, str(SERVICES[name])])
        images.append(img)

    if images:
        print()
        load_images(images)


@app.command()
def apply(
    overlay: str = typer.Argument("local", help="Overlay to apply (choices: local, prod)"),
) -> None:
    if overlay not in OVERLAYS:
        raise typer.BadParameter(f"Invalid overlay '{overlay}'. Choices: {', '.join(OVERLAYS)}")
    overlay_dir = INFRA / overlay
    if not overlay_dir.is_dir():
        raise typer.Exit(code=1)

    print(f"\n--- Applying overlay '{overlay}' ---")
    kustomize = "kustomize" if _has("kustomize") else None

    if kustomize:
        manifest = subprocess.check_output(
            [kustomize, "build", str(overlay_dir)], text=True
        )
        run(["kubectl", "apply", "-f", "-"], input=manifest, text=True)
    else:
        run(["kubectl", "apply", "-k", str(overlay_dir)])


@app.command()
def delete(
    overlay: str = typer.Argument("local", help="Overlay to delete (choices: local, prod)"),
) -> None:
    if overlay not in OVERLAYS:
        raise typer.BadParameter(f"Invalid overlay '{overlay}'. Choices: {', '.join(OVERLAYS)}")
    overlay_dir = INFRA / overlay
    if not overlay_dir.is_dir():
        raise typer.Exit(code=1)

    print(f"\n--- Deleting overlay '{overlay}' ---")
    kustomize = "kustomize" if _has("kustomize") else None

    if kustomize:
        manifest = subprocess.check_output(
            [kustomize, "build", str(overlay_dir)], text=True
        )
        run(["kubectl", "delete", "-f", "-"], input=manifest, text=True)
    else:
        run(["kubectl", "delete", "-k", str(overlay_dir)])


@app.command(name="port-forward")
def portforward(
    service: str = typer.Argument(..., help="Service to forward. Choices: api, web"),
    local_port: Optional[int] = typer.Option(None, "--local-port", help="Local port (default: 8080 for api, 3000 for web)"),
    remote_port: Optional[int] = typer.Option(None, "--remote-port", help="Service port (default: same as service default)"),
) -> None:
    if service not in ("api", "web"):
        raise typer.BadParameter(f"Invalid service '{service}'. Choices: api, web")
    if service == "bot":
        raise typer.BadParameter("bot has no service port to forward")

    default = SERVICE_PORTS[service]
    lp = local_port if local_port is not None else default
    rp = remote_port if remote_port is not None else default

    name = f"svc/control-plane-{service}"
    print(f"\n--- Port-forward {name} localhost:{lp} -> :{rp} ---")
    print("  Press Ctrl+C to stop.")
    try:
        run(["kubectl", "port-forward", name, f"{lp}:{rp}"])
    except KeyboardInterrupt:
        print()


def _has(cmd: str) -> bool:
    try:
        run_quiet([cmd, "version"])
        return True
    except Exception:
        return False


if __name__ == "__main__":
    app()
