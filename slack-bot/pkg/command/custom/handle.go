package custom

import (
	"fmt"
	"strings"

	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client"
	"github.com/immanoj16/edith/pkg/db"
)

// check if the current user has a alias for the current message, if yes, execute the commands
func (c *command) handle(ref msg.Ref, text string) bool {
	var commands string

	var customCommands []db.CustomCommand

	if ref.IsInternalMessage() {
		c.DB.Debug().Model(&db.CustomCommand{}).Find(&customCommands)
	} else {
		user := ref.GetDBUser()
		if user != nil {
			customCommands = ref.GetDBUser().CustomCommands
		}
	}

	for _, customCommand := range customCommands {
		if customCommand.Alias == text {
			commands = customCommand.Command
		}
	}

	if commands == "" {
		return false
	}

	c.SendMessage(ref, fmt.Sprintf("executing command: `%s`", commands))
	for _, command := range strings.Split(commands, ";") {
		message := client.HandleMessageWithDoneHandler(ref.WithText(command))
		message.Wait()
	}

	return true
}
