package main

import (
	"blog/cmd/server"

	"github.com/lpernett/godotenv"
)

func main() {
	godotenv.Load()
	server.Execute()
}
