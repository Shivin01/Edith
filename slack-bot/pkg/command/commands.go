package command

import (
	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/client"
	"github.com/Shivin01/Edith/slack-bot/pkg/command/admin"
	"github.com/Shivin01/Edith/slack-bot/pkg/command/cron"
	"github.com/Shivin01/Edith/slack-bot/pkg/command/custom"
	"github.com/Shivin01/Edith/slack-bot/pkg/command/edith"
	"github.com/Shivin01/Edith/slack-bot/pkg/command/jenkins"
	"github.com/Shivin01/Edith/slack-bot/pkg/command/queue"
	"github.com/Shivin01/Edith/slack-bot/pkg/command/weather"
	"github.com/Shivin01/Edith/slack-bot/pkg/config"
	"github.com/Shivin01/Edith/slack-bot/pkg/db"
)

// GetCommands returns the list of default command which are available
func GetCommands(slackClient client.SlackClient, cfg config.Config, db *db.DB) *bot.Commands {
	base := bot.BaseCommand{SlackClient: slackClient, DB: db}

	commands := &bot.Commands{DB: db}

	commands.AddCommand(
		NewSendMessageCommand(base),
		NewHelpCommand(base, commands),
		weather.NewWeatherCommand(base, cfg.OpenWeather),
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
