package project

var s Storage

type Storage interface {
	Get(project string, stage string) (*Deployment, error)
    GetAll() ([]Project)
    Store(project string, stage string, deployment Deployment) error
}

func GetStorage() Storage {

	if s == nil {
		s = fileStorage{}
	}

	return s
}