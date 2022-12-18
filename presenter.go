package main

import "fmt"

type Presenter struct {
}

func NewPresenter() *Presenter {
	return &Presenter{}
}

func (p *Presenter) List(list []ListItem) {
	for _, item := range list {
		fmt.Println(item.Label)
	}
}
