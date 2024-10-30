package comandos

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func ParserRmdisk(tokens []string) (string, error) {
	var ruta string

	// Itera sobre cada coincidencia encontrada

	for _, match := range tokens {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return "formato de parametro invalido: " + match, fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return "el path no puede estar vacio", errors.New("el path no puede estar vacío")
			}
			ruta = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "parametro desconocido: " + key, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	if ruta == "" {
		return "faltan parametros requeridos: -path", errors.New("faltan parámetros requeridos: -path")
	}

	// Abrir el archivo
	err := os.Remove(ruta)

	// Verificar si hubo algún error durante la lectura del archivo
	if err != nil {

		return "no se pudo eliminar el disco en la ruta: " + ruta + " el disco no existe", fmt.Errorf("no se pudo eliminar el archivo: %v", err)
	}

	return "SE ELIMINO CORRECTAMENTE EL DISCO EN LA RUTA: " + ruta, nil
}
