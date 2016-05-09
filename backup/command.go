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
		log.Printf("[Error] Startup failed: %s", err.Error())
		return 1
	}

	// build client
	if c.consulClient, err = c.config.consulConfig.New(); err != nil {
		log.Printf("[Error] Failed initialize consul client: %s", err.Error())
		return 1
	}

	// build transformer if needed
	if c.pathTransformer, err = common.NewTransformer(c.config.pathTransform); err != nil {
		log.Printf("[Error] Failed to initialize path transformer: %s", err.Error())
		return 1
	}

	// backup keys
	if count, err = c.backupKeys(); err != nil {
		log.Printf("[Error] Failed to backup data: %s", err.Error())
		return 1
	}

	// show success
	log.Printf("[Success] Backed up %d keys from %s%s to %s",
		count,
		c.config.consulConfig.Address,
		c.config.consulPrefix,
		c.config.fileName)

	// make sure they know to keep the sig
	fmt.Printf("Keep your backup (%s) and signature (%s.sig) files "+
		"in a safe place.\nYou will need both to restore your data.\n",
		c.config.fileName, c.config.fileName)

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

	Performs a backup operation against a consul cluster KV store.

Options:

	-file         Sets the backup filename (default: "consul.bak")
	-key          Passphrase for data encryption and signature validation (default: "password")
	-transform    Optional path transformation (oldPath,newPath...)
	-prefix       Optional prefix from under which all keys will be fetched
	-addr         Optional consul address and port (default: "127.0.0.1:8500")
	-scheme       Optional consul scheme ("http" or "https")
	-dc           Optional consul datacenter
	-token        Optional consul access token

`, c.Self)
}
