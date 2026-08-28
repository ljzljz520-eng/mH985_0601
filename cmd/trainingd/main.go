package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"trainingdesk/internal/api"
	"trainingdesk/internal/store"
)

func main() {
	dbPath := flag.String("db", "trainingdesk.db", "path to the bbolt database")
	addr := flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()
	s, err := store.Open(*dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()
	server := &http.Server{Addr: *addr, Handler: api.New(s).Handler()}
	fmt.Printf("trainingdesk listening on %s\n", *addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
