package main

import (
	"fmt"
	"github.com/sajimenezher_meli/meli-frescos-8/internal/application"
	"github.com/sajimenezher_meli/meli-frescos-8/pkg/loader"
	"os"
)

type test struct {
	Id   int
	Name string
}

func main() {
	app := application.Application{FolderDatabase: "docs/database"}
	app.SetEnvironment()

	buyerStruct1 := test{Id: 1, Name: "123451"}
	buyerStruct2 := test{Id: 7, Name: "124451"}
	buyerStructsSlice := []test{buyerStruct1, buyerStruct2}

	fileName := "text.json"
	storage := loader.NewJSONStorage[test](fmt.Sprintf("%s/%s", os.Getenv("folder_database"), fileName))

	err := storage.WriteAll(buyerStructsSlice)
	if err != nil {
		panic(err.Error())
	}

	buyers, err := storage.ReadAll()
	if err != nil {
		panic(err.Error())
	}

	// Nuevo objeto
	newObject := test{Id: 123, Name: "Samuel"}

	//Mapa con el nuevo objeto
	buyers[newObject.Id] = newObject

	newSlice := storage.MapToSlice(buyers)

	//Escribir nuevamente en la base de datos
	storage.WriteAll(newSlice)

	fmt.Println("buyers:")
	fmt.Println(buyers)

	buyer2 := buyers[7]
	fmt.Printf("Info buyer key 7, Id %d - Card %d", buyer2.Id, buyer2.Name)
}
