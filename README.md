# slackstatus
Simple Go library to post formatted status messages to a Slack channel

## Usage
```go
package main

import "github.com/pteich/slackstatus"

var slackmsg = slackstatus.Message{WebhookURL: "https://hooks.slack.com/services/XXXX", Username: "slackstatus", Channel: "#channelname", IconEmoji: ":monkey_face:", Footer: "Version 1.0.0"}

func main() {

	slackmsg.Send("Hello Slackstatus! Everything works fine.", slackstatus.COLOR_GOOD)
	slackmsg.Send("Oh crap, something went wrong!", slackstatus.COLOR_WARN)
	slackmsg.Send("Damn, we are in serious trouble!", slackstatus.COLOR_DANGER)
	slackmsg.Send("Ok.", "#439FE0")
  
}
```
