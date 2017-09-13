package consul

import (
	"crypto/tls"
	stdLog "log"
	"os"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-discover"
)

// package global logger
var logger *stdLog.Logger

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
func (c *Config) New() (*Client, error) {
	var ac *api.Config // upstream client configuration
	var client *Client // client wrapper
	var err error      // general error holder

	// init upstream config
	ac = api.DefaultConfig()

	// overwrite address if needed
	if c.Config.Address != "" {
		// check for cloud discovery
		if strings.Contains(c.Config.Address, "provider=") {
			var addrs []string // local holder
			// attempt service discovery
			if addrs, err = new(discover.Discover).Addrs(c.Config.Address, logger); err != nil {
				logger.Printf("[Error] Failed to disover cluster address: %s", err.Error())
				return nil, err
			}
			// select first returned server and let them know
			ac.Address = addrs[0]
			logger.Printf("[Info] Using %s for cluster address.", ac.Address)
		} else {
			// no discovery - pass on as set
			ac.Address = c.Config.Address
		}
	}

	// overwrite scheme if needed
	if c.Config.Scheme != "" {
		ac.Scheme = c.Config.Scheme
	}

	// overwrite dc if needed
	if c.Config.Datacenter != "" {
		ac.Datacenter = c.Config.Datacenter
	}

	// overwrite token if needed
	if c.Config.Token != "" {
		ac.Token = c.Config.Token
	}

	// configure if any TLS specific options were passed
	if c.TLS.CAFile != "" || c.TLS.CertFile != "" ||
		c.TLS.KeyFile != "" || c.TLS.InsecureSkipVerify {
		var tlsConfig *tls.Config // client TLS config
		// attempt to build tls config from passed options
		if tlsConfig, err = api.SetupTLSConfig(c.TLS); err != nil {
			return nil, err
		}
		// build a new http client and transport for each api client
		httpClient := cleanhttp.DefaultClient()
		httpTransport := cleanhttp.DefaultTransport()
		httpTransport.TLSClientConfig = tlsConfig
		httpClient.Transport = httpTransport
		// set client
		ac.HttpClient = httpClient
	}

	// init client wrapper
	client = new(Client)
	client.Client, err = api.NewClient(ac)

	// return client and error
	return client, err
}

func init() {
	logger = stdLog.New(os.Stderr, "", stdLog.LstdFlags)
}
