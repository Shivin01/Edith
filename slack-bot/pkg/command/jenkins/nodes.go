package jenkins

import (
	"context"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/immanoj16/edith/pkg/bot"
	matcher2 "github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/config"
	"sort"
)

const (
	iconStatusOnline  = "✔"
	iconStatusOffline = "🔴"
)

type nodesCommand struct {
	jenkinsCommand
	cfg config.Jenkins
}

// newNodesCommand lists all Jenkins nodes/slaves and the current number of running executors
func newNodesCommand(base jenkinsCommand, cfg config.Jenkins) bot.Command {
	return &nodesCommand{base, cfg}
}

func (c *nodesCommand) GetMatcher() matcher2.Matcher {
	return matcher2.NewTextMatcher("jenkins nodes", c.run)
}

func (c *nodesCommand) IsEnabled() bool {
	return c.jenkins != nil
}

func (c *nodesCommand) run(match matcher2.Result, message msg.Message) {
	ctx := context.TODO()
	nodes, err := c.jenkins.GetAllNodes(ctx)
	if err != nil {
		c.ReplyError(message, err)
		return
	}

	// sort nodes by name
	sort.Slice(nodes, func(i, j int) bool {
		return nodes[i].GetName() < nodes[j].GetName()
	})

	text := fmt.Sprintf("*<%s/computer/|%d Nodes>*\n", c.cfg.Host, len(nodes))
	var statusIcon string
	for _, node := range nodes {
		offline := node.Raw.Offline

		if offline {
			statusIcon = iconStatusOffline
		} else {
			statusIcon = iconStatusOnline
		}

		text += fmt.Sprintf(
			"• *<%s/computer/%s/|%s>* - status: %s - busy executors: %d/%d\n",
			c.cfg.Host,
			node.GetName(),
			node.GetName(),
			statusIcon,
			c.countBusyExecutors(node),
			len(node.Raw.Executors),
		)
	}

	c.SendMessage(message, text)
}

func (c *nodesCommand) countBusyExecutors(node *gojenkins.Node) int {
	busyNodes := 0
	for _, executor := range node.Raw.Executors {
		if executor.CurrentExecutable.Number != 0 {
			busyNodes++
		}
	}

	return busyNodes
}

func (c *nodesCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "jenkins nodes",
			Description: "Prints a list of all jenkins nodes",
			Examples: []string{
				"jenkins nodes",
			},
			Category: category,
		},
	}
}