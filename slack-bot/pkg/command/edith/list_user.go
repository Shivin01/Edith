package edith

import (
	"context"
	"fmt"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/db"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"strings"
)

// newPingCommand just prints a PING with the needed time from client->slack->edith server
func newListUserCommand(command edithCommand) bot.Command {
	return &listUserCommand{
		command,
	}
}

type listUserCommand struct {
	edithCommand
}

func (c *listUserCommand) GetMatcher() matcher.Matcher {
	return matcher.NewPrivateMatcher(
		c.SlackClient,
		matcher.NewManagerMatcher(
			c.SlackClient,
			matcher.NewRegexpMatcher(`list users (?P<action>local|server)`, c.run),
		),
	)
}

func (c listUserCommand) formatUsers(users []db.User) []slack.Block {
	sections := make([]slack.Block, 0)
	sections = append(sections, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*Users List*", false, false), nil, nil))
	for _, user := range users {
		fieldSlice := make([]*slack.TextBlockObject, 0)
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*ID:*\n%s", user.ID), false, false))
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Username:*\n%s", user.Username), false, false))
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Full Name:*\n%s", user.FullName), false, false))
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Designation:*\n%s", user.Designation), false, false))
		sections = append(sections, slack.NewSectionBlock(nil, fieldSlice, nil))
	}
	return sections
}

func (c *listUserCommand) getLocalUsers(message msg.Message) {
	var users []db.User
	if err := c.DB.Debug().Model(&db.User{}).Find(&users).Error; err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, couldn't retrieve db information."),
		)
		return
	}

	c.SlackClient.AddReaction("✅", message)
	c.SlackClient.SendBlockMessage(message, c.formatUsers(users))
}

func (c *listUserCommand) getServerUsers(message msg.Message) {
	res, err := c.client.GetUsers(context.TODO(), message.DBUser.AccessToken)
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while getting users from server."),
		)
		return
	}
	c.SlackClient.AddReaction("✅", message)

	sections := make([]slack.Block, 0)
	sections = append(sections, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*Users List*", false, false), nil, nil))
	for _, user := range res {
		fieldSlice := make([]*slack.TextBlockObject, 0)
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*ID:*\n%s", user.SlackID), false, false))
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Username:*\n%s", user.Username), false, false))
		name := strings.Trim(user.FirstName, " ")
		if user.MiddleName != "" {
			name += fmt.Sprintf(" %s", user.MiddleName)
		}
		if user.LastName != "" {
			name += fmt.Sprintf(" %s", user.LastName)
		}
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Full Name:*\n%s", name), false, false))
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Designation:*\n%s", user.Designation), false, false))
		sections = append(sections, slack.NewSectionBlock(nil, fieldSlice, nil))
	}

	c.SlackClient.SendBlockMessage(message, sections)
}

func (c *listUserCommand) run(match matcher.Result, message msg.Message) {
	switch match.GetString("action") {
	case "local":
		c.getLocalUsers(message)
	case "server":
		c.getServerUsers(message)
	}
}

func (c *listUserCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "list users <local|server>",
			Description: "list users",
			Examples: []string{
				"list users local",
				"list users server",
			},
			Category: category,
		},
	}
}
