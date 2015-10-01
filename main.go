package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var inputFile = flag.String("input", "consul.bak",
		"Input file for restore operations")
	var outputFile = flag.String("output", "consul.bak",
		"Output file for backup operations")
	var encKey = flag.String("key", "password",
		"Encryption key to be used for backup and restore operations")
	var consulAddr = flag.String("address", "",
		"Consul instance address (\"127.0.0.1:8500\")")
	var consulScheme = flag.String("scheme", "",
		"Optional consul instance scheme (\"http\" or \"https\")")
	var consulDc = flag.String("dc", "",
		"Optional consul datacenter label for backup and restore")
	var consulToken = flag.String("token", "",
		"Optional consul token to access the target cluster")
	var pathPrefix = flag.String("prefix", "/",
		"Optional prefix from under which all keys will be fetched or restored")
	var pathTransform = flag.String("transform", "",
		"Optional path transformation pairs for restore (source1,dest1,source2,dest2...)")
	var deleteTree = flag.Bool("delete", false,
		"Optionally delete all keys under destination prefix before restore")
	var doBackup = flag.Bool("backup", false,
		"Trigger backup operation")
	var doRestore = flag.Bool("restore", false,
		"Trigger restore operation")

	// parse flags
	flag.Parse()

	// doing a backup?
	if *doBackup {
		// build client
		c, err := buildClient(*consulAddr, *consulScheme, *consulDc, *consulToken)

		// check error
		if err != nil {
			log.Fatalf("Failed to build consul client: %s", err.Error())
		}

		// backup keys
		count, err := c.backupKeys(*outputFile, *encKey, *pathPrefix)

		// check error
		if err != nil {
			log.Fatalf("Failed to backup KV data: %s", err.Error())
		}

		// show success
		log.Printf("Success: Backed up %d keys from %s%s to %s",
			count, *consulAddr, *pathPrefix, *outputFile)

		// exit
		return
	}

	// doing a restore?
	if *doRestore {
		// build client
		c, err := buildClient(*consulAddr, *consulScheme, *consulDc, *consulToken)

		// check error
		if err != nil {
			log.Fatalf("Failed to build consul client: %s", err.Error())
		}

		// restore keys
		count, err := c.restoreKeys(*inputFile, *encKey, *pathPrefix, *pathTransform, *deleteTree)

		// check error
		if err != nil {
			log.Fatalf("Failed to restore data: %s", err.Error())
		}

		// show success
		log.Printf("Success: Restored %d keys from %s to %s%s",
			count, *inputFile, *consulAddr, *pathPrefix)

		// exit
		return
	}

	// print usage
	fmt.Fprintf(os.Stderr, "Usage: %s -h\n", os.Args[0])

	// exit
	return
}
