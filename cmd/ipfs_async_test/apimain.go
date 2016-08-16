package main

import (
	"sync"
	"time"

	ipfs_lib "github.com/ipfs/go-ipfs/cmd/ipfs_lib"
)

func main() {
	var wg sync.WaitGroup
	done := make(chan struct{})
	defer close(done)

	// path
	ipfs_lib.Ipfs_async_path("ipfs_home")

	// init
	callback := new(CallBack)
	ipfs_lib.Ipfs_async_init(callback)

	// daemon
	go func() {
		wg.Add(1)
		defer wg.Done()

		ipfs_lib.Ipfs_async_daemon()
		done <- struct{}{}
	}()

	// id
	// go func() {
	// 	wg.Add(1)
	// 	defer wg.Done()

	// 	// wait for daemon start
	// 	time.Sleep(15)

	// 	ipfs_lib.Ipfs_async_id(5)
	// }()

	// peerid
	// go func() {
	// 	wg.Add(1)
	// 	defer wg.Done()

	// 	// wait for daemon start
	// 	time.Sleep(15 * time.Second)

	// 	ipfs_lib.Ipfs_async_peerid("", 5)
	// 	ipfs_lib.Ipfs_async_peerid("QmeVm9QSUMxYa2CHxAPj5UpXkfLchezgZivn6XPSc11111", 5)
	// 	ipfs_lib.Ipfs_async_peerid("", 5)
	// }()

	// privkey
	// go func() {
	// 	wg.Add(1)
	// 	defer wg.Done()

	// 	// wait for daemon start
	// 	time.Sleep(15 * time.Second)
	// 	ipfs_lib.Ipfs_async_privkey("", 5)
	// 	ipfs_lib.Ipfs_async_privkey("mykey", 5)
	// 	ipfs_lib.Ipfs_async_privkey("", 5)
	// }()

	// add
	go func() {
		wg.Add(1)
		defer wg.Done()

		// wait for daemon start
		time.Sleep(15 * time.Second)
		ipfs_lib.Ipfs_async_add("QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn", "/", "apimain.go", 5)
	}()
	/*
		// move
		go func() {
			wg.Add(1)
			defer wg.Done()

			// wait for add done
			time.Sleep(15 * time.Second)
			ipfs_lib.IpfsAsyncAdd("QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn", "/xyz", "/zyx", 5, call)
		}()

			// get
			go func() {
				wg.Add(1)
				defer wg.Done()

				// wait for daemon start
				time.Sleep(10 * time.Second)
				ipfs_lib.IpfsAsyncGet("QmRVZmwRKGKVZprrqxCLHAiuqEwA9casjUA57e8pKufXNi", "getBlock", 5, call)
			}()

			// query
			go func() {
				wg.Add(1)
				defer wg.Done()

				// wait for daemon start
				time.Sleep(10 * time.Second)
				ipfs_lib.IpfsAsyncQuery("QmRVZmwRKGKVZprrqxCLHAiuqEwA9casjUA57e8pKufXNi", "/", 5, call)
			}()

			// delete
			go func() {
				wg.Add(1)
				defer wg.Done()

				// wait for move done
				time.Sleep(20 * time.Second)
				ipfs_lib.IpfsAsyncDelete("QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn", "/zyx", 5, call)
			}()
	*/

	// shutdown
	go func() {
		wg.Add(1)
		defer wg.Done()
		// wait for all goroutine done
		time.Sleep(30 * time.Second)
		ipfs_lib.Ipfs_async_shutdown()
	}()
	<-done
	wg.Wait()
}
