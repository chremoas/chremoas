// Package bot is a Hubot style bot that sits a microservice environment
package bot

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/micro/go-bot/command"
	"github.com/micro/go-bot/input"

	proto "github.com/chremoas/chremoas/proto"

	"github.com/chremoas/services-common/config"
	"golang.org/x/net/context"
	"io/ioutil"
	"strconv"
	"github.com/micro/go-micro/registry"
)

type bot struct {
	exit    chan bool
	ctx     *cli.Context
	service micro.Service

	sync.RWMutex
	inputs   map[string]input.Input
	commands map[string]command.Command
	services map[string]string
}

var (
	// Default server name
	Name = "go.micro.bot"
	// Namespace for commands
	Namespace = "go.micro.bot"
)

var DefaultFlags = []cli.Flag{
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
		Name: "registry_address",
		EnvVar: "MICRO_REGISTRY_ADDRESS",
		Usage: "The registry address and port <address>:<port>",
	},
	cli.StringFlag{
		Name:   "configuration_file",
		Usage:  "The yaml configuration file for the service being loaded",
		Value:  "/etc/auth-srv/application.yaml",
		EnvVar: "CONFIGURATION_FILE",
	},
}

var App *cli.App

func help(commands map[string]command.Command, serviceCommands []string) command.Command {
	usage := "help"
	desc := "Displays help for all known commands"

	var cmds []command.Command

	for _, cmd := range commands {
		cmds = append(cmds, cmd)
	}

	sort.Sort(sortedCommands{cmds})

	return command.NewCommand("help", usage, desc, func(args ...string) ([]byte, error) {
		response := []string{"\n"}
		for _, cmd := range cmds {
			response = append(response, fmt.Sprintf("%s - %s", cmd.Usage(), cmd.Description()))
		}
		response = append(response, serviceCommands...)
		return []byte(strings.Join(response, "\n")), nil
	})
}

func newBot(ctx *cli.Context, inputs map[string]input.Input, commands map[string]command.Command, service micro.Service) *bot {
	commands["^help$"] = help(commands, nil)

	return &bot{
		ctx:      ctx,
		exit:     make(chan bool),
		service:  service,
		commands: commands,
		inputs:   inputs,
		services: make(map[string]string),
	}
}

func (b *bot) loop(io input.Input) {
	log.Println("[bot][loop] starting", io.String())

	for {
		select {
		case <-b.exit:
			log.Println("[bot][loop] exiting", io.String())
			return
		default:
			if err := b.run(io); err != nil {
				log.Println("[bot][loop] error", err)
				time.Sleep(time.Second)
			}
		}
	}
}

func (b *bot) process(c input.Conn, ev input.Event) error {
	args := strings.Split(string(ev.Data), " ")
	if len(args) == 0 {
		return nil
	}

	b.RLock()
	defer b.RUnlock()

	// try built in command
	for pattern, cmd := range b.commands {
		// skip if it doesn't match
		if m, err := regexp.Match(pattern, ev.Data); err != nil || !m {
			continue
		}

		// matched, exec command
		rsp, err := cmd.Exec(args...)
		if err != nil {
			rsp = []byte("error executing cmd: " + err.Error())
		}

		// send response
		return c.Send(&input.Event{
			Meta: ev.Meta,
			From: ev.To,
			To:   ev.From,
			Type: input.TextEvent,
			Data: rsp,
		})
	}

	// no built in match
	// try service commands
	service := Namespace + "." + args[0]

	// is there a service for the command?
	if _, ok := b.services[service]; !ok {
		return nil
	}

	// make service request
	req := b.service.Client().NewRequest(service, "Command.Exec", &proto.ExecRequest{
		Sender: ev.From,
		Args:   args,
	})
	rsp := &proto.ExecResponse{}

	var response []byte

	// call service
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	if err := b.service.Client().Call(ctx, req, rsp); err != nil {
		response = []byte("error executing cmd: " + err.Error())
	} else if len(rsp.Error) > 0 {
		response = []byte("error executing cmd: " + rsp.Error)
	} else {
		response = rsp.Result
	}

	// send response
	return c.Send(&input.Event{
		Meta: ev.Meta,
		From: ev.To,
		To:   ev.From,
		Type: input.TextEvent,
		Data: response,
	})
}

