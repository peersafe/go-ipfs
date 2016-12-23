/*
Package core implements the IpfsNode object and related methods.

Packages underneath core/ provide a (relatively) stable, low-level API
to carry out most IPFS-related tasks.  For more details on the other
interfaces and how core/... fits into the bigger IPFS picture, see:

  $ godoc github.com/ipfs/go-ipfs
*/
package core

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	bstore "github.com/ipfs/go-ipfs/blocks/blockstore"
	bserv "github.com/ipfs/go-ipfs/blockservice"
	diag "github.com/ipfs/go-ipfs/diagnostics"
	exchange "github.com/ipfs/go-ipfs/exchange"
	bitswap "github.com/ipfs/go-ipfs/exchange/bitswap"
	bsnet "github.com/ipfs/go-ipfs/exchange/bitswap/network"
	rp "github.com/ipfs/go-ipfs/exchange/reprovide"
	mount "github.com/ipfs/go-ipfs/fuse/mount"
	merkledag "github.com/ipfs/go-ipfs/merkledag"
	mfs "github.com/ipfs/go-ipfs/mfs"
	namesys "github.com/ipfs/go-ipfs/namesys"
	ipnsrp "github.com/ipfs/go-ipfs/namesys/republisher"
	path "github.com/ipfs/go-ipfs/path"
	pin "github.com/ipfs/go-ipfs/pin"
	repo "github.com/ipfs/go-ipfs/repo"
	config "github.com/ipfs/go-ipfs/repo/config"
	nilrouting "github.com/ipfs/go-ipfs/routing/none"
	offroute "github.com/ipfs/go-ipfs/routing/offline"
	ft "github.com/ipfs/go-ipfs/unixfs"

	dht "gx/ipfs/QmNQPjpcXrwwwgDErKzKUm2xxhXCB3cuFgTHsrcCJ5uGbu/go-libp2p-kad-dht"
	p2phost "gx/ipfs/QmPTGbC34bPKaUm9wTxBo7zSCac7pDuG42ZmnXC718CKZZ/go-libp2p-host"
	ds "gx/ipfs/QmRWDav6mzWseLWeYfVd5fvUKiVe9xNH29YfMF438fG364/go-datastore"
	goprocess "gx/ipfs/QmSF8fPo3jgVBAy8fpdjjYqgG87dkJgUprRBHRd2tmfgpP/goprocess"
	mamask "gx/ipfs/QmSMZwvs3n4GBikZ7hKzT17c3bk65FmyZo2JqtJ16swqCv/multiaddr-filter"
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	b58 "gx/ipfs/QmT8rehPR3F6bmwL6zjUN8XpiDBFFpMP2myPdC6ApsWfJf/go-base58"
	mssmux "gx/ipfs/QmTfjLsou9ic6L4KqCcmbLSZcdiFu8q1v6njKp121pbbXx/go-smux-multistream"
	ma "gx/ipfs/QmUAQaWbKxGCUTuoQVvvicbQNZ9APF5pDGWyAZSe93AtKH/go-multiaddr"
	floodsub "gx/ipfs/QmV5jot2GfVXmgvetHExJCa2hprebf3AKjprZtuwaXSr1v/floodsub"
	addrutil "gx/ipfs/QmVDnc2zvyQm8LhT72n22THcshvH7j3qPMnhvjerQER62T/go-addr-util"
	spdy "gx/ipfs/QmWUNsat6Jb19nC5CiJCDXepTkxjdxi3eZqeoB6mrmmaGu/go-smux-spdystream"
	swarm "gx/ipfs/QmWfxnAiQ5TnnCgiX9ikVUKFNHRgGhbgKdx5DoKPELD7P4/go-libp2p-swarm"
	mplex "gx/ipfs/QmXGevGDVTqeKdisBzaxEK4CJZqfxeXiVSWLaXaVWcG5on/go-smux-multiplex"
	metrics "gx/ipfs/QmY2otvyPM2sTaDsczo7Yuosg98sUMCJ9qx1gpPaAPTS9B/go-libp2p-metrics"
	u "gx/ipfs/Qmb912gdngC1UWwTkhuW8knyRbcWeu5kqkxBpveLmW8bSr/go-ipfs-util"
	routing "gx/ipfs/QmbkGVaN9W6RYJK4Ws5FvMKXKDqdRQ5snhtaa92qP6L8eU/go-libp2p-routing"
	yamux "gx/ipfs/Qmbn7RYyWzBVXiUp9jZ1dA4VADHy9DtS7iZLwfhEUQvm3U/go-smux-yamux"
	discovery "gx/ipfs/QmbzCT1CwxVZ2ednptC9RavuJe7Bv8DDi2Ne89qUrA37XM/go-libp2p/p2p/discovery"
	p2pbhost "gx/ipfs/QmbzCT1CwxVZ2ednptC9RavuJe7Bv8DDi2Ne89qUrA37XM/go-libp2p/p2p/host/basic"
	rhost "gx/ipfs/QmbzCT1CwxVZ2ednptC9RavuJe7Bv8DDi2Ne89qUrA37XM/go-libp2p/p2p/host/routed"
	ping "gx/ipfs/QmbzCT1CwxVZ2ednptC9RavuJe7Bv8DDi2Ne89qUrA37XM/go-libp2p/p2p/protocol/ping"
	cid "gx/ipfs/QmcTcsTvfaeEBRFo1TkFgT8sRmgi1n1LTZpecfVP8fzpGD/go-cid"
	pstore "gx/ipfs/QmeXj9VAjmYQZxpmVz7VzccbJrpmr8qkCDSjfVNsPTWTYU/go-libp2p-peerstore"
	smux "gx/ipfs/QmeZBgYBHvxMukGK5ojg28BCNLB9SeXqT7XXg6o7r2GbJy/go-stream-muxer"
	peer "gx/ipfs/QmfMmLGoKzCHDN7cGgk64PJr4iipzidDRME8HABSJqvmhC/go-libp2p-peer"
	ic "gx/ipfs/QmfWDLQjGjVe4fr5CoztYW2DYYjRysMJrFe1RCsXLPTf46/go-libp2p-crypto"
	relaypin "github.com/ipfs/go-ipfs/remotecmd/relaypin"
	remotels "github.com/ipfs/go-ipfs/remotecmd/remotels"
	remotemsg "github.com/ipfs/go-ipfs/remotecmd/remotemsg"
	remotepin "github.com/ipfs/go-ipfs/remotecmd/remotepin"

)

