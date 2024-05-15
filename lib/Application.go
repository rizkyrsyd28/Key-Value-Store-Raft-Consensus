package lib

import "fmt"

type Application struct {
	data map[string]interface{}
}

func (app Application) Ping() {
	fmt.Println("Method Not Implemented")
}

func (app Application) Get() {
	fmt.Println("Method Not Implemented")
}

func (app Application) Set() {
	fmt.Println("Method Not Implemented")
}

func (app Application) Append() {
	fmt.Println("Method Not Implemented")
}

func (app Application) Del() {
	fmt.Println("Method Not Implemented")
}
