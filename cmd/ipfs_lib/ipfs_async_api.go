package ipfs_lib

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ipfs/go-ipfs/commands"
)

var (
	asyncApiIns    Instance
	ipfsAsyncPath  string
	globalCallBack IpfsCallBack
)

type Mycall struct {
	myCall func(string, error)
}

func (call *Mycall) Call(result string, err error) {
	call.myCall(result, err)
}

func initApi() {
	asyncApiIns = NewInstance()
}

func asyncCmd(cmd string, call commands.CallFunc) (int, string, error) {
	return asyncApiIns.AsyncApi(cmd, call)
}

func IpfsAsyncPath(path string) (int, string) {
	if path != "" {
		ipfsAsyncPath = path
		return SUCCESS, ""
	}
	return PARA_ERR, "path is nil"
}

func IpfsAsyncInit(call IpfsCallBack) (int, string) {
	if call != nil {
		globalCallBack = call
	} else {
		return PARA_ERR, "parameter IpfsCallBack is nil!"
	}

	cmd := strings.Join([]string{"ipfs", "init", "-e"}, cmdSep)
	fmt.Println(cmd)

	myCall := &Mycall{}
	myCall.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.Init(UNKOWN, err.Error())
			return
		}
		globalCallBack.Init(SUCCESS, "")
	}
	ret, str, err := ipfsAsyncCmd(cmd, myCall)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	index := strings.LastIndex(str, ":")
	str = strings.Trim(str[index+1:], " ")

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncDaemon() (int, string) {
	cmd := strings.Join([]string{"ipfs", "daemon"}, cmdSep)
	myCall := &Mycall{}
	myCall.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.Daemon(UNKOWN, err.Error())
			return
		}
		globalCallBack.Daemon(SUCCESS, "")
	}
	ret, str, err := ipfsAsyncCmd(cmd, myCall)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncShutDown() (int, string) {
	cmd := strings.Join([]string{"ipfs", "shutdown"}, cmdSep)
	myCall := &Mycall{}
	myCall.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.ShutDown(UNKOWN, err.Error())
			return
		}
		globalCallBack.ShutDown(SUCCESS, "")
	}
	ret, str, err := ipfsAsyncCmd(cmd, myCall)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncId(second int) (int, string) {
	cmd := strings.Join([]string{"ipfs", "id"}, cmdSep)
	myCall := &Mycall{}
	myCall.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.Id(UNKOWN, err.Error(), "")
			return
		}
		globalCallBack.Id(SUCCESS, "", result)
	}
	ret, str, err := ipfsAsyncCmdTime(cmd, second, myCall)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncAdd(root_hash, ipfs_path, os_path string, second int) (int, string) {
	var err error
	if root_hash, err = ipfsObjectHashCheck(root_hash); err != nil {
		fmt.Println("root_hash len not 46")
		globalCallBack.Add(PARA_ERR, "root_hash len not 46", "", ipfs_path, os_path, "")
		return PARA_ERR, ""
	}

	if len(ipfs_path) == 0 {
		fmt.Println("ipfs_path len is 0")
		globalCallBack.Add(PARA_ERR, "ipfs_path len is 0", "", ipfs_path, os_path, "")
		return PARA_ERR, ""
	}

	ipfs_path, err = ipfsPathClean(ipfs_path)
	if err != nil {
		fmt.Println(err)
		globalCallBack.Add(PARA_ERR, err.Error(), "", ipfs_path, os_path, "")
		return PARA_ERR, ""
	}

	var ret int

	myCall := &Mycall{}
	myCall.myCall = func(add_hash string, err error) {
		if err != nil {
			globalCallBack.Add(UNKOWN, err.Error(), "", ipfs_path, os_path, "")
			return
		}
		// do progress callback
		if !strings.Contains(add_hash, "Over") {
			results := strings.Split(add_hash, cmdSep)
			total, _ := strconv.ParseInt(results[0], 10, 64)
			current, _ := strconv.ParseInt(results[1], 10, 64)
			globalCallBack.Progress(ipfs_path, "", ADD_TYPE, total, current)
			return
		}
		// get real add_hash
		results := strings.Split(add_hash, cmdSep)
		add_hash = results[1]

		call := &Mycall{}
		call.myCall = func(result string, err error) {
			if err != nil {
				globalCallBack.Add(UNKOWN, err.Error(), "", ipfs_path, os_path, "")
				return
			}

			new_root := result
			globalCallBack.Add(SUCCESS, "", new_root, ipfs_path, os_path, add_hash)
		}
		ipfs_path = path.Clean(ipfs_path)
		cmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", root_hash, ipfs_path, add_hash}, cmdSep)
		ipfsAsyncCmdTime(cmd, second, call)
	}

	os_path, err = filepath.Abs(path.Clean(os_path))
	if err != nil {
		fmt.Println(err)
		return PARA_ERR, ""
	}

	fi, err := os.Lstat(os_path)
	if err != nil {
		fmt.Println(err)
		return PARA_ERR, ""
	}

	cmd := ""
	if fi.Mode().IsDir() {
		cmd = strings.Join([]string{"ipfs", "add", "--is-lib=true", "-r", os_path}, cmdSep)
	} else if fi.Mode().IsRegular() {
		cmd = strings.Join([]string{"ipfs", "add", "--is-lib=true", os_path}, cmdSep)
	} else {
		fmt.Println("Unkown file type!")
		return UNKOWN, ""
	}

	ret, _, err = ipfsAsyncCmdTime(cmd, second, myCall)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	return SUCCESS, ""
}

