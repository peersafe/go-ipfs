package commands

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"strings"
	"time"

	u "gx/ipfs/Qmb912gdngC1UWwTkhuW8knyRbcWeu5kqkxBpveLmW8bSr/go-ipfs-util"

	pstore "gx/ipfs/QmeXj9VAjmYQZxpmVz7VzccbJrpmr8qkCDSjfVNsPTWTYU/go-libp2p-peerstore"
	peer "gx/ipfs/QmfMmLGoKzCHDN7cGgk64PJr4iipzidDRME8HABSJqvmhC/go-libp2p-peer"

	cmds "github.com/ipfs/go-ipfs/commands"
	core "github.com/ipfs/go-ipfs/core"

	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
)

const kRemoteMsgTimeout = 20 * time.Second

type RemoteMsgResult struct {
	Success bool
	Text    string
}

var RemoteMsgCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Send format message to IPFS node.",
		Synopsis: `
ipfs remotemsg <peer ID> <peer KEY> <msg>
		`,
		ShortDescription: `
Send format message to IPFS node.	
		`,
	},
	Arguments: []cmds.Argument{
		cmds.StringArg("peer ID", true, false, "ID of peer to be notify"),
		cmds.StringArg("peer KEY", true, false, "Password of peer to be notify"),
		cmds.StringArg("msg ", true, true, "Msg to send"),
	},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {
			outChan, ok := res.Output().(<-chan interface{})
			if !ok {
				fmt.Println(reflect.TypeOf(res.Output()))
				return nil, u.ErrCast()
			}

			marshal := func(v interface{}) (io.Reader, error) {
				obj, ok := v.(*RemoteMsgResult)
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

		msg := ""
		if req.InvocContext().GetAsyncChan != nil {
			msg = req.Arguments()[2]
		} else {
			msg += `{"`
			msgArgs := req.Arguments()[2:]
			for k, v := range msgArgs {
				if k == len(msgArgs)-1 {
					result := strings.Split(v, ":")
					msg += result[0] + `":"` + result[1] + `"}`
					break
				}
				result := strings.Split(v, ":")
				msg += result[0] + `":"` + result[1] + `","`
			}
		}

		fmt.Printf("msg=%v\n", msg)

		outChan := remoteMsg(ctx, n, peerID, key, msg)
		res.SetOutput(outChan)
	},
	Type: RemoteMsgResult{},
}

func remoteMsg(ctx context.Context, n *core.IpfsNode, pid peer.ID, key, msg string) <-chan interface{} {
	outChan := make(chan interface{})
	go func() {
		defer close(outChan)

		// 添加需要通讯的节点到底层网络中
		if len(n.Peerstore.Addrs(pid)) == 0 {
			ctx, cancel := context.WithTimeout(ctx, kRemoteMsgTimeout)
			defer cancel()
			p, err := n.Routing.FindPeer(ctx, pid)
			if err != nil {
				outChan <- &RemoteMsgResult{
					Success: false,
					Text:    fmt.Sprintf("Peer lookup error: %s", err),
				}
				return
			}
			n.Peerstore.AddAddrs(p.ID, p.Addrs, pstore.TempAddrTTL)
		}

		ctx, cancel := context.WithTimeout(ctx, kRemoteMsgTimeout)
		defer cancel()
		remotemsg, err := n.Remotemsg.RemoteMsg(ctx, pid, key, msg)
		if err != nil {
			outChan <- &RemoteMsgResult{Success: false, Text: fmt.Sprintf("RemoteMsg error: %s", err)}
			return
		}

		select {
		case <-ctx.Done():
			outChan <- &RemoteMsgResult{
				Success: false,
				Text:    fmt.Sprintf("Remote node error"),
			}
		case err, ok := <-remotemsg:
			if !ok {
				outChan <- &RemoteMsgResult{
					Success: false,
					Text:    fmt.Sprintf("Client error: remotemsg chan is close"),
				}
				break
			}
			if err != nil {
				fmt.Println("client receive error:", err)
				outChan <- &RemoteMsgResult{
					Success: false,
					Text:    fmt.Sprint(err),
				}
				break
			}
			outChan <- &RemoteMsgResult{
				Success: true,
			}
		}

	}()
	return outChan
}