func (b *bot) run(io input.Input) error {
	log.Println("[bot][loop] connecting to", io.String())

	c, err := io.Stream()
	if err != nil {
		return err
	}

	for {
		select {
		case <-b.exit:
			log.Println("[bot][loop] closing", io.String())
			return c.Close()
		default:
			var recvEv input.Event
			// receive input
			if err := c.Recv(&recvEv); err != nil {
				return err
			}

			// only process TextEvent
			if recvEv.Type != input.TextEvent {
				continue
			}

			if len(recvEv.Data) == 0 {
				continue
			}

			if err := b.process(c, recvEv); err != nil {
				return err
			}
		}
	}
}

func (b *bot) start() error {
	log.Println("[bot] starting")

	// Start inputs
	for _, io := range b.inputs {
		log.Println("[bot] starting input", io.String())

		if err := io.Init(b.ctx); err != nil {
			return err
		}

		if err := io.Start(); err != nil {
			return err
		}

		go b.loop(io)
	}

	// start watcher
	go b.watch()

	return nil
}

func (b *bot) stop() error {
	log.Println("[bot] stopping")
	close(b.exit)

	// Stop inputs
	for _, io := range b.inputs {
		log.Println("[bot] stopping input", io.String())
		if err := io.Stop(); err != nil {
			log.Println("[bot]", err)
		}
	}

	return nil
}

func (b *bot) watch() {
	commands := map[string]command.Command{}
	services := map[string]string{}

	// copy commands
	b.RLock()
	for k, v := range b.commands {
		commands[k] = v
	}
	b.RUnlock()

	// getHelp retries usage and description from bot service commands
	getHelp := func(service string) (string, error) {
		// is within namespace?
		if !strings.HasPrefix(service, Namespace) {
			return "", fmt.Errorf("%s not within namespace", service)
		}

		if p := strings.TrimPrefix(service, Namespace); len(p) == 0 {
			return "", fmt.Errorf("%s not a service", service)
		}

		// get command help
		req := b.service.Client().NewRequest(service, "Command.Help", &proto.HelpRequest{})
		rsp := &proto.HelpResponse{}

		err := b.service.Client().Call(context.Background(), req, rsp)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s - %s", rsp.Usage, rsp.Description), nil
	}

	serviceList, err := b.service.Client().Options().Registry.ListServices()
	if err != nil {
		// log error?
		return
	}

	var serviceCommands []string

	// create service commands
	for _, service := range serviceList {
		h, err := getHelp(service.Name)
		if err != nil {
			continue
		}
		services[service.Name] = h
		serviceCommands = append(serviceCommands, h)
	}

	b.Lock()
	b.commands["^help$"] = help(commands, serviceCommands)
	b.services = services
	b.Unlock()

	w, err := b.service.Client().Options().Registry.Watch()
	if err != nil {
		// log error?
		return
	}

	go func() {
		<-b.exit
		w.Stop()
	}()

	// watch for changes to services
	for {
		res, err := w.Next()
		if err != nil {
			return
		}

		if res.Action == "delete" {
			delete(services, res.Service.Name)
		} else {
			h, err := getHelp(res.Service.Name)
			if err != nil {
				continue
			}
			services[res.Service.Name] = h
		}

		var serviceCommands []string
		for _, v := range services {
			serviceCommands = append(serviceCommands, v)
		}

		b.Lock()
		b.commands["^help$"] = help(commands, serviceCommands)
		b.services = services
		b.Unlock()
	}
}

func Run(ctx *cli.Context) {
	var inputs []string
	var conf *config.Configuration

	if len(ctx.String("configuration_file")) > 0 {
		conf = &config.Configuration{}

		conf.Load(ctx.String("configuration_file"))

		ctx = cliContextFromConfiguration(conf)
	}

	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.String("server_name")
	}

	if len(ctx.String("namespace")) > 0 {
		Namespace = ctx.String("namespace")
	}

	// Parse flags
	if len(ctx.String("inputs")) == 0 {
		log.Println("[bot] no inputs specified")
		cli.ShowAppHelp(ctx)
		os.Exit(1)
	}

	inputs = strings.Split(ctx.String("inputs"), ",")
	if len(inputs) == 0 {
		log.Println("[bot] no inputs specified")
		cli.ShowAppHelp(ctx)
		os.Exit(1)
	}

	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	ios := make(map[string]input.Input)
	cmds := make(map[string]command.Command)

	// take other commands
	for pattern, cmd := range command.Commands {
		if c, ok := cmds[pattern]; ok {
			log.Printf("[bot] command %s already registered for pattern %s\n", c.String(), pattern)
			continue
		}
		// register command
		cmds[pattern] = cmd
	}

	// Parse inputs
	for _, io := range inputs {
		i, ok := input.Inputs[io]
		if !ok {
			log.Printf("[bot] input %s not found\n", i)
			os.Exit(1)
		}
		ios[io] = i
	}

	// setup service
	service := micro.NewService(
		micro.Name(Name),
		micro.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
		micro.Registry(
			registry.NewRegistry(
				registry.Addrs(ctx.GlobalString("registry_address")),
			),
		),
	)

	// Start bot
	b := newBot(ctx, ios, cmds, service)

	if err := b.start(); err != nil {
		log.Println("error starting bot", err)
		os.Exit(1)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	// Stop bot
	if err := b.stop(); err != nil {
		log.Println("error stopping bot", err)
	}
}

