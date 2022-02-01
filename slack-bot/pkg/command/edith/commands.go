package edith

import (
	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/client/edith"
	"github.com/Shivin01/Edith/slack-bot/pkg/config"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// help category to group all Jenkins command
var category = bot.Category{
	Name:        "Edith",
	Description: "Interact with Jenkins jobs: Trigger builds, check job status or watch builds",
	HelpURL:     "http://localhost:8000/swagger/",
}

// base command to access Slack+Jenkins directly
type edithCommand struct {
	bot.BaseCommand
	client edith.Client
}

// GetCommands will return a list of available Jenkins commands...if the config is set!
func GetCommands(base bot.BaseCommand, cfg *config.Server) bot.Commands {
	var commands bot.Commands

	edithClient, err := edith.GetClient(cfg)
	if err != nil {
		log.Error(errors.Wrap(err, "Error while getting edith server client"))
		return commands
	}

	edithBase := edithCommand{
		base,
		edithClient,
	}

	commands.AddCommand(
		// Users command
		newTokenCommand(edithBase),
		newListUserCommand(edithBase),
		newAddUserCommand(edithBase),
		newGetTokenCommand(edithBase),
		newRemoveUserCommand(edithBase),
		newModifyUserCommand(edithBase),
		newEmployeeInfoCommand(edithBase),
		newMarkAttendanceCommand(edithBase),

		// announcements
		NewCelebrationsCommand(edithBase),
		NewNewsFeedCommand(edithBase),
		newMakeAnnouncementCommand(edithBase),

		// Leave commands
		newLeaveRequestCommand(edithBase),
		newLeaveListCommand(edithBase),
		newLeaveApproveCommand(edithBase),

		// Extra commands
		NewHolidayCommand(edithBase),
	)

	return commands
}
