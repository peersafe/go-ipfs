package commands

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"time"

	pstore "gx/ipfs/QmSZi9ygLohBUGyHMqE5N6eToPwqcg7bZQTULeVLFu7Q6d/go-libp2p-peerstore"
	peer "gx/ipfs/QmWXjJo15p4pzT7cayEwZi2sWgJqLnGDof6ZGMh9xBgU1p/go-libp2p-peer"

	u "gx/ipfs/QmZNVWh8LLjAavuQ2JXuFmuYH3C11xo988vSgp7UQrTRj1/go-ipfs-util"

	cmds "github.com/ipfs/go-ipfs/commands"
	core "github.com/ipfs/go-ipfs/core"
	path "github.com/ipfs/go-ipfs/path"

	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
)

const kRemotePinTimeout = 10 * time.Second

type RemotePinResult struct {
	Success bool
	Text    string
}

var RemotePinCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "notify given IPFS node,Pin specific objects to it's storage",
		Synopsis: `
ipfs remotepin <peer ID> <ipfs Path>
		`,
		ShortDescription: `
notify given <peer ID> IPFS node, Pin specific <ipfs path> object to it's local storage	
		`,
	},
	Arguments: []cmds.Argument{
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
				obj, ok := v.(*RemotePinResult)
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

		addr, peerID, err := ParsePeerParam(req.Arguments()[0])
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		if addr != nil {
			n.Peerstore.AddAddr(peerID, addr, pstore.TempAddrTTL) // temporary
		}

		key := req.Arguments()[1]
		matchstr := "^[a-zA-Z0-9-`=\\\\\\[\\];'\",./~!@#$%^&*()_+|{}:<>?]{8}$"
		if matched, err := regexp.MatchString(matchstr, key); err != nil || !matched {
			err = fmt.Errorf("peer key format error")
			res.SetError(err, cmds.ErrNormal)
			return
		}

		path := path.Path(req.Arguments()[2])
		if err := path.IsValid(); err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		outChan := remotePinPeer(ctx, n, peerID, key, path)
		res.SetOutput(outChan)
	},
	Type: RemotePinResult{},
}

func remotePinPeer(ctx context.Context, n *core.IpfsNode, pid peer.ID, key string, path path.Path) <-chan interface{} {
	outChan := make(chan interface{})
	go func() {
		defer close(outChan)

		// 添加需要通讯的节点到底层网络中
		if len(n.Peerstore.Addrs(pid)) == 0 {
			ctx, cancel := context.WithTimeout(ctx, kRemotePinTimeout)
			defer cancel()
			p, err := n.Routing.FindPeer(ctx, pid)
			if err != nil {
				outChan <- &RemotePinResult{
					Success: false,
					Text:    fmt.Sprintf("Peer lookup error: %s", err),
				}
				return
			}
			n.Peerstore.AddAddrs(p.ID, p.Addrs, pstore.TempAddrTTL)
		}

		ctx, cancel := context.WithTimeout(ctx, kRemotePinTimeout)
		defer cancel()
		// TAG remotePin.remotePin
		remotepin, err := n.Remotepin.Remotepin(ctx, pid, key, path)
		if err != nil {
			outChan <- &RemotePinResult{Success: false, Text: fmt.Sprintf("RemotePin error: %s", err)}
			return
		}

		select {
		case <-ctx.Done():
			outChan <- &RemotePinResult{
				Success: false,
				Text:    fmt.Sprintf("Remote node error"),
			}
		case err, ok := <-remotepin:
			if !ok {
				outChan <- &RemotePinResult{
					Success: false,
					Text:    fmt.Sprintf("Client error: remotepin chan is close"),
				}
				break
			}
			if err != nil {
				outChan <- &RemotePinResult{
					Success: false,
					Text:    fmt.Sprint(err),
				}
				break
			}
			outChan <- &RemotePinResult{
				Success: true,
			}
		}

	}()
	return outChan
}
