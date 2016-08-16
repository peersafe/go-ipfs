package main

import (
	"fmt"

	"github.com/ipfs/go-ipfs/cmd/ipfs_lib"
)

type CallBack struct {
}

func (call *CallBack) Init(code int, reason string) {
	fmt.Printf("func=[Init],code=[%v],reason=[%v]\n", code, reason)
}
func (call *CallBack) Daemon(code int, reason string) {
	fmt.Printf("func=[Daemon],code=[%v],reason=[%v]\n", code, reason)
}
func (call *CallBack) ShutDown(code int, reason string) {
	fmt.Printf("func=[ShutDown],code=[%v],reason=[%v]\n", code, reason)
}
func (call *CallBack) Id(code int, reason, id string) {
	fmt.Printf("func=[Id],code=[%v],reason=[%v],id=[%v]\n", code, reason, id)
	ipfs_lib.Ipfs_async_peerid(id, 10)
}
func (call *CallBack) Add(code int, reason, new_root, ipfs_path, file_path, add_hash string) {
	fmt.Printf("func=[Add],code=[%v],reason=[%v],new_root=[%v],ipfs_path=[%v],file_path=[%v],,add_hash=[%v]\n",
		code, reason, new_root, ipfs_path, file_path, add_hash)
}
func (call *CallBack) Delete(code int, reason, new_root, ipfs_path string) {
	fmt.Printf("func=[Delete],code=[%v],reason=[%v],new_root=[%v],ipfs_path=[%v]\n",
		code, reason, new_root, ipfs_path)
}
func (call *CallBack) Move(code int, reason, new_root, src_path, dst_path string) {
	fmt.Printf("func=[Move],code=[%v],reason=[%v],new_root=[%v],src_path=[%v],dst_path=[%v]\n",
		code, reason, new_root, src_path, dst_path)
}
func (call *CallBack) Share(code int, reason, object_hash, share_name, new_hash string) {
	fmt.Printf("func=[Share],code=[%v],reason=[%v],object_hash=[%v],share_name=[%v],new_hash=[%v]\n",
		code, reason, object_hash, share_name, new_hash)
}
func (call *CallBack) Get(code int, reason, share_hash, save_path string) {
	fmt.Printf("func=[Get],code=[%v],reason=[%v],share_hash=[%v],save_path=[%v]\n",
		code, reason, share_hash, save_path)
}
func (call *CallBack) Query(code int, reason, object_hash, ipfs_path, query_result string) {
	fmt.Printf("func=[Query],code=[%v],reason=[%v],object_hash=[%v],ipfs_path=[%v],query_result=[%v]\n",
		code, reason, object_hash, ipfs_path, query_result)

}
func (call *CallBack) Merge(code int, reason, new_root, ipfs_path, share_hash string) {
	fmt.Printf("func=[Merge],code=[%v],reason=[%v],new_root=[%v],ipfs_path=[%v],share_hash=[%v]\n",
		code, reason, new_root, ipfs_path, share_hash)

}
func (call *CallBack) PeerId(code int, reason, id string) {
	fmt.Printf("func=[PeerId],code=[%v],reason=[%v],id=[%v]\n",
		code, reason, id)

}
func (call *CallBack) PrivateKey(code int, reason, key string) {
	fmt.Printf("func=[PrivateKey],code=[%v],reason=[%v],key=[%v]\n",
		code, reason, key)
}
func (call *CallBack) Config(code int, reason, key, value string) {
	fmt.Printf("func=[Config],code=[%v],reason=[%v],key=[%v],value=[%v]\n",
		code, reason, key, value)
}
func (call *CallBack) Publish(code int, reason, object_hash, publish_hash string) {
	fmt.Printf("func=[Publish],code=[%v],reason=[%v],object_hash=[%v],publish_hash=[%v]\n",
		code, reason, object_hash, publish_hash)
}
func (call *CallBack) RemotePin(code int, reason, peer_id, peer_key, object_hash string) {
	fmt.Printf("func=[RemotePin],code=[%v],reason=[%v],peer_id=[%v],peer_key=[%v],object_hash=[%v]\n",
		code, reason, peer_id, peer_key, object_hash)
}
func (call *CallBack) Remotels(code int, reason, peer_id, peer_key, object_hash, ls_result string) {
	fmt.Printf("func=[Remotels],code=[%v],reason=[%v],peer_id=[%v],peer_key=[%v],object_hash=[%v],ls_result=[%v]\n",
		code, reason, peer_id, peer_key, object_hash, ls_result)
}
func (call *CallBack) RelayPin(code int, reason, relay_id, relay_key, peer_id, peer_key, object_hash string) {
	fmt.Printf("func=[RelayPin],code=[%v],reason=[%v],relay_id=[%v],relay_key=[%v],peer_id=[%v],peer_key=[%v],object_hash=[%v],ls_result=[%v]\n",
		code, reason, relay_id, relay_key, peer_id, peer_key, object_hash)
}
func (call *CallBack) ConnectPeer(code int, reason, peer_addr string) {
	fmt.Printf("func=[ConnectPeer],code=[%v],reason=[%v],peer_addr=[%v]\n",
		code, reason, peer_addr)
}
func (call *CallBack) Progress(ipfs_path, old_hash string, types int, total, current int64) {
	fmt.Printf("func=[Progress],ipfs_path=[%v],old_hash=[%v],types=[%v],total=[%v],current=[%v]\n",
		ipfs_path, old_hash, types, total, current)
}
