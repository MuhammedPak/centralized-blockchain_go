package main

import (
	"fmt"
	"math"
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
var mining = true

func CreateGenesisBlock(nodeId string) { //Genesisbloğunu oluşturur
	b := Block{}
	b.Datetime = time.Now()
	b.Merkleroot = Hashing("The Times 03/Jan/2009 Chancellor on brink of second bailout for banks")
	b.Index = 1
	b.Previous_hash = "0"
	b.Blockhash, b.Nonce, b.Datetime = ProofofWork(b.Merkleroot, b.Previous_hash)
	SaveBlock(b, nodeId)
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

func ValidatePoW(hash, previoushash, merkletree, time string, nonce int) bool { //Center noda gelen bloğun doğru olup olmadığını kontrol eder
	return hash == Hashing(merkletree+previoushash+string(nonce)+time)
}
func Run_Center_Node(newNode string) { //Center nodu ayağa kaldırır ve center node 3000 portundan dinlemeye başlar

	centernode := fmt.Sprintf("%s.db", newNode)
	CreateGenesisBlock(centernode)
	List("3000.db")
	Connections()

}
