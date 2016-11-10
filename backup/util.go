package backup

import (
	"flag"
	"fmt"
	cc "github.com/myENA/consul-backinator/common/config"
	ccns "github.com/myENA/consul-backinator/common/consul"
	"os"
	"strings"
)

// init instance configuration
func (c *Command) setupFlags(args []string) error {
	// init flagset
	cmdFlags := flag.NewFlagSet("backup", flag.ContinueOnError)
	cmdFlags.Usage = func() { fmt.Fprint(os.Stdout, c.Help()); os.Exit(0) }

	// declare flags
	cmdFlags.StringVar(&c.config.fileName, "file", "consul.bak",
		"Destination")
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

	// add shared flags
	cc.AddSharedConsulFlags(cmdFlags, c.config.consulConfig)

	// parse flags and ignore error
	if err := cmdFlags.Parse(args); err != nil {
		return nil
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
