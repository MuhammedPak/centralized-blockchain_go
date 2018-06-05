package main

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

const targetcount = 3
const targetstring = "0000000000000000000000000000000000000"

type Block struct {
	Datetime      time.Time
	Index         int64
	Previous_hash string
	Nonce         int
	Blockhash     string
	Merkleroot    string
}

var maxNonce = math.MaxInt64

func CreateGenesisBlock(nodeId string) { //Genesisbloğunu oluşturur
	b := Block{}
	b.Datetime = time.Now()
	b.Merkleroot = Hashing("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks")
	b.Index = 1
	b.Previous_hash = ""
	b.Blockhash, b.Nonce, b.Datetime = ProofofWork(b.Merkleroot, b.Previous_hash)
	SaveBlock(b, nodeId)
}

func Find_New_Block(nodeId string, transactionlist []string) Block { //Gelen Transaction listesine göre proof of workü başlatır yeni bloğu bulur

	var transactions []string
	next_block := Block{}
	list := transactionlist //Transaction listesini alıyoruz
	for _, item := range list {
		transactions = append(transactions, Hashing((time.Now().String())+item)) //Transaction listesindeki tüm transactionlerin hashini alıyoruz
	}
	merkletreehash := merkleroot(transactions)                       //Merkle Hashi hesaplamak için transaction listesini merkleroot a gönderiyoruz
	last_block_ındex := Last_Block_Index(nodeId)                     //Son bloğun indexini alıyoruz
	lastblock, _ := GetBlock(strconv.Itoa(last_block_ındex), nodeId) //Son bloğun indexi ile son bloğu getiriyoruz
	previoushash := lastblock.Blockhash
	fmt.Println("Proof of Work başladı")                                      //Son bloğun hashini gelecekte çıkacak olan bloğun previoushashıne ekliyoruz
	next_block_hash, nonce, time := ProofofWork(merkletreehash, previoushash) //proof ow Work ile nonce değeri ve blockhashi hesaplıyouz
	next_block.Blockhash = next_block_hash
	next_block.Datetime = time
	next_block.Nonce = nonce
	next_block.Index = int64(last_block_ındex + 1)
	next_block.Previous_hash = previoushash
	next_block.Merkleroot = merkletreehash
	return next_block //Blocğu database e kaydediyoruz

}
func ProofofWork(merkletree, previoushash string) (string, int, time.Time) { //Merklehash ,previous hashive zamanı alıyoruz nonce değerini buluyoruz
	target := ""
	nonce := 0
	var Powtime time.Time
	for mining == true {

		for nonce < maxNonce {
			Powtime = time.Now()
			target = Hashing(merkletree + previoushash + string(nonce) + Powtime.String())
			if targetstring[0:targetcount] == target[0:targetcount] {
				return target, nonce, Powtime
			} else {
				nonce++
			}
		}

	}
	return target, nonce, Powtime
}

func ValidatePoW(hash, previoushash, merkletree, time string, nonce int) bool {
	return hash == Hashing(merkletree+previoushash+string(nonce)+time)
}