func IpfsAsyncDelete(root_hash, ipfs_path string, second int) (int, string) {
	var err error
	if root_hash, err = ipfsObjectHashCheck(root_hash); err != nil {
		fmt.Println("root_hash len not 46")
		return PARA_ERR, ""
	}

	if len(ipfs_path) == 0 {
		fmt.Println("ipfs_path len is 0")
		return PARA_ERR, ""
	}

	ipfs_path, err = ipfsPathClean(ipfs_path)
	if err != nil {
		fmt.Println(err)
		return PARA_ERR, ""
	}

	cmd := strings.Join([]string{"ipfs", "object", "patch", "rm-link", root_hash, ipfs_path}, cmdSep)
	myCall := &Mycall{}
	myCall.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.Delete(UNKOWN, err.Error(), "", ipfs_path)
			return
		}

		fmt.Printf("async_delete result=[%v]\n", result)
		new_root := ""
		fmt.Println("TODO: get new_root from result>>>>>>>>>>>>>>")

		globalCallBack.Delete(SUCCESS, "", new_root, ipfs_path)
	}
	ret, str, err := ipfsAsyncCmdTime(cmd, second, myCall)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncMove(root_hash, ipfs_path_src, ipfs_path_des string, second int) (int, string) {
	var err error
	if root_hash, err = ipfsObjectHashCheck(root_hash); err != nil {
		fmt.Println("root_hash len not 46")
		return PARA_ERR, ""
	}

	if len(ipfs_path_src) == 0 {
		fmt.Println("ipfs_path_src len is 0")
		return PARA_ERR, ""
	}

	ipfs_path_src, err = ipfsPathClean(ipfs_path_src)
	if err != nil {
		fmt.Println(err)
		return PARA_ERR, ""
	}

	if len(ipfs_path_des) == 0 {
		fmt.Println("ipfs_path_des len is 0")
		return PARA_ERR, ""
	}

	ipfs_path_des, err = ipfsPathClean(ipfs_path_des)
	if err != nil {
		fmt.Println(err)
		return PARA_ERR, ""
	}

	object_path := ipfs_path_src
	if strings.HasPrefix(ipfs_path_src, "\"") && strings.HasSuffix(ipfs_path_src, "\"") {
		object_path = ipfs_path_src[1 : len(ipfs_path_src)-1]
	}

	myCall := &Mycall{}
	myCall.myCall = func(info string, err error) {
		if err != nil {
			globalCallBack.Move(UNKOWN, err.Error(), "", ipfs_path_src, ipfs_path_des)
			return
		}

		var nodeStat statInfo
		err = json.Unmarshal([]byte(info), &nodeStat)
		if err != nil {
			fmt.Println(err)
			return
		}
		nodeStat.Hash = strings.Trim(nodeStat.Hash, endsep)

		callPatch := &Mycall{}
		callPatch.myCall = func(patch string, err error) {
			if err != nil {
				globalCallBack.Move(UNKOWN, err.Error(), "", ipfs_path_src, ipfs_path_des)
				return
			}
			callResult := &Mycall{}
			callResult.myCall = func(result string, err error) {
				if err != nil {
					globalCallBack.Move(UNKOWN, err.Error(), "", ipfs_path_src, ipfs_path_des)
					return
				}

				fmt.Printf("async_move result=[%v]\n", result)
				new_root := ""
				fmt.Println("TODO: get new_root from result>>>>>>>>>>>>>>")
				globalCallBack.Move(SUCCESS, "", new_root, ipfs_path_src, ipfs_path_des)
			}
			newHash := strings.Trim(patch, endsep)
			delCmd := strings.Join([]string{"ipfs", "object", "patch", "rm-link", newHash, ipfs_path_src}, cmdSep)
			ipfsAsyncCmdTime(delCmd, second, callResult)
		}
		addCmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", root_hash, ipfs_path_des, nodeStat.Hash}, cmdSep)
		_, _, err = ipfsAsyncCmdTime(addCmd, second, callPatch)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	statCmd := strings.Join([]string{"ipfs", "object", "stat", "--is-lib=true", root_hash + "/" + object_path}, cmdSep)
	ret, _, err := ipfsAsyncCmdTime(statCmd, second, myCall)
	if err != nil {
		return ret, ""
	}

	return SUCCESS, ""
}

func IpfsAsyncShard(object_hash, share_name string, second int) (int, string) {
	var err error
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		fmt.Println("object_hash len not 46")
		return PARA_ERR, ""
	}

	if len(share_name) == 0 {
		fmt.Println("share_name len is 0")
		return PARA_ERR, ""
	}

	share_name, err = ipfsPathClean(share_name)
	if err != nil {
		fmt.Println(err)
		return PARA_ERR, ""
	}

	myCall := &Mycall{}
	myCall.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.Share(UNKOWN, err.Error(), object_hash, share_name, "")
			return
		}
		fmt.Printf("async_share result=[%v]\n", result)
		new_hash := ""
		fmt.Println("TODO: get new_root from result>>>>>>>>>>>>>>")

		globalCallBack.Share(SUCCESS, "", object_hash, share_name, new_hash)
	}
	cmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn", share_name, object_hash}, cmdSep)
	ret, str, err := ipfsAsyncCmdTime(cmd, second, myCall)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncGet(share_hash, os_path string, second int) int {
	var err error
	if share_hash, err = ipfsHashCheck(share_hash); err != nil {
		fmt.Println("share_hash format error")
		return PARA_ERR
	}
	if len(os_path) == 0 {
		fmt.Println("shard_name len is 0")
		return PARA_ERR
	}

	os_path, err = filepath.Abs(path.Clean(os_path))
	if err != nil {
		fmt.Println(err)
		return PARA_ERR
	}

	myCall := &Mycall{}
	myCall.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.Get(UNKOWN, err.Error(), share_hash, os_path)
			return
		}
		if result == "" {
			fmt.Println("result is nil!")
			return
		}

		// do progress callback
		if !strings.Contains(result, "Over") {
			results := strings.Split(result, cmdSep)
			total, _ := strconv.ParseInt(results[0], 10, 64)
			current, _ := strconv.ParseInt(results[1], 10, 64)
			globalCallBack.Progress(os_path, share_hash, GET_TYPE, total, current)
			return
		}

		// get real result
		results := strings.Split(result, cmdSep)
		result = results[1]

		globalCallBack.Get(SUCCESS, "", share_hash, os_path)
	}
	cmd := strings.Join([]string{"ipfs", "get", share_hash, "-o", os_path}, cmdSep)
	ret, _, err := ipfsAsyncCmdTime(cmd, second, myCall)
	if err != nil {
		fmt.Println(err)
		return ret
	}
	return SUCCESS
}

