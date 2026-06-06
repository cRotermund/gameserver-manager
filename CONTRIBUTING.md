# Contributing

## Technology Rules

Only these technologies are permitted without prior discussion:

- **Go** — backend services
- **TypeScript** — frontend and potentially the Discord bot
- **AWS** — compute and infrastructure
- **Docker** — all services must be containerized
- **Kubernetes** — deployment target
- **Kustomize** — Kubernetes configuration management

Anything outside this list requires approval via an issue or discussion. This keeps the stack small, maintainable, and cheap to run.

## Getting Started

> Full local setup is in progress. For now, explore the [API spec](src/services/control-plane-api/openapi.yaml) and [architecture doc](ARCHITECTURE.md) to understand the design.

As services are built, each will include a local development guide in its directory. For example:

```
src/services/control-plane-api/
├── openapi.yaml
├── acceptance-criteria.md
└── README.md          # ← setup, build, run instructions (coming)
```

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
4. The CI workflow (Docker build) must pass
5. At least one approving review is required before merge

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
infra/
└── k8s/
    ├── base/                   # Common Kustomize base
    ├── local/                  # Local dev overlays
    └── prod/                   # Production overlays
tests/                          # Integration / E2E tests
```

New services go under `src/services/<name>/`. Shared libraries can live in `src/pkg/` or `src/shared/`.

## Finding Work

Check the [issues tab](https://github.com/cRotermund/gameserver-manager/issues) for open work. Issues tagged `good first issue` are ideal for new contributors. If I want to tackle something not listed, open an issue to discuss it first.

## Questions?

Open a discussion or comment on a relevant issue. We're building this together — no question is too small.
