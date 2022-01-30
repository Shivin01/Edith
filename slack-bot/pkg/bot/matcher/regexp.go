package matcher

import (
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/bot/util"
	"regexp"
)

// NewRegexpMatcher matches a command based on a given regexp. 2 additional rules:
// - it's case insensitive
// - it always has to match the full line (adding ^ and $ implicitly)
func NewRegexpMatcher(regexpString string, run Runner) Matcher {
	return regexpMatcher{
		regexp: util.CompileRegexp(regexpString),
		run:    run,
	}
}

type regexpMatcher struct {
	regexp *regexp.Regexp
	run    Runner
}

func (m regexpMatcher) Match(message msg.Message) (Runner, Result) {
	matches := m.regexp.FindStringSubmatch(message.GetText())
	if len(matches) == 0 {
		return nil, nil
	}

	match := make(Result, len(m.regexp.SubexpNames()))
	for idx, name := range m.regexp.SubexpNames() {
		match[name] = matches[idx]
	}

	return m.run, match
}