const IpnsValidatorTag = "ipns"
const kSizeBlockstoreWriteCache = 100
const kReprovideFrequency = time.Hour * 12
const discoveryConnTimeout = time.Second * 30

var log = logging.Logger("core")

type mode int

const (
	// zero value is not a valid mode, must be explicitly set
	invalidMode mode = iota
	localMode
	offlineMode
	onlineMode
)

// IpfsNode is IPFS Core module. It represents an IPFS instance.
type IpfsNode struct {

	// Self
	Identity peer.ID // the local node's identity

	Repo repo.Repo

	// Local node
	Pinning    pin.Pinner // the pinning manager
	Mounts     Mounts     // current mount state, if any.
	PrivateKey ic.PrivKey // the local node's private Key

	// Services
	Peerstore  pstore.Peerstore     // storage for other Peer instances
	Blockstore bstore.GCBlockstore  // the block store (lower level)
	Blocks     bserv.BlockService   // the block service, get/add blocks.
	DAG        merkledag.DAGService // the merkle dag service, get/add objects.
	Resolver   *path.Resolver       // the path resolution system
	Reporter   metrics.Reporter
	Discovery  discovery.Service
	FilesRoot  *mfs.Root

	// Online
	PeerHost     p2phost.Host        // the network host (server+client)
	Bootstrapper io.Closer           // the periodic bootstrapper
	Routing      routing.IpfsRouting // the routing system. recommend ipfs-dht
	Exchange     exchange.Interface  // the block exchange + strategy (bitswap)
	Namesys      namesys.NameSystem  // the name system, resolves paths to hashes
	Diagnostics  *diag.Diagnostics   // the diagnostics service
	Ping         *ping.PingService
	Reprovider   *rp.Reprovider // the value reprovider system
	IpnsRepub    *ipnsrp.Republisher
	Relaypin     *relaypin.RelaypinService
	Remotepin    *remotepin.RemotepinService
	Remotels     *remotels.RemotelsService
	Remotemsg    *remotemsg.RemotemsgService

	Floodsub *floodsub.PubSub

	proc goprocess.Process
	ctx  context.Context

	mode         mode
	localModeSet bool
	Closer chan  struct{}
}

