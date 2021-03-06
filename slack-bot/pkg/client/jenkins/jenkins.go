package jenkins

import (
	"context"
	"time"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/util"
	"github.com/Shivin01/Edith/slack-bot/pkg/client"
	"github.com/bndr/gojenkins"
	log "github.com/sirupsen/logrus"
)

const (
	// iconRunning and the other ones are jenkins specific Slack reactions which are used for different kind of commands
	iconRunning = "🔄"
	iconSuccess = "✅"
	iconFailed  = "❌"
	iconPending = "☕"
	iconAborted = "black_circle_for_record"

	// we are polling the job every 10s-5min with a increasing delay
	minDelay = time.Second * 10
	maxDelay = time.Minute * 5
)

// JobResult holds the build result of a Build
type JobResult struct {
	build  *gojenkins.Build
	status string
}

// Job is a interface of gojenkins.Job
type Job interface {
	Poll(ctx context.Context) (int, error)
	GetLastBuild(ctx context.Context) (*gojenkins.Build, error)
	GetBuild(ctx context.Context, id int64) (*gojenkins.Build, error)
}

// WatchBuild will return a chan which is filled/closed when the build finished
func WatchBuild(build *gojenkins.Build) <-chan JobResult {
	resultChan := make(chan JobResult, 1)

	go func() {
		defer close(resultChan)

		delay := util.GetIncreasingDelay(minDelay, maxDelay)

		for {
			time.Sleep(delay.GetNextDelay())

			if !build.IsRunning(context.TODO()) {
				resultChan <- JobResult{
					status: build.GetResult(),
					build:  build,
				}

				return
			}
		}
	}()

	return resultChan
}

func processHooks(commands []string, ref msg.Ref, params Parameters) {
	for _, command := range commands {
		temp, err := util.CompileTemplate(command)
		if err != nil {
			log.Warn(err)
			continue
		}
		text, _ := util.EvalTemplate(temp, util.Parameters(params))
		client.HandleMessage(ref.WithText(text))
	}
}
