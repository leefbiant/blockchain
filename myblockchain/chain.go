package main

import (
	"time"
	"log"
	"net/url"
	// "github.com/golang/glog"
)

func (bc *BlockChain)ProofOfWork(last_proof int64) int64 {
	var proof int64 = 0
	for !ValidProof(last_proof, proof) {
		proof++
	}
	return proof
}

func (bc *BlockChain)AddrTransaction(trans *Transaction) string {
	bc.transactions = append(bc.transactions, *trans)
	log.Printf("New transaction From:%v To:%v Amount:%v", trans.From, trans.To, trans.Amount)
	return string("Transaction submission successful")
}

func CreateGennsisChain() *BlockChain {
	blockchain := &BlockChain {
		chain : make([]Block, 0),
		transactions : make([]Transaction, 0),
		nodes : make(map[string]bool),
	}
	block := Block{
		Index: 0,
		Timestamp: time.Now().UnixNano(),
		Proof: 1,
		Data: make([]Transaction, 0),
		Previous_hash: "0",
	}
	block.Hash = block.HasHBlock()
	blockchain.chain = append(blockchain.chain, block)
	return blockchain
}

func (bc *BlockChain)CreateBlock() Block {
	last_block := bc.chain[len(bc.chain)-1]
	// 添加交易数据, 生成区块
	new_block := Block {
		Index : last_block.Index + 1,
		Timestamp : time.Now().UnixNano(),
		Proof : bc.ProofOfWork(last_block.Proof),
		Data : bc.transactions,
		Previous_hash : last_block.Hash,
	}
	// 清空交易数据
	bc.transactions = nil
	new_block.Hash = new_block.HasHBlock()
	bc.chain = append(bc.chain, new_block)
	return new_block
}

func (bc *BlockChain)GetChain() []Block {
	return bc.chain;
}


func (bc *BlockChain) RegisterNode(address string) bool {
	u, err := url.Parse(address)
	if err != nil {
		log.Println("url.Parse failed:", err)
		return false
	}
	bc.nodes[u.Host] = true
	log.Println("add node:", u.Host)
	return true
}


func (bc *BlockChain) ResolveConflicts() bool {
	var nodes []string
	for key, _ := range bc.nodes {
		nodes = append(nodes, key)
	}
	log.Println("other nodes:", nodes)

	max_chain_length := len(bc.chain)
	new_chain := make([]Block, 0)
	for _, node_addr := range nodes {
		blockchaininfo, err := GetChain(node_addr)
		if err != nil {
			log.Printf("GetChain failed from:%v err:%v", node_addr, err)
			continue
		}
		log.Printf("start to check chain:%v", blockchaininfo.Chain)
		if (blockchaininfo.Length > max_chain_length) && !ValidChain(&blockchaininfo.Chain) {
			max_chain_length = blockchaininfo.Length
			new_chain = blockchaininfo.Chain
		}
	}
	if len(new_chain) > 0 {
		bc.chain = new_chain
		return true
	}
	return false
}
