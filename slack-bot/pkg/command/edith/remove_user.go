package edith

import (
	"context"
	"fmt"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/matcher"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
	"github.com/Shivin01/Edith/slack-bot/pkg/db"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// newRemoveUserCommand just prints a PING with the needed time from client->slack->edith server
func newRemoveUserCommand(command edithCommand) bot.Command {
	return &removeUserCommand{
		command,
	}
}

type removeUserCommand struct {
	edithCommand
}

func (c *removeUserCommand) GetMatcher() matcher.Matcher {
	return matcher.NewManagerMatcher(
		c.SlackClient,
		matcher.NewRegexpMatcher(`remove user <@(?P<user>[\w\-_\\/]+)>`, c.run),
		true,
	)
}

func (c *removeUserCommand) run(match matcher.Result, message msg.Message) {
	username := match.GetString("user")
	user := &db.User{}
	if err := c.DB.Debug().Model(&db.User{}).Where("ID = ?", username).Find(user).Error; err != nil {
		log.Errorf("got error: %s", err.Error())
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, couldn't get db information."),
		)
		return
	}

	err := c.client.DeleteUser(context.TODO(), user.ServerID, message.DBUser.AccessToken)
	if err != nil {
		log.Errorf("got error: %s", err.Error())
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, couldn't delete user."),
		)
		return
	}

	if err := c.DB.Debug().Model(&db.User{}).Where("ID = ?", username).Delete(&db.User{}).Error; err != nil {
		log.Errorf("got error: %s", err.Error())
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, couldn't delete db information."),
		)
		return
	}

	c.SlackClient.AddReaction("✅", message)
	c.SlackClient.SendMessage(message, fmt.Sprintf("successfully removed user %s", user.GetRealName()))
}

func (c *removeUserCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "remove user <@username>",
			Description: "remove user",
			Examples: []string{
				"remove user @abc",
			},
			Category: category,
		},
	}
}
