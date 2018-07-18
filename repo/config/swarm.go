package config

type SwarmConfig struct {
	AddrFilters             []string
	DisableBandwidthMetrics bool
	DisableNatPortMap       bool
	DisableRelay            bool
	EnableRelayHop          bool

	ConnMgr ConnMgr
	RedMgr  RedMgr
}

// ConnMgr defines configuration options for the libp2p connection manager
type ConnMgr struct {
	Type        string
	LowWater    int
	HighWater   int
	GracePeriod string
}

// RedMgr defines configuration options for the block data redundancy manager
type RedMgr struct {
	Type   string
	RedNum int
}
