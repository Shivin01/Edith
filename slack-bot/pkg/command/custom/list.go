package custom

import (
	"fmt"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot/matcher"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
)

func (c command) list(match matcher.Result, message msg.Message) {
	if len(message.DBUser.CustomCommands) == 0 {
		c.SendMessage(message, "No commands define yet. Use `add command 'your alias' 'command to execute'`")
		return
	}

	responseText := fmt.Sprintf("You defined %d commands:", len(message.DBUser.CustomCommands))
	for _, customCommand := range message.DBUser.CustomCommands {
		responseText += fmt.Sprintf("\n - %s: `%s`", customCommand.Alias, customCommand.Command)
	}

	c.SendMessage(message, responseText)
}
