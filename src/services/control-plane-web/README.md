# Control Plane Web

Browser-based dashboard for the Game Server Manager. Consumes the REST API.

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
