package main

import (
	"encoding/json"
	"log"
	"net"
	"net/rpc/jsonrpc"
)

type Args struct {
	X string
}

func main() {
	client, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	data, _ := json.Marshal(Transactionlist())
	args := &Args{string(data)}
	var reply string
	c := jsonrpc.NewClient(client)
	err = c.Call("Connection.Take_Transaction", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
}

func Transactionlist() []string {
	transactions := []string{}
	tx1 := "Send 1 BTC From Muhammed To Mehmetcan"
	tx2 := "Send 2 BTC From Muhammed To Mehmetcan"
	tx3 := "Send 3 BTC From Muhammed To Mehmetcan"
	tx4 := "Send 4 BTC From Muhammed To Mehmetcan"
	tx5 := "Send 5 BTC From Muhammed To Mehmetcan"
	transactions = append(transactions, tx1)
	transactions = append(transactions, tx2)
	transactions = append(transactions, tx3)
	transactions = append(transactions, tx4)
	transactions = append(transactions, tx5)
	return transactions
}
