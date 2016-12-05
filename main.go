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
	WebhookURL string
	Channel    string
	Username   string
	IconEmoji  string
	Footer     string
}

// COLOR_WARNING is a predefined color for a warning (yellow)
const ColorWarning string = "warning"

// COLOR_DANGER is a predefined color for a dangerous condition (red)
const ColorDanger string = "danger"

// COLOR_GOOD is a predefined color for a normal information (green)
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
