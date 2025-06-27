package main

import (
	"github.com/sajimenezher_meli/meli-frescos-8/internal/application"
)

type test struct {
	Id   int
	Name string
}

func main() {
	app := application.Application{FolderDatabase: "docs/database"}
	app.SetEnvironment()

	//prueba con Test Struct
	/*
		testStruct1 := test{Id: 1, Name: "Test"}
		testStruct2 := test{Id: 2, Name: "Test"}
		testsStructsSlice := []test{testStruct1, testStruct2}

		fileName := "text.json"
		storage := loader.NewJSONStorage[test](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), fileName))

		err := storage.WriteAll(testsStructsSlice)
		if err != nil {
			panic(err.Error())
		}



		fileName := "text.json"
		storage := loader.NewJSONStorage[test](
			fmt.Sprintf("%s/%s", os.Getenv("folder_database"), fileName),
		)

		data, err := storage.ReadAll()
		if err != nil {
			fmt.Println(err)
			return

		}
		fmt.Println(data)
	*/

}
