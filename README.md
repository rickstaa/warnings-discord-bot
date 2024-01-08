# Warnings Discord Bot

![Discord Logo](discord_logo.png)

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

1. Clone [this repository](https://github.com/rickstaa/warnings-discord-bot) to your local machine.
2. Modify the [config.json](config/config.json) file to include the keywords you want to monitor and the warning messages you want to send in response to keyword matches.
3. Setup a discord application (see [this guide](https://discordjs.guide/preparations/setting-up-a-bot-application.html#what-is-a-token-anyway)). Ensure that the [message content intent](https://discord.com/developers/docs/topics/gateway#list-of-intents) is enabled. Also, ensure that the `Send Messages` permission is requested on the URL Generator step.
4. Install the Golang dependencies using `go get`.
5. Build the bot using `go build`
6. Rename the `.env.template` file to `.env` and insert the required environmental variables.
7. Run the bot using `./warnings-discord-bot`.

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for more details.

## License

This project is licensed under the [MIT License](LICENSE).
# warnings-discord-bot
