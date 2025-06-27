package loader

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

// Storage define la interfaz genérica para persistencia de datos
type Storage[T any] interface {
	// ReadAll lee todos los elementos del almacenamiento
	ReadAll() ([]T, error)

	// WriteAll escribe todos los elementos al almacenamiento
	WriteAll(items []T) error
}

// StorageJSON implementa Storage[T] para archivos JSON
type StorageJSON[T any] struct {
	filepath string
}

// NewJSONStorage crea una nueva instancia de StorageJSON genérica
func NewJSONStorage[T any](filePath string) *StorageJSON[T] {
	return &StorageJSON[T]{
		filepath: filePath,
	}
}

// ReadAll implementa Storage[T].ReadAll
func (s *StorageJSON[T]) ReadAll() (map[int]T, error) {
	file, err := os.Open(s.filepath)
	if err != nil {
		return nil, fmt.Errorf("error abriendo el archivo: %w", err)
	}
	defer file.Close()

	var list []T
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&list); err != nil {
		return nil, fmt.Errorf("error decodificando JSON: %w", err)
	}

	result := make(map[int]T)
	for _, item := range list {
		// Usamos reflection para acceder al campo `Id`
		val := reflect.ValueOf(item)
		if val.Kind() == reflect.Struct {
			idField := val.FieldByName("Id")
			if idField.IsValid() && idField.Kind() == reflect.Int {
				id := int(idField.Int())
				result[id] = item
			} else {
				return nil, fmt.Errorf("el campo 'Id' no es válido o no es de tipo int en %v", item)
			}
		} else {
			return nil, fmt.Errorf("el tipo no es un struct válido")
		}
	}

	return result, nil
}

// WriteAll implementa Storage[T].WriteAll
// El tercer argumento (0644) son los permisos del archivo en formato octal:
// - 6: permisos de lectura y escritura para el propietario
// - 4: permisos de solo lectura para el grupo
// - 4: permisos de solo lectura para otros usuarios
func (s *StorageJSON[T]) WriteAll(items []T) error {
	// En este caso vamos a serializar los elementos a un formato json legible
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return fmt.Errorf("error al serializar elementos :%w", err)
	}
	// Escribimos el archivo para aplicar los cambios
	if err := os.WriteFile(s.filepath, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo el archivo :%w ", err)
	}
	return nil
}

// MapToSlice converts a generic map[int]T to a generic slice []T.
// It iterates through the values of the input map and appends them
// to a new slice, which is then returned.
func (s *StorageJSON[T]) MapToSlice(items map[int]T) []T {

	var itemsSlice = make([]T, 0)

	for _, value := range items {
		// Append each value from the map to the itemsSlice.
		itemsSlice = append(itemsSlice, value)
	}

	// Return the newly created slice containing all the values from the map.
	return itemsSlice
}
