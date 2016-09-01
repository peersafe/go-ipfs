package callback

var GlobalCallBack IpfsCallBack

type IpfsCallBack interface {
	Daemon(status int, err string)
	Add(uid, hash string, pos int, err string)
	Get(uid string, pos int, err string)
	Query(root_hash, ipfs_path, result string, err string)
	Publish(publish_hash string, err string)
	ConnectPeer(peer_addr string, err string)
	Message(peer_id, peer_key, msg string, err string)
}
