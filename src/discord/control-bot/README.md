# Discord Control Bot

Chat-based interface for managing game servers via Discord.

## Running locally

```sh
npm install
DISCORD_TOKEN=<your-token> npm run dev
```

## Docker

```sh
docker build -t control-bot .
docker run -e DISCORD_TOKEN=<your-token> control-bot
```
