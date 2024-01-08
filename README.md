[![Docker Build](https://github.com/rickstaa/warnings-discord-bot/actions/workflows/docker-build.yml/badge.svg)](https://github.com/rickstaa/warnings-discord-bot/actions/workflows/docker-build.yml)
[![Publish Docker image](https://github.com/rickstaa/warnings-discord-bot/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/rickstaa/warnings-discord-bot/actions/workflows/docker-publish.yml)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/rickstaa/warnings-discord-bot?logo=docker)
](https://hub.docker.com/r/rickstaa/warnings-discord-bot)
[![Latest Release](https://img.shields.io/github/v/release/rickstaa/warnings-discord-bot?label=latest%20release)](https://github.com/rickstaa/warnings-discord-bot/releases)

# Warnings Discord Bot

<img src="https://assets-global.website-files.com/6257adef93867e50d84d30e2/636e0b5061df29d55a92d945_full_logo_blurple_RGB.svg" width="200"><br>

**Warnings Discord Bot** is a simple Discord bot written in Go that monitors messages for specific keywords and responds with warning messages. It is designed to help maintain a respectful and safe chat environment within your Discord server.

## Features

- **Real-time Chat Monitoring**: Actively scans all chat messages for specified keywords and conditions, ensuring no inappropriate content slips through the cracks.
- **Automated Warning Messages**: Automatically issues customized warning messages when a message matches the specified keywords or conditions, helping to maintain a respectful and safe chat environment.
- **Flexible Configuration**: Allows you to easily specify the keywords and conditions to monitor for, as well as the corresponding warning messages, giving you full control over the bot's behavior.

## Prerequisites

Before you can run the bot, make sure you have the following:

- [Go](https://golang.org/) installed on your system.
- A [Discord bot token](https://discord.com/developers/applications) obtained by creating a Discord application and bot.

## How to Use

### Run bot locally

1. Clone [this repository](https://github.com/rickstaa/warnings-discord-bot) to your local machine.
2. Modify the [config.json](config/config.json) file to include the keywords you want to monitor and the warning messages you want to send in response to keyword matches.
3. Setup a discord application (see [this guide](https://discordjs.guide/preparations/setting-up-a-bot-application.html#what-is-a-token-anyway)). Ensure that the [message content intent](https://discord.com/developers/docs/topics/gateway#list-of-intents) is enabled. Also, ensure that the `Send Messages` and `Read Message History` permissions are requested on the URL Generator step.
4. Install the Golang dependencies using `go get`.
5. Build the bot using `go build`
6. Rename the `.env.template` file to `.env` and insert the required environmental variables.
7. Run the bot using `./warnings-discord-bot`.

### Running the bot with Docker

The Warnings Discord Bot can be run using the Docker image available on [Docker Hub](https://hub.docker.com/r/rickstaa/warnings-discord-bot). To pull and run the bot from Docker Hub, use the following command:

```bash
docker run --name warnings-discord-bot rickstaa/warnings-discord-bot:latest
```

Please ensure that you have a `.env` file in the current working directory that contains the required environmental variables. You can find an example of this file [here](./.env.template) or add the `DISCORD_BOT_TOKEN` as an environmental variable to the `docker run` command.

> [!NOTE]
> This repository also contains a [DockerFile](./Dockerfile) and [docker-compose.yml](./docker-compose.yml) file. These files can be used to build and run the bot locally. To do this, clone this repository and run `docker compose up` in the repository's root directory.

## Configuration

The bot's behaviour can be customized through the [config.json](config/config.json) file. This JSON file contains several fields that define the bot's keyword monitoring and response behaviour.

Here's a breakdown of the fields:

- **keyword_lists**: This is an array of objects, each representing a unique set of conditions for the bot to monitor. Each object in this array has the following properties:
  - **keywords**: An array of strings. The bot will monitor chat messages for these keywords.
  - **warning_message**: A string that defines the warning message the bot will send when a chat message matches the keywords.
  - **external_link_required**: A boolean value. If set to `true`, the bot will only issue a warning if the message contains both the specified keywords and an external link.
  - **excluded_roles**: An array of strings. If specified, the bot will not issue a warning if the message author has at least one of the roles in this list.
  - **required_roles**: An array of strings. If specified, the bot will only issue a warning if the message author has at least one of the roles in this list.

Please note that all string comparisons performed by the bot are case-insensitive.

## Contributing

We welcome contributions 🚀! Please see our [Contributing Guide](CONTRIBUTING.md) for more details.

## License

This project is licensed under the [MIT License](LICENSE).
