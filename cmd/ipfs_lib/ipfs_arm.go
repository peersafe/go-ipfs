package ipfs_lib

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/mitchellh/go-homedir/homedir"
	"github.com/ipfs/go-ipfs/cmd/ipfs_lib"
)

const hashLen int = 46
const keylen int = 1596
const separtor = "&X&"
const endsep = "\n"
const (
	errRet = 1
	sucRet = 0
)

func Ipfs_cmd_arm(cmd string) string {
	res, str, _ := Ipfs_cmd(cmd)

	str = strings.Trim(str, endsep)
	return string(res) + separtor + str
}

func Ipfs_path(path string) string {
	homedir.Home_Unix_Dir = path
}

func Ipfs_init(path string) string {
	cmd := "ipfs init -e"
	homedir.Home_Unix_Dir = path
	res, str, _ := Ipfs_cmd(cmd)
	str = strings.Trim(str, endsep)
	return string(res) + separtor + str
}

func Ipfs_daemon() string {
	cmd := "ipfs daemon"
	res, str, _ := Ipfs_cmd(cmd)
	str = strings.Trim(str, endsep)
	return string(res) + spartor + str
}

func Ipfs_id() string {
	cmd := "ipfs id"
	res, str, _ := Ipfs_cmd(cmd)
	str = strings.Trim(str, endsep)
	return string(res) + spartor + str
}

func Ipfs_peerid(new_id string) string {
	if len(new_id) != hashLen && len(new_id) != 0 {
		return errRet + separtor + ""
	}

	cmd := "ipfs config Identity.PeerID"
	_, peerId, err := Ipfs_cmd(cmd)
	if err != nil {
		return errRet + separtor + ""
	}

	if len(new_id) == hashLen {
		cmd += " " + new_id
		_, _, err := Ipfs_cmd(cmd)
		if err != nil {
			return errRet + separtor + ""
		}
		peerId = new_id
	}
	peerId = strings.Trim(peerId, endsep)
	return sucRet + spartor + peerId
}

func Ipfs_privkey(new_key string) string {
	if len(new_key) != keyLen && len(new_key) != 0 {
		return errRet + separtor + ""
	}

	cmd := "ipfs config Identity.PrivKey"
	_, key, err := Ipfs_cmd(cmd)
	if err != nil {
		return errRet + separtor + ""
	}

	if len(new_key) == hashLen {
		cmd += " " + new_key
		_, _, err := Ipfs_cmd(cmd)
		if err != nil {
			return errRet + separtor + ""
		}
		key = new_key
	}
	key = strings.Trim(key, endsep)
	return sucRet + spartor + key
}

func Ipfs_add(os_path string) string {
	var err error
	var addHash string
	if len(os_path) != 0 {
		os_path, err := filepath.Abs(path.Clean(os_path))
		if err != nil {
			return errRet + separtor + ""
		}

		fi, err := os.Lstat(os_path)
		cmdSuff := ""
		if fi.Mode().IsDir() {
			cmdSuff = "ipfs add -r "
		} else if fi.Mode().IsRegular() {
			cmdSuff = "ipfs add "
		} else {
			return errRet + separtor + ""
		}
		fmt.Println("add cmd", cmdSuff, os_path)
		res, addHash, _ = ipfs_lib.Ipfs_cmd(cmdSuff + os_path)
		if err != nil {
			return string(res) + separtor + ""
		}
		addHash = strings.Trim(addHash, endsep)
		return string(res) + separtor + addHash

	} else {
		return errRet + separtor + ""
	}
}

func Ipfs_get(object_hash, os_path string) string {
	if len(shard_hash) != hashLen {
		return errRet + separtor + ""
	}
	if len(os_path) == 0 {
		return errRet + separtor + ""
	}

	os_path, _ = filepath.Abs(path.Clean(os_path))

	cmd := "ipfs get " + shard_hash + " -o " + os_path
	fmt.Println("get cmd:", cmd)
	_, _, err := ipfs_lib.Ipfs_cmd(cmd)
	if err != nil {
		return errRet + separtor + ""
	}
	return sucRet + separtor + ""
}
