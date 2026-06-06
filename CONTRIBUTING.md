# Contributing

## Technology Rules

Only these technologies are permitted without prior discussion:

- **Go** — backend services
- **TypeScript** — frontend and Discord bot
- **AWS** — compute and infrastructure
- **Docker** — all services must be containerized
- **Kubernetes** — deployment target
- **Kustomize** — Kubernetes configuration management

Anything outside this list requires approval via an issue or discussion. This keeps the stack small, maintainable, and cheap to run.

## Getting Started

### Prerequisites

Install these before you start. Versions listed are the minimum required.

| Tool       | Version    | Install Link                                          | Required For        |
| ---------- | ---------- | ----------------------------------------------------- | ------------------- |
| Go         | >= 1.24    | [go.dev/dl](https://go.dev/dl/)                       | API service         |
| Node.js    | >= 22      | [nodejs.org](https://nodejs.org/)                     | Web, Discord bot    |
| Python     | >= 3.10    | [python.org](https://www.python.org/downloads/)       | `deploy` harness + `pip install -e .` |
| Docker     | latest     | See container runtime notes below                     | Container builds    |
| kubectl    | latest     | [kubernetes.io](https://kubernetes.io/docs/tasks/tools/) | K8s deployments  |
| kustomize  | latest     | [kustomize.io](https://kustomize.io/)                 | K8s configuration   |

**If you only need Path A (true local development),** Go and Node.js are enough. Docker, kubectl, kustomize, and Python are only needed for the container-based workflow.

For system design and architecture, see [docs/architecture.md](docs/architecture.md). For a concrete, step-by-step walkthrough of a known-working local setup (Podman + Minikube on Windows), see [docs/local-dev.md](docs/local-dev.md).

Clone the repo:

```sh
git clone https://github.com/cRotermund/gameserver-manager.git
cd gameserver-manager
pip install -e .
```

If `pip` warns that the scripts directory is not on PATH, the `deploy` command won't be found. Add the reported directory to your PATH, or use `python -m harness.cli` as a fallback (all examples below work with either approach).

Verify with:

```sh
deploy --help
```

### Two Development Workflows

---

#### Path A — Run services directly on your machine

Best for debugging, stepping through code in an IDE, and fast iteration. No containers, no Kubernetes.

**API (Go)**

```sh
cd src/services/control-plane-api
go run .                              # listens on :8080
```

The API reads `PORT` from the environment (defaults to `8080`). Set it if you need a different port:

```sh
PORT=9090 go run .
```

**Web frontend (TypeScript)**

```sh
cd src/services/control-plane-web
npm install
npm run dev                           # listens on :3000
```

`npm run dev` uses `tsx` for hot-reload during development. Use `npm run build && npm start` for a production-like run.

**Discord bot (TypeScript)**

```sh
cd src/discord/control-bot
npm install
DISCORD_TOKEN=<your-token> npm run dev
```

The bot requires a Discord application token. Create one at the [Discord Developer Portal](https://discord.com/developers/applications). Without a valid token the bot will exit immediately.

**IDE setup**

Most IDEs can run these commands as launch configurations:

- **VS Code**: Create `.vscode/launch.json` entries pointing at the `go run` / `npm run dev` commands. Set `cwd` to the service directory.
- **GoLand / IntelliJ**: Add a Go Build run configuration with the service directory as the module root.
- **WebStorm**: Add an npm run configuration targeting the `dev` script.

---

#### Path B — Run in containers with Kubernetes

Matches the production environment. All commands are issued from the repository root.

The [`deploy`](#) CLI (see [harness/](harness/)) is the single entry point for container-based development.

**1. Start a local Kubernetes cluster**

Pick one:

| Cluster            | Pros                                                         |
| ------------------ | ------------------------------------------------------------ |
| **Rancher Desktop** | Ships with Kubernetes built in. Easiest setup on Windows/macOS. |
| **minikube**       | Most configurable, good CI compatibility.                    |
| **kind**           | Runs K8s in Docker. Lightweight, fast start.                 |

The harness auto-detects whichever is running. For a full walkthrough of one combination known to work, see [docs/local-dev.md](docs/local-dev.md) (Podman + Minikube on Windows).

Ensure `kubectl` is configured to point at your cluster:

```sh
kubectl config current-context
```

**2. Build images**

```sh
# Build all services
deploy build

# Build a specific service
deploy build api

# Custom tag
deploy build --tag dev
```

Images are tagged as `localhost/<name>:latest` and automatically loaded into your cluster's image cache (no registry needed).

**3. Set up the Discord bot secret (optional)**

```sh
cp infra/k8s/local/bot.env.sample infra/k8s/local/bot.env
# Edit bot.env with your real token
```

**4. Deploy**

```sh
# Apply the local overlay
deploy apply

# Check that pods are running
kubectl get pods
```

The `local` overlay uses `localhost/`-prefixed images with `IfNotPresent` pull policy, so images are taken from your local Docker daemon directly.

**5. Access the services**

```sh
# Port-forward API to localhost:8080
deploy port-forward api

# Port-forward web UI to localhost:3000 (in another terminal)
deploy port-forward web
```

Custom ports:

```sh
deploy port-forward api --local-port 9090
```

**6. Tear down**

```sh
deploy delete
```

---

### Environmental Quirks

#### Container runtime: Docker vs Podman

The harness and Dockerfiles work with both. The harness auto-detects whichever is on your PATH, preferring Podman.

| Issue                                      | Solution                                                     |
| ------------------------------------------ | ------------------------------------------------------------ |
| Docker Desktop requires a license for commercial use | Use Podman or Rancher Desktop (free, open-source alternatives). |
| Podman builds may fail with SELinux errors | Disable SELinux labeling: `podman build --security-opt label=disable` (or configure your machine). |
| Podman cannot load images into kind        | `kind` expects `docker`. Export images as tarballs and load them manually: `podman save <img> | kind load image-archive /dev/stdin`. |
| Docker rate limits on Docker Hub pulls     | Log in to Docker Hub or configure a pull-through mirror.     |

#### Kubernetes cluster choice

| Issue                                      | Solution                                                     |
| ------------------------------------------ | ------------------------------------------------------------ |
| Rancher Desktop networking broken on VPN   | Use minikube with the `--driver=docker` flag instead.        |
| minikube `image load` fails with Podman    | The harness handles this automatically (exports to tar, loads via `minikube image load`). |
| kind naming depends on `kind-<name>` context| The harness detects this. If you use a nonstandard kind name, set `KUBECONFIG` explicitly. |
| WSL2 with Docker Desktop                   | Enable the WSL2 integration in Docker Desktop settings for your distro. |
| Port conflicts (8080, 3000)                | Use `--local-port` on `port-forward`, or set `PORT` for local runs. |

#### Windows-specific notes

- **Line endings:** The `.editorconfig` and `.gitattributes` enforce LF. If you see CRLF warnings, run `git config core.autocrlf input`.
- **Shell:** The harness and service scripts assume a POSIX-like shell (Git Bash, WSL2, or PowerShell). The commands in this guide use the format compatible with Git Bash / WSL2. For PowerShell, replace environment variable syntax (`$env:PORT=9090` instead of `PORT=9090`).
- **Node.js native modules:** Not currently used, but if added later they may require Visual Studio Build Tools on Windows.
- **Symlinks in node_modules:** Some Windows configurations require Developer Mode to be enabled for symlink creation during `npm install`.

### Troubleshooting

| Problem                                     | Likely Cause                                                 | Fix                                                           |
| ------------------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------- |
| `deploy: command not found`                 | Python's scripts directory not on PATH after `pip install`.  | Add the directory from the `pip` warning to PATH, or use `python -m harness.cli` instead. |
| `ErrImagePull` / `ImagePullBackOff`         | Image not available in cluster cache.                        | Run `deploy build` to rebuild into the local daemon. |
| `CrashLoopBackOff` on any pod               | The container exits immediately after starting — could be a code error, missing config, or bad env var. | Check the logs first: `kubectl logs deployment/<name>`. For the bot specifically, a missing `DISCORD_TOKEN` is the most common cause. |
| `kubectl get pods` shows nothing            | Wrong kubeconfig context.                                    | `kubectl config get-contexts` and switch to your local cluster. |
| `go run .` fails with "package not found"   | Not running from the service directory.                      | cd into the service directory first.                          |
| `npm run dev` fails on Windows              | Path too long for node_modules.                              | Enable long paths in Windows (`reg add HKLM\SYSTEM\CurrentControlSet\Control\FileSystem /v LongPathsEnabled /t REG_DWORD /d 1`). |
| "Neither podman nor docker" error from harness | Container runtime not installed or not on PATH.            | Install Docker Desktop, Podman, or Rancher Desktop.           |
| Port already in use                         | Another process or a previous `port-forward` still running.  | Kill the port-forward process or pick a different local port. |

## Development Workflow

### Branches

Branch from `main` using the naming convention:

```
feat/<issue-number>/<short-description>
```

Examples:
- `feat/1/api-design-draft`
- `feat/4/agentic-instrumentation`

### Commits

Use [conventional commits](https://www.conventionalcommits.org/). Every commit must reference an issue:

```
feat(resolves #1): first draft design of api
fix(#12): handle server not found in stop endpoint
docs(#3): add local development guide
```

Keep commits focused. One logical change per commit.

| Prefix      | When to use                                              |
| ----------- | -------------------------------------------------------- |
| `feat`      | New feature or capability                                |
| `fix`       | Bug fix                                                  |
| `docs`      | Documentation only (README, comments, OpenAPI spec)      |
| `chore`     | Maintenance tasks — deps, build scripts, tooling, CI     |
| `refactor`  | Code restructuring that doesn't change behavior          |
| `test`      | Adding or updating tests                                 |

#### Issue Linking

Every commit must include an issue reference — either a plain reference or a closing keyword. Which one to use depends on whether the commit completes the work:

| Syntax              | Effect                                                        |
| ------------------- | ------------------------------------------------------------- |
| `(#12)`             | Links to the issue without closing it. Use when the commit is one step of a multi-commit issue. |
| `(closes #12)`      | Closes the issue automatically when merged to `main`. Use when the commit (or PR) fully resolves the issue. |

Accepted closing keywords are `close`, `closes`, `closed`, `fix`, `fixes`, `fixed`, `resolve`, `resolves`, `resolved`. They all behave identically — pick one and be consistent in the PR.

**When to use each:**

```
# Multi-commit feature — only the final commit closes the issue
feat(#5): scaffold the new telemetry endpoint
feat(#5): add CPU and memory collectors
feat(resolves #5): wire up telemetry SSE stream

# Single-commit fix — closes immediately
fix(closes #8): correct ETag format on server detail response

# Pure refactor spread across commits — plain reference throughout, never closes
refactor(#12): extract AWS client to shared package
refactor(#12): migrate server service to shared client
```

The closing keyword is what matters during merge. If a single PR contains multiple commits, using a closing keyword on any one of them (usually the last) will close the issue. Using plain references on all commits means the issue stays open after merge and must be closed manually.

### Pull Requests

1. Create a feature branch from `main`
2. Make changes, following the conventions below
3. Open a PR against `main` with a clear description of what and why
4. At least one approving review is required before merge

## Code Style

An [`.editorconfig`](.editorconfig) is provided. Editors should respect it automatically.

| Language     | Indent | Tabs/Spaces |
| ------------ | ------ | ----------- |
| Go           | 4      | Tabs        |
| TypeScript   | 4      | Spaces      |
| JSON / YAML  | 4      | Spaces      |
| All files    | LF     | —           |

Additional guidelines:
- **Go:** Follow standard Go conventions. Use `gofmt`. Group imports (stdlib, third-party, internal).
- **TypeScript:** Prefer explicit types. Avoid `any`. Use functional components and hooks for UI.
- **General:** No commented-out code in commits. Trim trailing whitespace. End files with a newline.

## Testing

Write tests for new functionality. Patterns TBD — will be established as the first services are built. In the meantime:

- Go: table-driven tests with the standard `testing` package
- TypeScript: framework TBD (Vitest or Jest likely)
- Integration tests for API endpoints are strongly encouraged

## Directory Structure

```
src/
├── services/
│   ├── control-plane-api/      # Go REST API
│   └── control-plane-web/      # TypeScript web frontend
├── discord/
│   └── control-bot/            # Discord bot
harness/                        # Python CLI for building and deploying
docs/                           # Architecture and dev guides
infra/
└── k8s/
    ├── base/                   # Common Kustomize base
    ├── local/                  # Local dev overlays
    └── prod/                   # Production overlays
tests/                          # Integration / E2E tests
```

New services go under `src/services/<name>/`. Shared libraries can live in `src/pkg/` or `src/shared/`.

## Finding Work

Check the [issues tab](https://github.com/cRotermund/gameserver-manager/issues) for open work. Issues tagged `good first issue` are ideal for new contributors. If you want to tackle something not listed, open an issue to discuss it first.

## Questions?

Open a discussion or comment on a relevant issue. We're building this together — no question is too small.