// Given a loaded configuration, this function will map what it finds available a newly created context to be loaded
// by inputs and plugins.
// Available Arguments and how they map to configuration
//--server_name				How the bot registers itself				(conf.Namespace + "." + conf.Name)
//--inputs				Inputs to load on startup				(conf.Inputs[])
//--namespace				Set the namespace used by the bot to find commands	(conf.Namespace)
//--register_ttl "0"			Register TTL in seconds					(conf.Registry.RegisterTTL)
//--register_interval "0"		Register interval in seconds				(conf.Register.RegisterInterval)
//--configuration_file			The yaml configuration file				(no equivalent in the created context... this loads it :P)
//--hipchat_debug			Hipchat debug output					(?)
//--hipchat_username			Hipchat XMPP username					(?)
//--hipchat_password 			Hipchat XMPP password					(?)
//--hipchat_server "chat.hipchat.com"	Hipchat XMPP server					(?)
//--slack_debug				Slack debug output					(conf.Chat.Slack.Debug)
//--slack_token				Slack token						(conf.Chat.Slack.Token)
//--discord_token			Discord token						(conf.Chat.Discord.Token)
//--discord_whitelist			Discord Whitelist (seperated by ,)			(conf.Chat.Discord.WhiteList[])
//--discord_prefix "Micro "		Discord Prefix						(conf.Chat.Discord.Prefix)
//--help, -h				show help						(no equivalent)
func cliContextFromConfiguration(conf *config.Configuration) *cli.Context {
	arguments := []string{}

	if len(conf.Namespace) > 0 {
		arguments = append(arguments, "--namespace="+conf.Namespace+".cmd")
		arguments = append(arguments, "--server_name="+conf.Namespace+".chremoas")
	}

	if len(conf.Inputs) > 0 {
		inputs := ""
		for _, input := range conf.Inputs {
			if inputs != "" {
				inputs = inputs + "," + input
			} else {
				inputs = input
			}
		}
		arguments = append(arguments, "--inputs="+inputs)
	}

	arguments = append(arguments, "--register_ttl="+strconv.Itoa(conf.Registry.RegisterTTL))
	arguments = append(arguments, "--register_interval="+strconv.Itoa(conf.Registry.RegisterInterval))

	if conf.Chat.Slack.Debug {
		arguments = append(arguments, "--slack_debug")
	}
	if len(conf.Chat.Slack.Token) > 0 {
		arguments = append(arguments, "--slack_token="+conf.Chat.Slack.Token)
	}
	if len(conf.Chat.Discord.Token) > 0 {
		arguments = append(arguments, "--discord_token="+conf.Chat.Discord.Token)
	}
	if len(conf.Chat.Discord.WhiteList) > 0 {
		whitelist := ""
		for _, whitelisted := range conf.Chat.Discord.WhiteList {
			if whitelist != "" {
				whitelist = whitelist + "," + whitelisted
			} else {
				whitelist = whitelisted
			}
		}
		arguments = append(arguments, "--discord_whitelist"+whitelist)
	}
	if len(conf.Chat.Discord.Prefix) > 0 {
		arguments = append(arguments, "--discord_prefix="+conf.Chat.Discord.Prefix)
	}

	set := flagSet("config_set", App.Flags)
	set.SetOutput(ioutil.Discard)
	err := set.Parse(arguments)
	nerr := normalizeFlags(App.Flags, set)
	ctx := cli.NewContext(App, set, nil)

	if nerr != nil {
		fmt.Fprintln(App.Writer, nerr)
		cli.ShowAppHelp(ctx)
		return nil
	}

	if err != nil {
		if App.OnUsageError != nil {
			_ = App.OnUsageError(ctx, err, false)
			return nil
		} else {
			fmt.Fprintf(App.Writer, "%s\n\n", "Incorrect Usage.")
			cli.ShowAppHelp(ctx)
			return nil
		}
	}

	return ctx
}
