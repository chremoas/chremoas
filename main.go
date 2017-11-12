package main

import (
	"github.com/abaeve/chremoas/bot"
	"github.com/micro/go-bot/input"
	_ "github.com/micro/go-bot/input/slack"
	"github.com/micro/go-micro/cmd"
	_ "github.com/micro/go-plugins/micro/bot/input/discord"
)

var Version string = "1.0.0"

func main() {
	app := cmd.App()
	app.Flags = bot.DefaultFlags
	//app.Commands = append(app.Commands, bot.Commands()...)
	app.Action = bot.Run //func(context *cli.Context) { cli.ShowAppHelp(context) }

	// setup input flags
	for _, input := range input.Inputs {
		app.Flags = append(app.Flags, input.Flags()...)
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

	cmd.Init(
		cmd.Name("chremoas"),
		cmd.Description("A bot to kill the Dramiel"),
		cmd.Version(Version),
	)
}
