package notify

import (
	"bytes"
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

	b, _ := json.Marshal(p)
	_, err := http.Post(webhook, "application/json", bytes.NewBuffer(b))
	return err
}
