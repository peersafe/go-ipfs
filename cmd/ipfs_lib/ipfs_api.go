package ipfs_lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ipfs/go-ipfs/cmd/ipfs_lib/apiinterface"
)

const (
	keyLen int = 1596
	endsep     = "\n"
	errRet     = -1
	sucRet     = 0
	cmdSep     = "&X&"
)

var ipfsPath string

type statInfo struct {
	Hash string
}

var GApiCmd apicmd

type apicmd struct {
}

func (a apicmd) Cmd(str string, sec int) (int, string, error) {
	return ipfsCmdTime(str, sec)
}

func init() {
	apiinterface.GApiInterface = GApiCmd
}

func IpfsPath(path string) (int, string) {
	if path != "" {
		ipfsPath = path
		return sucRet, ""
	}
	return errRet, "path is nil"
}

func IpfsInit() (int, string) {
	cmd := strings.Join([]string{"ipfs", "init", "-e"}, cmdSep)
	fmt.Println(cmd)
	_, str, err := ipfsCmd(cmd)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	index := strings.LastIndex(str, ":")
	str = strings.Trim(str[index+1:], " ")

	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsDaemon() (int, string) {
	cmd := strings.Join([]string{"ipfs", "daemon"}, cmdSep)
	_, str, err := ipfsCmd(cmd)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsShutDown() (int, string) {
	cmd := strings.Join([]string{"ipfs", "shutdown"}, cmdSep)
	_, str, err := ipfsCmd(cmd)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsId(second int) (int, string) {
	cmd := strings.Join([]string{"ipfs", "id"}, cmdSep)
	_, str, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsAdd(root_hash, ipfs_path, os_path string, second int) (int, string) {
	var err error
	if root_hash, err = ipfsObjectHashCheck(root_hash); err != nil {
		fmt.Println("root_hash len not 46")
		return errRet, ""
	}

	if len(ipfs_path) == 0 {
		fmt.Println("ipfs_path len is 0")
		return errRet, ""
	}

	ipfs_path, err = ipfsPathClean(ipfs_path)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	var addHash string
	if len(os_path) != 0 {
		os_path, err := filepath.Abs(path.Clean(os_path))
		if err != nil {
			fmt.Println(err)
			return errRet, ""
		}

		fi, err := os.Lstat(os_path)
		if err != nil {
			fmt.Println(err)
			return errRet, ""
		}

		cmdSuff := ""
		if fi.Mode().IsDir() {
			cmdSuff = strings.Join([]string{"ipfs", "add", "--is-lib=true", "-r", os_path}, cmdSep)
		} else if fi.Mode().IsRegular() {
			cmdSuff = strings.Join([]string{"ipfs", "add", "--is-lib=true", os_path}, cmdSep)
		} else {
			return errRet, ""
		}

		_, addHash, err = ipfsCmdTime(cmdSuff, second)
		if err != nil {
			return errRet, ""
		}
	}

	ipfs_path = path.Clean(ipfs_path)
	cmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", root_hash, ipfs_path, addHash}, cmdSep)
	_, str, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsDelete(root_hash, ipfs_path string, second int) (int, string) {
	var err error
	if root_hash, err = ipfsObjectHashCheck(root_hash); err != nil {
		fmt.Println("root_hash len not 46")
		return errRet, ""
	}

	if len(ipfs_path) == 0 {
		fmt.Println("ipfs_path len is 0")
		return errRet, ""
	}

	ipfs_path, err = ipfsPathClean(ipfs_path)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	cmd := strings.Join([]string{"ipfs", "object", "patch", "rm-link", root_hash, ipfs_path}, cmdSep)
	_, str, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsMove(root_hash, ipfs_path_src, ipfs_path_des string, second int) (int, string) {
	var err error
	if root_hash, err = ipfsObjectHashCheck(root_hash); err != nil {
		fmt.Println("root_hash len not 46")
		return errRet, ""
	}

	if len(ipfs_path_src) == 0 {
		fmt.Println("ipfs_path_src len is 0")
		return errRet, ""
	}

	ipfs_path_src, err = ipfsPathClean(ipfs_path_src)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	if len(ipfs_path_des) == 0 {
		fmt.Println("ipfs_path_des len is 0")
		return errRet, ""
	}

	ipfs_path_des, err = ipfsPathClean(ipfs_path_des)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	object_path := ipfs_path_src
	if strings.HasPrefix(ipfs_path_src, "\"") && strings.HasSuffix(ipfs_path_src, "\"") {
		object_path = ipfs_path_src[1 : len(ipfs_path_src)-1]
	}

	statCmd := strings.Join([]string{"ipfs", "object", "stat", "--is-lib=true", root_hash + "/" + object_path}, cmdSep)
	_, statStr, err := ipfsCmdTime(statCmd, second)
	if err != nil {
		return errRet, ""
	}

	var nodeStat statInfo
	err = json.Unmarshal([]byte(statStr), &nodeStat)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}
	nodeStat.Hash = strings.Trim(nodeStat.Hash, endsep)

	addCmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", root_hash, ipfs_path_des, nodeStat.Hash}, cmdSep)
	_, newHash, err := ipfsCmdTime(addCmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	newHash = strings.Trim(newHash, endsep)
	delCmd := strings.Join([]string{"ipfs", "object", "patch", "rm-link", newHash, ipfs_path_src}, cmdSep)
	_, new_root_hash, err := ipfsCmdTime(delCmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	new_root_hash = strings.Trim(new_root_hash, endsep)
	return len(new_root_hash), new_root_hash
}

func IpfsShard(object_hash, shard_name string, second int) (int, string) {
	var err error
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		fmt.Println("object_hash len not 46")
		return errRet, ""
	}

	if len(shard_name) == 0 {
		fmt.Println("shard_name len is 0")
		return errRet, ""
	}

	shard_name, err = ipfsPathClean(shard_name)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	cmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn", shard_name, object_hash}, cmdSep)
	_, str, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsGet(shard_hash, os_path string, second int) int {
	var err error
	if shard_hash, err = ipfsHashCheck(shard_hash); err != nil {
		fmt.Println("shard_hash format error")
		return errRet
	}
	if len(os_path) == 0 {
		fmt.Println("shard_name len is 0")
		return errRet
	}

	os_path, err = filepath.Abs(path.Clean(os_path))
	if err != nil {
		fmt.Println(err)
		return errRet
	}

	cmd := strings.Join([]string{"ipfs", "get", shard_hash, "-o", os_path}, cmdSep)
	_, _, err = ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet
	}
	return sucRet
}

func IpfsQuery(object_hash, ipfs_path string, second int) (int, string) {
	var err error
	if object_hash, err = ipfsHashCheck(object_hash); err != nil {
		fmt.Println("object_hash len not 46")
		return errRet, ""
	}

	if !strings.HasPrefix(ipfs_path, "/") {
		fmt.Println("ipfs_path must preffix is -")
		return errRet, ""
	}
	ipfs_path = ipfs_path[1:]

	if len(ipfs_path) != 0 {
		statCmd := strings.Join([]string{"ipfs", "object", "stat", "--is-lib=true", object_hash + "/" + ipfs_path}, cmdSep)
		_, statStr, err := ipfsCmdTime(statCmd, second)
		if err != nil {
			fmt.Println(err)
			return errRet, ""
		}

		var nodeStat statInfo
		err = json.Unmarshal([]byte(statStr), &nodeStat)
		if err != nil {
			fmt.Println(err)
			return errRet, ""
		}
		object_hash = strings.Trim(nodeStat.Hash, endsep)
	}

	cmd := strings.Join([]string{"ipfs", "ls", "--is-lib=true", object_hash}, cmdSep)
	_, str, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsMerge(root_hash, ipfs_path, shard_hash string, second int) (int, string) {
	var err error
	if root_hash, err = ipfsObjectHashCheck(root_hash); err != nil {
		fmt.Println("root_hash len not 46")
		return errRet, ""
	}
	if shard_hash, err = ipfsObjectHashCheck(shard_hash); err != nil {
		fmt.Println("shard_hash len not 46")
		return errRet, ""
	}

	if len(ipfs_path) == 0 {
		fmt.Println("ipfs_path len is 0")
		return errRet, ""
	}

	ipfs_path, err = ipfsPathClean(ipfs_path)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	cmd := strings.Join([]string{"ipfs", "object", "patch", "add-link", root_hash, ipfs_path, shard_hash}, cmdSep)
	_, str, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsPeerid(new_id string, second int) (int, string) {
	if err := ipfsPeeridCheck(new_id); len(new_id) != 0 && err != nil {
		fmt.Println("new_id len is not 46 or is not 0")
		return errRet, ""
	}

	cmd := strings.Join([]string{"ipfs", "config", "Identity.PeerID"}, cmdSep)
	_, peeId, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	if len(new_id) != 0 {
		cmd = strings.Join([]string{cmd, new_id}, cmdSep)
		_, _, err := ipfsCmdTime(cmd, second)
		if err != nil {
			fmt.Println(err)
			return errRet, ""
		}
		peeId = new_id
	}

	peeId = strings.Trim(peeId, endsep)
	return len(peeId), peeId
}

func IpfsPrivkey(new_key string, second int) (int, string) {
	cmd := strings.Join([]string{"ipfs", "config", "Identity.PrivKey"}, cmdSep)
	_, key, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	if len(new_key) != 0 {
		cmd := strings.Join([]string{cmd, new_key}, cmdSep)
		_, _, err := ipfsCmdTime(cmd, second)
		if err != nil {
			fmt.Println(err)
			return errRet, ""
		}
		key = new_key
	}

	key = strings.Trim(key, endsep)
	return len(key), key
}

func IpfsPublish(object_hash string, second int) (int, string) {
	var err error
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		fmt.Println("object_hash len is not 52")
		return errRet, ""
	}

	cmd := strings.Join([]string{"ipfs", "name", "publish", "--is-lib=true", object_hash}, cmdSep)
	_, hash, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	hash = strings.Trim(hash, endsep)
	return len(hash), hash
}

func IpfsConnectPeer(peer_addr string, second int) (int, string) {
	if len(peer_addr) == 0 {
		fmt.Println("peer_addr len is 0")
		return errRet, ""
	}

	cmd := strings.Join([]string{"ipfs", "swarm", "connect", peer_addr}, cmdSep)
	_, str, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}
	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsConfig(key, value string) (int, string) {
	var cmd string
	if len(key) == 0 {
		cmd = strings.Join([]string{"ipfs", "config", "show"}, cmdSep)
	} else if len(key) != 0 && len(value) == 0 {
		cmd = strings.Join([]string{"ipfs", "config", key}, cmdSep)
	} else {
		cmd = strings.Join([]string{"ipfs", "config", key, value}, cmdSep)
	}

	_, str, err := ipfsCmd(cmd)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsRemotepin(peer_id, peer_key, object_hash string, second int) (int, string) {
	var err error
	if err = ipfsPeeridCheck(peer_id); err != nil {
		fmt.Println("peer_id len is not 46")
		return errRet, ""
	}
	if err = ipfsPeerkeyCheck(peer_key); err != nil {
		fmt.Println("peer_key len is not 8")
		return errRet, ""
	}
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		fmt.Println("object_hash format error")
		return errRet, ""
	}

	cmd := strings.Join([]string{"ipfs", "remotepin", peer_id, peer_key, object_hash}, cmdSep)

	_, str, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsRelaypin(relay_id, relay_key, peer_id, peer_key, object_hash string, second int) (int, string) {
	var err error
	if err = ipfsPeeridCheck(relay_id); err != nil {
		fmt.Println("relay_id len is not 46")
		return errRet, ""
	}

	if err = ipfsPeerkeyCheck(relay_key); err != nil {
		fmt.Println("relay_key len is not 8")
		return errRet, ""
	}

	if err = ipfsPeeridCheck(peer_id); err != nil {
		fmt.Println("peer_id len is not 46")
		return errRet, ""
	}

	if err = ipfsPeerkeyCheck(peer_key); err != nil {
		fmt.Println("peer_key len is not 8")
		return errRet, ""
	}

	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		fmt.Println("object_hash format error")
		return errRet, ""
	}

	cmd := strings.Join([]string{"ipfs", "relaypin", relay_id, relay_key, peer_id, peer_key, object_hash}, cmdSep)
	_, str, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsRemotels(peer_id, peer_key, object_hash string, second int) (int, string) {
	var err error
	if err = ipfsPeeridCheck(peer_id); err != nil {
		fmt.Println("peer_id len is not 46")
		return errRet, ""
	}
	if err = ipfsPeerkeyCheck(peer_key); err != nil {
		fmt.Println("peer_key len is not 8")
		return errRet, ""
	}
	if object_hash, err = ipfsObjectHashCheck(object_hash); err != nil {
		fmt.Println("object_hash format error")
		return errRet, ""
	}

	cmd := strings.Join([]string{"ipfs", "remotels", peer_id, peer_key, object_hash}, cmdSep)

	_, str, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	return len(str), str
}

func IpfsCmdApi(cmd string, second int) (int, string) {
	_, str, err := ipfsCmdTime(cmd, second)
	if err != nil {
		fmt.Println(err)
		return errRet, ""
	}

	str = strings.Trim(str, endsep)
	return len(str), str
}

func ipfsCmd(cmd string) (int, string, error) {
	return ipfsCmdTime(cmd, 0)
}

func ipfsCmdTime(cmd string, second int) (r int, s string, e error) {
	if len(strings.Trim(ipfsPath, " ")) > 0 {
		cmd = strings.Join([]string{cmd, "-c", ipfsPath}, cmdSep)
	}
	fmt.Println(cmd)
	return ipfsMain(cmd, second)
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

func ipfsPeerkeyCheck(peerkey string) error {
	matchstr := "^[a-zA-Z0-9-`=\\\\\\[\\];'\",./~!@#$%^&*()_+|{}:<>?]{8}$"
	if matched, err := regexp.MatchString(matchstr, peerkey); err != nil || !matched {
		return fmt.Errorf("invalid peerkey format")
	}
	return nil
}

func ipfsPeeridCheck(peerid string) error {
	matchstr := "^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{46}$"
	if matched, err := regexp.MatchString(matchstr, peerid); err != nil || !matched {
		return fmt.Errorf("invalid peerid format")
	}
	return nil
}

func ipfsObjectHashCheck(hash string) (string, error) {
	matchstr := "^((/ipfs/|addr://)?" +
		"[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{46})$"
	if matched, err := regexp.MatchString(matchstr, hash); err != nil || !matched {
		return "", fmt.Errorf("hash format error")
	}
	if strings.HasPrefix(hash, "addr://") {
		hash = strings.Replace(hash, "addr://", "/ipfs/", 1)
	}
	return hash, nil
}

func ipfsHashCheck(hash string) (string, error) {
	matchstr := "^((/ipfs/|/ipns/|peer://|addr://)?" +
		"[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]{46})$"
	if matched, err := regexp.MatchString(matchstr, hash); err != nil || !matched {
		return "", fmt.Errorf("hash format error")
	}
	if strings.HasPrefix(hash, "peer://") {
		hash = strings.Replace(hash, "peer://", "/ipns/", 1)
	} else if strings.HasPrefix(hash, "addr://") {
		hash = strings.Replace(hash, "addr://", "/ipfs/", 1)
	}
	return hash, nil
}
