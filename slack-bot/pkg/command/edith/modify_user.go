package edith

import (
	"context"
	"fmt"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/matcher"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
	"github.com/Shivin01/Edith/slack-bot/pkg/db"
	"github.com/pkg/errors"
)

// newPingCommand just prints a PING with the needed time from client->slack->edith server
func newModifyUserCommand(command edithCommand) bot.Command {
	return &modifyUserCommand{
		command,
	}
}

type modifyUserCommand struct {
	edithCommand
}

func (c *modifyUserCommand) GetMatcher() matcher.Matcher {
	return matcher.NewManagerMatcher(
		c.SlackClient,
		matcher.NewGroupMatcher(
			matcher.NewRegexpMatcher(`modify user <@(?P<user>[\w\-_\\/]+)> designation (?P<designation>hr|admin|dev|manager)`, c.modifyDesignation),
			matcher.NewRegexpMatcher(`modify user <@(?P<user>[\w\-_\\/]+)> birthdate (?P<date>\d{4}-\d{2}-\d{2})`, c.modifyBirthDate),
		),
		true,
	)
}

func (c *modifyUserCommand) modifyClient(match matcher.Result, message msg.Message) {
	user := match.GetString("user")
	clientName := match.GetString("name")

	u := &db.User{}
	if err := c.DB.Debug().Model(&db.User{}).Where("ID = ?", user).First(u).Error; err != nil {
		message.DBUser = nil
	}

	err := c.client.ModifyUser(context.TODO(), u.ServerID, map[string]interface{}{
		"client_name": clientName,
	}, message.DBUser.AccessToken)
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while creating user"),
		)
		return
	}

	if err := c.DB.Debug().Model(&db.User{}).Where("ID = ?", user).Update("client_name", clientName).Error; err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while updating user to database."),
		)
		return
	}
	c.SlackClient.AddReaction("✅", message)
	c.SlackClient.SendMessage(message, fmt.Sprintf("updated user %s", user))
}

func (c *modifyUserCommand) modifyBirthDate(match matcher.Result, message msg.Message) {
	user := match.GetString("user")
	date := match.GetString("date")

	u := &db.User{}
	if err := c.DB.Debug().Model(&db.User{}).Where("ID = ?", user).First(u).Error; err != nil {
		message.DBUser = nil
	}

	err := c.client.ModifyUser(context.TODO(), u.ServerID, map[string]interface{}{
		"birth_date": date,
	}, message.DBUser.AccessToken)
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while creating user"),
		)
		return
	}

	c.SlackClient.AddReaction("✅", message)
	c.SlackClient.SendMessage(message, fmt.Sprintf("updated user %s", u.Username))
}

func (c *modifyUserCommand) modifyDesignation(match matcher.Result, message msg.Message) {
	user := match.GetString("user")
	designation := match.GetString("designation")

	u := &db.User{}
	if err := c.DB.Debug().Model(&db.User{}).Where("ID = ?", user).First(u).Error; err != nil {
		message.DBUser = nil
	}

	err := c.client.ModifyUser(context.TODO(), u.ServerID, map[string]interface{}{
		"designation": designation,
	}, message.DBUser.AccessToken)
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while creating user"),
		)
		return
	}

	if err := c.DB.Debug().Model(&db.User{}).Where("ID = ?", user).Update("designation", designation).Error; err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while updating user to database."),
		)
		return
	}
	c.SlackClient.AddReaction("✅", message)
	c.SlackClient.SendMessage(message, fmt.Sprintf("updated user %s", u.GetRealName()))
}

func (c *modifyUserCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "modify user @username designation <hr|sales|dev|admin|manager>",
			Description: "modify existing user",
			Examples: []string{
				"modify user @jarvis designation sales",
			},
			Category: category,
		},
		{
			Command:     "modify user @username birthdate <birthdate>",
			Description: "modify existing user",
			Examples: []string{
				"modify user @jarvis birthdate 1995-04-20",
			},
			Category: category,
		},
	}
}
