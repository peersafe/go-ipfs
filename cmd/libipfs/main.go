package main

import "C"
import "unsafe"
import "github.com/ipfs/go-ipfs/cmd/ipfs_lib"

//export ipfs_cmd
func ipfs_cmd(cmd string,out_res **C.char) int {
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
