package cli

import (
	"errors"
	"fmt"
	"io"
)

type CommandHandler func(args []string) error

type App struct {
	in io.Reader
	out io.Writer
	dbPath string
	handlers map[string]CommandHandler
}

func NewApp(in io.Reader, out io.Writer, dbPath string) *App {
	app := &App{
		in: in,
		out: out,
		dbPath: dbPath,
		handlers: make(map[string]CommandHandler),
	}

	app.handlers["capture"] = app.handleCapture
	// app.handlers["list"] = app.handleList
	// app.handlers["show"] = app.handleShow
	return app
}

func (a* App) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: notes <capture|list|show>")
	}

	cmd := Command {
		Name: args[0],
		Args: args[1:],
	}
	return a.HandleCommand(cmd)
}

func (a* App) HandleCommand(cmd Command) error {
	handler, exists := a.handlers[cmd.Name]
	if !exists{
		err := fmt.Errorf("Invalid command: %s", cmd.Name)
		return err
	}

	return handler(cmd.Args)
}