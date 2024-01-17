package uds

import (
	"github.com/go-resty/resty/v2"
	"time"
)

const BaseUri = "https://api.uds.app/partner/v2/"

// Client
// Структура клиента UDS
type Client struct {
	client *resty.Client
}

func NewClient(clientID, apiKey string) *Client {
	client := resty.New().SetBaseURL(BaseUri).
		SetRetryCount(10).
		SetRetryWaitTime(time.Second).
		SetHeaders(map[string]string{
			"Accept":          "application/json",
			"Accept-Charset":  "utf-8",
			"Content-Type":    "application/json",
			"Accept-Language": "ru-RU, ru",
		}).
		SetBasicAuth(clientID, apiKey)

	return &Client{client}
}
