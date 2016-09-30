package ipfsmobile

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ipfs/go-ipfs/cmd/ipfs_lib"
	"github.com/ipfs/go-ipfs/cmd/ipfs_mobile/callback"
	uuid "gx/ipfs/QmcyaFHbyiZfoX5GTpcqqCPYmbjYNAhRDekXSJPFHdYNSV/go.uuid"
)

var (
	globalCallBack IpfsCallBack
	cmdSep         string = "&X&"
	loaderMap      map[string]chan struct{}
)

type IpfsCallBack interface {
	Daemon(status int, err string)
	Add(uid, hash string, pos int, err string)
	Get(uid string, pos int, err string)
	Query(root_hash, ipfs_path, result string, err string)
	Publish(publish_hash string, err string)
	ConnectPeer(peer_addr string, err string)
	Message(msg string, err string)
}

func IpfsInit(path string) error {
	_, err := ipfs_lib.IpfsAsyncInit(path)
	return err
}

func IpfsAsyncDaemon(path string, call IpfsCallBack) {
	if call != nil {
		globalCallBack = call
		callback.GlobalCallBack = call
	} else {
		globalCallBack.Daemon(1, "IpfsAsyncDaemon call parameter is nil!")
		return
	}
	outerCall := func(result string, err error) {
		if err != nil {
			globalCallBack.Daemon(1, err.Error())
			return
		}
		if result == "Start" {
			globalCallBack.Daemon(0, "")
		}
		if result == "Shutdown" {
			globalCallBack.Daemon(1, "")
		}
	}
	// uploaderMap and downloaderMap init
	loaderMap = make(map[string]chan struct{})
	ipfs_lib.IpfsAsyncDaemon(path, outerCall)
}

