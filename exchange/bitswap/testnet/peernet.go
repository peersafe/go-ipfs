package bitswap

import (
	bsnet "github.com/ipfs/go-ipfs/exchange/bitswap/network"
	mockrouting "github.com/ipfs/go-ipfs/routing/mock"
	testutil "github.com/ipfs/go-ipfs/thirdparty/testutil"
	ds "gx/ipfs/QmNgqJarToRiq2GBaPJhkmW4B5BxS5B74E1rkGvv2JoaTp/go-datastore"
	peer "gx/ipfs/QmWtbQU15LaB5B1JC2F7TV9P4K88vD3PpA4AJrwfCjhML8/go-libp2p-peer"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
	mockpeernet "gx/ipfs/Qmf4ETeAWXuThBfWwonVyFqGFSgTWepUDEr1txcctvpTXS/go-libp2p/p2p/net/mock"
)

type peernet struct {
	mockpeernet.Mocknet
	routingserver mockrouting.Server
}

func StreamNet(ctx context.Context, net mockpeernet.Mocknet, rs mockrouting.Server) (Network, error) {
	return &peernet{net, rs}, nil
}

func (pn *peernet) Adapter(p testutil.Identity) bsnet.BitSwapNetwork {
	client, err := pn.Mocknet.AddPeer(p.PrivateKey(), p.Address())
	if err != nil {
		panic(err.Error())
	}
	routing := pn.routingserver.ClientWithDatastore(context.TODO(), p, ds.NewMapDatastore())
	return bsnet.NewFromIpfsHost(client, routing)
}

func (pn *peernet) HasPeer(p peer.ID) bool {
	for _, member := range pn.Mocknet.Peers() {
		if p == member {
			return true
		}
	}
	return false
}

var _ Network = &peernet{}
