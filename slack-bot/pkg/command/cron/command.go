package cron

import (
	"strings"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/util"
	"github.com/Shivin01/Edith/slack-bot/pkg/client"
	"github.com/Shivin01/Edith/slack-bot/pkg/config"
	cronLib "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

// NewCronCommand registers cron which are configurable in the yaml config
func NewCronCommand(base bot.BaseCommand, crons []config.Cron) bot.Command {
	if len(crons) == 0 {
		return nil
	}

	cron := cronLib.New()
	cmd := &command{base, crons, cron}

	for _, cronCommand := range crons {
		_, err := cron.AddFunc(cronCommand.Schedule, cmd.getCallback(cronCommand))
		if err != nil {
			log.Error(err)
		}
	}

	cron.Start()
	log.Infof("Initialized %d crons", len(crons))

	return cmd
}

type command struct {
	bot.BaseCommand
	cfg  []config.Cron
	cron *cronLib.Cron
}

func (c *command) IsEnabled() bool {
	return len(c.cfg) > 0
}

func (c *command) getCallback(cron config.Cron) func() {
	return func() {
		for _, commandTemplate := range cron.Commands {
			command, err := util.CompileTemplate(commandTemplate)
			if err != nil {
				log.Error(err)
				continue
			}
			text, err := util.EvalTemplate(command, util.Parameters{})
			if err != nil {
				log.Error(err)
				continue
			}

			for _, line := range strings.Split(text, "\n") {
				newMessage := msg.Message{}
				newMessage.User = "cron"
				newMessage.Channel, _ = client.GetChannelIDAndName(cron.Channel)
				newMessage.Text = line
				newMessage.AdminMessage = cron.IsAdmin
				client.HandleMessage(newMessage)
			}
		}
	}
}
