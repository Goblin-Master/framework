package main

import (
	"framework/infrastructure/config"
)

func main() {
	config.LoadConfig()
	app, err := InitWeb()
	if err != nil {
		panic(err)
	}
	app.Run()
}
