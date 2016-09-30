package remotemsg

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	peer "gx/ipfs/QmWtbQU15LaB5B1JC2F7TV9P4K88vD3PpA4AJrwfCjhML8/go-libp2p-peer"
	ma "gx/ipfs/QmYzDkkgAEmrcNzFCiYo6L1dTX4EAG1gZkbtdbd9trL4vd/go-multiaddr"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
	host "gx/ipfs/Qmf4ETeAWXuThBfWwonVyFqGFSgTWepUDEr1txcctvpTXS/go-libp2p/p2p/host"
	inet "gx/ipfs/Qmf4ETeAWXuThBfWwonVyFqGFSgTWepUDEr1txcctvpTXS/go-libp2p/p2p/net"

	api "github.com/ipfs/go-ipfs/cmd/ipfs_lib/apiinterface"
	"github.com/ipfs/go-ipfs/cmd/ipfs_mobile/callback"
)

var log = logging.Logger("remotemsg")

const ID = "/ipfs/remotemsg"

type RemotemsgService struct {
	Host   host.Host
	Secret string
	ApiCmd api.Apier
}

func NewRemotemsgService(h host.Host, key string) *RemotemsgService {
	ps := &RemotemsgService{h, key, api.GApiInterface}
	h.SetStreamHandler(ID, ps.RemotemsgHandler)
	return ps
}

func (p *RemotemsgService) RemotemsgHandler(s inet.Stream) {
	for {
		slen := make([]byte, 2)
		_, err := io.ReadFull(s, slen)
		if err != nil {
			log.Errorf("remotels error:%v", err)
			return
		}

		newBuf := bytes.NewBuffer(slen)
		var bufLen uint16
		binary.Read(newBuf, binary.LittleEndian, &bufLen)

		buflen := int(bufLen)

		rbuf := make([]byte, buflen)
		_, err = io.ReadFull(s, rbuf)
		if err != nil {
			log.Errorf("remotemsg error:%v", err)
			return
		}

		var buf []byte = rbuf
		content, err := p.DecryptRequest(rbuf)
		if err != nil {
			buf = PKCS5Padding([]byte(fmt.Sprint(err)), len(rbuf))
		}
		fmt.Println(">>>>>>>>>>>>>remotemsg receive msg after decode =", string(content))
		err = p.remotemsg(content)
		if err != nil {
			buf = PKCS5Padding([]byte(fmt.Sprint(err)), len(rbuf))
		}
		_, err = s.Write(buf)
		if err != nil {
			log.Errorf("remotemsg error:%v", err)
			return
		}
	}
}

func (p *RemotemsgService) DecryptRequest(buf []byte) (rbuf []byte, err error) {
	fmt.Println("local peer secret=", p.Secret)

	md5hash_buf := buf[:md5.Size]
	crypted := buf[md5.Size:]

	orig, err := Decrypt(crypted, []byte(p.Secret))
	if err != nil {
		log.Errorf("remotemsg error:%v", err)
		return nil, err
	}

	md5hash := md5.Sum(orig)
	if !bytes.Equal(md5hash[:], md5hash_buf) {
		log.Errorf("Secret authentication failed")
		return nil, fmt.Errorf("Secret authentication failed")
	}

	rbuf = PKCS5UnPadding(orig)
	return rbuf, nil
}

type remoteMsg struct {
	MsgId          string `json:"uid"`
	Type           string `json:"type"`
	Hash           string `json:"hash"`
	FileName       string `json:"fileName"`
	MsgFromPeerId  string `json:"msg_from_peerid"`
	MsgFromPeerKey string `json:"msg_from_peerkey"`
	// if type="process",pos enable
	Pos int `json:"pos"`
	Ret int `json:"ret"`
	// for relaypin
	PeerId  string `json:"peer_id"`
	PeerKey string `json:"peer_key"`
	IsRelay bool   `json:"is_relay"`
}

