package client

import (
	"chain_sdk/base"
	"chain_sdk/session"
	"encoding/json"

	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/util/pathvar"
	"github.com/sirupsen/logrus"
)

type Client struct {
	*base.Client    `env:"-"`
	ConfigPath      string   `env:""`
	Username        string   `env:""`
	Organization    string   `env:""`
	ChannelID       string   `env:""`
	Chaincode       string   `env:""`
	TargetEndPoint  []string `env:""`
	OrdererEndPoint string   `env:""`
}

func (c *Client) Init() {
	if len(c.ConfigPath) == 0 {
		return
	}
	var options []session.ClientOption
	options = append(options, session.WithUser(c.Username))
	options = append(options, session.WithOrg(c.Organization))
	options = append(options, session.WithTargetEndpoint(c.TargetEndPoint))
	options = append(options, session.WithOrdererEndpoint(c.OrdererEndPoint))
	sess, err := session.New(pathvar.Subst(c.ConfigPath), options...)
	if err != nil {
		logrus.Fatalf("failed to new session: %v", err)
	}
	baseClient, err := base.NewClient(sess, c.ChannelID)
	if err != nil {
		logrus.Fatalf("failed to new baseClient: %v", err)
	}
	c.Client = baseClient
}

func (c *Client) CreateVehicle(args [][]byte) (string, error) {
	result, err := c.ChannelExecute(
		channel.Request{ChaincodeID: c.Chaincode, Fcn: "createVehicle", Args: args},
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetEndpoints(c.TargetEndPoint...),
	)
	if err != nil {
		return "", err
	}
	return string(result.TransactionID), nil
}

func (c *Client) FindVehicle(args [][]byte) (vehicle *Vehicle, err error) {
	result, err := c.ChannelQuery(
		channel.Request{ChaincodeID: c.Chaincode, Fcn: "findVehicle", Args: args},
		channel.WithRetry(retry.DefaultChannelOpts),
		channel.WithTargetEndpoints(c.TargetEndPoint...),
	)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(result.Payload, &vehicle)
	return
}