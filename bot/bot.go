package bot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	botApi     *tgbotapi.BotAPI
	webhookURL string
}

func NewTelegramBot(botApi *tgbotapi.BotAPI, webhookURL string) *TelegramBot {
	TelegramBot := &TelegramBot{
		botApi:     botApi,
		webhookURL: webhookURL,
	}

	return TelegramBot
}

func (b TelegramBot) sendMessage(msg tgbotapi.MessageConfig) error {
	_, err := b.botApi.Send(msg)
	return err
}

func (b TelegramBot) callWebhook(update tgbotapi.Update, message string) (result string, err error) {
	client := &http.Client{}
	chatId := strconv.Itoa(int(update.Message.Chat.ID))
	reqBody := map[string]string{
		"text":    message,
		"chat_id": chatId,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", b.webhookURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return "Webhook called", nil
}

func (b TelegramBot) handleMessage(update tgbotapi.Update) error {
	message := update.Message.Text
	userName := update.Message.From.FirstName
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	// voiceMsg := update.Message.Voice

	switch message {
	case "/start":
		msg.Text = "Hello there, " + userName
	default:
		webhookResult, err := b.callWebhook(update, message)
		msg.Text = webhookResult
		if err != nil {
			msg.Text = "Error calling webhook"
		}
	}

	return b.sendMessage(msg)
}
