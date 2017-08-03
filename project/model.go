package project

type Project struct {
	Name        string  `json:"name"`
	Deployments []Deployment `json:"deployments"`
}

type Deployment struct {
	Stage   string         `json:"stage"`
	Version string      `json:"version"`
}
