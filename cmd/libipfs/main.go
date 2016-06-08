package main

/*
#include <string.h>
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unsafe"

	"github.com/ipfs/go-ipfs/cmd/ipfs_lib"
)

const hashLen int = 46
const keyLen int = 1596
const endsep = "\n"

const (
	errRet = -1
	sucRet = 0
)

type statInfo struct {
	Hash string
}

//export ipfs_init
func ipfs_init(out_res *C.char) int {
	cmd := "ipfs init -e"
	_, str, err := ipfs_lib.Ipfs_cmd(cmd)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	str = strings.Trim(str, endsep)
	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

//export ipfs_daemon
func ipfs_daemon(out_res *C.char) int {
	cmd := "ipfs daemon"
	_, str, err := ipfs_lib.Ipfs_cmd(cmd)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	str = strings.Trim(str, endsep)
	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

//export ipfs_id
func ipfs_id(second int, out_res *C.char) int {
	cmd := "ipfs id"
	_, str, err := ipfs_lib.Ipfs_cmd_time(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	str = strings.Trim(str, endsep)
	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

//export ipfs_add
func ipfs_add(root_hash, ipfs_path, os_path string, second int, out_res *C.char) int {
	if len(root_hash) != hashLen {
		fmt.Println("root_hash len not 46")
		return errRet
	}

	if len(ipfs_path) == 0 {
		fmt.Println("ipfs_path len is 0")
		return errRet
	}

	ipfs_path, err := ipfsPathClean(ipfs_path)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	var addHash string
	if len(os_path) != 0 {
		os_path, err := filepath.Abs(path.Clean(os_path))
		if err != nil {
			fmt.Println(err)
			return errRet
		}

		fi, err := os.Lstat(os_path)
		if err != nil {
			fmt.Println(err)
			return errRet
		}

		cmdSuff := ""
		if fi.Mode().IsDir() {
			cmdSuff = "ipfs add -r "
		} else if fi.Mode().IsRegular() {
			cmdSuff = "ipfs add "
		} else {
			return errRet
		}

		fmt.Println(cmdSuff, os_path)
		_, addHash, err = ipfs_lib.Ipfs_cmd_time(cmdSuff+os_path, second)
		if err != nil {
			return errRet
		}
	}

	ipfs_path = path.Clean(ipfs_path)
	cmd := "ipfs object patch " + root_hash + " add-link " + ipfs_path + " " + addHash
	fmt.Println(cmd)
	_, str, err := ipfs_lib.Ipfs_cmd_time(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	str = strings.Trim(str, endsep)
	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

//export ipfs_delete
func ipfs_delete(root_hash, ipfs_path string, second int, out_res *C.char) int {
	if len(root_hash) != hashLen {
		fmt.Println("root_hash len is not 46")
		return errRet
	}

	if len(ipfs_path) == 0 {
		fmt.Println("ipfs_path len is 0")
		return errRet
	}

	ipfs_path, err := ipfsPathClean(ipfs_path)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	cmd := "ipfs object patch " + root_hash + " rm-link " + ipfs_path
	fmt.Println(cmd)
	_, str, err := ipfs_lib.Ipfs_cmd_time(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	str = strings.Trim(str, endsep)
	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

//export ipfs_move
func ipfs_move(root_hash, ipfs_path_src, ipfs_path_des string, second int, out_res *C.char) int {
	if len(root_hash) != hashLen {
		fmt.Println("root_hash len is not 46")
		return errRet
	}

	if len(ipfs_path_src) == 0 {
		fmt.Println("ipfs_path_src len is 0")
		return errRet
	}

	ipfs_path_src, err := ipfsPathClean(ipfs_path_src)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	if len(ipfs_path_des) == 0 {
		fmt.Println("ipfs_path_des len is 0")
		return errRet
	}

	ipfs_path_des, err = ipfsPathClean(ipfs_path_des)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	object_path := ipfs_path_src
	if strings.HasPrefix(ipfs_path_src, "\"") && strings.HasSuffix(ipfs_path_src, "\"") {
		object_path = ipfs_path_src[1 : len(ipfs_path_src)-1]
	}

	statCmd := "ipfs object stat " + root_hash + "/" + object_path
	fmt.Println(statCmd)
	_, statStr, err := ipfs_lib.Ipfs_cmd_time(statCmd, second)
	if err != nil {
		return errRet
	}

	var nodeStat statInfo
	err = json.Unmarshal([]byte(statStr), &nodeStat)
	if err != nil {
		fmt.Println(err)
		return errRet
	}
	nodeStat.Hash = strings.Trim(nodeStat.Hash, endsep)

	addCmd := "ipfs object patch " + root_hash + " add-link " + ipfs_path_des + " " + nodeStat.Hash
	fmt.Println(addCmd)
	_, newHash, err := ipfs_lib.Ipfs_cmd_time(addCmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	newHash = strings.Trim(newHash, endsep)
	delCmd := "ipfs object patch " + newHash + " rm-link " + ipfs_path_src
	fmt.Println(delCmd)
	_, new_root_hash, err := ipfs_lib.Ipfs_cmd_time(delCmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	new_root_hash = strings.Trim(new_root_hash, endsep)
	cs := unsafe.Pointer(C.CString(new_root_hash))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(new_root_hash)))
	C.free(cs)
	return len(new_root_hash)
}

//export ipfs_shard
func ipfs_shard(object_hash, shard_name string, second int, out_res *C.char) int {
	if len(object_hash) != hashLen {
		fmt.Println("object_hash len is not 46")
		return errRet
	}

	if len(shard_name) == 0 {
		fmt.Println("shard_name len is 0")
		return errRet
	}

	shard_name, err := ipfsPathClean(shard_name)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	cmd := "ipfs object patch QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn add-link " + shard_name + " " + object_hash
	fmt.Println(cmd)
	_, str, err := ipfs_lib.Ipfs_cmd_time(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	str = strings.Trim(str, endsep)
	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

//export ipfs_get
func ipfs_get(shard_hash, os_path string, second int) int {
	if len(shard_hash) != hashLen {
		fmt.Println("shard_hash len is not 46")
		return errRet
	}
	if len(os_path) == 0 {
		fmt.Println("shard_name len is 0")
		return errRet
	}

	os_path, err := filepath.Abs(path.Clean(os_path))
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	cmd := "ipfs get " + shard_hash + " -o " + os_path
	fmt.Println(cmd)
	_, _, err = ipfs_lib.Ipfs_cmd_time(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}
	return sucRet
}

//export ipfs_query
func ipfs_query(object_hash, ipfs_path string, second int, out_res *C.char) int {
	if len(object_hash) != hashLen {
		fmt.Println("object_hash len is not 46")
		return errRet
	}

	if !strings.HasPrefix(ipfs_path, "/") {
		fmt.Println("ipfs_path must preffix is -")
		return errRet
	}
	ipfs_path = ipfs_path[1:]

	if len(ipfs_path) != 0 {
		statCmd := "ipfs object stat " + object_hash + "/" + ipfs_path
		fmt.Println(statCmd)
		_, statStr, err := ipfs_lib.Ipfs_cmd_time(statCmd, second)
		if err != nil {
			fmt.Println(err)
			return errRet
		}

		var nodeStat statInfo
		err = json.Unmarshal([]byte(statStr), &nodeStat)
		if err != nil {
			fmt.Println(err)
			return errRet
		}
		object_hash = strings.Trim(nodeStat.Hash, endsep)
	}

	cmd := "ipfs ls " + object_hash
	fmt.Println(cmd)
	_, str, err := ipfs_lib.Ipfs_cmd_time(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	str = strings.Trim(str, endsep)
	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

//export ipfs_merge
func ipfs_merge(root_hash, ipfs_path, shard_hash string, second int, out_res *C.char) int {
	if len(root_hash) != hashLen {
		fmt.Println("root_hash len is not 46")
		return errRet
	}

	if len(shard_hash) != hashLen {
		fmt.Println("shard_hash len is not 46")
		return errRet
	}

	if len(ipfs_path) == 0 {
		fmt.Println("ipfs_path len is 0")
		return errRet
	}

	ipfs_path, err := ipfsPathClean(ipfs_path)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	cmd := "ipfs object patch " + root_hash + " add-link " + ipfs_path + " " + shard_hash
	fmt.Println(cmd)
	_, str, err := ipfs_lib.Ipfs_cmd_time(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	str = strings.Trim(str, endsep)
	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

//export ipfs_peerid
func ipfs_peerid(new_id string, second int, out_res *C.char) int {
	if len(new_id) != hashLen && len(new_id) != 0 {
		fmt.Println("new_id len is not 46 or is not 0")
		return errRet
	}

	cmd := "ipfs config Identity.PeerID"
	fmt.Println(cmd)
	_, peeId, err := ipfs_lib.Ipfs_cmd_time(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	if len(new_id) == hashLen {
		cmd += " " + new_id
		fmt.Println(cmd)
		_, _, err := ipfs_lib.Ipfs_cmd_time(cmd, second)
		if err != nil {
			fmt.Println(err)
			return errRet
		}
		peeId = new_id
	}

	peeId = strings.Trim(peeId, endsep)
	cs := unsafe.Pointer(C.CString(peeId))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(peeId)))
	C.free(cs)
	return len(peeId)
}

//export ipfs_privkey
func ipfs_privkey(new_key string, second int, out_res *C.char) int {
	if len(new_key) != keyLen && len(new_key) != 0 {
		fmt.Println("new_id len is not 1596 or is not 0")
		return errRet
	}

	cmd := "ipfs config Identity.PrivKey"
	fmt.Println(cmd)
	_, key, err := ipfs_lib.Ipfs_cmd_time(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	if len(new_key) == hashLen {
		cmd += " " + new_key
		fmt.Println(cmd)
		_, _, err := ipfs_lib.Ipfs_cmd_time(cmd, second)
		if err != nil {
			fmt.Println(err)
			return errRet
		}
		key = new_key
	}

	key = strings.Trim(key, endsep)
	cs := unsafe.Pointer(C.CString(key))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(key)))
	C.free(cs)
	return len(key)
}

//export ipfs_publish
func ipfs_publish(object_hash string, second int, out_res *C.char) int {
	if len(object_hash) != hashLen {
		fmt.Println("object_hash len is not 46")
		return errRet
	}

	cmd := "ipfs name publish /ipfs/" + object_hash
	fmt.Println(cmd)
	_, hash, err := ipfs_lib.Ipfs_cmd_time(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	hash = strings.Trim(hash, endsep)
	cs := unsafe.Pointer(C.CString(hash))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(hash)))
	C.free(cs)
	return len(hash)
}

//export ipfs_config
func ipfs_config(key, value string, out_res *C.char) int {
	var cmd string
	if len(key) == 0 {
		cmd = "ipfs config show"
	} else if len(key) != 0 && len(value) == 0 {
		cmd = "ipfs config " + key
	} else {
		cmd = "ipfs config " + key + " " + value
	}
	fmt.Println(cmd)

	_, str, err := ipfs_lib.Ipfs_cmd(cmd)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	str = strings.Trim(str, endsep)
	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

//export ipfs_cmd
func ipfs_cmd(cmd string, second int, out_res *C.char) int {
	_, str, err := ipfs_lib.Ipfs_cmd_time(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	str = strings.Trim(str, endsep)
	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

func ipfsPathClean(ipfsPath string) (string, error) {
	if !strings.HasPrefix(ipfsPath, "/") {
		return "", errors.New("must prefix is /")
	}

	path := ipfsPath[1:]
	if strings.HasPrefix(path, "-") {
		path = "\"" + path + "\""
	}
	return path, nil
}

// main roadmap:
// - parse the commandline to get a cmdInvocation
// - if user requests, help, print it and exit.
// - run the command invocation
// - output the response
// - if anything fails, print error, maybe with help
func main() {
}
