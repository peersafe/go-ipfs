package commands

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"time"

	pstore "gx/ipfs/QmeXj9VAjmYQZxpmVz7VzccbJrpmr8qkCDSjfVNsPTWTYU/go-libp2p-peerstore"
	peer "gx/ipfs/QmfMmLGoKzCHDN7cGgk64PJr4iipzidDRME8HABSJqvmhC/go-libp2p-peer"

	u "gx/ipfs/Qmb912gdngC1UWwTkhuW8knyRbcWeu5kqkxBpveLmW8bSr/go-ipfs-util"

	cmds "github.com/ipfs/go-ipfs/commands"
	core "github.com/ipfs/go-ipfs/core"
	path "github.com/ipfs/go-ipfs/path"

	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
)

const kRemoteLsTimeout = 10 * time.Second

type RemoteLsResult struct {
	Success bool
	Text    string
}

var RemoteLsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "notify given IPFS node,Ls specific objects to it's storage",
		Synopsis: `
ipfs remotels <peer ID> <ipfs Path>
		`,
		ShortDescription: `
notify given <peer ID> IPFS node, Ls specific <ipfs path> object to it's local storage	
		`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("peer ID", true, false, "ID of peer to be notify"),
		cmds.StringArg("peer KEY", true, false, "Password of peer to be notify"),
		cmds.StringArg("ipfs Path", true, false, "Path to object(s) to be ls"),
	},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			outChan, ok := res.Output().(<-chan interface{})
			if !ok {
				fmt.Println(reflect.TypeOf(res.Output()))
				return nil, u.ErrCast()
			}

			marshal := func(v interface{}) (io.Reader, error) {
				obj, ok := v.(*RemoteLsResult)
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

		outChan := remoteLsPeer(ctx, n, peerID, key, path)
		res.SetOutput(outChan)
	},
	Type: RemoteLsResult{},
}

func remoteLsPeer(ctx context.Context, n *core.IpfsNode, pid peer.ID, key string, path path.Path) <-chan interface{} {
	outChan := make(chan interface{})
	go func() {
		defer close(outChan)

		// 添加需要通讯的节点到底层网络中
		if len(n.Peerstore.Addrs(pid)) == 0 {
			ctx, cancel := context.WithTimeout(ctx, kRemoteLsTimeout)
			defer cancel()
			p, err := n.Routing.FindPeer(ctx, pid)
			if err != nil {
				outChan <- &RemoteLsResult{
					Success: false,
					Text:    fmt.Sprintf("Peer lookup error: %s", err),
				}
				return
			}
			n.Peerstore.AddAddrs(p.ID, p.Addrs, pstore.TempAddrTTL)
		}

		ctx, cancel := context.WithTimeout(ctx, kRemoteLsTimeout)
		defer cancel()
		remotels, err := n.Remotels.Remotels(ctx, pid, key, path)
		if err != nil {
			outChan <- &RemoteLsResult{Success: false, Text: fmt.Sprintf("RemoteLs error: %s", err)}
			return
		}

		select {
		case <-ctx.Done():
			outChan <- &RemoteLsResult{
				Success: false,
				Text:    fmt.Sprintf("Remote node error"),
			}
		case err, ok := <-remotels:
			if !ok {
				outChan <- &RemoteLsResult{
					Success: false,
					Text:    fmt.Sprintf("Client error: remotels chan is close"),
				}
				break
			}
			if err != nil {
				outChan <- &RemoteLsResult{
					Success: false,
					Text:    fmt.Sprint(err),
				}
				break
			}
			outChan <- &RemoteLsResult{
				Success: true,
			}
		}

	}()
	return outChan
}