func (ps *RemotemsgService) remotemsg(content []byte) error {
	// json unmarshal
	msg := new(remoteMsg)
	err := json.Unmarshal(content, msg)
	if err != nil {
		log.Errorf("remotemsg error=[%v]\n", err)
		return err
	}

	// store form peerID and privateKey
	fromPeerId, fromPeerKey := msg.MsgFromPeerId, msg.MsgFromPeerKey

	successServerMsg := func(types string) error {
		returnMsg := &remoteMsg{
			Type: types,
			Hash: msg.Hash,
		}
		data, err := json.Marshal(returnMsg)
		if err != nil {
			log.Errorf("remotemsg error:%v", err)
			return err
		}
		pid, err := peer.IDB58Decode(fromPeerId)
		if err != nil {
			log.Errorf("remotemsg error:%v", err)
			return err
		}
		if _, err := ps.RemoteMsg(context.TODO(), pid, fromPeerKey, string(data)); err != nil {
			log.Errorf("remotemsg error:%v", err)
			return err
		}
		return nil
	}

	file, _ := exec.LookPath(os.Args[0])
	app := filepath.Clean(file)
	app = filepath.ToSlash(app)
	app = filepath.Base(app)
	log.Debugf("remotemsg file[%v]", file)
	path, _ := filepath.Abs(file)

	if strings.HasSuffix(app, "ipfs") || strings.HasSuffix(app, "ipfs.exe") {
		switch msg.Type {
		case "remotepin":
			cmd := exec.Cmd{
				Path: path,
				Args: []string{"ipfs", "get", msg.Hash, "-o", "/dev/null"},
			}
			err := cmd.Run()
			if err != nil {
				log.Errorf("remotemsg>>>remotepin error:%v", err)
				return err
			}
			return successServerMsg("rRemotepin")
		case "remotels":
			cmd := exec.Cmd{
				Path: path,
				Args: []string{"ipfs", "ls", msg.Hash, "--timeout=30s"},
			}
			err := cmd.Run()
			if err != nil {
				log.Errorf("remotemsg>>>remotels error:%v", err)
				return err
			}
			return successServerMsg("rRemotels")
		case "relaypin":
			if msg.IsRelay { // local is master
				log.Debug("local is relay")
				cmd := exec.Cmd{
					Path: path,
					Args: []string{"ipfs", "ls", msg.Hash},
				}
				err := cmd.Run()
				if err != nil {
					log.Errorf("remotemsg->>>relaypin error:%v", err)
					return err
				}
				err = ps.relayPeer(msg.PeerId, msg.PeerKey, msg.Hash)
				if err != nil {
					log.Errorf("remotemsg->>>relaypin error:%v", err)
					return err
				}
			} else { // local is slave
				log.Debug("local is work")
				cmd := exec.Cmd{
					Path: path,
					Args: []string{"ipfs", "pin", "add", msg.Hash},
				}
				err := cmd.Run()
				if err != nil {
					log.Errorf("remotemsg->>>relaypin local work error:%v", err)
					return err
				}
			}
			return successServerMsg("rRelaypin")
		case "rRemotepin":
			fmt.Println("Remotepin command exec successfully!")
		case "rRemotels":
			fmt.Println("Remotels command exec successfully!")
		case "rRelaypin":
			fmt.Println("Relaypin command exec successfully!")
		}
	} else {
		callback.GlobalCallBack.Message(string(content), "")
	}
	return nil
}

