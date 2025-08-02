package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func Start(config Config) {
	bot, err := tgbotapi.NewBotAPI(config.telegramApiToken)
	if err != nil {
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

		if update.Message.From.ID != config.authorizedUserID {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			msg.Text = "Unauthorized user"
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		} else {
			err = TelegramBot.handleMessage(update)
			if err != nil {
				panic(err)
			}
		}
	}
}
