package common

import (
	"flag"
)

// AddSharedConsulFlags adds shared flags for Consul related options
func AddSharedConsulFlags(cmdFlags *flag.FlagSet, consulConfig *ConsulConfig) {
	cmdFlags.StringVar(&consulConfig.Address, "addr", "",
		"Optional consul address and port")
	cmdFlags.StringVar(&consulConfig.Scheme, "scheme", "",
		"Optional consul scheme")
	cmdFlags.StringVar(&consulConfig.Datacenter, "dc", "",
		"Optional consul datacenter")
	cmdFlags.StringVar(&consulConfig.Token, "token", "",
		"Optional consul access token")

	// TLS Settings
	cmdFlags.StringVar(&consulConfig.CACert, "ca-cert", "",
		"Optional path to a PEM encoded CA cert file to use to verify consul")
	cmdFlags.StringVar(&consulConfig.CAPath, "ca-path", "",
		"Optional path to a directory of PEM encoded CA cert files to verify consul")
	cmdFlags.StringVar(&consulConfig.CertFile, "client-cert", "",
		"Optional path to a PEM encoded client certificate for TLS authentication to consul")
	cmdFlags.StringVar(&consulConfig.KeyFile, "client-key", "",
		"Optional path to an unencrypted PEM encoded private key matching the client certificate from -client-cert")
	cmdFlags.BoolVar(&consulConfig.InsecureSkipVerify, "tls-skip-verify", false,
		"Optional bool for verifying a TLS certificate. This is highly not recommended")
}
