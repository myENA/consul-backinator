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
	var aclData common.BackupACLData // acl tokens
	var roleCount int                // role count
	var policyCount int              // policy count
	var tokenCount int               // token count
	var data []byte                  // read json data
	var err error                    // general error holder

	// read json data from source
	if data, err = common.ReadData(c.config.aclFileName, c.config.cryptKey); err != nil {
		return 0, err
	}

	// decode data
	if err = json.Unmarshal(data, &aclData); err != nil {
		return 0, err
	}

	// loop through acl roles
	for _, role := range aclData.Roles {
		// write role
		if _, _, err = c.consulClient.ACL().RoleCreate(role, nil); err != nil {
			c.Log.Printf("[Warning] Failed to restore ACL role %s: %v",
				role.Name, err)
		} else {
			// success - increment count
			roleCount++
		}
	}

	// loop through acl policies
	for _, policy := range aclData.Policies {
		if _, _, err := c.consulClient.ACL().PolicyCreate(policy, nil); err != nil {
			c.Log.Printf("[Warning] Failed to restore ACL policy %s: %v",
				policy.Name, err)
		} else {
			// success - increment count
			policyCount++
		}
	}

	// loop through acl tokens
	for _, token := range aclData.Tokens {
		// write token
		if _, _, err = c.consulClient.ACL().TokenCreate(token, nil); err != nil {
			c.Log.Printf("[Warning] Failed to restore ACL token with description %q: %v",
				token.Description, err)
		} else {
			// success - increment count
			tokenCount++
		}
	}

	// return restored acl items count - no error
	return roleCount + policyCount + tokenCount, nil
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
