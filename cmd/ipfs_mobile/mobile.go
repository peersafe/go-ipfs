// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hello is a trivial package for gomobile bind example.
package ipfsmobile

import (
	ipfslib "github.com/ipfs/go-ipfs/cmd/ipfs_lib"
	"github.com/ipfs/go-ipfs/commands"
)

type CmdCall interface {
	commands.CallFunc
}

func IpfsCmd(cmd string, second int, call CmdCall) string {
	return ipfslib.Ipfs_async_cmd_arm(cmd, second, call)
}

func IpfsPath(path string) string {
	return ipfslib.Ipfs_async_path(path)
}

func IpfsInit(call CmdCall) string {
	return ipfslib.Ipfs_async_init(call)
}

func IpfsDaemon(call CmdCall) string {
	return ipfslib.Ipfs_async_daemon(call)
}

func IpfsShutdown(call CmdCall) string {
	return ipfslib.Ipfs_async_shutdown(call)
}

func IpfsId(second int, call CmdCall) string {
	return ipfslib.Ipfs_async_id(second, call)
}

func IpfsPeerid(new_id string, second int, call CmdCall) string {
	return ipfslib.Ipfs_async_peerid(new_id, second, call)
}

func IpfsPrivkey(new_key string, second int, call CmdCall) string {
	return ipfslib.Ipfs_async_privkey(new_key, second, call)
}

func IpfsAdd(os_path string, second int, call CmdCall) string {
	return ipfslib.Ipfs_async_add(os_path, second, call)
}

func IpfsGet(object_hash, os_path string, second int, call CmdCall) string {
	return ipfslib.Ipfs_async_get(object_hash, os_path, second, call)
}

func IpfsPublish(object_hash string, second int, call CmdCall) string {
	return ipfslib.Ipfs_async_publish(object_hash, second, call)
}

func IpfsConfig(key, value string, call CmdCall) string {
	return ipfslib.Ipfs_async_config(key, value, call)
}

func IpfsRemotepin(peer_id, peer_key, object_hash string, second int, call CmdCall) string {
	return ipfslib.Ipfs_async_remotepin(peer_id, peer_key, object_hash, second, call)
}

func IpfsRemotels(peer_id, peer_key, object_hash string, second int, call CmdCall) string {
	return ipfslib.Ipfs_async_remotels(peer_id, peer_key, object_hash, second, call)
}

func IpfsConnectpeer(remote_peer string, second int, call CmdCall) string {
	return ipfslib.Ipfs_async_connectpeer(remote_peer, second, call)
}
