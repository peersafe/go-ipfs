package remotehandle

type Remotepin interface {
	RemotePin(string) error
}
