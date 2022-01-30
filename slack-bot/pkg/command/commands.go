package command

import (
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/client"
	"github.com/immanoj16/edith/pkg/command/admin"
	"github.com/immanoj16/edith/pkg/command/cron"
	"github.com/immanoj16/edith/pkg/command/custom"
	"github.com/immanoj16/edith/pkg/command/edith"
	"github.com/immanoj16/edith/pkg/command/jenkins"
	"github.com/immanoj16/edith/pkg/command/queue"
	"github.com/immanoj16/edith/pkg/command/weather"
	"github.com/immanoj16/edith/pkg/config"
	"github.com/immanoj16/edith/pkg/db"
)

// GetCommands returns the list of default command which are available
func GetCommands(slackClient client.SlackClient, cfg config.Config, db *db.DB) *bot.Commands {
	base := bot.BaseCommand{SlackClient: slackClient, DB: db}

	commands := &bot.Commands{DB: db}

	commands.AddCommand(
		NewSendMessageCommand(base),
		NewHelpCommand(base, commands),
		weather.NewWeatherCommand(base, cfg.OpenWeather),
		cron.NewCronCommand(base, cfg.Crons),
		custom.GetCommand(base),
		queue.NewQueueCommand(base),
		queue.NewListCommand(base),
		cron.NewCronCommand(base, cfg.Crons),
	)

	commands.Merge(admin.GetCommands(base, &cfg))
	commands.Merge(edith.GetCommands(base, &cfg.Server))

	// jenkins
	commands.Merge(jenkins.GetCommands(cfg.Jenkins, base))

	return commands
}
