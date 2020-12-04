package slackstatus

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/pteich/errors"
	"github.com/pteich/go-timeout-httpclient"
)

const (
	defaultRetrySec = 1
	maxRetry        = 10
)

var (
	errTooManyRequests = errors.New("too many requests")
)

// Message defines the message that should be send to Slack
type Message struct {
	WebhookURL       string `mapstructure:"webhook_url"`
	Channel          string `mapstructure:"channel"`
	Username         string `mapstructure:"username"`
	IconEmoji        string `mapstructure:"icon_emoji"`
	Footer           string `mapstructure:"footer"`
	RetryRatelimited bool   `mapstructure:"retry_ratelimited"`
	RetryBackground  bool   `mapstructure:"retry_background"`
	httpClient       *http.Client
	retryCount       int
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

	if msg.httpClient == nil {
		msg.createHttpClient()
	}

	req, err := http.NewRequest(http.MethodPost, msg.WebhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-type", "application/json")

	doRequest := func() (int, error) {
		resp, err := msg.httpClient.Do(req)
		if err != nil {
			return 0, err
		}

		defer resp.Body.Close()
		responseBody, _ := ioutil.ReadAll(resp.Body)

		if resp.StatusCode == http.StatusTooManyRequests {
			retryAfter := resp.Header.Get("retry-after")
			retrySec, err := strconv.Atoi(retryAfter)
			if err != nil || retrySec == 0 {
				retrySec = defaultRetrySec
			}
			return retrySec, errTooManyRequests
		}

		if resp.StatusCode != http.StatusOK {
			return 0, errors.Errorf("got wrong status %d. body: %s", resp.StatusCode, responseBody)
		}

		return 0, nil
	}

	retrySec, err := doRequest()
	if errors.Is(err, errTooManyRequests) && msg.RetryRatelimited {
		if msg.retryCount < maxRetry {
			return errors.Wrapf(err, "retry limit %d reached", maxRetry)
		}
		retryFunc := func(message string, color string) error {
			msg.retryCount++
			<-time.After(time.Duration(retrySec) * time.Second)
			return msg.Send(message, color)
		}

		if msg.RetryBackground {
			go retryFunc(message, color)
		} else {
			return retryFunc(message, color)
		}
	}

	return err
}

func (msg *Message) createHttpClient() {
	msg.httpClient = timeouthttp.NewClient(timeouthttp.Config{
		RequestTimeout: 5,
		ConnectTimeout: 5,
	})
}
