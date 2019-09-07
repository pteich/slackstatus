package slackstatus

type field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type attachment struct {
	Fallback   *string  `json:"fallback"`
	Color      *string  `json:"color"`
	PreText    *string  `json:"pretext"`
	AuthorName *string  `json:"author_name"`
	AuthorLink *string  `json:"author_link"`
	AuthorIcon *string  `json:"author_icon"`
	Title      *string  `json:"title"`
	TitleLink  *string  `json:"title_link"`
	Text       *string  `json:"text"`
	ImageURL   *string  `json:"image_url"`
	Fields     []*field `json:"fields"`
	Footer     *string  `json:"footer"`
	FooterIcon *string  `json:"footer_icon"`
}

type payload struct {
	Parse       string       `json:"parse,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	Attachments []attachment `json:"attachments,omitempty"`
}

func composeMessage(slackmessage *Message, text string, color string) payload {

	slackAttachment := attachment{
		Color:  &color,
		Text:   &text,
		Footer: &slackmessage.Footer,
	}

	return payload{
		Username:    slackmessage.Username,
		IconEmoji:   slackmessage.IconEmoji,
		Channel:     slackmessage.Channel,
		Attachments: []attachment{slackAttachment},
	}

}
