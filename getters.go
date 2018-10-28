package kv

import (
	"fmt"
	"regexp"

	"github.com/dgraph-io/badger"
)

// GetOne : the Query response from the given key string. for something more
// dynamic, try GetPrefix("prefix string") or the slower GetMatches("RegExPattern"),
// warning GetMatches will loop over entire Key store
func (l *Logger) GetOne(key []byte) (response Response) {
	// create the empty response
	response = Response{}
	// query and set the error if one is returned
	response.Error = l.KV.View(func(txn *badger.Txn) (err error) {
		// get the item from the KV
		if item, err := txn.Get(key); err != nil {
			err = fmt.Errorf("unable to get the key %s from the KV Database %v", key, err)
		} else {
			// set the key that was queried
			response.Record.Key = key
			// copy the value outside of the scope
			response.Record.Value, err = item.ValueCopy(nil)
		}
		return
	})
	return
}

// GetMatches : Gets all the values where the keys match the given regex pattern
func (l *Logger) GetMatches(regExPattern string) (responseCollection ResponseCollection, err error) {
	if matchingKeys, err := l.GetMatchingKeys(regExPattern); err == nil {
		responseCollection = l.GetMany(matchingKeys)
	}
	return
}

// GetMany : Returns all of the results for the given keyset
func (l *Logger) GetMany(keys [][]byte) (responseCollection ResponseCollection) {
	// create the var to hold all; the responses
	responseCollection = ResponseCollection{}
	responseCollection.Error = l.KV.View(func(txn *badger.Txn) (err error) {
		for _, key := range keys {
			if item, err := txn.Get(key); err != nil {
				break
			} else {
				response := Response{}
				// set the response key
				response.Record.Key = key
				// copy the response value, if error break loop and return
				if response.Record.Value, err = item.ValueCopy(nil); err != nil {
					break
				}
				// append response to responses
				responseCollection.Responses = append(responseCollection.Responses, response)
			}
		}
		return
	})
	return
}

// GetMatchingKeys : Returns all the keys that match the given regex pattern
func (l *Logger) GetMatchingKeys(regExPattern string) (keys [][]byte, err error) {
	// create the new slice of byte slices
	keys = [][]byte{}
	// generate the pattern
	if pattern, err := regexp.Compile(regExPattern); err == nil {
		// open the view
		err = l.KV.View(func(txn *badger.Txn) (err error) {
			// set to not prefetch the record values
			opts := badger.DefaultIteratorOptions
			opts.PrefetchValues = false
			// create a new iterator
			it := txn.NewIterator(opts)
			// close the iterator at the end
			defer it.Close()
			// for each entry
			for it.Rewind(); it.Valid(); it.Next() {
				// get the item
				item := it.Item()
				// get the key
				key := item.Key()
				// test against the pattern
				if pattern.Match(key) {
					// record the key
					keys = append(keys, key)
				}
			}
			return
		})
	}
	// if there wer no results that matched, return the error that none were found
	if len(keys) < 1 && err == nil {
		err = fmt.Errorf("could not match any keys with the pattern %s", regExPattern)
	}
	return
}

// GetPrefix : Returns all values where the key is prefixed with (prefix string)
func (l *Logger) GetPrefix(prefix string) (responses ResponseCollection) {
	// create the collection of responses
	responses = ResponseCollection{}
	// convert the prefix to byte slice
	bytePrefix := []byte(prefix)
	// set the error if one exists
	responses.Error = l.KV.View(func(txn *badger.Txn) (err error) {
		// create the new iterator
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		// close iterator at end
		defer it.Close()
		// loop through the values matching the prefix
		for it.Seek(bytePrefix); it.ValidForPrefix(bytePrefix); it.Next() {
			// set the item
			item := it.Item()
			// create a new response object
			response := Response{}
			// save the values in the response
			response.Record.Key = item.Key()
			// copy the value and if cant return err
			if response.Record.Value, err = item.ValueCopy(nil); err != nil {
				return
			}
			// append the response to the slice of Responses
			responses.Responses = append(responses.Responses, response)
		}
		return
	})
	return
}
