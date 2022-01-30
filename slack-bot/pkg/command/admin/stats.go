package admin

import (
	"fmt"
	"github.com/immanoj16/edith/pkg/bot"
	"github.com/immanoj16/edith/pkg/bot/matcher"
	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/bot/util"
	"github.com/immanoj16/edith/pkg/config"
	"github.com/immanoj16/edith/pkg/db"
	"runtime"
	"strings"
	"time"
)

// bots uptime
var startTime = time.Now()

// newStatsCommand shows a bunch of runtime statistics of the edith (admin-only)
func newStatsCommand(base bot.BaseCommand, cfg *config.Config) bot.Command {
	return &statsCommand{base, cfg}
}

type statsCommand struct {
	bot.BaseCommand
	cfg *config.Config
}

func (c *statsCommand) GetMatcher() matcher.Matcher {
	return matcher.NewAdminMatcher(
		c.SlackClient,
		matcher.NewTextMatcher("edith stats", c.stats),
	)
}

func (c *statsCommand) stats(match matcher.Result, message msg.Message) {
	result := statsResult{}
	result.addLine("Here are some current stats:")

	c.collectStats(&result)
	c.SendMessage(message, result.String())
}

func (c *statsCommand) collectStats(result *statsResult) {
	m := &runtime.MemStats{}
	runtime.ReadMemStats(m)

	stat := db.Stat{}
	if err := c.DB.Debug().Model(&db.Stat{}).First(&stat).Error; err != nil {
		// handle error
	}

	result.addNewSection("Overall stats")
	result.addValue("Total commands executed", util.FormatInt(stat.TotalCommands))
	result.addValue("Unknown commands", util.FormatInt(stat.UnknownCommands))
	result.addValue("Unauthorized commands", util.FormatInt(stat.UnauthorizedCommands))
	result.addValue("Handled interactions/buttons", util.FormatInt(stat.Interactions))

	result.addNewSection("Server runtime")
	result.addValue("Registered crons", util.FormatInt(len(c.cfg.Crons)))
	//result.addValue("Queued commands", util.FormatInt(queue.CountCurrentJobs()))
	result.addValue("Goroutines", util.FormatInt(runtime.NumGoroutine()))
	result.addValue("Mem Alloc", util.FormatBytes(m.Alloc))
	result.addValue("Mem Sys", util.FormatBytes(m.Sys))
	result.addValue("Uptime", util.FormatDuration(time.Since(startTime)))
	result.addValue("NumGC (since start)", util.FormatInt(int(m.NumGC)))
	result.addValue("Bot Version", "beta")
	result.addValue("Go Version", runtime.Version())
}

type statsResult struct {
	strings.Builder
}

func (s *statsResult) addNewSection(section string) {
	s.addLine(fmt.Sprintf("*%s*:", section))
}

func (s *statsResult) addValue(name string, value string) {
	s.addLine(fmt.Sprintf("• %s: %s", name, value))
}

func (s *statsResult) addLine(line string) {
	_, _ = s.WriteString(line + "\n")
}

func (c *statsCommand) GetHelp() []bot.Help {
	return []bot.Help{
		{
			Command:     "edith stats",
			Description: "display runtime stats from the bots, like total processed commands etc",
		},
	}
}
