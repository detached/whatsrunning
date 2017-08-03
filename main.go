package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/detached/whatsrunning/project"
	"github.com/detached/whatsrunning/dashboard"
)

func main() {

	log.Println("== WhatsRunning? ==")

	router := mux.NewRouter()

	dashboard.RegisterHandler(router)
	project.RegisterHandler(router)
	project.RegisterWebsocket(router)

	var server = ":8080"

	log.Printf("Starting on %s\n", server)

	log.Fatal(http.ListenAndServe(server, router))
}
