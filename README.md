# kv

Some basic helpers abstracted from larger library to improve code readability. Currently no tests for this abstraction also, no documentation provided on this abstraction.

Check code comments for usage.

```go
package main

import (
	"fmt"

	"github.com/benluxford/kv"
)

// Logger :
var Logger *kv.Logger

type someStruct struct {
	Name, Description string
	Age               int
}

func main() {
	// create connection to database
	if Logger, err := kv.Connect("tmp"); err == nil {

		// create the data to be persisted to the db
		dataStruct := someStruct{"Ben Luxford", "that random guyee'", 6}
		// convert the structure to bytes slice
		structInBytes, err := kv.StructToBytes(dataStruct)
		if err != nil {
			panic(err)
		}

		// create a record to be saved into the database
		record := kv.Record{}
		// create the key
		record.Key = []byte("example key")
		record.Value = structInBytes

		// set the value in the KV database
		if res := Logger.SetOne(record); res.Error != nil {
			panic(res.Error)
		}

		// create a var for decoding the bytes into after retrieved from the database
		var decodeInto someStruct
		// get the key and value from the database
		if res := Logger.GetOne(record.Key); res.Error != nil {
			panic(res.Error)
		} else {
			// create a new decoder
			decoder := kv.Decoder(res.Record.Value)
			// decode into the desired struct
			if err = decoder.Decode(&decodeInto); err != nil {
				panic(err)
			}
		}

		// disconnect from the database
		if err = Logger.Disconnect(); err != nil {
			panic(err)
		}

		// print the desired fields from the struct
		fmt.Println(decodeInto.Name, decodeInto.Description, decodeInto.Age)

	} else {
		panic(err)
	}
}
```