package edith

import (
	"fmt"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/db"
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
	return matcher.NewPrivateMatcher(
		c.SlackClient,
		matcher.NewManagerMatcher(
			c.SlackClient,
			matcher.NewRegexpMatcher(`modify user (?P<designation>hr|admin|dev|manager)`, c.run),
		),
	)
}

func (c *modifyUserCommand) run(match matcher.Result, message msg.Message) {
	designation := match.GetString("designation")
	//request := &edith.AddUserRequest{
	//	Username:    user.Profile.DisplayName,
	//	FirstName:   user.Profile.FirstName,
	//	LastName:    user.Profile.LastName,
	//	Password:    "admin@123",
	//	PhoneNumber: user.Profile.Phone,
	//	Skills:      []string{},
	//	SlackID:     user.ID,
	//	Email:       fmt.Sprintf("%s@gmail.com", user.Profile.DisplayName),
	//}
	//
	//response, err := c.client.AddUser(context.TODO(), request)
	//if err != nil {
	//	c.SlackClient.AddReaction("❌", message)
	//	c.SlackClient.ReplyError(
	//		message,
	//		errors.New("sorry, error while creating user"),
	//	)
	//	return
	//}

	if err := c.DB.Debug().Model(&db.User{}).Where("ID = ?", message.DBUser.ID).Update("designation", designation).Error; err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while saving user to database."),
		)
		return
	}
	c.SlackClient.AddReaction("✅", message)
	c.SlackClient.SendMessage(message, fmt.Sprintf("updated user %s", message.DBUser.Username))
}

func (c *modifyUserCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "modify user @username hr|sales|dev|admin",
			Description: "modify existing user",
			Examples: []string{
				"modify user @jarvis dev",
				"modify user @jarvis sales",
				"modify user @jarvis hr",
			},
			Category: category,
		},
	}
}
