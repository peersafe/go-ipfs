package commands

import (
	"io"
	"strings"

	logging "QmQg1J6vikuXF9oDvm4wpdeAUvvkVEKW1EYDw9HhTMnP2b/go-log"

	cmds "github.com/ipfs/go-ipfs/commands"
	unixfs "github.com/ipfs/go-ipfs/core/commands/unixfs"
)

var log = logging.Logger("core/commands")

type TestOutput struct {
	Foo string
	Bar int
}

const (
	ApiOption = "api"
)

var Root = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "global p2p merkle-dag filesystem",
		Synopsis: `
ipfs [<flags>] <command> [<arg>] ...
`,
		ShortDescription: `
BASIC COMMANDS

    init          Initialize ipfs local configuration
    add <path>    Add an object to ipfs
    cat <ref>     Show ipfs object data
    get <ref>     Download ipfs objects
    ls <ref>      List links from an object
    refs <ref>    List hashes of links from an object

DATA STRUCTURE COMMANDS

    block         Interact with raw blocks in the datastore
    object        Interact with raw dag nodes
    file          Interact with Unix filesystem objects

ADVANCED COMMANDS

    daemon        Start a long-running daemon process
    mount         Mount an ipfs read-only mountpoint
    resolve       Resolve any type of name
    name          Publish or resolve IPNS names
    dns           Resolve DNS links
    pin           Pin objects to local storage
    repo gc       Garbage collect unpinned objects

NETWORK COMMANDS

    id            Show info about ipfs peers
    bootstrap     Add or remove bootstrap peers
    swarm         Manage connections to the p2p network
    dht           Query the dht for values or peers
    ping          Measure the latency of a connection
    diag          Print diagnostics

TOOL COMMANDS

    config        Manage configuration
    version       Show ipfs version information
    update        Download and apply go-ipfs updates
    commands      List all available commands

Use 'ipfs <command> --help' to learn more about each command.

ipfs uses a repository in the local file system. By default, the repo is located
at ~/.ipfs. To change the repo location, set the $IPFS_PATH environment variable:

    export IPFS_PATH=/path/to/ipfsrepo
`,
	},
	Options: []cmds.Option{
		cmds.StringOption("config", "c", "Path to the configuration file to use"),
		cmds.BoolOption("debug", "D", "Operate in debug mode"),
		cmds.BoolOption("help", "Show the full command help text"),
		cmds.BoolOption("h", "Show a short version of the command help text"),
		cmds.BoolOption("local", "L", "Run the command locally, instead of using the daemon"),
		cmds.StringOption(ApiOption, "Use a specific API instance (defaults to /ip4/127.0.0.1/tcp/5001)"),
	},
}

// commandsDaemonCmd is the "ipfs commands" command for daemon
var CommandsDaemonCmd = CommandsCmd(Root)

var rootSubcommands = map[string]*cmds.Command{
	"add":       AddCmd,
	"block":     BlockCmd,
	"bootstrap": BootstrapCmd,
	"cat":       CatCmd,
	"commands":  CommandsDaemonCmd,
	"config":    ConfigCmd,
	"dht":       DhtCmd,
	"diag":      DiagCmd,
	"dns":       DNSCmd,
	"get":       GetCmd,
	"id":        IDCmd,
	"log":       LogCmd,
	"ls":        LsCmd,
	"mount":     MountCmd,
	"name":      NameCmd,
	"object":    ObjectCmd,
	"pin":       PinCmd,
	"remotepin": RemotePinCmd,
	"remotels":  RemoteLsCmd,
	"relaypin":  RelayPinCmd,
	"ping":      PingCmd,
	"refs":      RefsCmd,
	"repo":      RepoCmd,
	"resolve":   ResolveCmd,
	"stats":     StatsCmd,
	"swarm":     SwarmCmd,
	"tar":       TarCmd,
	"tour":      tourCmd,
	"file":      unixfs.UnixFSCmd,
	"update":    ExternalBinary(),
	"version":   VersionCmd,
	"bitswap":   BitswapCmd,
}

// RootRO is the readonly version of Root
var RootRO = &cmds.Command{}

var CommandsDaemonROCmd = CommandsCmd(RootRO)

var RefsROCmd = &cmds.Command{}

var rootROSubcommands = map[string]*cmds.Command{
	"block": &cmds.Command{
		Subcommands: map[string]*cmds.Command{
			"stat": blockStatCmd,
			"get":  blockGetCmd,
		},
	},
	"cat":      CatCmd,
	"commands": CommandsDaemonROCmd,
	"get":      GetCmd,
	"ls":       LsCmd,
	"name": &cmds.Command{
		Subcommands: map[string]*cmds.Command{
			"resolve": IpnsCmd,
		},
	},
	"object": &cmds.Command{
		Subcommands: map[string]*cmds.Command{
			"data":  objectDataCmd,
			"links": objectLinksCmd,
			"get":   objectGetCmd,
			"stat":  objectStatCmd,
			"patch": objectPatchCmd,
		},
	},
	"refs": RefsROCmd,
	//"resolve": ResolveCmd,
	"version": VersionCmd,
}

func init() {
	*RootRO = *Root

	// sanitize readonly refs command
	*RefsROCmd = *RefsCmd
	RefsROCmd.Subcommands = map[string]*cmds.Command{}

	Root.Subcommands = rootSubcommands
	RootRO.Subcommands = rootROSubcommands
}

type MessageOutput struct {
	Message string
}

func MessageTextMarshaler(res cmds.Response) (io.Reader, error) {
	return strings.NewReader(res.Output().(*MessageOutput).Message), nil
}
