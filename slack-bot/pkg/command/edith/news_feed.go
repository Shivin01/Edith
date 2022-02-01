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
		false,
	)
}

func (c *newsFeedCommand) run(match matcher.Result, message msg.Message) {
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
		c.SlackClient.SendMessage(message, "there are no newsfeeds at the moment, I will notify you if there any.")
		return
	}

	var fields [][]string
	for _, user := range newsFeeds {
		name := strings.Trim(user.FirstName, " ")
		if user.MiddleName != "" {
			name += fmt.Sprintf(" %s", user.MiddleName)
		}
		if user.LastName != "" {
			name += fmt.Sprintf(" %s", user.LastName)
		}

		fields = append(fields, []string{
			fmt.Sprintf("_Welcome to our new joiner_ :tada:"),
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Joining Date:*\t%s", user.JoiningDate),
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Full Name:*\t%s", name),
			fmt.Sprintf("\t"),
			fmt.Sprintf("*Designation:*\t%s", user.Designation),
		})
	}

	headerSection := client.GetTextBlock("*Feed*")
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
