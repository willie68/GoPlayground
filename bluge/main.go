package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/blugelabs/bluge"
)

const path = "./data/bluge"

var config bluge.Config

func main() {
	config = bluge.DefaultConfig(path)

	write()

	search()
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func write() {
	writer, err := bluge.OpenWriter(config)
	if err != nil {
		log.Fatalf("error opening writer: %v", err)
	}
	defer writer.Close()
	count := 0
	// index some data
	batch := bluge.NewBatch()
	for x := 0; x < 100; x++ {
		fmt.Printf(".")
		for i := 0; i < 1000; i++ {
			doc := bluge.NewDocument(fmt.Sprintf("ex%.8d", count)).
				AddField(bluge.NewTextField("From", "marty.schoch@gmail.com").StoreValue()).
				AddField(bluge.NewTextField("Body", fmt.Sprintf("bleve indexing is easy %d", count)).StoreValue()).
				AddField(bluge.NewTextField("Message", randSeq(128)).StoreValue())

			batch.Update(doc.ID(), doc)
			count++
		}
		err := writer.Batch(batch)
		batch.Reset()
		if err != nil {
			log.Fatalf("error updating document: %v", err)
		}
	}
	fmt.Println()
}

// reading and searching
func search() {
	reader, err := bluge.OpenReader(config)
	if err != nil {
		log.Fatalf("error getting index reader: %v", err)
	}
	defer reader.Close()

	query := bluge.NewMatchQuery("bleve").SetField("Body")
	request := bluge.NewTopNSearch(200000, query).
		WithStandardAggregations()
	documentMatchIterator, err := reader.Search(context.Background(), request)
	if err != nil {
		log.Fatalf("error executing search: %v", err)
	}
	count := 0
	match, err := documentMatchIterator.Next()
	for err == nil && match != nil {
		err = match.VisitStoredFields(func(field string, value []byte) bool {
			if field == "_id" {
				fmt.Printf("match: %s\n", string(value))
			}
			if field == "Body" {
				fmt.Printf("match: %s\n", string(value))
			}
			if field == "Message" {
				fmt.Printf("match: %s\n", string(value))
			}
			return true
		})
		if err != nil {
			log.Fatalf("error loading stored fields: %v", err)
		}
		match, err = documentMatchIterator.Next()
		count++
	}
	if err != nil {
		log.Fatalf("error iterator document matches: %v", err)
	}
	log.Printf("found %d docs", count)
}
