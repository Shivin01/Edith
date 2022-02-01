package custom

import (
	"fmt"

	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/db"
	"github.com/pkg/errors"
)

func (c command) delete(match matcher.Result, message msg.Message) {
	alias := match.GetString("alias")

	if err := c.DB.Debug().Model(&db.CustomCommand{}).Where("user_refer = ?", message.DBUser.ID).Where("alias = ?", alias).Delete(&db.CustomCommand{}).Error; err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while delete custom command"),
		)
		return
	}

	var customCommands []db.CustomCommand
	if err := c.DB.Debug().Model(&db.CustomCommand{}).Where("user_refer = ?", message.DBUser.ID).Find(&customCommands).Error; err != nil {
		c.SlackClient.AddReaction("❌", message)
		c.SlackClient.ReplyError(
			message,
			errors.New("sorry, error while gettting custom commands"),
		)
		return
	}

	message.DBUser.CustomCommands = customCommands

	c.SendMessage(message, fmt.Sprintf("Okay, I deleted command: `%s`", alias))
}
