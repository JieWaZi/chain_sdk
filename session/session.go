package session

import (
	"bytes"
	"io"
	"math/rand"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	contextApi "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/pkg/errors"
)

var logModules = [...]string{"chainsdk/resource", "chainsdk/binary"}

// Session provides a central location to create service clients from
type Session struct {
	Options Options
	fabsdk  *fabsdk.FabricSDK
}

// Options provides the means to control how a service context is created
type Options struct {
	Organization    string
	Username        string
	OrdererEndpoint string
	TargetEndPoint  []string
}

// FromReader creates a new instance of the Session from a io.Reader
func FromReader(in io.Reader, options ...ClientOption) (*Session, error) {
	if in == nil {
		return nil, errors.New("The in reader can't be nil")
	}
	return initFromReader(in, "yaml", options...)
}

// FromRaw creates a new instance of the Session from a byte array
func FromRaw(configBytes []byte, options ...ClientOption) (*Session, error) {
	if configBytes == nil {
		return nil, errors.New("The config bytes can't be nil")
	}
	buf := bytes.NewBuffer(configBytes)
	return initFromReader(buf, "yaml", options...)
}

func initFromReader(in io.Reader, configType string, options ...ClientOption) (*Session, error) {
	cfg := config.FromReader(in, configType)
	fabsdk, err := fabsdk.New(cfg)
	if err != nil {
		return nil, err
	}

	backend, err := fabsdk.Config()
	if err != nil {
		return nil, err
	}
	setLogLevel(backend)

	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	return &Session{Options: opts, fabsdk: fabsdk}, nil
}

// New creates a new instance of the Session
func New(cfgFile string, options ...ClientOption) (*Session, error) {
	if cfgFile == "" {
		return nil, errors.New("Config file must be provided")
	}

	cfg := config.FromFile(cfgFile)
	fabsdk, err := fabsdk.New(cfg)
	if err != nil {
		return nil, err
	}

	backend, err := fabsdk.Config()
	if err != nil {
		return nil, err
	}
	setLogLevel(backend)

	opts := Options{}
	for _, option := range options {
		option(&opts)
	}

	return &Session{Options: opts, fabsdk: fabsdk}, nil
}

// Context creates and returns context client which has all the necessary providers
func (sess *Session) Context() contextApi.ClientProvider {
	return sess.fabsdk.Context(fabsdk.WithUser(sess.Options.Username), fabsdk.WithOrg(sess.Options.Organization))
}

// ChannelContext creates and returns channel context
func (sess *Session) ChannelContext(channel string) contextApi.ChannelProvider {
	return sess.fabsdk.ChannelContext(channel, fabsdk.WithUser(sess.Options.Username), fabsdk.WithOrg(sess.Options.Organization))
}

// DiscoverLocalPeers queries the local peers for the given MSP context and returns all of the peers.
func (sess *Session) DiscoverLocalPeers() ([]fab.Peer, error) {
	orgUserContext := sess.Context()
	ctx, err := contextImpl.NewLocal(orgUserContext)
	if err != nil {
		return nil, errors.Wrap(err, "error creating local context")
	}

	mspid := ctx.Identifier().MSPID
	peers, err := ctx.LocalDiscoveryService().GetPeers()
	if err != nil {
		return nil, errors.Wrapf(err, "error getting peers for MSP [%s]", mspid)
	}
	return peers, nil
}

// SelectRandomPeer select random peer
func (sess *Session) SelectRandomPeer() (fab.Peer, error) {
	peers, err := sess.DiscoverLocalPeers()
	if err != nil {
		return nil, err
	}
	if len(peers) == 0 {
		return nil, errors.New("couldn't find any peer")
	}

	randomNumber := rand.Intn(len(peers))
	return peers[randomNumber], nil
}

// Close frees up caches and connections being maintained by the SDK
func (sess *Session) Close() {
	if sess.fabsdk != nil {
		sess.fabsdk.Close()
	}
}

func setLogLevel(backend core.ConfigBackend) {
	loggingLevelString, _ := backend.Lookup("client.logging.level")
	logLevel := logging.INFO
	if loggingLevelString != nil {
		var err error
		logLevel, err = logging.LogLevel(loggingLevelString.(string))
		if err != nil {
			panic(err)
		}
	}

	for _, logModule := range logModules {
		logging.SetLevel(logModule, logLevel)
	}
}
