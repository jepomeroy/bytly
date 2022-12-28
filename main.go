package main

import (
	"bytly/model"
	"bytly/server"
)

func main() {
	model.Setup()
	server.Setup()
}
