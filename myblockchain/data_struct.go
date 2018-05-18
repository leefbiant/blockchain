package main

// for chain
type Transaction struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Amount    int32  `json:"amount"`
}

type Block struct {
	Index int64 `json:"index"`
	Timestamp int64 `json:"timestamp"`
	Data []Transaction  `json:"data"`
	Proof  int64    `proof:"proof"`
	Previous_hash string  `json:"previous_hash"`
	Hash string `json:"hash"`
}

type BlockChain struct {
	chain []Block
	transactions []Transaction
	nodes map[string]bool
}

// for http
type handler struct {
	blockchain *BlockChain
}

type response struct {
	value      interface{}
	statusCode int
	err        error
}

type blockchainInfo struct {
	Length int     `json:"length"`
	Chain  []Block `json:"chain"`
}
