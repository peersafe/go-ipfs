package relaypin

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	logging "gx/ipfs/QmSpJByNKFX1sCsHBEp3R73FL4NF6FnQTEGyNAXHm2GS52/go-log"
	peer "gx/ipfs/QmfMmLGoKzCHDN7cGgk64PJr4iipzidDRME8HABSJqvmhC/go-libp2p-peer"
	ma "gx/ipfs/QmUAQaWbKxGCUTuoQVvvicbQNZ9APF5pDGWyAZSe93AtKH/go-multiaddr"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"
	p2phost "gx/ipfs/QmdML3R42PRSwnt46jSuEts9bHSqLctVYEjJqMR3UYV8ki/go-libp2p-host"
	inet "gx/ipfs/QmdXimY9QHaasZmw6hWojWnCJvfgxETjZQfg9g6ZrA9wMX/go-libp2p-net"

	api "github.com/ipfs/go-ipfs/cmd/ipfs_lib/apiinterface"

	path "github.com/ipfs/go-ipfs/path"
)

var log = logging.Logger("relaypin")

const ID = "/ipfs/relaypin"

type RelaypinService struct {
	Host   p2phost.Host
	Secret string
	ApiCmd api.Apier
}

type RelaypinMsg struct {
	delay bool
	peer  string
	key   string
	path  string
}

func (m *RelaypinMsg) toByte() []byte {
	str := fmt.Sprintf("%v|%v|%v|%v", m.delay, m.peer, m.key, m.path)
	return []byte(str)
}

func (m *RelaypinMsg) toMsg(buf []byte) {
	str := string(buf)
	sz := strings.Split(str, "|")
	m.delay = (sz[0] == "true")
	m.peer = sz[1]
	m.key = sz[2]
	m.path = sz[3]
}

func NewRelaypinService(h p2phost.Host, key string) *RelaypinService {
	ps := &RelaypinService{Host: h, Secret: key, ApiCmd: api.GApiInterface}
	h.SetStreamHandler(ID, ps.RelaypinService)
	return ps
}

func (p *RelaypinService) RelaypinService(s inet.Stream) {
	for {
		slen := make([]byte, 1)
		_, err := io.ReadFull(s, slen)
		if err != nil {
			log.Errorf("relaypin error:%v", err)
			return
		}

		blen := int(slen[0])

		rbuf := make([]byte, blen)
		_, err = io.ReadFull(s, rbuf)
		if err != nil {
			log.Errorf("relaypin error:%v", err)
			return
		}

		log.Debugf("relaypinService recv len=[%d] reallen=[%d] rbuf=[%v]", blen, len(rbuf), rbuf)

		var buf []byte = rbuf
		msgbuf, err := p.decryptRequest(rbuf)
		if err != nil {
			log.Errorf("relaypinService decryptRequest error:%v", err)
			buf = PKCS5Padding([]byte(fmt.Sprint(err)), blen)
		} else {
			msg := p.parseMsg(msgbuf)
			err = p.relayRequest(msg)
			if err != nil {
				buf = PKCS5Padding([]byte(fmt.Sprint(err)), blen)
			}
		}

		_, err = s.Write(buf)
		if err != nil {
			log.Errorf("relaypin error:%v", err)
			return
		}
	}
}

func (p *RelaypinService) decryptRequest(buf []byte) (rbuf []byte, err error) {
	md5hash_buf := buf[:md5.Size]
	crypted := buf[md5.Size:]

	orig, err := Decrypt(crypted, []byte(p.Secret))
	if err != nil {
		log.Errorf("relaypin error:%v", err)
		return nil, err
	}

	md5hash := md5.Sum(orig)
	if !bytes.Equal(md5hash[:], md5hash_buf) {
		log.Debug("Secret authentication failed")
		return nil, fmt.Errorf("Secret authentication failed")
	}

	rbuf = PKCS5UnPadding(orig)
	return rbuf, nil
}

