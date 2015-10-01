package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
	"strings"
)

// consul kvp separator
const consulSeparator = "/"

func (c *client) backupKeys(file, key, prefix string) (int, error) {
	// get data and key count from consul api
	data, count, err := c.getKeys(prefix)

	// check error
	if err != nil {
		return 0, err
	}

	// write data
	if err := writeFile(file, data, buildKey(key)); err != nil {
		return 0, err
	}

	// return key count - no error
	return count, nil
}

func (c *client) restoreKeys(file, key, prefix, transform string, deltree bool) (int, error) {
	var replacer *strings.Replacer
	var kvps api.KVPairs

	// restore keys
	bytes, err := readFile(file, buildKey(key))

	// check error
	if err != nil {
		return 0, err
	}

	// check transformation request
	if transform != "" {
		// split strings
		split := strings.Split(transform, ",")
		// check count
		if (len(split) % 2) != 0 {
			return 0, fmt.Errorf("Odd transform count: %d", len(split))
		}
		// build replacer
		replacer = strings.NewReplacer(split...)
	}

	// decode data
	if err := json.Unmarshal(bytes, &kvps); err != nil {
		return 0, err
	}

	// delete tree before restore if requested
	if deltree {
		if _, err := c.KV().DeleteTree(prefix, nil); err != nil {
			return 0, err
		}
	}

	// loop through keys
	for _, kv := range kvps {
		// transform path if we have a valid replacer
		if replacer != nil {
			// split path and key with strings because
			// the path package will trim a trailing / which
			// breaks empty folders present in the kvp store
			split := strings.Split(kv.Key, consulSeparator)
			// get and check length ... only continue if we actually
			// have a path we may want to transform
			if length := len(split); length > 1 {
				// isolate and replace path
				rpath := replacer.Replace(strings.Join(split[:length-1], consulSeparator))
				// join replaced path with key and update
				kv.Key = strings.Join([]string{rpath, split[length-1]}, consulSeparator)
			}
		}
		// attempt write key
		if _, err = c.KV().Put(kv, nil); err != nil {
			log.Printf("WARNING: Failed to restore %s: %s",
				kv.Key, err.Error())
		}
	}

	// return key count - no error
	return len(kvps), nil
}
