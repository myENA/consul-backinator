package consul

import (
	"crypto/tls"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
)

// Separator is the consul kvp separator
const Separator = "/"

// Config contains consul client configuration and TLSConfig in a single struct
type Config struct {
	api.Config
	TLS *api.TLSConfig
}

// Client contains a consul client implementation
type Client struct {
	*api.Client
}

// New returns an initialized consul client
func (cc *Config) New() (*Client, error) {
	var c *api.Config  // upstream client configuration
	var client *Client // client wrapper
	var err error      // general error holder

	// init upstream config
	c = api.DefaultConfig()

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

	// configure if any TLS specific options were passed
	if cc.TLS.CAFile != "" || cc.TLS.CertFile != "" || cc.TLS.KeyFile != "" || cc.TLS.InsecureSkipVerify {
		var tlsConfig *tls.Config // client TLS config
		// attempt to build tls config from passed options
		if tlsConfig, err = api.SetupTLSConfig(cc.TLS); err != nil {
			return nil, err
		}
		// build a new http client and transport
		httpClient := cleanhttp.DefaultClient()
		httpTransport := cleanhttp.DefaultTransport()
		httpTransport.TLSClientConfig = tlsConfig
		httpClient.Transport = httpTransport

		// set client
		c.HttpClient = httpClient
	}

	// init client wrapper
	client = new(Client)
	client.Client, err = api.NewClient(c)

	// return client and error
	return client, err
}
