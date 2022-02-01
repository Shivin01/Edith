package edith

import (
	"context"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/matcher"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
	"github.com/pkg/errors"
)

// newLeaveRequestCommand used to request leaves
func newLeaveApproveCommand(command edithCommand) bot.Command {
	return &leaveApproveCommand{
		command,
	}
}

type leaveApproveCommand struct {
	edithCommand
}

func (c *leaveApproveCommand) GetMatcher() matcher.Matcher {
	return matcher.NewAuthorizedMatcher(
		c.SlackClient,
		matcher.NewRegexpMatcher(`approve leave (?P<leave_id>\d+)`, c.run),
		true,
	)
}

func (c *leaveApproveCommand) run(match matcher.Result, message msg.Message) {
	leaveId := match.GetInt("leave_id")

	err := c.client.ApproveLeave(context.TODO(), leaveId, message.DBUser.AccessToken)
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while approval of leave..."),
		)
		return
	}
	c.SlackClient.AddReaction("✅", message)
	c.SlackClient.SendMessage(message, "Successfully appoved leave request.")
}

func (c *leaveApproveCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "approve leave <leave_id>",
			Description: "approve leave request",
			Examples: []string{
				"approve leave 1",
			},
			Category: category,
		},
	}
}
