package backup

import (
	"fmt"
	ccns "github.com/myENA/consul-backinator/common/consul"
	ct "github.com/myENA/consul-backinator/common/transformer"
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
	consulConfig  *ccns.Config
}

// Command is a Command implementation that runs the backup operation
type Command struct {
	Self            string
	config          *config
	consulClient    *ccns.Client
	pathTransformer *ct.PathTransformer
}

// Run is a function to run the command
func (c *Command) Run(args []string) int {
	var err error // error holder
	var count int // key counter

	// init config
	c.config = new(config)

	// init consul config
	c.config.consulConfig = new(ccns.Config)

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
	if c.consulClient, err = c.config.consulConfig.New(); err != nil {
		log.Printf("[Error] Failed initialize consul client: %s", err.Error())
		return 1
	}

	// build transformer if needed
	if c.pathTransformer, err = ct.New(c.config.pathTransform); err != nil {
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
		log.Printf("[Success] Backed up %d keys from %s/%s to %s",
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

	-file            Destination filename or S3 location (default: "consul.bak")
	-key             Passphrase for data encryption and signature validation (default: "password")
	-nokv            Do not attempt to backup kv data
	-acls            Optional backup filename or S3 location for acl tokens
	-transform       Optional path transformation (oldPath,newPath...)
	-prefix          Optional prefix from under which all keys will be fetched
	-addr            Optional consul address and port (default: "127.0.0.1:8500")
	-scheme          Optional consul scheme ("http" or "https")
	-dc              Optional consul datacenter
	-token           Optional consul access token
	-ca-cert         Optional path to a PEM encoded CA cert file
	-client-cert     Optional path to a PEM encoded client certificate
	-client-key      Optional path to an unencrypted PEM encoded private key
	-tls-skip-verify Optional bool for verifying a TLS certificate (not recommended)

Please see documentation on GitHub for a detailed explanation of all options.
https://github.com/myENA/consul-backinator

`, c.Self)
}
