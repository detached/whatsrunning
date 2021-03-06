package project

import (
	"encoding/json"
	"net/http"
	"log"
	"strings"

	"github.com/gorilla/mux"
)

func RegisterHandler(r *mux.Router) {
	log.Println("GET/PUT on /api/project/{project}/stage/{stage}")
	r.HandleFunc("/api/project/{project}/stage/{stage}", getDeploymentInfoHandler).Methods("GET")
	r.HandleFunc("/api/project/{project}/stage/{stage}", putDeploymentInfoHandler).Methods("PUT")
}

func getDeploymentInfoHandler(w http.ResponseWriter, r *http.Request) {
	var project, stage = parsePath(r)

	log.Printf("Get %s on stage %s \n", project, stage)

	var deployment, err = GetStorage().Get(project, stage)

	if err != nil {
		log.Printf("Error while getting %s on %s: %s \n", project, stage, err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(deployment)
}

func putDeploymentInfoHandler(w http.ResponseWriter, r *http.Request) {
	var project, stage = parsePath(r)
	var deployment Deployment

	if err := json.NewDecoder(r.Body).Decode(&deployment); err != nil {
		log.Printf("Error while reading body: %s \n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Save project %s on stage %s: %s\n", project, stage, deployment)

	deployment.Stage = stage

	if err := GetStorage().Store(project, stage, deployment); err != nil {
		log.Printf("Error while storing %s on %s: %s \n", project, stage, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func parsePath(r *http.Request) (string, string) {
	p := mux.Vars(r)
	return removeDangerousChars(p["project"]), removeDangerousChars(p["stage"])
}

func removeDangerousChars(text string) string {
	t := strings.Replace(text, ".", "", -1)
	t = strings.Replace(t, "/", "", -1)
	t = strings.Replace(t, "~", "", -1)
	return t
}