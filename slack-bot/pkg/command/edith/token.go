package edith

import (
	"context"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/matcher"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
	"github.com/Shivin01/Edith/slack-bot/pkg/db"
)

// newTokenCommand
func newTokenCommand(base edithCommand) bot.Command {
	return &tokenCommand{
		base,
	}
}

type tokenCommand struct {
	edithCommand
}

func (c *tokenCommand) GetMatcher() matcher.Matcher {
	return matcher.NewAuthorizedMatcher(
		c.SlackClient,
		matcher.NewRegexpMatcher(`token (?P<user>[\w\-_\\/]+) (?P<password>[\w\-_\\/!@#\$%\^\&*\)\(+=\.]+)`, c.run),
		true,
	)
}

func (c *tokenCommand) run(match matcher.Result, message msg.Message) {
	username := match.GetString("user")
	password := match.GetString("password")
	res, err := c.client.GetToken(context.TODO(), username, password)
	if err != nil {
		// handle error
	}

	if err := c.DB.Debug().Model(&db.User{}).Where("ID = ?", message.User).
		Select("AccessToken", "RefreshToken").
		Updates(map[string]interface{}{
			"access_token": res.Token,
		}).Error; err != nil {
		// handle error
	}
}
