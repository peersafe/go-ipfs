// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hello is a trivial package for gomobile bind example.
package ipfsmobile

import ipfslib "github.com/ipfs/go-ipfs/cmd/ipfs_lib"

type IpfsCallBack interface {
	ipfslib.IpfsCallBack
}

func IpfsPath(path string) string {
	return ipfslib.Ipfs_async_path(path)
}

func IpfsInit(call IpfsCallBack) string {
	return ipfslib.Ipfs_async_init(call)
}

func IpfsDaemon() string {
	return ipfslib.Ipfs_async_daemon()
}

func IpfsShutdown() string {
	return ipfslib.Ipfs_async_shutdown()
}

func IpfsId(second int) string {
	return ipfslib.Ipfs_async_id(second)
}

func IpfsPeerid(new_id string, second int) string {
	return ipfslib.Ipfs_async_peerid(new_id, second)
}

func IpfsPrivkey(new_key string, second int) string {
	return ipfslib.Ipfs_async_privkey(new_key, second)
}

func IpfsAdd(os_path string, second int) string {
	return ipfslib.Ipfs_async_add(os_path, second)
}

func IpfsGet(object_hash, os_path string, second int) string {
	return ipfslib.Ipfs_async_get(object_hash, os_path, second)
}

func IpfsPublish(object_hash string, second int) string {
	return ipfslib.Ipfs_async_publish(object_hash, second)
}

func IpfsConfig(key, value string) string {
	return ipfslib.Ipfs_async_config(key, value)
}

func IpfsRemotepin(peer_id, peer_key, object_hash string, second int) string {
	return ipfslib.Ipfs_async_remotepin(peer_id, peer_key, object_hash, second)
}

func IpfsRemotels(peer_id, peer_key, object_hash string, second int) string {
	return ipfslib.Ipfs_async_remotels(peer_id, peer_key, object_hash, second)
}

func IpfsConnectpeer(remote_peer string, second int) string {
	return ipfslib.Ipfs_async_connectpeer(remote_peer, second)
}
