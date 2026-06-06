#!/usr/bin/env python3
"""Build and deploy game-server-manager services to a local Kubernetes cluster."""

from __future__ import annotations

import argparse
import os
import shlex
import subprocess
import sys
import tempfile
from pathlib import Path

ROOT = Path(__file__).resolve().parent
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


def detect_cluster() -> str | None:
    """Return the local cluster type or None if none is detected."""
    for name, ctx in [("minikube", "minikube"), ("kind", "kind-"), ("rancher-desktop", "rancher-desktop")]:
        try:
            kubectl_current = run_quiet(["kubectl", "config", "current-context"])
            if name == "kind" and kubectl_current.startswith(ctx):
                return "kind"
            if kubectl_current == ctx:
                return name
        except Exception:
            pass

    try:
        run_quiet(["kubectl", "config", "current-context"])
        return "kubectl"
    except Exception:
        return None


def load_images(images: list[str]) -> None:
    """Push images into the local cluster cache if needed."""
    cluster = detect_cluster()
    cli = container_cli()

    if cluster == "minikube":
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

    elif cluster == "kind":
        for img in images:
            run(["kind", "load", "docker-image", img])

    # Docker Desktop and Rancher Desktop share the host Docker daemon —
    # images are already available.
    elif cluster in (None, "kubectl", "rancher-desktop"):
        pass

    else:
        print(f"  (skipped image load — unknown cluster {cluster!r})")


# ---------------------------------------------------------------------------
# commands
# ---------------------------------------------------------------------------


def cmd_build(services: list[str], tag: str) -> None:
    chosen = [s for s in SERVICES if not services or s in services]
    if not chosen:
        sys.exit("No services matched — choices: " + ", ".join(SERVICES))

    images: list[str] = []
    for name in chosen:
        img = f"localhost/{IMAGES[name]}:{tag}"
        print(f"\n--- Building {name} ({img}) ---")
        run([container_cli(), "build", "-t", img, str(SERVICES[name])])
        images.append(img)

    if images:
        print()
        load_images(images)


def cmd_apply(overlay: str) -> None:
    overlay_dir = INFRA / overlay
    if not overlay_dir.is_dir():
        sys.exit(f"overlay not found: {overlay_dir}")

    print(f"\n--- Applying overlay '{overlay}' ---")
    # Use a temp kustomize build + apply to support older kubectl
    # implementations that don't have -k pointing at a directory.
    kustomize = "kustomize" if _has("kustomize") else None

    if kustomize:
        manifest = subprocess.check_output(
            [kustomize, "build", str(overlay_dir)], text=True
        )
        run(["kubectl", "apply", "-f", "-"], input=manifest, text=True)
    else:
        run(["kubectl", "apply", "-k", str(overlay_dir)])


def cmd_delete(overlay: str) -> None:
    overlay_dir = INFRA / overlay
    if not overlay_dir.is_dir():
        sys.exit(f"overlay not found: {overlay_dir}")

    print(f"\n--- Deleting overlay '{overlay}' ---")
    kustomize = "kustomize" if _has("kustomize") else None

    if kustomize:
        manifest = subprocess.check_output(
            [kustomize, "build", str(overlay_dir)], text=True
        )
        run(["kubectl", "delete", "-f", "-"], input=manifest, text=True)
    else:
        run(["kubectl", "delete", "-k", str(overlay_dir)])


def cmd_portforward(service: str, local_port: int | None, remote_port: int | None) -> None:
    if service == "bot":
        sys.exit("bot has no service port to forward")

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


# ---------------------------------------------------------------------------
# cli
# ---------------------------------------------------------------------------


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    sub = parser.add_subparsers(dest="command", required=True)

    # build
    bp = sub.add_parser("build", help="Build Docker images")
    bp.add_argument(
        "services",
        nargs="*",
        choices=list(SERVICES),
        help="Services to build (default: all)",
    )
    bp.add_argument(
        "--tag", default="latest", help="Image tag (default: latest)"
    )

    # apply
    ap = sub.add_parser("apply", help="Apply a Kustomize overlay")
    ap.add_argument(
        "overlay",
        choices=OVERLAYS,
        default="local",
        nargs="?",
        help="Overlay to apply (default: local)",
    )

    # delete
    dp = sub.add_parser("delete", help="Delete a Kustomize overlay")
    dp.add_argument(
        "overlay",
        choices=OVERLAYS,
        default="local",
        nargs="?",
        help="Overlay to delete (default: local)",
    )

    # port-forward
    pp = sub.add_parser("port-forward", help="Forward a service port locally")
    pp.add_argument("service", choices=SERVICES, help="Service to forward")
    pp.add_argument(
        "--local-port",
        type=int,
        default=None,
        help="Local port (default: service default — 8080 for api, 3000 for web)",
    )
    pp.add_argument(
        "--remote-port",
        type=int,
        default=None,
        help="Service port (default: service default)",
    )

    args = parser.parse_args()

    if args.command == "build":
        cmd_build(args.services, args.tag)
    elif args.command == "apply":
        cmd_apply(args.overlay)
    elif args.command == "delete":
        cmd_delete(args.overlay)
    elif args.command == "port-forward":
        cmd_portforward(args.service, args.local_port, args.remote_port)


if __name__ == "__main__":
    main()