// Mounts defines what the node's mount state is. This should
// perhaps be moved to the daemon or mount. It's here because
// it needs to be accessible across daemon requests.
type Mounts struct {
	Ipfs mount.Mount
	Ipns mount.Mount
}

func (n *IpfsNode) startOnlineServices(ctx context.Context, routingOption RoutingOption, hostOption HostOption, do DiscoveryOption, pubsub, mplex bool) error {

	if n.PeerHost != nil { // already online.
		return errors.New("node already online")
	}

	// load private key
	if err := n.LoadPrivateKey(); err != nil {
		return err
	}

	// get undialable addrs from config
	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}
	var addrfilter []*net.IPNet
	for _, s := range cfg.Swarm.AddrFilters {
		f, err := mamask.NewMask(s)
		if err != nil {
			return fmt.Errorf("incorrectly formatted address filter in config: %s", s)
		}
		addrfilter = append(addrfilter, f)
	}

	if !cfg.Swarm.DisableBandwidthMetrics {
		// Set reporter
		n.Reporter = metrics.NewBandwidthCounter()
	}

	tpt := makeSmuxTransport(mplex)

	peerhost, err := hostOption(ctx, n.Identity, n.Peerstore, n.Reporter, addrfilter, tpt, cfg)
	if err != nil {
		return err
	}

	if err := n.startOnlineServicesWithHost(ctx, peerhost, routingOption); err != nil {
		return err
	}

	// Ok, now we're ready to listen.
	if err := startListening(ctx, n.PeerHost, cfg); err != nil {
		return err
	}

	n.Reprovider = rp.NewReprovider(n.Routing, n.Blockstore)

	if cfg.Reprovider.Interval != "0" {
		interval := kReprovideFrequency
		if cfg.Reprovider.Interval != "" {
			dur, err := time.ParseDuration(cfg.Reprovider.Interval)
			if err != nil {
				return err
			}

			interval = dur
		}

		go n.Reprovider.ProvideEvery(ctx, interval)
	}

	if pubsub {
		n.Floodsub = floodsub.NewFloodSub(ctx, peerhost)
	}

	// setup local discovery
	if do != nil {
		service, err := do(ctx, n.PeerHost)
		if err != nil {
			log.Error("mdns error: ", err)
		} else {
			service.RegisterNotifee(n)
			n.Discovery = service
		}
	}

	return n.Bootstrap(DefaultBootstrapConfig)
}

func makeSmuxTransport(mplexExp bool) smux.Transport {
	mstpt := mssmux.NewBlankTransport()

	ymxtpt := &yamux.Transport{
		AcceptBacklog:          8192,
		ConnectionWriteTimeout: time.Second * 10,
		KeepAliveInterval:      time.Second * 30,
		EnableKeepAlive:        true,
		MaxStreamWindowSize:    uint32(1024 * 512),
		LogOutput:              ioutil.Discard,
	}

	mstpt.AddTransport("/yamux/1.0.0", ymxtpt)

	mstpt.AddTransport("/spdy/3.1.0", spdy.Transport)

	if mplexExp {
		mstpt.AddTransport("/mplex/6.7.0", mplex.DefaultTransport)
	}

	// Allow muxer preference order overriding
	if prefs := os.Getenv("LIBP2P_MUX_PREFS"); prefs != "" {
		mstpt.OrderPreference = strings.Fields(prefs)
	}

	return mstpt
}

