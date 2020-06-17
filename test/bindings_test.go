package test

import (
	"github.com/dimonrus/goweb"
	"testing"
)

func TestNewConnectionBindings(t *testing.T) {
	b := goweb.NewConnectionBindings()
	if b == nil {
		t.Fatal("must be not nil")
	}
}

func TestConnectionBindings_GetBindingId(t *testing.T) {
	b := goweb.NewConnectionBindings()
	id := b.GetBindingId("bid")
	if id != nil {
		t.Fatal("must be a nil")
	}
}

func TestConnectionBindings_GetConnectionId(t *testing.T) {
	b := goweb.NewConnectionBindings()
	c := b.GetConnectionId("cid")
	if c != "" {
		t.Fatal("must be an empty string")
	}
}

func TestConnectionBindings_Bind(t *testing.T) {
	b := goweb.NewConnectionBindings()
	b = b.Bind("bid", "cid")
	cid := b.GetConnectionId("bid")
	if cid == "" {
		t.Fatal("cid is not eq cid")
	}
}

func TestConnectionBindings_UnBind(t *testing.T) {
	b := goweb.NewConnectionBindings()
	b = b.Bind("bid", "cid")
	bid := b.GetBindingId("cid")
	if bid == nil {
		t.Fatal("cid is not eq cid")
	}
	b.UnBind("bid", "cid")
	cid := b.GetConnectionId("bid")
	if cid != "" {
		t.Fatal("cid must be unbound")
	}
}

func TestConnectionBindings_GetBindingIdentifiers(t *testing.T) {
	b := goweb.NewConnectionBindings()
	b = b.Bind("bid", "cid")
	b = b.Bind("bid2", "cid2")
	ids := b.GetBindingIdentifiers()
	if ids[0] == "bid" || ids[0] == "bid2" {
		return
	}
	t.Fatal("its wrong")
}

func BenchmarkBindings(b *testing.B) {
	bindings := goweb.NewConnectionBindings()
	for i := 0; i < b.N; i++ {
		bindings.Bind("bid", "cid")
		bindings.UnBind("bid", "cid")
	}
	b.ReportAllocs()
}
