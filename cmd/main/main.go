package main

import (
	"auth-service/internal/app"
	"flag"
)

func main() {
	configPath := flag.String("configPath", "./config.json", "")

	flag.Parse()

	app.Run(configPath)
}
