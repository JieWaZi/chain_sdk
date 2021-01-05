package base

import (
	"encoding/hex"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel/invoke"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type Client struct {
	session       *Session
	channelClient *channel.Client
	ledgerClient  *ledger.Client
	resmgmtClient *resmgmt.Client
}

func NewClient(sess *Session, channelID string) (c *Client, err error) {
	c = &Client{session: sess}
	resmgmtClient, err := resmgmt.New(sess.ClientProvider())
	if err != nil {
		return nil, fmt.Errorf("failed to new resmgmtClient: %v", err)
	}
	c.resmgmtClient = resmgmtClient

	channelProvider := sess.ChannelProvider(channelID)
	channelClient, err := channel.New(channelProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to new channelClient: %v", err)
	}
	c.channelClient = channelClient
	ledgerClient, err := ledger.New(channelProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to new ledgerClient: %v", err)
	}
	c.ledgerClient = ledgerClient

	return c, nil
}

func (c *Client) Close() {
	if c.session != nil {
		c.session.Close()
	}
}

/*
---------------------------------------------
-------------channelClient-------------------
---------------------------------------------
*/

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

func (c *Client) ChannelInvokeHandler(
	handler invoke.Handler,
	request channel.Request,
	options ...channel.RequestOption,
) (channel.Response, error) {
	return c.channelClient.InvokeHandler(handler, request, options...)
}

/*
---------------------------------------------
-------------ledgerClient--------------------
---------------------------------------------
*/

func (c *Client) LedgerQueryConfig(options ...ledger.RequestOption) (fab.ChannelCfg, error) {
	return c.ledgerClient.QueryConfig(options...)
}

func (c *Client) LedgerQueryConfigBlock(options ...ledger.RequestOption) (*common.Block, error) {
	return c.ledgerClient.QueryConfigBlock(options...)
}

func (c *Client) LedgerQueryBlockchainInfo(options ...ledger.RequestOption) (*BlockchainInfo, error) {
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

func (c *Client) LedgerQueryBlock(blockNumber uint64, options ...ledger.RequestOption) (*Block, error) {
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

func (c *Client) LedgerQueryBlockByHash(blockHash string, options ...ledger.RequestOption) (*Block, error) {
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

func (c *Client) LedgerQueryBlockByTxID(txID string, options ...ledger.RequestOption) (*Block, error) {
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

func (c *Client) LedgerQueryTransaction(txID string, options ...ledger.RequestOption) (*ProcessedTransaction, error) {
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

/*
---------------------------------------------
-------------resmgmtClient-------------------
---------------------------------------------
*/

func (c *Client) ResmgmtQueryChannels(options ...resmgmt.RequestOption) (*ChannelQueryResponse, error) {
	respFrom, err := c.resmgmtClient.QueryChannels(options...)
	if err != nil {
		return nil, fmt.Errorf("failed to call resmgmt QueryChannels: %v", err)
	}
	respTo, err := DecodeChannelQueryResponse(respFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to decode channelQueryResponse: %v", err)
	}
	return respTo, nil
}

func (c *Client) ResmgmtJoinChannel(channelID string, options ...resmgmt.RequestOption) error {
	return c.resmgmtClient.JoinChannel(channelID, options...)
}

func (c *Client) ResmgmtSaveChannel(req resmgmt.SaveChannelRequest, options ...resmgmt.RequestOption) (*SaveChannelResponse, error) {
	respFrom, err := c.resmgmtClient.SaveChannel(req, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to call resmgmt SaveChannel: %v", err)
	}
	respTo, err := DecodeSaveChannelResponse(respFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to decode saveChannelResponse: %v", err)
	}
	return respTo, nil
}

func (c *Client) ResmgmtInstallCC(req resmgmt.InstallCCRequest, options ...resmgmt.RequestOption) ([]InstallCCResponse, error) {
	respFrom, err := c.resmgmtClient.InstallCC(req, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to call resmgmt InstallCC: %v", err)
	}
	respTo, err := DecodeInstallCCResponse(respFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to decode installCCResponse: %v", err)
	}
	return respTo, nil
}

func (c *Client) ResmgmtInstantiateCC(channelID string, req resmgmt.InstantiateCCRequest, options ...resmgmt.RequestOption) (*InstantiateCCResponse, error) {
	respFrom, err := c.resmgmtClient.InstantiateCC(channelID, req, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to call resmgmt InstantiateCC: %v", err)
	}
	respTo, err := DecodeInstantiateCCResponse(respFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to decode instantiateCCResponse: %v", err)
	}
	return respTo, nil
}

func (c *Client) ResmgmtUpgradeCC(channelID string, req resmgmt.UpgradeCCRequest, options ...resmgmt.RequestOption) (*UpgradeCCResponse, error) {
	respFrom, err := c.resmgmtClient.UpgradeCC(channelID, req, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to call resmgmt UpgradeCC: %v", err)
	}
	respTo, err := DecodeUpgradeCCResponse(respFrom)
	if err != nil {
		return nil, fmt.Errorf("failed to decode upgradeCCResponse: %v", err)
	}
	return respTo, nil
}
