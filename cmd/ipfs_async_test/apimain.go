package main

import (
	"fmt"
	"sync"
	"time"

	ipfs_lib "github.com/ipfs/go-ipfs/cmd/ipfs_lib"
)

type MyCall struct {
	call func(string, error)
}

func (c *MyCall) Call(result string, err error) {
	c.call(result, err)
}

func ipfsAsyncPeerid() {
	fmt.Println("======================================== ipfsAsyncPeerid =========================================")
	st := make(chan struct{})
	call := new(MyCall)
	call.call = func(str string, err error) {
		fmt.Printf("********Call back*******[%v][%v]\n\n\n************************************************************\n\n", str, err)
		st <- struct{}{}

	}

	// peerid
	if ret, str := ipfs_lib.IpfsAsyncPeerid("", 5, call); ret != ipfs_lib.SUCCESS {
		fmt.Println(">>>>>>>>", str)
		go call.call("", nil)
	}
	<-st

	if ret, str := ipfs_lib.IpfsAsyncPeerid("hello-heipi", 5, call); ret != ipfs_lib.SUCCESS {
		fmt.Println(">>>>>>>>", str)
		go call.call("", nil)
	}
	<-st

	if ret, str := ipfs_lib.IpfsAsyncPeerid("QmV3wPSkRkwnLckMyFYFEvap4jUv36jm71BkuCX6Tqufbv", 5, call); ret != ipfs_lib.SUCCESS {
		fmt.Println(">>>>>>>>", str)
		go call.call("", nil)
	}
	<-st

	if ret, str := ipfs_lib.IpfsAsyncPeerid("QmV3wPSkRkwnLckMyFYFEvap4jUv36jm71BkuCX6TqufbV", 5, call); ret != ipfs_lib.SUCCESS {
		fmt.Println(">>>>>>>>", str)
		go call.call("", nil)
	}
	<-st

	// peerid
	if ret, str := ipfs_lib.IpfsAsyncPeerid("", 5, call); ret != ipfs_lib.SUCCESS {
		fmt.Println(">>>>>>>>", str)
		go call.call("", nil)
	}
	<-st

	fmt.Println("==================================================================================================")
	fmt.Println("==================================================================================================")
}

func main() {
	var wg sync.WaitGroup
	done := make(chan struct{}, 1)
	st := make(chan struct{})

	// async init
	ipfs_lib.InitApi()

	// path
	ipfs_lib.IpfsAsyncPath("ipfs_home")

	callinit := new(MyCall)
	callinit.call = func(str string, err error) {
		fmt.Println("init call")
	}

	// init
	ipfs_lib.IpfsAsyncInit(callinit)

	// daemon
	go func() {
		wg.Add(1)
		defer wg.Done()

		ipfs_lib.IpfsAsyncDaemon(callinit)
		done <- struct{}{}
	}()

	// id
	time.Sleep(15 * time.Second)

	call := new(MyCall)
	call.call = func(str string, err error) {
		fmt.Printf("********Call back*******[%v][%v]\n\n\n************************************************************\n\n", str, err)
		st <- struct{}{}

	}

	if ret, str := ipfs_lib.IpfsAsyncId(5, call); ret != ipfs_lib.SUCCESS {
		fmt.Println(">>>>>>>>", str)
		go call.call("", nil)
	}
	<-st

	ipfsAsyncPeerid()

	// // privkey
	// ipfs_lib.IpfsAsyncPrivkey("", 5, call)
	// <-st

	// ipfs_lib.IpfsAsyncPrivkey("mykey", 5, call)
	// <-st

	// ipfs_lib.IpfsAsyncPrivkey("", 5, call)
	// <-st

	// // add
	// ipfs_lib.IpfsAsyncAdd("QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn", "/xyz", "apimain.go", 5, call)
	// <-st

	// // move
	// ipfs_lib.IpfsAsyncAdd("QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn", "/xyz", "/zyx", 5, call)
	// <-st

	// // get
	// ipfs_lib.IpfsAsyncGet("QmRVZmwRKGKVZprrqxCLHAiuqEwA9casjUA57e8pKufXNi", "./getBlock", 5, call)
	// <-st

	// // query
	// ipfs_lib.IpfsAsyncQuery("QmRVZmwRKGKVZprrqxCLHAiuqEwA9casjUA57e8pKufXNi", "/", 5, call)
	// <-st

	// // delete
	// ipfs_lib.IpfsAsyncDelete("QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn", "/zyx", 5, call)
	// <-st

	// // shutdown
	// ipfs_lib.IpfsAsyncShutDown(call)
	// <-st

	wg.Wait()
	<-done
}
