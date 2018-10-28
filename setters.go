package kv

import (
	"github.com/dgraph-io/badger"
)

// SetOne : Saves an individual record to the KV store, it is suggested that if you are to record
// high frequency fo records that you use SetMany(map[key][]byte(value)), collect the values then
// persist to the Key, Value store
func (l *Logger) SetOne(record Record) (response Response) {
	// Create the empty response
	response = Response{}
	// create the update connection, set the error if failure
	response.Error = l.KV.Update(func(txn *badger.Txn) (err error) {
		// set the individual key, value pair
		err = txn.Set(record.Key, record.Value)
		return
	})
	return
}

// SetMany : Takes a map of string, byte slice [key]value. arranges collection for save and persists
// all elements to the Key, Value store
func (l *Logger) SetMany(records []Record) (response Response) {
	// create the empty response
	response = Response{}
	// create a fresh transaction
	txn := l.KV.NewTransaction(true)
	// for all of the values in the collection
	for _, record := range records {
		// if there are too many in the queue
		if err := txn.Set(record.Key, record.Value); err == badger.ErrTxnTooBig {
			// commit the current queue
			if response.Error = txn.Commit(); response.Error != nil {
				return
			}
			// start a fresh queue
			txn = l.KV.NewTransaction(true)
			// place the failed item in fresh queue
			if response.Error = txn.Set(record.Key, record.Value); response.Error != nil {
				return
			}
			// continue to add items to the fresh queue
		}
	}
	// commit the leftover amount of transactions
	response.Error = txn.Commit()
	return
}
