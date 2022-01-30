package admin

import (
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/config"
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
