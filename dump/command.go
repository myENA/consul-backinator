package dump

import (
	"fmt"
	"github.com/myENA/consul-backinator/common"
	"log"
)

// primary configuration
type config struct {
	fileName      string
	cryptKey      string
	pathTransform string
	plainDump     bool
}

// Command is a Command implementation that runs the backup operation
type Command struct {
	Self            string
	config          *config
	pathTransformer *common.PathTransformer
}

// Run is a function to run the command
func (c *Command) Run(args []string) int {
	var err error // error holder

	// init config
	c.config = new(config)

	// setup flags
	if err = c.setupFlags(args); err != nil {
		log.Printf("[Error] Startup failed: %s", err.Error())
		return 1
	}

	// build transformer if needed
	if c.pathTransformer, err = common.NewTransformer(c.config.pathTransform); err != nil {
		log.Printf("[Error] Failed to initialize path transformer: %s", err.Error())
		return 1
	}

	// dump data
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

	Dump and optionally decode the contents of a backup file to stdout.

Options:

	-file         Source filename (default: "consul.bak")
	-key          Passphrase for data encryption and signature validation (default: "password")
	-transform    Optional path transformation (oldPath,newPath...)
	-plain        Dump only the key and decoded value

`, c.Self)
}