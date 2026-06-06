# Game Server Manager

Manage game server infrastructure in AWS without touching the console — from a web UI or directly in Discord.

## What It Does

Game Server Manager provides a simple way to control game servers running on AWS. Start, stop, reboot, and monitor servers without granting anyone AWS console access.

| Channel          | What you can do                                           |
| ---------------- | --------------------------------------------------------- |
| **Web App**      | Full dashboard — server list, start/stop, telemetry, logs |
| **Discord Bot**  | Quick commands from chat — `!start`, `!stop`, `!status`   |
| **REST API**     | Programmatic control for automation and integrations      |

## Architecture at a Glance

```
 Discord Bot ──┐
               ├──► REST API ──► AWS (EC2 / ECS) Game Servers
  Web App ─────┘        │
                        ├──► Prometheus / Grafana (metrics)
                        └──► OpenTelemetry (traces)
```

All components run as containers on Kubernetes, defined with Kustomize.

## Features

- **Start / Stop / Reboot** game servers on demand
- **Real-time observability** — CPU, memory, disk, network, connected clients
- **Per-server process inspection** — see what's running
- **Async operations** with polling — know when a start/stop completes
- **OAuth2 authentication** — scoped access (read-only or read-write)
- **Rate limiting** — protects against runaway clients or abuse
- **Audit logging** — every control action is tied to a user

### On the Roadmap

- Telemetry streaming (SSE) for live dashboards
- Log retrieval and tail streaming
- Server registration and de-registration via API
- Discord bot with natural-language command parsing

## Tech Stack

| Layer         | Technology                        |
| ------------- | --------------------------------- |
| Backend API   | Go (Dockerized)                   |
| Web Frontend  | TypeScript                        |
| Discord Bot   | TypeScript                        |
| Auth          | OAuth2 (client credentials)       |
| Observability | Prometheus, Grafana, OpenTelemetry|
| Infra         | AWS (EC2/ECS), Kubernetes, Kustomize |

## Getting Started

### Prerequisites

| Tool       | Version    | Why                             |
| ---------- | ---------- | ------------------------------- |
| Go         | >= 1.24    | API service                     |
| Node.js    | >= 22      | Web frontend & Discord bot      |
| Python     | >= 3.10    | `deploy` harness (K8s only)  |
| Docker     | latest     | Container builds                |
| kubectl    | latest     | K8s deployments                 |
| kustomize  | latest     | K8s configuration               |

Kubernetes is only required for the container-based workflow. For true local development (stepping through code in an IDE), you only need Go and Node.js.

### Two Development Workflows

Pick the one that fits your style.

---

**Path A — Run services directly on your machine**

Ideal for debugging, stepping through code, and fast iteration loops. No containers, no Kubernetes.

```sh
# Terminal 1 — API (Go)
cd src/services/control-plane-api
go run .                          # listens on :8080

# Terminal 2 — Web frontend (TypeScript)
cd src/services/control-plane-web
npm install
npm run dev                       # listens on :3000

# Terminal 3 — Discord bot (optional, needs a token)
cd src/discord/control-bot
npm install
DISCORD_TOKEN=<your-token> npm run dev
```

Open `http://localhost:3000` for the web UI, or point your HTTP client at `http://localhost:8080` for the API.

---

**Path B — Run in containers with Kubernetes**

Matches the production environment. Uses the `deploy` CLI as a one-stop harness.

```sh
# 1. Start your local cluster (minikube, kind, or Rancher Desktop)

# 2. Install Python dependencies
pip install -e .

# 3. Build all images
deploy build

# 4. Deploy to the cluster
deploy apply

# 5. Port-forward the services
deploy port-forward api     # API at localhost:8080
deploy port-forward web     # Web at localhost:3000
```

The harness auto-detects your container CLI (podman or docker) and cluster type. For the Discord bot, copy `infra/k8s/local/bot.env.sample` to `bot.env` with your token before deploying.

---

Full development environment details, troubleshooting, and environmental quirks are covered in [CONTRIBUTING.md](CONTRIBUTING.md).

## Contributing

Contributions welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for setup, conventions, and workflow. The [architecture doc](docs/architecture.md) covers system design in depth.

## License

[MIT](https://choosealicense.com/licenses/mit/)
