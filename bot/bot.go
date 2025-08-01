package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	botApi *tgbotapi.BotAPI
}

func NewTelegramBot(botApi *tgbotapi.BotAPI) *TelegramBot {
	TelegramBot := &TelegramBot{
		botApi: botApi,
	}

	return TelegramBot
}

func (b TelegramBot) buildResponse(message string) (response string) {
	// TODO: get the message and send to Webhook URL with json in the body
	switch message {
	case "/start":
		return "Hello there!"
	default:
		return "Calling webhook..."
	}
}

func Start(config Config) {
	bot, err := tgbotapi.NewBotAPI(config.telegramApiToken)
	if err != nil {
		// TODO: Retry message or handle error better
		panic(err)
	}

	TelegramBot := NewTelegramBot(bot)
	bot.Debug = true

	// Create a new UpdateConfig struct with an offset of 0. Offsets are used
	// to make sure Telegram knows we've handled previous values and we don't
	// need them repeated.
	updateConfig := tgbotapi.NewUpdate(0)

	// Tell Telegram we should wait up to 30 seconds on each request for an
	// update. This way we can get information just as quickly as making many
	// frequent requests without having to send nearly as many.
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		// Only look at messages for now and discard any other updates.
		if update.Message == nil {
			continue
		}

		message := update.Message.Text
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if update.Message.From.ID != config.authorizedUserID {
			msg.Text = "Unauthorized user"
		} else {
			msgText := TelegramBot.buildResponse(message)
			if strings.Contains(msgText, "<pre>") {
				msg.ParseMode = tgbotapi.ModeHTML
			}

			msg.Text = msgText
		}

		if _, err := bot.Send(msg); err != nil {
			// TODO: Retry message or handle error better
			panic(err)
		}
	}
}
