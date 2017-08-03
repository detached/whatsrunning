package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/detached/whatsrunning/config"
	"github.com/detached/whatsrunning/project"
	"github.com/detached/whatsrunning/dashboard"
	"os"
)

func main() {

	var confFile string

	if len(os.Args) < 2 {
		confFile = "config.json"
	} else {
		confFile = os.Args[1]
	}

	log.Println("== WhatsRunning? ==")

	config.LoadConfig(confFile)

	router := mux.NewRouter()

	dashboard.RegisterHandler(router)
	project.RegisterHandler(router)
	project.RegisterWebsocket(router)

	log.Println("Starting on", config.D.Server)

	log.Fatal(http.ListenAndServe(config.D.Server, router))
}
