package remotehandle

type Remotepin interface {
	RemotePin(string) error
}

type Remotels interface {
	RemoteLs(string) error
}
