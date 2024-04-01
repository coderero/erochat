package cassd

import (
	"time"

	"github.com/gocql/gocql"
)

func NewSession(host, username, password, keyspace string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(host)
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 10
	cluster.Timeout = time.Second * 10
	cluster.MaxWaitSchemaAgreement = time.Second * 10
	cluster.NumConns = 10
	cluster.DisableInitialHostLookup = true
	cluster.IgnorePeerAddr = true
	cluster.ReconnectInterval = time.Second * 10
	cluster.MaxPreparedStmts = 1000
	cluster.MaxRoutingKeyInfo = 1000
	cluster.PageSize = 5000
	cluster.SerialConsistency = gocql.LocalSerial
	cluster.SocketKeepalive = time.Second * 10

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}
