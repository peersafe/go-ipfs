package main

import (
	"fmt"
	"time"

	ipfs_mobile "github.com/ipfs/go-ipfs/cmd/ipfs_mobile"
)

var done chan struct{}

const (
	PATH = "ipfs_home"
)

func main() {
	done = make(chan struct{}, 1)

	// init
	callback := new(CallBack)
	fmt.Println(ipfs_mobile.IpfsInit(PATH))
	time.Sleep(1 * time.Second)

	// daemon
	go ipfs_mobile.IpfsAsyncDaemon(PATH, callback)

	// add
	add_uid := ipfs_mobile.IpfsAsyncAdd("apimain.go", 5)
	fmt.Println("func=[IpfsAsyncAdd],uid= ", add_uid)

	<-done
	close(done)
}
