package bot

import (
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type mockBotAPI struct {
	SendFunc func(c tgbotapi.Chattable) (tgbotapi.Message, error)
}

func (m *mockBotAPI) Send(c tgbotapi.Chattable) (tgbotapi.Message, error) {
	if m.SendFunc != nil {
		return m.SendFunc(c)
	}
	return tgbotapi.Message{}, nil
}

func (m *mockBotAPI) GetFileDirectURL(f string) (string, error) {
	return "", nil
}

type mockWebhookClient struct {
	PostFunc func(body map[string]string) (string, error)
}

func (mwc *mockWebhookClient) Post(body map[string]string) (string, error) {
	if mwc.PostFunc != nil {
		return mwc.PostFunc(body)
	}
	return "", nil
}

var (
	ErrSimulatedSendMessage = errors.New("simulated error from sendMessage")
	ErrSimulatedWebhook     = errors.New("simulated error from webhookClient")
)

func TestTelegramBot_sendMessage(t *testing.T) {
	cases := []struct {
		Title               string
		message             string
		expectedError       error
		SendMessageResponse string
		SendMessageError    error
	}{
		{
			Title:         "Success sendMessage",
			message:       "testing success",
			expectedError: nil,
		},
		{
			Title:            "Error on sendMessage",
			message:          "testing error",
			expectedError:    ErrSimulatedSendMessage,
			SendMessageError: ErrSimulatedSendMessage,
		},
	}

	for _, test := range cases {
		mock := &mockBotAPI{
			SendFunc: func(c tgbotapi.Chattable) (tgbotapi.Message, error) {
				return tgbotapi.Message{}, test.SendMessageError
			},
		}
		mockWebhookClient := &mockWebhookClient{
			PostFunc: func(body map[string]string) (string, error) {
				return "", nil
			},
		}

		bot := NewTelegramBot(mock, mockWebhookClient)
		msg := tgbotapi.NewMessage(123, test.message)
		got := bot.sendMessage(msg)
		want := test.expectedError

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	}
}

func TestTelegramBot_handleMessage(t *testing.T) {
	cases := []struct {
		Title               string
		Message             *tgbotapi.Message
		expectedError       error
		SendMessageResponse string
		SendMessageError    error
		SendWebhookError    error
	}{
		{
			Title: "Success with message text /start",
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 123456,
				},
				Text: "/start",
				From: &tgbotapi.User{
					FirstName: "User name",
				},
			},
			expectedError: nil,
		},
		{
			Title: "SendMessageError with message text /start",
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 123456,
				},
				Text: "/start",
				From: &tgbotapi.User{
					FirstName: "User name",
				},
			},
			expectedError:    ErrSimulatedSendMessage,
			SendMessageError: ErrSimulatedSendMessage,
		},
		{
			Title: "Success with any message text",
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 123456,
				},
				Text: "Random message",
			},
			expectedError: nil,
		},
		{
			Title: "SendWebhookError with any message text",
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 123456,
				},
				Text: "Random message",
			},
			expectedError:    nil,
			SendWebhookError: ErrSimulatedWebhook,
		},
		{
			Title: "SendMessageError with any message text",
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 123456,
				},
				Text: "Random message",
			},
			expectedError:    ErrSimulatedSendMessage,
			SendMessageError: ErrSimulatedSendMessage,
		},
		{
			Title: "Success with voice message",
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 123456,
				},
				Voice: &tgbotapi.Voice{
					FileID: "audio_file_id",
				},
			},
			expectedError: nil,
		},
		{
			Title: "SendWebhookError with voice message",
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 123456,
				},
				Voice: &tgbotapi.Voice{
					FileID: "audio_file_id",
				},
			},
			expectedError:    nil,
			SendWebhookError: ErrSimulatedWebhook,
		},
		{
			Title: "SendMessageError with voice message",
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: 123456,
				},
				Voice: &tgbotapi.Voice{
					FileID: "audio_file_id",
				},
			},
			expectedError:    ErrSimulatedSendMessage,
			SendMessageError: ErrSimulatedSendMessage,
		},
	}

	for _, test := range cases {
		mock := &mockBotAPI{
			SendFunc: func(c tgbotapi.Chattable) (tgbotapi.Message, error) {
				return tgbotapi.Message{}, test.SendMessageError
			},
		}

		mockWebhookClient := &mockWebhookClient{
			PostFunc: func(body map[string]string) (string, error) {
				return "", test.SendWebhookError
			},
		}

		bot := NewTelegramBot(mock, mockWebhookClient)
		update := tgbotapi.Update{
			Message: test.Message,
		}
		got := bot.handleMessage(update)
		want := test.expectedError

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	}
}
