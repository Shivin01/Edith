package matcher

import (
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/msg"
	log "github.com/sirupsen/logrus"
)

// WildcardMatcher will pass all messages into the runner. The runner needs to check if the event is relevant or not
// if the runner returns "true", the command is marked as executed and no other commands are checked
func WildcardMatcher(run wildcardRunner) Matcher {
	return wildcardMatcher{run}
}

type wildcardRunner func(ref msg.Ref, text string) bool

type wildcardMatcher struct {
	run wildcardRunner
}

func (m wildcardMatcher) Match(message msg.Message) (Runner, Result) {
	if m.run(message, message.GetText()) {
		// command got executed! -> further commands will be ignored
		return func(match Result, message msg.Message) {
			log.Infof("successfully handled command %s", message.GetText())
		}, Result{}
	}

	return nil, nil
}
