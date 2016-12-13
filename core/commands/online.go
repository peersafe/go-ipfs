package commands

import cmds "github.com/ipfs/go-ipfs/commands"

type Online struct {
	Result string
}

var OnlineCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Online daemon gracefully.",
		ShortDescription: "Online daemon gracefully.",
		LongDescription:  "Online daemon gracefully.",
	},

	Options:     []cmds.Option{},
	Subcommands: map[string]*cmds.Command{},
	Run:         onlineFunc,
	Type:        Online{},
}

func onlineFunc(req cmds.Request, res cmds.Response) {
	result := "Online Daemon!"

	online := Online{
		Result: result,
	}
	res.SetOutput(&online)
	node, _ := req.InvocContext().GetNode()
	node.Restart()
}
