package main

import (
	"log"
	"net/http"
)

func main() {
	m, err := NewModel(1 << 8)
	if err != nil {
		panic(err)
	}

	// Run the model for ever
	m.Run()

	mux := http.NewServeMux()
	mux.Handle("/", m)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