func IpfsAsyncQuery(object_hash, ipfs_path string, second int) (int, string) {
	var err error
	if object_hash, err = ipfsHashCheck(object_hash); err != nil {
		fmt.Println("object_hash len not 46")
		return PARA_ERR, ""
	}

	if !strings.HasPrefix(ipfs_path, "/") {
		fmt.Println("ipfs_path must preffix is -")
		return PARA_ERR, ""
	}
	ipfs_path = ipfs_path[1:]

	if len(ipfs_path) != 0 {
		callStat := &Mycall{}
		callStat.myCall = func(result string, err error) {
			if err != nil {
				globalCallBack.Query(UNKOWN, err.Error(), object_hash, ipfs_path, "")
				return
			}

			var nodeStat statInfo
			err = json.Unmarshal([]byte(result), &nodeStat)
			if err != nil {
				fmt.Println(err)
				globalCallBack.Query(UNKOWN, err.Error(), object_hash, ipfs_path, "")
				return
			}
			object_hash = strings.Trim(nodeStat.Hash, endsep)

			call := &Mycall{}
			call.myCall = func(result string, err error) {
				if err != nil {
					globalCallBack.Query(UNKOWN, err.Error(), object_hash, ipfs_path, "")
					return
				}
				globalCallBack.Query(SUCCESS, "", object_hash, ipfs_path, result)
			}
			cmd := strings.Join([]string{"ipfs", "ls", "--is-lib=true", object_hash}, cmdSep)
			ipfsAsyncCmdTime(cmd, second, call)
		}
		statCmd := strings.Join([]string{"ipfs", "object", "stat", "--is-lib=true", object_hash + "/" + ipfs_path}, cmdSep)
		ret, _, err := ipfsAsyncCmdTime(statCmd, second, callStat)
		if err != nil {
			fmt.Println(err)
			return ret, ""
		}

		return SUCCESS, ""
	}

	call := &Mycall{}
	call.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.Query(UNKOWN, err.Error(), object_hash, ipfs_path, "")
			return
		}
		globalCallBack.Query(SUCCESS, "", object_hash, ipfs_path, result)
	}
	cmd := strings.Join([]string{"ipfs", "ls", "--is-lib=true", object_hash}, cmdSep)
	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncMerge(root_hash, ipfs_path, share_hash string, second int) (int, string) {
	var err error
	if root_hash, err = ipfsObjectHashCheck(root_hash); err != nil {
		fmt.Println("root_hash len not 46")
		return PARA_ERR, ""
	}
	if share_hash, err = ipfsObjectHashCheck(share_hash); err != nil {
		fmt.Println("share_hash len not 46")
		return PARA_ERR, ""
	}

	if len(ipfs_path) == 0 {
		fmt.Println("ipfs_path len is 0")
		return PARA_ERR, ""
	}

	ipfs_path, err = ipfsPathClean(ipfs_path)
	if err != nil {
		fmt.Println(err)
		return PARA_ERR, ""
	}

	call := &Mycall{}
	call.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.Merge(UNKOWN, err.Error(), "", ipfs_path, share_hash)
			return
		}

		fmt.Printf("async_merge result=[%v]\n", result)
		new_root := ""
		fmt.Println("TODO: get new_root from result>>>>>>>>>>>>>>")

		globalCallBack.Merge(SUCCESS, "", new_root, ipfs_path, share_hash)
	}
	cmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", root_hash, ipfs_path, share_hash}, cmdSep)
	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncPeerid(new_id string, second int) (int, string) {
	if err := ipfsPeeridCheck(new_id); len(new_id) != 0 && err != nil {
		fmt.Println("new_id len is not 46 or is not 0")
		globalCallBack.PeerId(PARA_ERR, "new_id len is not 46 or is not 0", "")
		return PARA_ERR, ""
	}

	cmd := strings.Join([]string{"ipfs", "config", "Identity.PeerID"}, cmdSep)
	if len(new_id) != 0 {
		cmd = strings.Join([]string{cmd, new_id}, cmdSep)
	}

	call := &Mycall{}
	call.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.PeerId(UNKOWN, err.Error(), "")
			return
		}
		if new_id != "" {
			globalCallBack.PeerId(SUCCESS, "", new_id)
			return
		}
		globalCallBack.PeerId(SUCCESS, "", result)
	}
	ret, peeId, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}
	if len(new_id) != 0 {
		peeId = new_id
	}

	peeId = strings.Trim(peeId, endsep)
	return SUCCESS, peeId
}

