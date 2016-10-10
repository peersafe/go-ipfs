package main

import (
	"fmt"

	ipfs_mobile "github.com/ipfs/go-ipfs/cmd/ipfs_mobile"
)

var done chan struct{}

const (
	// PATH = "ipfs_home"
	PATH = "/Users/sunzhiming/.ipfs"
)

func main() {
	// done = make(chan struct{}, 1)
	// defer close(done)

	// // init
	// callback := new(CallBack)
	// fmt.Println(ipfs_mobile.IpfsInit(PATH))
	// time.Sleep(1 * time.Second)

	// // daemon
	// go ipfs_mobile.IpfsAsyncDaemon(PATH, callback)

	// <-done
	// config test
	fmt.Println(ipfs_mobile.IpfsInit(PATH))
	ret, e := ipfs_mobile.IpfsConfig("Identity.IsMobile", "true")
	if e != nil {
		fmt.Println("func=[IpfsAsyncConfig],err= ", e)
		return
	}
	fmt.Println("11111func=[IpfsAsyncConfig],ret= ", ret)

	ret, e = ipfs_mobile.IpfsConfig("Identity.IsMobile", "")
	if e != nil {
		fmt.Println("func=[IpfsAsyncConfig],err= ", e)
		return
	}
	fmt.Println("22222func=[IpfsAsyncConfig],ret= ", ret)
}
