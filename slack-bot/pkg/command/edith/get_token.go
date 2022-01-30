package edith

import (
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
)

// newGetTokenCommand just prints a PING with the needed time from client->slack->edith server
func newGetTokenCommand(command edithCommand) bot.Command {
	return &getToken{
		command,
	}
}

type getToken struct {
	edithCommand
}

func (c *getToken) GetMatcher() matcher.Matcher {
	return matcher.NewPrivateMatcher(
		c.SlackClient,
		matcher.NewAuthorizedMatcher(
			c.SlackClient,
			matcher.NewTextMatcher("get token", c.run),
		),
	)
}

func (c *getToken) run(match matcher.Result, message msg.Message) {
	c.SlackClient.AddReaction("✅", message)
	c.SlackClient.SendMessage(message, message.DBUser.AccessToken)
}

func (c *getToken) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "get token",
			Description: "get token",
			Examples: []string{
				"get token",
			},
			Category: category,
		},
	}
}
