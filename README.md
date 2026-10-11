# Game Server Manager

Game Server Manager is a Go control-plane API, a TypeScript web health service,
and a Go Discord bot for managing AWS game-server infrastructure. The current
repository provides the API's initial server-control endpoints and a minimal
Discord hello-world bot; the broader dashboard and command experience remains
under development.

## Current Components

| Component | Current behavior |
| --- | --- |
| **REST API** | Lists servers and starts, stops, or reboots them through AWS EC2. |
| **Web service** | Serves a health endpoint on port `3000`. |
| **Discord bot** | Replies to messages that mention it with a hello-world response. |
| **Harness** | Builds images and applies or deletes local and production Kustomize overlays. |

The API routes are currently:

- `GET /health`
- `GET /v1/servers`
- `GET /v1/servers/{serverId}`
- `POST /v1/servers/{serverId}/start`
- `POST /v1/servers/{serverId}/stop`
- `POST /v1/servers/{serverId}/reboot`

## Architecture

```
Discord Bot ──┐
              ├──> REST API ──> AWS EC2 game servers
Web Service ──┘
```

All services are containerized. Kubernetes manifests use Kustomize with a
local overlay and a production overlay. See [docs/architecture.md](docs/architecture.md)
for the current implementation and planned architecture.

## Planned Work

- Browser dashboard for server management and telemetry
- Discord commands backed by the REST API
- OAuth2 authentication and scoped authorization
- Async operation resources and polling
- Metrics, tracing, logs, and audit history
- Additional API endpoints for telemetry, processes, and logs

## Tech Stack

| Layer | Technology |
| --- | --- |
| Backend API | Go |
| Discord bot | Go |
| Web service | TypeScript and Node.js |
| Deployment harness | Go |
| Infrastructure | AWS, Kubernetes, Kustomize |
| Containers | Docker-compatible runtime and distroless images where applicable |

## Getting Started

### Prerequisites

| Tool | Version | Why |
| --- | --- | --- |
| Go | >= 1.26 | API, Discord bot, and deployment harness |
| Node.js | >= 22 | Web service |
| Docker or Podman | Current release | Container builds |
| kubectl | Current release | Kubernetes access |
| kustomize | Current release, optional | Standalone manifest rendering; the harness falls back to `kubectl apply -k` when unavailable |

### Run Services Directly

API:

```sh
cd src/services/control-plane-api
# Set PORT, the three GAMESERVER_* tag names, and AWS credentials first.
go run .
```

Web service:

```sh
cd src/services/control-plane-web
npm install
npm run dev
```

Discord bot:

```sh
cd src/discord/control-bot
DISCORD_TOKEN=<your-token> go run .
```

The bot also loads a `.env` file from its working directory. Never commit a
real token.

### Run With Kubernetes

Start a local Kubernetes cluster, then run these commands from the repository
root:

```sh
# Build all service images and load them when the cluster supports it.
./harness.sh build

# Create local secret files from the examples and fill in real values.
cp infra/k8s/local/api.env.example infra/k8s/local/api.env
cp infra/k8s/local/bot.env.example infra/k8s/local/bot.env

# Apply the local overlay.
./harness.sh apply

# Port-forward the API or web service.
./harness.sh port-forward api
./harness.sh port-forward web
```

The local secret files are gitignored. The Discord bot requires
`DISCORD_TOKEN`; the API also requires its AWS and game-server configuration.
See [CONTRIBUTING.md](CONTRIBUTING.md) and [docs/local-dev.md](docs/local-dev.md)
for setup details.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development conventions and
[docs/architecture.md](docs/architecture.md) for system design.

## License

[MIT](https://choosealicense.com/licenses/mit/)
