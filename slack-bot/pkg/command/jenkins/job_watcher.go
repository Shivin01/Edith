package jenkins

import (
	"context"
	"fmt"
	"github.com/immanoj16/edith/pkg/bot"
	matcher2 "github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client/jenkins"
)

const (
	actionWatch   = "watch"
	actionUnwatch = "unwatch"
)

// newJobWatcherCommand initialize a new command to watch for any jenkins job
func newJobWatcherCommand(base jenkinsCommand) bot.Command {
	return &watcherCommand{
		base,
		make(map[string]chan bool),
	}
}

type watcherCommand struct {
	jenkinsCommand
	stopper map[string]chan bool
}

func (c *watcherCommand) GetMatcher() matcher2.Matcher {
	return matcher2.NewRegexpMatcher(`(?P<action>watch|unwatch) job (?P<job>[\w\-_\\/]+)`, c.run)
}

func (c *watcherCommand) run(match matcher2.Result, message msg.Message) {
	action := match.GetString("action")
	jobName := match.GetString("job")
	if action == actionWatch {
		stop := make(chan bool, 1)
		ctx := context.TODO()
		// todo use context.WithCancel instead of stopper chan
		c.stopper[jobName+message.GetUser()] = stop
		builds, err := jenkins.WatchJob(ctx, c.jenkins, jobName, stop)
		if err != nil {
			c.ReplyError(message, err)
			return
		}

		c.SendMessage(message, fmt.Sprintf("Okay, I'll watch %s\nUnwatch via `unwatch %s`", jobName, jobName))

		go func() {
			for build := range builds {
				c.SendMessage(message, fmt.Sprintf(
					"*%s*: %s #%d: %s",
					build.GetResult(),
					jobName,
					build.GetBuildNumber(),
					build.GetUrl()),
				)
			}
		}()
	}

	if action == actionUnwatch {
		if stop, ok := c.stopper[jobName+message.User]; ok {
			stop <- true
		}

		c.SendMessage(message, fmt.Sprintf("Okay, you just unwatched %s", jobName))
	}
}

func (c *watcherCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "watch job <job>",
			Description: "watch jenkins jobs and informs about successful/error jobs",
			Examples: []string{
				"watch job MyJobName",
				"unwatch job MyJobName",
			},
			Category: category,
		},
		{
			Command:     "unwatch job <job>",
			Description: "unwatch jenkins jobs and informs about successful/error jobs",
			Examples: []string{
				"watch job MyJobName",
				"unwatch job MyJobName",
			},
			Category: category,
		},
	}
}
