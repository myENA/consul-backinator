package restore

import (
	"encoding/json"
	"github.com/hashicorp/consul/api"
	"github.com/myENA/consul-backinator/common"
	"log"
)

// read keys from a backup file and restore to consul
func (c *Command) restoreKeys() (int, error) {
	var kvps api.KVPairs  // decoded kv pairs
	var count int         // key count
	var data []byte       // read json data
	var s3 *common.S3Info // s3 info struct
	var err error         // general error holder

	// check filename
	if s3, err = common.GetS3Info(c.config.fileName); err == nil {
		// read data from s3
		if data, err = s3.Read(c.config.cryptKey); err != nil {
			return 0, err
		}
	} else {
		// read json data from file
		if data, err = common.ReadFile(c.config.fileName, c.config.cryptKey); err != nil {
			return 0, err
		}
	}

	// decode data
	if err = json.Unmarshal(data, &kvps); err != nil {
		return 0, err
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
			log.Printf("[Warning] Failed to restore key %s: %s",
				kv.Key, err.Error())
		}
	}

	// return key count - no error
	return count, nil
}

// read acl tokens from a backup file and restore to consul
func (c *Command) restoreAcls() (int, error) {
	var acls []*api.ACLEntry // decoded acl tokens
	var count int            // key count
	var data []byte          // read json data
	var s3 *common.S3Info    // s3 info struct
	var err error            // general error holder

	// check filename
	if s3, err = common.GetS3Info(c.config.aclFileName); err == nil {
		// read data from s3
		if data, err = s3.Read(c.config.cryptKey); err != nil {
			return 0, err
		}
	} else {
		// read json data from file
		if data, err = common.ReadFile(c.config.aclFileName, c.config.cryptKey); err != nil {
			return 0, err
		}
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
