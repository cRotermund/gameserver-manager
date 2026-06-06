# Game Server Manager

Manage my game server infrastructure in AWS without touching the console — from a web UI or directly in Discord.

## What It Does

Game Server Manager gives me and my friends a simple way to control game servers running on AWS. Start, stop, reboot, and monitor servers without granting anyone AWS console access.

| Channel          | What I can do                                             |
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

| Layer        | Technology                        |
| ------------ | --------------------------------- |
| Backend API  | Go (Dockerized)                   |
| Web Frontend | TypeScript                        |
| Discord Bot  | Go or TypeScript (TBD)            |
| Auth         | OAuth2 (client credentials)       |
| Observability| Prometheus, Grafana, OpenTelemetry|
| Infra        | AWS (EC2/ECS), Kubernetes, Kustomize |
| CI/CD        | GitHub Actions, GHCR              |

## Getting Started

> This project is in early active development. Local setup instructions are coming soon.

Once the first components land, I'll find detailed setup instructions in each service directory under `src/services/`.

In the meantime, I can explore the [API specification](src/services/control-plane-api/openapi.yaml) and [acceptance criteria](src/services/control-plane-api/acceptance-criteria.md).

## Contributing

I welcome contributions! See [CONTRIBUTING.md](CONTRIBUTING.md) for setup, conventions, and workflow. Check the [architecture doc](ARCHITECTURE.md) for a deeper dive into how the system is designed.

## License

[MIT](https://choosealicense.com/licenses/mit/)
