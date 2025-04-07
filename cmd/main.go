package main

import (
	"ddd/infrastructure/config"
)

func main() {
	config.LoadConfig()
	r, _ := InitWeb()
	r.Run(config.Conf.App.Link())
}
