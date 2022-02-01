package edith

import (
	"context"
	"fmt"
	"strings"

	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client"
	"github.com/immanoj16/edith/pkg/client/edith"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
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
	return matcher.NewAuthorizedMatcher(
		c.SlackClient,
		matcher.NewRegexpMatcher(`show info <@(?P<user>[\w\-_\\/]+)>`, c.run),
		true,
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

	name := strings.Trim(user.FirstName, " ")
	if user.MiddleName != "" {
		name += fmt.Sprintf(" %s", user.MiddleName)
	}
	if user.LastName != "" {
		name += fmt.Sprintf(" %s", user.LastName)
	}

	var fields [][]string
	fields = append(fields, []string{
		fmt.Sprintf("*ID:*\t%s", user.SlackID),
		fmt.Sprintf("\t"),
		fmt.Sprintf("*Username:*\t%s", user.Username),
		fmt.Sprintf("\t"),
		fmt.Sprintf("*Full Name:*\t%s", name),
		fmt.Sprintf("\t"),
		fmt.Sprintf("*Designation:*\t%s", user.Designation),
	})
	headerSection := client.GetTextBlock("*User Information*")
	blocks := make([]slack.Block, 0, len(fields)+1)
	blocks = append(blocks, headerSection)
	for _, elements := range fields {
		textBlocks := make([]*slack.TextBlockObject, 0, len(elements))
		for _, element := range elements {
			textBlocks = append(textBlocks, slack.NewTextBlockObject("mrkdwn", element, false, false))
		}
		blocks = append(blocks, slack.NewSectionBlock(nil, textBlocks, nil))
	}

	c.SendBlockMessage(message, blocks)
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
