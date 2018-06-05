package main

import (
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

	args := &Args{""}
	var reply string
	c := jsonrpc.NewClient(client)
	err = c.Call("Connection.Start_Validation", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)

	}

}
