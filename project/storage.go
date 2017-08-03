package project

import (
	"log"
	"os"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"github.com/detached/whatsrunning/config"
)

func get(project string, stage string) (*Deployment, error) {

	var filePath = getDeploymentFile(project, stage)

	f, err := os.Open(filePath)

	if err != nil {
		log.Printf("Error while open file %s: %s", filePath, err)
		return nil, err
	}

	var deployment Deployment
	err = json.NewDecoder(f).Decode(&deployment)

	if err != nil {
		log.Printf("Error while reading file %s: %s", filePath, err)
		return nil, err
	}

	return &deployment, nil
}

func GetAll() ([]Project) {

	entries, err := ioutil.ReadDir(config.StoragePath)

	if err != nil {
		log.Fatal("Can't read storage dir: ", err)
	}

	var projects []Project

	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		if entry.IsDir() {
			projects = append(projects, createProject(entry))
		}
	}

	return projects
}

func store(project string, stage string, deployment Deployment) error {

	var deploymentFile = getDeploymentFile(project, stage)
	err := os.MkdirAll(filepath.Dir(deploymentFile), 0760)

	if err != nil {
		log.Printf("Error while creating dir: %s\n", err)
		return err
	}

	f, err := os.Create(deploymentFile)

	if err != nil {
		log.Printf("Error while creating file: %s\n", err)
		return err
	}

	defer f.Close()

	err = json.NewEncoder(f).Encode(deployment)

	if err != nil {
		log.Printf("Cant write file:%s\n", err)
		return err
	}

	UpdateChannel <- wrapWithProject(project, deployment)

	return nil
}

func wrapWithProject(projectName string, deployment Deployment) Project {

	var p Project
	p.Name = projectName
	p.Deployments = append(p.Deployments, deployment)
	return p
}

func createProject(fi os.FileInfo) Project {

	var p Project

	p.Name = fi.Name()

	projectDir := filepath.Join(config.StoragePath, fi.Name())
	entries, _ := ioutil.ReadDir(projectDir)

	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		if !entry.IsDir() {

			d, err := createDeployment(filepath.Join(projectDir, entry.Name()))

			if err == nil {
				p.Deployments = append(p.Deployments, *d)
			} else {
				log.Printf("Error creating deployment: %s", err)
			}
		}
	}

	return p
}

func createDeployment(filePath string) (*Deployment, error) {

	f, err := os.OpenFile(filePath, os.O_RDONLY, 0760)

	if err != nil {
		return nil, err
	}

	var deployment Deployment
	err = json.NewDecoder(f).Decode(&deployment)

	return &deployment, err
}

func getDeploymentFile(project string, stage string) string {
	return filepath.Join(config.StoragePath, project, stage+".json")
}