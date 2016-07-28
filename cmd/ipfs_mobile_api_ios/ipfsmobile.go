// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hello is a trivial package for gomobile bind example.
package ipfsmobileios

import (
	ipfslib "github.com/ipfs/go-ipfs/cmd/ipfs_lib"
)

func IpfsCmd(path string, cmd string, second int) string {
	ipfslib.Ipfs_path(path)
	return ipfslib.Ipfs_cmd_arm(cmd, second)
}

func IpfsPath(path string) string {
	return ipfslib.Ipfs_path(path)
}

func IpfsInit(path string) string {
	return ipfslib.Ipfs_init(path)
}

func IpfsDaemon() string {
	return ipfslib.Ipfs_daemon()
}

func IpfsId(path string, second int) string {
	ipfslib.Ipfs_path(path)
	return ipfslib.Ipfs_id(second)
}

func IpfsPeerid(path string, new_id string, second int) string {
	ipfslib.Ipfs_path(path)
	return ipfslib.Ipfs_peerid(new_id, second)
}

func IpfsPrivkey(path string, new_key string, second int) string {
	ipfslib.Ipfs_path(path)
	return ipfslib.Ipfs_privkey(new_key, second)
}

func IpfsAdd(path string, os_path string, second int) string {
	ipfslib.Ipfs_path(path)
	return ipfslib.Ipfs_add(os_path, second)
}

func IpfsGet(path string, object_hash, os_path string, second int) string {
	ipfslib.Ipfs_path(path)
	return ipfslib.Ipfs_get(object_hash, os_path, second)
}

func IpfsPublish(path string, object_hash string, second int) string {
	ipfslib.Ipfs_path(path)
	return ipfslib.Ipfs_publish(object_hash, second)
}

func IpfsRemotepin(path, peer_id, peer_key, object_hash string, second int) string {
	ipfslib.Ipfs_path(path)
	return ipfslib.Ipfs_remotepin(peer_id, peer_key, object_hash, second)
}
