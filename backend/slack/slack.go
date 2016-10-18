package slack

import (
	"log"

	slacklib "github.com/bluele/slack"

	"app/models"
)

type Env struct {
	StaticHost string
	hook       *slacklib.WebHook
}

func New(staticHost string, hook *slacklib.WebHook) *Env {
	return &Env{staticHost, hook}
}

func (e *Env) Send(channel string, user *models.User, text string) error {
	payload := &slacklib.WebHookPostPayload{
		Channel:  "#" + channel,
		Username: user.Nickname,
		IconUrl:  "https://www.gravatar.com/avatar/" + user.Gravatar,
		Text: "*Issued by:* <https://" + e.StaticHost + "/users/" + user.Id + "|" +
			user.Nickname + ">\n" + text,
	}
	if e.hook == nil {
		log.Printf("slack.Channel: %v\n", payload.Channel)
		log.Printf("slack.Username: %v\n", payload.Username)
		log.Printf("slack.IconUrl: %v\n", payload.IconUrl)
		log.Printf("slack.Text: %v\n", payload.Text)
		return nil
	}

	return e.hook.PostMessage(payload)
}
