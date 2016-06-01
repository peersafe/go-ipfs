package ipfs_lib

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/mitchellh/go-homedir/homedir"
	"github.com/ipfs/go-ipfs/cmd/ipfs_lib"
)

const separtor = "&X&"
const (
	errRet = 1
	sucRet = 0
)

func Ipfs_cmd_arm(cmd string) string {
	res, str, _ := Ipfs_cmd(cmd)

	return string(res) + separtor + str
}

func Ipfs_path(path string) string {
	homedir.Home_Unix_Dir = path
}

func Ipfs_init(path string) string {
	cmd := "ipfs init -e"
	homedir.Home_Unix_Dir = path
	res, str, _ := Ipfs_cmd(cmd)
	return string(res) + separtor + str
}

func Ipfs_daemon() string {
	cmd := "ipfs daemon"
	res, str, _ := Ipfs_cmd(cmd)
	return string(res) + spartor + str
}

func Ipfs_id() string {
	cmd := "ipfs id"
	res, str, _ := Ipfs_cmd(cmd)
	return string(res) + spartor + str
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
		return string(res) + separtor + addHash

	} else {
		return errRet + separtor + ""
	}
}

func Ipfs_get(object_hash, os_path string) string {
	if len(shard_hash) != 46 {
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
