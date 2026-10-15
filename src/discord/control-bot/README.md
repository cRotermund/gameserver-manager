# Discord Control Bot

The Discord control bot lets users query and change the state of game servers
through the Game Server Manager control-plane API. It responds to messages that
mention the bot and use the `!` command prefix.

## Commands

| Command | Description |
| --- | --- |
| `!help` | Show the available commands. |
| `!listservers` | List the available game servers. |
| `!serverstatus <name>` | Show a server's current status. |
| `!startserver <name>` | Request that a server start, then poll for the target state. |
| `!stopserver <name>` | Request that a server stop, then poll for the target state. |

Example:

```text
@GameServerBot !serverstatus valheim
```

Server names are currently parsed as a single argument, so names containing
spaces are not supported.

## Configuration

The bot reads configuration from the process environment. It also attempts to
load a `.env` file from its working directory for local development.

| Variable | Required | Description |
| --- | --- | --- |
| `DISCORD_TOKEN` | Yes | Discord bot token. |
| `GSM_API_HOST` | Yes | Hostname of the control-plane API. |
| `GSM_API_PORT` | Yes | Port of the control-plane API. |
| `GSM_API_TIMEOUT_MS` | No | Control-plane request timeout in milliseconds. Defaults to `5000`. |

Create the bot token in the [Discord Developer Portal](https://discord.com/developers/applications).
Never commit a real token or `.env` file.

## Local Development

From this directory, configure the required variables and run the module:

```sh
go run .
```

The bot requires a reachable control-plane API and a Discord application with
the required gateway intents enabled.

## Build

From this directory:

```sh
go build -o control-bot .
```

Run the checks with:

```sh
go test ./...
go vet ./...
```

## Docker

Build from the repository root because the Dockerfile copies both the bot and
the shared control-plane client package:

```sh
docker build -f src/discord/control-bot/Dockerfile -t control-bot .
```

Run the container with the required configuration:

```sh
docker run --rm \
  -e DISCORD_TOKEN=<your-token> \
  -e GSM_API_HOST=<control-plane-host> \
  -e GSM_API_PORT=<control-plane-port> \
  control-bot
```

The Dockerfile produces a static Linux binary in a non-root distroless image.
