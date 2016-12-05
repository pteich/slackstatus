package slackstatus

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Message defines the message that should be send to Slack
type Message struct {
	WebhookURL string `mapstructure:"webhook_url"`
	Channel    string `mapstructure:"channel"`
	Username   string `mapstructure:"username"`
	IconEmoji  string `mapstructure:"icon_emoji"`
	Footer     string `mapstructure:"footer"`
}

// ColorWarning is a predefined color for a warning (yellow)
const ColorWarning string = "warning"

// ColorDanger is a predefined color for a dangerous condition (red)
const ColorDanger string = "danger"

// ColorGood is a predefined color for a normal information (green)
const ColorGood string = "good"

// Send sends the given message to Slack
func (msg *Message) Send(message string, color string) error {

	body, err := json.Marshal(composeMessage(msg, message, color))
	if err != nil {
		return err
	}

	resp, err := http.Post(msg.WebhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		responseBody, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(responseBody))
	}

	return nil
}
