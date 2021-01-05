package base

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/util/pathvar"
	"log"
)

type ChainClient struct {
	*Client         `env:"-"`
	ConfigPath      string   `env:""`
	Username        string   `env:""`
	Organization    string   `env:""`
	ChannelID       string   `env:""`
	ChainCode       string   `env:""`
	TargetEndPoint  []string `env:""`
	OrdererEndPoint string   `env:""`
}

func (c *ChainClient) Init() {
	if len(c.ConfigPath) == 0 {
		return
	}
	var options []ClientOption
	options = append(options,
		WithUser(c.Username),
		WithOrg(c.Organization),
		WithTargetEndpoint(c.TargetEndPoint),
		WithOrdererEndpoint(c.OrdererEndPoint))
	sess, err := New(pathvar.Subst(c.ConfigPath), options...)
	if err != nil {
		log.Fatalf("create session err:%s", err.Error())
	}
	baseClient, err := NewClient(sess, c.ChannelID)
	if err != nil {
		log.Fatalf("new client err:%s", err.Error())
	}
	c.Client = baseClient
}
