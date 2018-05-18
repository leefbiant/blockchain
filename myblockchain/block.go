package main

import (
	"log"
	"bytes"
	"encoding/json"
	"encoding/binary"
)

func (block *Block)HasHBlock() string {
	var buf bytes.Buffer
	jsonblock, marshalErr := json.Marshal(block)
	if marshalErr != nil {
		log.Fatalf("Could not marshal block: %s", marshalErr.Error())
	}
	hashingErr := binary.Write(&buf, binary.BigEndian, jsonblock) 
	if hashingErr != nil {
		log.Fatalf("Could not hash block: %s", hashingErr.Error())   
	}
	return ComputeHashSha256(buf.Bytes())   
}

