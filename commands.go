package main

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/myENA/consul-backinator/backup"
	"github.com/myENA/consul-backinator/dump"
	"github.com/myENA/consul-backinator/restore"
)

// available commands
var cliCommands map[string]cli.CommandFactory

// init command factory
func init() {
	// register sub commands
	cliCommands = map[string]cli.CommandFactory{
		"backup": func() (cli.Command, error) {
			return &backup.Command{
				Self: os.Args[0],
			}, nil
		},
		"restore": func() (cli.Command, error) {
			return &restore.Command{
				Self: os.Args[0],
			}, nil
		},
		"dump": func() (cli.Command, error) {
			return &dump.Command{
				Self: os.Args[0],
			}, nil
		},
	}
}
