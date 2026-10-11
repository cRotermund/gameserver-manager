# Control Plane Web

TypeScript and Express web service for the Game Server Manager. The current
implementation exposes a health endpoint; the browser dashboard is planned.

## Quick start

```sh
npm install
npm run dev                           # listens on :3000
```

See [CONTRIBUTING.md](../../../CONTRIBUTING.md) for full development environment setup, including the container-based workflow with Kubernetes.

## Building

```sh
npm run build
```

## Docker

```sh
docker build -t control-plane-web .
docker run -p 3000:3000 control-plane-web
```

The health endpoint is available at `GET /health`.
