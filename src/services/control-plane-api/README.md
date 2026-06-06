# Control Plane REST API

The backbone of the Game Server Manager. All server state mutations and queries flow through this API.

## Quick start

```sh
go run .                              # listens on :8080
```

See [CONTRIBUTING.md](../../../CONTRIBUTING.md) for full development environment setup, including the container-based workflow with Kubernetes.

## API

Full OpenAPI 3.1 spec: [`openapi.yaml`](./openapi.yaml)

Acceptance criteria: [`acceptance-criteria.md`](./acceptance-criteria.md)

## Building

```sh
CGO_ENABLED=0 go build -o server .
```

## Docker

```sh
docker build -t control-plane-api .
docker run -p 8080:8080 control-plane-api
```
