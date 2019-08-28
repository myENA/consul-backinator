package dump

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/myENA/consul-backinator/common"
)

// dumpData reads data from a backup file and prints to stdout
func (c *Command) dumpData() error {
	var kvps api.KVPairs // kv pairs
	var data []byte      // read json data
	var err error        // general error holder

	// read json data from source
	if data, err = common.ReadData(c.config.fileName, c.config.cryptKey); err != nil {
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

	// assume kv data - plain dump really only makes sense for kv data
	if err = json.Unmarshal(data, &kvps); err != nil {
		return err
	}

	// loop through and print keys
	for _, kv := range kvps {
		fmt.Printf("Key: %s\n%s\n", kv.Key, kv.Value)
	}

	// okay
	return nil
}
