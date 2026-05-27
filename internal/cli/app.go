package cli

import (
	"errors"
	"io"
	"fmt"
	"slices"
)

var allowedCommands = []string{"capture", "list", "show"}

type App struct {
	in io.Reader
	out io.Writer
	dbPath string
}

func NewApp(in io.Reader, out io.Writer, dbPath string) *App {
	return &App{
		in: in,
		out: out,
		dbPath: dbPath,
	}
}

func (a* App) Run(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: notes <capture|list|show>")
	}

	if slices.Contains(allowedCommands, args[0]) {
		fmt.Println("I exist.")
		return nil
	} else {
		return fmt.Errorf("Unknown command: %s", args[0])
	}

}