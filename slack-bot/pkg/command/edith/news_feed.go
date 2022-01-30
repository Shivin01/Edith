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

// NewNewsFeedCommand is able to send a message to any user/channel
func NewNewsFeedCommand(command edithCommand) bot.Command {
	return &newsFeedCommand{command}
}

type newsFeedCommand struct {
	edithCommand
}

func (c *newsFeedCommand) GetMatcher() matcher.Matcher {
	return matcher.NewAuthorizedMatcher(
		c.SlackClient,
		matcher.NewTextMatcher("newsfeed", c.run),
	)
}

func (c *newsFeedCommand) run(match matcher.Result, message msg.Message) {
	newsFeeds, err := c.client.GetNewsFeeds(context.TODO(), message.DBUser.AccessToken)
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while getting news feeds."),
		)
		return
	}

	if len(newsFeeds) == 0 {
		c.SlackClient.AddReaction("✅", message)
		c.SlackClient.SendMessage(message, "there are no newsfeeds at the moment, I will notify you if there any, thank you.")
		return
	}

	c.SlackClient.AddReaction("✅", message)
	sections := make([]slack.Block, 0)
	sections = append(sections, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*News Feed*", false, false), nil, nil))
	for _, user := range newsFeeds {
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
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Designation:*\n%s", user.Designation), false, false))
		sections = append(sections, slack.NewSectionBlock(nil, fieldSlice, nil))
	}

	c.SlackClient.SendBlockMessage(message, sections)
}

func (c *newsFeedCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "newsfeed",
			Description: "show current news feed",
			Examples: []string{
				"newsfeed",
			},
			Category: category,
		},
	}
}
