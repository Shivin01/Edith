package cron

import (
	"fmt"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/bot/util"
	"strings"
	"time"
)

func (c *command) GetMatcher() matcher.Matcher {
	return matcher.NewAuthorizedMatcher(
		c.SlackClient,
		matcher.NewTextMatcher("list crons", c.ListCrons),
	)
}

func (c *command) ListCrons(match matcher.Result, message msg.Message) {
	text := fmt.Sprintf("*%d crons:*\n", len(c.cfg))

	now := time.Now()
	for i, entry := range c.cron.Entries() {
		text += fmt.Sprintf(
			" - `%s`, next in %s (`%s`)\n",
			c.cfg[i].Schedule,
			util.FormatDuration(entry.Next.Sub(now)),
			strings.Join(c.cfg[i].Commands, "; "),
		)
	}

	c.SendMessage(message, text)
}

func (c *command) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "list crons",
			Description: "list the registered crons and the next execution time",
			HelpURL:     "",
			Examples: []string{
				"list crons",
			},
		},
	}
}
