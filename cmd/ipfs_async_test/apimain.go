package main

import (
	"fmt"
	"time"
	"sync"

	ipfs_lib "github.com/ipfs/go-ipfs/cmd/ipfs_lib"
)

type MyCall struct {
}

func (call *MyCall) Call(result string, err error) {
	fmt.Printf("result=%v,=======,err=%v\n", result, err)
}

func main() {
	call := new(MyCall)
	var wg sync.WaitGroup
	done :=make(chan struct{})

	// async init
	ipfs_lib.InitApi()

	// path
	ipfs_lib.IpfsAsyncPath("ipfs_home")

	// init
	ipfs_lib.IpfsAsyncInit(call)

	// daemon
	go func() {
		wg.Add(1)
		defer wg.Done()
		

		ipfs_lib.IpfsAsyncDaemon(call)
		done <- struct{}{}
	}()

	// id
	go func () {
		wg.Add(1)
		defer wg.Done()

		// wait for daemon start
		time.Sleep(10*time.Second)
		ipfs_lib.IpfsAsyncId(5,call)
	}()

	// peerid
	go func () {
		wg.Add(1)
		defer wg.Done()

		// wait for daemon start
		time.Sleep(10*time.Second)
		ipfs_lib.IpfsAsyncPeerid("",5,call)

		ipfs_lib.IpfsAsyncPeerid("hello-heipi",5,call)

		ipfs_lib.IpfsAsyncPeerid("",5,call)
	}()

	// privkey
	go func () {
		wg.Add(1)
		defer wg.Done()

		// wait for daemon start
		time.Sleep(10*time.Second)
		ipfs_lib.IpfsAsyncPrivkey("",5,call)

		ipfs_lib.IpfsAsyncPrivkey("mykey",5,call)

		ipfs_lib.IpfsAsyncPrivkey("",5,call)
	}()

	// add
	go func(){
		wg.Add(1)
		defer wg.Done()

		// wait for daemon start
		time.Sleep(10*time.Second)
		ipfs_lib.IpfsAsyncAdd("QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn","/xyz","apimain.go",5,call)
	}()

	// move
	go func(){
		wg.Add(1)
		defer wg.Done()

		// wait for add done 
		time.Sleep(15*time.Second)
		ipfs_lib.IpfsAsyncAdd("QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn","/xyz","/zyx",5,call)
	}()

	// get
	go func(){
		wg.Add(1)
		defer wg.Done()

		// wait for daemon start 
		time.Sleep(10*time.Second)
		ipfs_lib.IpfsAsyncGet("QmRVZmwRKGKVZprrqxCLHAiuqEwA9casjUA57e8pKufXNi","./getBlock",5,call)
	}()

	// query
	go func(){
		wg.Add(1)
		defer wg.Done()

		// wait for daemon start
		time.Sleep(10*time.Second)
		ipfs_lib.IpfsAsyncQuery("QmRVZmwRKGKVZprrqxCLHAiuqEwA9casjUA57e8pKufXNi","/",5,call)
	}()
	

	// delete
	go func(){
		wg.Add(1)
		defer wg.Done()

		// wait for move done
		time.Sleep(20*time.Second)
		ipfs_lib.IpfsAsyncDelete("QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn","/zyx",5,call)
	}()
	
	// shutdown
	go func() {
		wg.Add(1)
		defer wg.Done()
		// wait for all goroutine done
		time.Sleep(30 * time.Second)
		ipfs_lib.IpfsAsyncShutDown(call)
	}()
	wg.Wait()
	<-done
}

