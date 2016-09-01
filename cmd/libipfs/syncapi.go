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

//export ipfs_path
func ipfs_path(path string, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsPath(path)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_init
func ipfs_init(out_res *C.char) int {
	ret, str := ipfs_lib.IpfsInit()
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_daemon
func ipfs_daemon(out_res *C.char) int {
	ret, str := ipfs_lib.IpfsDaemon()
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_shutdown
func ipfs_shutdown(out_res *C.char) int {
	ret, str := ipfs_lib.IpfsShutDown()
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_id
func ipfs_id(second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsId(second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_add
func ipfs_add(root_hash, ipfs_path, os_path string, second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsAdd(root_hash, ipfs_path, os_path, second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_delete
func ipfs_delete(root_hash, ipfs_path string, second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsDelete(root_hash, ipfs_path, second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_move
func ipfs_move(root_hash, ipfs_path_src, ipfs_path_des string, second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsMove(root_hash, ipfs_path_src, ipfs_path_des, second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_shard
func ipfs_shard(object_hash, shard_name string, second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsShard(object_hash, shard_name, second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_get
func ipfs_get(shard_hash, os_path string, second int) int {
	return ipfs_lib.IpfsGet(shard_hash, os_path, second)
}

//export ipfs_query
func ipfs_query(object_hash, ipfs_path string, second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsQuery(object_hash, ipfs_path, second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_merge
func ipfs_merge(root_hash, ipfs_path, shard_hash string, second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsMerge(root_hash, ipfs_path, shard_hash, second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_peerid
func ipfs_peerid(new_id string, second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsPeerid(new_id, second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_privkey
func ipfs_privkey(new_key string, second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsPrivkey(new_key, second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_publish
func ipfs_publish(object_hash string, second int, out_res *C.char) int {
	object_hash = "/ipfs/" + object_hash
	ret, str := ipfs_lib.IpfsPublish(object_hash, second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_remotepin
func ipfs_remotepin(remote_peer, peer_key, object_hash string, second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsRemotepin(remote_peer, peer_key, object_hash, second)
	if str != "" {
		goStringToChar(str, out_res)

	}
	return ret
}

//export ipfs_relaypin
func ipfs_relaypin(relay_peer, relay_key, remote_peer, peer_key, object_hash string, second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsRelaypin(relay_peer, relay_key, remote_peer, peer_key, object_hash, second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_remotels
func ipfs_remotels(remote_peer, peer_key, object_hash string, second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsRemotels(remote_peer, peer_key, object_hash, second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_connectpeer
func ipfs_connectpeer(remote_peer string, second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsConnectPeer(remote_peer, second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_config
func ipfs_config(key, value string, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsConfig(key, value)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

//export ipfs_cmd
func ipfs_cmd(cmd string, second int, out_res *C.char) int {
	ret, str := ipfs_lib.IpfsCmdApi(cmd, second)
	if str != "" {
		goStringToChar(str, out_res)
	}
	return ret
}

func goStringToChar(str string, out_res *C.char) {
	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
}
