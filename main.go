package main

import "os"

type Feedback struct {
	Action Action
	Item   ListItem
}

type Action int

const (
	ListAction = Action(iota)
	GetAction
	ListParentAction
	QuitAction
)

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

type Controller interface {
	Run() Feedback
}

type Config struct {
	Client Client
	Path   string
}

func main() {
	view, err := NewView()
	if err != nil {
		panic(err)
	}
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

	controller := NewListController(list, view)
	handle(controller.Run(), func() {
		view.Close()
		os.Exit(0)
	})
}

func handle(feedback Feedback, quit func()) {
	switch feedback.Action {
	case ListAction:
	case GetAction:
	case ListParentAction:
	case QuitAction:
		quit()
	}
}

func readConfig() (Config, error) {
	return Config{
		NewVaultClient(),
		"secret/",
	}, nil
}
