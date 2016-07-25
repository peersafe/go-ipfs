package remotels

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

	logging "QmQg1J6vikuXF9oDvm4wpdeAUvvkVEKW1EYDw9HhTMnP2b/go-log"

	host "github.com/ipfs/go-ipfs/p2p/host"
	inet "github.com/ipfs/go-ipfs/p2p/net"
	peer "github.com/ipfs/go-ipfs/p2p/peer"
	path "github.com/ipfs/go-ipfs/path"
)

var log = logging.Logger("remotels")

const ID = "/ipfs/remotels"

type RemotelsService struct {
	Host   host.Host
	Secret string
}

func NewRemotelsService(h host.Host, key string) *RemotelsService {
	ps := &RemotelsService{h, key}
	h.SetStreamHandler(ID, ps.RemotelsHandler)
	return ps
}

func (p *RemotelsService) RemotelsHandler(s inet.Stream) {
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

		err = p.remoteLs(string(path))
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

func (p *RemotelsService) DecryptRequest(buf []byte) (rbuf []byte, err error) {
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

func (ps *RemotelsService) remoteLs(fpath string) error {
	go func() {
		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		cmd := exec.Cmd{
			Path: path,
			Args: []string{"ipfs", "ls", fpath},
		}
		err := cmd.Run()
		if err != nil {
			log.Debug(err)
		}

	}()
	return nil
}

func (ps *RemotelsService) Remotels(ctx context.Context, p peer.ID, key string, path path.Path) (<-chan error, error) {
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
			_, err := remotels(s, key, path)
			if err != nil {
				log.Debugf("remotels error: %s", err)
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

func remotels(s inet.Stream, key string, path path.Path) (time.Duration, error) {
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
