package webhook

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type WebhookClient struct {
	webhookURL string
}

func NewWebhookClient(webhookURL string) *WebhookClient {
	WebhookClient := &WebhookClient{
		webhookURL: webhookURL,
	}

	return WebhookClient
}

func (wc WebhookClient) Post(reqBody map[string]string) (result string, err error) {
	client := &http.Client{}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", wc.webhookURL, bytes.NewBuffer(jsonBody))
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
