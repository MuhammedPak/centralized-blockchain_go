package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Args struct {
	X               string
	Transactionlist []string
}
type New_Block struct {
	newblock string
}

type Connection struct{}

var NodeList []string

func Connections() {//Center nodu 3000 portundan ayağa kaldırır
	cal := new(Connection)
	server := rpc.NewServer()
	server.Register(cal)
	server.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	listener, e := net.Listen("tcp", ":3000")
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
func (t *Connection) Connection_Requests(args *Args, reply *string) error {//Minerdan gelen ağa bağlanma isteğini alır 
	newnode := args.X

	if CheckNodeList(newnode) == true {
		NodeList = append(NodeList, newnode)
		fmt.Println("You can connect to network")
		*reply = "You can connect to network"

	} else {
		fmt.Println("Node address already exists")
		*reply = "You can't connect to network"
	}
	return nil
}
func (t *Connection) Send_All_Block_To_Miner(args *Args, reply *string) error {//Tüm blockları ağa yen katılan miner a gönderir
	allblock := AllBlock()
	data, _ := json.Marshal(allblock)
	*reply = string(data)

	return nil
}

func (t *Connection) Take_Transaction(args *Args, reply *string) error {//Gelen transaction listesini alarak miner a gönderir

	*reply = "Transaction listesi alındı"
	Tx_list := string_to_json(args.X)
	Start_Mining(Tx_list)
	return nil
}

func CheckNodeList(nodename string) bool {//Gelen node ismiyle aynı isimde miner varmı kontrol eder
	for _, item := range NodeList {
		if item == nodename {
			return false
		}
	}
	return true
}

func Start_Mining(Tx_list []string) {//Transaction listesi alındıktan sonra Minerlara gönderilerek proof of work başlatılır
	fmt.Println("Transaction listesi tüm nodlara gönderiliyor ")
	client, err := net.Dial("tcp", "127.0.0.1:3001")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	args := &Args{"Start Mining", Tx_list}
	var reply string
	c := jsonrpc.NewClient(client)
	err = c.Call("Connection.Start_Mining", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	block := string_to_json_block(reply)
	fmt.Println("Yeni Block Bulundu")
	SaveBlock(block, "3000.db")
	List("3000.db")

	Send_New_Block_to_Miner(block)

}

func Send_New_Block_to_Miner(block Block) {//Yeni gelen doğrulanmış bloğu miner a gönderir
	client, err := net.Dial("tcp", "127.0.0.1:3001")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	data, _ := json.Marshal(block)
	args := &Args{string(data), nil}
	var reply string
	c := jsonrpc.NewClient(client)
	err = c.Call("Connection.Add_New_Block", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Println("Yeni Block Miner a gönderildi")
}
func Validate_Miner() {//Ana noda Validate isteği geldikten sonra Minerın blockları kontrol edilmek için minera istek atar
	var control bool
	client, err := net.Dial("tcp", "127.0.0.1:3001")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	args := &Args{"Validation", nil}
	var reply string
	c := jsonrpc.NewClient(client)
	err = c.Call("Connection.Validation", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	Miner_Blocks := string_to_json_Blocklist(reply)
	CenterNode_Blocks := AllBlock()
	for center, miner := range Miner_Blocks {
		if CenterNode_Blocks[center].Blockhash != miner.Blockhash {
			fmt.Println(CenterNode_Blocks[center].Blockhash, "===", miner.Blockhash)
			control = false
			break
		} else {
			control = true
		}
	}
	fmt.Println(len(CenterNode_Blocks), " ", len(Miner_Blocks))
	if len(CenterNode_Blocks) != len(Miner_Blocks) {
		control = false
	}
	if control == true {
		reply = "Miner is valid"
		fmt.Println("Miner Doğrulandı")
	} else {
		reply = "Miner is not valid"
		fmt.Println("Miner Değiştirirlmiş")
		Send_Validate_Blocks_to_Miner()

	}
	// Synchronous call

}
func Send_Validate_Blocks_to_Miner() {//Eğer Minerdaki data değiştirilmiş ise doğru block listesi tekrar minera gönderilir
	client, err := net.Dial("tcp", "127.0.0.1:3001")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	all_block := AllBlock()
	data, _ := json.Marshal(all_block)
	args := &Args{string(data), nil}
	var reply string
	c := jsonrpc.NewClient(client)
	err = c.Call("Connection.Add_Validate_Blocks", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Println("Doğrulanmış blocklar gönderildi")
}
func (t *Connection) Start_Validation(args *Args, reply *string) error {
	fmt.Println("Validation isteği alındı")
	Validate_Miner()
	return nil
}
func string_to_json_block(data string) Block {
	var blocklist Block
	_ = json.Unmarshal([]byte(data), &blocklist)
	return blocklist
}
func string_to_json(data string) []string {
	var blocklist []string
	_ = json.Unmarshal([]byte(data), &blocklist)
	return blocklist
}
func string_to_json_Blocklist(data string) []Block {
	var blocklist []Block
	_ = json.Unmarshal([]byte(data), &blocklist)
	return blocklist
}
