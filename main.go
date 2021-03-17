package main

import (
	"fmt"

	"github.com/micro/go-bot/input"
	_ "github.com/micro/go-bot/input/discord"
	_ "github.com/micro/go-bot/input/hipchat"
	_ "github.com/micro/go-bot/input/slack"
	"github.com/micro/go-micro/config/cmd"
	"go.uber.org/zap"

	chremoasPrometheus "github.com/chremoas/services-common/prometheus"

	"github.com/chremoas/chremoas/bot"
)

var (
	Version = "SET ME YOU KNOB"
	logger  *zap.Logger
)

func main() {
	var err error

	// TODO pick stuff up from the config
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	logger.Info("Initialized logger")

	go chremoasPrometheus.PrometheusExporter(logger)

	app := cmd.App()
	app.Flags = bot.DefaultFlags
	//app.Commands = append(app.Commands, bot.Commands()...)
	app.Action = bot.Run //func(context *cli.Context) { cli.ShowAppHelp(context) }

	// setup input flags
	for _, myInput := range input.Inputs {
		app.Flags = append(app.Flags, myInput.Flags()...)
	}

	for _, p := range bot.Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			app.Commands = append(app.Commands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			app.Flags = append(app.Flags, flags...)
		}
	}

	//setup(app)

	bot.App = app

	err = cmd.Init(
		cmd.Name("chremoas"),
		cmd.Description("A bot to kill the Dramiel"),
		cmd.Version(Version),
	)

	if err != nil {
		fmt.Printf("There was an error running the app: %v\n", err.Error())
	}
}
