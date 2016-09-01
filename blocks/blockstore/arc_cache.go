package blockstore

import (
	"github.com/ipfs/go-ipfs/blocks"
	key "github.com/ipfs/go-ipfs/blocks/key"
	ds "gx/ipfs/QmNgqJarToRiq2GBaPJhkmW4B5BxS5B74E1rkGvv2JoaTp/go-datastore"
	lru "gx/ipfs/QmVYxfoJQiZijTgPNHCHgHELvQpbsJNTg6Crmc3dQkj3yy/golang-lru"
	context "gx/ipfs/QmZy2y8t9zQH2a1b8q2ZSLKp17ATuJoCNxxyMFG5qFExpt/go-net/context"

)

const (
	RemoveChanSize = 1024
)

type arccache struct {
	arc        *lru.ARCCache
	blockstore Blockstore
	removeKeys chan key.Key
}

func arcCached(bs Blockstore, lruSize int) (*arccache, chan key.Key, error) {
	arc, err := lru.NewARC(lruSize)
	if err != nil {
		return nil, nil, err
	}

	removeKeys := make(chan key.Key, RemoveChanSize)
	return &arccache{arc: arc, blockstore: bs, removeKeys: removeKeys}, removeKeys, nil
}

func (b *arccache) DeleteBlock(k key.Key) error {
	if has, ok := b.hasCached(k); ok && !has {
		return ErrNotFound
	}

	b.arc.Remove(k) // Invalidate cache before deleting.
	err := b.blockstore.DeleteBlock(k)
	switch err {
	case nil, ds.ErrNotFound, ErrNotFound:
		b.arc.Add(k, false)
		return err
	default:
		return err
	}
	return nil
}

// if ok == false has is inconclusive
// if ok == true then has respons to question: is it contained
func (b *arccache) hasCached(k key.Key) (has bool, ok bool) {
	if k == "" {
		// Return cache invalid so the call to blockstore happens
		// in case of invalid key and correct error is created.
		return false, false
	}

	h, ok := b.arc.Get(k)
	if ok {
		return h.(bool), true
	}
	return false, false
}

func (b *arccache) Has(k key.Key) (bool, error) {
	if has, ok := b.hasCached(k); ok {
		return has, nil
	}

	res, err := b.blockstore.Has(k)
	if res {
		removeKeys := b.arc.Add(k, res)
		b.handleRemoveKeys(removeKeys)
	}
	return res, err
}

func (b *arccache) Get(k key.Key) (blocks.Block, error) {
	if has, ok := b.hasCached(k); ok && !has {
		return nil, ErrNotFound
	}

	var removeKeys []interface{}
	bl, err := b.blockstore.Get(k)
	if bl == nil && err == ErrNotFound {
		removeKeys = b.arc.Add(k, false)
	} else if bl != nil {
		removeKeys = b.arc.Add(k, true)
	}
	b.handleRemoveKeys(removeKeys)

	return bl, err
}

func (b *arccache) Put(bl blocks.Block) error {
	if has, ok := b.hasCached(bl.Key()); ok && has {
		return nil
	}

	err := b.blockstore.Put(bl)
	if err == nil {
		removeKeys := b.arc.Add(bl.Key(), true)
		b.handleRemoveKeys(removeKeys)
	}
	return err
}

func (b *arccache) PutMany(bs []blocks.Block) error {
	var good []blocks.Block
	for _, block := range bs {
		if has, ok := b.hasCached(block.Key()); !ok || (ok && !has) {
			good = append(good, block)
		}
	}
	err := b.blockstore.PutMany(bs)
	if err != nil {
		return err
	}
	for _, block := range bs {
		removeKeys := b.arc.Add(block.Key(), true)
		b.handleRemoveKeys(removeKeys)
	}
	return nil
}

func (b *arccache) AllKeysChan(ctx context.Context) (<-chan key.Key, error) {
	return b.blockstore.AllKeysChan(ctx)
}

func (b *arccache) GCLock() Unlocker {
	return b.blockstore.(GCBlockstore).GCLock()
}

func (b *arccache) PinLock() Unlocker {
	return b.blockstore.(GCBlockstore).PinLock()
}

func (b *arccache) GCRequested() bool {
	return b.blockstore.(GCBlockstore).GCRequested()
}

func (b *arccache) handleRemoveKeys(removeKeys []interface{}) {
	for _, v := range removeKeys {
		if v != nil {
			removeKey := v.(key.Key)
			if removeKey.String() != "" {
				b.removeKeys <- removeKey
			}
		}
	}
}
