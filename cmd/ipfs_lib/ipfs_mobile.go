package ipfs_lib

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
)

const preLen int = 6
const hashLen int = 46
const keyLen int = 1596
const separtor = "&X&"
const endsep = "\n"
const (
	errRet = 1
	sucRet = 0
)

func Ipfs_cmd_arm(cmd string, second int) string {
	res, str, _ := Ipfs_cmd_time(cmd, second)

	str = strings.Trim(str, endsep)
	return fmt.Sprintf("%d%s%s", res, separtor, str)
}

func Ipfs_path(path string) string {
	homedir.Home_Unix_Dir = path
	return fmt.Sprintf("%d%s%s", sucRet, separtor, "")
}

func Ipfs_init(path string) string {
	cmd := "ipfs init -e"
	homedir.Home_Unix_Dir = path
	res, str, _ := Ipfs_cmd(cmd)
	str = strings.Trim(str, endsep)
	return fmt.Sprintf("%d%s%s", res, separtor, str)
}

func Ipfs_daemon() string {
	cmd := "ipfs daemon"
	res, str, _ := Ipfs_cmd(cmd)
	str = strings.Trim(str, endsep)
	return fmt.Sprintf("%d%s%s", res, separtor, str)
}

func Ipfs_id() string {
	cmd := "ipfs id"
	res, str, _ := Ipfs_cmd(cmd)
	str = strings.Trim(str, endsep)
	return fmt.Sprintf("%d%s%s", res, separtor, str)
}

func Ipfs_peerid(new_id string) string {
	if len(new_id) != hashLen && len(new_id) != 0 {
		return fmt.Sprintf("%d%s%s", errRet, separtor, "")
	}

	cmd := "ipfs config Identity.PeerID"
	_, peerId, err := Ipfs_cmd(cmd)
	if err != nil {
		return fmt.Sprintf("%d%s%s", errRet, separtor, "")
	}

	if len(new_id) == hashLen {
		cmd += " " + new_id
		_, _, err := Ipfs_cmd(cmd)
		if err != nil {
			return fmt.Sprintf("%d%s%s", errRet, separtor, "")
		}
		peerId = new_id
	}
	peerId = strings.Trim(peerId, endsep)
	return fmt.Sprintf("%d%s%s", sucRet, separtor, peerId)
}

func Ipfs_privkey(new_key string) string {
	if len(new_key) != keyLen && len(new_key) != 0 {
		return fmt.Sprintf("%d%s%s", errRet, separtor, "")
	}

	cmd := "ipfs config Identity.PrivKey"
	_, key, err := Ipfs_cmd(cmd)
	if err != nil {
		return fmt.Sprintf("%d%s%s", errRet, separtor, "")
	}

	if len(new_key) == hashLen {
		cmd += " " + new_key
		_, _, err := Ipfs_cmd(cmd)
		if err != nil {
			return fmt.Sprintf("%d%s%s", errRet, separtor, "")
		}
		key = new_key
	}
	key = strings.Trim(key, endsep)
	return fmt.Sprintf("%d%s%s", sucRet, separtor, key)
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
		res, addHash, err := Ipfs_cmd_time(cmdSuff+os_path, second)
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
	if len(shard_hash) == 0 {
		return fmt.Sprintf("%d%s%s", errRet, separtor, "")
	}
	if len(os_path) == 0 {
		return fmt.Sprintf("%d%s%s", errRet, separtor, "")
	}

	os_path, _ = filepath.Abs(path.Clean(os_path))

	cmd := "ipfs get " + shard_hash + " -o " + os_path
	fmt.Println("get cmd:", cmd)
	_, _, err := Ipfs_cmd_time(cmd, second)
	if err != nil {
		return fmt.Sprintf("%d%s%s", errRet, separtor, "")
	}
	return fmt.Sprintf("%d%s%s", sucRet, separtor, "")
}

func Ipfs_publish(object_hash string, second int) string {
	if len(object_hash) != hashLen+preLen {
		return fmt.Sprintf("%d%s%s", errRet, separtor, "")
	}

	cmd := "ipfs name publish " + object_hash
	fmt.Println(cmd)
	res, hash, err := Ipfs_cmd_time(cmd, second)
	if err != nil {
		return fmt.Sprintf("%d%s%s", errRet, separtor, "")
	}
	hash = strings.Trim(hash, endsep)
	return fmt.Sprintf("%d%s%s", res, separtor, hash)
}
