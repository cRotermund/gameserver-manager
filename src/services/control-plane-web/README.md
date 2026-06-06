# Control Plane Web

Browser-based dashboard for the Game Server Manager. Consumes the REST API.

## Running locally

```sh
npm install
npm run dev   # starts on :3000
```

## Building

```sh
npm run build
```

## Docker

```sh
docker build -t control-plane-web .
docker run -p 3000:3000 control-plane-web
```
