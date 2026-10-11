# Control Plane REST API

The Go control-plane API for querying and changing managed AWS game servers.

## Quick start

```sh
go run .                              # listens on :8080
```

See [CONTRIBUTING.md](../../../CONTRIBUTING.md) for full development environment setup, including the container-based workflow with Kubernetes.

## API

Full OpenAPI 3.1 spec: [`openapi.yaml`](./openapi.yaml)

Acceptance criteria: [`docs/api-acceptance-criteria.md`](../../../docs/api-acceptance-criteria.md)

Current routes are `GET /health`, `GET /v1/servers`,
`GET /v1/servers/{serverId}`, and `POST` start, stop, and reboot routes under
`/v1/servers/{serverId}`.

## Building

```sh
CGO_ENABLED=0 go build -o server .
```

## Docker

```sh
docker build -t control-plane-api .
docker run --env-file <path-to-api.env> -p 8080:8080 control-plane-api
```

The container requires `PORT`, `GAMESERVER_ISSERVER_TAG_KEY`,
`GAMESERVER_NAME_TAG_KEY`, and `GAMESERVER_DESCRIPTOR_TAG_KEY`. It also needs
AWS SDK credentials and permissions to describe and control the target EC2
instances. Kubernetes supplies these values through the Deployment and Secret.
