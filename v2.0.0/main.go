package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dgraph-io/badger"
)

func main() {
	path := flag.String("d", "../db", "path containing the badger data")
	flag.Parse()

	bdb, err := badger.Open(badger.DefaultOptions(*path).WithLogger(nil))
	if err != nil {
		log.Fatal(err)
	}
	db := database{bdb}
	defer db.close()

	s := bufio.NewScanner(os.Stdin)
	for fmt.Printf("> "); s.Scan(); fmt.Printf("> ") {
		line := strings.TrimSpace(s.Text())
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}

		var err error
		switch first := fields[0]; first {
		case "bye", "exit", "quit":
			fmt.Println("good bye")
			return
		case "set":
			err = db.set(fields)
		case "get":
			err = db.get(fields)
		case "delete":
			err = db.delete(fields)
		case "ls":
			err = db.list(fields)
		default:
			err = fmt.Errorf("unknown command %q", first)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	if err := s.Err(); err != nil {
		log.Fatalf("reading from stdin: %v", err)
	}
}

type database struct{ db *badger.DB }

func (db database) close() {
	if err := db.db.Close(); err != nil {
		log.Fatal(err)
	}
}

func (db database) get(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("get expects one key as parameter: 'get key'")
	}
	key := args[1]

	var val string
	err := db.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(b []byte) error { val = string(b); return nil })
	})
	if err != nil {
		return err
	}
	fmt.Println(val)
	return nil
}

func (db database) set(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("set expects one key and one value as parameters: 'set key value'")
	}
	key := args[1]
	val := args[2]

	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(val))
	})
}

func (db database) delete(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("delete expects one key as parameter: 'delete key'")
	}
	key := args[1]

	return db.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (db database) list(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("ls doesn't expect any parameters")
	}
	return db.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.IteratorOptions{})
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			fmt.Printf("%s\n", it.Item().Key())
		}
		return nil
	})
}

