package test

import (
	"github.com/dimonrus/goweb"
	"net"
	"testing"
)

func TestNewConnections(t *testing.T) {
	c := goweb.NewConnections()
	if c == nil {
		t.Fatal("must be not nill")
	}
}

func TestConnections_Get(t *testing.T) {
	c := goweb.NewConnections()
	cid := c.Get("cid")
	if cid != nil {
		t.Fatal("must be nil")
	}
}

func TestConnections_Len(t *testing.T) {
	c := goweb.NewConnections()
	if c.Len() != 0 {
		t.Fatal("must be 0")
	}
}

func TestConnections_Set(t *testing.T) {
	c := goweb.NewConnections()
	con := net.TCPConn{}
	c.Set("cid", &con)
	cn := c.Get("cid")
	if cn == nil {
		t.Fatal("cn must not be a nil")
	}
	cn = c.Get("cid0")
	if cn != nil {
		t.Fatal("cn must be a nil")
	}
}

func TestConnections_Unset(t *testing.T) {
	c := goweb.NewConnections()
	con := net.TCPConn{}
	c.Set("cid", &con)
	c.Unset("cid")
	cn := c.Get("cid")
	if cn != nil {
		t.Fatal("cn must be a nil")
	}
}

func TestConnections_GetIdentifiers(t *testing.T) {
	c := goweb.NewConnections()
	con := net.TCPConn{}
	c.Set("cid", &con)
	c.Set("cida", &con)

	keys := c.GetIdentifiers()
	if len(keys) != 2 {
		t.Fatal("wrong ids")
	}
}

func BenchmarkConnections_Unset(b *testing.B) {
	c := goweb.NewConnections()
	con := net.TCPConn{}
	for i := 0; i < b.N; i++ {
		c.Set("cid", &con)
		c.Unset("cid")
	}
	b.ReportAllocs()
}
