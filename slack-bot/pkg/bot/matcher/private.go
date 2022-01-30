package matcher

import (
	"errors"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client"
	"strings"
)

// NewPrivateMatcher is a wrapper to only executable by a whitelisted admins user
func NewPrivateMatcher(slackClient client.SlackClient, matcher Matcher) Matcher {
	return privateMatcher{matcher, slackClient}
}

type privateMatcher struct {
	matcher     Matcher
	slackClient client.SlackClient
}

func (m privateMatcher) Match(message msg.Message) (Runner, Result) {
	if !strings.HasPrefix(message.Channel, "D") && !message.InternalMessage {
		return func(match Result, message msg.Message) {
			m.slackClient.AddReaction("❌", message)
			m.slackClient.ReplyError(
				message,
				errors.New("sorry, this command can be used in private channel only (edith)"),
			)
		}, Result{}
	}

	run, result := m.matcher.Match(message)
	if run == nil {
		// the wrapped command didn't match...ignore
		return nil, nil
	}

	return run, result
}
