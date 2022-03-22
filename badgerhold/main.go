package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/timshannon/badgerhold/v4"
)

type Blob struct {
	Properties map[string]interface{}
}

func main() {
	insert := false

	options := badgerhold.DefaultOptions
	options.Dir = "data"
	options.ValueDir = "data"

	store, err := badgerhold.Open(options)
	if err != nil {
		// handle error
		log.Fatal(err)
	}
	defer store.Close()
	bd := Blob{
		Properties: make(map[string]interface{}),
	}
	bd.Properties["X-tenant"] = "Murks"

	if insert {
		err = store.Insert("key", &bd)
		if err != nil {
			// handle error
			log.Fatal(err)
		}
	}

	var result []Blob
	store.Find(&result, badgerhold.Where(badgerhold.Key).Contains("key"))
	for _, bd := range result {
		js, _ := json.Marshal(bd)
		fmt.Printf("Result: %s\r\n", js)
	}
	fmt.Println("-----")
	store.Find(&result, badgerhold.Where("Properties").HasKey("X-tenant").And("Properties").Contains("Murks"))
	for _, bd := range result {
		js, _ := json.Marshal(bd)
		fmt.Printf("Result: %s\r\n", js)
	}
}
