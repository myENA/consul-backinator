package config

import (
	"flag"
	"github.com/hashicorp/consul/api"
	ccns "github.com/myENA/consul-backinator/common/consul"
)

// AddSharedConsulFlags adds flags shared by multiple command implementations
func AddSharedConsulFlags(cmdFlags *flag.FlagSet, consulConfig *ccns.Config) {
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
	consulConfig.TLS = new(api.TLSConfig)

	// TLS settings
	cmdFlags.StringVar(&consulConfig.TLS.CAFile, "ca-cert", "",
		"Optional path to a PEM encoded CA cert file")
	cmdFlags.StringVar(&consulConfig.TLS.CertFile, "client-cert", "",
		"Optional path to a PEM encoded client certificate")
	cmdFlags.StringVar(&consulConfig.TLS.KeyFile, "client-key", "",
		"Optional path to an unencrypted PEM encoded private key")
	cmdFlags.BoolVar(&consulConfig.TLS.InsecureSkipVerify, "tls-skip-verify", false,
		"Optional bool for verifying a TLS certificate (not recommended)")
}
