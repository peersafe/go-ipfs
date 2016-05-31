package main

import "C"
import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"unsafe"

	"github.com/ipfs/go-ipfs/cmd/ipfs_lib"
)

//export ipfs_init
func ipfs_init(out_res **C.char) int {
	cmd := "ipfs init -e"
	_, str, err := ipfs_lib.Ipfs_cmd(cmd)
	fmt.Println("[[", str, "]]")

	cs := unsafe.Pointer(C.CString(str))
	*out_res = (*C.char)(cs)
	if err == nil {
		return len(str)
	} else {
		return -1
	}
}

//export ipfs_daemon
func ipfs_daemon(out_res **C.char) int {
	cmd := "ipfs daemon"
	_, str, err := ipfs_lib.Ipfs_cmd(cmd)
	fmt.Println("[[", str, "]]")

	cs := unsafe.Pointer(C.CString(str))
	*out_res = (*C.char)(cs)
	if err == nil {
		return len(str)
	} else {
		return -1
	}
}

//export ipfs_add
func ipfs_add(root_hash, ipfs_path, os_path string, out_res **C.char) int {
	if len(root_hash) != 46 {
		out_res = nil
		fmt.Println("error 1")
		return -1
	}

	if len(ipfs_path) == 0 {
		out_res = nil
		fmt.Println("error 2")
		return -1
	}

	var err error
	var addHash string
	if len(os_path) != 0 {
		os_path, err := filepath.Abs(path.Clean(os_path))
		if err != nil {
			out_res = nil
			return -1
		}

		fi, err := os.Lstat(os_path)
		cmdSuff := ""
		if fi.Mode().IsDir() {
			cmdSuff = "ipfs add -r "
		} else if fi.Mode().IsRegular() {
			cmdSuff = "ipfs add "
		} else {
			out_res = nil
			return -1
		}
		fmt.Println("add cmd", cmdSuff, os_path)
		_, addHash, err = ipfs_lib.Ipfs_cmd(cmdSuff + os_path)
		if err != nil {
			out_res = nil
			return -1
		}
	}

	ipfs_path = path.Clean(ipfs_path)
	cmd := "ipfs object patch " + root_hash + " add-link " + ipfs_path + " " + addHash
	fmt.Println("object cmd", cmd)
	_, str, err := ipfs_lib.Ipfs_cmd(cmd)
	if err != nil {
		out_res = nil
		return -1
	}

	cs := unsafe.Pointer(C.CString(str))
	*out_res = (*C.char)(cs)
	if err == nil {
		return len(str)
	} else {
		return -1
	}
}

//export ipfs_cmd
func ipfs_cmd(cmd string, out_res **C.char) int {
	_, str, err := ipfs_lib.Ipfs_cmd(cmd)

	cs := unsafe.Pointer(C.CString(str))
	*out_res = (*C.char)(cs)
	if err == nil {
		return len(str)
	} else {
		return -1
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
