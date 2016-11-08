package backup

import (
	"encoding/json"
	"errors"
	"github.com/hashicorp/consul/api"
	"github.com/myENA/consul-backinator/common"
)

// fetch keys from the store and write to a backup file
func (c *Command) backupKeys() (int, error) {
	var kvps api.KVPairs       // list of requested kv pairs
	var opts *api.QueryOptions // client query options
	var count int              // key count
	var data []byte            // read keys
	var s3 *common.S3Info      // s3 info struct
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

	// check options
	if s3, err = common.GetS3Info(c.config.fileName); err == nil {
		// write data to s3
		if err = s3.Write(c.config.cryptKey, data); err != nil {
			return 0, err
		}
	} else {
		// write data to file
		if err = common.WriteFile(c.config.fileName, c.config.cryptKey, data); err != nil {
			return 0, err
		}
	}

	// return key count - no error
	return count, nil
}

// fetch acl tokens from the cluster and write to a backup file
func (c *Command) backupAcls() (int, error) {
	var acls []*api.ACLEntry   // list of acl tokens
	var opts *api.QueryOptions // client query options
	var count int              // key count
	var data []byte            // read keys
	var s3 *common.S3Info      // s3 info struct
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

	// check options
	if s3, err = common.GetS3Info(c.config.aclFileName); err == nil {
		// write data to s3
		if err = s3.Write(c.config.cryptKey, data); err != nil {
			return 0, err
		}
	} else {
		// write data to file
		if err = common.WriteFile(c.config.aclFileName, c.config.cryptKey, data); err != nil {
			return 0, err
		}
	}

	// return token count - no error
	return count, nil
}
