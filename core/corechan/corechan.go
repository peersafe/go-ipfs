package corechan

import (
	"time"

	"gx/ipfs/QmQopLATEYMNg7dVqZRNDfeE2S1yKy8zrRh5xnYiuqeZBn/goprocess"
	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"

	cmds "github.com/ipfs/go-ipfs/commands"
	asyncchan "github.com/ipfs/go-ipfs/commands/asyncchan"
	core "github.com/ipfs/go-ipfs/core"
)

var log = logging.Logger("core/asyncserver")

func Serve(node *core.IpfsNode, asyncChan *(<-chan cmds.Request), handler asyncchan.Handler) error {
	// if the server exits beforehand
	var serverError error
	serverExited := make(chan struct{})

	node.Process().Go(func(p goprocess.Process) {
		//serverError = http.Serve(lis, handler)
		log.Debugf(">>>>>>>>>>>>>>>>> chan select start")
		for {
			req, ok := <-*asyncChan

			if !ok {
				break
			}

			go func() {
				handler.ServeAsyncChan(req)
			}()
		}
		log.Debugf(">>>>>>>>>>>>>>>>> chan select stop")

		close(serverExited)
	})

	// wait for server to exit.
	select {
	case <-serverExited:

	// if node being closed before server exits, close server
	case <-node.Process().Closing():
		log.Infof("server at async terminating...")

	outer:
		for {
			// wait until server exits
			select {
			case <-serverExited:
				// if the server exited as we are closing, we really dont care about errors
				serverError = nil
				break outer
			case <-time.After(5 * time.Second):
				log.Infof("waiting for server at async to terminate...")
			}
		}
	}

	log.Infof("server at async terminated")
	return serverError
}
