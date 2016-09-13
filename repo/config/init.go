package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/rand"

	ci "gx/ipfs/QmVoi5es8D5fNHZDqoW6DgDAEPEV5hQp8GBz161vZXiwpQ/go-libp2p-crypto"
	peer "gx/ipfs/QmWtbQU15LaB5B1JC2F7TV9P4K88vD3PpA4AJrwfCjhML8/go-libp2p-peer"
)

func Init(out io.Writer, nBitsForKeypair int, ismobile bool) (*Config, error) {
	identity, err := identityConfig(out, nBitsForKeypair)
	if err != nil {
		return nil, err
	}
	if ismobile {
		identity.IsMobile = "true"
	}

	bootstrapPeers, err := DefaultBootstrapPeers()
	if err != nil {
		return nil, err
	}

	datastore, err := datastoreConfig()
	if err != nil {
		return nil, err
	}

	conf := &Config{

		// setup the node's default addresses.
		// NOTE: two swarm listen addrs, one tcp, one utp.
		Addresses: Addresses{
			Swarm: []string{
				"/ip4/0.0.0.0/tcp/4001",
				// "/ip4/0.0.0.0/udp/4002/utp", // disabled for now.
				"/ip6/::/tcp/4001",
			},
			API:     "/ip4/127.0.0.1/tcp/5001",
			Gateway: "/ip4/127.0.0.1/tcp/8080",
		},

		Datastore: datastore,
		Bootstrap: BootstrapPeerStrings(bootstrapPeers),
		Identity:  identity,
		Discovery: Discovery{MDNS{
			Enabled:    true,
			Interval:   10,
			ServiceTag: DefaultMDNSServiceTag,
		}},

		// setup the node mount points.
		Mounts: Mounts{
			IPFS: "/ipfs",
			IPNS: "/ipns",
		},

		Ipns: Ipns{
			ResolveCacheSize: 128,
		},

		Gateway: Gateway{
			RootRedirect: "",
			Writable:     false,
			PathPrefixes: []string{},
			HTTPHeaders: map[string][]string{
				"Access-Control-Allow-Origin":  []string{"*"},
				"Access-Control-Allow-Methods": []string{"GET"},
				"Access-Control-Allow-Headers": []string{"X-Requested-With"},
			},
		},
		RemoteMultiplex: RemoteMultiplex{
			Master:  false,
			TryTime: DefaultTryTime,
			Slave:   DefaultSlave,
			MaxPin:  DefaultMaxPin,
		},
		Reprovider: Reprovider{
			Interval: "12h",
		},
	}

	return conf, nil
}

func datastoreConfig() (Datastore, error) {
	dspath, err := DataStorePath("")
	if err != nil {
		return Datastore{}, err
	}
	return Datastore{
		Path:               dspath,
		Type:               "leveldb",
		StorageMax:         "10GB",
		StorageGCWatermark: 90, // 90%
		GCPeriod:           "1h",
		HashOnRead:         false,
		BloomFilterSize:    512 * 8 * 1024,
		ARCCacheSize:       64 * 1024,
	}, nil
}

// identityConfig initializes a new identity.
func identityConfig(out io.Writer, nbits int) (Identity, error) {
	// TODO guard higher up
	ident := Identity{}
	if nbits < 1024 {
		return ident, errors.New("Bitsize less than 1024 is considered unsafe.")
	}

	fmt.Fprintf(out, "generating %v-bit RSA keypair...", nbits)
	sk, pk, err := ci.GenerateKeyPair(ci.RSA, nbits)
	if err != nil {
		return ident, err
	}
	fmt.Fprintf(out, "done\n")

	// currently storing key unencrypted. in the future we need to encrypt it.
	// TODO(security)
	skbytes, err := sk.Bytes()
	if err != nil {
		return ident, err
	}
	ident.PrivKey = base64.StdEncoding.EncodeToString(skbytes)

	id, err := peer.IDFromPublicKey(pk)
	if err != nil {
		return ident, err
	}
	ident.PeerID = id.Pretty()
	fmt.Fprintf(out, "peer identity: %s\n", ident.PeerID)

	secretLen := 8
	privKeyLen := len(ident.PrivKey)
	bSecret := make([]byte, secretLen)
	for i := 0; i < secretLen; i++ {
		index := rand.Intn(privKeyLen)
		bSecret[i] = ident.PrivKey[index]
	}
	ident.Secret = string(bSecret)

	// default
	ident.IsMobile = "false"

	return ident, nil
}
