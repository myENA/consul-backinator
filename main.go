package main

import (
	"fmt"
	"log"
	"os"
)

// it all starts here
func main() {
	var c *config // main config
	var count int // key count
	var err error // general error holder

	// init main config
	if c, err = initConfig(); err != nil {
		log.Fatalf("[Error] Startup failed: %s", err.Error())
	}

	// doing a backup?
	if c.backupReq {
		// backup keys
		if count, err = c.backupKeys(); err != nil {
			log.Fatalf("[Error] Failed to backup data: %s", err.Error())
		}

		// show success
		log.Printf("[Success] Backed up %d keys from %s%s to %s",
			count, c.consulAddr, c.consulPrefix, c.fileName)

		// make sure they know to keep the sig
		fmt.Printf("Keep your backup (%s) and signature (%s.sig) files "+
			"in a safe place.\nYou will need both to restore your data.\n",
			c.fileName, c.fileName)

		// exit
		return
	}

	// doing a restore?
	if c.restoreReq {
		// restore keys
		if count, err = c.restoreKeys(); err != nil {
			log.Fatalf("[Error] Failed to restore data: %s", err.Error())
		}

		// show success
		log.Printf("[Success] Restored %d keys from %s to %s%s",
			count, c.fileName, c.consulAddr, c.consulPrefix)

		// exit
		return
	}

	// print usage
	fmt.Printf("Usage: %s -h\n", os.Args[0])

	// exit
	return
}
