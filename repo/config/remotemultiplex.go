package config

// RemoteMultiplex contains options for the HTTP gateway server.
type RemoteMultiplex struct {
	Master  bool
	TryTime int
	Slave   [][]string
	MaxPin  int
}

var DefaultSlave = [][]string{
	{"/ip4/115.159.105.185/tcp/40001/ipfs/QmPkFbxAQ7DeKD5VGSh9HQrdS574pyNzDmxJeGrRJxoucF", "TkACALc9"},
	{"/ip4/119.29.67.136/tcp/40001/ipfs/QmTGkgHSsULk8p3AKTAqKixxidZQXFyF7mCURcutPqrwjQ", "E8C6UgG9"},
	{"/ip4/219.223.222.4/tcp/40001/ipfs/Qmf96ojxn2i8QPZ83FbutnwGjffEXsV4VaoFGzuC3YEwwY", "qpzKK98T"},
}

var DefaultMaxPin = 5

var DefaultTryTime = 10
