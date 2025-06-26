package application

import "os"

type Application struct {
	PortServer     string
	FolderDatabase string
}

func (app *Application) SetEnvironment() {

	if err := os.Setenv("folder_database", app.FolderDatabase); err != nil {
		panic("can't set environment")
	}

	if err := os.Setenv("port_server", app.PortServer); err != nil {
		panic("can't set environment")
	}
}

func (app *Application) InitApplication() {

}
