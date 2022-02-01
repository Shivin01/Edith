package queue

import (
	"fmt"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/matcher"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
	"github.com/Shivin01/Edith/slack-bot/pkg/client"
	log "github.com/sirupsen/logrus"
)

// NewQueueCommand is able to execute a command when another blocking process is done
// e.g. have a running jenkins job and using "then reply done!" to get a information later
func NewQueueCommand(base bot.BaseCommand) bot.Command {
	executeFallbackCommand(base.DB)

	return &thenCommand{
		base,
	}
}

type thenCommand struct {
	bot.BaseCommand
}

func (c *thenCommand) GetMatcher() matcher.Matcher {
	return matcher.NewRegexpMatcher("(?i:queue|then) (?P<command>.*)", c.run)
}

func (c *thenCommand) run(match matcher.Result, message msg.Message) {
	runningCommand, found := runningCommands[message.GetUniqueKey()]
	if !found {
		c.ReplyError(
			message,
			fmt.Errorf("you have to call this command when another long running command is already running"),
		)
		return
	}

	command := match.GetString("command")
	c.AddReaction(waitIcon, message)

	go func() {
		runningCommand.Wait()

		c.AddReaction(doneIcon, message)

		// trigger new command
		client.HandleMessage(message.WithText(command))

		log.Infof("[Queue] Blocking command is over, eval message: %s", command)
	}()
}
