package asyncchan

import (
	cmds "github.com/ipfs/go-ipfs/commands"

	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
)

type client struct {
}

func NewClient() cmds.Client {
	return &client{}
}

func (c *client) Send(req cmds.Request) (cmds.Response, error) {
	log.Debugf(">>>>>Send[%v]", req)

	if req.Context() == nil {
		log.Warningf("no context set in request")
		if err := req.SetRootContext(context.TODO()); err != nil {
			return nil, err
		}
	}

	send, _, _ := req.InvocContext().GetAsyncChan()

	*send <- req

	res := cmds.NewResponse(req)

	return res, nil
}
