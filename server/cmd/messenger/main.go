package main

import (
	"messenger/internal/app"
	"messenger/internal/config"
)

func main() {
	cfg := config.MustLoad()

	app.Run(cfg)
}
