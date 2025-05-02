package container

type Config struct {
	Name    string
	Command []string
}

type Container struct {
	ID     string
	Name   string
	Status string
}

func Run(config *Config) error {
	return nil
}
