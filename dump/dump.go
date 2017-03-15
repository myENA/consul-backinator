package dump

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/myENA/consul-backinator/common"
	"os"
)

// dumpData reads data from a backup file and prints to stdout
func (c *Command) dumpData() error {
	var kvps api.KVPairs                       // kv pairs
	var acls []*api.ACLEntry                   // acl entries
	var queries []*api.PreparedQueryDefinition // query definitions
	var data []byte                            // read json data
	var err error                              // general error holder

	// read json data from source
	if data, err = common.ReadData(kvFileName, cryptKey); err != nil {
		return err
	}

	// check plain
	if !isPlain {
		// write payload
		os.Stdout.Write(data)
		// write a blank line
		os.Stdout.WriteString("\n")
		// all done
		return nil
	}

	switch {
	case isACL:
		// decode acl data
		if err = json.Unmarshal(data, &acls); err != nil {
			return err
		}
		// loop through and print acls
		for _, acl := range acls {
			fmt.Printf("Token: %s (%s)\n%s\n", acl.Name, acl.Type, acl.Rules)
		}
	case isQuery:
		// decode acl data
		if err = json.Unmarshal(data, &queries); err != nil {
			return err
		}
		// loop through and print query definitions (not very helpful...)
		for _, query := range queries {
			fmt.Printf("Query: %s %s\n", query.ID, query.Token)
		}
	default:
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
