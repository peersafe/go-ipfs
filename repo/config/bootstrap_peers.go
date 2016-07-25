package config

import (
	"errors"
	"fmt"

	iaddr "github.com/ipfs/go-ipfs/util/ipfsaddr"
)

// DefaultBootstrapAddresses are the hardcoded bootstrap addresses
// for ipfs. they are nodes run by the ipfs team. docs on these later.
// As with all p2p networks, bootstrap is an important security concern.
//
// Note: this is here -- and not inside cmd/ipfs/init.go -- because of an
// import dependency issue. TODO: move this into a config/default/ package.
var DefaultBootstrapAddresses = []string{
	"/ip4/101.201.40.124/tcp/4001/ipfs/QmZDYAhmMDtnoC6XZRw8R1swgoshxKvXDA9oQF97AYkPZc",
	"/ip4/108.61.161.202/tcp/4001/ipfs/QmaXZQYJeFTyytgdMpguU3CFtTj82EX6CEYyNoRakFqqkJ",
	"/ip4/139.129.99.7/tcp/4001/ipfs/QmcwRr8wfdtQTru2R3L2cx6BZBxdVpBgRfEULQQRmwRJtK",
	"/ip4/115.159.105.185/tcp/4001/ipfs/QmPkFbxAQ7DeKD5VGSh9HQrdS574pyNzDmxJeGrRJxoucF",
	"/ip4/119.29.67.136/tcp/4001/ipfs/QmTGkgHSsULk8p3AKTAqKixxidZQXFyF7mCURcutPqrwjQ",
	"/ip4/45.32.70.172/tcp/4001/ipfs/QmZYf9EUMYk9V2DTBpoJRmSotFRPxQPdozj8CXtYPvXDyU",
	"/ip4/101.201.220.73/tcp/4001/ipfs/QmeVGtbVRrz4m6ioPQfqimgiSQFpiGibF5tmbsxGW95Gdm",
	"/ip4/219.223.222.4/tcp/4001/ipfs/Qmf96ojxn2i8QPZ83FbutnwGjffEXsV4VaoFGzuC3YEwwY",
}

// BootstrapPeer is a peer used to bootstrap the network.
type BootstrapPeer iaddr.IPFSAddr

// ErrInvalidPeerAddr signals an address is not a valid peer address.
var ErrInvalidPeerAddr = errors.New("invalid peer address")

func (c *Config) BootstrapPeers() ([]BootstrapPeer, error) {
	return ParseBootstrapPeers(c.Bootstrap)
}

// DefaultBootstrapPeers returns the (parsed) set of default bootstrap peers.
// if it fails, it returns a meaningful error for the user.
// This is here (and not inside cmd/ipfs/init) because of module dependency problems.
func DefaultBootstrapPeers() ([]BootstrapPeer, error) {
	ps, err := ParseBootstrapPeers(DefaultBootstrapAddresses)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse hardcoded bootstrap peers: %s
This is a problem with the ipfs codebase. Please report it to the dev team.`, err)
	}
	return ps, nil
}

func (c *Config) SetBootstrapPeers(bps []BootstrapPeer) {
	c.Bootstrap = BootstrapPeerStrings(bps)
}

func ParseBootstrapPeer(addr string) (BootstrapPeer, error) {
	ia, err := iaddr.ParseString(addr)
	if err != nil {
		return nil, err
	}
	return BootstrapPeer(ia), err
}

func ParseBootstrapPeers(addrs []string) ([]BootstrapPeer, error) {
	peers := make([]BootstrapPeer, len(addrs))
	var err error
	for i, addr := range addrs {
		peers[i], err = ParseBootstrapPeer(addr)
		if err != nil {
			return nil, err
		}
	}
	return peers, nil
}

func BootstrapPeerStrings(bps []BootstrapPeer) []string {
	bpss := make([]string, len(bps))
	for i, p := range bps {
		bpss[i] = p.String()
	}
	return bpss
}
