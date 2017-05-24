package restore

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/myENA/consul-backinator/common"
)

// restoreKeys reads keys from a backup file and restores them to consul
func (c *Command) restoreKeys() (int, error) {
	var kvps api.KVPairs // decoded kv pairs
	var count int        // key count
	var data []byte      // read json data
	var err error        // general error holder

	// read json data from source
	if data, err = common.ReadData(c.config.fileName, c.config.cryptKey); err != nil {
		return 0, err
	}

	// decode data
	if err = json.Unmarshal(data, &kvps); err != nil {
		return 0, err
	}

	// transform paths
	c.pathTransformer.Transform(kvps)

	// set count
	count = len(kvps)

	// set to passed prefix
	myPrefix := c.config.consulPrefix
	// check prefix
	if c.config.consulPrefix == "/" {
		myPrefix = "" // special case for root
	}

	// delete tree before restore if requested
	if c.config.delTree {

		// send the delete request
		if _, err := c.consulClient.KV().DeleteTree(myPrefix, nil); err != nil {
			return count, err
		}
	}

	// loop through keys
	for _, kv := range kvps {
		// filter prefix
		if myPrefix != "" && !strings.HasPrefix(kv.Key, myPrefix) {
			continue
		}
		// write key
		if _, err = c.consulClient.KV().Put(kv, nil); err != nil {
			log.Printf("[Warning] Failed to restore key %s: %s",
				kv.Key, err.Error())
		}
	}

	// return key count - no error
	return count, nil
}

// restoreACLs reads acl tokens from a backup file and restores them to consul
func (c *Command) restoreACLs() (int, error) {
	var acls []*api.ACLEntry // acl tokens
	var count int            // token count
	var data []byte          // read json data
	var err error            // general error holder

	// read json data from source
	if data, err = common.ReadData(c.config.aclFileName, c.config.cryptKey); err != nil {
		return 0, err
	}

	// decode data
	if err = json.Unmarshal(data, &acls); err != nil {
		return 0, err
	}

	// set count
	count = len(acls)

	// loop through acls
	for _, acl := range acls {
		// write token
		if _, _, err = c.consulClient.ACL().Create(acl, nil); err != nil {
			log.Printf("[Warning] Failed to restore ACL token %s: %s",
				acl.Name, err.Error())
		}
	}

	// return acl count - no error
	return count, nil
}

// restoreQueries reads query definitions from a backup file and restores them to consul
func (c *Command) restoreQueries() (int, error) {
	var queries []*api.PreparedQueryDefinition // query definitions
	var count int                              // query count
	var data []byte                            // read json data
	var err error                              // general error holder

	// read json data from source
	if data, err = common.ReadData(c.config.queryFileName, c.config.cryptKey); err != nil {
		return 0, err
	}

	// decode data
	if err = json.Unmarshal(data, &queries); err != nil {
		return 0, err
	}

	// set count
	count = len(queries)

	// loop through acls
	for _, query := range queries {
		// write query definitions
		if _, _, err = c.consulClient.PreparedQuery().Create(query, nil); err != nil {
			log.Printf("[Warning] Failed to restore query definition %s: %s",
				query.ID, err.Error())
		}
	}

	// return query count - no error
	return count, nil
}
