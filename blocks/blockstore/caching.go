package blockstore

import (
	"errors"

	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"

	cid "gx/ipfs/QmfSc2xehWmWLnwwYR91Y8QF4xdASypTFVknutoKQS3GHp/go-cid"
)

// Next to each option is it aproximate memory usage per unit
type CacheOpts struct {
	HasBloomFilterSize   int // 1 bit
	HasBloomFilterHashes int // No size, 7 is usually best, consult bloom papers
	HasARCCacheSize      int // 32 bytes
}

func DefaultCacheOpts() CacheOpts {
	return CacheOpts{
		HasBloomFilterSize:   512 * 8 * 1024, // 512K Memory
		HasBloomFilterHashes: 7,
		HasARCCacheSize:      64 * 1024, // 2M Memory 2G Disk
	}
}

func CachedBlockstore(bs GCBlockstore,
	ctx context.Context, opts CacheOpts) (cbs GCBlockstore, cids chan *cid.Cid, err error) {
	cbs = bs

	if opts.HasBloomFilterSize < 0 || opts.HasBloomFilterHashes < 0 ||
		opts.HasARCCacheSize < 0 {
		return nil, nil, errors.New("all options for cache need to be greater than zero")
	}

	if opts.HasBloomFilterSize != 0 && opts.HasBloomFilterHashes == 0 {
		return nil, nil, errors.New("bloom filter hash count can't be 0 when there is size set")
	}
	if opts.HasBloomFilterSize != 0 {
		cbs, err = bloomCached(cbs, ctx, opts.HasBloomFilterSize, opts.HasBloomFilterHashes)
	}
	if opts.HasARCCacheSize > 0 {
		cbs, cids, err = arcCached(cbs, opts.HasARCCacheSize)
	}
	return cbs, cids, err
}
