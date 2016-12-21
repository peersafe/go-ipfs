package commands

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"time"

	pstore "gx/ipfs/QmSZi9ygLohBUGyHMqE5N6eToPwqcg7bZQTULeVLFu7Q6d/go-libp2p-peerstore"
	peer "gx/ipfs/QmfMmLGoKzCHDN7cGgk64PJr4iipzidDRME8HABSJqvmhC/go-libp2p-peer"

	u "gx/ipfs/QmZNVWh8LLjAavuQ2JXuFmuYH3C11xo988vSgp7UQrTRj1/go-ipfs-util"

	cmds "github.com/ipfs/go-ipfs/commands"
	core "github.com/ipfs/go-ipfs/core"
	path "github.com/ipfs/go-ipfs/path"

	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
)

const kRelayPinTimeout = 10 * time.Second

type RelayPinResult struct {
	Success bool
	Text    string
}

var RelayPinCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "notify given IPFS node through common delay node,Pin specific objects to it's storage",
		Synopsis: `
ipfs relaypin <relay ID> <relay Key> <peer ID> <peer Key> <ipfs Path>
		`,
		ShortDescription: `
notify given <peer ID> IPFS node, Pin specific <ipfs path> object to it's local storage	, delay by <relay ID> <relay Key>
		`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("relay ID", true, false, "ID of relay peer to be notify"),
		cmds.StringArg("relay KEY", true, false, "Password of relay peer to be notify"),
		cmds.StringArg("peer ID", true, false, "ID of peer to be notify"),
		cmds.StringArg("peer KEY", true, false, "Password of peer to be notify"),
		cmds.StringArg("ipfs Path", true, false, "Path to object(s) to be pinned"),
	},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			outChan, ok := res.Output().(<-chan interface{})
			if !ok {
				fmt.Println(reflect.TypeOf(res.Output()))
				return nil, u.ErrCast()
			}

			marshal := func(v interface{}) (io.Reader, error) {
				obj, ok := v.(*RelayPinResult)
				if !ok {
					return nil, u.ErrCast()
				}

				buf := new(bytes.Buffer)
				if obj.Success {
					fmt.Fprintf(buf, "OK\n")
					return buf, nil
				} else {
					fmt.Fprintf(buf, "FAIL")
					return buf, fmt.Errorf("%s", obj.Text)
				}
			}

			return &cmds.ChannelMarshaler{
				Channel:   outChan,
				Marshaler: marshal,
				Res:       res,
			}, nil
		},
	},
	Run: func(req cmds.Request, res cmds.Response) {
		ctx := req.Context()
		n, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		// Must be online!
		if !n.OnlineMode() {
			res.SetError(errNotOnline, cmds.ErrClient)
			return
		}

		relayid := req.Arguments()[0]
		relaykey := req.Arguments()[1]
		peerid := req.Arguments()[2]
		peerkey := req.Arguments()[3]
		fpath := req.Arguments()[4]
		log.Debugf(">>>>>>>[%v][%v][%v][%v][%v]", relayid, relaykey, peerid, peerkey, fpath)

		addr, relay, err := ParsePeerParam(relayid)
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		_, _, err = ParsePeerParam(peerid)
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		if addr != nil {
			n.Peerstore.AddAddr(relay, addr, pstore.TempAddrTTL) // temporary
		}

		matchstr := "^[a-zA-Z0-9-`=\\\\\\[\\];'\",./~!@#$%^&*()_+|{}:<>?]{8}$"
		if matched, err := regexp.MatchString(matchstr, relaykey); err != nil || !matched {
			err = fmt.Errorf("relay key format error")
			res.SetError(err, cmds.ErrNormal)
			return
		}

		if matched, err := regexp.MatchString(matchstr, peerkey); err != nil || !matched {
			err = fmt.Errorf("peer key format error")
			res.SetError(err, cmds.ErrNormal)
			return
		}

		path := path.Path(fpath)
		if err := path.IsValid(); err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		log.Debug(">>>>>>start relayPinPeer")
		outChan := relayPinPeer(ctx, n, relay, relaykey, peerid, peerkey, path)
		res.SetOutput(outChan)
	},
	Type: RelayPinResult{},
}

func relayPinPeer(ctx context.Context, n *core.IpfsNode, relay peer.ID, relaykey, peerid, peerkey string, path path.Path) <-chan interface{} {
	outChan := make(chan interface{})
	go func() {
		defer close(outChan)

		// 添加需要通讯的节点到底层网络中
		if len(n.Peerstore.Addrs(relay)) == 0 {
			ctx, cancel := context.WithTimeout(ctx, kRelayPinTimeout)
			defer cancel()
			p, err := n.Routing.FindPeer(ctx, relay)
			if err != nil {
				outChan <- &RelayPinResult{
					Success: false,
					Text:    fmt.Sprintf("Peer lookup error: %s", err),
				}
				return
			}
			n.Peerstore.AddAddrs(p.ID, p.Addrs, pstore.TempAddrTTL)
		}

		ctx, cancel := context.WithTimeout(ctx, kRelayPinTimeout)
		defer cancel()
		log.Debugf(">>>>>>Relaypin.Relaypin call[%v][%v][%v][%v][%v]", relay, relaykey, peerid, peerkey, path)
		relaypin, err := n.Relaypin.Relaypin(ctx, relay, relaykey, peerid, peerkey, path, true)
		log.Debug(">>>>>>>>Relaypin return err=", err)
		if err != nil {
			outChan <- &RelayPinResult{Success: false, Text: fmt.Sprintf("RelayPin error: %s", err)}
			return
		}

		select {
		case <-ctx.Done():
			outChan <- &RelayPinResult{
				Success: false,
				Text:    fmt.Sprintf("Relay node error"),
			}
		case err, ok := <-relaypin:
			log.Debug(">>>>>>>>Relaypin relaypin return err=", err)
			if !ok {
				outChan <- &RelayPinResult{
					Success: false,
					Text:    fmt.Sprintf("Client error: relaypin chan is close"),
				}
				break
			}
			if err != nil {
				outChan <- &RelayPinResult{
					Success: false,
					Text:    fmt.Sprint(err),
				}
				log.Debug(">>>>>>>>>>>>>> Over")
				break
			}
			outChan <- &RelayPinResult{
				Success: true,
			}
		}

	}()
	return outChan
}
