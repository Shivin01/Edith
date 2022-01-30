package edith

import (
	"context"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/db"
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
	return matcher.NewPrivateMatcher(
		c.SlackClient,
		matcher.NewManagerMatcher(
			c.SlackClient,
			matcher.NewRegexpMatcher(`remove user <@(?P<user>[\w\-_\\/]+)>`, c.run),
		),
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