func IpfsAsyncPrivkey(new_key string, second int) (int, string) {
	cmd := strings.Join([]string{"ipfs", "config", "Identity.PrivKey"}, cmdSep)

	if len(new_key) != 0 {
		cmd = strings.Join([]string{cmd, new_key}, cmdSep)
	}
	call := &Mycall{}
	call.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.PrivateKey(UNKOWN, err.Error(), "")
			return
		}
		if new_key != "" {
			globalCallBack.PrivateKey(SUCCESS, "", new_key)
			return
		}
		globalCallBack.PrivateKey(SUCCESS, "", result)
	}
	ret, key, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}
	if len(new_key) != 0 {
		key = new_key
	}

	key = strings.Trim(key, endsep)
	return SUCCESS, key
}

func IpfsAsyncPublish(object_hash string, second int) (int, string) {
	var err error
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		fmt.Println("object_hash len is not 52")
		return PARA_ERR, ""
	}

	call := &Mycall{}
	call.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.Publish(UNKOWN, err.Error(), object_hash, "")
			return
		}
		globalCallBack.Publish(SUCCESS, "", object_hash, result)
	}
	cmd := strings.Join([]string{"ipfs", "name", "publish", "--is-lib=true", object_hash}, cmdSep)
	ret, hash, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	hash = strings.Trim(hash, endsep)
	return SUCCESS, hash
}

