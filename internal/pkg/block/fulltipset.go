package block

import (
	"github.com/ipfs/go-cid"
)

// FullTipSet is an expanded version of the TipSet that contains all the blocks and messages
type FullTipSet struct {
	Blocks []*FullBlock
	tipset *TipSet
	cids   []cid.Cid
}

func NewFullTipSet(blks []*FullBlock) *FullTipSet {
	return &FullTipSet{
		Blocks: blks,
	}
}

func (fts *FullTipSet) Cids() []cid.Cid {
	if fts.cids != nil {
		return fts.cids
	}

	var cids []cid.Cid
	for _, b := range fts.Blocks {
		cids = append(cids, b.Cid())
	}
	fts.cids = cids

	return cids
}

// TipSet returns a narrower view of this FullTipSet elliding the block
// messages.
func (fts *FullTipSet) TipSet() *TipSet {
	if fts.tipset != nil {
		// FIXME: fts.tipset is actually never set. Should it memoize?
		return fts.tipset
	}

	var headers []*Block
	for _, b := range fts.Blocks {
		headers = append(headers, b.Header)
	}

	ts, err := NewTipSet(headers...)
	if err != nil {
		panic(err)
	}

	return ts
}