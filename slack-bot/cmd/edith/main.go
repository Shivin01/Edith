package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Shivin01/Edith/slack-bot/pkg/bot"
	"github.com/Shivin01/Edith/slack-bot/pkg/bot/util"
	"github.com/Shivin01/Edith/slack-bot/pkg/client"
	"github.com/Shivin01/Edith/slack-bot/pkg/command"
	"github.com/Shivin01/Edith/slack-bot/pkg/config"
	"github.com/Shivin01/Edith/slack-bot/pkg/db"
	log "github.com/sirupsen/logrus"
)

// main entry point for the edith application. Listens on incoming slack messages and handles them
func main() {
	var configFile string
	var verbose bool
	var showConfig bool
	flag.StringVar(&configFile, "config", "config.yaml", "Path to config.yaml. Can be a directory which will load all '*.yaml' inside")
	flag.BoolVar(&verbose, "verbose", false, "More verbose output")
	flag.BoolVar(&showConfig, "show-config", false, "Print the config as JSON and exit")
	flag.Parse()

	cfg, err := config.Load(configFile)
	checkError(err)

	if verbose {
		cfg.Logger.Level = "debug"
	}

	if showConfig {
		fmt.Println(config.Dump(cfg))
		os.Exit(0)
	}

	bot.InitLogger(cfg.Logger)
	log.Infof("Loaded config from %s", configFile)

	slackClient, err := client.GetSlackClient(cfg.Slack)
	checkError(err)

	database, err := db.NewDB(cfg.StoragePath)

	// get the list of all default commands
	commands := command.GetCommands(slackClient, cfg, database)

	b := bot.NewBot(cfg, slackClient, commands, database)
	err = b.Init()
	checkError(err)

	// start main loop!
	ctx := util.NewServerContext()
	b.Run(ctx)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
