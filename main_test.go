package main_test

import (
	"io/ioutil"
	stdLog "log"
	"os"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/sdk/testutil"
	"github.com/mitchellh/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/myENA/consul-backinator/command/backup"
	"github.com/myENA/consul-backinator/command/restore"
)

const (
	MyAwesomeToken = "c2011091e592a41d557b425c4da65241fce12c0c"
	MySecretKey    = "CorrectHorseBatteryStaple"
	appName        = "consul-backinator" // normally from version.go
	appVersion     = "test"              // normally from version.go
)

type BackinatorTestSuite struct {
	suite.Suite
	TestSource                  *testutil.TestServer
	TestTarget                  *testutil.TestServer
	TestSourceClient            *api.Client
	TestSourceClientConfig      *api.Config
	TestACLEntry                *api.ACLEntry
	TestPreparedQueryDefinition *api.PreparedQueryDefinition
	TestKeyFile                 string
	TestACLFile                 string
	TestQueryFile               string
}

func mktemp(prefix string) string {
	var file *os.File // temp file
	var err error     // error holder
	if file, err = ioutil.TempFile(os.TempDir(), prefix); err != nil {
		panic(err.Error())
	}
	defer file.Close()
	return file.Name()
}

func (suite *BackinatorTestSuite) SetupSuite() {
	var err error      // error holder
	var aclID string   // acl identifier
	var queryID string // query identifier

	// create source consul server
	if suite.TestSource, err = testutil.NewTestServerConfigT(
		suite.T(),
		func(c *testutil.TestServerConfig) {
			c.Datacenter = "test-source"
			c.ACLDatacenter = "test-source"
			c.ACLDefaultPolicy = "allow"
			c.ACLMasterToken = MyAwesomeToken
		}); err != nil {
		suite.T().Fatal(err)
	}

	// create target consul server
	if suite.TestTarget, err = testutil.NewTestServerConfigT(
		suite.T(),
		func(c *testutil.TestServerConfig) {
			c.Datacenter = "test-target"
			c.ACLDatacenter = "test-target"
			c.ACLDefaultPolicy = "allow"
			c.ACLMasterToken = MyAwesomeToken
		}); err != nil {
		suite.T().Fatal(err)
	}

	// build api client config
	suite.TestSourceClientConfig = api.DefaultConfig()
	suite.TestSourceClientConfig.Address = suite.TestSource.HTTPAddr
	suite.TestSourceClientConfig.Datacenter = suite.TestSource.Config.Datacenter
	suite.TestSourceClientConfig.Token = MyAwesomeToken

	// create api client
	if suite.TestSourceClient, err = api.NewClient(suite.TestSourceClientConfig); err != nil {
		suite.T().Fatal(err)
	}

	// populate dummy acl entry
	suite.TestACLEntry = &api.ACLEntry{
		Name:  "myCustomACL",
		Type:  "client",
		Rules: `key "" { policy = "read" } key "foo/" { policy = "write" } key "foo/private/" { policy = "deny" } operator = "read"`,
	}

	// push a custom acl
	aclID, _, err = suite.TestSourceClient.ACL().Create(suite.TestACLEntry, nil)

	// check return
	assert.NoError(suite.T(), err, "api acl operation returned error")
	assert.NotEmpty(suite.T(), aclID, "api acl operation returned empty")

	// populate dummy prepared query
	suite.TestPreparedQueryDefinition = &api.PreparedQueryDefinition{
		Name: "myCustomQuery",
		Service: api.ServiceQuery{
			Service: "testService",
		},
		DNS: api.QueryDNSOptions{
			TTL: "10s",
		},
	}

	// push the prepared query
	queryID, _, err = suite.TestSourceClient.PreparedQuery().Create(suite.TestPreparedQueryDefinition, nil)

	// check return
	assert.NoError(suite.T(), err, "api query operation returned error")
	assert.NotEmpty(suite.T(), queryID, "api query operation returned empty")

	// populate source kv data
	suite.TestSource.SetKVString(suite.T(), "key1", "value1")
	suite.TestSource.SetKVString(suite.T(), "folder1/key2", "value2")
	suite.TestSource.SetKVString(suite.T(), "folder2/key3", "value3")

	// check keys
	assert.Equal(suite.T(),
		"value1", suite.TestSource.GetKVString(suite.T(), "/key1"),
		"value2", suite.TestSource.GetKVString(suite.T(), "/folder1/key2"),
		"value3", suite.TestSource.GetKVString(suite.T(), "/folder2/key3"),
	)

	// setup temporary files
	suite.TestACLFile = mktemp(appName + ".acls")
	suite.TestQueryFile = mktemp(appName + ".pqs")
	suite.TestKeyFile = mktemp(appName + ".bak")
}

