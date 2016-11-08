package dump

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/myENA/consul-backinator/common"
	"os"
)

// read data from a backup file
func (c *Command) dumpData() error {
	var kvps api.KVPairs     // decoded kv pairs
	var acls []*api.ACLEntry // decoded acl entries
	var s3 *common.S3Info    // s3 info struct
	var data []byte          // read json data
	var err error            // general error holder

	// check filename
	if s3, err = common.GetS3Info(c.config.fileName); err == nil {
		// read data from s3
		if data, err = s3.Read(c.config.cryptKey); err != nil {
			return err
		}
	} else {
		// read json data from file
		if data, err = common.ReadFile(c.config.fileName, c.config.cryptKey); err != nil {
			return err
		}
	}

	// check plain
	if !c.config.plainDump {
		// write payload
		os.Stdout.Write(data)
		// write a blank line
		os.Stdout.WriteString("\n")
		// all done
		return nil
	}

	// acls
	if c.config.acls {
		// decode acl data
		if err = json.Unmarshal(data, &acls); err != nil {
			return err
		}
		// loop through and print acls
		for _, acl := range acls {
			fmt.Printf("Token: %s (%s)\n%s\n", acl.Name, acl.Type, acl.Rules)
		}
	} else {
		// decode kv data
		if err = json.Unmarshal(data, &kvps); err != nil {
			return err
		}
		// loop through and print keys
		for _, kv := range kvps {
			fmt.Printf("Key: %s\n%s\n", kv.Key, kv.Value)
		}
	}

	// okay
	return nil
}
