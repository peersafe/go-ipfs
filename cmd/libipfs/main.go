package main

/*
#include <string.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/ipfs/go-ipfs/cmd/ipfs_lib"
)

const hashLen int = 46
const keyLen int = 1596
const endsep = "\n"

const (
	errRet = -1
	sucRet = 0
)

type statInfo struct {
	Hash string
}

//export ipfs_init
func ipfs_init(out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsInit(); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_daemon
func ipfs_daemon(out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsDaemon(); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_id
func ipfs_id(second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsId(second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_add
func ipfs_add(root_hash, ipfs_path, os_path string, second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsAdd(root_hash, ipfs_path, os_path, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_delete
func ipfs_delete(root_hash, ipfs_path string, second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsDelete(root_hash, ipfs_path, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_move
func ipfs_move(root_hash, ipfs_path_src, ipfs_path_des string, second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsMove(root_hash, ipfs_path_src, ipfs_path_des, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_shard
func ipfs_shard(object_hash, shard_name string, second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsShard(object_hash, shard_name, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_get
func ipfs_get(shard_hash, os_path string, second int) int {
	return ipfs_lib.IpfsGet(shard_hash, os_path, second)
}

//export ipfs_query
func ipfs_query(object_hash, ipfs_path string, second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsQuery(object_hash, ipfs_path, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_merge
func ipfs_merge(root_hash, ipfs_path, shard_hash string, second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsMerge(root_hash, ipfs_path, shard_hash, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_peerid
func ipfs_peerid(new_id string, second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsPeerid(new_id, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_privkey
func ipfs_privkey(new_key string, second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsPrivkey(new_key, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_publish
func ipfs_publish(object_hash string, second int, out_res *C.char) int {
	object_hash = "/ipfs/" + object_hash
	if ret, str := ipfs_lib.IpfsPublish(object_hash, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_remotepin
func ipfs_remotepin(remote_peer, peer_key, object_hash string, second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsRemotepin(remote_peer, peer_key, object_hash, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_relaypin
func ipfs_relaypin(relay_peer, relay_key, remote_peer, peer_key, object_hash string, second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsRelaypin(relay_peer, relay_key, remote_peer, peer_key, object_hash, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_remotels
func ipfs_remotels(remote_peer, peer_key, object_hash string, second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsRemotels(remote_peer, peer_key, object_hash, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_connectpeer
func ipfs_connectpeer(remote_peer string, second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsConnectPeer(remote_peer, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_config
func ipfs_config(key, value string, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsConfig(key, value); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

//export ipfs_cmd
func ipfs_cmd(cmd string, second int, out_res *C.char) int {
	if ret, str := ipfs_lib.IpfsCmdApi(cmd, second); ret != errRet {
		cs := unsafe.Pointer(C.CString(str))
		C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
		C.free(cs)
		return ret
	}
	return errRet
}

// main roadmap:
// - parse the commandline to get a cmdInvocation
// - if user requests, help, print it and exit.
// - run the command invocation
// - output the response
// - if anything fails, print error, maybe with help
func main() {
}
