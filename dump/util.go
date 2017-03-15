package dump

import (
	"flag"
	"fmt"
	"os"
)

// setupFlags initializes the instance configuration
func (c *Command) setupFlags(args []string) error {
	var cmdFlags *flag.FlagSet // instance flagset

	// init config if needed
	if config == nil {
		config = new(configStruct)
	}

	// init flagset
	cmdFlags = flag.NewFlagSet("dump", flag.ContinueOnError)
	cmdFlags.Usage = func() { fmt.Fprint(os.Stdout, c.Help()); os.Exit(0) }

	// declare flags
	cmdFlags.StringVar(&config.fileName, "file", "consul.bak",
		"Destination file target")
	cmdFlags.StringVar(&config.cryptKey, "key", "password",
		"Passphrase for data encryption and signature validation")
	cmdFlags.BoolVar(&config.plainDump, "plain", false,
		"Dump a reduced set of information")
	cmdFlags.BoolVar(&config.acls, "acls", false,
		"Specified file is an ACL token backup file")
	cmdFlags.BoolVar(&config.queries, "queries", false,
		"Specified file is a prepared query backup file")

	// parse flags and ignore error
	if err := cmdFlags.Parse(args); err != nil {
		return nil
	}

	// always okay
	return nil
}
