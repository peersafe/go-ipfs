package ipfs_lib

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
)

const separtor = "&X&"

func Ipfs_cmd_arm(cmd string, second int) string {
	res, str := IpfsCmdApi(cmd, second)
	return fmt.Sprintf("%d%s%s", res, separtor, str)
}

func Ipfs_path(path string) string {
	homedir.Home_Unix_Dir = path
	return fmt.Sprintf("%d%s%s", sucRet, separtor, "")
}

func Ipfs_init(path string) string {
	homedir.Home_Unix_Dir = path
	res, str := IpfsInit()
	return fmt.Sprintf("%d%s%s", res, separtor, str)
}

func Ipfs_daemon() string {
	res, str := IpfsDaemon()
	return fmt.Sprintf("%d%s%s", res, separtor, str)
}

func Ipfs_config(key, value string) string {
	res, str := IpfsConfig(key, value)
	return fmt.Sprintf("%d%s%s", res, separtor, str)
}

func Ipfs_id(second int) string {
	res, str := IpfsId(second)
	return fmt.Sprintf("%d%s%s", res, separtor, str)
}

func Ipfs_peerid(new_id string, second int) string {
	res, str := IpfsPeerid(new_id, second)
	return fmt.Sprintf("%d%s%s", res, separtor, str)
}

func Ipfs_privkey(new_key string, second int) string {
	res, str := IpfsPrivkey(new_key, second)
	return fmt.Sprintf("%d%s%s", res, separtor, str)
}

func Ipfs_add(os_path string, second int) string {
	if len(os_path) != 0 {
		os_path, err := filepath.Abs(path.Clean(os_path))
		if err != nil {
			return fmt.Sprintf("%d%s%s", errRet, separtor, "")
		}

		fi, err := os.Lstat(os_path)
		cmdSuff := ""
		if fi.Mode().IsDir() {
			cmdSuff = "ipfs add -r "
		} else if fi.Mode().IsRegular() {
			cmdSuff = "ipfs add "
		} else {
			return fmt.Sprintf("%d%s%s", errRet, separtor, "")
		}
		fmt.Println("add cmd", cmdSuff, os_path)
		res, addHash, err := ipfsCmdTime(cmdSuff+os_path, second)
		if err != nil {
			return fmt.Sprintf("%d%s%s", res, separtor, "")
		}
		addHash = strings.Trim(addHash, endsep)
		return fmt.Sprintf("%d%s%s", res, separtor, addHash)

	} else {
		return fmt.Sprintf("%d%s%s", errRet, separtor, "")
	}
}

func Ipfs_get(shard_hash, os_path string, second int) string {
	res := IpfsGet(shard_hash, os_path, second)
	return fmt.Sprintf("%d%s%s", res, separtor, "")
}

func Ipfs_publish(object_hash string, second int) string {
	res, str := IpfsPublish(object_hash, second)
	return fmt.Sprintf("%d%s%s", res, separtor, str)
}

func Ipfs_remotepin(peer_id, object_hash string, second int) string {
	res, str := IpfsRemotepin(peer_id, object_hash, second)
	return fmt.Sprintf("%d%s%s", res, separtor, str)
}
