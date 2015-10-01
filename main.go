package main

import (
	"flag"
	"log"
)

func main() {
	var inputFile = flag.String("input", "consul.bak", "Input file")
	var outputFile = flag.String("output", "consul.bak", "Output file")
	var encKey = flag.String("key", "password", "Encryption key")
	var consulAddr = flag.String("address", "", "Consul instance address")
	var consulScheme = flag.String("scheme", "", "Consul instance scheme (http or https")
	var consulDc = flag.String("dc", "", "Consul datacenter")
	var consulToken = flag.String("token", "", "Consul access token")
	var pathPrefix = flag.String("prefix", "/", "Path prefix")
	//	var pathRegex = flag.String("regex", "", "Regex applied to data paths")
	var doBackup = flag.Bool("backup", false, "Request backup")
	var doRestore = flag.Bool("restore", false, "Request restoration")

	flag.Parse()

	// doing a backup?
	if *doBackup {
		// build client
		c, err := buildClient(*consulAddr, *consulScheme, *consulDc, *consulToken)

		// get data from consul api
		data, err := c.getKeys(*pathPrefix)

		// check error
		if err != nil {
			log.Printf("Failed to fetch keys: %s", err.Error())
		}

		// write data
		if err := writeFile(*outputFile, data, buildKey(*encKey)); err != nil {
			log.Fatalf("Failed to write backup data: %s", err.Error())
		}

		log.Printf("Backed up xxx keys to %s", *outputFile)

		// exit
		return
	}

	// doing a restore?
	if *doRestore {
		bytes, err := readFile(*inputFile, buildKey(*encKey))
		if err != nil {
			log.Fatalf("Failed to read restore data: %s", err.Error())
		}
		log.Printf("Read: %s", string(bytes))
		// exit
		return
	}

	// exit
	return
}
