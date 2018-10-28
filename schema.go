package kv

import (
	"github.com/dgraph-io/badger"
)

// Logger : the global connection and logger application
type Logger struct {
	KV *badger.DB
}

// Record : The basic Key, Value pair
type Record struct {
	Key   []byte
	Value []byte
}

// Response : The reply to any singular query
type Response struct {
	Record Record
	Error  error
}

// ResponseCollection : The collection of responses to a group query
type ResponseCollection struct {
	Responses []Response
	Error     error
}
