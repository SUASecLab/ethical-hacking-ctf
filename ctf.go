package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	sidecarUrl string

	mongoDBName   string
	mongoHostname string
	mongoPassword string
	mongoUsername string
)

func main() {
	sidecarUrl = os.Getenv("SIDECAR_URL")

	mongoDBName = os.Getenv("DB_NAME")
	mongoHostname = os.Getenv("DB_HOST")
	mongoPassword = os.Getenv("DB_PASSWORD")
	mongoUsername = os.Getenv("DB_USER")

	r := mux.NewRouter()
	r.HandleFunc("/", index)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Println("Starting up CTF service")
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatalln("CTF service failed: ", err)
	}
}
