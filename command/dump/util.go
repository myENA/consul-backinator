package dump

import (
	"flag"
	"fmt"
	"os"

	cc "github.com/myENA/consul-backinator/common/config"
)

// setupFlags initializes the instance configuration
func (c *Command) setupFlags(args []string) error {
	var cmdFlags *flag.FlagSet // instance flagset

	// init config if needed
	if c.config == nil {
		c.config = new(config)
	}

	// init flagset
	cmdFlags = flag.NewFlagSet("dump", flag.ContinueOnError)
	cmdFlags.Usage = func() { fmt.Fprint(os.Stdout, c.Help()); os.Exit(0) }

	// declare flags
	cmdFlags.StringVar(&c.config.fileName, "file", "consul.bak",
		"Destination file target")
	cmdFlags.StringVar(&c.config.cryptKey, "key", "password",
		"Passphrase for data encryption and signature validation")
	cmdFlags.BoolVar(&c.config.plainDump, "plain", false,
		"Dump a reduced set of key/value data")

	// parse flags and ignore error
	if err := cmdFlags.Parse(args); err != nil {
		return nil
	}

	// check for remaining garbage
	if cmdFlags.NArg() > 0 {
		return cc.ErrUnknownArg
	}

	// always okay
	return nil
}
