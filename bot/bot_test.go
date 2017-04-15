package bot

import (
	"errors"
	"flag"
	"strings"
	"testing"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-bot/command"
	"github.com/micro/go-bot/input"

	"github.com/abaeve/services-common/config"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/cmd"
	"github.com/micro/go-micro/registry/mock"
	_ "github.com/micro/go-plugins/micro/bot/input/discord"
	"os"
)

type testInput struct {
	send chan *input.Event
	recv chan *input.Event
	exit chan bool
}

func (t *testInput) Flags() []cli.Flag {
	return nil
}

func (t *testInput) Init(*cli.Context) error {
	return nil
}

func (t *testInput) Close() error {
	select {
	case <-t.exit:
	default:
		close(t.exit)
	}
	return nil
}

func (t *testInput) Send(event *input.Event) error {
	if event == nil {
		return errors.New("nil event")
	}

	select {
	case <-t.exit:
		return errors.New("connection closed")
	case t.send <- event:
		return nil
	}
}

func (t *testInput) Recv(event *input.Event) error {
	if event == nil {
		return errors.New("nil event")
	}

	select {
	case <-t.exit:
		return errors.New("connection closed")
	case ev := <-t.recv:
		*event = *ev
		return nil
	}

}

func (t *testInput) Start() error {
	return nil
}

func (t *testInput) Stop() error {
	return nil
}

func (t *testInput) Stream() (input.Conn, error) {
	return t, nil
}

func (t *testInput) String() string {
	return "test"
}

func TestBot(t *testing.T) {
	flagSet := flag.NewFlagSet("test", flag.ExitOnError)
	app := cli.NewApp()
	ctx := cli.NewContext(app, flagSet, nil)

	io := &testInput{
		send: make(chan *input.Event),
		recv: make(chan *input.Event),
		exit: make(chan bool),
	}

	inputs := map[string]input.Input{
		"test": io,
	}

	commands := map[string]command.Command{
		"^echo ": command.NewCommand("echo", "test usage", "test description", func(args ...string) ([]byte, error) {
			return []byte(strings.Join(args[1:], " ")), nil
		}),
	}

	service := micro.NewService(
		micro.Registry(mock.NewRegistry()),
	)

	bot := newBot(ctx, inputs, commands, service)

	if err := bot.start(); err != nil {
		t.Fatal(err)
	}

	// send command
	select {
	case io.recv <- &input.Event{
		Meta: map[string]interface{}{},
		Type: input.TextEvent,
		Data: []byte("echo test"),
	}:
	case <-time.After(time.Second):
		t.Fatal("timed out sending event")
	}

	// recv output
	select {
	case ev := <-io.send:
		if string(ev.Data) != "test" {
			t.Fatal("expected 'test', got: ", string(ev.Data))
		}
	case <-time.After(time.Second):
		t.Fatal("timed out receiving event")
	}

	if err := bot.stop(); err != nil {
		t.Fatal(err)
	}
}

func Test_cliContextFromConfiguration(t *testing.T) {
	App = cmd.App()
	App.Flags = DefaultFlags
	App.Writer = os.Stdout

	// setup input flags
	for _, input := range input.Inputs {
		App.Flags = append(App.Flags, input.Flags()...)
	}

	conf := config.Configuration{}
	err := conf.Load("../application.dist.yaml")

	if err != nil {
		t.Fatalf("Error while reading the configuration: (%s)", err)
	}

	ctx := cliContextFromConfiguration(&conf)

	if ctx == nil {
		t.Fatal("Context was nil")
	}

	if ctx.String("server_name") != conf.Name {
		t.Errorf("Expected server_name: (%s) but received: (%s)", conf.Name, ctx.String("server_name"))
	}

	if ctx.String("inputs") != "slack,discord" {
		t.Errorf("Expected inputs: (slack,discord) but received: (%s)", ctx.String("inputs"))
	}

	if ctx.String("namespace") != conf.Namespace {
		t.Errorf("Expected namespace: (%s) but received: (%s)", conf.Namespace, ctx.String("namespace"))
	}

	if ctx.Int("register_ttl") != conf.Registry.RegisterTTL {
		t.Errorf("Expected register_ttl: (20) but received: (%s)", ctx.Int("register_ttl"))
	}

	if ctx.Int("register_interval") != conf.Registry.RegisterInterval {
		t.Errorf("Expected register_interval: (slack,discord) but received: (%s)", ctx.Int("register_interval"))
	}

	//if ctx.String("hipchat_debug") != "slack,discord" {
	//	t.Errorf("Expected hipchat_debug: (slack,discord) but received: (%s)", ctx.String("hipchat_debug"))
	//}
	//
	//if ctx.String("hipchat_username") != "slack,discord" {
	//	t.Errorf("Expected hipchat_username: (slack,discord) but received: (%s)", ctx.String("hipchat_username"))
	//}
	//
	//if ctx.String("hipchat_password") != "slack,discord" {
	//	t.Errorf("Expected hipchat_password: (slack,discord) but received: (%s)", ctx.String("hipchat_password"))
	//}
	//
	//if ctx.String("hipchat_server") != "slack,discord" {
	//	t.Errorf("Expected hipchat_server: (slack,discord) but received: (%s)", ctx.String("hipchat_server"))
	//}

	if !ctx.Bool("slack_debug") {
		t.Errorf("Expected slack_debug: (%t) but received: (%t)", conf.Chat.Slack.Debug, ctx.Bool("slack_debug"))
	}

	if ctx.String("slack_token") != conf.Chat.Slack.Token {
		t.Errorf("Expected slack_token: (%s) but received: (%s)", conf.Chat.Slack.Token, ctx.String("slack_token"))
	}

	if ctx.String("discord_token") != conf.Chat.Discord.Token {
		t.Errorf("Expected discord_token: (%s) but received: (%s)", conf.Chat.Discord.Token, ctx.String("discord_token"))
	}

	if ctx.String("discord_whitelist") != "11234567890,2234567890,3234567890" {
		t.Errorf("Expected discord_whitelist: (11234567890,2234567890,3234567890) but received: (%s)", ctx.String("discord_whitelist"))
	}

	if ctx.String("discord_prefix") != conf.Chat.Discord.Prefix {
		t.Errorf("Expected discord_prefix: (!) but received: (%s)", ctx.String("discord_prefix"))
	}
}
