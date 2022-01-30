package edith

import (
	"context"
	"fmt"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client"
	"github.com/immanoj16/edith/pkg/client/edith"
	"github.com/pkg/errors"
)

// newMakeAnnouncementCommand used to make announcement
func newMakeAnnouncementCommand(command edithCommand) bot.Command {
	return &makeAnnouncementCommand{
		command,
	}
}

type makeAnnouncementCommand struct {
	edithCommand
}

func (c *makeAnnouncementCommand) GetMatcher() matcher.Matcher {
	return matcher.NewPrivateMatcher(
		c.SlackClient,
		matcher.NewAuthorizedMatcher(
			c.SlackClient,
			matcher.NewRegexpMatcher(`make announcement (?P<type>[\w\-_\\/]+) (?P<details>[\w\-_\\/]+)`, c.run),
		),
	)
}

func (c *makeAnnouncementCommand) run(match matcher.Result, message msg.Message) {
	t := match.GetString("type")
	details := match.GetString("details")

	res, err := c.client.MakeAnnouncement(context.TODO(), message.DBUser.AccessToken, edith.AnnouncementRequest{
		Type:   t,
		Detail: details,
	})
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while make announcements..."),
		)
		return
	}
	c.SlackClient.SendMessage(message, "Successfully added response, users will see in #devops channel.")
	channelID, _ := client.GetChannelIDAndName("#devops")
	c.SlackClient.NewPostMessage(message, channelID, fmt.Sprintf(":tada: *New Announcement by* _%s_ :tada:\n\n*Type:* %s\n*Details:* %s\n", message.DBUser.FullName, res.Type, res.Detail))
}

func (c *makeAnnouncementCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "make announcement <type_in_string> <details_in_string>",
			Description: "make announcement in company wide",
			Examples: []string{
				"make announcement abc def",
			},
			Category: category,
		},
	}
}
