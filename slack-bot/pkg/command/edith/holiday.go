package edith

import (
	"context"
	"fmt"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/pkg/errors"
	"github.com/slack-go/slack"
)

// NewHolidayCommand is able to send a message to any user/channel
func NewHolidayCommand(command edithCommand) bot.Command {
	return &holidayCommand{command}
}

type holidayCommand struct {
	edithCommand
}

func (c *holidayCommand) GetMatcher() matcher.Matcher {
	return matcher.NewAuthorizedMatcher(
		c.SlackClient,
		matcher.NewTextMatcher("list holiday", c.run),
	)
}

func (c *holidayCommand) run(match matcher.Result, message msg.Message) {
	holidayList, err := c.client.ListHoliday(context.TODO(), message.DBUser.AccessToken)
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while getting holiday list."),
		)
		return
	}

	c.SlackClient.AddReaction("✅", message)

	sections := make([]slack.Block, 0)
	sections = append(sections, slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*Holiday List*", false, false), nil, nil))
	for _, holiday := range holidayList {
		fieldSlice := make([]*slack.TextBlockObject, 0)
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Date:*\n%s", holiday.Date), false, false))
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Name:*\n%s", holiday.Name), false, false))
		fieldSlice = append(fieldSlice, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Description:*\n%s", holiday.Description), false, false))
		sections = append(sections, slack.NewSectionBlock(nil, fieldSlice, nil))
	}

	c.SlackClient.SendBlockMessage(message, sections)
}

func (c *holidayCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "list holiday",
			Description: "Show holiday list",
			Examples: []string{
				"list holiday",
			},
			Category: category,
		},
	}
}