func (ps *RemotemsgService) RemoteMsg(ctx context.Context, p peer.ID, key, msg string) (<-chan error, error) {
	s, err := ps.Host.NewStream(ctx, p, ID)
	if err != nil {
		return nil, err
	}

	out := make(chan error)
	go func() {
		defer close(out)
		select {
		case <-ctx.Done():
			return
		default:
			log.Debugf("remotemsg [%s][%s]", ID, msg)
			if msg == "" {
				return
			}
			_, err := remotemsg(s, key, msg)
			if err != nil {
				log.Errorf("remotemsg error:%v", err)
			}

			select {
			case out <- err:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, nil
}

func remotemsg(s inet.Stream, key, msg string) (time.Duration, error) {
	before := time.Now()

	orig := PKCS5Padding([]byte(msg), aes.BlockSize)
	md5hash := md5.Sum(orig)

	crypted, err := Encrypt(orig, []byte(key))
	if err != nil {
		return 0, err
	}

	buf := append(md5hash[:], crypted...)
	buflen := len(buf)

	slen := make([]byte, 2)
	binary.LittleEndian.PutUint16(slen, uint16(buflen))

	_, err = s.Write(slen)
	if err != nil {
		return 0, err
	}

	_, err = s.Write(buf)
	if err != nil {
		return 0, err
	}

	rbuf := make([]byte, buflen)
	_, err = io.ReadFull(s, rbuf)
	if err != nil {
		return 0, err
	}

	if !bytes.Equal(buf, rbuf) {
		str := PKCS5UnPadding(rbuf)
		return 0, errors.New(string(str))
	}

	return time.Now().Sub(before), nil
}

func (ps *RemotemsgService) relayPeer(p, key, fpath string) error {
	log.Debugf("relayPeer remote[%v]", p)

	_, pid, err := parsePeerParam(p)
	if err != nil {
		return fmt.Errorf("peer addr hash format error")
	}

	ps.Host.Network().ConnsToPeer(pid)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	out, err := ps.Relaypin(ctx, pid, key, p, key, fpath, false)
	if err != nil {
		return fmt.Errorf("routing remote error")
	}

	select {
	case reerr := <-out:
		if reerr != nil {
			return reerr
		}
	case <-ctx.Done():
		return fmt.Errorf("routing timeout")
	}
	return nil
}

func (ps *RemotemsgService) Relaypin(ctx context.Context, relayID peer.ID, relayKey, peerID, peerKey, hash string, relay bool) (<-chan error, error) {
	s, err := ps.Host.NewStream(ctx, relayID, ID)
	if err != nil {
		return nil, err
	}

	out := make(chan error)
	go func() {
		defer close(out)
		select {
		case <-ctx.Done():
			return
		default:
			log.Debugf("[%v][%v][%v][%v][%v]", relayID, relayKey, peerID, peerKey, hash)
			msg := &remoteMsg{
				Type:    "relaypin",
				Hash:    hash,
				PeerId:  peerID,
				PeerKey: peerKey,
				IsRelay: relay,
			}
			var err error
			data, err := json.Marshal(msg)
			if err != nil {
				log.Errorf("msg json marshal:%v", err)
			} else {
				_, err = remotemsg(s, relayKey, string(data))
				if err != nil {
					log.Errorf("call relaypin remote error:%v", err)
				}
			}
			select {
			case out <- err:
			case <-ctx.Done():
			}
		}
	}()

	return out, nil
}

func parsePeerParam(text string) (ma.Multiaddr, peer.ID, error) {
	// to be replaced with just multiaddr parsing, once ptp is a multiaddr protocol
	idx := strings.LastIndex(text, "/")
	if idx == -1 {
		pid, err := peer.IDB58Decode(text)
		if err != nil {
			return nil, "", err
		}

		return nil, pid, nil
	}

	addrS := text[:idx]
	peeridS := text[idx+1:]

	var maddr ma.Multiaddr
	var pid peer.ID

	// make sure addrS parses as a multiaddr.
	if len(addrS) > 0 {
		var err error
		maddr, err = ma.NewMultiaddr(addrS)
		if err != nil {
			return nil, "", err
		}
	}

	// make sure idS parses as a peer.ID
	var err error
	pid, err = peer.IDB58Decode(peeridS)
	if err != nil {
		return nil, "", err
	}

	return maddr, pid, nil
}

func Encrypt(src, key []byte) ([]byte, error) {

	bkey := sha256.Sum256(key)
	block, err := aes.NewCipher(bkey[:])
	if err != nil {
		return nil, err
	}

	// 验证输入参数
	// 必须为aes.Blocksize的倍数
	if len(src)%aes.BlockSize != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}

	encryptText := make([]byte, aes.BlockSize+len(src))

	iv := encryptText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(encryptText[aes.BlockSize:], src)

	return encryptText, nil
}

// AES解密
func Decrypt(src, key []byte) ([]byte, error) {

	bkey := sha256.Sum256(key)
	block, err := aes.NewCipher(bkey[:])
	if err != nil {
		return nil, err
	}

	// hex
	decryptText, err := hex.DecodeString(fmt.Sprintf("%x", string(src)))
	if err != nil {
		return nil, err
	}

	// 长度不能小于aes.Blocksize
	if len(decryptText) < aes.BlockSize {
		return nil, errors.New("crypto/cipher: ciphertext too short")
	}

	iv := decryptText[:aes.BlockSize]
	decryptText = decryptText[aes.BlockSize:]

	// 验证输入参数
	// 必须为aes.Blocksize的倍数
	if len(decryptText)%aes.BlockSize != 0 {
		return nil, errors.New("crypto/cipher: ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(decryptText, decryptText)

	return decryptText, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	if length == 0 {
		return nil
	}
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
