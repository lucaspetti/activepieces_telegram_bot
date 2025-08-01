package bot

type Config struct {
	telegramApiToken string
	authorizedUserID int64
	webhookURL       string
}

func NewConfig(
	authorizedUserID int64,
	telegramApiToken,
	webhookURL string,
) *Config {
	return &Config{
		telegramApiToken: telegramApiToken,
		authorizedUserID: authorizedUserID,
		webhookURL:       webhookURL,
	}
}
