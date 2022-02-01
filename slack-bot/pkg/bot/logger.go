package bot

import (
	"os"
	"strings"

	"github.com/immanoj16/edith/pkg/bot/msg"
	"github.com/immanoj16/edith/pkg/client"
	"github.com/immanoj16/edith/pkg/config"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
)

// InitLogger provides logger instance for the given config
func InitLogger(cfg config.Logger) {
	level, err := log.ParseLevel(cfg.Level)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(os.Stdout)
	log.SetLevel(level)

	if cfg.File != "" {
		log.AddHook(lfshook.NewHook(
			cfg.File,
			&log.TextFormatter{},
		))
	}
}

// get a log.Entry with some user related fields
func (b *Bot) getUserBasedLogger(ref msg.Ref) *log.Entry {
	_, username := client.GetUserIDAndName(ref.GetUser())

	var channel string
	if strings.HasPrefix(ref.GetChannel(), "D") {
		channel = "@" + username
	} else {
		channel = ref.GetChannel()
	}

	var user string
	if ref.IsInternalMessage() {
		user = "edith"
	} else {
		if ref.GetDBUser() == nil {
			user = ref.GetUser()
		} else {
			user = ref.GetDBUser().GetUsername()
		}
	}

	return log.
		WithField("channel", channel).
		WithField("user", user)
}
