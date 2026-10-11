# Architecture

## Current System

The repository currently contains three services and a Go deployment harness:

1. **Control-plane API** - Go HTTP service exposing the initial AWS EC2 server-control routes.
2. **Web service** - TypeScript and Express service exposing `/health`.
3. **Discord bot** - Go service that responds to messages mentioning the bot.
4. **Harness** - Go CLI for building images and applying Kustomize overlays.

The intended long-term interaction model is:

```
Web client or Discord bot --> Control-plane API --> AWS EC2 game servers
```

The web service and Discord bot do not yet call the API. They are foundations
for those integrations.

## Control-Plane API

Location: `src/services/control-plane-api`

The API currently registers these routes:

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | Health response |
| `GET` | `/v1/servers` | List managed servers |
| `GET` | `/v1/servers/{serverId}` | Get server details |
| `POST` | `/v1/servers/{serverId}/start` | Start a server |
| `POST` | `/v1/servers/{serverId}/stop` | Stop a server |
| `POST` | `/v1/servers/{serverId}/reboot` | Reboot a server |

The API loads its port, AWS game-server tag names, and AWS credentials from
the environment. See the service [README](../src/services/control-plane-api/README.md)
and [OpenAPI specification](../src/services/control-plane-api/openapi.yaml).

OAuth2, scoped authorization, async operation resources, conditional requests,
telemetry, process inspection, and logs are planned but are not implemented by
the current HTTP server.

## Web Service

Location: `src/services/control-plane-web`

The current TypeScript service is an Express process with one route:

```text
GET /health
```

A browser dashboard and API client are planned. The service currently does
not provide server management, telemetry, authentication, or an API proxy.

## Discord Bot

Location: `src/discord/control-bot`

The current Go bot:

- Connects using `DISCORD_TOKEN`.
- Ignores messages authored by itself.
- Responds when a user mentions the bot.
- Does not call the control-plane API.

Discord server commands, authorization, and natural-language parsing are
planned work.

## Deployment

All service images have Dockerfiles. The Go harness builds local images named
`localhost/<service>:<tag>` and applies the local or production Kustomize
overlay.

The local overlay deploys the API, web service, and bot. The production overlay
currently deploys the API and bot; it intentionally removes the web Deployment
and Service. Production secrets are sourced from AWS Systems Manager through
External Secrets Operator.

The GitHub Actions workflow currently publishes only the API image to GHCR. It
does not deploy to Kubernetes or synchronize an Argo CD application. Publishing
the web and bot images and automating deployment remain follow-up work.

## Observability And Authentication

Prometheus, Grafana, OpenTelemetry, OAuth2, rate limiting, and audit logging
are target capabilities, not currently deployed components in this repository.

## Design Direction

The intended architecture remains a REST API consumed by a web client and a
Discord bot, with Kubernetes and Kustomize for deployment. Planned work should
keep the API as the single source of truth for server operations and should
avoid adding direct AWS control paths to client-facing services.

## Open Questions

- What is the final game-server inventory and registration model?
- Which Discord commands and permission model should be exposed?
- Which OAuth2 provider should issue client credentials?
- Should observability be self-hosted or managed?
- Should server log streaming use SSE, Loki, or CloudWatch?
- What cost guardrails should stop idle servers?
