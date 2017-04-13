package main

import (
	"github.com/abaeve/chremoas/bot"
	"github.com/micro/cli"
	"github.com/micro/go-bot/input"
	_ "github.com/micro/go-bot/input/slack"
	"github.com/micro/go-micro/cmd"
	_ "github.com/micro/go-plugins/micro/bot/input/discord"
)

var version string = "1.0.0"

func main() {
	app := cmd.App()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "server_name",
			EnvVar: "MICRO_SERVER_NAME",
			Usage:  "Name of the server. go.micro.srv.example",
		},
		cli.StringFlag{
			Name:  "inputs",
			Usage: "Inputs to load on startup",
		},
		cli.StringFlag{
			Name:   "namespace",
			Usage:  "Set the namespace used by the bot to find commands e.g. com.example.bot",
			EnvVar: "MICRO_BOT_NAMESPACE",
		},
		cli.IntFlag{
			Name:   "register_ttl",
			EnvVar: "MICRO_REGISTER_TTL",
			Usage:  "Register TTL in seconds",
		},
		cli.IntFlag{
			Name:   "register_interval",
			EnvVar: "MICRO_REGISTER_INTERVAL",
			Usage:  "Register interval in seconds",
		},
		cli.StringFlag{
			Name:   "configuration_file",
			Usage:  "The yaml configuration file for the service being loaded",
			Value:  "/etc/auth-srv/application.yaml",
			EnvVar: "CONFIGURATION_FILE",
		},
	}
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

	cmd.Init(
		cmd.Name("chremoas"),
		cmd.Description("A bot to kill the Dramiel"),
		cmd.Version(version),
	)
}
