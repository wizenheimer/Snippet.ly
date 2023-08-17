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

type application struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
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

	// application struct to share loggers with handlers
	app := &application{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
	}

	// multiplexer and fileserver configuration
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir(cfg.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// routes
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	srv := &http.Server{
		Addr:     cfg.address,
		ErrorLog: errorLogger,
		Handler:  mux,
	}

	infoLogger.Printf("Serving@ http://%s", cfg.address)
	err := srv.ListenAndServe()
	errorLogger.Fatal(err)
}
