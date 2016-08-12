package asyncchan

import (
	"io/ioutil"
	"runtime/debug"

	logging "gx/ipfs/QmNQynaz7qfriSUJkiEZUrm2Wen1u3Kj9goZzWtrPyu7XR/go-log"
	"gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"

	cmds "github.com/ipfs/go-ipfs/commands"
)

var log = logging.Logger("commands/asyncchan")

// the internal handler for the API
type internalHandler struct {
	ctx  cmds.Context
	root *cmds.Command
}

type Handler interface {
	ServeAsyncChan(cmds.Request)
}

// The Handler struct is funny because we want to wrap our internal handler
// with CORS while keeping our fields.
type handler struct {
	internalHandler
}

func NewHandler(ctx cmds.Context, root *cmds.Command) Handler {
	// setup request logger
	ctx.ReqLog = new(cmds.ReqLog)

	// Wrap the internal handler with CORS handling-middleware.
	// Create a handler for the API.
	internal := internalHandler{
		ctx:  ctx,
		root: root,
	}
	return &handler{internal}
}

func (i handler) ServeAsyncChan(r cmds.Request) {
	// Call the CORS handler which wraps the internal handler.
	i.internalHandler.ServeAsyncChan(r)
}

func (i internalHandler) ServeAsyncChan(req cmds.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("a panic has occurred in the commands handler!")
			log.Error(r)

			debug.PrintStack()
		}
	}()

	// get the node's context to pass into the commands.
	node, err := i.ctx.GetNode()
	if err != nil {
		log.Errorf("cmds/http: couldn't GetNode():%s", err)
		req.CallFunc().Call("", err)
		return
	}

	ctx, cancel := context.WithCancel(node.Context())
	defer cancel()

	rlog := i.ctx.ReqLog.Add(req)
	defer rlog.Finish()

	//ps: take note of the name clash - commreqands.Context != context.Context
	req.SetInvocContext(i.ctx)

	err = req.SetRootContext(ctx)
	if err != nil {
		log.Errorf("setRootContext failed! %v", err)
		req.CallFunc().Call("", err)
		return
	}

	res := i.root.Call(req)
	if req.Command().PostRun != nil {
		req.Command().PostRun(req, res)
	}

	if err := res.Error(); err != nil {
		log.Errorf("root Call failed! %v", err)
		req.CallFunc().Call("", err)
		return
	}

	out, err := res.Reader()
	if err != nil {
		log.Errorf("res.Reader failed! %v", err)
		req.CallFunc().Call("", err)
		return
	}

	buf, err := ioutil.ReadAll(out)
	if err != nil {
		log.Errorf("iotuil.ReadAll failed! %v", err)
		req.CallFunc().Call("", err)
		return
	}

	req.CallFunc().Call(string(buf), nil)
}
