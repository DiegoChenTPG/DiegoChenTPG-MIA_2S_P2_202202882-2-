package comandos

import (
	estructuras "PROYECTO2/estructuras"
	global "PROYECTO2/global"
	utilidades "PROYECTO2/utilidades"
	"errors" // Paquete para manejar errores y crear nuevos errores con mensajes personalizados
	"fmt"    // Paquete para formatear cadenas y realizar operaciones de entrada/salida

	"strings" // Paquete para manipular cadenas, como unir, dividir, y modificar contenido de cadenas
)

// MOUNT estructura que representa el comando mount con sus parámetros
type MOUNT struct {
	path string // Ruta del archivo del disco
	name string // Nombre de la partición
}

/*
	mount -path=/home/Disco1.mia -name=Part1 #id=341a
	mount -path=/home/Disco2.mia -name=Part1 #id=342a
	mount -path=/home/Disco3.mia -name=Part2 #id=343a
*/

// CommandMount parsea el comando mount y devuelve una instancia de MOUNT
func ParserMount(tokens []string) (*MOUNT, string, error) {
	cmd := &MOUNT{} // Crea una nueva instancia de MOUNT

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
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, "el path no puede estar vacio", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-name":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return nil, "el nombre no puede estar vacio", errors.New("el nombre no puede estar vacío")
			}
			cmd.name = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, "parametro desconocido: " + key, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -path y -name hayan sido proporcionados
	if cmd.path == "" {
		return nil, "faltan parametros requeridos: -path", errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return nil, "faltan parametros requeridos: -name", errors.New("faltan parámetros requeridos: -name")
	}

	// Montamos la partición
	posible_error, err := commandMount(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, posible_error, err
	}

	return cmd, "Se monto la particion: " + cmd.name + " del disco ubicado en: " + cmd.path + " en el sistema", nil // Devuelve el comando MOUNT creado
}

func commandMount(mount *MOUNT) (string, error) {

	for _, mountInfo := range global.MountedPartitions {
		fmt.Println("COMPARACION MOUNT NAME")
		fmt.Println(mountInfo.Name)
		fmt.Println(mount.name)
		if mountInfo.Name == mount.name {
			return "Error: La partición" + mountInfo.Name + " ya está montada", errors.New("la partición ya está montada")
		}
	}

	// Crear una instancia de MBR
	var mbr estructuras.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(mount.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return "Error deserializando el MBR", err
	}

	// Buscar la partición con el nombre especificado
	partition, indexPartition, _ := mbr.GetPartitionByName(mount.name)
	fmt.Println("IMPRIMIENDO INDEX PARTITION")
	fmt.Println(indexPartition + 1)

	// Generar un id único para la partición
	idPartition, err := GenerateIdPartition(mount, indexPartition+1) //Le sumamos 1 para que los IDs partition empiezen en 1 y luego tenemos: 821A
	if err != nil {
		fmt.Println("Error generando el id de partición:", err)
		return "Error generando el id de la partición", err
	}
	fmt.Println(idPartition + " IMPRIMIENDO ID PARTITION")
	//  Guardar la partición montada en la lista de montajes globales

	global.MountedPartitions[idPartition] = global.MountInfo{
		Path: mount.path,
		Name: mount.name,
	}
	//global.MountedPartitions[idPartition] = mount.path

	// Modificamos la partición para indicar que está montada
	partition.MountPartition(indexPartition, idPartition)

	// Guardar la partición modificada en el MBR
	mbr.Mbr_partitions[indexPartition] = *partition

	// Serializar la estructura MBR en el archivo binario
	err = mbr.Serialize(mount.path)
	if err != nil {
		fmt.Println("Error serializando el MBR:", err)
		return "Error serializando el MBR", err
	}

	fmt.Println("IMPRIMIENTO PARTICIONES EN EL MBR")
	mbr.PrintPartitions()
	return "", nil
}

func GenerateIdPartition(mount *MOUNT, indexPartition int) (string, error) {
	// Asignar una letra a la partición
	letter, err := utilidades.GetLetter(mount.path)
	if err != nil {
		fmt.Println("Error obteniendo la letra:", err)
		return "", err
	}

	// Crear id de partición
	idPartition := fmt.Sprintf("%s%d%s", utilidades.Carnet, indexPartition, letter)

	return idPartition, nil
}
