package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"

	_ "modernc.org/sqlite" // Import go-sqlite3 library
)

type message struct {
	Id      int
	From    string
	Body    string
	Message string
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
	// SQLite is a file based database.

	log.Println("Creating sqlite-database.db...")
	file, err := os.Create("sqlite-database.db") // Create SQLite file
	if err != nil {
		log.Fatal(err.Error())
	}
	file.Close()
	log.Println("sqlite-database.db created")

	sqliteDatabase, err := sql.Open("sqlite", "./sqlite-database.db") // Open the created SQLite File
	if err != nil {
		log.Fatal(err.Error())
	}
	defer sqliteDatabase.Close() // Defer Closing the database
	createTable(sqliteDatabase)  // Create Database Tables

	count := 0
	// index some data
	for x := 0; x < 1; x++ {
		fmt.Printf(".")
		for i := 0; i < 1000; i++ {
			mymessage := message{
				Id:      count,
				From:    "marty.schoch@gmail.com",
				Body:    fmt.Sprintf("bleve indexing is easy %d", count),
				Message: randSeq(128),
			}
			insertMessage(sqliteDatabase, mymessage)
			count++
			if err != nil {
				log.Fatalf("error: %v", err)
			}
		}
	}
	fmt.Println()

	/*
		// INSERT RECORDS
		insertMessage(sqliteDatabase, "0001", "Liana Kim", "Bachelor")
		insertMessage(sqliteDatabase, "0002", "Glen Rangel", "Bachelor")
		insertMessage(sqliteDatabase, "0003", "Martin Martins", "Master")
		insertMessage(sqliteDatabase, "0004", "Alayna Armitage", "PHD")
		insertMessage(sqliteDatabase, "0005", "Marni Benson", "Bachelor")
		insertMessage(sqliteDatabase, "0006", "Derrick Griffiths", "Master")
		insertMessage(sqliteDatabase, "0007", "Leigh Daly", "Bachelor")
		insertMessage(sqliteDatabase, "0008", "Marni Benson", "PHD")
		insertMessage(sqliteDatabase, "0009", "Klay Correa", "Bachelor")
	*/
	log.Println("all messages")
	// DISPLAY INSERTED RECORDS
	displayMessages(sqliteDatabase)

	log.Println("search students")
	searchMessage(sqliteDatabase, "id=10")

	log.Println("search students")
	searchMessage(sqliteDatabase, "name LIKE 'Mar%' AND name LIKE '%Ben%' AND code > '0005'")
}

func createTable(db *sql.DB) {
	createStudentTableSQL := `CREATE TABLE student (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"from" TEXT,
		"body" TEXT,
		"message" TEXT		
	  );` // SQL Statement for Create Table

	log.Println("Create message table...")
	statement, err := db.Prepare(createStudentTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	log.Println("message table created")
}

// We are passing db reference connection from main to our method with other parameters
func insertMessage(db *sql.DB, msg message) {
	//log.Println("Inserting student record ...")
	insertStudentSQL := `INSERT INTO message(from, body, message) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertStudentSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(msg.From, msg.Body, msg.Message)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func displayMessages(db *sql.DB) {
	row, err := db.Query("SELECT * FROM student ORDER BY name")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var code string
		var name string
		var program string
		row.Scan(&id, &code, &name, &program)
		log.Println("Student: ", code, " ", name, " ", program)
	}
}

func searchMessage(db *sql.DB, query string) {
	row, err := db.Query(fmt.Sprintf("SELECT * FROM message WHERE %s ORDER BY id", query))
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	for row.Next() { // Iterate and fetch the records from result cursor
		var id int
		var from string
		var body string
		var message string
		row.Scan(&id, &from, &body, &message)
		log.Println("Msg: ", from, " ", body, " ", message)
	}
}
