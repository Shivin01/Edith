package edith

import (
	"context"
	"fmt"
	"strings"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/matcher"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
	"github.com/Shivin01/Edith/slack-bot/pkg/client"
	"github.com/Shivin01/Edith/slack-bot/pkg/db"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
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
	return matcher.NewManagerMatcher(
		c.SlackClient,
		matcher.NewGroupMatcher(
			matcher.NewTextMatcher("list users local", c.getLocalUsers),
			matcher.NewTextMatcher("list users server", c.getServerUsers),
		),
		true,
	)
}

func (c listUserCommand) formatUsers(users []db.User) []slack.Block {
	var fields [][]string
	for _, user := range users {
		fields = append(fields, []string{
			fmt.Sprintf("*ID:*\t%s", user.ID),
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Username:*\t%s", user.Username),
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Full Name:*\t%s", user.FullName),
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Designation:*\t%s", user.Designation),
		})
	}

	headerSection := client.GetTextBlock("*Users*")
	blocks := make([]slack.Block, 0, len(fields)+1)
	blocks = append(blocks, headerSection)
	for _, elements := range fields {
		textBlocks := make([]*slack.TextBlockObject, 0, len(elements))
		for _, element := range elements {
			textBlocks = append(textBlocks, slack.NewTextBlockObject("mrkdwn", element, false, false))
		}
		blocks = append(blocks, slack.NewSectionBlock(nil, textBlocks, nil))
	}
	return blocks
}

func (c *listUserCommand) getLocalUsers(match matcher.Result, message msg.Message) {
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

func (c *listUserCommand) getServerUsers(match matcher.Result, message msg.Message) {
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

	var fields [][]string
	for _, user := range res {
		name := strings.Trim(user.FirstName, " ")
		if user.MiddleName != "" {
			name += fmt.Sprintf(" %s", user.MiddleName)
		}
		if user.LastName != "" {
			name += fmt.Sprintf(" %s", user.LastName)
		}
		fields = append(fields, []string{
			fmt.Sprintf("*ID:*\t%s", user.SlackID),
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Username:*\t%s", user.Username),
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Full Name:*\t%s", name),
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Designation:*\t%s", user.Designation),
		})
	}

	headerSection := client.GetTextBlock("*Users*")
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
