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

var ErrSimulatedSendMessage = errors.New("simulated error from sendMessage")

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

		bot := NewTelegramBot(mock, "")
		msg := tgbotapi.NewMessage(123, test.message)
		got := bot.sendMessage(msg)
		want := test.expectedError

		if got != want {
			t.Errorf("got %s, want %s", got, want)
		}
	}
}