func setupDiscoveryOption(d config.Discovery) DiscoveryOption {
	if d.MDNS.Enabled {
		return func(ctx context.Context, h p2phost.Host) (discovery.Service, error) {
			if d.MDNS.Interval == 0 {
				d.MDNS.Interval = 5
			}
			// TODO(liliuhai): Remember to change mdns.go support set ServiceTag,modify by mdns.diff file.
			return discovery.NewMdnsService(ctx, h, time.Duration(d.MDNS.Interval)*time.Second, d.MDNS.ServiceTag)
		}
	} 
	return nil
}

func (n *IpfsNode) HandlePeerFound(p pstore.PeerInfo) {
	log.Warning("trying peer info: ", p)
	ctx, cancel := context.WithTimeout(n.Context(), discoveryConnTimeout)
	defer cancel()
	if err := n.PeerHost.Connect(ctx, p); err != nil {
		log.Warning("Failed to connect to peer found by discovery: ", err)
	}
}

// startOnlineServicesWithHost  is the set of services which need to be
// initialized with the host and _before_ we start listening.
func (n *IpfsNode) startOnlineServicesWithHost(ctx context.Context, host p2phost.Host, routingOption RoutingOption) error {
	// setup diagnostics service
	n.Diagnostics = diag.NewDiagnostics(n.Identity, host)
	n.Ping = ping.NewPingService(host)

	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}
	n.Remotepin = remotepin.NewRemotepinService(host, cfg.Identity.Secret, cfg.RemoteMultiplex)
	n.Remotels = remotels.NewRemotelsService(host, cfg.Identity.Secret)
	n.Relaypin = relaypin.NewRelaypinService(host, cfg.Identity.Secret)
	n.Remotemsg = remotemsg.NewRemotemsgService(host, cfg.Identity.Secret, cfg.RemoteMultiplex)

	// setup routing service
	r, err := routingOption(ctx, host, n.Repo.Datastore())
	if err != nil {
		return err
	}
	n.Routing = r

	// Wrap standard peer host with routing system to allow unknown peer lookups
	n.PeerHost = rhost.Wrap(host, n.Routing)

	localIsMobile := false
	if cfg.Identity.IsMobile == "true" {
		localIsMobile = true
	}
	// setup exchange service
	const alwaysSendToPeer = true // use YesManStrategy
	bitswapNetwork := bsnet.NewFromIpfsHost(n.PeerHost, n.Routing)
	n.Exchange = bitswap.New(ctx, n.Identity, bitswapNetwork, n.Blockstore, alwaysSendToPeer, localIsMobile)

	size, err := n.getCacheSize()
	if err != nil {
		return err
	}

	// setup name system
	n.Namesys = namesys.NewNameSystem(n.Routing, n.Repo.Datastore(), size)

	// setup ipns republishing
	err = n.setupIpnsRepublisher()
	if err != nil {
		return err
	}

	return nil
}

func (n *IpfsNode) Restart() {
	rcfg, err := n.Repo.Config()
	if err != nil {
		//return err
	}

	do := setupDiscoveryOption(rcfg.Discovery)
	// LLH pubsub mplex set false !!!!!TODO!!!!!!
	n.Online(n.Context(), DHTOption, DefaultHostOption, do, false, false)
	n.Blocks = bserv.New(n.Blockstore, n.Exchange)
	n.DAG.SetBlockService(n.Blocks)
}

