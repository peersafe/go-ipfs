package ipfs_lib

import "fmt"

type IpfsCallBack interface {
	// Path(code int, reason string)
	Init(code int, reason string)
	Daemon(code int, reason string)
	ShutDown(code int, reason string)
	Id(code int, reason, id string)
	Add(code int, reason, new_root, ipfs_path, file_path, add_hash string)
	Delete(code int, reason, new_root, ipfs_path string)
	Move(code int, reason, new_root, src_path, dst_path string)
	Share(code int, reason, object_hash, share_name, new_hash string)
	Get(code int, reason, share_hash, save_path string)
	Query(code int, reason, object_hash, ipfs_path, query_result string)
	Merge(code int, reason, new_root, ipfs_path, share_hash string)
	PeerId(code int, reason, id string)
	PrivateKey(code int, reason, key string)
	Config(code int, reason, key, value string)
	Publish(code int, reason, object_hash, publish_hash string)
	RemotePin(code int, reason, peer_id, peer_key, object_hash string)
	Remotels(code int, reason, peer_id, peer_key, object_hash, ls_result string)
	RelayPin(code int, reason, relay_id, relay_key, peer_id, peer_key, object_hash string)
	ConnectPeer(code int, reason, peer_addr string)
	Progress(ipfs_path, old_hash string, types int, total, current int64)
}

const (
	ADD_TYPE int = iota
	GET_TYPE
)

func Ipfs_async_path(path string) string {
	res, str := IpfsAsyncPath(path)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_init(call IpfsCallBack) string {
	initApi()
	res, str := IpfsAsyncInit(call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

// func Ipfs_async_cmd_arm(cmd string, second int) string {
// 	res, str := IpfsAsyncCmdApi(cmd, second)
// 	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
// }

func Ipfs_async_daemon() string {
	res, str := IpfsAsyncDaemon()
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_shutdown() string {
	res, str := IpfsAsyncShutDown()
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_config(key, value string) string {
	res, str := IpfsAsyncConfig(key, value)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_id(second int) string {
	res, str := IpfsAsyncId(second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_peerid(new_id string, second int) string {
	res, str := IpfsAsyncPeerid(new_id, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_privkey(new_key string, second int) string {
	res, str := IpfsAsyncPrivkey(new_key, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_add(root_hash, ipfs_path, os_path string, second int) string {
	res, str := IpfsAsyncAdd(root_hash, ipfs_path, os_path, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_get(share_hash, os_path string, second int) string {
	res := IpfsAsyncGet(share_hash, os_path, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, "")
}

func Ipfs_async_publish(object_hash string, second int) string {
	res, str := IpfsAsyncPublish(object_hash, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_remotepin(peer_id, peer_key, object_hash string, second int) string {
	res, str := IpfsAsyncRemotepin(peer_id, peer_key, object_hash, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_remotels(peer_id, peer_key, object_hash string, second int) string {
	res, str := IpfsAsyncRemotels(peer_id, peer_key, object_hash, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_connectpeer(remote_peer string, second int) string {
	res, str := IpfsAsyncConnectPeer(remote_peer, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}
