// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hello is a trivial package for gomobile bind example.
package ipfsmobileand

import (
	ipfslib "github.com/ipfs/go-ipfs/cmd/ipfs_lib"
)

func IpfsCmd(cmd string, second int) string {
	return ipfslib.Ipfs_cmd_arm(cmd, second)
}

func IpfsPath(path string) string {
	return ipfslib.Ipfs_path(path)
}

func IpfsInit() string {
	return ipfslib.Ipfs_init()
}

func IpfsDaemon() string {
	return ipfslib.Ipfs_daemon()
}

func IpfsShutdown() string {
	return ipfslib.Ipfs_shutdown()
}

func IpfsId(second int) string {
	return ipfslib.Ipfs_id(second)
}

func IpfsPeerid(new_id string, second int) string {
	return ipfslib.Ipfs_peerid(new_id, second)
}

func IpfsPrivkey(new_key string, second int) string {
	return ipfslib.Ipfs_privkey(new_key, second)
}

func IpfsAdd(os_path string, second int) string {
	return ipfslib.Ipfs_add(os_path, second)
}

func IpfsGet(object_hash, os_path string, second int) string {
	return ipfslib.Ipfs_get(object_hash, os_path, second)
}

func IpfsPublish(object_hash string, second int) string {
	return ipfslib.Ipfs_publish(object_hash, second)
}

func IpfsRemotepin(peer_id, peer_key, object_hash string, second int) string {
	return ipfslib.Ipfs_remotepin(peer_id, peer_key, object_hash, second)
}

func IpfsConnectpeer(remote_peer string, second int) string {
	return ipfslib.Ipfs_connectpeer(remote_peer, second)
}
