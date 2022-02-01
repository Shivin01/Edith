package admin

import (
	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/config"
)

// GetCommands will return a list of available admin commands
func GetCommands(base bot.BaseCommand, cfg *config.Config) bot.Commands {
	var commands bot.Commands

	commands.AddCommand(
		newStatsCommand(base, cfg),
		newBotLogCommand(base, cfg),
		newPingCommand(base),
	)

	return commands
}
