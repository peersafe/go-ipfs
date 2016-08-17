package ipfsmobile

import (
	"errors"
	"strconv"
	"strings"

	"github.com/ipfs/go-ipfs/cmd/ipfs_lib"
	uuid "gx/ipfs/QmcyaFHbyiZfoX5GTpcqqCPYmbjYNAhRDekXSJPFHdYNSV/go.uuid"
)

const (
	ADD_TYPE int = iota
	GET_TYPE
)

var (
	globalCallBack IpfsCallBack
	cmdSep         string = "&X&"
)

type IpfsCallBack interface {
	Daemon(status int, err error)
	Add(uid, hash string, pos int, err error)
	Get(uid string, pos int, err error)
	Query(root_hash, ipfs_path, result string, err error)
	Publish(publish_hash string, err error)
	ConnectPeer(peer_addr string, err error)
}

func IpfsInit(path string) error {
	_, err := ipfs_lib.IpfsAsyncInit(path)
	return err
}

func IpfsAsyncDaemon(path string, call IpfsCallBack) {
	if call != nil {
		globalCallBack = call
	} else {
		globalCallBack.Daemon(1, errors.New("error: IpfsAsyncDaemon call parameter is nil!"))
		return
	}
	outerCall := func(result string, err error) {
		if err != nil {
			globalCallBack.Daemon(1, err)
			return
		}
		if result == "Start" {
			globalCallBack.Daemon(0, nil)
		}
		if result == "Shutdown" {
			globalCallBack.Daemon(1, nil)
		}
	}
	ipfs_lib.IpfsAsyncDaemon(path, outerCall)
}

