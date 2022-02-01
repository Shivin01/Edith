package edith

import (
	"context"

	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// newPingCommand just prints a PING with the needed time from client->slack->edith server
func newMarkAttendanceCommand(command edithCommand) bot.Command {
	return &markAttendanceCommand{
		command,
	}
}

type markAttendanceCommand struct {
	edithCommand
}

func (c *markAttendanceCommand) GetMatcher() matcher.Matcher {
	return matcher.NewAuthorizedMatcher(
		c.SlackClient,
		matcher.NewTextMatcher("mark attendance", c.run),
		true,
	)
}

func (c *markAttendanceCommand) run(match matcher.Result, message msg.Message) {
	err := c.client.MarkAttendance(context.TODO(), message.DBUser.AccessToken)
	if err != nil {
		log.Errorf("got error: %s", err.Error())
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while marking attendance for today."),
		)
		return
	}
	c.SlackClient.AddReaction("✅", message)
	c.SlackClient.SendMessage(message, "Marked attendance for today.")
}

func (c *markAttendanceCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "mark attendance",
			Description: "mark attendance for the user",
			Examples: []string{
				"mark attendance",
			},
			Category: category,
		},
	}
}
