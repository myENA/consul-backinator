package backup

import (
	"encoding/json"
	"errors"
	"strings"

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

	// check for exclusions
	if c.config.pathExclude != "" {
		// loop through keys
		for idx, kv := range kvps {
			match := false // reset matcher
		inner: // mark inner loop
			for _, pe := range strings.Split(c.config.pathExclude, ",") {
				if strings.HasPrefix(kv.Key, strings.TrimPrefix(pe, "/")) {
					match = true // we found a match
					break inner  // stop checking excludes
				}
			}
			if match {
				// shuffle index to end of slice
				if idx < len(kvps) {
					kvps[idx] = kvps[len(kvps)-1]
				}
				kvps = kvps[:len(kvps)-1] // remove element
			}
		}
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

// backupACLTokens fetches new (1.4.0+) consul acl tokens and writes them to a backup file
func (c *Command) backupACLTokens() (int, error) {
	var aclTokenList []*api.ACLTokenListEntry // new-style acl token entries
	var aclTokens []*api.ACLToken             // new-style acl tokens
	var opts *api.QueryOptions                // client query options
	var count int                             // token count
	var data []byte                           // read tokens
	var err error                             // general error holder

	// build query options
	opts = &api.QueryOptions{
		AllowStale:        false,
		RequireConsistent: true,
	}

	// check for new-style acl token entries
	if aclTokenList, _, err = c.consulClient.ACL().TokenList(opts); err != nil {
		return 0, err
	}

	// fetch tokens
	for _, entry := range aclTokenList {
		if token, _, err := c.consulClient.ACL().TokenRead(entry.AccessorID, opts); err == nil {
			aclTokens = append(aclTokens, token)
		}
	}

	// check count
	if count = len(aclTokens); count == 0 {
		return 0, errors.New("No tokens found")
	}

	// encode and return
	if data, err = json.MarshalIndent(aclTokens, "", "  "); err != nil {
		return 0, err
	}

	// write data to destination
	if err = common.WriteData(c.config.aclFileName, c.config.cryptKey, data); err != nil {
		return 0, err
	}

	// return token count - no error
	return count, nil
}

// backupACLPolicies fetches new (1.4.0+) consul acl policies and writes them to a backup file
func (c *Command) backupACLPolicies() (int, error) {
	var aclPolicyList []*api.ACLPolicyListEntry // new-style acl policy token entries
	var aclPolicies []*api.ACLPolicy            // new-style acl policies
	var opts *api.QueryOptions                  // client query options
	var count int                               // token count
	var data []byte                             // read tokens
	var err error                               // general error holder

	// build query options
	opts = &api.QueryOptions{
		AllowStale:        false,
		RequireConsistent: true,
	}

	// check for new-style policy entries
	if aclPolicyList, _, err = c.consulClient.ACL().PolicyList(opts); err != nil {
		return 0, err
	}

	// fetch policies
	for _, entry := range aclPolicyList {
		if policy, _, err := c.consulClient.ACL().PolicyRead(entry.ID, opts); err == nil {
			aclPolicies = append(aclPolicies, policy)
		}
	}

	// check count
	if count = len(aclPolicies); count == 0 {
		return 0, errors.New("No policies found")
	}

	// encode and return
	if data, err = json.MarshalIndent(aclPolicies, "", "  "); err != nil {
		return 0, err
	}

	// write data to destination
	if err = common.WriteData(c.config.aclPolicyFileName, c.config.cryptKey, data); err != nil {
		return 0, err
	}

	// return token count - no error
	return count, nil
}

// backupLegacyACLs fetches legacy consul acl tokens and writes them to a backup file
func (c *Command) backupLegacyACLs() (int, error) {
	var acls []*api.ACLEntry   // list of acl tokens
	var opts *api.QueryOptions // client query options
	var count int              // token count
	var data []byte            // read tokens
	var err error              // general error holder

	// build query options
	opts = &api.QueryOptions{
		AllowStale:        false,
		RequireConsistent: true,
	}

	// check for legacy acl tokens
	if acls, _, err = c.consulClient.ACL().List(opts); err != nil {
		return 0, err
	}

	// check legacy count
	if count = len(acls); count == 0 {
		return 0, errors.New("No legacy tokens found")
	}

	// encode and return
	if data, err = json.MarshalIndent(acls, "", "  "); err != nil {
		return 0, err
	}

	// write data to destination
	if err = common.WriteData(c.config.legacyACLFileName, c.config.cryptKey, data); err != nil {
		return 0, err
	}

	// return token count - no error
	return count, nil
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
