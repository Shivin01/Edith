package queue

import (
	"fmt"
	"text/template"
	"time"

	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/bot/util"
	"github.com/immanoj16/edith/pkg/client"
	"github.com/immanoj16/edith/pkg/db"
	"github.com/slack-go/slack"
)

const processingReaction = "eyes"

type listCommand struct {
	bot.BaseCommand
}

type filterFunc func(msg.Message) bool

// NewListCommand prints the list of all queued commands (blocking commands like running Jenkins jobs)
func NewListCommand(base bot.BaseCommand) bot.Command {
	return &listCommand{
		base,
	}
}

func (c *listCommand) GetMatcher() matcher.Matcher {
	return matcher.NewAuthorizedMatcher(
		c.SlackClient,
		matcher.NewGroupMatcher(
			matcher.NewTextMatcher("list queue", c.listAll),
			matcher.NewTextMatcher("list queue in channel", c.listChannel),
		),
		false,
	)
}

// ListAll shows a list of all queued commands
func (c *listCommand) listAll(match matcher.Result, message msg.Message) {
	c.sendList(message, func(event msg.Message) bool {
		return true
	})
}

// ListChannel shows a list of all queued commands within the current channel
func (c *listCommand) listChannel(match matcher.Result, message msg.Message) {
	c.sendList(message, func(queuedEvent msg.Message) bool {
		return message.GetChannel() == queuedEvent.GetChannel()
	})
}

// format a block-based message with all matching commands
func (c *listCommand) sendList(message msg.Message, filter filterFunc) {
	// add :eyes: temporary because loading the list might take some seconds
	c.AddReaction(processingReaction, message)
	defer c.RemoveReaction(processingReaction, message)

	count, queueBlocks := c.getQueueAsBlocks(message, filter)

	blocks := []slack.Block{
		client.GetTextBlock(fmt.Sprintf("*%d queued commands*", count)),
	}
	blocks = append(blocks, queueBlocks...)

	// replace the original message when it's triggered by the "refresh" button
	var msgOptions []slack.MsgOption
	if message.IsUpdatedMessage() {
		msgOptions = append(msgOptions, slack.MsgOptionUpdate(message.Timestamp))
	}

	c.SendBlockMessage(message, blocks, msgOptions...)
}

// loads all matching queue entries and format them into slack.Block
func (c *listCommand) getQueueAsBlocks(message msg.Message, filter filterFunc) (count uint, blocks []slack.Block) {
	now := time.Now()
	var keys []db.FallbackQueue
	c.DB.Debug().Model(&db.FallbackQueue{}).Find(&keys)

	for _, key := range keys {
		queuedEvent := msg.Message{
			MessageRef: msg.MessageRef{
				Channel:         key.Channel,
				User:            key.User,
				Timestamp:       key.Timestamp,
				Thread:          key.Thread,
				InternalMessage: key.InternalMessage,
				UpdatedMessage:  key.UpdatedMessage,
			},
			Text: key.Text,
		}
		if !filter(message) {
			continue
		}

		count++
		_, userName := client.GetUserIDAndName(queuedEvent.User)

		messageTime := queuedEvent.GetTime()
		timeAgo := now.Sub(messageTime)
		text := fmt.Sprintf(
			"*%s* (<%s|%s, %s ago>): ```%s``` %s",
			userName,
			client.GetSlackArchiveLink(queuedEvent),
			messageTime.Format(time.Stamp),
			util.FormatDuration(timeAgo),
			queuedEvent.GetText(),
			c.getReactions(queuedEvent),
		)

		textBlock := client.GetTextBlock(text)
		blocks = append(
			blocks,
			textBlock,
		)
	}

	// add "Updated at..." time if there was an update
	if message.IsUpdatedMessage() {
		blocks = append(
			blocks,
			client.GetContextBlock(fmt.Sprintf("Updated at: %s", time.Now().Format(time.Stamp))),
		)
	}

	// add "Refresh" button
	blocks = append(
		blocks,
		slack.NewActionBlock(
			"",
			client.GetInteractionButton("Refresh :arrows_counterclockwise:", message.GetText()),
		),
	)

	return count, blocks
}

func (c *listCommand) getReactions(ref msg.Ref) string {
	formattedReactions := ""
	msgRef := slack.NewRefToMessage(ref.GetChannel(), ref.GetTimestamp())
	reactions, _ := c.GetReactions(msgRef, slack.NewGetReactionsParameters())

	for _, reaction := range reactions {
		formattedReactions += ":" + reaction.Name + ":"
	}
	return formattedReactions
}

func (c *listCommand) GetTemplateFunction(database *db.DB) template.FuncMap {
	return template.FuncMap{
		"countBackgroundJobs": func() int {
			var keys []db.FallbackQueue
			database.Debug().Model(&db.FallbackQueue{}).Find(&keys)
			return len(keys)
		},
		"countBackgroundJobsInChannel": func(channel string) int {
			count := 0
			var keys []db.FallbackQueue
			database.Debug().Model(&db.FallbackQueue{}).Find(&keys)
			for _, key := range keys {
				if key.Channel != channel {
					continue
				}
				count++
			}

			return count
		},
	}
}

func (c *listCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "list queue",
			Description: "list all queued commands",
			Examples: []string{
				"list queue",
			},
		},
		{
			Command:     "list queue in channel",
			Description: "list queued commands in current channel",
			Examples: []string{
				"list queue in channel",
			},
		},
	}
}
