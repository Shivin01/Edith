package edith

import (
	"context"
	"fmt"
	"strings"

	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client"
	"github.com/immanoj16/edith/pkg/db"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

// NewCelebrationsCommand is able to send a message to any user/channel
func NewCelebrationsCommand(command edithCommand) bot.Command {
	return &celebrationsCommand{command}
}

type celebrationsCommand struct {
	edithCommand
}

func (c *celebrationsCommand) GetMatcher() matcher.Matcher {
	return matcher.NewAuthorizedMatcher(
		c.SlackClient,
		matcher.NewTextMatcher("celebrations", c.run),
		false,
	)
}

func (c *celebrationsCommand) run(match matcher.Result, message msg.Message) {

	if message.DBUser == nil {
		user := &db.User{}
		if err := c.DB.Debug().Model(&db.User{}).First(&user).Error; err != nil {
			c.SlackClient.AddReaction("❌", message)
			c.SlackClient.ReplyError(
				message,
				errors.New("sorry, error while getting user token to get celebrations."),
			)
			return
		}
		message.DBUser = user
	}

	celebrations, err := c.client.GetCelebrations(context.TODO(), message.DBUser.AccessToken)
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while getting news feeds."),
		)
		return
	}

	if len(celebrations) == 0 {
		c.SlackClient.AddReaction("✅", message)
		c.SlackClient.SendMessage(message, "there are no celebrations at the moment, I will notify you if there any, thank you.")
		return
	}

	c.SlackClient.AddReaction("✅", message)

	var fields [][]string
	for _, user := range celebrations {
		name := strings.Trim(user.FirstName, " ")
		if user.MiddleName != "" {
			name += fmt.Sprintf(" %s", user.MiddleName)
		}
		if user.LastName != "" {
			name += fmt.Sprintf(" %s", user.LastName)
		}

		var field []string
		switch user.Type {
		case "anniversary":
			field = append(field, fmt.Sprintf("_Happy anniversary!_ :tada:"))
			break
		case "birthday":
			field = append(field, fmt.Sprintf("_Wishing you happy birthday_ :tada:"))
			break
		}

		field = append(field, []string{
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Full Name:*\t%s", name),
			fmt.Sprintf("\t"),
		}...)

		fields = append(fields, field)
	}

	headerSection := client.GetTextBlock("*Today's celebrations*")
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

func (c *celebrationsCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "celebrations",
			Description: "show current celebrations if any",
			Examples: []string{
				"celebrations",
			},
			Category: category,
		},
	}
}