func (p *RelaypinService) parseMsg(buf []byte) RelaypinMsg {
	msg := RelaypinMsg{}
	msg.toMsg(buf)
	return msg
}

func (p *RelaypinService) relayRequest(msg RelaypinMsg) error {
	if msg.delay { // local is master
		log.Debug("local is delay")

		go func() {
			file, _ := exec.LookPath(os.Args[0])
			path, _ := filepath.Abs(file)
			cmd := exec.Cmd{
				Path: path,
				Args: []string{"ipfs", "ls", msg.path},
			}
			err := cmd.Run()
			if err != nil {
				log.Errorf("relaypin ls work error:%v", err)
			}
		}()

		err := p.RelayPeer(msg.peer, msg.key, msg.path)
		if err != nil {
			log.Errorf("relay request errer:%v", err)
		}

		return err

	} else { // local is slave
		log.Debug("local is work")

		go func() {
			file, _ := exec.LookPath(os.Args[0])
			app := filepath.Clean(file)
			app = filepath.ToSlash(app)
			app = filepath.Base(app)
			log.Debugf("relayRequest local work file[%v]", file)
			if strings.HasSuffix(app, "ipfs") || strings.HasSuffix(app, "ipfs.exe") {
				// use ipfs
				path, _ := filepath.Abs(file)
				cmd := exec.Cmd{
					Path: path,
					Args: []string{"ipfs", "pin", "add", msg.path},
				}
				err := cmd.Run()
				if err != nil {
					log.Errorf("relaypin local work error:%v", err)
				}
			} else {
				// use libipfs
				_, _, err := p.ApiCmd.Cmd(strings.Join([]string{"ipfs", "pin", "add", msg.path}, "&X&"), 0)
				if err != nil {
					log.Errorf("relaypin local work error:%v", err)
				}
			}
		}()
	}

	return nil
}

func (ps *RelaypinService) RelayPeer(p, key, fpath string) error {
	log.Debugf("relayPeer remote[%v]", p)

	_, pid, err := parsePeerParam(p)
	if err != nil {
		return fmt.Errorf("peer addr hash format error")
	}

	ps.Host.Network().ConnsToPeer(pid)

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	out, err := ps.Relaypin(ctx, pid, key, p, key, path.Path(fpath), false)
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

func (ps *RelaypinService) Relaypin(ctx context.Context, p peer.ID, key, peer, peerkey string, path path.Path, delay bool) (<-chan error, error) {
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
			log.Debugf("[%v][%v][%v][%v][%v]", p, key, peer, peerkey, path)
			_, err := relaypin(s, key, peer, peerkey, path, delay)
			if err != nil {
				log.Errorf("call relaypin remote error:%v", err)
			}
			select {
			case out <- err:
			case <-ctx.Done():
			}
		}
	}()

	return out, nil
}

func relaypin(s inet.Stream, key, peer, peerkey string, path path.Path, relay bool) (time.Duration, error) {
	before := time.Now()
	if !strings.HasPrefix(string(path), "/ipfs/") {
		path = "/ipfs/" + path
	}

	msg := RelaypinMsg{relay, peer, peerkey, string(path)}

	orig := PKCS5Padding(msg.toByte(), aes.BlockSize)
	md5hash := md5.Sum(orig)
	crypted, err := Encrypt(orig, []byte(key))
	if err != nil {
		return 0, err
	}
	buf := append(md5hash[:], crypted...)
	buflen := len(buf)

	log.Debugf("relaypin msg len[%d] buf[%v]", buflen, buf)

	slen := make([]byte, 1)
	slen[0] = byte(buflen)
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

	log.Debugf("relaypin msg recv len[%d] rbuf[%v]", len(rbuf), rbuf)

	if !bytes.Equal(buf, rbuf) {
		str := PKCS5UnPadding(rbuf)
		return 0, errors.New(string(str))
	}

	return time.Now().Sub(before), nil
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
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
