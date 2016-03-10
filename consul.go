package main

import (
	"encoding/json"
	"github.com/hashicorp/consul/api"
	"log"
	"os"
)

// build consul client
func (c *config) buildClient() error {
	var cc *api.Config // client configuration
	var err error      // general error holder

	// build default client config
	cc = api.DefaultConfig()

	// overwrite address if needed
	if c.consulAddr != "" {
		cc.Address = c.consulAddr
	}

	// overwrite scheme if needed
	if c.consulScheme != "" {
		cc.Scheme = c.consulScheme
	}

	// overwrite dc if needed
	if c.consulDc != "" {
		cc.Datacenter = c.consulDc
	}

	// overwrite token if needed
	if c.consulToken != "" {
		cc.Token = c.consulToken
	}

	// populate client wrapper
	c.consulClient, err = api.NewClient(cc)

	// return last error
	return err
}

// fetch keys from the store and write to a backup file
func (c *config) backupKeys() (int, error) {
	var kvps api.KVPairs       // list of requested kv pairs
	var opts *api.QueryOptions // client query options
	var count int              // key count
	var data []byte            // read keys
	var err error              // general error holder

	// build query options
	opts = &api.QueryOptions{
		AllowStale:        false,
		RequireConsistent: true,
	}

	// get all keys
	if kvps, _, err = c.consulClient.KV().List(c.consulPrefix, opts); err != nil {
		return count, err
	}

	// transform paths
	c.transformPaths(kvps)

	// set count
	count = len(kvps)

	// encode and return
	if data, err = json.MarshalIndent(kvps, "", "  "); err != nil {
		return count, err
	}

	// write data
	if err = c.writeBackupFile(data); err != nil {
		return count, err
	}

	// return key count - no error
	return count, nil
}

// read keys from a backup file and restore to consul
func (c *config) restoreKeys() (int, error) {
	var kvps api.KVPairs // decoded kv pairs
	var count int        // key count
	var data []byte      // read json data
	var err error        // general error holder

	// read json data from file
	if data, err = c.readBackupFile(); err != nil {
		return count, err
	}

	// doing a dump?
	if c.dataDump {
		// dump data
		if err = dumpData(data, c.plainDump); err != nil {
			return count, err
		}
		// exit clean
		os.Exit(0)
	}

	// decode data
	if err = json.Unmarshal(data, &kvps); err != nil {
		return count, err
	}

	// transform paths
	c.transformPaths(kvps)

	// set count
	count = len(kvps)

	// delete tree before restore if requested
	if c.delTree {
		// set delete prefix to passed prefix
		deletePrefix := c.consulPrefix
		// check prefix
		if c.consulPrefix == "/" {
			deletePrefix = "" // special case for root
		}
		// send the delete request
		if _, err := c.consulClient.KV().DeleteTree(deletePrefix, nil); err != nil {
			return count, err
		}
	}

	// loop through keys
	for _, kv := range kvps {
		// write key
		if _, err = c.consulClient.KV().Put(kv, nil); err != nil {
			log.Printf("[Warning] Failed to restore %s: %s",
				kv.Key, err.Error())
		}
	}

	// return key count - no error
	return count, nil
}
