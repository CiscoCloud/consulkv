package main

import (
	"os"

	"github.com/CiscoCloud/consulkv/command"
	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.BasicUi{ Writer: os.Stdout }

	Commands = map[string]cli.CommandFactory{
		"delete": func() (cli.Command, error) {
			return &command.DeleteCommand{
				UI: ui,
			}, nil
		},

		"read": func() (cli.Command, error) {
			return &command.ReadCommand{
				UI: ui,
			}, nil
		},

		"write": func() (cli.Command, error) {
			return &command.WriteCommand{
				UI: ui,
			}, nil
		},

		"lock": func() (cli.Command, error) {
			return &command.LockCommand{
				UI: ui,
			}, nil
		},

		"unlock": func() (cli.Command, error) {
			return &command.UnlockCommand{
				UI: ui,
			}, nil
		},
	}
}
