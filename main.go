package main

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const homePath = "/.lil"

func main() {

	screen, err := initScreen()

	if err != nil {
		panic("Can't create view: " + err.Error())
	}

	var once sync.Once
	gracefulClose := func() {
		once.Do(func() {
			screen.Fini()
		})
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, os.Kill, syscall.SIGTERM)
		<-sigChan
		gracefulClose()
		os.Exit(0)
	}()

	defer func() {
		if err != nil {
			if errors.Is(err, Quit) {
				return
			}
			log.Println("Exit with error: " + err.Error())
			os.Exit(1)
		}
	}()
	defer gracefulClose()

	view := NewView(screen)
	input := NewKeyboardInput(screen)

	for {
		var userHome string
		userHome, err = os.UserHomeDir()
		if err != nil {
			return
		}
		clientListController := NewClientListController(input, view, userHome+homePath)
		var client Client
		client, err = clientListController.Run()
		if err != nil {
			return
		}
		controller := NewController(input, view, client)
		err = controller.Run()
		if err != nil {
			return
		}
	}
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
