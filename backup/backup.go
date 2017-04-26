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
	if kvps, _, err = c.consulClient.KV().List(consulPrefix, opts); err != nil {
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
	if err = common.WriteData(kvFileName, cryptKey, data); err != nil {
		return 0, err
	}

	// return key count - no error
	return count, nil
}

// backupACLs fetches acl tokens consul and writes them to a backup file
func (c *Command) backupACLs() (int, error) {
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

	// get all acl tokens
	if acls, _, err = c.consulClient.ACL().List(opts); err != nil {
		return 0, err
	}

	// set count
	count = len(acls)

	// check count
	if count == 0 {
		return 0, errors.New("No tokens found")
	}

	// encode and return
	if data, err = json.MarshalIndent(acls, "", "  "); err != nil {
		return 0, err
	}

	// write data to destination
	if err = common.WriteData(aclFileName, cryptKey, data); err != nil {
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
	if err = common.WriteData(queryFileName, cryptKey, data); err != nil {
		return 0, err
	}

	// return query count - no error
	return count, nil
}
