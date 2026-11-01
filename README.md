# Kudasai_Bot: an anime news bot for Telegram

Get the latest anime, manga and light novel news in your Telegram chat by sending a single command. Kudasai_Bot is a Go program that fetches news from ([SomosKudasai](https://somoskudasai.com/)) RSS feed. The news is sent to your Telegram chat through your bot.

![Kudasai_Bot demo](assets/Kudasai_Bot.gif)

**Built using**: Go, ([viper](https://github.com/spf13/viper)) (environment variables), ([lumberjack](https://github.com/natefinch/lumberjack)) (log rotation) and Docker.

## Why Kudasai_Bot?

I used to follow a Telegram channel that posted anime news in Spanish, but it is no longer maintained. So, if you wanted to quickly check news about your favorite anime/manga, you would have to navigate to the *SomosKudasai* website and browse each category one by one.

That's why I decided to create Kudasai_Bot: it reads the site's RSS feed, keeps only the anime, manga, and light novel categories, sorts the news by publication date, and delivers it to your Telegram chat.

## Installation

### Prerequisites

1. [Create your Telegram bot](https://core.telegram.org/bots#how-do-i-create-a-bot) and copy its **bot token**.
2. Get your **chat ID** ([How to get Telegram Bot Chat ID](https://gist.github.com/nafiesl/4ad622f344cd1dc3bb1ecbe468ff9f8a)).
3. Install [Docker](https://docs.docker.com/get-docker/) (option 1) or [Go](https://go.dev/doc/install) (option 2).

Then, you have two options:

### Option one: Docker

- Pull the image:

> docker pull nesar9/kudasai


- Run the container:

> docker run -d -e token_id=XXXXXX -e chat_id=XXXXXX nesar9/kudasai


### Option two: Go build/run

- Clone the repo:

> git clone https://github.com/arturogood17/kudasai_bot

> cd kudasai_bot


- Run it directly with:

> token_id=XXXXXX chat_id=XXXXXX go run cmd/kudasai_bot/main.go


- Build the binary and run it:

> go build -o kudasai_bot cmd/kudasai_bot/main.go

> token_id=XXXXXX chat_id=XXXXXX ./kudasai_bot


## Usage

Once the bot is running, open the chat with your bot in Telegram and send one of the following commands:

- `fetch`: The bot will reply with the latest news from SomosKudasai (filtered by anime/manga/light novels).
- `fetch all`: The bot will reply with all the news from SomosKudasai (without filtering).
- `reset`: Clears the log file of already-sent news, so the filtered news can be fetched again.

Each news item is sent as a message with its image, title, and link.

Right now it fetches the news of the following categories:

### Supported categories

- Anime
- Manga
- Light novels

## Roadmap

- [X] Add a `reset` command to resend all the filtered news
- [ ] Let users fetch news from a specific category
- [X] Let users fetch news from all of the site's categories
