package core

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"syscall"
	"time"

	bstore "github.com/ipfs/go-ipfs/blocks/blockstore"
	key "github.com/ipfs/go-ipfs/blocks/key"
	bserv "github.com/ipfs/go-ipfs/blockservice"
	offline "github.com/ipfs/go-ipfs/exchange/offline"
	dag "github.com/ipfs/go-ipfs/merkledag"
	path "github.com/ipfs/go-ipfs/path"
	pin "github.com/ipfs/go-ipfs/pin"
	repo "github.com/ipfs/go-ipfs/repo"
	cfg "github.com/ipfs/go-ipfs/repo/config"

	ds "gx/ipfs/QmNgqJarToRiq2GBaPJhkmW4B5BxS5B74E1rkGvv2JoaTp/go-datastore"
	dsync "gx/ipfs/QmNgqJarToRiq2GBaPJhkmW4B5BxS5B74E1rkGvv2JoaTp/go-datastore/sync"
	goprocessctx "gx/ipfs/QmQopLATEYMNg7dVqZRNDfeE2S1yKy8zrRh5xnYiuqeZBn/goprocess/context"
	pstore "gx/ipfs/QmSZi9ygLohBUGyHMqE5N6eToPwqcg7bZQTULeVLFu7Q6d/go-libp2p-peerstore"
	retry "gx/ipfs/QmV71dCdjYivZwta22a7YvNmeS86r1EpnBjq9XAD21rBV9/retry-datastore"
	ci "gx/ipfs/QmVoi5es8D5fNHZDqoW6DgDAEPEV5hQp8GBz161vZXiwpQ/go-libp2p-crypto"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
)

type BuildCfg struct {
	// If online is set, the node will have networking enabled
	Online bool

	// If permament then node should run more expensive processes
	// that will improve performance in long run
	Permament bool

	// If NilRepo is set, a repo backed by a nil datastore will be constructed
	NilRepo bool

	Routing RoutingOption
	Host    HostOption
	Repo    repo.Repo
}

func (cfg *BuildCfg) fillDefaults() error {
	if cfg.Repo != nil && cfg.NilRepo {
		return errors.New("cannot set a repo and specify nilrepo at the same time")
	}

	if cfg.Repo == nil {
		var d ds.Datastore
		d = ds.NewMapDatastore()
		if cfg.NilRepo {
			d = ds.NewNullDatastore()
		}
		r, err := defaultRepo(dsync.MutexWrap(d))
		if err != nil {
			return err
		}
		cfg.Repo = r
	}

	if cfg.Routing == nil {
		cfg.Routing = DHTOption
	}

	if cfg.Host == nil {
		cfg.Host = DefaultHostOption
	}

	return nil
}

func defaultRepo(dstore repo.Datastore) (repo.Repo, error) {
	c := cfg.Config{}
	priv, pub, err := ci.GenerateKeyPairWithReader(ci.RSA, 1024, rand.Reader)
	if err != nil {
		return nil, err
	}

	data, err := pub.Hash()
	if err != nil {
		return nil, err
	}

	privkeyb, err := priv.Bytes()
	if err != nil {
		return nil, err
	}

	c.Bootstrap = cfg.DefaultBootstrapAddresses
	c.Addresses.Swarm = []string{"/ip4/0.0.0.0/tcp/4001"}
	c.Identity.PeerID = key.Key(data).B58String()
	c.Identity.PrivKey = base64.StdEncoding.EncodeToString(privkeyb)

	return &repo.Mock{
		D: dstore,
		C: c,
	}, nil
}

func NewNode(ctx context.Context, cfg *BuildCfg) (*IpfsNode, error) {
	if cfg == nil {
		cfg = new(BuildCfg)
	}

	err := cfg.fillDefaults()
	if err != nil {
		return nil, err
	}

	n := &IpfsNode{
		mode:      offlineMode,
		Repo:      cfg.Repo,
		ctx:       ctx,
		Peerstore: pstore.NewPeerstore(),
		Closer:    make(chan struct{}),
	}
	if cfg.Online {
		n.mode = onlineMode
	}

	// TODO: this is a weird circular-ish dependency, rework it
	n.proc = goprocessctx.WithContextAndTeardown(ctx, n.teardown)
	if err := setupNode(ctx, n, cfg); err != nil {
		n.Close()
		return nil, err
	}

	return n, nil
}

