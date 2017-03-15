package backup

import (
	"flag"
	"fmt"
	cc "github.com/myENA/consul-backinator/common/config"
	ccns "github.com/myENA/consul-backinator/common/consul"
	"os"
	"strings"
)

// setupFlags initializes the instance configuration
func (c *Command) setupFlags(args []string) error {
	var cmdFlags *flag.FlagSet // instance flagset
	var err error              // error holder

	// init consul config if needed
	if consulConfig == nil {
		consulConfig = new(ccns.Config)
	}

	// init flagset
	cmdFlags = flag.NewFlagSet("backup", flag.ContinueOnError)
	cmdFlags.Usage = func() { fmt.Fprint(os.Stdout, c.Help()); os.Exit(0) }

	// declare flags
	cmdFlags.StringVar(&kvFileName, "file", "consul.bak",
		"Destination")
	cmdFlags.StringVar(&cryptKey, "key", "password",
		"Passphrase for data encryption and signature validation")
	cmdFlags.BoolVar(&skipKV, "nokv", false,
		"Do not attempt to backup kv data")
	cmdFlags.StringVar(&aclFileName, "acls", "",
		"Optional backup filename for acl tokens")
	cmdFlags.StringVar(&queryFileName, "queries", "",
		"Optional backup filename for query definitions")
	cmdFlags.StringVar(&pathTransformation, "transform", "",
		"Optional path transformation")
	cmdFlags.StringVar(&consulPrefix, "prefix", "/",
		"Optional prefix from under which all keys will be fetched")

	// add shared flags
	cc.AddSharedConsulFlags(cmdFlags, consulConfig)

	// parse flags and ignore error
	if err = cmdFlags.Parse(args); err != nil {
		return nil
	}

	// populate potentially missing config items
	cc.AddEnvDefaults(consulConfig)

	// fixup prefix per upstream issue 2403
	// https://github.com/hashicorp/consul/issues/2403
	consulPrefix = strings.TrimPrefix(consulPrefix,
		ccns.Separator)

	// always okay
	return nil
}
