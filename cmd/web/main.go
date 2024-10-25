package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/marcin-sieminski/AuthenticationService/models"
)

var (
	session *scs.SessionManager
	counts  int64
)

type config struct {
	api  string
	port int
}

type templateData struct {
	API             string
	IsAuthenticated bool
	UserID          int
	UserName        string
}

type application struct {
	config   config
	Session  *scs.SessionManager
	DB       models.DBModel
	infoLog  *log.Logger
	errorLog *log.Logger
}

func main() {
	fmt.Println("Starting front end service on port 80")

	var cfg = config{
		api:  "http://localhost:81",
		port: 80}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Store = pgxstore.New(pool)

	conn := connectToDB()
	if conn == nil {
		log.Panic("Brak połączenia z bazą danych")
	}

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		config:   cfg,
		Session:  session,
		DB:       models.DBModel{DB: conn},
		infoLog:  infoLog,
		errorLog: errorLog,
	}

	go app.ListenToWsChannel()

	err = http.ListenAndServe(":80", app.routes())
	if err != nil {
		log.Panic(err)
	}
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 5 {
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(2 * time.Second)
		continue
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
