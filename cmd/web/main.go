package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	address     string
	staticDir   string
}

func main() {
	// logger configuration
	infoLogger := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	// application struct to share loggers with handlers
	app := &application{
		infoLogger:  infoLogger,
		errorLogger: errorLogger,
	}

	// server configurations
	flag.StringVar(&app.address, "addr", "localhost:4000", "HTTP network address")
	flag.StringVar(&app.staticDir, "static", "./ui/static", "Static Directory for Assets")
	flag.Parse()

	mux := app.routes()

	srv := &http.Server{
		Addr:     app.address,
		ErrorLog: errorLogger,
		Handler:  mux,
	}

	infoLogger.Printf("Serving@ http://%s", app.address)
	err := srv.ListenAndServe()
	errorLogger.Fatal(err)
}
