package angela

import (
	"hash/fnv"
	"testing"
)

func TestNewTree(t *testing.T) {
	tr := NewTree(128)
	t.Log(tr.Root.Hash(false))
	tr.Insert(testItem{"foo"})
	t.Log(tr.Root.hash)
	if tr.Root.hash != tr.Root.Hash(false) {
		t.Errorf("hash mismatch")
	}
	t.Log(tr.Root.Hash(false))
	var found bool
	for _, b := range tr.finalBranches {
		if _, ok := b.items["foo"]; ok {
			found = true
		}
	}
	if !found {
		t.Errorf("expected found, got !found")
	}
}

type testItem struct {
	id string
}

func (i testItem) ID() string {
	return i.id
}

func (i testItem) Hash() uint32 {
	h := fnv.New32()
	h.Write([]byte(i.id))
	return h.Sum32()
}
