package backup

import (
	"fmt"
	"github.com/myENA/consul-backinator/common"
	"log"
)

// primary configuration
type config struct {
	fileName      string
	cryptKey      string
	noKV          bool
	aclFileName   string
	pathTransform string
	consulPrefix  string
	consulConfig  *common.ConsulConfig
}

// Command is a Command implementation that runs the backup operation
type Command struct {
	Self            string
	config          *config
	consulClient    *common.ConsulClient
	pathTransformer *common.PathTransformer
}

// Run is a function to run the command
func (c *Command) Run(args []string) int {
	var err error // error holder
	var count int // key counter

	// init config
	c.config = new(config)

	// init consul config
	c.config.consulConfig = new(common.ConsulConfig)

	// setup flags
	if err = c.setupFlags(args); err != nil {
		log.Printf("[Error] Init failed: %s", err.Error())
		return 1
	}

	// sanity check
	if c.config.noKV && c.config.aclFileName == "" {
		log.Printf("[Error] Passing 'nokv' and an empty 'acls' file " +
			"doesn't make any sense.  You should specify an 'acls' file " +
			"when using the 'nokv' option.")
		return 1
	}

	// build client
	if c.consulClient, err = c.config.consulConfig.NewClient(); err != nil {
		log.Printf("[Error] Failed initialize consul client: %s", err.Error())
		return 1
	}

	// build transformer if needed
	if c.pathTransformer, err = common.NewTransformer(c.config.pathTransform); err != nil {
		log.Printf("[Error] Failed to initialize path transformer: %s", err.Error())
		return 1
	}

	// backup keys unless otherwise requested
	if !c.config.noKV {
		if count, err = c.backupKeys(); err != nil {
			log.Printf("[Error] Failed to backup key data: %s", err.Error())
			return 1
		}

		// show success
		log.Printf("[Success] Backed up %d keys from %s%s to %s",
			count,
			c.config.consulConfig.Address,
			c.config.consulPrefix,
			c.config.fileName)
	}

	// backup acls if requested
	if c.config.aclFileName != "" {
		if count, err = c.backupAcls(); err != nil {
			log.Printf("[Error] Failed to backup ACL tokens: %s", err.Error())
			return 1
		}

		// show success
		log.Printf("[Success] Backed up %d ACL tokens from %s to %s",
			count,
			c.config.consulConfig.Address,
			c.config.aclFileName)
	}

	// make sure they know to keep the sig
	fmt.Print("Keep your backup and signature files " +
		"in a safe place.\nYou will need both to restore your data.\n")

	// exit clean
	return 0
}

// Synopsis shows the command summary
func (c *Command) Synopsis() string {
	return "Perform a backup operation"
}

// Help shows the detailed command options
func (c *Command) Help() string {
	return fmt.Sprintf(`Usage: %s backup [options]

	Performs a backup operation against a consul cluster.

Options:

	-file            Sets the backup filename (default: "consul.bak")
	-key             Passphrase for data encryption and signature validation (default: "password")
	-nokv            Do not attempt to backup kv data
	-acls            Optional backup filename for acl tokens
	-transform       Optional path transformation (oldPath,newPath...)
	-prefix          Optional prefix from under which all keys will be fetched
	-addr            Optional consul address and port (default: "127.0.0.1:8500")
	-scheme          Optional consul scheme ("http" or "https")
	-dc              Optional consul datacenter
	-token           Optional consul access token
	-ca-cert         Optional path to a PEM encoded CA cert file to use to verify consul
	-ca-path         Optional path to a directory of PEM encoded CA cert files to verify consul
	-client-cert     Optional path to a PEM encoded client certificate for TLS authentication to consul
	-client-key      Optional path to an unencrypted PEM encoded private key matching the client certificate from -client-cert
	-tls-skip-verify Optional bool for verifying a TLS certificate. This is highly not recommended

`, c.Self)
}
