package test

import (
	"fmt"
	"github.com/dimonrus/goweb"
	"net/http"
	"testing"
	"time"
)

func TestNewXTransactionInfo(t *testing.T) {
	tx := goweb.NewXTransactionInfo("someID", 4, time.Now())
	if tx.GetTTL() != 4 {
		t.Fatal("ttl must be 4")
	}
}

func TestExtractXTransactionInfo(t *testing.T) {
	r := &http.Request{Header: make(http.Header)}
	r.Header.Add(goweb.HeaderXTransactionId, "someID")
	r.Header.Add(goweb.HeaderXTransactionTTL, "4")
	r.Header.Add(goweb.HeaderXTransactionTime, time.Now().Format(time.RFC3339Nano))
	info, e := goweb.ExtractXTransactionInfo(r)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(info.GetId(), info.GetStartedAt(), info.GetTTL())

	info.AddLatency(time.Second)
	info.AddLatency(time.Second + time.Millisecond * 100)

	if len(info.GetLatencies()) != 2 {
		t.Fatal("must be 2 latencies")
	}

	fmt.Println(info.GetAvgLatencies())
	info.ResetLatencies()

	type testCase []struct{
		key string
		value string
	}

	tests := []testCase{
		{
			{key: goweb.HeaderXTransactionId, value: "foo"},
			{goweb.HeaderXTransactionTime, time.Now().Format(time.RFC3339Nano)},
		},
		{
			{key: goweb.HeaderXTransactionId, value: "foo"},
			{goweb.HeaderXTransactionTTL, "4"},
		},
		{
			{goweb.HeaderXTransactionTTL, "4"},
		},
		{
			{key: goweb.HeaderXTransactionId, value: "foo"},
			{goweb.HeaderXTransactionTTL, "dfg"},
		},
		{
			{key: goweb.HeaderXTransactionId, value: "foo"},
			{goweb.HeaderXTransactionTTL, "4"},
			{goweb.HeaderXTransactionTime, "dfg"},
		},

	}
	for i, x := range tests {
		t.Run(fmt.Sprintf("Subtest %v", i), func(t *testing.T) {
			r := &http.Request{Header: make(http.Header)}
			for _, h := range x {
				r.Header.Add(h.key, h.value)
			}
			_, e := goweb.ExtractXTransactionInfo(r)
			if e == nil {
				t.Fatal("must be an error")
			}
		})
	}
}
