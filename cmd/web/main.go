package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/wizenheimer/snippet.ly/internal/models"

	"github.com/go-sql-driver/mysql"
)

type application struct {
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	snippet       *models.SnippetModel
	address       string
	staticDir     string
	dsn           string
	templateCache map[string]*template.Template
}

func main() {
	// logger configuration
	infoLogger := log.New(os.Stdout, "INFO:\t", log.Ldate|log.Ltime)
	errorLogger := log.New(os.Stderr, "ERROR:\t", log.Ldate|log.Ltime|log.Lshortfile)

	// database configuration
	cfg := mysql.Config{
		User:                 "user",
		Passwd:               "password",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "snippetly",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	dsn := cfg.FormatDSN()

	infoLogger.Printf("Opening connection pool to database %s", dsn)
	db, err := openDB(dsn)
	if err != nil {
		errorLogger.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLogger.Fatal(err)
	}

	// application struct to share loggers with handlers
	app := &application{
		infoLogger:    infoLogger,
		errorLogger:   errorLogger,
		snippet:       &models.SnippetModel{DB: db},
		dsn:           dsn,
		templateCache: templateCache,
	}

	// server configurations
	flag.StringVar(&app.address, "addr", "localhost:4000", "HTTP network address")
	flag.StringVar(&app.staticDir, "static", "./ui/static", "Static Directory for Assets")
	flag.Parse()

	srv := &http.Server{
		Addr:     app.address,
		ErrorLog: errorLogger,
		Handler:  app.routes(),
	}

	infoLogger.Printf("Serving@ http://%s", app.address)
	err = srv.ListenAndServe()
	errorLogger.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
