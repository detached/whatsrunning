package dashboard

import (
	"log"
	"net/http"
	"html/template"
	"github.com/gorilla/mux"
	"github.com/detached/whatsrunning/project"
	"github.com/detached/whatsrunning/config"
)

var dTemplate *template.Template

func RegisterHandler(r *mux.Router) {
	log.Println("GET on /")
	r.HandleFunc("/", requestHandler).Methods("GET")
	r.PathPrefix("/dist").Handler(http.StripPrefix("/dist", http.FileServer(http.Dir(config.D.ContentDir))))
}

func requestHandler(w http.ResponseWriter, r *http.Request) {

	if dTemplate == nil {
		dTemplate = template.Must(template.ParseFiles(config.D.DashboardTemplate))
	}

	if err := dTemplate.Execute(w, view{Projects: project.GetStorage().GetAll()}); err != nil {

		log.Printf("Error while execute template %s: %s\n", config.D.DashboardTemplate, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

type view struct {
	Projects        []project.Project
}