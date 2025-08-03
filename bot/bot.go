package bot

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// BotAPI defines the interface for sending messages, to allow mocking in tests.
type BotAPI interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
	GetFileDirectURL(fileID string) (string, error)
}

// TelegramBot uses a BotAPI interface instead of the concrete *tgbotapi.BotAPI
type TelegramBot struct {
	botApi     BotAPI
	webhookURL string
}

func NewTelegramBot(botApi BotAPI, webhookURL string) *TelegramBot {
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

func (b TelegramBot) sendWebhookAudio(update tgbotapi.Update, fileID string) (result string, err error) {
	chatId := strconv.Itoa(int(update.Message.Chat.ID))
	directURL, err := b.botApi.GetFileDirectURL(fileID)
	if err != nil {
		return "", err
	}

	reqBody := map[string]string{
		"audio_url": directURL,
		"chat_id":   chatId,
	}

	return b.callWebhook(reqBody)
}

func (b TelegramBot) sendWebhookText(update tgbotapi.Update, message string) (result string, err error) {
	chatId := strconv.Itoa(int(update.Message.Chat.ID))
	reqBody := map[string]string{
		"text":    message,
		"chat_id": chatId,
	}

	return b.callWebhook(reqBody)
}

func (b TelegramBot) callWebhook(reqBody map[string]string) (result string, err error) {
	client := &http.Client{}
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
	message := update.Message

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

	if message.Text == "/start" {
		userName := update.Message.From.FirstName
		msg.Text = "Hello there, " + userName
	} else if message.Voice != nil {
		voiceFileID := update.Message.Voice.FileID
		webhookResult, err := b.sendWebhookAudio(update, voiceFileID)
		msg.Text = webhookResult
		if err != nil {
			msg.Text = "Error calling webhook"
		}
	} else {
		webhookResult, err := b.sendWebhookText(update, message.Text)
		msg.Text = webhookResult
		if err != nil {
			msg.Text = "Error calling webhook"
		}
	}

	return b.sendMessage(msg)
}
