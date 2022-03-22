package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/blevesearch/bleve"
)

type message struct {
	Id      string
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
	path := "r:/example.bleve"
	os.RemoveAll(path)
	// open a new index
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(path, mapping)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	count := 0
	// index some data
	for x := 0; x < 100; x++ {
		fmt.Printf(".")
		batch := index.NewBatch()
		for i := 0; i < 1000; i++ {
			mymessage := message{
				Id:      fmt.Sprintf("ex%.8d", count),
				From:    "marty.schoch@gmail.com",
				Body:    fmt.Sprintf("bleve indexing is easy %d", count),
				Message: randSeq(128),
			}
			err = batch.Index(mymessage.Id, mymessage)
			count++
			if err != nil {
				log.Fatalf("error: %v", err)
			}
		}
		err := index.Batch(batch)
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}
	fmt.Println()

	// search for some text
	query := bleve.NewMatchQuery("easy")
	search := bleve.NewSearchRequest(query)
	searchResults, err := index.Search(search)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("result: %v", searchResults)

	query = bleve.NewMatchQuery("99999")

	search = bleve.NewSearchRequest(query)
	searchResults, err = index.Search(search)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("result: %v", searchResults)
}
