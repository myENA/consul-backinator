package backup

import (
	"flag"
	"fmt"
	"github.com/myENA/consul-backinator/common"
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
	cmdFlags.BoolVar(&c.config.noKV, "nokv", false,
		"Do not attempt to backup kv data")
	cmdFlags.StringVar(&c.config.aclFileName, "acls", "",
		"Optional backup filename for acl tokens")
	cmdFlags.StringVar(&c.config.pathTransform, "transform", "",
		"Optional path transformation")
	cmdFlags.StringVar(&c.config.consulPrefix, "prefix", "/",
		"Optional prefix from under which all keys will be fetched")

	// Add shared Consul flags
	common.AddSharedConsulFlags(cmdFlags, c.config.consulConfig)

	// parse flags and ignore error
	if err := cmdFlags.Parse(args); err != nil {
		return nil
	}

	// always okay
	return nil
}
