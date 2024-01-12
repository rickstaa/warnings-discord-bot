[![Docker Build](https://github.com/rickstaa/warnings-discord-bot/actions/workflows/docker-build.yml/badge.svg)](https://github.com/rickstaa/warnings-discord-bot/actions/workflows/docker-build.yml)
[![Publish Docker image](https://github.com/rickstaa/warnings-discord-bot/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/rickstaa/warnings-discord-bot/actions/workflows/docker-publish.yml)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/rickstaa/warnings-discord-bot?logo=docker)
](https://hub.docker.com/r/rickstaa/warnings-discord-bot)
[![Latest Release](https://img.shields.io/github/v/release/rickstaa/warnings-discord-bot?label=latest%20release)](https://github.com/rickstaa/warnings-discord-bot/releases)

# Warnings Discord Bot

<img src="https://assets-global.website-files.com/6257adef93867e50d84d30e2/636e0b5061df29d55a92d945_full_logo_blurple_RGB.svg" width="200"><br>

**Warnings Discord Bot** is a Discord bot written in Go that monitors messages for specific keywords or regular expressions and responds with warning messages. It can also send a welcome warning message to new members. It is designed to help maintain a respectful and safe chat environment within your Discord server.

![image](https://github.com/rickstaa/warnings-discord-bot/assets/17570430/fd29682c-bf16-4c90-b8f8-2c0b6ac5d8f4)

## Features

- **Real-time Chat Monitoring**: Actively scans all chat messages for specified keywords or regular expressions, ensuring no inappropriate content slips through the cracks.
- **Automated Warning Messages**: Automatically issues customized warning messages when a message matches the specified keywords or conditions, helping to maintain a respectful and safe chat environment.
- **Welcome Warning Message**: Sends a customizable warning message to new members when they join the server.
- **Flexible Configuration**: Allows you to easily specify the keywords, regular expressions, and conditions to monitor for, as well as the corresponding warning messages, giving you full control over the bot's behaviour.

## Prerequisites

Before you can run the bot, make sure you have the following:

- [Go](https://golang.org/) installed on your system.
- A [Discord bot token](https://discord.com/developers/applications) obtained by creating a Discord application and bot.

## How to Use

### Run bot locally

1. Clone [this repository](https://github.com/rickstaa/warnings-discord-bot) to your local machine.
2. Modify the [config.json](config/config.json) file to include the keywords you want to monitor and the warning messages you want to send in response to keyword matches.
3. Setup a discord application (see [this guide](https://discordjs.guide/preparations/setting-up-a-bot-application.html#what-is-a-token-anyway)). Ensure that the [message content and guild members intents](https://discord.com/developers/docs/topics/gateway#list-of-intents) are enabled. Also, ensure that the `Send Messages` and `Read Message History` permissions are requested on the URL Generator step.
4. Install the Golang dependencies using `go get`.
5. Build the bot using `go build`
6. Rename the `.env.template` file to `.env` and insert the required environmental variables.
7. Run the bot using `./warnings-discord-bot`.

### Running the bot with Docker

The Warnings Discord Bot can be run using the Docker image available on [Docker Hub](https://hub.docker.com/r/rickstaa/warnings-discord-bot). To pull and run the bot from Docker Hub, use the following command:

```bash
docker run --name warnings-discord-bot rickstaa/warnings-discord-bot:latest
```

Please ensure you have a `.env` file in the current working directory containing the required environmental variables. You can find an example of this file [here](./.env.template) or add the `DISCORD_BOT_TOKEN` as an environmental variable to the `docker run` command.

> [!NOTE]
> This repository also contains a [DockerFile](./Dockerfile) and [docker-compose.yml](./docker-compose.yml) file. These files can be used to build and run the bot locally. To do this, clone this repository and run `docker compose up` in the repository's root directory.

## Configuration

The bot's behaviour can be customized through the [config.json](config/config.json) file. This JSON file contains several fields that define the bot's keyword and regular expression monitoring, response behaviour, and join warning message.

Here's a breakdown of the fields:

- **join_warning_message**: A string that defines the message the bot will send to a user when they join the server. If this field is empty, no message will be sent.

- **alert_rules**: An array of objects, each representing a unique alert rule. An alert rule defines conditions that, when met, trigger a bot warning. Each object has the following properties:
  - **keywords**: An array of strings for the bot to monitor. Used only if no regex is provided.
  - **regex_patterns**: An array of regular expressions for the bot to monitor. Overrides keywords if provided.
  - **warning_message**: A string that defines the warning message the bot will send when a chat message matches the keywords or regular expressions.
  - **external_link_required**: A boolean value. If set to `true`, the bot will only issue a warning if the message contains both the specified keywords or regular expressions and an external link.
  - **excluded_roles**: An array of strings. If specified, the bot will not issue a warning if the message author has at least one of the roles in this list.
  - **required_roles**: An array of strings. If specified, the bot will only issue a warning if the message author has at least one of the roles in this list.
  - **omit_members_older_than_days**: An integer value. If set to a positive number, the bot will not issue a warning if the message author has been a member of the community for more days than the specified value. If set to 0 or a negative number, this rule is ignored.

Please note that all keyword and regular expression comparisons performed by the bot are case-insensitive.

> [!IMPORTANT]\
> Remember to escape backslashes when writing regular expressions in the JSON configuration file. For example, `\b(word)\b` should be written as `\\b(word)\\b`.

## Contributing

We welcome contributions 🚀! Please see our [Contributing Guide](CONTRIBUTING.md) for more details.

## License

This project is licensed under the [MIT License](LICENSE).
