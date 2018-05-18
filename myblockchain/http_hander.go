package main

import (
	"fmt"
	"log"
	"encoding/json"
	"io"
	"net/http"
)

func NewHandler(blockchain *BlockChain) http.Handler {
	h := handler{blockchain}
	mux := http.NewServeMux()
	fmt.Println("NewHandler :", h.blockchain.chain)
	mux.HandleFunc("/mine", BuildResponse(h.Mine)) // 生成区块, 将交易添加到区块上 (挖矿)
	mux.HandleFunc("/chain", BuildResponse(h.BlockChain)) // 获取链
	mux.HandleFunc("/txion", BuildResponse(h.NewTransaction)) // 添加交易信息

	mux.HandleFunc("/nodes/register", BuildResponse(h.RegisterNode)) // 注册全节点
	mux.HandleFunc("/nodes/getchainnode", BuildResponse(h.GetChainNode)) // 查询区块链节点
	mux.HandleFunc("/nodes/resolve", BuildResponse(h.ResolveConflicts)) // 同步区块链 
	return mux
}

func BuildResponse(h func(io.Writer, *http.Request) response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := h(w, r)
		msg := resp.value
		if resp.err != nil {
			msg = resp.err.Error()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.statusCode)
		fmt.Printf("resp msg:%v\n", msg)
		if err := json.NewEncoder(w).Encode(msg); err != nil {
			log.Printf("could not encode response to output: %v", err)
		}
	}
}


func (h *handler) Mine(w io.Writer, r *http.Request) response {
	if r.Method != http.MethodPost {
		return response {
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowd", r.Method),
		}
	}
	block := h.blockchain.CreateBlock()
	resp := map[string]interface{}{"message": "New Block Forged", "block": block}
	return response{resp, http.StatusOK, nil}
}

func (h *handler) NewTransaction(w io.Writer, r *http.Request) response {
	if r.Method != http.MethodPost {
		return response {
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowd", r.Method),
		}
	}
	var tx Transaction
	dec := json.NewDecoder(r.Body);
	dec.DisallowUnknownFields()
	err := dec.Decode(&tx)
	if err != nil {
		log.Printf("there was an error when trying to add a transaction %v\n", err)
		return response {
			nil,
			http.StatusInternalServerError,
			fmt.Errorf("fail to add transaction to the blockchain"),
		}
	}
	h.blockchain.AddrTransaction(&tx)
	resp := map[string]interface{}{"message": "Transaction submission successful"}
  return response{resp, http.StatusOK, nil}

}

func (h *handler) BlockChain(w io.Writer, r *http.Request) response {
	if r.Method != http.MethodGet {
		return response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowd", r.Method),
		}
	}
	log.Println("BlockChain requested")
	chain := h.blockchain.GetChain()
	resp := map[string]interface{}{"chain": chain, "length": len(h.blockchain.chain)}
	log.Println("BlockChain :", resp)
	return response{resp, http.StatusOK, nil}
}

func (h *handler) RegisterNode(w io.Writer, r *http.Request) response {
	if r.Method != http.MethodPost {
		return response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowd", r.Method),
		}
	}
	log.Println("BlockChain requested")
	var body map[string][]string
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("there was an error when trying to RegisterNode%v\n", err)
		return response {
			nil,
			http.StatusInternalServerError,
			fmt.Errorf("fail to add transaction to the blockchain"),
		}
	}
	for _, node := range body["nodes"] {
		h.blockchain.RegisterNode(node)
	}

	var keys []string
	for key, _ := range h.blockchain.nodes {
		keys = append(keys, key)
	}
	resp := map[string]interface{}{
		"message": "New nodes have been added",
		"nodes":   keys,
	}
	status := http.StatusCreated
	return response{resp, status, err}
}

func (h *handler) GetChainNode(w io.Writer, r *http.Request) response {
	if r.Method != http.MethodPost {
		return response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowd", r.Method),
		}
	}
	var keys []string
	for key, _ := range h.blockchain.nodes {
		keys = append(keys, key)
	}
	resp := map[string]interface{}{
		"message": "chain nodes",
		"nodes":   keys,
	}
	status := http.StatusCreated
	return response{resp, status, nil}
}

func (h *handler) ResolveConflicts(w io.Writer, r *http.Request) response {
	if r.Method != http.MethodPost {
		return response{
			nil,
			http.StatusMethodNotAllowed,
			fmt.Errorf("method %s not allowd", r.Method),
		}
	}
	msg := "Our chain is authoritative"
	if h.blockchain.ResolveConflicts() {
		msg = "Our chain was replaced" 
	}
	resp := map[string]interface{}{"message": msg, "chain": h.blockchain.chain} 
	status := http.StatusCreated
	return response{resp, status, nil}
}
