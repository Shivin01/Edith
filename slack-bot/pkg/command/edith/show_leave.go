package edith

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/matcher"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
	"github.com/Shivin01/Edith/slack-bot/pkg/client/edith"
	"github.com/Shivin01/Edith/slack-bot/pkg/db"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
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
	return matcher.NewManagerMatcher(
		c.SlackClient,
		matcher.NewGroupMatcher(
			matcher.NewTextMatcher("list leaves", c.run),
			matcher.NewTextMatcher("list leaves for approval", c.getNotApprovedLeaves),
		),
		true,
	)
}

func (c *leaveListCommand) handleList(message msg.Message, leaves []*edith.LeaveResponse, showId bool) {
	c.SlackClient.AddReaction("✅", message)

	if len(leaves) == 0 {
		c.SlackClient.SendMessage(message, "there are no leaves yet, I will tell you whenever you get leave requests")
		return
	}

	sections := make([]slack.Block, 0)
	sections = append(sections, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*Leaves*", false, false), nil, nil))
	for _, leave := range leaves {
		fieldSlice := make([]*slack.TextBlockObject, 0)
		if showId {
			fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*ID:*\n%d", leave.ID), false, false))
		}
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Start Date:*\n%v", time.Unix(int64(leave.StartDateTime), 0).Format("2006-01-02")), false, false))
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*End Date:*\n%v", time.Unix(int64(leave.StopDateTime), 0).Format("2006-01-02")), false, false))
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

	if message.IsAdminMessage() {
		var users []db.User
		if err := c.DB.Debug().Model(&db.User{}).Find(&users).Error; err != nil {
			c.SlackClient.AddReaction("❌", message)
			c.SlackClient.ReplyError(
				message,
				errors.New("sorry, couldn't retrieve db information."),
			)
			return
		}
		for _, user := range users {
			if user.IsAdmin() {
				message.User = user.ID
				c.SlackClient.SendBlockEphemeralMessage(message, sections)
			}
		}
		return
	}

	c.SlackClient.SendBlockMessage(message, sections)
}

func (c *leaveListCommand) run(match matcher.Result, message msg.Message) {
	if message.User == "cron" {
		user := &db.User{}
		if err := c.DB.Debug().Model(&db.User{}).Where("designation IN ?", []string{"admin", "hr", "manager"}).First(&user).Error; err != nil {
			c.SlackClient.ReplyError(
				message,
				errors.New("sorry, couldn't retrieve db information."),
			)
			return
		}

		message.DBUser = user
	}

	res, err := c.client.ListLeaves(context.TODO(), message.DBUser.AccessToken)
	if err != nil {
		if !message.AdminMessage {
			c.SlackClient.AddReaction("❌", message)
		}
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while getting leaves list from server."),
		)
		return
	}

	c.handleList(message, res, false)
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

	c.handleList(message, res, true)
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
