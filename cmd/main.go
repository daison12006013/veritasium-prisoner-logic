package main

import (
	"os"

	prisoner "github.com/daison12006013/veritasium-prisoner-logic"

	"github.com/lucidfy/lucid/pkg/errors"
	cli "github.com/urfave/cli/v2"
)

func main() {
	consoleApplication(prisoner.Prisoner().Command)
}

func consoleApplication(cmds ...*cli.Command) {
	app := &cli.App{
		Name:     "Run",
		Usage:    "A console command runner for lucid!",
		Commands: cmds,
	}

	err := app.Run(os.Args)
	if errors.Handler("error running run console command", err) {
		panic(err)
	}
}
