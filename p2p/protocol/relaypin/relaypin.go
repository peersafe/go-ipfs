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

	context "github.com/ipfs/go-ipfs/Godeps/_workspace/src/golang.org/x/net/context"

	ma "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-multiaddr"
	host "github.com/ipfs/go-ipfs/p2p/host"
	inet "github.com/ipfs/go-ipfs/p2p/net"
	peer "github.com/ipfs/go-ipfs/p2p/peer"
	path "github.com/ipfs/go-ipfs/path"
	logging "github.com/ipfs/go-ipfs/vendor/QmQg1J6vikuXF9oDvm4wpdeAUvvkVEKW1EYDw9HhTMnP2b/go-log"
)

var log = logging.Logger("relaypin")

const ID = "/ipfs/relaypin"

type RelaypinService struct {
	Host   host.Host
	Secret string
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
	log.Debug(str)
	sz := strings.Split(str, "|")
	m.delay = (sz[0] == "true")
	m.peer = sz[1]
	m.key = sz[2]
	m.path = sz[3]
	log.Debugf("[%v][%v][%v][%v]", m.delay, m.peer, m.key, m.path)
}

func NewRelaypinService(h host.Host, key string) *RelaypinService {
	ps := &RelaypinService{Host: h, Secret: key}
	h.SetStreamHandler(ID, ps.RelaypinService)
	return ps
}

func (p *RelaypinService) RelaypinService(s inet.Stream) {
	for {
		slen := make([]byte, 1)
		_, err := io.ReadFull(s, slen)
		if err != nil {
			log.Debug(err)
			return
		}

		blen := int(slen[0])

		rbuf := make([]byte, blen)
		_, err = io.ReadFull(s, rbuf)
		if err != nil {
			log.Debug(err)
			return
		}
		log.Debugf(">>>>>>>>>>>>>> recv len=[%d] rbuf->[%v]", len(rbuf), rbuf)

		var buf []byte = rbuf
		msgbuf, err := p.decryptRequest(rbuf)
		if err != nil {
			buf = PKCS5Padding([]byte(fmt.Sprint(err)), blen)
		}

		log.Debug(">>>>>>>>>>>>>> recv msgbuf->", msgbuf)

		msg := p.parseMsg(msgbuf)

		err = p.relayRequest(msg)
		if err != nil {
			log.Debugf(">>>>>>relayRequest is error error=[%v]len=[%d]", err, blen)
			buf = PKCS5Padding([]byte(fmt.Sprint(err)), blen)
		}

		_, err = s.Write(buf)
		if err != nil {
			log.Debug(err)
			return
		}
	}
}

func (p *RelaypinService) decryptRequest(buf []byte) (rbuf []byte, err error) {
	md5hash_buf := buf[:md5.Size]
	crypted := buf[md5.Size:]

	orig, err := Decrypt(crypted, []byte(p.Secret))
	if err != nil {
		log.Debug(err)
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

		err := p.relayPeer(msg.peer, msg.key, msg.path)
		log.Debugf(">>>>>>delay request err=[%v]", err)
		return err

	} else { // local is slave
		log.Debug("local is work")

		go func() {
			file, _ := exec.LookPath(os.Args[0])
			path, _ := filepath.Abs(file)
			log.Debugf("path [%s]", path)
			cmd := exec.Cmd{
				Path: path,
				Args: []string{"ipfs", "pin", "add", msg.path},
			}
			err := cmd.Run()
			if err != nil {
				log.Debug(err)
			}
		}()
	}

	return nil
}

func (ps *RelaypinService) relayPeer(p, key, fpath string) error {
	log.Debugf(">>>>>>>>>relayPeer p[%v]", p)
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
		log.Debug(">>>>>>>>> Relaypin return reerr=", reerr)
		if reerr != nil {
			return reerr
		}
	case <-ctx.Done():
		return fmt.Errorf("routing timeout")
	}
	return nil
}

func (ps *RelaypinService) Relaypin(ctx context.Context, p peer.ID, key, peer, peerkey string, path path.Path, delay bool) (<-chan error, error) {
	log.Debugf(">>>>>>>>[%v][%v][%v][%v][%v]", p, key, peer, peerkey, path)
	s, err := ps.Host.NewStream(ID, p)
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
			log.Debug(">>>>>>>>>> call relaypin")
			_, err := relaypin(s, key, peer, peerkey, path, delay)
			select {
			case out <- err:
				log.Debug(">>>>>>>>>> call relaypin call back err=", err)
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
	log.Debugf(">>>>>>>[%v]", msg)

	orig := PKCS5Padding(msg.toByte(), aes.BlockSize)
	md5hash := md5.Sum(orig)
	crypted, err := Encrypt(orig, []byte(key))
	if err != nil {
		return 0, err
	}
	buf := append(md5hash[:], crypted...)
	buflen := len(buf)

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

	log.Debugf(">>>>>rbuf[%v][%d]", rbuf, len(rbuf))
	if !bytes.Equal(buf, rbuf) {
		str := PKCS5UnPadding(rbuf)
		log.Debugf(">>>>>relaypin error")
		return 0, errors.New(string(str))
	}

	log.Debugf(">>>>>relaypin success")
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
