package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type config struct {
	address   string
	staticDir string
}

func main() {
	// server configurations
	var cfg config
	flag.StringVar(&cfg.address, "addr", "localhost:4000", "HTTP network address")
	flag.StringVar(&cfg.staticDir, "static", "./ui/static", "Static Directory for Assets")
	flag.Parse()

	// logger configuration
	infoLogger := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	// multiplexer and fileserver configuration
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// routes
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	infoLogger.Printf("Serving@ http://%s", cfg.address)
	err := http.ListenAndServe(cfg.address, mux)
	errorLogger.Fatal(err)
}
