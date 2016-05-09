package backup

import (
	"flag"
	"fmt"
	"os"
)

// init instance configuration
func (c *Command) setupFlags(args []string) error {
	// init flagset
	cmdFlags := flag.NewFlagSet("backup", flag.ContinueOnError)
	cmdFlags.Usage = func() { fmt.Fprint(os.Stdout, c.Help()); os.Exit(0) }

	// declare flags
	cmdFlags.StringVar(&c.config.fileName, "file", "consul.bak",
		"Backup filename")
	cmdFlags.StringVar(&c.config.cryptKey, "key", "password",
		"Passphrase for data encryption and signature validation")
	cmdFlags.StringVar(&c.config.pathTransform, "transform", "",
		"Optional path transformation")
	cmdFlags.StringVar(&c.config.consulPrefix, "prefix", "/",
		"Optional prefix from under which all keys will be fetched")
	cmdFlags.StringVar(&c.config.consulConfig.Address, "addr", "",
		"Optional consul address and port")
	cmdFlags.StringVar(&c.config.consulConfig.Scheme, "scheme", "",
		"Optional consul scheme")
	cmdFlags.StringVar(&c.config.consulConfig.Datacenter, "dc", "",
		"Optional consul datacenter")
	cmdFlags.StringVar(&c.config.consulConfig.Token, "token", "",
		"Optional consul access token")

	// parse flags and ignore error
	if err := cmdFlags.Parse(args); err != nil {
		return nil
	}

	// always okay
	return nil
}
