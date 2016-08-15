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

var asyncApiIns Instance

func InitApi() {
	asyncApiIns = NewInstance()
}

func AsyncCmd(cmd string, call commands.CallFunc) (int, string, error) {
	return asyncApiIns.AsyncApi(cmd, call)
}

var ipfsAsyncPath string

func IpfsAsyncPath(path string) (int, string) {
	if path != "" {
		ipfsAsyncPath = path
		return SUCCESS, ""
	}
	return PARA_ERR, "path is nil"
}

func IpfsAsyncInit(call commands.CallFunc) (int, string) {
	cmd := strings.Join([]string{"ipfs", "init", "-e"}, cmdSep)
	fmt.Println(cmd)
	ret, str, err := ipfsAsyncCmd(cmd, call)
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

func IpfsAsyncDaemon(call commands.CallFunc) (int, string) {
	cmd := strings.Join([]string{"ipfs", "daemon"}, cmdSep)
	ret, str, err := ipfsAsyncCmd(cmd, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncShutDown(call commands.CallFunc) (int, string) {
	cmd := strings.Join([]string{"ipfs", "shutdown"}, cmdSep)
	ret, str, err := ipfsAsyncCmd(cmd, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncId(second int, call commands.CallFunc) (int, string) {
	cmd := strings.Join([]string{"ipfs", "id"}, cmdSep)
	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

type Mycall struct {
	myCall func(string, error)
}

func (call *Mycall) Call(result string, err error) {
	call.myCall(result, err)
}

func IpfsAsyncAdd(root_hash, ipfs_path, os_path string, second int, call commands.CallFunc) (int, string) {
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

	var ret int

	myCall := &Mycall{}
	myCall.myCall = func(result string, err error) {
		if err != nil {
			call.Call(result, err)
		}

		ipfs_path = path.Clean(ipfs_path)
		cmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", root_hash, ipfs_path, result}, cmdSep)
		ipfsAsyncCmdTime(cmd, second, call)
	}
	if len(os_path) != 0 {
		os_path, err := filepath.Abs(path.Clean(os_path))
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
	}
	return SUCCESS, ""
}

func IpfsAsyncDelete(root_hash, ipfs_path string, second int, call commands.CallFunc) (int, string) {
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
	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncMove(root_hash, ipfs_path_src, ipfs_path_des string, second int, call commands.CallFunc) (int, string) {
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
	myCall.myCall = func(result string, err error) {
		if err != nil {
			call.Call(result, err)
		}

		var nodeStat statInfo
		err = json.Unmarshal([]byte(result), &nodeStat)
		if err != nil {
			fmt.Println(err)
			return
		}
		nodeStat.Hash = strings.Trim(nodeStat.Hash, endsep)

		callPatch := &Mycall{}
		callPatch.myCall = func(result string, err error) {
			if err != nil {
				call.Call(result, err)
			}
			newHash := strings.Trim(result, endsep)
			delCmd := strings.Join([]string{"ipfs", "object", "patch", "rm-link", newHash, ipfs_path_src}, cmdSep)
			ipfsAsyncCmdTime(delCmd, second, call)
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

func IpfsAsyncShard(object_hash, shard_name string, second int, call commands.CallFunc) (int, string) {
	var err error
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		fmt.Println("object_hash len not 46")
		return PARA_ERR, ""
	}

	if len(shard_name) == 0 {
		fmt.Println("shard_name len is 0")
		return PARA_ERR, ""
	}

	shard_name, err = ipfsPathClean(shard_name)
	if err != nil {
		fmt.Println(err)
		return PARA_ERR, ""
	}

	cmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn", shard_name, object_hash}, cmdSep)
	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncGet(shard_hash, os_path string, second int, call commands.CallFunc) int {
	var err error
	if shard_hash, err = ipfsHashCheck(shard_hash); err != nil {
		fmt.Println("shard_hash format error")
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

	cmd := strings.Join([]string{"ipfs", "get", shard_hash, "-o", os_path}, cmdSep)
	ret, _, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret
	}
	return SUCCESS
}

func IpfsAsyncQuery(object_hash, ipfs_path string, second int, call commands.CallFunc) (int, string) {
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
				call.Call(result, err)
			}

			var nodeStat statInfo
			err = json.Unmarshal([]byte(result), &nodeStat)
			if err != nil {
				fmt.Println(err)
				return
			}
			object_hash = strings.Trim(nodeStat.Hash, endsep)
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

	cmd := strings.Join([]string{"ipfs", "ls", "--is-lib=true", object_hash}, cmdSep)
	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncMerge(root_hash, ipfs_path, shard_hash string, second int, call commands.CallFunc) (int, string) {
	var err error
	if root_hash, err = ipfsObjectHashCheck(root_hash); err != nil {
		fmt.Println("root_hash len not 46")
		return PARA_ERR, ""
	}
	if shard_hash, err = ipfsObjectHashCheck(shard_hash); err != nil {
		fmt.Println("shard_hash len not 46")
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

	cmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", root_hash, ipfs_path, shard_hash}, cmdSep)
	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncPeerid(new_id string, second int, call commands.CallFunc) (int, string) {
	if err := ipfsPeeridCheck(new_id); len(new_id) != 0 && err != nil {
		fmt.Println("new_id len is not 46 or is not 0")
		return PARA_ERR, ""
	}

	cmd := strings.Join([]string{"ipfs", "config", "Identity.PeerID"}, cmdSep)

	if len(new_id) != 0 {
		cmd = strings.Join([]string{cmd, new_id}, cmdSep)
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

func IpfsAsyncPrivkey(new_key string, second int, call commands.CallFunc) (int, string) {
	cmd := strings.Join([]string{"ipfs", "config", "Identity.PrivKey"}, cmdSep)

	if len(new_key) != 0 {
		cmd = strings.Join([]string{cmd, new_key}, cmdSep)
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

func IpfsAsyncPublish(object_hash string, second int, call commands.CallFunc) (int, string) {
	var err error
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		fmt.Println("object_hash len is not 52")
		return PARA_ERR, ""
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

func IpfsAsyncConnectPeer(peer_addr string, second int, call commands.CallFunc) (int, string) {
	if len(peer_addr) == 0 {
		fmt.Println("peer_addr len is 0")
		return PARA_ERR, ""
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

func IpfsAsyncConfig(key, value string, call commands.CallFunc) (int, string) {
	var cmd string
	if len(key) == 0 {
		cmd = strings.Join([]string{"ipfs", "config", "show"}, cmdSep)
	} else if len(key) != 0 && len(value) == 0 {
		cmd = strings.Join([]string{"ipfs", "config", key}, cmdSep)
	} else {
		cmd = strings.Join([]string{"ipfs", "config", key, value}, cmdSep)
	}

	ret, str, err := ipfsAsyncCmd(cmd, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncRemotepin(peer_id, peer_key, object_hash string, second int, call commands.CallFunc) (int, string) {
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

	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncRelaypin(relay_id, relay_key, peer_id, peer_key, object_hash string, second int, call commands.CallFunc) (int, string) {
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

	cmd := strings.Join([]string{"ipfs", "relaypin", relay_id, relay_key, peer_id, peer_key, object_hash}, cmdSep)
	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncRemotels(peer_id, peer_key, object_hash string, second int, call commands.CallFunc) (int, string) {
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

	cmd := strings.Join([]string{"ipfs", "remotels", peer_id, peer_key, object_hash}, cmdSep)

	ret, str, err := ipfsAsyncCmdTime(cmd, second, call)
	if err != nil {
		fmt.Println(err)
		return ret, ""
	}

	str = strings.Trim(str, endsep)
	return SUCCESS, str
}

func IpfsAsyncCmdApi(cmd string, second int, call commands.CallFunc) (int, string) {
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
	return AsyncCmd(cmd, call)
}
