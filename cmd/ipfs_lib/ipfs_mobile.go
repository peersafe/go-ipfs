package ipfs_lib

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
)

func Ipfs_cmd_arm(cmd string, second int) string {
	res, str := IpfsCmdApi(cmd, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_path(path string) string {
	homedir.Home_Unix_Dir = path
	return fmt.Sprintf("%d%s%s", sucRet, cmdSep, "")
}

func Ipfs_init(path string) string {
	homedir.Home_Unix_Dir = path
	res, str := IpfsInit()
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_daemon() string {
	res, str := IpfsDaemon()
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_config(key, value string) string {
	res, str := IpfsConfig(key, value)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_id(second int) string {
	res, str := IpfsId(second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_peerid(new_id string, second int) string {
	res, str := IpfsPeerid(new_id, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_privkey(new_key string, second int) string {
	res, str := IpfsPrivkey(new_key, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_add(os_path string, second int) string {
	if len(os_path) != 0 {
		os_path, err := filepath.Abs(path.Clean(os_path))
		if err != nil {
			return fmt.Sprintf("%d%s%s", errRet, cmdSep, "")
		}

		fi, err := os.Lstat(os_path)
		cmdSuff := ""
		if fi.Mode().IsDir() {
			cmdSuff = strings.Join([]string{"ipfs", "add", "--is-lib=true", "-r", os_path}, cmdSep)
		} else if fi.Mode().IsRegular() {
			cmdSuff = strings.Join([]string{"ipfs", "add", "--is-lib=true", os_path}, cmdSep)
		} else {
			return mt.Sprintf("%d%s%s", errRet, cmdSep, "")
		}
		res, addHash, err := ipfsCmdTime(cmdSuff, second)
		if err != nil {
			return fmt.Sprintf("%d%s%s", res, cmdSep, "")
		}
		addHash = strings.Trim(addHash, endsep)
		return fmt.Sprintf("%d%s%s", res, cmdSep, addHash)

	} else {
		return fmt.Sprintf("%d%s%s", errRet, cmdSep, "")
	}
}

func Ipfs_get(shard_hash, os_path string, second int) string {
	res := IpfsGet(shard_hash, os_path, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, "")
}

func Ipfs_publish(object_hash string, second int) string {
	res, str := IpfsPublish(object_hash, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_remotepin(peer_id, peer_key, object_hash string, second int) string {
	res, str := IpfsRemotepin(peer_id, peer_key, object_hash, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_remotels(peer_id, peer_key, object_hash string, second int) string {
	res, str := IpfsRemotels(peer_id, peer_key, object_hash, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_connectpeer(remote_peer string, second int) string {
	res, str := IpfsConnectPeer(remote_peer, second)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}