func (n *IpfsNode) Offline() error {
	log.Debug("core is shutting down...")
	// owned objects are closed in this teardown to ensure that they're closed
	// regardless of which constructor was used to add them to the node.
	var closers []io.Closer

	if n.PeerHost == nil {
		return nil
	}

	// NOTE: The order that objects are added(closed) matters, if an object
	// needs to use another during its shutdown/cleanup process, it should be
	// closed before that other object

	//n.Closer<-new struct{}

	if n.Blocks != nil {
		closers = append(closers, n.Blocks)
	}

	if n.Bootstrapper != nil {
		closers = append(closers, n.Bootstrapper)
	}

	if n.Discovery != nil {
		closers = append(closers, n.Discovery)
	}

	if n.Exchange != nil {
		closers = append(closers, n.Exchange)
	}

	if dht, ok := n.Routing.(*dht.IpfsDHT); ok {
		closers = append(closers, dht.Process())
	}

	if n.PeerHost != nil {
		closers = append(closers, n.PeerHost)
	}

	//IpnsRepub  *ipnsrp.Republisher   LLH ?????

	var errs []error
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	n.PeerHost = nil
	return nil
}

func (n *IpfsNode) Online(ctx context.Context, routingOption RoutingOption, hostOption HostOption, do DiscoveryOption, pubsub, mplex bool) error {
	// get undialable addrs from config
	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}
	var addrfilter []*net.IPNet
	for _, s := range cfg.Swarm.AddrFilters {
		f, err := mamask.NewMask(s)
		if err != nil {
			return fmt.Errorf("incorrectly formatted address filter in config: %s", s)
		}
		addrfilter = append(addrfilter, f)
	}

	tpt := makeSmuxTransport(mplex)

	peerhost, err := hostOption(ctx, n.Identity, n.Peerstore, n.Reporter, addrfilter, tpt, cfg)
	if err != nil {
		return err
	}

	if err := n.startOnlineServicesWithHost(ctx, peerhost, routingOption); err != nil {
		return err
	}

	// Ok, now we're ready to listen.
	if err := startListening(ctx, n.PeerHost, cfg); err != nil {
		return err
	}

	n.Reprovider = rp.NewReprovider(n.Routing, n.Blockstore)

	if cfg.Reprovider.Interval != "0" {
		interval := kReprovideFrequency
		if cfg.Reprovider.Interval != "" {
			dur, err := time.ParseDuration(cfg.Reprovider.Interval)
			if err != nil {
				return err
			}

			interval = dur
		}

		go n.Reprovider.ProvideEvery(ctx, interval)
	}

	// setup local discovery
	if do != nil {
		service, err := do(ctx, n.PeerHost)
		if err != nil {
			log.Error("mdns error: ", err)
		} else {
			service.RegisterNotifee(n)
			n.Discovery = service
		}
	}

	return n.Bootstrap(DefaultBootstrapConfig)
}

// getCacheSize returns cache life and cache size
func (n *IpfsNode) getCacheSize() (int, error) {
	cfg, err := n.Repo.Config()
	if err != nil {
		return 0, err
	}

	cs := cfg.Ipns.ResolveCacheSize
	if cs == 0 {
		cs = 128
	}
	if cs < 0 {
		return 0, fmt.Errorf("cannot specify negative resolve cache size")
	}
	return cs, nil
}

func (n *IpfsNode) setupIpnsRepublisher() error {
	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}

	n.IpnsRepub = ipnsrp.NewRepublisher(n.Routing, n.Repo.Datastore(), n.Peerstore)
	n.IpnsRepub.AddName(n.Identity)

	if cfg.Ipns.RepublishPeriod != "" {
		d, err := time.ParseDuration(cfg.Ipns.RepublishPeriod)
		if err != nil {
			return fmt.Errorf("failure to parse config setting IPNS.RepublishPeriod: %s", err)
		}

		if !u.Debug && (d < time.Minute || d > (time.Hour*24)) {
			return fmt.Errorf("config setting IPNS.RepublishPeriod is not between 1min and 1day: %s", d)
		}

		n.IpnsRepub.Interval = d
	}

	if cfg.Ipns.RecordLifetime != "" {
		d, err := time.ParseDuration(cfg.Ipns.RepublishPeriod)
		if err != nil {
			return fmt.Errorf("failure to parse config setting IPNS.RecordLifetime: %s", err)
		}

		n.IpnsRepub.RecordLifetime = d
	}

	n.Process().Go(n.IpnsRepub.Run)

	return nil
}

