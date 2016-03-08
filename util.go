package main

import (
	//"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
	"strings"
)

// consul kvp separator
const consulSeparator = "/"

// our instance configuration
type config struct {
	fileName      string
	cryptKey      string
	pathTransform string
	pathReplacer  *strings.Replacer
	delTree       bool
	backupReq     bool
	restoreReq    bool
	dataDump      bool
	plainDump     bool
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
	flag.StringVar(&c.fileName, "file", "consul.bak",
		"File for backup and restore operations")
	flag.StringVar(&c.cryptKey, "key", "password",
		"Passphrase used for data encryption and signature validation")
	flag.StringVar(&c.pathTransform, "transform", "",
		"Optional folder path transformation (oldPath,newPath...)")
	flag.BoolVar(&c.dataDump, "dump", false,
		"Dump backup file contents to stdout and exit when used with "+
			"the -restore option")
	flag.BoolVar(&c.plainDump, "plain", false,
		"Dump only the key and decoded value to stdout when used with "+
			"the -restore and -dump options")
	flag.BoolVar(&c.delTree, "delete", false,
		"Delete all keys under the destination prefix before restore")
	flag.BoolVar(&c.backupReq, "backup", false,
		"Trigger backup operation")
	flag.BoolVar(&c.restoreReq, "restore", false,
		"Trigger restore operation")
	flag.StringVar(&c.consulAddr, "addr", "",
		"Optional consul instance address and port (\"127.0.0.1:8500\")")
	flag.StringVar(&c.consulScheme, "scheme", "",
		"Optional consul instance scheme (\"http\" or \"https\")")
	flag.StringVar(&c.consulDc, "dc", "",
		"Optional consul datacenter label for backup and restore")
	flag.StringVar(&c.consulToken, "token", "",
		"Optional consul token to access the target cluster")
	flag.StringVar(&c.consulPrefix, "prefix", "/",
		"Optional prefix from under which all keys will be fetched")

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

// dump decoded data
func dumpData(data []byte, plain bool) error {
	var kvps api.KVPairs // decoded kv pairs
	//var dd []byte        // decoded data
	var err error // general error holder

	if !plain {
		// write payload
		os.Stdout.Write(data)
		// write a blank line
		os.Stdout.WriteString("\n")
		// all done
		return nil
	}

	// decode data
	if err = json.Unmarshal(data, &kvps); err != nil {
		return err
	}

	// loop through and print data
	for _, kv := range kvps {
		fmt.Printf("Key: %s\n%s\n", kv.Key, kv.Value)
	}

	// okay
	return nil
}
