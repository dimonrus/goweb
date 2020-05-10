package test

import (
	"fmt"
	"github.com/dimonrus/goweb"
	"testing"
)

func TestNewConnectionEventListeners(t *testing.T) {
	ev := goweb.NewConnectionEventListeners()
	if ev == nil {
		t.Fatal("event can ot be nill")
	}
}

func TestConnectionEventListeners_Register(t *testing.T) {
	ev := goweb.NewConnectionEventListeners()
	event := ev.Register("cid", "some_event")
	go func() {
		ev.Dispatch("cid", "some_event")
	}()
	<-event
	fmt.Println("event done")
}

func TestConnectionEventListeners_Unregister(t *testing.T) {
	ev := goweb.NewConnectionEventListeners()
	ev.Register("cid", "some_event")
	ev.Register("cid", "some_other_event")
	ev.Unregister("cid", "some_event")
	events := ev.Get("cid")
	if len(events) != 1 {
		t.Fatal("events must not exists")
	}
}

func TestConnectionEventListeners_Get(t *testing.T) {
	ev := goweb.NewConnectionEventListeners()
	ev.Register("cid", "some_event")
	events := ev.Get("cid")
	if len(events) == 0 {
		t.Fatal("events must exists")
	}
}

func TestConnectionEventListeners_UnregisterConnection(t *testing.T) {
	ev := goweb.NewConnectionEventListeners()
	ev.Register("cid", "some_event")
	ev.UnregisterConnection("cid")
	events := ev.Get("cid")
	if len(events) != 0 {
		t.Fatal("events must not exists")
	}
}

func BenchmarkConnectionEventListeners_UnregisterConnection(b *testing.B) {
	ev := goweb.NewConnectionEventListeners()
	for i := 0; i < b.N; i++ {
		ev.Register("cid", "some_event")
		ev.Unregister("cid", "some_event")
	}
	b.ReportAllocs()
}