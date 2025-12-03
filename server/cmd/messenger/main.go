package main

import (
	"messanger/internal/app"
	"messanger/internal/config"
)

func main() {
	cfg := config.MustLoad()

	app.Run(cfg)
}