// Process returns the Process object
func (n *IpfsNode) Process() goprocess.Process {
	return n.proc
}

// Close calls Close() on the Process object
func (n *IpfsNode) Close() error {
	return n.proc.Close()
}

// Context returns the IpfsNode context
func (n *IpfsNode) Context() context.Context {
	if n.ctx == nil {
		n.ctx = context.TODO()
	}
	return n.ctx
}

// teardown closes owned children. If any errors occur, this function returns
// the first error.
func (n *IpfsNode) teardown() error {
	log.Debug("core is shutting down...")
	// owned objects are closed in this teardown to ensure that they're closed
	// regardless of which constructor was used to add them to the node.
	var closers []io.Closer

	// NOTE: The order that objects are added(closed) matters, if an object
	// needs to use another during its shutdown/cleanup process, it should be
	// closed before that other object

	if n.FilesRoot != nil {
		closers = append(closers, n.FilesRoot)
	}

	if n.Exchange != nil {
		closers = append(closers, n.Exchange)
	}

	if n.Mounts.Ipfs != nil && !n.Mounts.Ipfs.IsActive() {
		closers = append(closers, mount.Closer(n.Mounts.Ipfs))
	}
	if n.Mounts.Ipns != nil && !n.Mounts.Ipns.IsActive() {
		closers = append(closers, mount.Closer(n.Mounts.Ipns))
	}

	if dht, ok := n.Routing.(*dht.IpfsDHT); ok {
		closers = append(closers, dht.Process())
	}

	if n.Blocks != nil {
		closers = append(closers, n.Blocks)
	}

	if n.Bootstrapper != nil {
		closers = append(closers, n.Bootstrapper)
	}

	if n.PeerHost != nil {
		closers = append(closers, n.PeerHost)
	}

	// Repo closed last, most things need to preserve state here
	closers = append(closers, n.Repo)

	var errs []error
	for _, closer := range closers {
		if err := closer.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func (n *IpfsNode) OnlineMode() bool {
	switch n.mode {
	case onlineMode:
		return true
	default:
		return false
	}
}

func (n *IpfsNode) SetLocal(isLocal bool) {
	if isLocal {
		n.mode = localMode
	}
	n.localModeSet = true
}

func (n *IpfsNode) LocalMode() bool {
	if !n.localModeSet {
		// programmer error should not happen
		panic("local mode not set")
	}
	switch n.mode {
	case localMode:
		return true
	default:
		return false
	}
}

func (n *IpfsNode) Bootstrap(cfg BootstrapConfig) error {

	// TODO what should return value be when in offlineMode?
	if n.Routing == nil {
		return nil
	}

	if n.Bootstrapper != nil {
		n.Bootstrapper.Close() // stop previous bootstrap process.
	}

	// if the caller did not specify a bootstrap peer function, get the
	// freshest bootstrap peers from config. this responds to live changes.
	if cfg.BootstrapPeers == nil {
		cfg.BootstrapPeers = func() []pstore.PeerInfo {
			ps, err := n.loadBootstrapPeers()
			if err != nil {
				log.Warning("failed to parse bootstrap peers from config")
				return nil
			}
			return ps
		}
	}

	var err error
	n.Bootstrapper, err = Bootstrap(n, cfg)
	return err
}

func (n *IpfsNode) loadID() error {
	if n.Identity != "" {
		return errors.New("identity already loaded")
	}

	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}

	cid := cfg.Identity.PeerID
	if cid == "" {
		return errors.New("identity was not set in config (was 'ipfs init' run?)")
	}
	if len(cid) == 0 {
		return errors.New("no peer ID in config! (was 'ipfs init' run?)")
	}

	n.Identity = peer.ID(b58.Decode(cid))
	return nil
}

func (n *IpfsNode) GetKey(name string) (ic.PrivKey, error) {
	if name == "self" {
		return n.PrivateKey, nil
	} else {
		return n.Repo.Keystore().Get(name)
	}
}

func (n *IpfsNode) LoadPrivateKey() error {
	if n.Identity == "" || n.Peerstore == nil {
		return errors.New("loaded private key out of order.")
	}

	if n.PrivateKey != nil {
		return errors.New("private key already loaded")
	}

	cfg, err := n.Repo.Config()
	if err != nil {
		return err
	}

	sk, err := loadPrivateKey(&cfg.Identity, n.Identity)
	if err != nil {
		return err
	}

	n.PrivateKey = sk
	n.Peerstore.AddPrivKey(n.Identity, n.PrivateKey)
	n.Peerstore.AddPubKey(n.Identity, sk.GetPublic())
	return nil
}

func (n *IpfsNode) loadBootstrapPeers() ([]pstore.PeerInfo, error) {
	cfg, err := n.Repo.Config()
	if err != nil {
		return nil, err
	}

	parsed, err := cfg.BootstrapPeers()
	if err != nil {
		return nil, err
	}
	return toPeerInfos(parsed), nil
}

func (n *IpfsNode) loadFilesRoot() error {
	dsk := ds.NewKey("/local/filesroot")
	pf := func(ctx context.Context, c *cid.Cid) error {
		return n.Repo.Datastore().Put(dsk, c.Bytes())
	}

	var nd *merkledag.ProtoNode
	val, err := n.Repo.Datastore().Get(dsk)

	switch {
	case err == ds.ErrNotFound || val == nil:
		nd = ft.EmptyDirNode()
		_, err := n.DAG.Add(nd)
		if err != nil {
			return fmt.Errorf("failure writing to dagstore: %s", err)
		}
	case err == nil:
		c, err := cid.Cast(val.([]byte))
		if err != nil {
			return err
		}

		rnd, err := n.DAG.Get(n.Context(), c)
		if err != nil {
			return fmt.Errorf("error loading filesroot from DAG: %s", err)
		}

		pbnd, ok := rnd.(*merkledag.ProtoNode)
		if !ok {
			return merkledag.ErrNotProtobuf
		}

		nd = pbnd
	default:
		return err
	}

	mr, err := mfs.NewRoot(n.Context(), n.DAG, nd, pf)
	if err != nil {
		return err
	}

	n.FilesRoot = mr
	return nil
}

// SetupOfflineRouting loads the local nodes private key and
// uses it to instantiate a routing system in offline mode.
// This is primarily used for offline ipns modifications.
func (n *IpfsNode) SetupOfflineRouting() error {
	if n.Routing != nil {
		// Routing was already set up
		return nil
	}
	err := n.LoadPrivateKey()
	if err != nil {
		return err
	}

	n.Routing = offroute.NewOfflineRouter(n.Repo.Datastore(), n.PrivateKey)

	size, err := n.getCacheSize()
	if err != nil {
		return err
	}

	n.Namesys = namesys.NewNameSystem(n.Routing, n.Repo.Datastore(), size)

	return nil
}

func loadPrivateKey(cfg *config.Identity, id peer.ID) (ic.PrivKey, error) {
	sk, err := cfg.DecodePrivateKey("passphrase todo!")
	if err != nil {
		return nil, err
	}

	id2, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return nil, err
	}

	if id2 != id {
		return nil, fmt.Errorf("private key in config does not match id: %s != %s", id, id2)
	}

	return sk, nil
}

