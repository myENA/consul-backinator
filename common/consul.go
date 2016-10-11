package common

import (
	"crypto/tls"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-rootcerts"
	"net/http"
)

// ConsulConfig contains consul client configuration
type ConsulConfig struct {
	api.Config
	// Can be merged into api.TLSConfig in Consul >= 0.7.0
	CACert             string
	CAPath             string
	CertFile           string
	KeyFile            string
	InsecureSkipVerify bool
}

// ConsulClient contains a consul client implementation
type ConsulClient struct {
	*api.Client
}

// NewClient returns an initialized consul client
func (cc *ConsulConfig) NewClient() (*ConsulClient, error) {
	var c *api.Config        // upstream client configuration
	var client *ConsulClient // client wrapper
	var err error            // general error holder

	// init upstream config
	c = api.DefaultConfig()

	// configure an http.Client for TLS if any configs were set
	if cc.CACert != "" || cc.CAPath != "" || cc.CertFile != "" || cc.KeyFile != "" || cc.InsecureSkipVerify {
		httpClient, err := configureHTTPClient(cc.CACert, cc.CAPath, cc.CertFile, cc.KeyFile, cc.InsecureSkipVerify)
		if err != nil {
			return nil, err
		}
		c.HttpClient = httpClient
	}

	// overwrite address if needed
	if cc.Address != "" {
		c.Address = cc.Address
	}

	// overwrite scheme if needed
	if cc.Scheme != "" {
		c.Scheme = cc.Scheme
	}

	// overwrite dc if needed
	if cc.Datacenter != "" {
		c.Datacenter = cc.Datacenter
	}

	// overwrite token if needed
	if cc.Token != "" {
		c.Token = cc.Token
	}

	// init client wrapper
	client = new(ConsulClient)
	client.Client, err = api.NewClient(c)

	// return client and error
	return client, err
}

func configureHTTPClient(caCert, caPath, certFile, keyFile string, insecureSkipVerify bool) (*http.Client, error) {
	httpClient := new(http.Client)

	if caCert != "" && caPath != "" {

	}

	return httpClient, nil
}
