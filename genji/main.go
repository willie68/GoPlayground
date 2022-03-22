package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/genjidb/genji"
	"github.com/genjidb/genji/document"
	"github.com/genjidb/genji/types"
)

type message struct {
	Id      int    `genji:"id"`
	Tenant  string `genji:"tenant"`
	Origin  string `genji:"origin"`
	Body    string `genji:"body"`
	Message string `genji:"message"`
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

const DBFile = "r:\\genji.db"
const doImport = false

func main() {

	if doImport {
		os.Remove(DBFile) // I delete the file to avoid duplicated records.
	}
	// SQLite is a file based database.

	log.Println("Creating genji.db...")
	log.Println("search")
	start := time.Now()
	db, err := genji.Open(DBFile)
	duration := time.Since(start)
	log.Printf("load duration: %v", duration)
	if err != nil {
		log.Fatal(err)
	}
	// Don't forget to close the database when you're done
	defer db.Close()

	// Attach context, e.g. (*http.Request).Context().
	db = db.WithContext(context.Background())

	// Create a table. Schemas are optional, you don't need to specify one if not needed
	// or you can create a table with constraints on certain fields

	if doImport {
		createTable(db)
	}
	// Create an index
	//    err = db.Exec("CREATE INDEX idx_user_city_zip ON user (address.city, address.zipcode)")

	if doImport {
		importMany(db)
	}

	// Query some documents
	log.Println("search")
	start = time.Now()
	searchMessage(db, "origin LIKE 'w.kla%' AND id > 10 AND id < 12 AND tenant='tnt000000'")
	duration = time.Since(start)
	log.Printf("search duration: %v", duration)

	//	start = time.Now()
	//	searchMessage(db, "message LIKE '%abcx%'")
	//	duration = time.Since(start)
	//	log.Printf("duration: %v", duration)
}

func createTable(db *genji.DB) {
	err := db.Exec(`CREATE TABLE messages (id INT PRIMARY KEY)`)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func importMany(db *genji.DB) {
	st, err := db.Prepare(`INSERT INTO messages VALUES ?`)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	count := 0

	// index some data
	start := time.Now()
	for x := 0; x < 100; x++ {
		fmt.Printf(".")
		tnt := fmt.Sprintf("tnt%.6d", x)
		for i := 0; i < 10000; i++ {
			mymessage := message{
				Id:      count,
				Tenant:  tnt,
				Origin:  "w.klaas@gmx.de",
				Body:    fmt.Sprintf("indexing is easy %d", count),
				Message: randSeq(128),
			}
			err := insertMessage(st, mymessage)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			count++
		}
	}
	duration := time.Since(start)
	log.Printf("duration: %v", duration)
	fmt.Println()
}

// We are passing db reference connection from main to our method with other parameters
func insertMessage(st *genji.Statement, msg message) error {
	return st.Exec(&msg)
}

func searchMessage(db *genji.DB, query string) {
	queryStr := fmt.Sprintf("SELECT * FROM messages WHERE %s;", query)
	log.Println("Query:", queryStr)
	res, err := db.Query(queryStr)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Close()
	count := 0
	// Iterate over the results
	err = res.Iterate(func(d types.Document) error {
		// It is also possible to scan the results into a structure
		var msg message
		err = document.StructScan(d, &msg)
		if err != nil {
			return err
		}

		jsonStr, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		fmt.Println(string(jsonStr))
		count++
		return nil
	})
	log.Println("found ", count, " hits")
}
