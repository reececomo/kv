package kv

// KV Connection : manages the connection to the key value database or its instance
// IMPORTANT: There can only be one connection to a KV store at a time, if the connection is
// required my multiple packages, consider fragmentation of databases or setting the instance as global.

import (
	"runtime"

	"github.com/dgraph-io/badger"
)

// Connect : the initial connection to the Key Value store database,
// returns the Logger struct with the connection and helper functions.
// Remember to logger.Disconnect() when all logging is complete.
func Connect(databaseFolderLocation string) (logger *Logger, err error) {
	// Create the new struct
	logger = new(Logger)
	// Set the badger options
	options := badger.DefaultOptions
	options.Dir = databaseFolderLocation
	options.ValueDir = databaseFolderLocation
	// Set truncate to true - windows only
	if runtime.GOOS == "windows" {
		options.Truncate = true
	}
	logger.KV, err = badger.Open(options)
	return
}

// Disconnect : removes the connection to the Key Value store
// this should be done when the application no longer needs the instance of
// the logging service or never if your application logs for life
func (l *Logger) Disconnect() error { return l.KV.Close() }