func IpfsAsyncConnectPeer(peer_addr string, second int) (int, string) {
	if len(peer_addr) == 0 {
		fmt.Println("peer_addr len is 0")
		return PARA_ERR, ""
	}

	call := &Mycall{}
	call.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.ConnectPeer(UNKOWN, err.Error(), peer_addr)
			return
		}
		globalCallBack.ConnectPeer(SUCCESS, "", peer_addr)
	}
	cmd := strings.Join([]string{"ipfs", "swarm", "connect", peer_addr}, cmdSep)
	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}
	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncConfig(key, value string) (int, string) {
	var cmd string
	if len(key) == 0 {
		cmd = strings.Join([]string{"ipfs", "config", "show"}, cmdSep)
	} else if len(key) != 0 && len(value) == 0 {
		cmd = strings.Join([]string{"ipfs", "config", key}, cmdSep)
	} else {
		cmd = strings.Join([]string{"ipfs", "config", key, value}, cmdSep)
	}

	call := &Mycall{}
	call.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.Config(UNKOWN, err.Error(), key, value)
			return
		}
		globalCallBack.Config(SUCCESS, "", key, value)
	}
	ret, str, err := ipfsAsyncCmd(cmd, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncRemotepin(peer_id, peer_key, object_hash string, second int) (int, string) {
	var err error
	if err = ipfsPeeridCheck(peer_id); err != nil {
		fmt.Println("peer_id len is not 46")
		return PARA_ERR, ""
	}
	if err = ipfsPeerkeyCheck(peer_key); err != nil {
		fmt.Println("peer_key len is not 8")
		return PARA_ERR, ""
	}
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		fmt.Println("object_hash format error")
		return PARA_ERR, ""
	}

	cmd := strings.Join([]string{"ipfs", "remotepin", peer_id, peer_key, object_hash}, cmdSep)
	call := &Mycall{}
	call.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.RemotePin(UNKOWN, err.Error(), peer_id, peer_key, object_hash)
			return
		}
		globalCallBack.RemotePin(SUCCESS, "", peer_id, peer_key, object_hash)
	}
	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncRelaypin(relay_id, relay_key, peer_id, peer_key, object_hash string, second int) (int, string) {
	var err error
	if err = ipfsPeeridCheck(relay_id); err != nil {
		fmt.Println("relay_id len is not 46")
		return PARA_ERR, ""
	}

	if err = ipfsPeerkeyCheck(relay_key); err != nil {
		fmt.Println("relay_key len is not 8")
		return PARA_ERR, ""
	}

	if err = ipfsPeeridCheck(peer_id); err != nil {
		fmt.Println("peer_id len is not 46")
		return PARA_ERR, ""
	}

	if err = ipfsPeerkeyCheck(peer_key); err != nil {
		fmt.Println("peer_key len is not 8")
		return PARA_ERR, ""
	}

	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		fmt.Println("object_hash format error")
		return PARA_ERR, ""
	}

	call := &Mycall{}
	call.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.RelayPin(UNKOWN, err.Error(), relay_id, relay_key, peer_id, peer_key, object_hash)
			return
		}
		globalCallBack.RelayPin(SUCCESS, "", relay_id, relay_key, peer_id, peer_key, object_hash)
	}
	cmd := strings.Join([]string{"ipfs", "relaypin", relay_id, relay_key, peer_id, peer_key, object_hash}, cmdSep)
	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncRemotels(peer_id, peer_key, object_hash string, second int) (int, string) {
	var err error
	if err = ipfsPeeridCheck(peer_id); err != nil {
		fmt.Println("peer_id len is not 46")
		return PARA_ERR, ""
	}
	if err = ipfsPeerkeyCheck(peer_key); err != nil {
		fmt.Println("peer_key len is not 8")
		return PARA_ERR, ""
	}
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		fmt.Println("object_hash format error")
		return PARA_ERR, ""
	}

	call := &Mycall{}
	call.myCall = func(result string, err error) {
		if err != nil {
			globalCallBack.Remotels(UNKOWN, err.Error(), peer_id, peer_key, object_hash, "")
			return
		}
		globalCallBack.Remotels(SUCCESS, "", peer_id, peer_key, object_hash, result)
	}
	cmd := strings.Join([]string{"ipfs", "remotels", peer_id, peer_key, object_hash}, cmdSep)
	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func ipfsAsyncCmd(cmd string, call commands.CallFunc) (int, string, error) {
	return ipfsAsyncCmdTime(cmd, 0, call)
}

func ipfsAsyncCmdTime(cmd string, second int, call commands.CallFunc) (r int, s string, e error) {
	if len(strings.Trim(ipfsAsyncPath, " ")) > 0 {
		if second != 0 {
			timeout := "--timeout=" + strconv.Itoa(second) + "s"
			cmd = strings.Join([]string{cmd, "-c", ipfsAsyncPath, timeout}, cmdSep)
		} else {
			cmd = strings.Join([]string{cmd, "-c", ipfsAsyncPath}, cmdSep)

		}
	}
	fmt.Println(cmd)
	return asyncCmd(cmd, call)
}
