package config

type Discovery struct {
	MDNS MDNS
}

type MDNS struct {
	Enabled bool

	// Time in seconds between discovery rounds
	Interval   int
	ServiceTag string
}

const DefaultMDNSServiceTag = "_peersafe.ipfs-discovery._udp"
