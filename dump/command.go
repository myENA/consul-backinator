package dump

import (
	"fmt"
	"log"
)

// primary configuration
type config struct {
	fileName      string
	cryptKey      string
	pathTransform string
	plainDump     bool
	acls          bool
}

// Command is a Command implementation that runs the backup operation
type Command struct {
	Self   string
	config *config
}

// Run is a function to run the command
func (c *Command) Run(args []string) int {
	var err error // error holder

	// init config
	c.config = new(config)

	// setup flags
	if err = c.setupFlags(args); err != nil {
		log.Printf("[Error] Init failed: %s", err.Error())
		return 1
	}

	// dump data or acls
	if err = c.dumpData(); err != nil {
		log.Printf("[Error] Failed to dump data: %s", err.Error())
		return 1
	}

	// exit clean
	return 0
}

// Synopsis shows the command summary
func (c *Command) Synopsis() string {
	return "Dump a backup file"
}

// Help shows the detailed command options
func (c *Command) Help() string {
	return fmt.Sprintf(`Usage: %s dump [options]

	Dump the contents of a backup file to stdout.

Options:

	-file         Source filename (default: "consul.bak")
	-key          Passphrase for data encryption and signature validation (default: "password")
	-plain        Dump a reduced set of information
	-acls         Specified file is an ACL token backup file (only applicable if decoding)

Please see documentation on GitHub for a detailed explanation of all options.
https://github.com/myENA/consul-backinator

`, c.Self)
}
