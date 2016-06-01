package commands

import (
	"bytes"
	"encoding/json"
	"io"

	cmds "github.com/ipfs/go-ipfs/commands"
	core "github.com/ipfs/go-ipfs/core"
	merkledag "github.com/ipfs/go-ipfs/merkledag"
	path "github.com/ipfs/go-ipfs/path"
	unixfs "github.com/ipfs/go-ipfs/unixfs"
	unixfspb "github.com/ipfs/go-ipfs/unixfs/pb"
)

type LsLink struct {
	Name, Hash string
	Size       uint64
	Type       unixfspb.Data_DataType
}

type LsObject struct {
	Hash  string
	Links []LsLink
}

type LsOutput struct {
	Objects []LsObject
}

var LsCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "List links from an object.",
		ShortDescription: `
Retrieves the object named by <ipfs-or-ipns-path> and displays the links
it contains, with the following format:

  <link base58 hash> <link size in bytes> <link name>
`,
	},

	Arguments: []cmds.Argument{
		cmds.StringArg("ipfs-path", true, true, "The path to the IPFS object(s) to list links from").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.BoolOption("headers", "v", "Print table headers (Hash, Name, Size)"),
	},
	Run: func(req cmds.Request, res cmds.Response) {
		node, err := req.InvocContext().GetNode()
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		// get options early -> exit early in case of error
		if _, _, err := req.Option("headers").Bool(); err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		paths := req.Arguments()

		var dagnodes []*merkledag.Node
		for _, fpath := range paths {
			dagnode, err := core.Resolve(req.Context(), node, path.Path(fpath))
			if err != nil {
				res.SetError(err, cmds.ErrNormal)
				return
			}
			dagnodes = append(dagnodes, dagnode)
		}

		output := make([]LsObject, len(req.Arguments()))
		for i, dagnode := range dagnodes {
			output[i] = LsObject{
				Hash:  paths[i],
				Links: make([]LsLink, len(dagnode.Links)),
			}
			for j, link := range dagnode.Links {
				link.Node, err = link.GetNode(req.Context(), node.DAG)
				if err != nil {
					res.SetError(err, cmds.ErrNormal)
					return
				}
				d, err := unixfs.FromBytes(link.Node.Data)
				if err != nil {
					res.SetError(err, cmds.ErrNormal)
					return
				}
				output[i].Links[j] = LsLink{
					Name: link.Name,
					Hash: link.Hash.B58String(),
					Size: link.Size,
					Type: d.GetType(),
				}
			}
		}

		res.SetOutput(&LsOutput{output})
	},
	Marshalers: cmds.MarshalerMap{
		cmds.Text: func(res cmds.Response) (io.Reader, error) {

			output := res.Output().(*LsOutput)
			buf := new(bytes.Buffer)
			enc := json.NewEncoder(buf)
			enc.Encode(output)

			return buf, nil
		},
	},
	Type: LsOutput{},
}
