package main

/*
#include <string.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/ipfs/go-ipfs/cmd/ipfs_mobile"
)

//export IpfsInit
func IpfsInit(path string) int {
	err := ipfsmobile.IpfsInit(path)
	if err != nil {
		return UNKOWN
	}
	return SUCCESS
}

//export IpfsAsyncDaemon
func IpfsAsyncDaemon(path string,
	cb_daemon unsafe.Pointer,
	cb_add unsafe.Pointer,
	cb_get unsafe.Pointer,
	cb_query unsafe.Pointer,
	cb_publish unsafe.Pointer,
	cb_connectpeer unsafe.Pointer,
	cb_message unsafe.Pointer) {

	call := caller{cb_daemon, cb_add, cb_get, cb_query, cb_publish, cb_connectpeer, cb_message}
	ipfsmobile.IpfsAsyncDaemon(path, call)
}

//export IpfsShutdown
func IpfsShutdown() int {
	err := ipfsmobile.IpfsShutdown()
	if err != nil {
		return UNKOWN
	}
	return SUCCESS
}

//export IpfsAsyncAdd
func IpfsAsyncAdd(os_path string, second int) *C.char {
	uid := ipfsmobile.IpfsAsyncAdd(os_path, second)
	return C.CString(uid)
}

//export IpfsDelete
func IpfsDelete(root_hash, ipfs_path string, second int) (new_root *C.char, retErr int) {
	n_root, err := ipfsmobile.IpfsDelete(root_hash, ipfs_path, second)
	new_root, retErr = C.CString(n_root), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsMove
func IpfsMove(root_hash, ipfs_src_path, ipfs_dst_path string, second int) (new_root *C.char, retErr int) {
	n_root, err := ipfsmobile.IpfsMove(root_hash, ipfs_src_path, ipfs_dst_path, second)
	new_root, retErr = C.CString(n_root), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsShare
func IpfsShare(object_hash, share_name string, sencond int) (new_hash *C.char, retErr int) {
	n_hash, err := ipfsmobile.IpfsShare(object_hash, share_name, sencond)
	new_hash, retErr = C.CString(n_hash), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsAsyncGet
func IpfsAsyncGet(share_hash, save_path string, second int) *C.char {
	uid := ipfsmobile.IpfsAsyncGet(share_hash, save_path, second)
	return C.CString(uid)
}

//export IpfsAsyncQuery
func IpfsAsyncQuery(object_hash, ipfs_path string, second int) {
	ipfsmobile.IpfsAsyncQuery(object_hash, ipfs_path, second)
}

//export IpfsMerge
func IpfsMerge(root_hash, ipfs_path, share_hash string, second int) (new_root *C.char, retErr int) {
	n_root, err := ipfsmobile.IpfsMerge(root_hash, ipfs_path, share_hash, second)
	new_root, retErr = C.CString(n_root), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsPeerid
func IpfsPeerid(new_id string, second int) (id *C.char, retErr int) {
	str, err := ipfsmobile.IpfsPeerid(new_id, second)
	id, retErr = C.CString(str), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsPrivkey
func IpfsPrivkey(new_key string, second int) (key *C.char, retErr int) {
	str, err := ipfsmobile.IpfsPrivkey(new_key, second)
	key, retErr = C.CString(str), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsAsyncPublish
func IpfsAsyncPublish(object_hash string, second int) {
	ipfsmobile.IpfsAsyncPublish(object_hash, second)
}

//export IpfsAsyncConnectpeer
func IpfsAsyncConnectpeer(peer_addr string, second int) {
	ipfsmobile.IpfsAsyncConnectpeer(peer_addr, second)
}

//export IpfsConfig
func IpfsConfig(key, value string) (retValue *C.char, retErr int) {
	str, err := ipfsmobile.IpfsConfig(key, value)
	retValue, retErr = C.CString(str), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsRemotepin
func IpfsRemotepin(peer_id, peer_key, object_hash string, second int) (retErr int) {
	err := ipfsmobile.IpfsRemotepin(peer_id, peer_key, object_hash, second)
	retErr = SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsRemotels
func IpfsRemotels(peer_id, peer_key, object_hash string, second int) (lsResult *C.char, retErr int) {
	str, err := ipfsmobile.IpfsRemotels(peer_id, peer_key, object_hash, second)
	lsResult, retErr = C.CString(str), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsMessage
func IpfsMessage(peer_id, peer_key, msg string) {
	ipfsmobile.IpfsAsyncMessage(peer_id, peer_key, msg)
}

//export IpfsCancel
func IpfsCancel(uuid string) {
	ipfsmobile.IpfsCancel(uuid)
}
