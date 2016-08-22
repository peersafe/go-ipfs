package main

import (
	"fmt"

	ipfs_mobile "github.com/ipfs/go-ipfs/cmd/ipfs_mobile"
)

type CallBack struct {
}

const (
	Tab = "===================================================================================="
)

func (call *CallBack) Daemon(status int, err error) {
	fmt.Println(Tab)
	fmt.Printf("func=[Daemon],status=[%v],err=[%v]\n", status, err)
	if status == 0 {
		fmt.Println("Daemon start...")
		ipfs_mobile.IpfsAsyncConnectpeer("/ip4/101.201.40.124/tcp/40001/ipfs/QmZDYAhmMDtnoC6XZRw8R1swgoshxKvXDA9oQF97AYkPZc", 5)
	}
	if status == 1 {
		fmt.Println("Daemon shutdown...")
		done <- struct{}{}
	}
}

func (call *CallBack) Add(uid, add_hash string, pos int, err error) {
	fmt.Println(Tab)
	fmt.Printf("func=[Add],uid=[%v],add_hash=[%v],pos=[%v],err=[%v]\n",
		uid, add_hash, pos, err)
	if add_hash == "" {
		fmt.Printf("uid=%v, ==================== process=%v %\n", uid, pos)
		return
	}
	fmt.Printf("uid=%v, ==================== process=%v %, add_hash=%v\n", uid, pos, add_hash)
	err = ipfs_mobile.IpfsShutdown()
	if err != nil {
		fmt.Println("func=[IpfsAsyncShutdown],err= ", err)
		return
	}
	fmt.Println("func=[IpfsAsyncShutdown], Over")
}

func (call *CallBack) Get(uid string, pos int, err error) {
	fmt.Println(Tab)
	fmt.Printf("func=[Get],uid=[%v],pos=[%v],err=[%v]\n",
		uid, pos, err)
	fmt.Printf("uid=%v, ==================== process=%v %\n", uid, pos)

	if pos != 100 {
		return
	}

	// add
	add_uid := ipfs_mobile.IpfsAsyncAdd("apimain.go", 5)
	fmt.Println("func=[IpfsAsyncAdd],uid= ", add_uid)
}

func (call *CallBack) Query(object_hash, ipfs_path, result string, err error) {
	fmt.Println(Tab)
	fmt.Printf("func=[Query],object_hash=[%v],ipfs_path=[%v],result=[%v],err=[%v]\n",
		object_hash, ipfs_path, result, err)
}

func (call *CallBack) Publish(publish_hash string, err error) {
	fmt.Println(Tab)
	fmt.Printf("func=[Publish],publish_hash=[%v],err=[%v]\n",
		publish_hash, err)
}

func (call *CallBack) ConnectPeer(peer_addr string, err error) {
	fmt.Println(Tab)
	fmt.Printf("func=[ConnectPeer],peer_addr=[%v],err=[%v]\n",
		peer_addr, err)
	if err != nil {
		fmt.Println(err)
		return
	}
	id, err := ipfs_mobile.IpfsPeerid("", 5)
	if err != nil {
		fmt.Println("func=[IpfsAsyncPeerid],err=", err)
		return
	}
	fmt.Println("func=[IpfsAsyncPeerid],id= ", id)

	uid := ipfs_mobile.IpfsAsyncGet("QmfJ6DFC8pTv72JLKzdE1q9LLv1hGdqeUafZ9SFgXWY1kK", "test", 60)
	fmt.Println("func=[IpfsAsyncGet],uid= ", uid)
}

func (call *CallBack) Message(peer_id, passwd, msg string, err error) {
	fmt.Println(Tab)
	fmt.Printf("func=[Message],peer_id=[%v],passwd=[%v],msg=[%v],err=[%v]\n",
		peer_id, passwd, msg, err)
	if err != nil {
		fmt.Println(err)
		return
	}
}
