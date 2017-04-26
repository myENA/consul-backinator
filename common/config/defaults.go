package config

import (
	"os"

	ccns "github.com/myENA/consul-backinator/common/consul"
)

// AddEnvDefaults attempts to populates missing config information from environment variables
func AddEnvDefaults(consulConfig *ccns.Config) {
	// this is used in a few print statements - so we want it populated
	if consulConfig.Address == "" {
		consulConfig.Address = os.Getenv("CONSUL_HTTP_ADDR")
	}
}
