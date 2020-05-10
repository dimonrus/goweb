package goweb

import (
	"github.com/dimonrus/porterr"
	"net/http"
	"strconv"
	"time"
)

const (
	HeaderXTransactionId   = "x-transaction-id"
	HeaderXTransactionTTL  = "x-transaction-ttl"
	HeaderXTransactionTime = "x-transaction-time"
)

// Transaction alias for type string
type XTransactionId string

// Transaction info
type XTransactionInfo struct {
	// Id of transaction
	id XTransactionId
	// When transaction must be done
	ttl int
	// When transaction was started
	startedAt time.Time
	// latencies list per heartbeat
	latencies []time.Duration
}

// Return X Transaction Id
func (x *XTransactionInfo) GetId() XTransactionId {
	return x.id
}

// Return X Transaction TTL
func (x *XTransactionInfo) GetTTL() int {
	return x.ttl
}

// Return X Transaction Started at
func (x *XTransactionInfo) GetStartedAt() time.Time {
	return x.startedAt
}

// Return X Transaction Latencies
func (x *XTransactionInfo) GetLatencies() []time.Duration {
	return x.latencies
}

// Calculate avg for latencies
func (x *XTransactionInfo) GetAvgLatencies() time.Duration {
	var avg time.Duration
	for _, v := range x.latencies {
		avg += v
	}
	return avg / time.Duration(len(x.latencies))
}

// Add X Transaction Latency
func (x *XTransactionInfo) AddLatency(l time.Duration) *XTransactionInfo {
	x.latencies = append(x.latencies, l)
	return x
}

// Reset X Transaction Latencies
func (x *XTransactionInfo) ResetLatencies() *XTransactionInfo {
	x.latencies = make([]time.Duration, 0)
	return x
}

// New XTransactionInfo
func NewXTransactionInfo(id XTransactionId, ttl int, startedAt time.Time) *XTransactionInfo {
	return &XTransactionInfo{
		id:        id,
		ttl:       ttl,
		startedAt: startedAt,
	}
}

// Get form context x transaction information
func ExtractXTransactionInfo(r *http.Request) (*XTransactionInfo, porterr.IError) {
	transactionId := r.Header.Get(HeaderXTransactionId)
	if transactionId == "" {
		return nil, porterr.New(porterr.PortErrorTransaction, HeaderXTransactionId+" header is required")
	}
	ttl := r.Header.Get(HeaderXTransactionTTL)
	if ttl == "" {
		return nil, porterr.New(porterr.PortErrorTransaction, HeaderXTransactionTTL+" header is required")
	}
	ttlValue, err := strconv.Atoi(ttl)
	if err != nil {
		return nil, porterr.New(porterr.PortErrorTransaction, HeaderXTransactionTTL+" value must be an int")
	}
	transactionTime := r.Header.Get(HeaderXTransactionTime)
	if transactionTime == "" {
		return nil, porterr.New(porterr.PortErrorTransaction, HeaderXTransactionTime+" header is required")
	}
	transactionTimeUnixNano, err := time.Parse(time.RFC3339Nano, transactionTime)
	if err != nil {
		return nil, porterr.New(porterr.PortErrorTransaction, HeaderXTransactionTime+" value must be in RFC3339Nano format")
	}
	return NewXTransactionInfo(XTransactionId(transactionId), ttlValue, transactionTimeUnixNano), nil
}
