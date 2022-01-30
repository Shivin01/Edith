package bot

import (
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/bot/util"
	"github.com/immanoj16/edith/pkg/client"
	"github.com/immanoj16/edith/pkg/db"
	"sync"
)

var lock sync.Mutex

// Command is the main command struct which needs to provide the matcher and the actual executed action
type Command interface {
	GetMatcher() matcher.Matcher
}

// BaseCommand is base struct which is handy for all commands, as a direct Slack communication is possible
type BaseCommand struct {
	client.SlackClient
	DB *db.DB
}

// Conditional checks if the command should be activated. E.g. if dependencies are not present or it's disabled in the config
type Conditional interface {
	IsEnabled() bool
}

// HelpProvider can be provided by a command to add information within "help" command
type HelpProvider interface {
	// GetHelp each command should provide information, like a description or examples
	GetHelp() []Help
}

// Commands is a wrapper of a list of commands. Only the first matched command will be executed
type Commands struct {
	commands []Command
	matcher  []matcher.Matcher // precompiled matcher objects
	compiled bool
	DB       *db.DB
}

// GetHelp returns the help for ALL included commands
func (c *Commands) GetHelp() []Help {
	help := make([]Help, 0)

	for _, command := range c.commands {
		if helpCommand, ok := command.(HelpProvider); ok {
			help = append(help, helpCommand.GetHelp()...)
		}
	}

	return help
}

func (c *Commands) Run(message msg.Message) bool {
	c.compile()

	for _, command := range c.matcher {
		run, match := command.Match(message)
		if match != nil {
			// this is needed for ConditionMatcher: runner gets already executed in the matcher itself!
			if run != nil {
				run(match, message)
			}

			// only the first command is executed -> abort here
			return true
		}
	}

	return false
}

// AddCommand registers a command to the command list
func (c *Commands) AddCommand(commands ...Command) {
	for _, command := range commands {
		if command == nil {
			continue
		}

		if condition, ok := command.(Conditional); ok {
			if !condition.IsEnabled() {
				// command is disabled!
				continue
			}
		}

		// register i
		if provider, ok := command.(util.TemplateFunctionProvider); ok {
			util.RegisterFunctions(provider.GetTemplateFunction(c.DB))
		}

		c.commands = append(c.commands, command)
	}

	c.compiled = false
}

// Merge two list of commands
func (c *Commands) Merge(commands Commands) {
	c.AddCommand(commands.commands...)
}

// Count the registered/valid commands
func (c *Commands) Count() int {
	return len(c.commands)
}

func (c *Commands) compile() {
	if c.compiled {
		return
	}

	// make sure only one process is creating the compiled list
	lock.Lock()
	defer lock.Unlock()

	if !c.compiled {
		c.matcher = make([]matcher.Matcher, len(c.commands))
		for i, command := range c.commands {
			c.matcher[i] = command.GetMatcher()
		}
		c.compiled = true
	}
}