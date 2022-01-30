package edith

import (
	"context"
	"fmt"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client/edith"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"strings"
)

// newPingCommand just prints a PING with the needed time from client->slack->edith server
func newEmployeeInfoCommand(command edithCommand) bot.Command {
	return &employeeInfoCommand{
		command,
	}
}

type employeeInfoCommand struct {
	edithCommand
}

func (c *employeeInfoCommand) GetMatcher() matcher.Matcher {
	return matcher.NewPrivateMatcher(
		c.SlackClient,
		matcher.NewAuthorizedMatcher(
			c.SlackClient,
			matcher.NewRegexpMatcher(`show info <@(?P<user>[\w\-_\\/]+)>`, c.run),
		),
	)
}

func (c *employeeInfoCommand) run(match matcher.Result, message msg.Message) {
	username := match.GetString("user")
	var users []*edith.UserResponse
	var err error
	if message.DBUser.IsAdmin() {
		users, err = c.client.GetMinimalUser(context.TODO(), username, message.DBUser.AccessToken)
		if err != nil {
			log.Errorf("got error: %s", err.Error())
			c.SlackClient.AddReaction("❌", message)
			c.SlackClient.ReplyError(
				message,
				errors.New("sorry, error while getting users from server."),
			)
			return
		}
	} else {
		users, err = c.client.GetUser(context.TODO(), username, message.DBUser.AccessToken)
		if err != nil {
			log.Errorf("got error: %s", err.Error())
			c.SlackClient.AddReaction("❌", message)
			c.SlackClient.ReplyError(
				message,
				errors.New("sorry, error while getting users from server."),
			)
			return
		}
	}

	c.SlackClient.AddReaction("✅", message)
	if len(users) == 0 {
		c.SlackClient.SendMessage(message, fmt.Sprintf("couldn't found slack user with %s ID, please create with `add user`", username))
		return
	}

	if len(users) > 1 {
		c.SlackClient.SendMessage(message, fmt.Sprintf("got multiple users with %s ID, please check with server admin", username))
		return
	}

	user := users[0]
	sections := make([]slack.Block, 0)
	sections = append(sections, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*User Information*", false, false), nil, nil))
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

	c.SlackClient.SendBlockMessage(message, sections)
}

func (c *employeeInfoCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "show info <@username>",
			Description: "show info about the user",
			Examples: []string{
				"show info @user",
			},
			Category: category,
		},
	}
}
