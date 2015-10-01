package main

import (
	"errors"
	"flag"
	"github.com/hashicorp/consul/api"
	"log"
	"strings"
)

// consul kvp separator
const consulSeparator = "/"

// our instance configuration
type config struct {
	inFile        string
	outFile       string
	cryptKey      string
	pathTransform string
	pathReplacer  *strings.Replacer
	delTree       bool
	backupReq     bool
	restoreReq    bool
	consulAddr    string
	consulScheme  string
	consulDc      string
	consulToken   string
	consulPrefix  string
	consulClient  *api.Client
}

// bad count
var ErrBadTransform = errors.New("Path transformation list not even. " +
	"Transformations must be specified as pairs.")

// init instance configuration
func initConfig() (*config, error) {
	var c = new(config) // instance configuration
	var err error       // general error holder

	// declare flags
	flag.StringVar(&c.inFile, "in", "consul.bak",
		"Input file for restore operations")
	flag.StringVar(&c.outFile, "out", "consul.bak",
		"Output file for backup operations")
	flag.StringVar(&c.cryptKey, "key", "password",
		"Encryption key used to secure the destination file on backup "+
			"and read the input file on restore")
	flag.StringVar(&c.pathTransform, "transform", "",
		"Optional path transformation to be applied on backup and restore "+
			"(oldPath,newPath...)")
	flag.BoolVar(&c.delTree, "delete", false,
		"Optionally delete all keys under the destination prefix before restore")
	flag.BoolVar(&c.backupReq, "backup", false,
		"Trigger backup operation")
	flag.BoolVar(&c.restoreReq, "restore", false,
		"Trigger restore operation")
	flag.StringVar(&c.consulAddr, "addr", "",
		"Consul instance address and port (\"127.0.0.1:8500\")")
	flag.StringVar(&c.consulScheme, "scheme", "",
		"Optional consul instance scheme (\"http\" or \"https\")")
	flag.StringVar(&c.consulDc, "dc", "",
		"Optional consul datacenter label for backup and restore")
	flag.StringVar(&c.consulToken, "token", "",
		"Optional consul token to access the target cluster")
	flag.StringVar(&c.consulPrefix, "prefix", "/",
		"Optional prefix from under which all keys will be fetched or restored")

	// parse flags
	flag.Parse()

	// build replacer
	if c.pathTransform != "" {
		// split strings
		split := strings.Split(c.pathTransform, ",")
		// check count
		if (len(split) % 2) != 0 {
			return c, ErrBadTransform
		}
		// build replacer
		c.pathReplacer = strings.NewReplacer(split...)
	}

	// populate client
	err = c.buildClient()

	// return configuration and last error
	return c, err
}

// perform path transformation if needed
func (c *config) transformPaths(kvps api.KVPairs) {
	// check replacer - return immediately if not valid
	if c.pathReplacer == nil {
		return
	}

	// loop through keys
	for _, kv := range kvps {
		// split path and key with strings because
		// the path package will trim a trailing / which
		// breaks empty folders present in the kvp store
		split := strings.Split(kv.Key, consulSeparator)
		// get and check length ... only continue if we actually
		// have a path we may want to transform
		if length := len(split); length > 1 {
			// isolate and replace path
			rpath := c.pathReplacer.Replace(strings.Join(split[:length-1], consulSeparator))
			// join replaced path with key
			newKey := strings.Join([]string{rpath, split[length-1]}, consulSeparator)
			// check keys
			if kv.Key != newKey {
				// log change
				log.Printf("[Transform] %s -> %s", kv.Key, newKey)
				// update key
				kv.Key = newKey
			}
		}
	}
}
