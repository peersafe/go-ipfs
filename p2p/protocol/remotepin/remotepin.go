package remotepin

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
	"strings"
	"time"

	ma "github.com/ipfs/go-ipfs/Godeps/_workspace/src/github.com/jbenet/go-multiaddr"
	context "github.com/ipfs/go-ipfs/Godeps/_workspace/src/golang.org/x/net/context"
	iaddr "github.com/ipfs/go-ipfs/util/ipfsaddr"

	localhandle "github.com/ipfs/go-ipfs/core/remotehandle"
	host "github.com/ipfs/go-ipfs/p2p/host"
	inet "github.com/ipfs/go-ipfs/p2p/net"
	peer "github.com/ipfs/go-ipfs/p2p/peer"
	path "github.com/ipfs/go-ipfs/path"
	"github.com/ipfs/go-ipfs/repo/config"
	logging "github.com/ipfs/go-ipfs/vendor/QmQg1J6vikuXF9oDvm4wpdeAUvvkVEKW1EYDw9HhTMnP2b/go-log"
)

var log = logging.Logger("remotepin")

const ID = "/ipfs/remotepin"

type RemotepinService struct {
	Host            host.Host
	LocalHandle     localhandle.Remotepin
	Secret          string
	RemoteMultiplex config.RemoteMultiplex
	currentPin      chan []string
	pinQueue        chan string
}

func NewRemotepinService(h host.Host, handler localhandle.Remotepin, key string, remu config.RemoteMultiplex) *RemotepinService {
	if remu.MaxPin < 1 {
		remu.MaxPin = 1
	}
	ps := &RemotepinService{Host: h, LocalHandle: handler, Secret: key, RemoteMultiplex: remu}
	if remu.Master {
		ps.currentPin = make(chan []string, len(remu.Slave))
		for _, v := range ps.RemoteMultiplex.Slave {
			ps.currentPin <- v
		}
		ps.pinQueue = make(chan string, 100)
		go ps.Run()
	} else {
		ps.currentPin = make(chan []string, remu.MaxPin)
	}
	h.SetStreamHandler(ID, ps.RemotepinHandler)
	return ps
}

func (p *RemotepinService) RemotepinHandler(s inet.Stream) {
	for {
		slen := make([]byte, 1)
		_, err := io.ReadFull(s, slen)
		if err != nil {
			log.Debug(err)
			return
		}

		len := int(slen[0])

		rbuf := make([]byte, len)
		_, err = io.ReadFull(s, rbuf)
		if err != nil {
			log.Debug(err)
			return
		}

		var buf []byte = rbuf
		path, err := p.DecryptRequest(rbuf)
		if err != nil {
			buf = PKCS5Padding([]byte(fmt.Sprint(err)), len)
		}

		err = p.MultiplexRequest(string(path))
		if err != nil {
			buf = PKCS5Padding([]byte(fmt.Sprint(err)), len)
		}

		_, err = s.Write(buf)
		if err != nil {
			log.Debug(err)
			return
		}
	}
}

func (ps *RemotepinService) Run() {
	tryMax := len(ps.currentPin)
	tryTime := time.Duration(ps.RemoteMultiplex.TryTime) * time.Second
	for {
		fpath := <-ps.pinQueue
		var tryNo = 0
		for {
			strPeer := <-ps.currentPin
			tryNo++
			err := ps.delay(strPeer, fpath)
			if err != nil && err.Error() == "max pin" {
				log.Debugf("[%s] delay [%s] max pin", fpath, strPeer)
				if tryNo >= tryMax {
					select {
					case <-time.After(tryTime):
					}
					tryNo = 0
				}
				ps.currentPin <- strPeer
			} else if err != nil {
				ps.currentPin <- strPeer
				log.Debugf("[%s] delay [%s] error [%s]", fpath, strPeer, err)
			} else {
				ps.currentPin <- strPeer
				log.Debugf("[%s] delay [%s] success", fpath, strPeer)
				break
			}
		}
	}
}

func (p *RemotepinService) DecryptRequest(buf []byte) (rbuf []byte, err error) {
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

func (p *RemotepinService) MultiplexRequest(fpath string) error {
	if p.RemoteMultiplex.Master { // local is master
		log.Debug("local is master")
		p.pinQueue <- fpath

	} else { // local is slave
		log.Debug("local is slave")
		select {
		case p.currentPin <- []string{}:
			log.Debugf("current Pin go num [%d]", len(p.currentPin))
		default:
			return errors.New("max pin")
		}

		go func() {
			defer func() {
				<-p.currentPin
			}()
			p.LocalHandle.RemotePin(fpath)
		}()
	}

	return nil
}

func (ps *RemotepinService) delay(peerInfo []string, fpath string) error {
	addr, err := parseAddresses(peerInfo[0])
	if err != nil {
		return err
	}

	ps.Host.Network().ConnsToPeer(addr.ID())

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	key := peerInfo[1]

	out, err := ps.Remotepin(ctx, addr.ID(), key, path.Path(fpath))
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

func (ps *RemotepinService) Remotepin(ctx context.Context, p peer.ID, key string, path path.Path) (<-chan error, error) {
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
			_, err := remotepin(s, key, path)
			select {
			case out <- err:
			case <-ctx.Done():
			}
		}
	}()

	return out, nil
}

func remotepin(s inet.Stream, key string, path path.Path) (time.Duration, error) {
	before := time.Now()
	if !strings.HasPrefix(string(path), "/ipfs/") {
		path = "/ipfs/" + path
	}

	orig := PKCS5Padding([]byte(path), aes.BlockSize)
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

	if !bytes.Equal(buf, rbuf) {
		str := PKCS5UnPadding(rbuf)
		return 0, errors.New(string(str))
	}

	return time.Now().Sub(before), nil
}

func parseAddresses(addr string) (iaddrs iaddr.IPFSAddr, err error) {
	iaddrs, err = iaddr.ParseString(addr)
	if err != nil {
		return nil, fmt.Errorf("invalid peer address: " + err.Error())
	}
	return
}

func peersWithAddresses(addr string) (pis []peer.PeerInfo, err error) {
	iaddr, err := parseAddresses(addr)
	if err != nil {
		return nil, err
	}

	pis = append(pis, peer.PeerInfo{
		ID:    iaddr.ID(),
		Addrs: []ma.Multiaddr{iaddr.Transport()},
	})
	return pis, nil
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
