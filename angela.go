package angela

import (
	"encoding/binary"
	"hash/fnv"

	"github.com/google/uuid"
)

func NewTree(buckets uint64) *Tree {
	i := buckets
	if (i & (i - 1)) != 0 {
		panic("number of branches must be base 2")
	}
	// tree depth for a base 2 number = setBits(n-1)
	i = i - 1
	i = i - ((i >> 1) & 0x5555555555555555)
	i = (i & 0x3333333333333333) + ((i >> 2) & 0x3333333333333333)
	i = ((i + (i >> 4)) & 0x0F0F0F0F0F0F0F0F)
	n := (i * (0x0101010101010101)) >> 56
	r := NewBranch(0, 128)
	t := &Tree{Root: r, buckets: 128}
	t.populate(r, n-1, 0)
	final := 0
	for i := range t.branches {
		if t.branches[i].final {
			final++
		}
	}
	t.Root.Hash(false)
	return t
}

type Tree struct {
	Root          *Branch
	branches      map[string]*Branch
	depth         uint64
	buckets       uint32
	finalBranches []*Branch
}

func (t *Tree) Insert(i Item) {
	t.Root.insert(i, t.buckets)
}

func (t *Tree) populate(b *Branch, n, depth uint64) {
	if t.branches == nil {
		t.branches = make(map[string]*Branch)
	}
	start := b.start
	bucketRange := (b.end - b.start) / 2
	for i := range b.b {
		b.b[i] = NewBranch(start, start+bucketRange)
		start += bucketRange
		t.branches[b.b[i].id] = b.b[i]
		if n == depth {
			t.finalBranches = append(t.finalBranches, b.b[i])
			b.b[i].final = true
			continue
		}
		t.populate(b.b[i], n, depth+1)
	}
}

func NewBranch(start, end uint32) *Branch {
	return &Branch{start: start, end: end, id: uuid.New().String()}
}

type Branch struct {
	id         string
	hash       uint32
	b          [2]*Branch
	final      bool
	start, end uint32
	items      map[string]Item
}

func (b *Branch) insert(i Item, n uint32) bool {
	if b.final {
		if b.items == nil {
			b.items = make(map[string]Item)
		}
		b.items[i.ID()] = i
		b.Hash(true)
		return true
	}
	bn := i.Hash() % n
	if !(bn >= b.start && bn <= b.end) {
		return false
	}
	for j := range b.b {
		if b.b[j].insert(i, n) {
			b.Hash(true)
			return true
		}
	}
	panic("lost in the tree!")
	return false
}

func (b *Branch) Hash(shallow bool) uint32 {
	h := fnv.New32()
	buf := make([]byte, 4)
	if b.final {
		for i := range b.items {
			binary.LittleEndian.PutUint32(buf, b.items[i].Hash())
			h.Write(buf)
		}
		b.hash = h.Sum32()
		return b.hash
	}
	bh := b.b[0].hash
	if !shallow {
		bh = b.b[0].Hash(false)
	}
	binary.LittleEndian.PutUint32(buf, bh)
	h.Write(buf)
	bh = b.b[1].hash
	if !shallow {
		bh = b.b[1].Hash(false)
	}
	binary.LittleEndian.PutUint32(buf, bh)
	h.Write(buf)
	b.hash = h.Sum32()
	return b.hash
}

type Item interface {
	ID() string
	Hash() uint32
}
