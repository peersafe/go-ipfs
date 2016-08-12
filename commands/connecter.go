package commands

type Client interface {
	Send(req Request) (Response, error)
}
