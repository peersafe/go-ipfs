package main

import (
	"fmt"
	"time"

	ipfs_mobile "github.com/ipfs/go-ipfs/cmd/ipfs_mobile"
)

type CallBack struct {
}

const (
	Tab = "===================================================================================="
)

// func (call *CallBack) Daemon(status int, err string) {
// 	fmt.Println(Tab)
// 	fmt.Printf("func=[Daemon],status=[%v],err=[%v]\n", status, err)
// 	if status == 0 {
// 		fmt.Println("Daemon start...")
// 		ipfs_mobile.IpfsAsyncConnectpeer("/ip4/101.201.40.124/tcp/40001/ipfs/QmZDYAhmMDtnoC6XZRw8R1swgoshxKvXDA9oQF97AYkPZc", 5)
// 	}
// 	if status == 1 {
// 		fmt.Println("Daemon shutdown...")
// 		done <- struct{}{}
// 	}
// }

func (call *CallBack) Daemon(status int, err string) {
	fmt.Println(Tab)
	fmt.Printf("func=[Daemon],status=[%v],err=[%v]\n", status, err)
	if status == 0 {
		fmt.Println("Daemon start...")

		// conncet
		ipfs_mobile.IpfsAsyncConnectpeer("/ip4/172.16.154.129/tcp/4001/ipfs/QmeNgHawAonsK2uAYLMZP5TAa395DdzHUzoPYHgEv3khez", 5)

		// 	add
		add_uid := ipfs_mobile.IpfsAsyncGet("QmVvvSWZK3csra9QbFMqUrkXtyx1FeSuERpRDpwhJSYkoz", "test", 30)
		fmt.Println("func=[IpfsAsyncAdd],uid= ", add_uid)

		time.Sleep(2 * time.Second)

		// cancel
		ipfs_mobile.IpfsCancel(add_uid)

		shutdown()
	}
	if status == 1 {
		fmt.Println("Daemon shutdown...")
		done <- struct{}{}
	}
}

func (call *CallBack) Add(uid, add_hash string, pos int, err string) {
	fmt.Println(Tab)
	fmt.Printf("func=[Add],uid=[%v],add_hash=[%v],pos=[%v],err=[%v]\n",
		uid, add_hash, pos, err)
	fmt.Printf("uid=%v, ==================== process=%v %\n", uid, pos)
	if pos != 100 {
		return
	}

	// config test
	peerid, e := ipfs_mobile.IpfsConfig("Identity.PeerID", "")
	if e != nil {
		fmt.Println("func=[IpfsAsyncConfig],err= ", e)
		return
	}
	fmt.Println("func=[IpfsAsyncConfig],peerid= ", peerid)
	peerkey, e := ipfs_mobile.IpfsConfig("Identity.Secret", "")
	if e != nil {
		fmt.Println("func=[IpfsAsyncConfig],err= ", e)
		return
	}
	fmt.Println("func=[IpfsAsyncConfig],peerkey= ", peerkey)

	peerid, e = ipfs_mobile.IpfsPeerid("", 5)
	if e != nil {
		fmt.Println("func=[IpfsPeerid],err= ", e)
		return
	}
	fmt.Println("func=[IpfsPeerid],peerid= ", peerid)

	peerkey, e = ipfs_mobile.IpfsPrivkey("", 5)
	if e != nil {
		fmt.Println("func=[IpfsPrivkey],err= ", e)
		return
	}
	fmt.Println("func=[IpfsPrivkey],peerkey= ", peerkey)

	// remotemsg test
	// msg := `{"type":"remotepin","hash":"` + add_hash + `","msg_from_peerid":"` + peerid + `","msg_from_peerkey":"` + peerkey + `"}`
	// ipfs_mobile.IpfsAsyncMessage("", "", msg)

	e = ipfs_mobile.IpfsShutdown()
	if e != nil {
		fmt.Println("func=[IpfsAsyncShutdown],err= ", e)
		return
	}
	fmt.Println("func=[IpfsAsyncShutdown], Over")
}

func (call *CallBack) Get(uid string, pos int, err string) {
	fmt.Println(Tab)
	fmt.Printf("func=[Get],uid=[%v],pos=[%v],err=[%v]\n",
		uid, pos, err)
	fmt.Printf("uid=%v, ==================== process=%v %\n", uid, pos)
	if err != "" {
		// shutdown()
		return
	}
	if pos != 100 {
		return
	}

	// shutdown()
	// // add
	// add_uid := ipfs_mobile.IpfsAsyncAdd("apimain.go", 30)
	// fmt.Println("func=[IpfsAsyncAdd],uid= ", add_uid)
}

func shutdown() {
	e := ipfs_mobile.IpfsShutdown()
	if e != nil {
		fmt.Println("func=[IpfsAsyncShutdown],err= ", e)
		return
	}
	fmt.Println("func=[IpfsAsyncShutdown], Over")
}

func (call *CallBack) Query(object_hash, ipfs_path, result string, err string) {
	fmt.Println(Tab)
	fmt.Printf("func=[Query],object_hash=[%v],ipfs_path=[%v],result=[%v],err=[%v]\n",
		object_hash, ipfs_path, result, err)
}

func (call *CallBack) Publish(publish_hash string, err string) {
	fmt.Println(Tab)
	fmt.Printf("func=[Publish],publish_hash=[%v],err=[%v]\n",
		publish_hash, err)
}

func (call *CallBack) ConnectPeer(peer_addr string, err string) {
	fmt.Println(Tab)
	fmt.Printf("func=[ConnectPeer],peer_addr=[%v],err=[%v]\n",
		peer_addr, err)
	if err != "" {
		fmt.Println(err)
		return
	}
	// id, e := ipfs_mobile.IpfsPeerid("", 5)
	// if e != nil {
	// 	fmt.Println("func=[IpfsAsyncPeerid],err=", err)
	// 	return
	// }
	// fmt.Println("func=[IpfsAsyncPeerid],id= ", id)

	// uid := ipfs_mobile.IpfsAsyncGet("QmfJ6DFC8pTv72JLKzdE1q9LLv1hGdqeUafZ9SFgXWY1kK", "test", 60)
	// fmt.Println("func=[IpfsAsyncGet],uid= ", uid)
}

func (call *CallBack) Message(peer_id, passwd, msg string, err string) {
	fmt.Println(Tab)
	fmt.Printf("func=[Message],peer_id=[%v],passwd=[%v],msg=[%v],err=[%v]\n",
		peer_id, passwd, msg, err)
	if err != "" {
		fmt.Println(err)
		return
	}
}
