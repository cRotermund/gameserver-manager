# Discord Control Bot

Chat-based interface for managing game servers via Discord.

## Quick start

```sh
npm install
DISCORD_TOKEN=<your-token> npm run dev
```

Create a token at the [Discord Developer Portal](https://discord.com/developers/applications). See [CONTRIBUTING.md](../../../CONTRIBUTING.md) for full development environment setup, including the container-based workflow with Kubernetes.

## Docker

```sh
docker build -t control-bot .
docker run -e DISCORD_TOKEN=<your-token> control-bot
```
