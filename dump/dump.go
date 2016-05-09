package dump

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/myENA/consul-backinator/common"
	"os"
)

// read keys from a backup file and restore to consul
func (c *Command) dumpData() error {
	var kvps api.KVPairs // decoded kv pairs
	var data []byte      // read json data
	var err error        // general error holder

	// read json data from file
	if data, err = common.ReadFile(c.config.fileName, c.config.cryptKey); err != nil {
		return err
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

	// decode data
	if err = json.Unmarshal(data, &kvps); err != nil {
		return err
	}

	// loop through and print data
	for _, kv := range kvps {
		fmt.Printf("Key: %s\n%s\n", kv.Key, kv.Value)
	}

	// okay
	return nil
}
