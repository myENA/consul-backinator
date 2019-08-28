package backup

import (
	"fmt"
	stdLog "log"

	ccns "github.com/myENA/consul-backinator/common/consul"
	ct "github.com/myENA/consul-backinator/common/transformer"
)

// primary configuration
type config struct {
	fileName          string
	cryptKey          string
	noKV              bool
	aclFileName       string
	aclPolicyFileName string
	legacyACLFileName string
	queryFileName     string
	pathTransform     string
	consulPrefix      string
	pathExclude       string
	consulConfig      *ccns.Config
}

// Command is a Command implementation that runs the backup operation
type Command struct {
	Self            string
	Log             *stdLog.Logger
	config          *config
	consulClient    *ccns.Client
	pathTransformer *ct.PathTransformer
}

// Run is a function to run the command
func (c *Command) Run(args []string) int {
	var err error // error holder
	var count int // key counter

	// setup flags
	if err = c.setupFlags(args); err != nil {
		c.Log.Printf("[Error] Setup failed: %s", err.Error())
		return 1
	}

	// sanity check
	if c.config.noKV && c.config.aclFileName == "" && c.config.queryFileName == "" {
		c.Log.Printf("[Error] Passing 'nokv' without an 'acls' or 'queries' file " +
			"doesn't make any sense.  You should specify an 'acls' or 'queries' file " +
			"when using the 'nokv' option.")
		return 1
	}

	// warn if backing up ACLs without policies
	if c.config.aclFileName != "" && c.config.aclPolicyFileName == "" {
		c.Log.Printf("[Warning] Backing up ACL tokens but not policies.  " +
			"You must specify a 'policies' file to backup new-style ACL policies")
	}

	// warn in change of behavior for ACL option
	if c.config.aclFileName != "" {
		c.Log.Printf("[Warning] The behavior of the 'acls' option has changed.  " +
			"This option now only backs-up new (1.4.0+) ACL tokens.  " +
			"To backup legacy ACLs please use the 'legacy-acls' option.")
	}

	// build client
	if c.consulClient, err = c.config.consulConfig.New(); err != nil {
		c.Log.Printf("[Error] Failed initialize consul client: %s", err.Error())
		return 1
	}

	// build transformer if needed
	if c.pathTransformer, err = ct.New(c.config.pathTransform); err != nil {
		c.Log.Printf("[Error] Failed to initialize path transformer: %s", err.Error())
		return 1
	}

	// backup keys unless otherwise requested
	if !c.config.noKV {
		if count, err = c.backupKeys(); err != nil {
			c.Log.Printf("[Error] Failed to backup key data: %s", err.Error())
			return 1
		}

		// show success
		c.Log.Printf("[Success] Backed up %d keys from %s/%s to %s",
			count,
			c.config.consulConfig.Address,
			c.config.consulPrefix,
			c.config.fileName)
	}

	// backup new-style acls if requested
	if c.config.aclFileName != "" {
		if count, err = c.backupACLTokens(); err != nil {
			c.Log.Printf("[Error] Failed to backup ACL tokens: %s", err.Error())
			return 1
		}

		// show success
		c.Log.Printf("[Success] Backed up %d ACL tokens from %s to %s",
			count,
			c.config.consulConfig.Address,
			c.config.aclFileName)
	}

	// backup new-style acl policies if requested
	if c.config.aclPolicyFileName != "" {
		if count, err = c.backupACLPolicies(); err != nil {
			c.Log.Printf("[Error] Failed to backup ACL policies: %s", err.Error())
			return 1
		}

		// show success
		c.Log.Printf("[Success] Backed up %d ACL policies from %s to %s",
			count,
			c.config.consulConfig.Address,
			c.config.aclPolicyFileName)
	}

	// backup legacy acls if requested
	if c.config.legacyACLFileName != "" {
		if count, err = c.backupLegacyACLs(); err != nil {
			c.Log.Printf("[Error] Failed to legacy backup ACL tokens: %s", err.Error())
			return 1
		}

		// show success
		c.Log.Printf("[Success] Backed up %d legacy ACL tokens from %s to %s",
			count,
			c.config.consulConfig.Address,
			c.config.legacyACLFileName)
	}

	// backup query definitions if requested
	if c.config.queryFileName != "" {
		if count, err = c.backupQueries(); err != nil {
			c.Log.Printf("[Error] Failed to backup query definitions: %s", err.Error())
			return 1
		}

		// show success
		c.Log.Printf("[Success] Backed up %d query definitions from %s to %s",
			count,
			c.config.consulConfig.Address,
			c.config.queryFileName)
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
	-acls            Optional backup filename or S3 location for new (1.4.0+) acl tokens
	-policies        Optional backup filename or S3 location for new (1.4.0+) acl policies
	-legacy-acls     Optional backup filename or S3 location for legacy acl tokens
	-queries         Optional backup filename or S3 location for prepared queries
	-transform       Optional path transformation (oldPath,newPath...)
	-prefix          Optional prefix from under which all keys will be fetched
	-exclude         Optional list of excluded paths (pathOne,PathTwo...)
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
