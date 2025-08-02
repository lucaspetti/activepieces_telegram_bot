package bot

import (
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

func (b TelegramBot) sendMessage(msg tgbotapi.MessageConfig) error {
	_, err := b.botApi.Send(msg)
	return err
}

func (b TelegramBot) handleMessage(update tgbotapi.Update) error {
	// TODO: get the message and send to Webhook URL with json in the body
	message := update.Message.Text
	userName := update.Message.From.FirstName
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	// voiceMsg := update.Message.Voice

	switch message {
	case "/start":
		msg.Text = "Hello there, " + userName
		return b.sendMessage(msg)
	default:
		msg.Text = "Calling webhook..."
		err := b.sendMessage(msg)
		if err != nil {
			return err
		}

		msg.Text = "Done"
		return b.sendMessage(msg)
	}
}
