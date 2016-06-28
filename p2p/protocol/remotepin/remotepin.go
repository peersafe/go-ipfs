package remotepin

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"time"

	context "github.com/ipfs/go-ipfs/Godeps/_workspace/src/golang.org/x/net/context"

	localhandle "github.com/ipfs/go-ipfs/core/remotehandle"
	host "github.com/ipfs/go-ipfs/p2p/host"
	inet "github.com/ipfs/go-ipfs/p2p/net"
	peer "github.com/ipfs/go-ipfs/p2p/peer"
	path "github.com/ipfs/go-ipfs/path"
	logging "github.com/ipfs/go-ipfs/vendor/QmQg1J6vikuXF9oDvm4wpdeAUvvkVEKW1EYDw9HhTMnP2b/go-log"
)

var log = logging.Logger("remotepin")

const RemotepinSize = 46 + 6 // /ipfs/....

const ID = "/ipfs/remotepin"

type RemotepinService struct {
	Host        host.Host
	LocalHandle localhandle.Remotepin
}

func NewRemotepinService(h host.Host, handler localhandle.Remotepin) *RemotepinService {
	ps := &RemotepinService{h, handler}
	h.SetStreamHandler(ID, ps.RemotepinHandler)
	return ps
}

func (p *RemotepinService) RemotepinHandler(s inet.Stream) {
	buf := make([]byte, RemotepinSize)

	for {
		_, err := io.ReadFull(s, buf)
		if err != nil {
			log.Debug(err)
			return
		}

		err = p.LocalHandle.RemotePin(string(buf))
		if err != nil {
			buffer := bytes.NewBuffer(buf)
			buffer.WriteString("0")
			buf = buffer.Bytes()
		}

		_, err = s.Write(buf)
		if err != nil {
			log.Debug(err)
			return
		}
	}
}

func (ps *RemotepinService) Remotepin(ctx context.Context, p peer.ID, path path.Path) (<-chan time.Duration, error) {
	s, err := ps.Host.NewStream(ID, p)
	if err != nil {
		return nil, err
	}

	out := make(chan time.Duration)
	go func() {
		defer close(out)
		select {
		case <-ctx.Done():
			return
		default:
			t, err := remotepin(s, path)
			if err != nil {
				log.Debugf("remotepin error: %s", err)
				return
			}

			select {
			case out <- t:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

func remotepin(s inet.Stream, path path.Path) (time.Duration, error) {
	before := time.Now()
	if !strings.HasPrefix(string(path), "/ipfs/") {
		path = "/ipfs/" + path
	}
	_, err := s.Write([]byte(path))
	if err != nil {
		return 0, err
	}

	rbuf := make([]byte, RemotepinSize)
	_, err = io.ReadFull(s, rbuf)
	if err != nil {
		return 0, err
	}

	if !bytes.Equal([]byte(path), rbuf) {
		return 0, errors.New("remotepin packet was incorrect!")
	}

	return time.Now().Sub(before), nil
}
