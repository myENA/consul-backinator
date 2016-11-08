package common

import (
	"flag"
	"github.com/hashicorp/consul/api"
	"os"
)

// AddSharedConsulFlags adds flags shared by multiple command implementations
func AddSharedConsulFlags(cmdFlags *flag.FlagSet, consulConfig *ConsulConfig) {
	// client flags
	cmdFlags.StringVar(&consulConfig.Address, "addr", "",
		"Optional consul address and port")
	cmdFlags.StringVar(&consulConfig.Scheme, "scheme", "",
		"Optional consul scheme")
	cmdFlags.StringVar(&consulConfig.Datacenter, "dc", "",
		"Optional consul datacenter")
	cmdFlags.StringVar(&consulConfig.Token, "token", "",
		"Optional consul access token")

	// init tls struct
	consulConfig.tls = new(api.TLSConfig)

	// TLS settings
	cmdFlags.StringVar(&consulConfig.tls.CAFile, "ca-cert", "",
		"Optional path to a PEM encoded CA cert file")
	cmdFlags.StringVar(&consulConfig.tls.CertFile, "client-cert", "",
		"Optional path to a PEM encoded client certificate")
	cmdFlags.StringVar(&consulConfig.tls.KeyFile, "client-key", "",
		"Optional path to an unencrypted PEM encoded private key")
	cmdFlags.BoolVar(&consulConfig.tls.InsecureSkipVerify, "tls-skip-verify", false,
		"Optional bool for verifying a TLS certificate (not reccomended)")
}

// AddEnvDefaults attempts to populates missing config information from environment variables
func AddEnvDefaults(consulConfig *ConsulConfig) {
	// this is used in a few print statements - so we want it populated
	if consulConfig.Address == "" {
		consulConfig.Address = os.Getenv("CONSUL_HTTP_ADDR")
	}
}
