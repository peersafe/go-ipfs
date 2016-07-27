package ipfs_lib

import (
	"strings"
	"testing"
)

var (
	rootHash string
)

func TestIpfsInit(t *testing.T) {
	ret, result := IpfsInit()
	t.Logf("IpfsInit result:\n %s\n", result)

	index := strings.LastIndex(result, ":")
	rootHash = strings.Trim(result[index+1:], " ")
	t.Logf("peer identity:%s", rootHash)

	if ret <= 0 {
		t.Errorf("IpfsInit failed!")
	}
}

// func TestIpfsDaemon(t *testing.T) {
// 	go func() {
// 		ret, result := IpfsDaemon()
// 		t.Logf("IpfsDaemon result:\n %s\n", result)
// 		if ret <= 0 {
// 			t.Errorf("IpfsDaemon failed!")
// 		}
// 	}()
// }

func TestIpfsId(t *testing.T) {
	ret, result := IpfsId(0)
	t.Logf("IpfsId result:\n %s\n", result)
	if ret <= 0 {
		t.Errorf("IpfsId failed!")
	}
}

// func TestIpfaAdd(t *testing.T) {
// 	ret, result := IpfsAdd()
// 	t.Logf("IpfsAdd result:\n %s\n", result)
// 	if ret <= 0 {
// 		t.Errorf("IpfsAdd failed!")
// 	}
// }