func (suite *BackinatorTestSuite) TearDownSuite() {
	suite.T().Log("Shutting down consul servers...")
	suite.TestSource.Stop() // stop source consul server
	suite.TestTarget.Stop() // stop target consul server
	suite.T().Log("Done!")
	suite.T().Log("Removing temporary files ...")
	os.Remove(suite.TestKeyFile)
	os.Remove(suite.TestKeyFile + ".sig")
	os.Remove(suite.TestACLFile)
	os.Remove(suite.TestACLFile + ".sig")
	os.Remove(suite.TestQueryFile)
	os.Remove(suite.TestQueryFile + ".sig")
	suite.T().Log("Done!")
}

func (suite *BackinatorTestSuite) Test01Backup() {
	var c *cli.CLI // cli object
	var status int // exit status
	var err error  // error holder

	// init and populate cli object
	c = cli.NewCLI(appName, appVersion)
	c.Args = []string{
		"backup",
		"-file",
		suite.TestKeyFile,
		"-key",
		MySecretKey,
		"-acls",
		suite.TestACLFile,
		"-queries",
		suite.TestQueryFile,
		"-addr",
		suite.TestSource.HTTPAddr,
		"-dc",
		suite.TestSource.Config.Datacenter,
		"-token",
		MyAwesomeToken,
	}
	c.Commands = map[string]cli.CommandFactory{
		"backup": func() (cli.Command, error) {
			return &backup.Command{
				Self: "test-backup",
				Log:  stdLog.New(os.Stderr, "", stdLog.LstdFlags),
			}, nil
		},
	}
	// run command
	status, err = c.Run()

	// check results
	assert.NoError(suite.T(), err, "operation returned error")
	assert.Equal(suite.T(), status, 0, "operation exited non-zero")
}

func (suite *BackinatorTestSuite) Test02Restore() {
	var c *cli.CLI // cli object
	var status int // exit status
	var err error  // error holder

	// init and populate cli object
	c = cli.NewCLI(appName, appVersion)
	c.Args = []string{
		"restore",
		"-file",
		suite.TestKeyFile,
		"-key",
		MySecretKey,
		"-acls",
		suite.TestACLFile,
		"-queries",
		suite.TestQueryFile,
		"-addr",
		suite.TestTarget.HTTPAddr,
		"-dc",
		suite.TestTarget.Config.Datacenter,
		"-token",
		MyAwesomeToken,
	}
	c.Commands = map[string]cli.CommandFactory{
		"restore": func() (cli.Command, error) {
			return &restore.Command{
				Self: "test-restore",
				Log:  stdLog.New(os.Stderr, "", stdLog.LstdFlags),
			}, nil
		},
	}
	// run command
	status, err = c.Run()

	// check results
	assert.NoError(suite.T(), err, "operation returned error")
	assert.Equal(suite.T(), status, 0, "operation exited non-zero")
}

func (suite *BackinatorTestSuite) Test03VerifyTarget() {
	// TODO Read data from target server and verify
	// it matches the original source
}

func TestBackinatorTestSuite(t *testing.T) {
	suite.Run(t, new(BackinatorTestSuite))
}
