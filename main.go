package main

type ListItem struct {
	Label string
	Path  string
	Final bool
}

type Client interface {
	List(string) ([]ListItem, error)
	Get(string) (string, error)
	HasParent(path string) bool
	Parent(path string) string
}

type Config struct {
	Client Client
	Path   string
}

func main() {
	presenter := NewPresenter()
	config, err := readConfig()
	if err != nil {
		panic("Can't read config: " + err.Error())
	}
	client := config.Client
	path := config.Path

	list, err := client.List(path)
	if err != nil {
		panic("Can't list: " + err.Error())
	}
	presenter.List(list)
}

func readConfig() (Config, error) {
	return Config{
		NewVaultClient(),
		"secret/",
	}, nil
}
