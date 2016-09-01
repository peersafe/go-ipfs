package commands

import cmds "github.com/ipfs/go-ipfs/commands"

type ShutDown struct {
	Result string
}

var ShutdownCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Shutdown daemon gracefully.",
		ShortDescription: "Shutdown daemon gracefully.",
		LongDescription:  "Shutdown daemon gracefully.",
	},

	Options:     []cmds.Option{},
	Subcommands: map[string]*cmds.Command{},
	Run:         shutdownFunc,
	Type:        ShutDown{},
}

func shutdownFunc(req cmds.Request, res cmds.Response) {
	result := "Shutdown Daemon!"

	shutdown := ShutDown{
		Result: result,
	}
	res.SetOutput(&shutdown)
	node, _ := req.InvocContext().GetNode()

	if req.InvocContext().GetAsyncChan != nil {
		send, _, _ := req.InvocContext().GetAsyncChan()
		if send != nil {
			close(*send)
		}
	}

	node.Close()
}
