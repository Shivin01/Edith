package admin

import (
	"bytes"
	"fmt"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/config"
	"os"
	"strings"
)

const logChars = 4000

// newBotLogCommand prints the recent edith.log as slack command
func newBotLogCommand(base bot.BaseCommand, cfg *config.Config) bot.Command {
	return &botLogCommand{
		base,
		cfg,
	}
}

type botLogCommand struct {
	bot.BaseCommand
	cfg *config.Config
}

func (c *botLogCommand) GetMatcher() matcher.Matcher {
	return matcher.NewAdminMatcher(
		c.SlackClient,
		matcher.NewTextMatcher("edith log", c.showBotLog),
	)
}

func (c *botLogCommand) showBotLog(match matcher.Result, message msg.Message) {
	log := c.readFile(c.cfg.Logger.File, logChars)
	parts := strings.SplitN(string(log), "\n", 2)
	if len(parts) <= 1 {
		c.SendMessage(message, "No logs so far")
		return
	}

	c.SendMessage(message, fmt.Sprintf("The most recent messages:\n```%s```", parts[1]))
}

func (c *botLogCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "edith log",
			Description: "displays the last log entries of the edith",
			Examples: []string{
				"edith log",
			},
		},
	}
}

// get the last X characters from the given file
func (c *botLogCommand) readFile(filename string, chars int64) []byte {
	buf := make([]byte, chars)
	file, err := os.Open(filename) // #nosec
	if err != nil {
		return buf
	}
	defer file.Close()

	stat, _ := os.Stat(filename)
	start := stat.Size() - chars
	if start < 0 {
		start = 0
	}

	file.ReadAt(buf, start)

	return bytes.Trim(buf, "\x00")
}
