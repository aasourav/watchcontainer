package notify

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
)

type SlackPayload struct {
	Text    string `json:"text"`
	Channel string `json:"channel,omitempty"`
}

func Slack(webhook string, msg string, channel string) error {
	p := SlackPayload{
		Text:    msg,
		Channel: channel,
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	b, _ := json.Marshal(p)
	_, err := client.Post(webhook, "application/json", bytes.NewBuffer(b))
	return err
}
