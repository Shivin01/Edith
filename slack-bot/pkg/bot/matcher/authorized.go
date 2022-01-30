package matcher

import (
	"errors"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client"
)

// NewAuthorizedMatcher is a wrapper to only executable by a whitelisted admins user
func NewAuthorizedMatcher(slackClient client.SlackClient, matcher Matcher) Matcher {
	return authorizedMatcher{matcher, slackClient}
}

type authorizedMatcher struct {
	matcher     Matcher
	slackClient client.SlackClient
}

func (m authorizedMatcher) Match(message msg.Message) (Runner, Result) {
	run, result := m.matcher.Match(message)
	if run == nil {
		// the wrapped command didn't match...ignore
		return nil, nil
	}

	if message.DBUser != nil {
		return run, result
	}

	slackUser := m.slackClient.GetUserDetails(message.GetUser(), message)
	if slackUser.IsOwner || slackUser.IsAdmin {
		return run, result
	}

	return func(match Result, message msg.Message) {
		m.slackClient.AddReaction("❌", message)
		m.slackClient.ReplyError(
			message,
			errors.New("sorry, you are not authorized, create token first"),
		)
	}, Result{}
}
