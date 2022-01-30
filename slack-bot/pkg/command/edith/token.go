package edith

import (
	"context"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/db"
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
	return matcher.NewPrivateMatcher(
		c.SlackClient,
		matcher.NewAuthorizedMatcher(
			c.SlackClient,
			matcher.NewRegexpMatcher(`token (?P<user>[\w\-_\\/]+) (?P<password>[\w\-_\\/!@#\$%\^\&*\)\(+=\.]+)`, c.run),
		),
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

func (c *tokenCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "token <username> <password>",
			Description: "Save token in the database to use in the future until expire",
			Examples: []string{
				"token manojm manojm",
			},
			Category: category,
		},
	}
}
