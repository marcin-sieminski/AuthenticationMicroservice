package main

import (
	"fmt"
	"log"
	"net/http"
)

var (
	cfg = config{
		api:  "http://localhost:81",
		port: 80}
	app = &application{
		config: cfg,
	}
)

type config struct {
	api  string
	port int
}

type application struct {
	config config
}

type templateData struct {
	API             string
	IsAuthenticated int
}

func main() {
	fmt.Println("Starting front end service on port 80")
	err := http.ListenAndServe(":80", app.routes())
	if err != nil {
		log.Panic(err)
	}
}
