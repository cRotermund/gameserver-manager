# Architecture

## Overview

Game Server Manager follows a three-tier architecture:

1. **Control Plane REST API** — the single source of truth for server state and operations
2. **Web Application** — a browser-based dashboard for management and observability
3. **Discord Bot** — a chat-based interface for quick, convenient control

```
┌─────────────────────────────────────────────────────┐
│                     Users / Friends                   │
└──────────┬──────────────────────┬────────────────────┘
           │                      │
     ┌─────▼─────┐          ┌─────▼─────┐
     │  Web App   │          │ Discord   │
     │ (Browser)  │          │  Bot      │
     └─────┬──────┘          └─────┬─────┘
           │                      │
           └──────────┬───────────┘
                      │ HTTPS (OAuth2)
              ┌───────▼────────┐
              │   REST API     │
              │   (Go, K8s)    │
              └───┬────────────┘
                  │
      ┌───────────┼───────────┐
      │           │           │
 ┌────▼────┐ ┌────▼────┐ ┌────▼────┐
 │  AWS    │ │ Prom/Graf│ │  oTel   │
 │EC2/ECS  │ │  ana     │ │Collector│
 └─────────┘ └─────────┘ └─────────┘
```

## Primary Motivation

Extremely cheap to run — ideally free-tier or near-zero cost — on top of pre-existing Kubernetes infrastructure.

## Subsystems

### REST API (`src/services/control-plane-api`)

The backbone of the system. All state mutations and queries flow through the API.

- **Language:** Go, compiled to a static binary, running in a distroless Docker container
- **Auth:** OAuth2 client credentials with scoped access (`control:ro`, `control:rw`, `observability:ro`)
- **API style:** RESTful, versioned under `/v1`, documented with OpenAPI 3.1
- **Async operations:** Start/stop/reboot return `202 Accepted` with an operation resource that can be polled via `/operations/{id}`
- **Caching:** Conditional requests with `ETag` and `Cache-Control` — cheap endpoints cached longer, in-progress operations cached briefly

#### Key Endpoints

| Method | Path                           | Scope            | Description                    |
| ------ | ------------------------------ | ---------------- | ------------------------------ |
| GET    | `/v1/servers`                  | `control:ro`     | List all managed servers       |
| GET    | `/v1/servers/{id}`             | `control:ro`     | Server detail + resource usage |
| GET    | `/v1/servers/{id}/processes`   | `control:ro`     | Running process table          |
| POST   | `/v1/servers/{id}/start`       | `control:rw`     | Start a server (async)         |
| POST   | `/v1/servers/{id}/stop`        | `control:rw`     | Stop a server (async)          |
| POST   | `/v1/servers/{id}/reboot`      | `control:rw`     | Quick OS reboot (async)        |
| GET    | `/v1/operations/{id}`          | `control:ro`     | Poll operation status          |

Full spec: [`src/services/control-plane-api/openapi.yaml`](src/services/control-plane-api/openapi.yaml)

### Web Application (`src/services/control-plane-web`)

A browser-based UI built in TypeScript that consumes the REST API.

- **Server management:** List servers, view status, issue start/stop/reboot
- **Observability dashboard:** CPU, memory, disk, network, connected clients
- **API health monitor:** Simple availability indicator for the control plane itself
- **Auth:** Handles OAuth2 flow and token management transparently for the user

### Discord Bot (`src/discord/control-bot`)

A lightweight chat interface for Discord servers.

- **Commands:** `!start`, `!stop`, `!status` (exact set TBD)
- **Auth:** Uses a service account with OAuth2 client credentials
- **Future:** In-memory semantic analysis for natural-language command interpretation (experimental)

### Observability

- **OpenTelemetry:** Distributed tracing across API, web, and bot
- **Prometheus:** Scrapes metrics endpoints from all services
- **Grafana:** Dashboards for system health, server telemetry, and operational metrics

## Data Flow

### Control Operations (Start / Stop / Reboot)

```
Client ──POST /v1/servers/{id}/start──► API
                                          │
                                          ├── Validate auth & state
                                          ├── Create Operation record (status: pending)
                                          ├── Return 202 + Location: /v1/operations/{id}
                                          │
                                          ├── [Async] Issue AWS API call
                                          ├── [Async] Poll AWS for completion
                                          ├── [Async] Update Operation (in_progress → completed)
                                          └── [Async] Update Server status (starting → running)
```

### Read Operations (Status / Telemetry)

```
Client ──GET /v1/servers/{id}──► API
                                  │
                                  ├── Check If-None-Match / ETag
                                  ├── Return 304 if unchanged
                                  └── Return 200 + fresh data + new ETag + Cache-Control
```

## Deployment

```
GitHub Push ──► GitHub Actions ──► Build Docker image ──► Push to GHCR
                                                              │
                                                              ▼
                                              Kubernetes (Kustomize)
                                              ├── base/       (common config)
                                              ├── local/      (dev overlays)
                                              └── prod/       (production overlays)
```

All services are containerized and deployed to Kubernetes. Kustomize overlays handle environment-specific differences (credentials, replica counts, resource limits). No manual infrastructure changes — everything is Infrastructure as Code.

## Authentication & Authorization

- **Scheme:** OAuth2 client credentials flow
- **Token endpoint:** served by the API itself (`/oauth/token`)
- **Scopes:**
  - `control:ro` — read server state and operations
  - `control:rw` — full control (start, stop, reboot)
  - `observability:ro` — read telemetry and logs
- **Granularity:** Per-endpoint. An endpoint can accept multiple scopes (e.g., GET servers accepts both `control:ro` and `control:rw`).

## Design Decisions

| Decision                     | Rationale                                                                 |
| ---------------------------- | ------------------------------------------------------------------------- |
| REST over GraphQL/gRPC       | Broader client compatibility, easier caching, simpler for a small API     |
| Async operations (202 + poll)| Server starts/stops take minutes; blocking HTTP calls are impractical     |
| Conditional requests (ETag)  | Reduces bandwidth and server load for polling clients                     |
| Kustomize over Helm          | Lighter weight, better fits a small project, still declarative            |
| Distroless containers        | Minimal attack surface, smaller images, faster pulls                      |

## Open Questions

- Exact game server inventory model — register via API, IaC, or both?
- Discord bot — Go or TypeScript? Shared auth client available in both?
- Observability pipeline — self-hosted Prometheus/Grafana or managed (Grafana Cloud, AWS Managed Prometheus)?
- Server log streaming — push (SSE from API) or pull (ship to Loki/CloudWatch)?
- Cost guardrails — automatic shutdown after idle period? Budget alerts?
