package edith

import (
	"context"
	"fmt"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client/edith"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"strings"
	"time"
)

// newPingCommand just prints a PING with the needed time from client->slack->edith server
func newLeaveListCommand(command edithCommand) bot.Command {
	return &leaveListCommand{
		command,
	}
}

type leaveListCommand struct {
	edithCommand
}

func (c *leaveListCommand) GetMatcher() matcher.Matcher {
	return matcher.NewPrivateMatcher(
		c.SlackClient,
		matcher.NewManagerMatcher(
			c.SlackClient,
			matcher.NewGroupMatcher(
				matcher.NewTextMatcher("list leaves", c.run),
				matcher.NewTextMatcher("list leaves for approval", c.getNotApprovedLeaves),
			),
		),
	)
}

func (c *leaveListCommand) handleList(message msg.Message, leaves []*edith.LeaveResponse) {
	c.SlackClient.AddReaction("✅", message)

	sections := make([]slack.Block, 0)
	sections = append(sections, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*Leaves*", false, false), nil, nil))
	for _, leave := range leaves {
		fieldSlice := make([]*slack.TextBlockObject, 0)
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*ID:*\n%d", leave.ID), false, false))
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Start Date:*\n%v", time.Unix(int64(leave.StartDateTime), 0)), false, false))
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*End Date:*\n%v", time.Unix(int64(leave.StopDateTime), 0)), false, false))
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Type:*\n%s", leave.LeaveType), false, false))
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Kind:*\n%s", leave.Kind), false, false))
		name := strings.Trim(leave.Employee.FirstName, " ")
		if leave.Employee.MiddleName != "" {
			name += fmt.Sprintf(" %s", leave.Employee.MiddleName)
		}
		if leave.Employee.LastName != "" {
			name += fmt.Sprintf(" %s", leave.Employee.LastName)
		}
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Full Name:*\n%s", name), false, false))
		sections = append(sections, slack.NewSectionBlock(nil, fieldSlice, nil))
	}

	c.SlackClient.SendBlockMessage(message, sections)
}

func (c *leaveListCommand) run(match matcher.Result, message msg.Message) {
	res, err := c.client.ListLeaves(context.TODO(), message.DBUser.AccessToken)
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while getting leaves list from server."),
		)
		return
	}

	c.handleList(message, res)
}

func (c *leaveListCommand) getNotApprovedLeaves(match matcher.Result, message msg.Message) {
	res, err := c.client.ListLeavesForApproval(context.TODO(), message.DBUser.AccessToken)
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while getting leaves list from server."),
		)
		return
	}

	c.handleList(message, res)
}

func (c *leaveListCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "list leaves",
			Description: "check who are on leave",
			Examples: []string{
				"list leaves",
			},
			Category: category,
		},
		{
			Command:     "list leaves for approval",
			Description: "check who are on leave",
			Examples: []string{
				"list leaves for approval",
			},
			Category: category,
		},
	}
}
