package main

import (
	"encoding/json"
	"github.com/hashicorp/consul/api"
)

type client struct {
	*api.Client
}

func buildClient(address, scheme, dc, token string) (*client, error) {
	// build default config
	config := api.DefaultConfig()
	// overwrite address if needed
	if address != "" {
		config.Address = address
	}
	// overwrite scheme if needed
	if scheme != "" {
		config.Scheme = scheme
	}
	// overwrite dc if needed
	if dc != "" {
		config.Datacenter = dc
	}
	// overwrite token if needed
	if token != "" {
		config.Token = token
	}
	// build client
	c, err := api.NewClient(config)
	// return wrapped client
	return &client{Client: c}, err
}

func (c *client) getKeys(prefix string) ([]byte, error) {
	// build options
	opts := &api.QueryOptions{
		RequireConsistent: true,
	}
	// get all keys
	pairs, _, err := c.KV().List(prefix, opts)
	// check error
	if err != nil {
		return nil, err
	}
	// encode and return
	return json.Marshal(pairs)
}
