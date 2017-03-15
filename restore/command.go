package restore

import (
	"fmt"
	ccns "github.com/myENA/consul-backinator/common/consul"
	ct "github.com/myENA/consul-backinator/common/transformer"
	"log"
)

// Command is a Command implementation that runs the backup operation
type Command struct {
	Self            string
	consulClient    *ccns.Client
	pathTransformer *ct.PathTransformer
}

// global configuration variables
var (
	kvFileName         string
	cryptKey           string
	skipKV             bool
	aclFileName        string
	queryFileName      string
	pathTransformation string
	delTree            bool
	consulPrefix       string
	consulConfig       *ccns.Config
)

// Run is a function to run the command
func (c *Command) Run(args []string) int {
	var err error // error holder
	var count int // key counter

	// setup flags
	if err = c.setupFlags(args); err != nil {
		log.Printf("[Error] Init failed: %s", err.Error())
		return 1
	}

	// sanity check
	if skipKV && aclFileName == "" {
		log.Printf("[Error] Passing 'nokv' and an empty 'acls' file " +
			"doesn't make any sense.  You should specify an 'acls' file " +
			"when using the 'nokv' option.")
		return 1
	}

	// build client
	if c.consulClient, err = consulConfig.New(); err != nil {
		log.Printf("[Error] Failed initialize consul client: %s", err.Error())
		return 1
	}

	// build transformer if needed
	if c.pathTransformer, err = ct.New(pathTransformation); err != nil {
		log.Printf("[Error] Failed to initialize path transformer: %s", err.Error())
		return 1
	}

	// restore keys unless otherwise requested
	if !skipKV {
		if count, err = c.restoreKeys(); err != nil {
			log.Printf("[Error] Failed to restore kv data: %s", err.Error())
			return 1
		}

		// show success
		log.Printf("[Success] Restored %d keys from %s to %s/%s",
			count,
			kvFileName,
			consulConfig.Address,
			consulPrefix)
	}

	// restore acls if requested
	if aclFileName != "" {
		if count, err = c.restoreACLs(); err != nil {
			log.Printf("[Error] Failed to restore ACL tokens: %s", err.Error())
			return 1
		}

		// show success
		log.Printf("[Success] Restored %d ACL tokens from %s to %s",
			count,
			aclFileName,
			consulConfig.Address)
	}

	// exit clean
	return 0
}

// Synopsis shows the command summary
func (c *Command) Synopsis() string {
	return "Perform a restore operation"
}

// Help shows the detailed command options
func (c *Command) Help() string {
	return fmt.Sprintf(`Usage: %s restore [options]

	Performs a restore operation against a consul cluster.

Options:

	-file            Source filename or S3 location (default: "consul.bak")
	-key             Passphrase for data encryption and signature validation (default: "password")
	-nokv            Do not attempt to restore kv data
	-acls            Optional source filename or S3 location for acl tokens
	-queries         Optional source filename or S3 location for query definitions
	-transform       Optional path transformation (oldPath,newPath...)
	-delete          Delete all keys under specified prefix prior to restoration (default: false)
	-prefix          Prefix for delete operation
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
