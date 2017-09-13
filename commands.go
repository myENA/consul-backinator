package main

import (
	stdLog "log"
	"os"

	"github.com/mitchellh/cli"
	"github.com/myENA/consul-backinator/backup"
	"github.com/myENA/consul-backinator/dump"
	"github.com/myENA/consul-backinator/restore"
)

// package global logger
var logger *stdLog.Logger

// available commands
var cliCommands map[string]cli.CommandFactory

// init command factory
func init() {
	// init logger
	logger = stdLog.New(os.Stderr, "", stdLog.LstdFlags)

	// register sub commands
	cliCommands = map[string]cli.CommandFactory{
		"backup": func() (cli.Command, error) {
			return &backup.Command{
				Self: os.Args[0],
				Log:  logger,
			}, nil
		},
		"restore": func() (cli.Command, error) {
			return &restore.Command{
				Self: os.Args[0],
				Log:  logger,
			}, nil
		},
		"dump": func() (cli.Command, error) {
			return &dump.Command{
				Self: os.Args[0],
				Log:  logger,
			}, nil
		},
	}
}
