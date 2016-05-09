package restore

import (
	"encoding/json"
	"github.com/hashicorp/consul/api"
	"github.com/myENA/consul-backinator/common"
	"log"
)

// read keys from a backup file and restore to consul
func (c *Command) restoreKeys() (int, error) {
	var kvps api.KVPairs // decoded kv pairs
	var count int        // key count
	var data []byte      // read json data
	var err error        // general error holder

	// read json data from file
	if data, err = common.ReadFile(c.config.fileName, c.config.cryptKey); err != nil {
		return count, err
	}

	// decode data
	if err = json.Unmarshal(data, &kvps); err != nil {
		return count, err
	}

	// transform paths
	c.pathTransformer.Transform(kvps)

	// set count
	count = len(kvps)

	// delete tree before restore if requested
	if c.config.delTree {
		// set delete prefix to passed prefix
		deletePrefix := c.config.consulPrefix
		// check prefix
		if c.config.consulPrefix == "/" {
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
