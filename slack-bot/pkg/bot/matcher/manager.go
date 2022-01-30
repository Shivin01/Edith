package matcher

import (
	"errors"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client"
	"strings"
)

// NewManagerMatcher is a wrapper to only executable by a whitelisted admins user
func NewManagerMatcher(slackClient client.SlackClient, matcher Matcher) Matcher {
	return managerMatcher{matcher, slackClient}
}

type managerMatcher struct {
	matcher     Matcher
	slackClient client.SlackClient
}

func (m managerMatcher) Match(message msg.Message) (Runner, Result) {
	run, result := m.matcher.Match(message)
	if run == nil {
		// the wrapped command didn't match...ignore
		return nil, nil
	}

	if message.DBUser != nil &&
		(strings.ToLower(message.DBUser.Designation) == "hr" ||
			strings.ToLower(message.DBUser.Designation) == "manager" ||
			strings.ToLower(message.DBUser.Designation) == "admin") {
		// valid admin -> execute the wrapped command
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
			errors.New("sorry, you are not allowed to execute this command"),
		)
	}, Result{}
}
