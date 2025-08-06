# Activepieces Telegram Bot

- Bot created to simplify interactions with self-hosted activepieces

### Background

If you have a self-hosted instance of [Activepieces](https://github.com/activepieces/activepieces) running on a host that is not exposed, you can use a telegram bot and make it call a webhook on that same instance whenever a message is received.
This repository has the basic setup for that

### HTTP Request to Webhook

- When a text message is received, the bot sends a POST request to the webhook URL with this payload:
```json
{
    "text": "The message text content",
    "chat_id": "123456"
}
```

- When a voice message is received, this is what is sent:
```json
{
    "audio_url": "url_for_the_audio_message",
    "chat_id": "123456"
}
```

Then you can create a flow that starts with a webhook and handle it from there

## Setup

- From Docker image

```
docker pull lucaspetti/activepieces_telegram_bot:latest

# Make sure to set the correct env variables
docker run -d --name activepieces_telegram_bot \
  -e WEBHOOK_URL=$ACTIVEPIECES_FLOW_WEBHOOK_URL \
  -e TELEGRAM_BOT_TOKEN=$YOUR_TELEGRAM_BOT_TOKEN \
  lucaspetti/activepieces_telegram_bot:latest start
```

- With docker compose

```
# Copy and set your env variables in the .env file
cp .env.example .env

# The docker-compose file is an example with a service called "sample_bot", which can be renamed
docker compose up -d
```
