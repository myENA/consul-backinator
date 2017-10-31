package restore

import (
	"flag"
	"fmt"
	"os"
	"strings"

	cc "github.com/myENA/consul-backinator/common/config"
	ccns "github.com/myENA/consul-backinator/common/consul"
)

// setupFlags initializes the instance configuration
func (c *Command) setupFlags(args []string) error {
	var cmdFlags *flag.FlagSet // instance flagset

	// init config if needed
	if c.config == nil {
		c.config = new(config)
	}

	// init consul config if needed
	if c.config.consulConfig == nil {
		c.config.consulConfig = new(ccns.Config)
	}

	// init flagset
	cmdFlags = flag.NewFlagSet("restore", flag.ContinueOnError)
	cmdFlags.Usage = func() { fmt.Fprint(os.Stdout, c.Help()); os.Exit(0) }

	// declare flags
	cmdFlags.StringVar(&c.config.fileName, "file", "consul.bak",
		"Source")
	cmdFlags.StringVar(&c.config.cryptKey, "key", "password",
		"Passphrase for data encryption and signature validation")
	cmdFlags.BoolVar(&c.config.noKV, "nokv", false,
		"Do not attempt to restore kv data")
	cmdFlags.StringVar(&c.config.aclFileName, "acls", "",
		"Optional source filename for acl tokens")
	cmdFlags.StringVar(&c.config.queryFileName, "queries", "",
		"Optional source filename for query definitions")
	cmdFlags.StringVar(&c.config.pathTransform, "transform", "",
		"Optional path transformation")
	cmdFlags.BoolVar(&c.config.delTree, "delete", false,
		"Delete all keys under specified prefix")
	cmdFlags.StringVar(&c.config.consulPrefix, "prefix", "/",
		"Prefix for delete operation")

	// add shared flags
	cc.AddSharedConsulFlags(cmdFlags, c.config.consulConfig)

	// parse flags and ignore error
	if err := cmdFlags.Parse(args); err != nil {
		return nil
	}

	// check for remaining garbage
	if cmdFlags.NArg() > 0 {
		return cc.ErrUnknownArg
	}

	// populate potentially missing config items
	cc.AddEnvDefaults(c.config.consulConfig)

	// fixup prefix per upstream issue 2403
	// https://github.com/hashicorp/consul/issues/2403
	c.config.consulPrefix = strings.TrimPrefix(c.config.consulPrefix,
		ccns.Separator)

	// always okay
	return nil
}
