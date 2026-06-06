# Control Plane REST API

The backbone of the Game Server Manager. All server state mutations and queries flow through this API.

## Running locally

```sh
go run .   # starts on :8080
```

## Building

```sh
CGO_ENABLED=0 go build -o server .
```

## Docker

```sh
docker build -t control-plane-api .
docker run -p 8080:8080 control-plane-api
```

## API

Full OpenAPI spec: [`openapi.yaml`](./openapi.yaml)
