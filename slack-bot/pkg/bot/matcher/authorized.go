package matcher

import (
	"errors"
	"strings"

	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client"
)

type AuthorizationFunc func(message msg.Message) bool

// NewAuthorizedMatcher is a wrapper to only executable by a user who is authorized by edith server
func NewAuthorizedMatcher(slackClient client.SlackClient, matcher Matcher, private bool) Matcher {
	return authorizedMatcher{
		matcher,
		slackClient,
		private,
		func(message msg.Message) bool {
			return message.DBUser != nil
		},
		"sorry, you are not authorized, create token first",
	}
}

// NewAdminMatcher is a wrapper to only executable by a user who is admin for the edith server
func NewAdminMatcher(slackClient client.SlackClient, matcher Matcher, private bool) Matcher {
	return authorizedMatcher{
		matcher,
		slackClient,
		private,
		func(message msg.Message) bool {
			return message.DBUser != nil && strings.ToLower(message.DBUser.Designation) == "admin"
		},
		"sorry, you are no admin and not allowed to execute this command",
	}
}

// NewManagerMatcher is a wrapper to only executable by a user who can be either hr, manager or admin
func NewManagerMatcher(slackClient client.SlackClient, matcher Matcher, private bool) Matcher {
	return authorizedMatcher{
		matcher,
		slackClient,
		private,
		func(message msg.Message) bool {
			return message.DBUser != nil &&
				(strings.ToLower(message.DBUser.Designation) == "hr" ||
					strings.ToLower(message.DBUser.Designation) == "manager" ||
					strings.ToLower(message.DBUser.Designation) == "admin")
		},
		"sorry, you are not allowed to execute this command",
	}
}

type authorizedMatcher struct {
	matcher           Matcher
	slackClient       client.SlackClient
	private           bool
	authorizationFunc AuthorizationFunc
	errorMsg          string
}

func (m authorizedMatcher) Match(message msg.Message) (Runner, Result) {
	run, result := m.matcher.Match(message)
	if run == nil {
		// the wrapped command didn't match...ignore
		return nil, nil
	}

	if m.private && !strings.HasPrefix(message.Channel, "D") && !message.InternalMessage {
		return func(match Result, message msg.Message) {
			m.slackClient.AddReaction("❌", message)
			m.slackClient.ReplyError(
				message,
				errors.New("sorry, this command can be used in private channel only (edith)"),
			)
		}, Result{}
	}

	if m.authorizationFunc(message) {
		return run, result
	}

	if message.User == "cron" {
		return run, result
	}

	slackUser := m.slackClient.GetUserDetails(message.GetUser(), message)
	if message.DBUser == nil && (slackUser.IsOwner || slackUser.IsAdmin) {
		return run, result
	}

	return func(match Result, message msg.Message) {
		m.slackClient.AddReaction("❌", message)
		m.slackClient.ReplyError(
			message,
			errors.New(m.errorMsg),
		)
	}, Result{}
}
