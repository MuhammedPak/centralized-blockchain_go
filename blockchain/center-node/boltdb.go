package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/boltdb/bolt"
)

var db *bolt.DB
var open bool

func Open(dbname string) error { //Datadb ye bağlanmak için kullanılır
	var err error
	_, filename, _, _ := runtime.Caller(0)          // get full path of this file
	dbfile := path.Join(path.Dir(filename), dbname) //Data db dosyasına bağlanır
	config := &bolt.Options{Timeout: 1 * time.Second}
	db, err = bolt.Open(dbfile, 0600, config)
	if err != nil {
		log.Fatal(err)
	}
	open = true
	return nil
}

func Close() { //Data db ile bağlantıyı sonlandırmak için
	open = false
	db.Close()
}

func (p *Block) encode() ([]byte, error) { //Nesne bytearray tipine dönüştürülür
	enc, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}
	return enc, nil
}
func decode(data []byte) (*Block, error) { //Byte array tipindeki veri nesneye dönüştürülür
	var b *Block
	err := json.Unmarshal(data, &b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func List(nodeId string) { //Listedeki tüm nesneleri dönderir
	Open(nodeId)
	defer Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(nodeId)).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("Index=%s, value=%s\n", k, v)
		}
		return nil
	})
}

func GetBlock(id, dbname string) (*Block, error) { //İD sone göre nesne dönderir
	Open(dbname)
	defer Close()
	if !open {
		return nil, fmt.Errorf("db must be opened before saving!")
	}
	var p *Block
	err := db.View(func(tx *bolt.Tx) error {
		var err error
		b := tx.Bucket([]byte(dbname))
		k := []byte(id)
		p, err = decode(b.Get(k))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Could not get Block ID %s", id)
		return nil, err
	}
	return p, nil
}

func SaveBlock(b Block, nodeId string) error { //Blockları kaydeder
	Open(nodeId)
	defer Close()
	err := db.Update(func(tx *bolt.Tx) error {
		bucketname, err := tx.CreateBucketIfNotExists([]byte(nodeId)) //Eklenecek bucket adı
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		enc, err := b.encode() //Veri byte dönüştürülür
		if err != nil {
			return fmt.Errorf("could not encode Block %s: %s", b.Index, err)
		}
		encIndex, _ := json.Marshal(b.Index)
		err = bucketname.Put([]byte(string(encIndex)), enc) //Byte tipindeki veri ve idsi bucket a eklenir
		return err
	})
	return err
}
func SaveBlockchain(nodeId string) error { //Nodun blockları tutacağı veritabanını oluşturur
	Open(nodeId)
	defer Close()
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(nodeId)) //Eklenecek bucket adı
		if err != nil {
			fmt.Println("%s zaten mevcut", nodeId)
			return nil
		} else {
			fmt.Println("%v oluşturuldu", nodeId)
		}
		return nil
	})
	return err

}

func Last_Block_Index(nodeId string) int { //Son bloğun indexini dönderir
	Open(nodeId)
	defer Close()
	var last_block_ındex int
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte(nodeId)).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("Index=%s, value=%s\n", k, v)
			_ = json.Unmarshal(k, &last_block_ındex)
		}

		return nil
	})

	return last_block_ındex
}

func dbExists(dbFile string) bool { //Gelen verinin olup olmadığını kontrol eder
	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return false
	} //Thanks Jeiwan
	return true
}

func AllBlock() []Block {//Tüm bockları dönderir
	blocks := []Block{}
	block := Block{}
	Open("3000.db")
	defer Close()
	db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket([]byte("3000.db")).Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			json.Unmarshal(v, &block)
			blocks = append(blocks, block)
		}
		return nil
	})
	return blocks
}
