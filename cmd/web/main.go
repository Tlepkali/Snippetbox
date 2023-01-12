package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/Tlepkali/snippetbox/pkg/models/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Addr      string
	StaticDir string
}

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	snippets      *sqlite.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	cfg := new(Config)
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "./ui/static", "Path to static assets")
	flag.Parse()
	dsn := flag.String("dsn", "db/snippetbox?parseTime=true", "Sql Database connection")

	infoF, err := os.OpenFile("logs/info.log", os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	defer infoF.Close()

	errorF, err := os.OpenFile("logs/error.log", os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		log.Fatal(err)
	}
	defer errorF.Close()

	infoLog := log.New(infoF, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(errorF, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		snippets:      &sqlite.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     cfg.Addr,
		ErrorLog: errorLog,
		Handler:  app.routes(cfg),
	}

	infoLog.Printf("Starting server on %s", cfg.Addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
