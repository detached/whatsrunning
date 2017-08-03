package dashboard

import (
	"log"
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
	"github.com/detached/whatsrunning/project"
	"github.com/detached/whatsrunning/config"
)

var (
	dashboardTemplate = template.Must(template.ParseFiles(config.DashboardTemplate))
)

func RegisterHandler(r *mux.Router) {
	log.Println("GET on /")
	r.HandleFunc("/", requestHandler).Methods("GET")
	r.PathPrefix("/dist").Handler(http.StripPrefix("/dist", http.FileServer(http.Dir(config.ContentDir))))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {

	var v view
	v.Projects = project.GetAll()
	v.WebsocketUrl = config.WebsocketUrl

	err := dashboardTemplate.Execute(w, v)

	if err != nil {
		log.Printf("Error while execute template %s: %s\n", config.DashboardTemplate, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

type view struct {
	Projects        []project.Project
	WebsocketUrl    string
}