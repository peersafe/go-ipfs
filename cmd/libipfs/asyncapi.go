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
	// memcpy for C lib
	ipfsPath := []byte(path)

	ipfsmobile.IpfsInit(string(ipfsPath))
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

	// memcpy for C lib
	ipfsPath := []byte(path)

	call := caller{cb_daemon, cb_add, cb_get, cb_query, cb_publish, cb_connectpeer, cb_message}
	ipfsmobile.IpfsAsyncDaemon(string(ipfsPath), call)
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
	// memcpy for C lib
	osPath := []byte(os_path)

	uid := ipfsmobile.IpfsAsyncAdd(string(osPath), second)
	return C.CString(uid)
}

//export IpfsDelete
func IpfsDelete(root_hash, ipfs_path string, second int) (new_root *C.char, retErr int) {
	// memcpy for C lib
	rootHash := []byte(root_hash)
	ipfsPath := []byte(ipfs_path)

	n_root, err := ipfsmobile.IpfsDelete(string(rootHash), string(ipfsPath), second)
	new_root, retErr = C.CString(n_root), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsMove
func IpfsMove(root_hash, ipfs_src_path, ipfs_dst_path string, second int) (new_root *C.char, retErr int) {
	// memcpy for C lib
	rootHash := []byte(root_hash)
	ipfsSrcPath := []byte(ipfs_src_path)
	ipfsDstPath := []byte(ipfs_dst_path)

	n_root, err := ipfsmobile.IpfsMove(string(rootHash), string(ipfsSrcPath), string(ipfsDstPath), second)
	new_root, retErr = C.CString(n_root), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsShare
func IpfsShare(object_hash, share_name string, sencond int) (new_hash *C.char, retErr int) {
	// memcpy for C lib
	objectHash := []byte(object_hash)
	shareName := []byte(share_name)

	n_hash, err := ipfsmobile.IpfsShare(string(objectHash), string(shareName), sencond)
	new_hash, retErr = C.CString(n_hash), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsAsyncGet
func IpfsAsyncGet(share_hash, save_path string, second int) *C.char {
	// memcpy for C lib
	shareHash := []byte(share_hash)
	savePath := []byte(save_path)

	uid := ipfsmobile.IpfsAsyncGet(string(shareHash), string(savePath), second)
	return C.CString(uid)
}

//export IpfsAsyncQuery
func IpfsAsyncQuery(object_hash, ipfs_path string, second int) {
	// memcpy for C lib
	rootHash := []byte(object_hash)
	ipfsPath := []byte(ipfs_path)

	ipfsmobile.IpfsAsyncQuery(string(rootHash), string(ipfsPath), second)
}

//export IpfsQuery
func IpfsQuery(object_hash, ipfs_path string, second int) (result *C.char, retErr int) {
	// memcpy for C lib
	rootHash := []byte(object_hash)
	ipfsPath := []byte(ipfs_path)

	queryResult, err := ipfsmobile.IpfsQuery(string(rootHash), string(ipfsPath), second)
	result, retErr = C.CString(queryResult), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsMerge
func IpfsMerge(root_hash, ipfs_path, share_hash string, second int) (new_root *C.char, retErr int) {
	// memcpy for C lib
	rootHash := []byte(root_hash)
	ipfsPath := []byte(ipfs_path)
	shareHash := []byte(share_hash)

	n_root, err := ipfsmobile.IpfsMerge(string(rootHash), string(ipfsPath), string(shareHash), second)
	new_root, retErr = C.CString(n_root), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsPeerid
func IpfsPeerid(new_id string, second int) (id *C.char, retErr int) {
	// memcpy for C lib
	newId := []byte(new_id)

	str, err := ipfsmobile.IpfsPeerid(string(newId), second)
	id, retErr = C.CString(str), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsPrivkey
func IpfsPrivkey(new_key string, second int) (key *C.char, retErr int) {
	// memcpy for C lib
	newKey := []byte(new_key)

	str, err := ipfsmobile.IpfsPrivkey(string(newKey), second)
	key, retErr = C.CString(str), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsAsyncPublish
func IpfsAsyncPublish(object_hash string, second int) {
	// memcpy for C lib
	objectHash := []byte(object_hash)

	ipfsmobile.IpfsAsyncPublish(string(objectHash), second)
}

//export IpfsAsyncConnectpeer
func IpfsAsyncConnectpeer(peer_addr string, second int) {
	// memcpy for C lib
	peerAddr := []byte(peer_addr)

	ipfsmobile.IpfsAsyncConnectpeer(string(peerAddr), second)
}

//export IpfsConfig
func IpfsConfig(key, value string) (retValue *C.char, retErr int) {
	// memcpy for C lib
	k := []byte(key)
	v := []byte(value)

	str, err := ipfsmobile.IpfsConfig(string(k), string(v))
	retValue, retErr = C.CString(str), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsRemotepin
func IpfsRemotepin(peer_id, peer_key, object_hash string, second int) (retErr int) {
	// memcpy for C lib
	peerId := []byte(peer_id)
	peerKey := []byte(peer_key)
	objectHash := []byte(object_hash)

	err := ipfsmobile.IpfsRemotepin(string(peerId), string(peerKey), string(objectHash), second)
	retErr = SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsRemotels
func IpfsRemotels(peer_id, peer_key, object_hash string, second int) (lsResult *C.char, retErr int) {
	// memcpy for C lib
	peerId := []byte(peer_id)
	peerKey := []byte(peer_key)
	objectHash := []byte(object_hash)

	str, err := ipfsmobile.IpfsRemotels(string(peerId), string(peerKey), string(objectHash), second)
	lsResult, retErr = C.CString(str), SUCCESS
	if err != nil {
		retErr = UNKOWN
	}
	return
}

//export IpfsMessage
func IpfsMessage(peer_id, peer_key, msg string) int {
	// memcpy for C lib
	peerId := []byte(peer_id)
	peerKey := []byte(peer_key)
	msgs := []byte(msg)

	return ipfsmobile.IpfsAsyncMessage(string(peerId), string(peerKey), string(msgs))
}

//export IpfsCancel
func IpfsCancel(uuid string) {
	// memcpy for C lib
	uid := []byte(uuid)

	ipfsmobile.IpfsCancel(string(uid))
}
