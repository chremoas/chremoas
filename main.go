package main

import (
	"fmt"
	"github.com/chremoas/chremoas/bot"
	"github.com/micro/go-bot/input"
	_ "github.com/micro/go-bot/input/hipchat"
	_ "github.com/micro/go-bot/input/slack"
	_ "github.com/micro/go-micro/agent/input/discord"
	"github.com/micro/go-micro/config/cmd"
)

var Version = "SET ME YOU KNOB"

func main() {
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

	err := cmd.Init(
		cmd.Name("chremoas"),
		cmd.Description("A bot to kill the Dramiel"),
		cmd.Version(Version),
	)

	if err != nil {
		fmt.Printf("There was an error running the app: %v\n", err.Error())
	}
}
