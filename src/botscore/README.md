# Sneat bots

At the moment only Telegram bots are supported.

## Testing Telegram bots locally

To test Telegram bots locally you need to set environment 2 sets of variables.

### 1. A comma separate list of bot IDs and assigned parameters like:

```shell
export SNEAT_TG_DEV_BOTS=<BOT1_ID>:<BOT1_PROFILE>,<BOT2_ID>:<BOT2_PROFILE>...
```

You should use bot IDs of bots you registered yourself with [@BotFather](https://t.me/BotFather).

To register bot use some reverse-proxy (like ngrok) and call `/bot/tg/set-webhook?code=<BOT_ID>` endpoint.
More details at https://github.com/bots-go-framework/bots-fw-telegram

Possible values for <BOT#_PROFILE> are:

- [`sneat`](../botprofiles/sneatbot) - profile for [@SneatBot](https://t.me/SneatBot)
- [`listus`](../botprofiles/listusbot) - profile for [@Listus_Bot](https://t.me/Listus_bot)

### 2. Telegram bot token for each bot:

```shell
export TELEGRAM_BOT_TOKEN_<BOT_ID_1>=<TELEGRAM_BOT_TOKEN_1>
export TELEGRAM_BOT_TOKEN_<BOT_ID_2>=<TELEGRAM_BOT_TOKEN_2>
```
