package edith

import (
	"context"
	"fmt"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client/edith"
	"github.com/pkg/errors"
	"strings"
	"time"
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
	return matcher.NewPrivateMatcher(
		c.SlackClient,
		matcher.NewAuthorizedMatcher(
			c.SlackClient,
			matcher.NewRegexpMatcher(`take leave (?P<startingDate>\d{4}-\d{2}-\d{2}) (?P<days>\d+) (?P<type>casual|sick|maternity) (?P<kind>full_day|half_day|restricted)`, c.run),
		),
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

	if !startingDate.After(time.Now()) {
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
	c.SlackClient.SendMessage(message, "Successfully requested leave, please wait for it to be approved.")
}

func (c *leaveRequestCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "take leave <YYYY-MM-DD> <days> <casual|sick|maternity> <full_day|half_day|restricted>",
			Description: "apply leaves, cancel leave or modify existing leave",
			Examples: []string{
				"take leave 2022-02-12 3 sick full_day",
				"take leave 2022-02-12 3 casual full_day",
			},
			Category: category,
		},
	}
}
