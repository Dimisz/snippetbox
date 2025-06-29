package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/dimisz/snippetbox/internal/models"
	"github.com/go-playground/form/v4"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorLog       *log.Logger
	infoLog        *log.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	templateCache, err := newTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}
	// initialize a decoder instace
	formDecoder := form.NewDecoder()

	// use scs.New() to initialize a new session manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	// set Secure attribute on session cookie
	// so user browser sends cookie only via HTTPS
	sessionManager.Cookie.Secure = true

	// initialize our application
	app := &application{
		errorLog:       errorLog,
		infoLog:        infoLog,
		snippets:       &models.SnippetModel{DB: db},
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// tlsConfig struct to hold non-default TLS settings
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}

	// custom server struct to channel all error logs thru custom logger
	srv := &http.Server{
		Addr:      *addr,
		ErrorLog:  errorLog,
		Handler:   app.routes(),
		TLSConfig: tlsConfig,
		// Idle, Read, and Write timeouts
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	// err = srv.ListenAndServe()
	// for HTTPS use ListenAndServeTLS
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")

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
