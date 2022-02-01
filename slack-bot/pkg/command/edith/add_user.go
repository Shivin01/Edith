package edith

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/matcher"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
	"github.com/Shivin01/Edith/slack-bot/pkg/client/edith"
	"github.com/Shivin01/Edith/slack-bot/pkg/db"
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
		matcher.NewRegexpMatcher(`add user <@(?P<user>[\w\-_\\/]+)> (?P<designation>hr|sales|dev|admin|manager) (?P<gender>male|female|others)`, c.run),
		true,
	)
}

func (c *addUserCommand) run(match matcher.Result, message msg.Message) {
	username := match.GetString("user")
	designation := match.GetString("designation")
	gender := match.GetString("gender")
	user := c.SlackClient.GetUserDetails(username, message)
	team := c.SlackClient.GetTeamDetails(message)

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
		ClientName:  team.Name,
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
		ClientName:  team.Name,
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
