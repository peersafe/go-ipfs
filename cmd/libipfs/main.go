package main

/*
#include <string.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"unsafe"

	"github.com/ipfs/go-ipfs/cmd/ipfs_lib"
)

const (
	errRet = -1
	sucRet = 0
)

//export ipfs_init
func ipfs_init(out_res *C.char) int {
	cmd := "ipfs init -e"
	_, str, err := ipfs_lib.Ipfs_cmd(cmd)
	fmt.Println("[[", str, "]]")
	if err != nil {
		return errRet
	}

	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

//export ipfs_daemon
func ipfs_daemon(out_res *C.char) int {
	cmd := "ipfs daemon"
	_, str, err := ipfs_lib.Ipfs_cmd(cmd)
	fmt.Println("[[", str, "]]")
	if err != nil {
		return errRet
	}

	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

//export ipfs_add
func ipfs_add(root_hash, ipfs_path, os_path string, out_res *C.char) int {
	if len(root_hash) != 46 {
		fmt.Println("error 1")
		return errRet
	}

	if len(ipfs_path) == 0 {
		fmt.Println("error 2")
		return errRet
	}

	var err error
	var addHash string
	if len(os_path) != 0 {
		os_path, err := filepath.Abs(path.Clean(os_path))
		if err != nil {
			return errRet
		}

		fi, err := os.Lstat(os_path)
		cmdSuff := ""
		if fi.Mode().IsDir() {
			cmdSuff = "ipfs add -r "
		} else if fi.Mode().IsRegular() {
			cmdSuff = "ipfs add "
		} else {
			return errRet
		}
		fmt.Println("add cmd", cmdSuff, os_path)
		_, addHash, err = ipfs_lib.Ipfs_cmd(cmdSuff + os_path)
		if err != nil {
			return errRet
		}
	}

	ipfs_path = path.Clean(ipfs_path)
	cmd := "ipfs object patch " + root_hash + " add-link " + ipfs_path + " " + addHash
	fmt.Println("object cmd", cmd)
	_, str, err := ipfs_lib.Ipfs_cmd(cmd)
	if err != nil {
		return errRet
	}

	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

//export ipfs_shard
func ipfs_shard(object_hash, shard_name string, out_res *C.char) int {
	if len(object_hash) != 46 {
		return errRet
	}

	if len(shard_name) == 0 {
		return errRet
	}

	cmd := "ipfs object patch QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn add-link " + shard_name + " " + object_hash
	fmt.Println("object cmd", cmd)
	_, str, err := ipfs_lib.Ipfs_cmd(cmd)
	if err != nil {
		return errRet
	}

	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	return len(str)
}

//export ipfs_get
func ipfs_get(shard_hash, os_path string) int {
	if len(shard_hash) != 46 {
		return errRet
	}
	if len(os_path) == 0 {
		return errRet
	}

	os_path, _ = filepath.Abs(path.Clean(os_path))

	cmd := "ipfs get " + shard_hash + " -o " + os_path
	fmt.Println("get cmd:", cmd)
	_, _, err := ipfs_lib.Ipfs_cmd(cmd)
	if err != nil {
		return errRet
	}
	return sucRet
}

//export ipfs_cmd
func ipfs_cmd(cmd string, out_res *C.char) int {
	_, str, err := ipfs_lib.Ipfs_cmd(cmd)

	cs := unsafe.Pointer(C.CString(str))
	C.memcpy(unsafe.Pointer(out_res), cs, C.size_t(len(str)))
	C.free(cs)
	if err == nil {
		return len(str)
	} else {
		return errRet
	}
}

// main roadmap:
// - parse the commandline to get a cmdInvocation
// - if user requests, help, print it and exit.
// - run the command invocation
// - output the response
// - if anything fails, print error, maybe with help
func main() {
}
