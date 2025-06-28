package main

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/application"
)

type test struct {
	Id   int
	Name string
}

func main() {
	app := application.Application{
		PortServer:     "localhost:8080",
		FolderDatabase: "docs/database",
	}

	app.SetEnvironment()
	app.InitApplication()
}
