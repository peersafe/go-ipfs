package main

/*
typedef void (*cb_daemon)(int, int);
void Daemon(cb_daemon fn, int status, int ret) {
	fn(status, ret);
}
typedef void (*cb_add)(char*, char*, int, int);
void Add(cb_add fn, char* uid, char* hash, int pos, int ret) {
	fn(uid, hash, pos, ret);
}
typedef void (*cb_get)(char*, int, int);
void Get(cb_get fn, char* uid, int pos, int ret) {
	fn(uid, pos, ret);
}
typedef void (*cb_query)(char*, char*, char*, int);
void Query(cb_query fn, char* root_hash, char* ipfs_path, char* result, int ret) {
	fn(root_hash, ipfs_path, result, ret);
}
typedef void (*cb_publish)(char*, int);
void Publish(cb_publish fn, char* publish_hash, int ret) {
	fn(publish_hash, ret);
}
typedef void (*cb_connectpeer)(char *, int);
void ConnectPeer(cb_connectpeer fn, char* peer_addr, int ret) {
	fn(peer_addr, ret);
}
typedef void (*cb_message)(char *, char *, char *, int);
void Message(cb_message fn, char* peer_id, char* peer_key, char* msg, int ret) {
	fn(peer_id, peer_key, msg, ret);
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

type caller struct {
	cbdaemon      unsafe.Pointer
	cbadd         unsafe.Pointer
	cbget         unsafe.Pointer
	cbquery       unsafe.Pointer
	cbpublish     unsafe.Pointer
	cbconnectpeer unsafe.Pointer
	cbmessage     unsafe.Pointer
}

func (c caller) Daemon(status int, err string) {
	ret := SUCCESS
	if err != "" {
		fmt.Println("[Daemon] error:", err)
		ret = UNKOWN
	}
	fn := C.cb_daemon(c.cbdaemon)
	C.Daemon(fn, C.int(status), C.int(ret))
}

func (c caller) Add(uid, hash string, pos int, err string) {
	ret := SUCCESS

	if err != "" {
		ret = UNKOWN
		if err == "timeout" {
			ret = TIMEOUT
		}
		fmt.Println("[Add] error:", err)
	}

	fn := C.cb_add(c.cbadd)
	C.Add(fn, C.CString(uid), C.CString(hash), C.int(pos), C.int(ret))
}

func (c caller) Get(uid string, pos int, err string) {
	ret := SUCCESS
	if err != "" {
		ret = UNKOWN
		if err == "timeout" {
			ret = TIMEOUT
		}
		fmt.Println("[Get] error:", err)
	}
	fn := C.cb_get(c.cbget)
	C.Get(fn, C.CString(uid), C.int(pos), C.int(ret))
}

func (c caller) Query(root_hash, ipfs_path, result string, err string) {
	ret := SUCCESS
	if err != "" {
		fmt.Println("[Query] error:", err)
		ret = UNKOWN
	}
	fn := C.cb_query(c.cbquery)
	C.Query(fn, C.CString(root_hash), C.CString(ipfs_path), C.CString(result), C.int(ret))
}

func (c caller) Publish(publish_hash string, err string) {
	ret := SUCCESS
	if err != "" {
		fmt.Println("[Publish] error:", err)
		ret = UNKOWN
	}
	fn := C.cb_publish(c.cbpublish)
	C.Publish(fn, C.CString(publish_hash), C.int(ret))
}

func (c caller) ConnectPeer(peer_addr string, err string) {
	ret := SUCCESS
	if err != "" {
		fmt.Println("[ConnectPeer] error:", err)
		ret = UNKOWN
	}
	fn := C.cb_connectpeer(c.cbconnectpeer)
	C.ConnectPeer(fn, C.CString(peer_addr), C.int(ret))
}

func (c caller) Message(peer_id, peer_key, msg string, err string) {
	ret := SUCCESS
	if err != "" {
		fmt.Println("[Message] error:", err)
		ret = UNKOWN
	}
	fn := C.cb_message(c.cbmessage)
	C.Message(fn, C.CString(peer_id), C.CString(peer_key), C.CString(msg), C.int(ret))
}
