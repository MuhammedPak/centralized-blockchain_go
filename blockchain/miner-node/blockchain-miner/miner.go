package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"
)

var mining = true

type Args struct {
	X               string
	Transactionlist []string
}
type Connection struct{}
type New_Block struct {
	newblock string
}

func Run_Miner_Node(newNode string) { //Yeni node oluştutut center node a istek atar

	Handshake(newNode)  //Center noda ağa bağlanma isteği yollar center node kabul ederse tüm blockları center node tan alır
	Miner_Connections() //Miner nodu 3001 portundan ayağa kaldırır

}

func Miner_Connections() { //Miner nodu 3001 portundan ayağa kaldırır
	cal := new(Connection)
	server := rpc.NewServer()
	server.Register(cal)
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listener, e := net.Listen("tcp", ":3001")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	for {
		if conn, err := listener.Accept(); err != nil {
			log.Fatal("accept error: " + err.Error())
		} else {
			log.Printf("new connection established\n")
			go server.ServeCodec(jsonrpc.NewServerCodec(conn))
		}
	}
}

func (t *Connection) Start_Mining(args *Args, reply *string) error {
	fmt.Println("Transaction listesi alındı")

	new_block := Find_New_Block("3001.db", args.Transactionlist)
	fmt.Println("Yeni Block Bulundu")
	data, _ := json.Marshal(new_block)
	*reply = string(data)
	fmt.Println("Yeni Block CenterNode a gönderildi")
	fmt.Println(new_block)
	return nil
}
func (t *Connection) Add_New_Block(args *Args, reply *string) error {
	fmt.Println("Yeni Block alındı")
	new_block := string_to_json_block(args.X)
	SaveBlock(new_block, "3001.db")
	*reply = "thanx"
	List("3001.db")
	return nil
}
func (t *Connection) Add_Validate_Blocks(args *Args, reply *string) error {
	DeleteBucket()
	fmt.Println("Doğrulanmış blocklar alındı")
	new_block := string_to_json_blocklist(args.X)
	for _, item := range new_block {
		SaveBlock(item, "3001.db")
	}
	*reply = "thanx"
	List("3001.db")
	return nil
}

func Handshake(newnode string) { //Center noda ağa bağlanma isteği yollar center node kabul ederse tüm blockları center node tan alır
	client, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	args := &Args{newnode, nil}
	var reply string
	c := jsonrpc.NewClient(client)
	err = c.Call("Connection.Connection_Requests", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Println("Node address already exists")
	if reply == "You can connect to network" {
		Take_All_Blocks(newnode)
	}

}
func Take_All_Blocks(newnode string) { //Tüm blockları center node tan almak için center noda istek atar
	client, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	args := &Args{newnode, nil}
	var reply string
	c := jsonrpc.NewClient(client)
	err = c.Call("Connection.Send_All_Block_To_Miner", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	if dbExists(newnode) == false {

		SaveBlockchain(newnode)
		blocks := string_to_json_blocklist(reply)
		for _, item := range blocks {
			SaveBlock(item, newnode)
		}
		List("3001.db")
	}

}
func (t *Connection) Validation(args *Args, reply *string) error { //Center node tan gelen Validate mesajını alır ve blocklarını ana noda gönderir
	allblock := AllBlock()
	data, _ := json.Marshal(allblock)
	*reply = string(data)

	return nil
}
func Add_Changed_Data() { //Değiştirilmiş Block eklemek için kullanıyoruz
	var blocks Block
	blocks.Blockhash = "1"
	blocks.Datetime = time.Now()
	blocks.Index = 23
	blocks.Merkleroot = "1"
	blocks.Nonce = 1
	blocks.Previous_hash = "1a"
	SaveBlock(blocks, "3001.db")
}
func string_to_json_block(data string) Block {
	var blocklist Block
	_ = json.Unmarshal([]byte(data), &blocklist)
	return blocklist
}
func string_to_json_blocklist(data string) []Block {
	var blocklist []Block
	_ = json.Unmarshal([]byte(data), &blocklist)
	return blocklist
}
