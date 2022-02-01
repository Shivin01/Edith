package edith

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client/edith"
	"github.com/pkg/errors"
)

// newLeaveRequestCommand used to request leaves
func newLeaveRequestCommand(command edithCommand) bot.Command {
	return &leaveRequestCommand{
		command,
	}
}

type leaveRequestCommand struct {
	edithCommand
}

func (c *leaveRequestCommand) GetMatcher() matcher.Matcher {
	return matcher.NewAuthorizedMatcher(
		c.SlackClient,
		matcher.NewRegexpMatcher(`apply leave (?P<startingDate>\d{4}-\d{2}-\d{2}) (?P<days>\d+) (?P<type>casual|sick|maternity) (?P<kind>full_day|half_day|restricted)`, c.run),
		true,
	)
}

func (c *leaveRequestCommand) run(match matcher.Result, message msg.Message) {
	t := match.GetString("startingDate")
	days := match.GetInt("days")
	leaveType := strings.ToUpper(match.GetString("type"))
	kind := strings.ToUpper(match.GetString("kind"))
	startingDate, err := time.Parse("2006-01-02", t)
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New(fmt.Sprintf("invalid starting date: %s", err.Error())),
		)
		return
	}

	if !startingDate.After(time.Now()) && time.Now().Format("2006-01-02") != startingDate.Format("2006-01-02") {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("Future dates can be applied for leave"),
		)
		return
	}

	err = c.client.RequestForLeave(context.TODO(), message.DBUser.AccessToken, edith.LeaveRequest{
		StartDateTime: startingDate.Unix(),
		StopDateTime:  startingDate.Add(time.Hour * 24 * time.Duration(days)).Unix(),
		Kind:          kind,
		LeaveType:     leaveType,
	})
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while request for leave..."),
		)
		return
	}
	c.SlackClient.AddReaction("✅", message)
	c.SlackClient.SendMessage(message, "Successfully requested leave, please wait for it to be approved.")

	// TODO: Send message to be approved
}

func (c *leaveRequestCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "appy leave <YYYY-MM-DD> <days> <casual|sick|maternity> <full_day|half_day|restricted>",
			Description: "apply leaves, cancel leave or modify existing leave",
			Examples: []string{
				"apply leave 2022-02-12 3 sick full_day",
				"apply leave 2022-02-12 3 casual full_day",
			},
			Category: category,
		},
	}
}
