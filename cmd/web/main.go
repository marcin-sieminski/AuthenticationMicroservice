package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

var (
	session *scs.SessionManager
	cfg     = config{api: "http://localhost:81"}
	app     = &application{
		config:  cfg,
		Session: scs.New(),
	}
)

type config struct {
	api string
}

type application struct {
	config  config
	Session *scs.SessionManager
}

type templateData struct {
	API string
}

func main() {
	fmt.Println("Starting front end service on port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Panic(err)
	}
}
