package edith

import (
	"context"
	"fmt"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/matcher"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
	"github.com/Shivin01/Edith/slack-bot/pkg/client"
	"github.com/Shivin01/Edith/slack-bot/pkg/db"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

// newPingCommand just prints a PING with the needed time from client->slack->edith server
func newListClientCommand(command edithCommand) bot.Command {
	return &listClientCommand{
		command,
	}
}

type listClientCommand struct {
	edithCommand
}

func (c *listClientCommand) GetMatcher() matcher.Matcher {
	return matcher.NewManagerMatcher(
		c.SlackClient,
		matcher.NewGroupMatcher(
			matcher.NewTextMatcher("list clients local", c.getLocalClients),
			matcher.NewTextMatcher("list clients server", c.getServerClients),
		),
		true,
	)
}

func (c listClientCommand) format(clients []db.Client) []slack.Block {
	var fields [][]string
	for _, c := range clients {
		fields = append(fields, []string{
			fmt.Sprintf("*ID:*\t%d", c.ID),
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Name:*\t%s", c.Name),
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Registered Name:*\t%s", c.RegisteredName),
		})
	}

	headerSection := client.GetTextBlock("*Clients*")
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

func (c *listClientCommand) getLocalClients(match matcher.Result, message msg.Message) {
	var clients []db.Client
	if err := c.DB.Debug().Model(&db.Client{}).Find(&clients).Error; err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, couldn't retrieve db information."),
		)
		return
	}

	c.SlackClient.AddReaction("✅", message)
	c.SendBlockMessage(message, c.format(clients))
}

func (c *listClientCommand) getServerClients(match matcher.Result, message msg.Message) {
	res, err := c.client.GetClients(context.TODO())
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
	for _, c := range res {
		fields = append(fields, []string{
			fmt.Sprintf("*ID:*\t%d", c.ID),
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Name:*\t%s", c.Name),
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Registered Name:*\t%s", c.RegisteredName),
		})
	}

	headerSection := client.GetTextBlock("*Clients*")
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

func (c *listClientCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "list clients <local|server>",
			Description: "list clients",
			Examples: []string{
				"list clients local",
				"list clients server",
			},
			Category: category,
		},
	}
}
