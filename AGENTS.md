# AGENTS.md — AI Coding Assistant Instructions

## Project Context

Game Server Manager is a three-tier system for managing AWS game servers without console access. Users interact via a **web app** or **Discord bot**, both of which talk to a central **REST API** that orchestrates AWS resources.

The project is early-stage. Design docs and API specs exist, but most service code has not been written yet. You are helping build this system from the ground up.

## Tech Stack (Non-Negotiable)

| Layer         | Technology                    |
| ------------- | ----------------------------- |
| Backend       | Go (Dockerized, distroless)   |
| Frontend      | TypeScript                    |
| Infra         | AWS (EC2/ECS)                 |
| Deployment    | Kubernetes + Kustomize        |
| CI/CD         | GitHub Actions, GHCR          |
| Auth          | OAuth2 client credentials     |
| Observability | Prometheus, Grafana, oTel     |

**Do not introduce** languages, frameworks, or infrastructure tools outside this list without explicit approval.

## Directory Conventions

```
harness/                 ← Development harness (CLI for building, deploying, port-forwarding)
src/services/<name>/     ← Each service gets its own directory
src/discord/<name>/      ← Discord bot lives here
src/pkg/                 ← Shared Go libraries (create as needed)
infra/k8s/base/          ← Common Kustomize manifests
infra/k8s/local/         ← Local dev overlays
infra/k8s/prod/          ← Production overlays
tests/                   ← Integration and E2E tests
```

When creating a new service, include at minimum: code, a Dockerfile, and a README with local setup instructions.

## Code Style

Reference [`.editorconfig`](.editorconfig) at all times:

- **Go:** tabs, indent 4, LF line endings. Follow standard Go conventions (`gofmt`, grouped imports: stdlib → third-party → internal). Use the standard `testing` package for tests — prefer table-driven tests.
- **TypeScript:** spaces, indent 4, LF line endings. Explicit types, no `any`. Functional components and hooks for UI code.
- **YAML/JSON:** spaces, indent 4, LF.
- **All files:** no trailing whitespace, end with a final newline, no commented-out code.

## Git Conventions

- **Branches:** `feat/<issue-number>/<short-description>` — branch from `main`
- **Commits:** Conventional Commits format. Reference issues when applicable:
  ```
  feat(resolves #1): description
  fix(#12): description
  docs: description
  ```
- **PRs:** One logical change per PR. Must pass CI. Requires review.

## API Design Rules

The REST API is the single source of truth. When extending it:

- Version under `/v1`. Keep the OpenAPI spec in sync (`src/services/control-plane-api/openapi.yaml`).
- Mutations return `202 Accepted` with an operation resource for async tracking (server operations take time).
- Read endpoints use conditional requests (`ETag`, `If-None-Match`, `Cache-Control`).
- Errors follow the `ApiError` schema: `{ "error": "human-readable", "code": "MACHINE_READABLE" }`.
- Scopes are per-endpoint and checked against OAuth2 client credentials.

## What to Do

- Write idiomatic Go and TypeScript that follows existing patterns
- Containerize every service — include a Dockerfile
- Keep infrastructure declarative (Kustomize overlays in `infra/k8s/`)
- Update the OpenAPI spec when API endpoints change
- Write tests alongside implementation code
- Respect the acceptance criteria in `src/services/control-plane-api/acceptance-criteria.md`

## What NOT to Do

- Do not add new languages or frameworks (no Python, Rust, React alternatives, Helm, Terraform, etc.) without asking
- Do not commit secrets, `.env` files, or IDE configs (see [`.gitignore`](.gitignore))
- Do not leave commented-out code in commits
- Do not write blocking HTTP calls for long-running operations — always use the async operation pattern
- Do not introduce infrastructure changes outside the Kustomize workflow
- Do not add dependencies without a clear justification tied to the project's goals