func listenAddresses(cfg *config.Config) ([]ma.Multiaddr, error) {
	var listen []ma.Multiaddr
	for _, addr := range cfg.Addresses.Swarm {
		maddr, err := ma.NewMultiaddr(addr)
		if err != nil {
			return nil, fmt.Errorf("Failure to parse config.Addresses.Swarm: %s", cfg.Addresses.Swarm)
		}
		listen = append(listen, maddr)
	}

	return listen, nil
}

type HostOption func(ctx context.Context, id peer.ID, ps pstore.Peerstore, bwr metrics.Reporter, fs []*net.IPNet, tpt smux.Transport, cfg *config.Config) (p2phost.Host, error)

var DefaultHostOption HostOption = constructPeerHost

// isolates the complex initialization steps
func constructPeerHost(ctx context.Context, id peer.ID, ps pstore.Peerstore, bwr metrics.Reporter, fs []*net.IPNet, tpt smux.Transport, cfg *config.Config) (p2phost.Host, error) {
	ismobile := false
	if cfg.Identity.IsMobile == "true" {
		ismobile = true
	}
	// no addresses to begin with. we'll start later.
	swrm, err := swarm.NewSwarmWithProtector(ctx, nil, id, ps, nil, tpt, bwr, ismobile)
	if err != nil {
		return nil, err
	}

	network := (*swarm.Network)(swrm)

	for _, f := range fs {
		network.Swarm().Filters.AddDialFilter(f)
	}

	host := p2pbhost.New(network, p2pbhost.NATPortMap, bwr)

	return host, nil
}

