package main

import (
	"fmt"
	"flag"
	"github.com/syndtr/goleveldb/leveldb"
)

var dbpath string

// go run get_ldb_datas.go -dbpath  /mydbpath
func init() {
	flag.StringVar(&dbpath,"dbpath", "", "Path to LevelDB")
}


func readAll(db *leveldb.DB) {

	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		kbyes := iter.Key()

		key   := string(kbyes)
		fmt.Println("key:",key)
		fmt.Println("keyBytes:",kbyes)
		fmt.Println("value:", string(iter.Value()))
		fmt.Println("valueBytes:", iter.Value())
		fmt.Println("===========================================")
	}
	iter.Release()
}

func main() {
	flag.Parse()
	if  dbpath  == "" {
		fmt.Printf("ERROR: dbpath could be empty\n")
		return
	}

	db, err := leveldb.OpenFile(dbpath, nil)
	if err != nil {
		fmt.Printf("ERROR: Cannot open LevelDB from [%s], with error=[%v]\n", dbpath, err);
	}
	defer db.Close()

	readAll(db)

}
