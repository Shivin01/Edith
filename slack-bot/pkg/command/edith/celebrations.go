package edith

import (
	"context"
	"fmt"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
	"strings"
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
	)
}

func (c *celebrationsCommand) run(match matcher.Result, message msg.Message) {
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
	sections := make([]slack.Block, 0)
	sections = append(sections, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*Celebrations*", false, false), nil, nil))
	for _, user := range celebrations {
		fieldSlice := make([]*slack.TextBlockObject, 0)
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*ID:*\n%s", user.ID), false, false))
		name := strings.Trim(user.FirstName, " ")
		if user.MiddleName != "" {
			name += fmt.Sprintf(" %s", user.MiddleName)
		}
		if user.LastName != "" {
			name += fmt.Sprintf(" %s", user.LastName)
		}
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Full Name:*\n%s", name), false, false))
		sections = append(sections, slack.NewSectionBlock(nil, fieldSlice, nil))
	}

	c.SlackClient.SendBlockMessage(message, sections)
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
