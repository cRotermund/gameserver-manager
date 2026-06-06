# Local Development Quickstart — Podman + Minikube on Windows

This is a concrete, step-by-step guide to getting a working local development environment on **Windows with Podman and Minikube**. It documents the exact setup known to work for maintainers.

If you use a different container runtime or Kubernetes emulator, see [CONTRIBUTING.md](../CONTRIBUTING.md) for general guidance and alternative options.

## Prerequisites

### 1. Install Chocolatey

Chocolatey is a package manager for Windows. Open an **administrator** PowerShell terminal and run:

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force
[System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072
Invoke-Expression ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

Close and reopen the terminal, then verify:

```powershell
choco --version
```

Full docs: [chocolatey.org/install](https://chocolatey.org/install)

### 2. Install Podman Desktop (Manual)

Download and install from [podman-desktop.io](https://podman-desktop.io/). Podman Desktop bundles the `podman` CLI and a management GUI.

### 3. Install Remaining Tools via Chocolatey

Run the following from an **administrator** PowerShell terminal:

```powershell
# Kubernetes
choco install minikube -y
choco install kubernetes-cli -y

# Docker CLI (routes through Podman — see next step)
choco install docker-cli -y

# Languages and runtimes
choco install golang -y
choco install nodejs -y
choco install python -y
```

`choco install` prompts for confirmation on each package; `-y` skips the prompts.

**kustomize** is bundled with `kubernetes-cli` (kubectl >= 1.14). If your version doesn't include it, run `choco install kustomize -y`.

### 4. Alias Docker CLI to Podman

The Docker CLI is installed but no Docker daemon is running. Point it at Podman's Docker-compatible socket so that any tool expecting `docker` commands works transparently.

```powershell
# Set DOCKER_HOST to Podman's named pipe
setx DOCKER_HOST "npipe:////./pipe/podman-machine-default"
```

Replace `podman-machine-default` with your Podman machine name if different. Run `podman machine list` to see the name.

Close and reopen terminals for `setx` to take effect, then verify:

```powershell
docker version
# Should report Podman's version info, not Docker Engine
```

### 5. Verify

Close and reopen all terminals (Chocolatey and `setx` update PATH and env during install), then:

```powershell
podman version
minikube version
kubectl version --client
kustomize version
go version
node --version
python --version
```

## Initialize Minikube with Podman

```powershell
# Start a Minikube cluster using the Podman driver
minikube start --driver=podman

# Verify the cluster is running
minikube status
kubectl config current-context   # should print "minikube"
kubectl get nodes                # should show one Ready node
```

After the cluster is running, install the **Kubernetes** extension in Podman Desktop. This lets you browse pods, services, and deployments from the Podman Desktop UI instead of switching to a terminal for everything.

## Clone and Build

```powershell
git clone https://github.com/cRotermund/gameserver-manager.git
cd gameserver-manager
pip install -e .
```

Build all service images. The harness auto-detects Podman and Minikube, handles image loading into the cluster automatically.

```powershell
deploy build
```

Expected output:

```
--- Building api (localhost/control-plane-api:latest) ---
>>> podman build -t localhost/control-plane-api:latest src\services\control-plane-api
...
--- Building web (localhost/control-plane-web:latest) ---
>>> podman build -t localhost/control-plane-web:latest src\services\control-plane-web
...
--- Building bot (localhost/control-bot:latest) ---
>>> podman build -t localhost/control-bot:latest src\discord\control-bot
```

## Set Up the Discord Bot (Optional)

If you want the bot to run:

```powershell
Copy-Item infra\k8s\local\bot.env.sample infra\k8s\local\bot.env
# Edit infra\k8s\local\bot.env and replace with your real token
```

If you skip this, the bot pod will crash-loop. That's fine for API+Web development.

## Deploy

```powershell
deploy apply
```

Check that pods are running:

```powershell
kubectl get pods
```

Expected output:

```
NAME                                 READY   STATUS    RESTARTS   AGE
control-plane-api-xxxxxxxxxx-xxxxx   1/1     Running   0          30s
control-plane-web-xxxxxxxxxx-xxxxx   1/1     Running   0          30s
```

(The bot may show `CrashLoopBackOff` if you skipped the token step.)

## Access the Services

In separate terminals:

```powershell
# Terminal 1 — API
deploy port-forward api
# API available at http://localhost:8080

# Terminal 2 — Web UI
deploy port-forward web
# Web UI available at http://localhost:3000
```

Verify the API is responding:

```powershell
curl http://localhost:8080/health
# {"status":"healthy"}
```

## True Local Development (No Containers)

For debugging and stepping through code without Kubernetes:

```powershell
# Terminal 1 — API
cd src\services\control-plane-api
go run .

# Terminal 2 — Web UI
cd src\services\control-plane-web
npm install
npm run dev

# Terminal 3 — Bot (optional)
cd src\discord\control-bot
npm install
$env:DISCORD_TOKEN = "your-token"
npm run dev
```

## Tear Down

```powershell
# Remove deployments from the cluster
deploy delete

# Stop Minikube (preserves state for next session)
minikube stop

# Delete the cluster entirely (fresh start next time)
minikube delete
```

## Common Issues

### `minikube start` fails: "driver 'podman' not found"

Minikube may not detect Podman automatically. Ensure Podman is running first:

```powershell
podman machine start
minikube start --driver=podman
```

### Podman builds fail with permission errors

On Windows, Podman runs in a Linux VM. If you see filesystem permission errors, ensure the repo is cloned inside WSL2 or use the Podman volume mount instead of host bind mounts.

```powershell
# If using WSL2, clone the repo inside WSL:
wsl
cd ~
git clone https://github.com/cRotermund/gameserver-manager.git
cd gameserver-manager
# Then run podman and minikube from WSL as well
```

### `deploy build` fails: "Neither podman nor docker"

The harness can't find a container runtime on your PATH. Ensure Podman Desktop is running and `podman` is available:

```powershell
podman version
```

If Podman Desktop is installed but `podman` isn't on PATH, restart your terminal or add Podman's installation directory to your PATH manually.

### Minikube can't pull images: ImagePullBackOff

The harness should have loaded images into Minikube automatically. If you still see this error, rebuild and ensure the cluster is the active context:

```powershell
kubectl config use-context minikube
deploy build
deploy apply
```

### Port already in use

If port 8080 or 3000 is occupied, use a custom local port:

```powershell
deploy port-forward api --local-port 9090
# Then access the API at http://localhost:9090
```

### Podman machine takes too much disk space

Podman's Linux VM can grow over time. Reclaim space:

```powershell
podman machine stop
podman system prune --all --volumes
podman machine start
```

### Minikube on WSL2 with Podman

This combination works but requires Minikube to be installed inside WSL2 (not on the Windows host). Install Minikube inside your WSL distro and run everything from the WSL terminal.