// startListening on the network addresses
func startListening(ctx context.Context, host p2phost.Host, cfg *config.Config) error {
	listenAddrs, err := listenAddresses(cfg)
	if err != nil {
		return err
	}

	// make sure we error out if our config does not have addresses we can use
	log.Debugf("Config.Addresses.Swarm:%s", listenAddrs)
	filteredAddrs := addrutil.FilterUsableAddrs(listenAddrs)
	log.Debugf("Config.Addresses.Swarm:%s (filtered)", filteredAddrs)
	if len(filteredAddrs) < 1 {
		return fmt.Errorf("addresses in config not usable: %s", listenAddrs)
	}

	// Actually start listening:
	if err := host.Network().Listen(filteredAddrs...); err != nil {
		return err
	}

	// list out our addresses
	addrs, err := host.Network().InterfaceListenAddresses()
	if err != nil {
		return err
	}
	log.Infof("Swarm listening at: %s", addrs)
	return nil
}

func constructDHTRouting(ctx context.Context, host p2phost.Host, dstore repo.Datastore) (routing.IpfsRouting, error) {
	dhtRouting := dht.NewDHT(ctx, host, dstore)
	dhtRouting.Validator[IpnsValidatorTag] = namesys.IpnsRecordValidator
	dhtRouting.Selector[IpnsValidatorTag] = namesys.IpnsSelectorFunc
	return dhtRouting, nil
}

func constructClientDHTRouting(ctx context.Context, host p2phost.Host, dstore repo.Datastore) (routing.IpfsRouting, error) {
	dhtRouting := dht.NewDHTClient(ctx, host, dstore)
	dhtRouting.Validator[IpnsValidatorTag] = namesys.IpnsRecordValidator
	dhtRouting.Selector[IpnsValidatorTag] = namesys.IpnsSelectorFunc
	return dhtRouting, nil
}

type RoutingOption func(context.Context, p2phost.Host, repo.Datastore) (routing.IpfsRouting, error)

type DiscoveryOption func(context.Context, p2phost.Host) (discovery.Service, error)

var DHTOption RoutingOption = constructDHTRouting
var DHTClientOption RoutingOption = constructClientDHTRouting
var NilRouterOption RoutingOption = nilrouting.ConstructNilRouting
