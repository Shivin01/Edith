package edith

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client/edith"
	"github.com/immanoj16/edith/pkg/db"
	"github.com/pkg/errors"
)

// newPingCommand just prints a PING with the needed time from client->slack->edith server
func newAddUserCommand(command edithCommand) bot.Command {
	return &addUserCommand{
		command,
	}
}

type addUserCommand struct {
	edithCommand
}

func (c *addUserCommand) GetMatcher() matcher.Matcher {
	return matcher.NewManagerMatcher(
		c.SlackClient,
		matcher.NewRegexpMatcher(`add user <@(?P<user>[\w\-_\\/]+)> (?P<designation>hr|admin|dev|manager) (?P<gender>male|female|others) (?P<client>[\w\-_\\/]+)`, c.run),
		true,
	)
}

func (c *addUserCommand) run(match matcher.Result, message msg.Message) {
	username := match.GetString("user")
	designation := match.GetString("designation")
	clientName := match.GetString("client")
	gender := match.GetString("gender")
	user := c.SlackClient.GetUserDetails(username, message)

	spew.Dump(user)

	const password = "admin@123"
	request := edith.AddUserRequest{
		Username:    user.Profile.DisplayName,
		FirstName:   user.Profile.FirstName,
		LastName:    user.Profile.LastName,
		Password:    password,
		Password1:   password,
		Password2:   password,
		PhoneNumber: user.Profile.Phone,
		Skills:      []string{},
		SlackID:     user.ID,
		Email:       fmt.Sprintf("%s@gmail.com", user.Profile.DisplayName),
		Designation: designation,
		JoiningDate: time.Now().Format("2006-01-02"),
		ClientName:  clientName,
		Gender:      strings.ToUpper(gender),
	}

	response, err := c.client.AddUser(context.TODO(), request)
	if err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while creating user"),
		)
		return
	}

	dbUser := db.User{
		ID:          user.ID,
		Username:    response.User.Username,
		FullName:    user.RealName,
		AccessToken: response.Token,
		Designation: designation,
		ServerID:    response.User.Pk,
		ClientName:  clientName,
		Gender:      strings.ToUpper(gender),
	}

	if err := c.DB.Debug().Model(&db.User{}).Create(&dbUser).Error; err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while saving user to database."),
		)
		return
	}
	c.SlackClient.AddReaction("✅", message)
	c.SlackClient.SendMessage(message, fmt.Sprintf("created user %s, please change password.", dbUser.Username))
}

func (c *addUserCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "add user @username <hr|sales|dev|admin> <gender> <client>",
			Description: "add user to root permissions",
			Examples: []string{
				"add user @jarvis dev female client",
				"add user @jarvis sales female client",
				"add user @jarvis hr male client",
				"add user @jarvis hr male client",
			},
			Category: category,
		},
	}
}
