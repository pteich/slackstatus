package slackstatus

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type Message struct {
	WebhookURL string
	Channel    string
	Username   string
	IconEmoji  string ":monkey_face:"
	Footer     string ""
}

const COLOR_WARN string = "warn"
const COLOR_DANGER string = "danger"
const COLOR_GOOD string = "good"

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
