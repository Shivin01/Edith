package custom

import (
	"fmt"

	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/db"
	"github.com/pkg/errors"
)

func (c *command) add(match matcher.Result, message msg.Message) {
	alias := match.GetString("alias")
	command := match.GetString("command")

	found := false
	for _, customCommand := range message.DBUser.CustomCommands {
		if customCommand.Alias == alias {
			found = true
		}
	}

	if !found {
		customCommand := db.CustomCommand{
			Alias:     alias,
			Command:   command,
			UserRefer: message.DBUser.ID,
		}

		if err := c.DB.Debug().Model(&db.CustomCommand{}).Create(&customCommand).Error; err != nil {
			c.SlackClient.AddReaction("❌", message)
			c.SlackClient.ReplyError(
				message,
				errors.New("sorry, error while saving custom command to database."),
			)
			return
		}

		message.DBUser.CustomCommands = append(message.DBUser.CustomCommands, customCommand)
	}

	c.SendMessage(
		message,
		fmt.Sprintf("Added command: `%s`. Just use `%s` in future.", command, alias),
	)
}
