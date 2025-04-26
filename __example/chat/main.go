package main

import (
	"github.com/abolfazlalz/emitter"
	"log"
	"net/http"
)

func main() {
	srv := emitter.NewServer(nil)
	go srv.StartListen()
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		srv.Handler(writer, request, nil)
	})

	fs := http.FileServer(http.Dir("./static"))

	// Serve the static files when accessing the root URL ("/")
	http.Handle("/", fs)

	log.Printf("starting server on :8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
