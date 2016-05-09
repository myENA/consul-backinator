package common

import (
	"github.com/hashicorp/consul/api"
)

// ConsulConfig contains consul client configuration
type ConsulConfig struct {
	api.Config
}

// ConsulClient contains a consul client implementation
type ConsulClient struct {
	*api.Client
}

// New returns an initialized consul client
func (cc *ConsulConfig) New() (*ConsulClient, error) {
	var c *api.Config        // upstream client configuration
	var client *ConsulClient // client wrapper
	var err error            // general error holder

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

	// init client wrapper
	client = new(ConsulClient)
	client.Client, err = api.NewClient(c)

	// return client and error
	return client, err
}
