package restore

import (
	"encoding/json"
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
		// filter by prefix
		if myPrefix != "" && !strings.HasPrefix(kv.Key, myPrefix) {
			continue
		}
		// write key
		if _, err = c.consulClient.KV().Put(kv, nil); err != nil {
			c.Log.Printf("[Warning] Failed to restore key %s: %s",
				kv.Key, err.Error())
		} else {
			// success - increment count
			count++
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

	// loop through acls
	for _, acl := range acls {
		// write token
		if _, _, err = c.consulClient.ACL().Create(acl, nil); err != nil {
			c.Log.Printf("[Warning] Failed to restore ACL token %s: %s",
				acl.Name, err.Error())
		} else {
			// success - increment count
			count++
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

	// loop through queries
	for _, query := range queries {
		var existing []*api.PreparedQueryDefinition // existing query definitions
		// check for existing query
		if existing, _, err = c.consulClient.PreparedQuery().Get(query.ID, nil); err != nil && existing != nil {
			// update existing query
			if _, err = c.consulClient.PreparedQuery().Update(query, nil); err != nil {
				c.Log.Printf("[Warning] Failed to update existing query definition %s: %s",
					query.ID, err.Error())
			} else {
				// success - increment count
				count++
			}
		} else {
			// remove id from backed-up query before creating
			query.ID = ""
			// attempt to create non-existent query
			if _, _, err = c.consulClient.PreparedQuery().Create(query, nil); err != nil {
				c.Log.Printf("[Warning] Failed to create missing query definition %s: %s",
					query.ID, err.Error())
			} else {
				// success - increment count
				count++
			}
		}
	}

	// return query count - no error
	return count, nil
}
