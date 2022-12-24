package backup

import (
	"encoding/json"
	"errors"

	"github.com/hashicorp/consul/api"
	"github.com/myENA/consul-backinator/common"
)

// backupKeys fetches key/value pairs from consul and writes them to a backup file
func (c *Command) backupKeys() (int, error) {
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
	if kvps, _, err = c.consulClient.KV().List(c.config.consulPrefix, opts); err != nil {
		return 0, err
	}

	// transform paths
	c.pathTransformer.Transform(kvps)

	// set count
	count = len(kvps)

	// check count
	if count == 0 {
		return 0, errors.New("No keys found")
	}

	// encode and return
	if data, err = json.MarshalIndent(kvps, "", "  "); err != nil {
		return 0, err
	}

	// write data to destination
	if err = common.WriteData(c.config.fileName, c.config.cryptKey, data); err != nil {
		return 0, err
	}

	// return key count - no error
	return count, nil
}

// backupACLs fetches acl tokens consul and writes them to a backup file
func (c *Command) backupACLs() (int, error) {
	var aclData *common.BackupACLData // storage for acl data
	var opts *api.QueryOptions        // client query options
	var roleCount int                 // role count
	var policyCount int               // policy count
	var tokenCount int                // token count
	var data []byte
	var err error // general error holder

	// build backup acl data type
	aclData = &common.BackupACLData{
		Roles:    make([]*api.ACLRole, 0),
		Policies: make([]*api.ACLPolicy, 0),
		Tokens:   make([]*api.ACLToken, 0),
	}

	// build query options
	opts = &api.QueryOptions{
		AllowStale:        false,
		RequireConsistent: true,
	}

	// get all acl roles
	if aclData.Roles, _, err = c.consulClient.ACL().RoleList(opts); err != nil {
		return 0, err
	}
	// get all acl policies
	if policies, _, err := c.consulClient.ACL().PolicyList(opts); err != nil {
		return 0, err
	} else {
		for _, policy := range policies {
			if aclPolicy, _, err := c.consulClient.ACL().PolicyRead(policy.ID, opts); err != nil {
				return 0, err
			} else {
				aclData.Policies = append(aclData.Policies, aclPolicy)
			}
		}
	}
	// get all acl tokens
	if tokens, _, err := c.consulClient.ACL().TokenList(opts); err != nil {
		return 0, err
	} else {
		for _, token := range tokens {
			if aclToken, _, err := c.consulClient.ACL().TokenRead(token.AccessorID, opts); err != nil {
				return 0, err
			} else {
				aclData.Tokens = append(aclData.Tokens, aclToken)
			}
		}
	}

	// set counts
	roleCount = len(aclData.Roles)
	policyCount = len(aclData.Policies)
	tokenCount = len(aclData.Tokens)

	// check tokenCount
	if tokenCount == 0 && roleCount == 0 && policyCount == 0 {
		return 0, errors.New("no acl tokens, roles, or policies found")
	}

	// encode and return
	if data, err = json.MarshalIndent(aclData, "", "  "); err != nil {
		return 0, err
	}

	// write acls to destination
	if err = common.WriteData(c.config.aclFileName, c.config.cryptKey, data); err != nil {
		return 0, err
	}

	// return token tokenCount - no error
	return roleCount + policyCount + tokenCount, nil
}

// backupQueries fetches prepared query definitions from consul and writes them to a backup file
func (c *Command) backupQueries() (int, error) {
	var queries []*api.PreparedQueryDefinition // list of query definitions
	var opts *api.QueryOptions                 // client query options
	var count int                              // query count
	var data []byte                            // read definitions
	var err error                              // general error holder

	// build query options
	opts = &api.QueryOptions{
		AllowStale:        false,
		RequireConsistent: true,
	}

	// get all query definitions
	if queries, _, err = c.consulClient.PreparedQuery().List(opts); err != nil {
		return 0, err
	}

	// set count
	count = len(queries)

	// check count
	if count == 0 {
		return 0, errors.New("No query definitions found")
	}

	// encode and return
	if data, err = json.MarshalIndent(queries, "", "  "); err != nil {
		return 0, err
	}

	// write data to destination
	if err = common.WriteData(c.config.queryFileName, c.config.cryptKey, data); err != nil {
		return 0, err
	}

	// return query count - no error
	return count, nil
}
