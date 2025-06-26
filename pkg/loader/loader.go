package loader

import (
	"encoding/json"
	"fmt"
	"os"
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
func (s *StorageJSON[T]) ReadAll() ([]T, error) {
	// 1. Leer el archivo
	data, err := os.ReadFile(s.filepath)
	if err != nil {
		return nil, fmt.Errorf("Error al leer el archivo %s , %w", s.filepath, err)
	}
	//2. Deserializar el archivo
	var items []T
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, fmt.Errorf("Error al deserializar el archivo %s : %w", s.filepath, err)
	}
	return items, nil
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
		return fmt.Errorf("Error al serializar elementos :%w", err)
	}
	// Escribimos el archivo para aplicar los cambios
	if err := os.WriteFile(s.filepath, data, 0644); err != nil {
		return fmt.Errorf("Error escribiendo el archivo :%w ", err)
	}
	return nil
}
