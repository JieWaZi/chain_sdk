package base

import (
	"chain_sdk/session"
	"encoding/hex"
	"fmt"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type Client struct {
	session       *session.Session
	channelClient *channel.Client
	ledgerClient  *ledger.Client
}

func NewClient(sess *session.Session, channelID string) (c *Client, err error) {
	c = &Client{session: sess}
	channelProvider := sess.ChannelContext(channelID)
	channelClient, err := channel.New(channelProvider)
	if err != nil {
		return nil, err
	}
	c.channelClient = channelClient
	ledgerClient, err := ledger.New(channelProvider)
	if err != nil {
		return nil, err
	}
	c.ledgerClient = ledgerClient
	return c, nil
}

func (c *Client) Close() {
	if c.session != nil {
		c.session.Close()
	}
}

func (c *Client) ChannelExecute(
	request channel.Request,
	options ...channel.RequestOption,
) (channel.Response, error) {
	return c.channelClient.Execute(request, options...)
}

func (c *Client) ChannelQuery(
	request channel.Request,
	options ...channel.RequestOption,
) (channel.Response, error) {
	return c.channelClient.Query(request, options...)
}

func (c *Client) QueryChannelConfig(options ...ledger.RequestOption) (fab.ChannelCfg, error) {
	return c.ledgerClient.QueryConfig(options...)
}

func (c *Client) QueryConfigBlock(options ...ledger.RequestOption) (*common.Block, error) {
	return c.ledgerClient.QueryConfigBlock(options...)
}

func (c *Client) QueryBlockchainInfo(options ...ledger.RequestOption) (*BlockchainInfo, error) {
	respFrom, err := c.ledgerClient.QueryInfo(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to call ledger QueryInfo: %v", err)
	}
	blockchainInfo := respFrom.BCI
	if blockchainInfo == nil {
		return nil, fmt.Errorf("nil blockchain info")
	}
	respTo, err := DecodeBlockchainInfo(blockchainInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to decode blockchainInfo(%+v): %v", blockchainInfo, err)
	}
	return respTo, nil
}

func (c *Client) QueryBlock(blockNumber uint64, options ...ledger.RequestOption) (*Block, error) {
	respFrom, err := c.ledgerClient.QueryBlock(blockNumber, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to call ledger QueryBlock: %v", err)
	}
	respTo, err := DecodeBlock(respFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to decode block: %v", err)
	}
	return respTo, nil
}

func (c *Client) QueryBlockByHash(blockHash string, options ...ledger.RequestOption) (*Block, error) {
	blockHashBytes, err := hex.DecodeString(blockHash)
	if err != nil {
		return nil, fmt.Errorf("failed to decode blockHash(%s): %v", blockHash, err)
	}
	respFrom, err := c.ledgerClient.QueryBlockByHash(blockHashBytes, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to call ledger QueryBlockByHash: %v", err)
	}
	respTo, err := DecodeBlock(respFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to decode block: %v", err)
	}
	return respTo, nil
}

func (c *Client) QueryBlockByTxID(txID string, options ...ledger.RequestOption) (*Block, error) {
	respFrom, err := c.ledgerClient.QueryBlockByTxID(fab.TransactionID(txID), options...)
	if err != nil {
		return nil, fmt.Errorf("failed to call ledger QueryBlockByTxID: %v", err)
	}
	respTo, err := DecodeBlock(respFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to decode block: %v", err)
	}
	return respTo, nil
}

func (c *Client) QueryTransaction(txID string, options ...ledger.RequestOption) (*ProcessedTransaction, error) {
	respFrom, err := c.ledgerClient.QueryTransaction(fab.TransactionID(txID), options...)
	if err != nil {
		return nil, fmt.Errorf("failed to call ledger QueryTransaction: %v", err)
	}
	respTo, err := DecodeProcessedTransaction(respFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to decode processedTransaction: %v", err)
	}
	return respTo, nil
}
