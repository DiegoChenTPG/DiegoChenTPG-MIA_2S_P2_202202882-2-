package comandos

import (
	estructuras "PROYECTO2/estructuras"
	global "PROYECTO2/global"
	"errors"
	"fmt"
	"strings"
)

// MOUNT estructura que representa el comando mount con sus parámetros
type UNMOUNT struct {
	id string
}

/*
	unmount -id=821A
*/

// CommandMount parsea el comando mount y devuelve una instancia de MOUNT
func ParserUnmount(tokens []string) (*UNMOUNT, string, error) {
	cmd := &UNMOUNT{} // Crea una nueva instancia de MOUNT

	// Itera sobre cada coincidencia encontrada

	for _, match := range tokens {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, "formato de parámetro inválido: " + match, fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {
		case "-id":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, "el id no puede estar vacio", errors.New("el path no puede estar vacío")
			}
			cmd.id = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, "parametro desconocido: " + key, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -path y -name hayan sido proporcionados
	if cmd.id == "" {
		return nil, "faltan parametros requeridos: -id", errors.New("faltan parámetros requeridos: -path")
	}

	// Montamos la partición
	posible_error, err := commandUnmount(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, posible_error, err
	}

	return cmd, "Se desmonto la particion con el id: " + cmd.id, nil // Devuelve el comando MOUNT creado
}

func commandUnmount(unmount *UNMOUNT) (string, error) {
	fmt.Println("ESTAMOS EN commandoUnmount")
	// Verifica si la partición está montada buscando el ID en la lista global
	mountInfo, exists := global.MountedPartitions[unmount.id]
	if !exists {
		return "Error: la partición no está montada o el ID es incorrecto", errors.New("la partición no está montada o el ID es incorrecto")
	}

	// Obtener el índice o información relevante del MBR/EBR, si aplica
	var mbr estructuras.MBR
	err := mbr.Deserialize(mountInfo.Path)
	if err != nil {
		return "Error deserializando el MBR para desmontar", err
	}

	// Buscar la partición específica en el MBR según su ID o índice
	partition, indexPartition, _ := mbr.GetPartitionByName(mountInfo.Name)

	// Restablecer correlativo u otros valores de la partición
	partition.ResetPartition() // `Reset` debe ser una función para restaurar el estado inicial de la partición
	// Actualizar el MBR y eliminar la partición de los montajes globales
	mbr.Mbr_partitions[indexPartition] = *partition
	err = mbr.Serialize(mountInfo.Path)
	if err != nil {
		return "Error al actualizar el MBR tras desmontar", err
	}

	err2 := global.DeleteMountedPartition(unmount.id)
	if err2 != nil {
		return "Error al eliminar la partición de la lista de montajes", err
	}

	fmt.Println("IMPRIMIENTO PARTICIONES EN EL MBR en unmount")
	mbr.PrintPartitions()

	return "", nil
}
