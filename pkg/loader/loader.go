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
	// Usamos os.OpenFile con las banderas O_RDONLY (solo lectura), O_CREATE (crear si no existe)
	// y 0644 para los permisos del archivo (lectura/escritura para el propietario, solo lectura para otros).
	file, err := os.OpenFile(s.filepath, os.O_RDONLY|os.O_CREATE, 0644)

	if err != nil {
		return nil, fmt.Errorf("error abriendo o creando el archivo: %w", err)
	}
	defer file.Close()

	// Obtener el tamaño del archivo para verificar si está vacío.
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo información del archivo: %w", err)
	}

	// Si el archivo está vacío, devuelve un mapa vacío para indicar que no hay datos.
	if fileInfo.Size() == 0 {
		return make(map[int]T), nil
	}

	var list []T
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&list); err != nil {
		// Manejar el error si el JSON no es válido, pero el archivo sí existe.
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
