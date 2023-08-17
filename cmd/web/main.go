package main

import (
	"flag"
	"log"
	"net/http"
)

type config struct {
	address   string
	staticDir string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.address, "addr", "localhost:4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static", "./ui/static", "Static Directory for Assets")
	flag.Parse()

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Printf("Serving@ http://%s", cfg.address)

	err := http.ListenAndServe(cfg.address, mux)
	log.Fatal(err)
}
