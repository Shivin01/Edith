package admin

import (
	"fmt"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/bot/util"
	"time"
)

// newPingCommand just prints a PING with the needed time from client->slack->edith server
func newPingCommand(base bot.BaseCommand) bot.Command {
	return &pingCommand{
		base,
	}
}

type pingCommand struct {
	bot.BaseCommand
}

func (c *pingCommand) GetMatcher() matcher.Matcher {
	return matcher.NewTextMatcher("ping", c.ping)
}

func (c *pingCommand) ping(match matcher.Result, message msg.Message) {
	c.SendMessage(message, fmt.Sprintf(
		"PONG in %s",
		util.FormatDuration(time.Since(message.GetTime())),
	))
}

func (c *pingCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "ping",
			Description: "just prints a PING with the needed time from client->slack->edith server",
			Examples: []string{
				"ping",
			},
		},
	}
}
