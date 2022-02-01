package cron

import (
	"fmt"
	"strings"
	"time"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/matcher"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/util"
)

func (c *command) GetMatcher() matcher.Matcher {
	return matcher.NewAdminMatcher(
		c.SlackClient,
		matcher.NewTextMatcher("list crons", c.ListCrons),
		false,
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
