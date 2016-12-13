package commands

import cmds "github.com/ipfs/go-ipfs/commands"

type Offline struct {
	Result string
}

var OfflineCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Offline daemon gracefully.",
		ShortDescription: "Offline daemon gracefully.",
		LongDescription:  "Offline daemon gracefully.",
	},

	Options:     []cmds.Option{},
	Subcommands: map[string]*cmds.Command{},
	Run:         offlineFunc,
	Type:        Offline{},
}

func offlineFunc(req cmds.Request, res cmds.Response) {
	result := "Offline Daemon!"

	offline := Offline{
		Result: result,
	}
	res.SetOutput(&offline)
	node, _ := req.InvocContext().GetNode()
	node.Offline()
}
