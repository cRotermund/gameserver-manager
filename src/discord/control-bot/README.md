# Discord Control Bot

Minimal Go Discord bot for the Game Server Manager. The current implementation
responds to messages that mention the bot with a hello-world message.

## Quick Start

Set `DISCORD_TOKEN` and run the Go module:

```sh
DISCORD_TOKEN=<your-token> go run .
```

The bot also loads a `.env` file from its working directory. Create a token at
the [Discord Developer Portal](https://discord.com/developers/applications).
Never commit a real token.

## Build

```sh
go build -o control-bot .
```

## Docker

```sh
docker build -t control-bot .
docker run -e DISCORD_TOKEN=<your-token> control-bot
```

The Dockerfile produces a static Linux binary in a non-root distroless image.
