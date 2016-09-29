package ipfs_lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	homedir "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
	"github.com/ipfs/go-ipfs/cmd/ipfs_lib/apiinterface"
	"github.com/ipfs/go-ipfs/commands"
)

var asyncApiIns Instance

type apiAsyncCmd struct {
}

func (a *apiAsyncCmd) Cmd(str string, sec int) (int, string, error) {
	return ipfsAsyncCmdTime(str, sec, func(result string, err error) {
		if err != nil {
			fmt.Printf("apiAsyncCmd cmd err=[%v]\n", err)
		}
		fmt.Printf("apiAsyncCmd cmd result=[%v]\n")
	})
}

func IpfsAsyncInit(path string) (string, error) {
	if path != "" {
		asyncApiIns = NewInstance(path)
		homedir.Home_Unix_Dir = path
	} else {
		return "", errors.New("path is nil!")
	}

	cmd := strings.Join([]string{"ipfs", "init", "-e"}, cmdSep)
	fmt.Println(cmd)
	call := func(result string, err error) {}
	_, str, err := ipfsAsyncCmd(cmd, call)
	if err != nil {
		return "", err
	}
	return str, nil
}

func IpfsAsyncDaemon(path string, outerCall commands.RequestCB) {
	if path == "" {
		outerCall("", errors.New("path is nil!"))
		return
	}
	if path != "" && asyncApiIns == nil {
		asyncApiIns = NewInstance(path)
		homedir.Home_Unix_Dir = path
	}

	// init apiinterface for remote cmds
	if apiinterface.GApiInterface == nil {
		apiinterface.GApiInterface = new(apiAsyncCmd)
	}

	cmd := strings.Join([]string{"ipfs", "daemon"}, cmdSep)
	call := func(result string, err error) {
		outerCall(result, err)
	}
	_, _, err := ipfsAsyncCmd(cmd, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncShutDown(outerCall commands.RequestCB) {
	cmd := strings.Join([]string{"ipfs", "shutdown"}, cmdSep)
	call := func(result string, err error) {
		outerCall(result, err)
	}
	_, _, err := ipfsAsyncCmd(cmd, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncId(second int, outerCall commands.RequestCB) {
	cmd := strings.Join([]string{"ipfs", "id"}, cmdSep)
	call := func(result string, err error) {
		outerCall(result, err)
	}
	_, _, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncAdd(os_path string, second int, outerCall commands.RequestCB, cancel chan struct{}) {
	call := func(add_hash string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}

		add_hash = stringTrim(add_hash)
		// do progress callback
		if !strings.Contains(add_hash, "Over") {
			outerCall(add_hash, err)
			return
		}

		// get real add_hash
		results := strings.Split(add_hash, cmdSep)
		add_hash = results[1]
		outerCall(add_hash, nil)
	}

	var err error
	os_path, err = filepath.Abs(path.Clean(os_path))
	if err != nil {
		outerCall("", err)
		return
	}

	fi, err := os.Lstat(os_path)
	if err != nil {
		outerCall("", err)
		return
	}

	cmd := ""
	if fi.Mode().IsDir() {
		cmd = strings.Join([]string{"ipfs", "add", "--is-lib=true", "-r", os_path}, cmdSep)
	} else if fi.Mode().IsRegular() {
		cmd = strings.Join([]string{"ipfs", "add", "--is-lib=true", os_path}, cmdSep)
	} else {
		outerCall("", errors.New("Unkown file type!"))
		return
	}

	_, _, err = ipfsAsyncCmdWithCancel(cmd, second, call, cancel)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncDelete(root_hash, ipfs_path string, second int, outerCall commands.RequestCB) {
	var err error
	if root_hash, err = ipfsObjectHashCheck(root_hash); err != nil {
		outerCall("", errors.New("root_hash len not 46"))
		return
	}

	if len(ipfs_path) == 0 {
		outerCall("", errors.New("ipfs_path len is 0"))
		return
	}

	ipfs_path, err = ipfsPathClean(ipfs_path)
	if err != nil {
		outerCall("", err)
		return
	}

	call := func(result string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		new_root := stringTrim(result)
		outerCall(new_root, nil)
	}
	cmd := strings.Join([]string{"ipfs", "object", "patch", "rm-link", root_hash, ipfs_path}, cmdSep)
	_, _, err = ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncMove(root_hash, ipfs_path_src, ipfs_path_des string, second int, outerCall commands.RequestCB) {
	var err error
	if root_hash, err = ipfsObjectHashCheck(root_hash); err != nil {
		outerCall("", errors.New("root_hash len not 46"))
		return
	}

	if len(ipfs_path_src) == 0 {
		outerCall("", errors.New("ipfs_path_src len is 0"))
		return
	}

	ipfs_path_src, err = ipfsPathClean(ipfs_path_src)
	if err != nil {
		outerCall("", err)
		return
	}

	if len(ipfs_path_des) == 0 {
		outerCall("", errors.New("ipfs_path_des len is 0"))
		return
	}

	ipfs_path_des, err = ipfsPathClean(ipfs_path_des)
	if err != nil {
		outerCall("", err)
		return
	}

	object_path := ipfs_path_src
	if strings.HasPrefix(ipfs_path_src, "\"") && strings.HasSuffix(ipfs_path_src, "\"") {
		object_path = ipfs_path_src[1 : len(ipfs_path_src)-1]
	}

	call := func(info string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		info = stringTrim(info)

		var nodeStat statInfo
		err = json.Unmarshal([]byte(info), &nodeStat)
		if err != nil {
			outerCall("", err)
			return
		}
		nodeStat.Hash = strings.Trim(nodeStat.Hash, endsep)

		call2 := func(patch string, err error) {
			if err != nil {
				outerCall("", err)
				return
			}

			call3 := func(result string, err error) {
				if err != nil {
					outerCall("", err)
					return
				}
				new_root := stringTrim(result)
				outerCall(new_root, nil)
			}
			newHash := strings.Trim(patch, endsep)
			delCmd := strings.Join([]string{"ipfs", "object", "patch", "rm-link", newHash, ipfs_path_src}, cmdSep)
			_, _, err = ipfsAsyncCmdTime(delCmd, second, call3)
			if err != nil {
				outerCall("", err)
				return
			}
		}
		addCmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", root_hash, ipfs_path_des, nodeStat.Hash}, cmdSep)
		_, _, err = ipfsAsyncCmdTime(addCmd, second, call2)
		if err != nil {
			outerCall("", err)
			return
		}
	}
	statCmd := strings.Join([]string{"ipfs", "object", "stat", "--is-lib=true", root_hash + "/" + object_path}, cmdSep)
	_, _, err = ipfsAsyncCmdTime(statCmd, second, call)
	if err != nil {
		outerCall("", err)
		return
	}
}

func IpfsAsyncShare(object_hash, share_name string, second int, outerCall commands.RequestCB) {
	var err error
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		outerCall("", errors.New("object_hash len not 46"))
		return
	}

	if len(share_name) == 0 {
		outerCall("", errors.New("share_name len is 0"))
		return
	}

	share_name, err = ipfsPathClean(share_name)
	if err != nil {
		outerCall("", err)
		return
	}

	call := func(result string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		new_hash := stringTrim(result)
		outerCall(new_hash, nil)
	}
	cmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn", share_name, object_hash}, cmdSep)
	_, _, err = ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncGet(share_hash, os_path string, second int, outerCall commands.RequestCB, cancel chan struct{}) {
	var err error
	if share_hash, err = ipfsHashCheck(share_hash); err != nil {
		outerCall("", errors.New("share_hash format error"))
		return
	}

	call := func(result string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		result = stringTrim(result)
		// do progress callback
		if result != "" && !strings.Contains(result, "Over") {
			outerCall(result, nil)
			return
		}
		outerCall(result, nil)
	}

	var cmd string
	if len(os_path) != 0 {
		os_path, err = filepath.Abs(path.Clean(os_path))
		if err != nil {
			outerCall("", err)
			return
		}
		cmd = strings.Join([]string{"ipfs", "get", share_hash, "-o", os_path}, cmdSep)
	} else { // remtote cmd pin add
		cmd = strings.Join([]string{"ipfs", "get", share_hash, "-o", "/dev/null"}, cmdSep)
	}

	_, _, err = ipfsAsyncCmdWithCancel(cmd, second, call, cancel)
	if err != nil {
		outerCall("", err)
		return
	}
}

func IpfsAsyncQuery(object_hash, ipfs_path string, second int, outerCall commands.RequestCB) {
	fmt.Printf("ipfs_lib IpfsAsyncQuery [%s] [%s]\n", object_hash, ipfs_path)
	var err error
	if object_hash, err = ipfsHashCheck(object_hash); err != nil {
		outerCall("", errors.New("object_hash len not 46"))
		return
	}

	if !strings.HasPrefix(ipfs_path, "/") {
		outerCall("", errors.New("ipfs_path must preffix is /"))
		return
	}
	ipfs_path = ipfs_path[1:]

	if len(ipfs_path) != 0 {
		callStat := func(result string, err error) {
			if err != nil {
				outerCall("", err)
				return
			}
			result = stringTrim(result)

			var nodeStat statInfo
			err = json.Unmarshal([]byte(result), &nodeStat)
			if err != nil {
				outerCall("", err)
				return
			}
			object_hash = strings.Trim(nodeStat.Hash, endsep)

			callls := func(result string, err error) {
				if err != nil {
					outerCall("", err)
					return
				}
				result = stringTrim(result)
				outerCall(result, nil)
			}
			cmd := strings.Join([]string{"ipfs", "ls", "--is-lib=true", object_hash}, cmdSep)
			_, _, err = ipfsAsyncCmdTime(cmd, second, callls)
			if err != nil {
				outerCall("", err)
				return
			}
		}
		statCmd := strings.Join([]string{"ipfs", "object", "stat", "--is-lib=true", object_hash + "/" + ipfs_path}, cmdSep)
		_, _, err := ipfsAsyncCmdTime(statCmd, second, callStat)
		if err != nil {
			outerCall("", err)
			return
		}
		return
	}

	call := func(result string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		result = stringTrim(result)
		outerCall(result, nil)
	}
	cmd := strings.Join([]string{"ipfs", "ls", "--is-lib=true", object_hash}, cmdSep)
	_, _, err = ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncMerge(root_hash, ipfs_path, share_hash string, second int, outerCall commands.RequestCB) {
	var err error
	if root_hash, err = ipfsObjectHashCheck(root_hash); err != nil {
		outerCall("", errors.New("object_hash len not 46"))
		return
	}
	if share_hash, err = ipfsObjectHashCheck(share_hash); err != nil {
		outerCall("", errors.New("object_hash len not 46"))
		return
	}

	if len(ipfs_path) == 0 {
		outerCall("", errors.New("ipfs_path len is 0"))
		return
	}

	ipfs_path, err = ipfsPathClean(ipfs_path)
	if err != nil {
		outerCall("", err)
		return
	}

	call := func(result string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		new_root := stringTrim(result)
		outerCall(new_root, nil)
	}
	cmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", root_hash, ipfs_path, share_hash}, cmdSep)
	_, _, err = ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncPeerid(new_id string, second int, outerCall commands.RequestCB) {
	if err := ipfsPeeridCheck(new_id); len(new_id) != 0 && err != nil {
		outerCall("", errors.New("new_id len is not 46 or is not 0"))
		return
	}

	cmd := strings.Join([]string{"ipfs", "config", "Identity.PeerID"}, cmdSep)
	if len(new_id) != 0 {
		cmd = strings.Join([]string{cmd, new_id}, cmdSep)
	}

	call := func(result string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		if new_id != "" {
			outerCall(new_id, nil)
			return
		}
		result = stringTrim(result)
		outerCall(result, nil)
	}
	_, _, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncPrivkey(new_key string, second int, outerCall commands.RequestCB) {
	cmd := strings.Join([]string{"ipfs", "config", "Identity.PrivKey"}, cmdSep)
	if len(new_key) != 0 {
		cmd = strings.Join([]string{cmd, new_key}, cmdSep)
	}
	call := func(result string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		if new_key != "" {
			outerCall(new_key, nil)
			return
		}
		result = stringTrim(result)
		outerCall(result, nil)
	}
	_, _, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncPublish(object_hash string, second int, outerCall commands.RequestCB) {
	var err error
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		outerCall("", errors.New("object_hash len is not 46"))
		return
	}

	call := func(result string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		result = stringTrim(result)
		outerCall(result, nil)
	}
	cmd := strings.Join([]string{"ipfs", "name", "publish", "--is-lib=true", object_hash}, cmdSep)
	_, _, err = ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		outerCall("", nil)
	}
}

func IpfsAsyncConnectPeer(peer_addr string, second int, outerCall commands.RequestCB) {
	if len(peer_addr) == 0 {
		outerCall("", errors.New("peer_addr len is 0"))
		return
	}

	call := func(result string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		outerCall(result, nil)
	}
	cmd := strings.Join([]string{"ipfs", "swarm", "connect", peer_addr}, cmdSep)
	_, _, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncConfig(key, value string, outerCall commands.RequestCB) {
	var cmd string
	if len(key) == 0 {
		cmd = strings.Join([]string{"ipfs", "config", "show"}, cmdSep)
	} else if len(key) != 0 && len(value) == 0 {
		cmd = strings.Join([]string{"ipfs", "config", key}, cmdSep)
	} else {
		cmd = strings.Join([]string{"ipfs", "config", key, value, "--json"}, cmdSep)
	}

	call := func(result string, err error) {
		fmt.Println("IpfsAsyncConfig result = ", result)
		if err != nil {
			outerCall("", err)
			return
		}
		outerCall(result, nil)
	}
	_, _, err := ipfsAsyncCmd(cmd, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncMessage(peer_id, peer_key, msg string, outerCall commands.RequestCB) {
	if err := ipfsPeeridCheck(peer_id); err != nil {
		outerCall("", errors.New("peer_id len is not 46"))
		return
	}
	if err := ipfsPeerkeyCheck(peer_key); err != nil {
		outerCall("", errors.New("peer_key len is not 8"))
		return
	}
	if len(msg) == 0 {
		outerCall("", errors.New("msg len is 0"))
		return
	}
	cmd := strings.Join([]string{"ipfs", "remotemsg", peer_id, peer_key, msg}, cmdSep)
	call := func(result string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		outerCall(result, nil)
	}
	_, _, err := ipfsAsyncCmdTime(cmd, 0, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncRemotepin(peer_id, peer_key, object_hash string, second int, outerCall commands.RequestCB) {
	var err error
	if err = ipfsPeeridCheck(peer_id); err != nil {
		outerCall("", errors.New("peer_id len is not 46"))
		return
	}
	if err = ipfsPeerkeyCheck(peer_key); err != nil {
		outerCall("", errors.New("peer_key len is not 8"))
		return
	}
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		outerCall("", errors.New("object_hash format error"))
		return
	}

	cmd := strings.Join([]string{"ipfs", "remotepin", peer_id, peer_key, object_hash}, cmdSep)
	call := func(result string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		outerCall(result, nil)
	}
	_, _, err = ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncRemotels(peer_id, peer_key, object_hash string, second int, outerCall commands.RequestCB) {
	var err error
	if err = ipfsPeeridCheck(peer_id); err != nil {
		outerCall("", errors.New("peer_id len is not 46"))
		return
	}
	if err = ipfsPeerkeyCheck(peer_key); err != nil {
		outerCall("", errors.New("peer_key len is not 8"))
		return
	}
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		outerCall("", errors.New("object_hash format error"))
		return
	}

	call := func(result string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		result = stringTrim(result)
		outerCall(result, nil)
	}
	cmd := strings.Join([]string{"ipfs", "remotels", peer_id, peer_key, object_hash}, cmdSep)
	_, _, err = ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		outerCall("", err)
	}
}

func IpfsAsyncRelaypin(relay_id, relay_key, peer_id, peer_key, object_hash string, second int, outerCall commands.RequestCB) {
	var err error
	if err = ipfsPeeridCheck(relay_id); err != nil {
		outerCall("", errors.New("relay_id len is not 46"))
		return
	}

	if err = ipfsPeerkeyCheck(relay_key); err != nil {
		outerCall("", errors.New("relay_key len is not 8"))
		return
	}

	if err = ipfsPeeridCheck(peer_id); err != nil {
		outerCall("", errors.New("peer_id len is not 46"))
		return
	}

	if err = ipfsPeerkeyCheck(peer_key); err != nil {
		outerCall("", errors.New("peer_key len is not 8"))
		return
	}

	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		outerCall("", errors.New("object_hash format error"))
		return
	}

	call := func(result string, err error) {
		if err != nil {
			outerCall("", err)
			return
		}
		outerCall(result, nil)
	}
	cmd := strings.Join([]string{"ipfs", "relaypin", relay_id, relay_key, peer_id, peer_key, object_hash}, cmdSep)
	_, _, err = ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		outerCall("", err)
	}
}

func ipfsAsyncCmd(cmd string, call commands.RequestCB) (int, string, error) {
	return ipfsAsyncCmdTime(cmd, 0, call)
}

func ipfsAsyncCmdTime(cmd string, second int, call commands.RequestCB) (r int, s string, e error) {
	if nil == asyncApiIns {
		call("", fmt.Errorf("Deamon not run"))
	}
	ipfsAsyncPath := asyncApiIns.AsyncPath()
	if len(strings.Trim(ipfsAsyncPath, " ")) > 0 {
		if second != 0 {
			timeout := "--timeout=" + strconv.Itoa(second) + "s"
			cmd = strings.Join([]string{cmd, "-c", ipfsAsyncPath, timeout}, cmdSep)
		} else {
			cmd = strings.Join([]string{cmd, "-c", ipfsAsyncPath}, cmdSep)

		}
	}
	fmt.Println(cmd)
	return asyncApiIns.AsyncApi(cmd, call, nil)
}

func ipfsAsyncCmdWithCancel(cmd string, second int, call commands.RequestCB, cancel chan struct{}) (r int, s string, e error) {
	if nil == asyncApiIns {
		call("", fmt.Errorf("Deamon not run"))
	}
	ipfsAsyncPath := asyncApiIns.AsyncPath()
	if len(strings.Trim(ipfsAsyncPath, " ")) > 0 {
		cmd = strings.Join([]string{cmd, "-c", ipfsAsyncPath}, cmdSep)
	}
	fmt.Println(cmd)
	return asyncApiIns.AsyncApi(cmd, call, cancel)
}

func stringTrim(src string) string {
	return strings.Trim(src, endsep)
}
