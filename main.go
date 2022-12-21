package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Client interface {
	RootPath() Path
	List(Path) ([]Path, error)
	Get(Path) (Content, error)
}

func main() {
	client, err := buildClient()
	if err != nil {
		panic("Can't read config: " + err.Error())
	}
	screen, err := initScreen()

	if err != nil {
		panic("Can't create view: " + err.Error())
	}

	var once sync.Once
	close := func() {
		once.Do(func() {
			screen.Fini()
		})
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-sigChan
		close()
		os.Exit(0)
	}()

	defer func() {
		if err != nil {
			log.Println("Exit with error: " + err.Error())
			os.Exit(1)
		}
	}()
	defer close()

	view := NewView(screen)
	input := NewInput(screen)
	controller := NewController(input, view, client)
	err = controller.Run()
}

func initScreen() (tcell.Screen, error) {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("can't create presenter screen: %w", err)
	}
	if err = screen.Init(); err != nil {
		return nil, fmt.Errorf("can't create presenter screen: %w", err)
	}
	screen.Clear()
	return screen, nil
}

func buildClient() (Client, error) {
	client, err := NewVaultClient("secret/")
	return client, err
}
