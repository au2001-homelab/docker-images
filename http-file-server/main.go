package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

func main() {
	addr := flag.String("addr", ":8000", "address to listen on")
	dir := flag.String("dir", "/shares", "directory to serve")
	flag.Parse()

	srv := &http.Server{
		Addr:              *addr,
		Handler:           http.FileServer(http.Dir(*dir)),
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("serving %s on %s", *dir, *addr)
	log.Fatal(srv.ListenAndServe())
}
