package common

import (
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
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

// configureHTTPClient uses the given TLS files and creates an http.Client which can be used by Consul.
// Note: In Consul 0.7.0 this can likely be replaced with api.SetupTLSConfig(tlsConfig)
func configureHTTPClient(caCert, caPath, certFile, keyFile string, insecureSkipVerify bool) (*http.Client, error) {
	httpClient := cleanhttp.DefaultClient()
	transport := cleanhttp.DefaultTransport()

	// configure a TLSConfig with the provided CAs
	tlsConfig := &tls.Config{}
	err := rootcerts.ConfigureTLS(tlsConfig, &rootcerts.Config{
		CAFile: caCert,
		CAPath: caPath,
	})
	if err != nil {
		return nil, err
	}

	// configure the TLSConfig with the provided client cert and key
	if certFile != "" && keyFile != "" {
		var err error
		clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	} else if certFile != "" || keyFile != "" {
		return nil, fmt.Errorf("Both client cert and client key must be provided")
	}

	// configure the TLSConfig with InsecureSkipVerify
	tlsConfig.InsecureSkipVerify = insecureSkipVerify

	transport.TLSClientConfig = tlsConfig
	httpClient.Transport = transport
	return httpClient, nil
}
