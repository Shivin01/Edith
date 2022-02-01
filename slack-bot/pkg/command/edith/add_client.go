package edith

import (
	"context"
	"fmt"

	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client/edith"
	"github.com/immanoj16/edith/pkg/db"
	"github.com/pkg/errors"
)

// newPingCommand just prints a PING with the needed time from client->slack->edith server
func newAddClientCommand(command edithCommand) bot.Command {
	return &addClientCommand{
		command,
	}
}

type addClientCommand struct {
	edithCommand
}

func (c *addClientCommand) GetMatcher() matcher.Matcher {
	return matcher.NewManagerMatcher(
		c.SlackClient,
		matcher.NewRegexpMatcher(`add client (?P<name>[\w\-_\\/]+) (?P<leave>\d+) (?P<notice_period_count>\d+)`, c.run),
		true,
	)
}

func (c *addClientCommand) run(match matcher.Result, message msg.Message) {
	request := edith.AddClientRequest{
		Name:              match.GetString("name"),
		RegisteredName:    match.GetString("name"),
		LeaveCount:        match.GetInt("leave"),
		NoticePeriodCount: match.GetInt("notice_period_count"),
	}

	res, err := c.client.AddClient(context.TODO(), request)
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while creating user"),
		)
		return
	}

	dbClient := db.Client{
		Name:              res.Name,
		RegisteredName:    res.RegisteredName,
		LeaveCount:        res.LeaveCount,
		NoticePeriodCount: res.NoticePeriodCount,
		ServerID:          res.ID,
	}

	if err := c.DB.Debug().Model(&db.Client{}).Create(&dbClient).Error; err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while saving user to database."),
		)
		return
	}

	c.SlackClient.AddReaction("✅", message)
	c.SlackClient.SendMessage(message, fmt.Sprintf("created client %s", request.Name))
}

func (c *addClientCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "add client <name> <leave_count> <notice_period_count>",
			Description: "add client to the server",
			Examples: []string{
				"add client client_name 20 60",
			},
			Category: category,
		},
	}
}
