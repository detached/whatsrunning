package project

import (
	"log"
	"os"
	"io/ioutil"
	"path/filepath"
	"encoding/json"
	"github.com/detached/whatsrunning/config"
)

type fileStorage struct{}

func (f fileStorage) Get(project string, stage string) (*Deployment, error) {

	var filePath = f.getDeploymentFile(project, stage)

	file, err := os.Open(filePath)

	if err != nil {
		log.Printf("Error while open file %s: %s", filePath, err)
		return nil, err
	}

	var deployment Deployment

	if err = json.NewDecoder(file).Decode(&deployment); err != nil {
		log.Printf("Error while reading file %s: %s", filePath, err)
		return nil, err
	}

	return &deployment, nil
}

func (f fileStorage) GetAll() ([]Project) {

	entries, err := ioutil.ReadDir(config.D.StoragePath)

	if err != nil {
		log.Fatal("Can't read storage dir: ", err)
	}

	var projects []Project

	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		if entry.IsDir() {
			projects = append(projects, f.createProject(entry))
		}
	}

	return projects
}

func (f fileStorage) Store(project string, stage string, deployment Deployment) error {

	var deploymentFile = f.getDeploymentFile(project, stage)

	if err := os.MkdirAll(filepath.Dir(deploymentFile), 0760); err != nil {
		log.Printf("Error while creating dir: %s\n", err)
		return err
	}

	file, err := os.Create(deploymentFile)

	if err != nil {
		log.Printf("Error while creating file: %s\n", err)
		return err
	}

	defer file.Close()

	if err = json.NewEncoder(file).Encode(deployment); err != nil {
		log.Printf("Cant write file:%s\n", err)
		return err
	}

	select {
	case UpdateChannel <- f.wrapWithProject(project, deployment):
	default:
	}

	return nil
}

func (f fileStorage) wrapWithProject(projectName string, deployment Deployment) Project {

	var p Project
	p.Name = projectName
	p.Deployments = append(p.Deployments, deployment)
	return p
}

func (f fileStorage) createProject(fi os.FileInfo) Project {

	var p Project

	p.Name = fi.Name()

	projectDir := filepath.Join(config.D.StoragePath, fi.Name())
	entries, _ := ioutil.ReadDir(projectDir)

	for i := 0; i < len(entries); i++ {
		entry := entries[i]
		if !entry.IsDir() {

			d, err := f.createDeployment(filepath.Join(projectDir, entry.Name()))

			if err == nil {
				p.Deployments = append(p.Deployments, *d)
			} else {
				log.Printf("Error creating deployment: %s", err)
			}
		}
	}

	return p
}

func (fileStorage) createDeployment(filePath string) (*Deployment, error) {

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0760)

	if err != nil {
		return nil, err
	}

	var deployment Deployment
	err = json.NewDecoder(file).Decode(&deployment)

	return &deployment, err
}

func (fileStorage) getDeploymentFile(project string, stage string) string {
	return filepath.Join(config.D.StoragePath, project, stage+".json")
}
