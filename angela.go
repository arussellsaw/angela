package angela

import (
	"encoding/binary"
	"hash/fnv"
	"log"
	"math/bits"

	"github.com/google/uuid"
)

func NewTree(buckets uint32) *Tree {
	if (buckets & (buckets - 1)) != 0 {
		panic("number of branches must be base 2")
	}
	n := bits.TrailingZeros32(buckets)
	r := NewBranch(buckets)
	t := &Tree{Root: r}
	t.populate(r, uint32(n), 0)
	t.Root.Hash(false)
	return t
}

type Tree struct {
	Root          *Branch
	finalBranches map[string]*Branch
}

func (t *Tree) Insert(i Item) {
	t.Root.Insert(i)
}

func (t *Tree) Get(id string) (Item, bool) {
	return t.Root.Get(id)
}

func (t *Tree) Delete(id string) bool {
	return t.Root.Delete(id)
}

func (t *Tree) Scan(fn func(Item) error) error {
	for i := range t.finalBranches {
		for _, v := range t.finalBranches[i].items {
			if err := fn(v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Tree) populate(b *Branch, n, depth uint32) {
	t.finalBranches = make(map[string]*Branch)
	for i := range b.b {
		b.b[i] = NewBranch(b.buckets / 2)
		if n == depth {
			b.b[i].final = true
			t.finalBranches[b.b[i].id] = b.b[i]
			continue
		}
		t.populate(b.b[i], n, depth+1)
	}
}

func NewBranch(buckets uint32) *Branch {
	return &Branch{buckets: buckets, id: uuid.New().String()}
}

type Branch struct {
	id      string
	buckets uint32
	hash    uint32
	b       [2]*Branch
	final   bool
	items   map[string]Item
}

func (b *Branch) Insert(i Item) {
	h := fnv.New32()
	h.Write([]byte(i.ID()))
	b.doInsert(i, h.Sum32()%b.buckets, uint32(31-bits.LeadingZeros32(b.buckets)))
}

func (b *Branch) doInsert(i Item, n, depth uint32) {
	if b.final {
		if b.items == nil {
			b.items = make(map[string]Item)
		}
		log.Printf("branch %s inserting %s %v", b.id, i.ID(), i.Hash())
		b.items[i.ID()] = i
		b.Hash(true)
		return
	}
	depth = depth - 1
	if (n>>depth)&1 == 1 {
		b.b[1].doInsert(i, n, depth)
		b.Hash(true)
		return
	}
	b.b[0].doInsert(i, n, depth)
	b.Hash(true)
}

func (b *Branch) Get(id string) (Item, bool) {
	h := fnv.New32()
	h.Write([]byte(id))
	return b.doGet(id, h.Sum32()%b.buckets, uint32(31-bits.LeadingZeros32(b.buckets)))
}

func (b *Branch) doGet(id string, n, depth uint32) (Item, bool) {
	if b.final {
		if b.items == nil {
			return nil, false
		}
		i, ok := b.items[id]
		return i, ok
	}
	depth = depth - 1
	if (n>>depth)&1 == 1 {
		return b.b[1].doGet(id, n, depth)
	}
	return b.b[0].doGet(id, n, depth)
}

func (b *Branch) Delete(id string) bool {
	h := fnv.New32()
	h.Write([]byte(id))
	return b.doDelete(id, h.Sum32(), uint32(31-bits.LeadingZeros32(b.buckets)))
}

func (b *Branch) doDelete(id string, n, depth uint32) bool {
	if b.final {
		if b.items == nil {
			return false
		}
		delete(b.items, id)
		b.Hash(true)
		return true
	}
	depth = depth - 1
	if (n>>depth)&1 == 1 {
		ok := b.b[1].doDelete(id, n, depth)
		if ok {
			b.Hash(true)
		}
		return ok
	}
	ok := b.b[0].doDelete(id, n, depth)
	if ok {
		b.Hash(true)
	}
	return ok
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
	log.Print(b.hash, b.buckets)
	return b.hash
}

type Item interface {
	ID() string
	Hash() uint32
}
