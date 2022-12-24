package common

import "github.com/hashicorp/consul/api"

type BackupACLData struct {
	Roles    []*api.ACLRole
	Policies []*api.ACLPolicy
	Tokens   []*api.ACLToken
}
