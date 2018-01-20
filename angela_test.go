package angela

import (
	"hash/fnv"
	"testing"
)

func TestNewTree(t *testing.T) {
	tr := NewTree(8192)
	t.Log(tr.Root.Hash(false))
	tr.Insert(testItem{"foo"})
	t.Log(tr.Root.hash)
	if tr.Root.hash != tr.Root.Hash(false) {
		t.Errorf("hash mismatch")
	}
	t.Log(tr.Root.Hash(false))
	i, ok := tr.Get("foo")
	if !ok {
		t.Errorf("expected ok, got !ok")
	}
	if i.ID() != "foo" {
		t.Errorf("expected foo, got %s", i.ID())
	}
	ok = tr.Delete("foo")
	if !ok {
		t.Errorf("expected ok, got !ok")
	}
	_, ok = tr.Get("foo")
	if ok {
		t.Errorf("expected !ok, got !ok")
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
