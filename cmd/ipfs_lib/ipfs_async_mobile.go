package ipfs_lib

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	homedir "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/mitchellh/go-homedir"
	"github.com/ipfs/go-ipfs/commands"
)

func Ipfs_async_path(path string) string {
	homedir.Home_Unix_Dir = path
	res, str := IpfsAsyncPath(path)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_init(call commands.CallFunc) string {
	InitApi()
	res, str := IpfsAsyncInit(call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_cmd_arm(cmd string, second int, call commands.CallFunc) string {
	res, str := IpfsAsyncCmdApi(cmd, second, call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_daemon(call commands.CallFunc) string {
	res, str := IpfsAsyncDaemon(call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_shutdown(call commands.CallFunc) string {
	res, str := IpfsAsyncShutDown(call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_config(key, value string, call commands.CallFunc) string {
	res, str := IpfsAsyncConfig(key, value, call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_id(second int, call commands.CallFunc) string {
	res, str := IpfsAsyncId(second, call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_peerid(new_id string, second int, call commands.CallFunc) string {
	res, str := IpfsAsyncPeerid(new_id, second, call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_privkey(new_key string, second int, call commands.CallFunc) string {
	res, str := IpfsAsyncPrivkey(new_key, second, call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_add(os_path string, second int, call commands.CallFunc) string {
	if len(os_path) != 0 {
		os_path, err := filepath.Abs(path.Clean(os_path))
		if err != nil {
			return fmt.Sprintf("%d%s%s", PARA_ERR, cmdSep, "")
		}

		fi, err := os.Lstat(os_path)
		cmdSuff := ""
		if fi.Mode().IsDir() {
			cmdSuff = strings.Join([]string{"ipfs", "add", "--is-lib=true", "-r", os_path}, cmdSep)
		} else if fi.Mode().IsRegular() {
			cmdSuff = strings.Join([]string{"ipfs", "add", "--is-lib=true", os_path}, cmdSep)
		} else {
			return fmt.Sprintf("%d%s%s", PARA_ERR, cmdSep, "")
		}
		res, addHash, err := ipfsAsyncCmdTime(cmdSuff, second, call)
		if err != nil {
			return fmt.Sprintf("%d%s%s", res, cmdSep, "")
		}
		addHash = strings.Trim(addHash, endsep)
		return fmt.Sprintf("%d%s%s", res, cmdSep, addHash)

	} else {
		return fmt.Sprintf("%d%s%s", PARA_ERR, cmdSep, "")
	}
}

func Ipfs_async_get(shard_hash, os_path string, second int, call commands.CallFunc) string {
	res := IpfsAsyncGet(shard_hash, os_path, second, call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, "")
}

func Ipfs_async_publish(object_hash string, second int, call commands.CallFunc) string {
	res, str := IpfsAsyncPublish(object_hash, second, call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_remotepin(peer_id, peer_key, object_hash string, second int, call commands.CallFunc) string {
	res, str := IpfsAsyncRemotepin(peer_id, peer_key, object_hash, second, call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_remotels(peer_id, peer_key, object_hash string, second int, call commands.CallFunc) string {
	res, str := IpfsAsyncRemotels(peer_id, peer_key, object_hash, second, call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}

func Ipfs_async_connectpeer(remote_peer string, second int, call commands.CallFunc) string {
	res, str := IpfsAsyncConnectPeer(remote_peer, second, call)
	return fmt.Sprintf("%d%s%s", res, cmdSep, str)
}
