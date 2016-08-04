package apiinterface

var GApiInterface Apier

type Apier interface {
	Cmd(string, int) (int, string, error)
}