func isTooManyFDError(err error) bool {
	perr, ok := err.(*os.PathError)
	if ok && perr.Err == syscall.EMFILE {
		return true
	}

	return false
}

func setupNode(ctx context.Context, n *IpfsNode, cfg *BuildCfg) error {
	// setup local peer ID (private key is loaded in online setup)
	if err := n.loadID(); err != nil {
		return err
	}

	rds := &retry.Datastore{
		Batching:    n.Repo.Datastore(),
		Delay:       time.Millisecond * 200,
		Retries:     6,
		TempErrFunc: isTooManyFDError,
	}

	var err error
	bs := bstore.NewBlockstore(rds)
	opts := bstore.DefaultCacheOpts()
	conf, err := n.Repo.Config()
	if err != nil {
		return err
	}

	opts.HasBloomFilterSize = conf.Datastore.BloomFilterSize
	if !cfg.Permament {
		opts.HasBloomFilterSize = 0
	}

	opts.HasARCCacheSize = conf.Datastore.ARCCacheSize
	if !cfg.Permament {
		opts.HasARCCacheSize = 0
	}

	blockstore, removeCids, err := bstore.CachedBlockstore(bs, ctx, opts)
	if err != nil {
		return err
	}
	n.Blockstore = blockstore

	rcfg, err := n.Repo.Config()
	if err != nil {
		return err
	}

	if rcfg.Datastore.HashOnRead {
		bs.HashOnRead(true)
	}

	if cfg.Online {
		do := setupDiscoveryOption(rcfg.Discovery)
		if err := n.startOnlineServices(ctx, cfg.Routing, cfg.Host, do); err != nil {
			return err
		}
	} else {
		n.Exchange = offline.Exchange(n.Blockstore)
	}

	n.Blocks = bserv.New(n.Blockstore, n.Exchange)
	n.DAG = dag.NewDAGService(n.Blocks)

	internalDag := dag.NewDAGService(bserv.New(n.Blockstore, offline.Exchange(n.Blockstore)))
	n.Pinning, err = pin.LoadPinner(n.Repo.Datastore(), n.DAG, internalDag)
	if err != nil {
		// TODO: we should move towards only running 'NewPinner' explicity on
		// node init instead of implicitly here as a result of the pinner keys
		// not being found in the datastore.
		// this is kinda sketchy and could cause data loss
		n.Pinning = pin.NewPinner(n.Repo.Datastore(), n.DAG, internalDag)
	}
	n.Resolver = &path.Resolver{DAG: n.DAG}

	err = n.loadFilesRoot()
	if err != nil {
		return err
	}

	if cfg.Online && removeCids != nil {
		go func() {
			for {
				select {
				case removeCid := <-removeCids:
					removeKey := key.B58KeyDecode(removeCid.String())
					fmt.Println("arcCache remoteKey=", removeKey)
					exist, err := n.Blockstore.Has(removeKey)
					if err != nil {
						log.Errorf("has block err:%v\n", err)
					}

					_, pinState, err := n.Pinning.IsPinned(removeCid)
					if err != nil {
						log.Errorf("ping removeKey err:%v\n", err)
					}
					if !pinState && exist {
						err := n.Blockstore.DeleteBlock(removeKey)
						if err != nil {
							log.Errorf("delete block err:%v\n", err)
						}
						fmt.Println("delete block:", removeKey)
					}
					fmt.Println("==============removeKey pin state=", pinState, "exist state=", exist)
				case <-n.Closer:
					close(removeCids)
					return
				}
			}
		}()
	}

	return nil
}
