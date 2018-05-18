package main

import (
	"fmt"
	"net/http"
	"crypto/sha256"
	"encoding/json"
)

func ComputeHashSha256(bytes []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(bytes))
}

func ValidProof(lastProof, proof int64) bool {
	guess := fmt.Sprintf("%d%d", lastProof, proof)
	sHash := ComputeHashSha256([]byte(guess))
	return sHash[:4] == "0000"
}

// 校验区块链
func ValidChain(chain *[]Block) bool {
	last_block := (*chain)[0]
	curr_index := 1
	for curr_index < len(*chain) {
		curr_block := (*chain)[curr_index]
		if curr_block.Previous_hash != last_block.HasHBlock() {
			return false
		}
		if !ValidProof(last_block.Proof, curr_block.Proof) {
			return false
		}
	}
	return true
}

func GetChain(address string) (blockchainInfo, error) {
	response, err := http.Get(fmt.Sprintf("http://%s/chain", address))
	if err == nil && response.StatusCode == http.StatusOK {
		var bi blockchainInfo
		if err := json.NewDecoder(response.Body).Decode(&bi); err != nil {
			return blockchainInfo{}, err
		}
		return bi, nil
	}
	return blockchainInfo{}, err
}