func IpfsShutdown() (retErr error) {
	sync := make(chan struct{})
	outerCall := func(result string, err error) {
		if err != nil {
			retErr = err
			sync <- struct{}{}
			return
		}
		retErr = nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncShutDown(outerCall)
	<-sync
	return
}

func IpfsAsyncAdd(os_path string, second int) string {
	uid := geneUuid()
	outerCall := func(result string, err error) {
		if err != nil {
			globalCallBack.Add(uid, "", 0, err)
			return
		}
		// do progress callback
		if !strings.Contains(result, "Over") && !strings.HasPrefix(result, "Qm") {
			results := strings.Split(result, cmdSep)
			total, _ := strconv.ParseInt(results[0], 10, 64)
			current, _ := strconv.ParseInt(results[1], 10, 64)
			pos := int((current / total) * 100)
			globalCallBack.Add(uid, "", pos, nil)
			return
		}

		add_hash := result
		globalCallBack.Add(uid, add_hash, 100, nil)
	}
	ipfs_lib.IpfsAsyncAdd(os_path, second, outerCall)
	return uid
}

func IpfsDelete(root_hash, ipfs_path string, second int) (new_root string, retErr error) {
	sync := make(chan struct{})
	outerCall := func(result string, err error) {
		if err != nil {
			new_root, retErr = "", err
			sync <- struct{}{}
			return
		}
		new_root, retErr = result, nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncDelete(root_hash, ipfs_path, second, outerCall)
	<-sync
	return
}

func IpfsMove(root_hash, ipfs_src_path, ipfs_dst_path string, second int) (new_root string, retErr error) {
	sync := make(chan struct{})
	outerCall := func(result string, err error) {
		if err != nil {
			new_root, retErr = "", err
			sync <- struct{}{}
			return
		}
		new_root, retErr = result, nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncMove(root_hash, ipfs_src_path, ipfs_dst_path, second, outerCall)
	<-sync
	return
}

func IpfsShare(object_hash, share_name string, sencond int) (new_hash string, retErr error) {
	sync := make(chan struct{})
	outerCall := func(result string, err error) {
		if err != nil {
			new_hash, retErr = "", err
			sync <- struct{}{}
			return
		}
		new_hash, retErr = result, nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncShard(object_hash, share_name, sencond, outerCall)
	<-sync
	return
}

func IpfsAsyncGet(share_hash, save_path string, second int) string {
	uid := geneUuid()
	outerCall := func(result string, err error) {
		if err != nil {
			globalCallBack.Get(uid, 0, err)
			return
		}
		// do progress callback
		if result != "" && !strings.Contains(result, "Over") {
			results := strings.Split(result, cmdSep)
			total, _ := strconv.ParseInt(results[0], 10, 64)
			current, _ := strconv.ParseInt(results[1], 10, 64)
			pos := int((current / total) * 100)
			globalCallBack.Get(uid, pos, nil)
			return
		}
		globalCallBack.Get(uid, 100, nil)
	}
	ipfs_lib.IpfsAsyncGet(share_hash, save_path, second, outerCall)
	return uid
}

func IpfsAsyncQuery(object_hash, ipfs_path string, second int) {
	outerCall := func(result string, err error) {
		if err != nil {
			globalCallBack.Query(object_hash, ipfs_path, "", err)
			return
		}

		globalCallBack.Query(object_hash, ipfs_path, result, err)
	}
	ipfs_lib.IpfsAsyncQuery(object_hash, ipfs_path, second, outerCall)
}

func IpfsMerge(root_hash, ipfs_path, share_hash string, second int) (new_root string, retErr error) {
	sync := make(chan struct{})
	outerCall := func(result string, err error) {
		if err != nil {
			new_root, retErr = "", err
			sync <- struct{}{}
			return
		}
		new_root, retErr = result, nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncMerge(root_hash, ipfs_path, share_hash, second, outerCall)
	<-sync
	return
}

func IpfsPeerid(new_id string, second int) (id string, retErr error) {
	sync := make(chan struct{})
	outerCall := func(result string, err error) {
		if err != nil {
			id, retErr = "", err
			sync <- struct{}{}
			return
		}
		id, retErr = result, nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncPeerid(new_id, second, outerCall)
	<-sync
	return
}

func IpfsPrivkey(new_key string, second int) (key string, retErr error) {
	sync := make(chan struct{})
	outerCall := func(result string, err error) {
		if err != nil {
			key, retErr = "", err
			sync <- struct{}{}
			return
		}
		key, retErr = result, nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncPrivkey(new_key, second, outerCall)
	<-sync
	return
}

func IpfsAsyncPublish(object_hash string, second int) {
	outerCall := func(result string, err error) {
		if err != nil {
			globalCallBack.Publish("", err)
			return
		}
		publish_hash := result
		globalCallBack.Publish(publish_hash, nil)
	}
	ipfs_lib.IpfsAsyncPublish(object_hash, second, outerCall)
}

func IpfsAsyncConnectpeer(peer_addr string, second int) {
	outerCall := func(result string, err error) {
		if err != nil {
			globalCallBack.ConnectPeer(peer_addr, err)
			return
		}
		globalCallBack.ConnectPeer(peer_addr, nil)
	}
	ipfs_lib.IpfsAsyncConnectPeer(peer_addr, second, outerCall)
}

func IpfsConfig(key, value string) (retValue string, retErr error) {
	sync := make(chan struct{})
	outerCall := func(result string, err error) {
		if err != nil {
			retValue, retErr = "", err
			sync <- struct{}{}
			return
		}
		retValue, retErr = result, nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncConfig(key, value, outerCall)
	<-sync
	return
}

func IpfsRemotepin(peer_id, peer_key, object_hash string, second int) (retErr error) {
	sync := make(chan struct{})
	outerCall := func(result string, err error) {
		if err != nil {
			retErr = err
			sync <- struct{}{}
			return
		}
		retErr = nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncRemotepin(peer_id, peer_key, object_hash, second, outerCall)
	<-sync
	return
}

func IpfsRemotels(peer_id, peer_key, object_hash string, second int) (lsResult string, retErr error) {
	sync := make(chan struct{})
	outerCall := func(result string, err error) {
		if err != nil {
			lsResult, retErr = "", err
			sync <- struct{}{}
			return
		}
		lsResult, retErr = result, nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncRemotels(peer_id, peer_key, object_hash, second, outerCall)
	<-sync
	return
}

func geneUuid() string {
	return uuid.NewV4().String()
}
