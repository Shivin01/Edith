package queue

import (
	"strings"
	"sync"

	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client"
	"github.com/immanoj16/edith/pkg/db"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	waitIcon = "coffee"
	doneIcon = "white_check_mark"
)

var mu sync.RWMutex

// AddRunningCommand registers a long running command, e.g. a running Jenkins job or watching a pull request
// it's doing following magic:
// - when using "trigger job XXX" and "then reply done" it will execute the "reply done" when the running command was finished
// - when the server got restarted, the fallbackCommand gets executed
func AddRunningCommand(message msg.Message, fallbackCommand string, database *db.DB) *RunningCommand {
	var queueKey string

	// store fallback command in storage:
	// when the bot restarts for any reason, it can recover, based on this fallback commands
	if fallbackCommand != "" {
		message.Text = fallbackCommand
		// add timestamp to the key to have a fix sorting by time
		queueKey = strings.ReplaceAll(message.Timestamp, ".", "") + "-" + message.GetUniqueKey()

		fallbackQueue := db.FallbackQueue{
			Channel:         message.Channel,
			User:            message.User,
			Timestamp:       message.Timestamp,
			Thread:          message.Thread,
			InternalMessage: message.InternalMessage,
			UpdatedMessage:  message.UpdatedMessage,
			QueueKey:        queueKey,
			Text:            message.Text,
		}

		if err := database.DB.Debug().Model(&db.FallbackQueue{}).Create(&fallbackQueue).Error; err != nil {
			log.Error(errors.Wrap(err, "error while saving fallback queue"))
		}
	}

	log.Infof("add a blocking process: %s", message.GetText())

	mu.Lock()
	defer mu.Unlock()

	key := message.GetUniqueKey()
	runningCommand := &RunningCommand{}
	runningCommand.wg.Add(1)
	runningCommands[key] = runningCommand

	go func() {
		// wait until blocking task is over
		runningCommand.Wait()

		mu.Lock()
		delete(runningCommands, key)
		mu.Unlock()

		if queueKey != "" {
			if err := database.DB.Debug().Model(&db.FallbackQueue{}).Where("queue_key = ?", queueKey).Delete(&db.FallbackQueue{}).Error; err != nil {
				log.Error(errors.Wrapf(err, "error while deleting processed queue entry %s", key))
			}
		}
	}()

	return runningCommand
}

// CountCurrentJobs will return the number of current pending/queued jobs
func CountCurrentJobs() int {
	mu.RLock()
	defer mu.RUnlock()

	return len(runningCommands)
}

func executeFallbackCommand(database *db.DB) {
	var keys []db.FallbackQueue
	database.DB.Debug().Model(&db.FallbackQueue{}).Find(&keys)

	if len(keys) == 0 {
		return
	}

	log.Infof("[Queue] Booted! I'll trigger %d command now", len(keys))

	for _, key := range keys {
		client.HandleMessage(msg.Message{
			MessageRef: msg.MessageRef{
				Channel:         key.Channel,
				User:            key.User,
				Timestamp:       key.Timestamp,
				Thread:          key.Thread,
				InternalMessage: key.InternalMessage,
				UpdatedMessage:  key.UpdatedMessage,
			},
			Text: key.Text,
		})
	}

	database.DB.Debug().Model(&db.FallbackQueue{}).Exec("DELETE FROM fallback_queues")
}