func IpfsShutdown() (retErr error) {
	sync := make(chan struct{}, 1)
	defer close(sync)
	outerCall := func(result string, err error) {
		if err != nil {
			retErr = err
			sync <- struct{}{}
			return
		}
		if result == "" {
			return
		}
		retErr = nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncShutDown(outerCall)
	<-sync
	return
}

type bakpos struct {
	pos  int
	done bool
}

func IpfsAsyncAdd(os_path string, second int) string {
	uid := geneUuid()
	bakPos := &bakpos{0, false}

	heartBeat := make(chan struct{})
	go func() {
		timer := time.NewTimer(time.Second * time.Duration(second))
		for {
			select {
			case <-heartBeat:
				timer = time.NewTimer(time.Second * time.Duration(second))
			case <-timer.C:
				globalCallBack.Add(uid, "", bakPos.pos, "timeout")
				ipfsDone(uid)
				return
			default:
				if bakPos.done {
					return
				}
			}
		}
	}()
	outerCall := func(result string, err error) {
		if err != nil {
			globalCallBack.Add(uid, "", bakPos.pos, err.Error())
			bakPos.done = true
			return
		}
		// do progress callback
		if !strings.Contains(result, "Over") && !strings.HasPrefix(result, "Qm") {
			results := strings.Split(result, cmdSep)
			total, _ := strconv.ParseFloat(results[0], 64)
			current, _ := strconv.ParseFloat(results[1], 64)
			pos := int((current / total) * 100)
			if pos == 100 || bakPos.pos == pos {
				return
			}

			heartBeat <- struct{}{}

			bakPos.pos = pos
			globalCallBack.Add(uid, "", pos, "")
			return
		}

		if !bakPos.done {
			add_hash := result
			globalCallBack.Add(uid, add_hash, 100, "")
			bakPos.done = true
			ipfsDone(uid)
		}
	}
	cancel := make(chan struct{})
	loaderMap[uid] = cancel
	ipfs_lib.IpfsAsyncAdd(os_path, second, outerCall, cancel)
	return uid
}

func IpfsDelete(root_hash, ipfs_path string, second int) (new_root string, retErr error) {
	sync := make(chan struct{})
	defer close(sync)
	outerCall := func(result string, err error) {
		if err != nil {
			new_root, retErr = "", err
			sync <- struct{}{}
			return
		}
		if result == "" {
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
	defer close(sync)
	outerCall := func(result string, err error) {
		if err != nil {
			new_root, retErr = "", err
			sync <- struct{}{}
			return
		}
		if result == "" {
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
	defer close(sync)

	outerCall := func(result string, err error) {
		if err != nil {
			new_hash, retErr = "", err
			sync <- struct{}{}
			return
		}
		if result == "" {
			return
		}
		new_hash, retErr = result, nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncShare(object_hash, share_name, sencond, outerCall)
	<-sync

	return
}

func IpfsAsyncGet(share_hash, save_path string, second int) string {
	uid := geneUuid()
	bakPos := &bakpos{0, false}
	heartBeat := make(chan struct{})
	go func() {
		timer := time.NewTimer(time.Duration(second) * time.Second)
		for {
			select {
			case <-heartBeat:
				timer = time.NewTimer(time.Duration(second) * time.Second)
			case <-timer.C:
				globalCallBack.Get(uid, bakPos.pos, "timeout")
				ipfsDone(uid)
				return
			default:
				if bakPos.done {
					return
				}
			}
		}
	}()
	outerCall := func(result string, err error) {
		if err != nil {
			globalCallBack.Get(uid, bakPos.pos, err.Error())
			bakPos.done = true
			return
		}
		// do progress callback
		if result != "" && !strings.Contains(result, "Over") {
			results := strings.Split(result, cmdSep)
			total, _ := strconv.ParseFloat(results[0], 64)
			current, _ := strconv.ParseFloat(results[1], 64)
			pos := int((current / total) * 100)
			if pos == 100 || bakPos.pos == pos {
				return
			}

			heartBeat <- struct{}{}

			bakPos.pos = pos
			globalCallBack.Get(uid, pos, "")
			return
		}
		if result == "" {
			return
		}

		if !bakPos.done {
			globalCallBack.Get(uid, 100, "")
			bakPos.done = true
			ipfsDone(uid)
		}
	}
	cancel := make(chan struct{})
	loaderMap[uid] = cancel
	ipfs_lib.IpfsAsyncGet(share_hash, save_path, second, outerCall, cancel)
	return uid
}

func IpfsAsyncQuery(object_hash, ipfs_path string, second int) {
	outerCall := func(result string, err error) {
		if err != nil {
			globalCallBack.Query(object_hash, ipfs_path, "", err.Error())
			return
		}
		globalCallBack.Query(object_hash, ipfs_path, result, "")
	}
	ipfs_lib.IpfsAsyncQuery(object_hash, ipfs_path, second, outerCall)
}

func IpfsQuery(object_hash, ipfs_path string, second int) (queryReuslt string, retErr error) {
	sync := make(chan struct{})
	defer close(sync)
	outerCall := func(result string, err error) {
		if err != nil {
			queryReuslt, retErr = "", err
			sync <- struct{}{}
			return
		}
		if result == "" {
			return
		}
		queryReuslt, retErr = result, nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncQuery(object_hash, ipfs_path, second, outerCall)
	<-sync

	return
}

func IpfsMerge(root_hash, ipfs_path, share_hash string, second int) (new_root string, retErr error) {
	sync := make(chan struct{})
	defer close(sync)
	outerCall := func(result string, err error) {
		if err != nil {
			new_root, retErr = "", err
			sync <- struct{}{}
			return
		}
		if result == "" {
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
	defer close(sync)
	outerCall := func(result string, err error) {
		if err != nil {
			id, retErr = "", err
			sync <- struct{}{}
			return
		}
		if result == "" {
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
	defer close(sync)
	outerCall := func(result string, err error) {
		if err != nil {
			key, retErr = "", err
			sync <- struct{}{}
			return
		}
		if result == "" {
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
			globalCallBack.Publish("", err.Error())
			return
		}
		publish_hash := result
		globalCallBack.Publish(publish_hash, "")
	}
	ipfs_lib.IpfsAsyncPublish(object_hash, second, outerCall)
}

func IpfsAsyncConnectpeer(peer_addr string, second int) {
	outerCall := func(result string, err error) {
		if err != nil {
			globalCallBack.ConnectPeer(peer_addr, err.Error())
			return
		}
		globalCallBack.ConnectPeer(peer_addr, "")
	}
	ipfs_lib.IpfsAsyncConnectPeer(peer_addr, second, outerCall)
}

func IpfsConfig(key, value string) (retValue string, retErr error) {
	sync := make(chan struct{}, 1)
	defer close(sync)
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
	defer close(sync)
	outerCall := func(result string, err error) {
		if err != nil {
			retErr = err
			sync <- struct{}{}
			return
		}
		if result == "" {
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
	defer close(sync)
	outerCall := func(result string, err error) {
		if err != nil {
			lsResult, retErr = "", err
			sync <- struct{}{}
			return
		}
		if result == "" {
			return
		}
		lsResult, retErr = result, nil
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncRemotels(peer_id, peer_key, object_hash, second, outerCall)
	<-sync
	return
}

func IpfsAsyncMessage(peer_id, peer_key, msg string) (ret int) {
	sync := make(chan struct{})
	defer close(sync)
	ret = 0
	outerCall := func(result string, err error) {
		fmt.Println("IpfsAsyncMessage err=", err)
		if err != nil {
			// secret failed for remotemsg
			if strings.Contains(err.Error(), "Secret authentication failed") {
				ret = -4
			}
			if strings.Contains(err.Error(), "dial attempt failed") {
				ret = -5
			}
			sync <- struct{}{}
			return
		}
		sync <- struct{}{}
	}
	ipfs_lib.IpfsAsyncMessage(peer_id, peer_key, msg, outerCall)
	<-sync
	return
}

func IpfsCancel(uuid string) {
	cancel, ok := loaderMap[uuid]
	if ok {
		cancel <- struct{}{}
		close(cancel)
		delete(loaderMap, uuid)
	}
}

func IpfsUuid() string {
	return geneUuid()
}

func ipfsDone(uuid string) {
	cancel, ok := loaderMap[uuid]
	if ok {
		close(cancel)
		delete(loaderMap, uuid)
	}
}

func geneUuid() string {
	return uuid.NewV4().String()
}

func IpfsPing(peer_id string) (ping bool) {
	sync := make(chan struct{})
	defer close(sync)
	outerCall := func(result string, err error) {
		if err != nil {
			sync <- struct{}{}
			ping = false
			return
		}
		if strings.Contains(result, "not found") {
			sync <- struct{}{}
			ping = false
			return
		}
		if strings.Contains(result, "time=") && strings.Contains(result, "ms") {
			sync <- struct{}{}
			ping = true
			return
		}
	}
	ipfs_lib.IpfsPing(peer_id, outerCall)
	<-sync
	return
}
