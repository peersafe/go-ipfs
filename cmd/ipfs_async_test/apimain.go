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

	// init
	ipfs_lib.InitApi()

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
		time.Sleep(5*time.Second)
		ipfs_lib.IpfsAsyncId(5,call)
	}()

	// add
	// go func(){
	// 	wg.Add(1)
	// 	defer wg.Done()

	// 	// wait for daemon start
	// 	time.Sleep(5*time.Second)
	// 	ipfs_lib.IpfsAsyncAdd("QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn","/xyz","apimain.go",5,call)
	// }()
	
	// shutdown
	go func() {
		wg.Add(1)
		defer wg.Done()
		// wait for all goroutine done
		time.Sleep(10 * time.Second)
		ipfs_lib.IpfsAsyncShutDown(call)
	}()
	wg.Wait()
	<-done
}

