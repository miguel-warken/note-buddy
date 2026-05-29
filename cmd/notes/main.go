package main

import (
	"os"
	"fmt"
	"note-buddy-main/internal/cli"
)

func main() {
	app := cli.NewApp(os.Stdin, os.Stdout, "notes.db")

	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}