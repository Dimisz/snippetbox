package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
}

func main() {
	// command-line flag to customize port
	addr := flag.String("addr", ":4000", "HTTP network address")
	// command-line flag for MySQL DSN str
	dsn := flag.String("dsn", "username:password@/snippetbox?parseTime=true", "MySQL datasource name")
	flag.Parse()

	// custom logger
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	// initialize our application
	app := &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	// custom server struct to channel all error logs thru custom logger
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	// pool initialized but no connections yet
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	// that's why need to ping to verify
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
