package main

import (
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/application"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
)

type test struct {
	Id   int
	Name string
}

func main() {
	app := application.Application{FolderDatabase: "docs/database"}
	app.SetEnvironment()

	/*
		testStruct := test{Id: 1, Name: "Test"}
		fileName := "text.json"
		err := loader.WriterFile(fileName, testStruct)
		if err != nil {
			panic(err.Error())
		}
	*/

	fileName := "text.json"

	data := loader.ReadFile[test](fileName)

	fmt.Println(data)

}
