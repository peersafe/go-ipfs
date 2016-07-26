package config

// RemoteMultiplex contains options for the HTTP gateway server.
type RemoteMultiplex struct {
	Master  bool
	TryTime int
	Slave   [][]string
	MaxPin  int
}

var DefaultSlave = [][]string{
	{"/ip4/101.201.40.124/tcp/4001/ipfs/QmZDYAhmMDtnoC6XZRw8R1swgoshxKvXDA9oQF97AYkPZc", "SuG2pVkw"},
	{"/ip4/219.223.222.4/tcp/4001/ipfs/QmS8DGkGkkcZMjb6MWUL9TSeSa3E4Jffi5FBVKeoYogYKv", "ieAJK5ar"},
}

var DefaultMaxPin = 5

var DefaultTryTime = 10
