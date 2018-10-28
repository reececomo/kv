package kv

import (
	"bytes"
	"encoding/gob"
)

// StructToBytes : Converts the given structure to a slice of bytes
func StructToBytes(s interface{}) (converted []byte, err error) {
	// create new buffer
	var buff bytes.Buffer
	// create encoder and set to buffer
	encoder := gob.NewEncoder(&buff)
	// encode to the buffer
	if err = encoder.Encode(s); err != nil {
		return
	}
	// return the encoded bytes
	converted = buff.Bytes()
	return
}

// Decoder : returns the decoder for converting bytes (gobs) back into structs
// create a var for decoding the bytes into
// 	var decodeInto exampleOne
// create a new decoder
// 	decoder := dblog.Decoder(rawBytes)
// decode into the desired struct
// 	decoder.Decode(&decodeInto)
// print the desired fields from the struct
// 	fmt.Println(decodeInto.Name)
func Decoder(rawBytes []byte) (decoder *gob.Decoder) {
	reader := bytes.NewReader(rawBytes)
	decoder = gob.NewDecoder(reader)
	return
}
