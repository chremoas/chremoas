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

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry/mock"
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
