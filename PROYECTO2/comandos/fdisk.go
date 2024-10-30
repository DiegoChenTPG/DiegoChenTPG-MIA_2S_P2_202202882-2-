package comandos

import (
	estructuras "PROYECTO2/estructuras"
	utilidades "PROYECTO2/utilidades"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type FDISK struct {
	size int
	unit string
	fit  string
	path string
	typ  string
	name string
}

func ParserFdisk(tokens []string) (*FDISK, string, error) {

	cmd := &FDISK{}

	// Itera sobre cada coincidencia encontrada
	for _, match := range tokens {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		if len(kv) != 2 {
			return nil, "formato de parametro invalido: " + match, fmt.Errorf("formato de parámetro inválido: %s", match)
		}
		key, value := strings.ToLower(kv[0]), kv[1]

		// Remove quotes from value if present
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = strings.Trim(value, "\"")
		}

		// Switch para manejar diferentes parámetros
		switch key {

		case "-size":
			// Convierte el valor del tamaño a un entero
			size, err := strconv.Atoi(value)
			if err != nil || size <= 0 {
				return nil, "el tamaño debe ser un número entero positivo", errors.New("el tamaño debe ser un número entero positivo")
			}
			cmd.size = size
		case "-unit":
			// Verifica que la unidad sea "K" o "M"
			value = strings.ToUpper(value)
			if value != "K" && value != "M" && value != "B" {
				return nil, "la unidad debe ser K o M", errors.New("la unidad debe ser K, M o B")
			}
			cmd.unit = strings.ToUpper(value)
		case "-fit":
			// Verifica que el ajuste sea "BF", "FF" o "WF"
			value = strings.ToUpper(value)
			if value != "BF" && value != "FF" && value != "WF" {
				return nil, "el ajuste debe ser BF, FF o WF", errors.New("el ajuste debe ser BF, FF o WF")
			}
			cmd.fit = value
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, "el path no puede estar vacío", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-type":
			// Verifica que el tipo sea "P", "E" o "L"
			value = strings.ToUpper(value)
			if value != "P" && value != "E" && value != "L" {
				return nil, "el tipo debe ser P, E o L", errors.New("el tipo debe ser P, E o L")
			}
			cmd.typ = value
		case "-name":
			// Verifica que el nombre no esté vacío
			fmt.Println(value + " nombre de la particion")
			if value == "" {
				return nil, "el nombre no puede estar vacío", errors.New("el nombre no puede estar vacío")
			}
			cmd.name = value
		case "-add":
			fmt.Println("implementar -add")
		case "-delete":
			fmt.Println("implementar -delete")
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, "parametro desconocido: " + key, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -size, -path y -name hayan sido proporcionados
	if cmd.size == 0 {
		return nil, "faltan parámetros requeridos: -size", errors.New("faltan parámetros requeridos: -size")
	}
	if cmd.path == "" {
		return nil, "faltan parámetros requeridos: -path", errors.New("faltan parámetros requeridos: -path")
	}
	if cmd.name == "" {
		return nil, "faltan parámetros requeridos: -name", errors.New("faltan parámetros requeridos: -name")
	}

	// Si no se proporcionó la unidad, se establece por defecto a "M"
	if cmd.unit == "" {
		cmd.unit = "K"
	}

	// Si no se proporcionó el ajuste, se establece por defecto a "FF"
	if cmd.fit == "" {
		cmd.fit = "WF"
	}

	// Si no se proporcionó el tipo, se establece por defecto a "P"
	if cmd.typ == "" {
		cmd.typ = "P"
	}
	var tipo_particion string
	// Crear la partición con los parámetros proporcionados
	tipo_particion, err := commandFdisk(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return cmd, "Error al crear la particion en: " + cmd.path, nil
	}

	return cmd, "Se creo una particion de tipo: " + tipo_particion + " llamada: " + cmd.name + " en el disco con la ruta: " + cmd.path, nil // Devuelve el comando FDISK creado

}

func commandFdisk(fdisk *FDISK) (string, error) {
	// Convertir el tamaño a bytes
	var tipo_particion string
	sizeBytes, err := utilidades.ConvertToBytes(fdisk.size, fdisk.unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return "", err
	}

	if fdisk.typ == "P" {
		// Crear partición primaria
		tipo_particion, err = createPrimaryPartition(fdisk, sizeBytes)
		if err != nil {
			fmt.Println("Error creando partición primaria:", err)
			return "Error creando particion primaria", err
		}
	} else if fdisk.typ == "E" {
		fmt.Println("Creando partición extendida...")
		/*tipo_particion, */ err = createExtendedPartition2(fdisk, sizeBytes)
		tipo_particion = "E"
		if err != nil {

			fmt.Println("Error creando la particion extendida: ", err)
			return "Error creando la particion extendida", err
		}

	} else if fdisk.typ == "L" {
		fmt.Println("Creando partición lógica...") // Les toca a ustedes implementar la partición lógica
		/*tipo_particion, */ err = createLogicPartition2(fdisk, sizeBytes)
		tipo_particion = "L"
		if err != nil {

			fmt.Println("Error creando la particion extendida: ", err)
			return "Error creando la particion extendida", err
		}
	}

	return tipo_particion, nil
}

func createPrimaryPartition(fdisk *FDISK, sizeBytes int) (string, error) {
	// Crear una instancia de MBR
	var mbr estructuras.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return "Error deserealizando el MBR", err
	}

	// Obtener la primera partición disponible
	availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition()
	if availablePartition == nil {
		fmt.Println("No hay particiones disponibles.")
		return "No hay particiones disponibles", errors.New("no hay particiones disponibles")
	}

	/* SOLO PARA VERIFICACIÓN */
	// Print para verificar que la partición esté disponible

	fmt.Println("\nPartición disponible:")
	availablePartition.Print()

	// Crear la partición con los parámetros proporcionados
	availablePartition.CreatePartition(startPartition, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name)

	// Print para verificar que la partición se haya creado correctamente

	fmt.Println("\nPartición creada (modificada):")
	availablePartition.Print()
	fmt.Println("PRIMERA PASADA")
	fmt.Println(indexPartition)
	// Colocar la partición en el MBR
	if availablePartition != nil {
		mbr.Mbr_partitions[indexPartition] = *availablePartition
	}

	// Imprimir las particiones del MBR

	fmt.Println("\nParticiones del MBR:")
	mbr.PrintPartitions()
	fmt.Println("AQUI SE TERMINA DE IMPRIMIR LA PARTE DE LA PARTICION PRIMARIA")
	fmt.Println("==============================")

	// Serializar el MBR en el archivo binario
	err = mbr.Serialize(fdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return "Primaria", nil
}

func createExtendedPartition2(fdisk *FDISK, sizeBytes int) error {
	// Crear una instancia de MBR
	var mbr estructuras.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	// Veremos si existe una partición extendida en el disco
	extendedPartition, startExtended, indexExtended := mbr.GetExtended()
	fmt.Println("Particion" + string(startExtended) + string(indexExtended))

	if extendedPartition == nil {
		// Obtener la primera partición disponible
		availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition()
		if availablePartition == nil {
			return errors.New("no hay particiones disponibles")
		}

		// Crear la partición con los parámetros proporcionados
		availablePartition.CreatePartition(startPartition, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name)

		// Colocar la partición en el MBR
		if availablePartition != nil {
			mbr.Mbr_partitions[indexPartition] = *availablePartition
		}

		// Serializar el MBR en el archivo binario
		err = mbr.Serialize(fdisk.path)
		if err != nil {
			fmt.Println("Error:", err)
		}

		err = createEBR2(fdisk, int64(startPartition))
		if err != nil {
			fmt.Println("Error creating EBR:", err)
			return err
		}

	} else {
		return errors.New("ya existe una partición extendida en el disco")
	}

	// Serializar el MBR en el archivo binario
	err = mbr.Serialize(fdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}

func createLogicPartition2(fdisk *FDISK, sizeBytes int) error {
	// Crear una instancia de MBR
	var mbr estructuras.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return err
	}

	// Veremos si existe una partición extendida en el disco
	extendedPartition, startExtended, indexExtended := mbr.GetExtended()
	fmt.Println("Particion" + string(startExtended) + string(indexExtended))

	offset := int64(startExtended)
	if extendedPartition != nil {
		for i := 0; i < 10; i++ {
			// Crear una instancia de MBR
			var ebr estructuras.EBR2

			// Deserializar la estructura MBR desde un archivo binario
			err := ebr.Deserialize2(fdisk.path, offset)
			if err != nil {
				fmt.Println("Error deserializando el MBR:", err)
				return err
			}

			if ebr.Ebr_next == -1 {
				inicio := ebr.Ebr_start + ebr.Ebr_size
				err = createEBR2(fdisk, int64(inicio))
				if err != nil {
					fmt.Println("Error creating EBR:", err)
					return err
				}

				ebr.Ebr_next = inicio

				// Serializar el MBR en el archivo binario
				err = ebr.Serialize2(fdisk.path, offset)
				if err != nil {
					fmt.Println("Error:", err)
				}

				break

			}
			offset = offset + int64(ebr.Ebr_next)

		}

	} else {
		return errors.New("no existe una partición extendida en el disco")
	}
	return nil
}

func createEBR2(fdisk *FDISK, offset int64) error {
	sizeBytes, errr := utilidades.ConvertToBytes(fdisk.size, fdisk.unit)
	if errr != nil {
		fmt.Println("Error converting size:", errr)
		return errr
	}

	ebr := &estructuras.EBR2{
		Ebr_mount: [1]byte{'9'},
		Ebr_fit:   [1]byte{'9'},
		Ebr_start: int32(offset) + 30,
		Ebr_size:  int32(sizeBytes),
		Ebr_next:  int32(-1),
		Ebr_name:  [16]byte{'0'},
	}

	// Serializar el MBR en el archivo
	err := ebr.Serialize2(fdisk.path, offset)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}

/*
func createExtendedPartition(fdisk *FDISK, sizeBytes int) (string, error) {
	// Crear una instancia de MBR
	var mbr estructuras.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(fdisk.path)
	if err != nil {
		fmt.Println("Error deserializando el MBR:", err)
		return "Error deserializando el MBR", err
	}

	// Obtener la primera partición disponible
	availablePartition, startPartition, indexPartition := mbr.GetFirstAvailablePartition()
	if availablePartition == nil {
		return "No hay particiones disponibles", fmt.Errorf("no hay particiones disponibles")
	}

	// Crear la partición extendida con los parámetros proporcionados
	availablePartition.CreatePartition(startPartition, sizeBytes, fdisk.typ, fdisk.fit, fdisk.name)

	// Ajustar el MBR con la partición extendida antes de serializar el EBR
	mbr.Mbr_partitions[indexPartition] = *availablePartition

	// Escribir el primer EBR al inicio de la partición extendida
	var ebr estructuras.EBR
	ebr.Part_status[0] = 0
	ebr.Part_fit[0] = fdisk.fit[0]
	ebr.Part_start = int32(startPartition) + int32(unsafe.Sizeof(ebr))
	ebr.Part_size = 0
	ebr.Part_next = -1
	copy(ebr.Part_name[:], []byte(fdisk.name))

	// Serializar el EBR en el disco
	err = ebr.Serialize(fdisk.path, (int64(startPartition)))
	if err != nil {
		return "Error serializando el primer EBR", fmt.Errorf("error serializando el primer EBR: %v", err)
	}

	// Serializar el MBR nuevamente en el disco
	err = mbr.Serialize(fdisk.path)
	if err != nil {
		return "Error serializando el MBR", fmt.Errorf("error serializando el MBR: %v", err)
	}

	// Prints de depuración
	/*
		fmt.Println("\nPartición creada (modificada) en la extendida:")
		availablePartition.Print()

		fmt.Println("Partición Extendida Creada:")
		fmt.Printf("Part_start: %d, Part_size: %d\n", availablePartition.Part_start, availablePartition.Part_size)

		fmt.Println("Primer EBR creado dentro de la partición extendida:")
		fmt.Println(ebr.Print())

		fmt.Println("IMPRIMIENDO PARTICIONES DEL MBR")
		mbr.PrintPartitions()
*/
/*
	return "Extendida", nil
}
*/

/*
func createLogicalPartition(fdisk *FDISK, sizeBytes int) (string, error) {
	// Crear una instancia de MBR
	var mbr estructuras.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(fdisk.path)
	if err != nil {
		return "", fmt.Errorf("error deserializando el MBR: %v", err)
	}

	// Verificar si existe una partición extendida
	extendedPartition, startExtended, _ := mbr.GetExtendedPartition()
	if extendedPartition == nil {
		return "", fmt.Errorf("no hay partición extendida disponible")
	} else {
		fmt.Println("PARTICION EXTENDIDA ENCONTRADA EN EL INDICE:", startExtended)
	}

	// Leer el primer EBR de la partición extendida
	var ebr estructuras.EBR
	ebrStart := int64(extendedPartition.Part_start) // offset del primer EBR
	err = ebr.Deserialize(fdisk.path, ebrStart)     // Leer el primer EBR en el offset de la partición extendida
	if err != nil {
		return "", fmt.Errorf("error leyendo el EBR: %v", err)
	}

	// Recorrer las particiones lógicas existentes
	lastEBR := &ebr
	for lastEBR.Part_next != -1 {
		ebrStart = int64(lastEBR.Part_next) // Siguiente EBR según el Part_next
		err = ebr.Deserialize(fdisk.path, ebrStart)
		if err != nil {
			return "", fmt.Errorf("error leyendo el EBR: %v", err)
		}
		lastEBR = &ebr
	}

	// Mensaje para revisar los valores antes de calcular el nuevo EBR
	fmt.Printf("Último EBR -> Start: %d, Size: %d, Next: %d\n", lastEBR.Part_start, lastEBR.Part_size, lastEBR.Part_next)

	// Calcular el inicio de la nueva partición lógica
	newEBRStart := int64(lastEBR.Part_start) + int64(lastEBR.Part_size) + int64(binary.Size(ebr))
	fmt.Printf("Calculando nueva partición lógica -> Nuevo EBR Start: %d\n", newEBRStart)
	fmt.Printf("Tamaño del EBR (usando unsafe.Sizeof): %d\n", unsafe.Sizeof(ebr))

	// Crear una nueva partición lógica
	var newEBR estructuras.EBR
	newEBR.Part_status[0] = 0              // Partición no montada
	newEBR.Part_fit[0] = fdisk.fit[0]      // Tipo de ajuste
	newEBR.Part_start = int32(newEBRStart) // El inicio de la nueva partición
	newEBR.Part_size = int32(sizeBytes)    // Tamaño de la nueva partición
	newEBR.Part_next = -1                  // No hay más particiones por ahora
	copy(newEBR.Part_name[:], fdisk.name)  // Asignar nombre a la nueva partición lógica

	// Actualizar el campo Part_next del último EBR con el nuevo inicio
	lastEBR.Part_next = int32(newEBRStart)                         // Apuntar al nuevo EBR
	err = lastEBR.Serialize(fdisk.path, int64(lastEBR.Part_start)) // Serializar el EBR anterior en su posición
	if err != nil {
		return "", fmt.Errorf("error actualizando el EBR anterior: %v", err)
	}

	// Guardar el nuevo EBR en el disco en su nueva posición
	err = newEBR.Serialize(fdisk.path, newEBRStart) // Serializar el nuevo EBR en su posición
	if err != nil {
		return "", fmt.Errorf("error serializando el nuevo EBR: %v", err)
	}



	fmt.Println("Nueva partición lógica creada:")
	fmt.Println(newEBR.Print())

	return "Lógica", nil
}
*/

/*
func createLogicalPartition(fdisk *FDISK, sizeBytes int) (string, error) {
	var mbr estructuras.MBR
	err := mbr.Deserialize(fdisk.path)
	if err != nil {
		return "", fmt.Errorf("error deserializando el MBR: %v", err)
	}

	extendedPartition, _, _ := mbr.GetExtendedPartition()
	if extendedPartition == nil {
		return "", fmt.Errorf("no hay partición extendida disponible")
	}

	var ebr estructuras.EBR
	ebrStart := int64(extendedPartition.Part_start)
	var lastEBR *estructuras.EBR
	var newEBRStart int64

	// Buscar el último EBR o un espacio libre
	for {
		err := ebr.Deserialize(fdisk.path, ebrStart)
		if err != nil {
			return "", fmt.Errorf("error leyendo el EBR: %v", err)
		}

		if ebr.Part_status[0] == 0 && ebr.Part_size == 0 {
			// Encontramos un EBR vacío, usaremos este espacio
			newEBRStart = ebrStart
			break
		}

		if ebr.Part_next == -1 {
			// Llegamos al final de la cadena de EBRs
			newEBRStart = int64(ebr.Part_start) + int64(ebr.Part_size)
			lastEBR = &ebr
			break
		}

		ebrStart = int64(ebr.Part_next)
		lastEBR = &ebr
	}

	// Verificar si hay espacio suficiente
	if newEBRStart+int64(sizeBytes) > int64(extendedPartition.Part_start)+int64(extendedPartition.Part_size) {
		return "", fmt.Errorf("no hay suficiente espacio en la partición extendida")
	}

	// Crear una nueva partición lógica
	var newEBR estructuras.EBR
	newEBR.CreateLogicalPartition(
		int(newEBRStart)+int(unsafe.Sizeof(newEBR)),
		sizeBytes,
		-1, // No hay siguiente EBR por ahora
		string(fdisk.fit[0]),
		fdisk.name,
	)

	// Actualizar el EBR anterior si existe
	if lastEBR != nil {
		lastEBR.Part_next = int32(newEBRStart)
		err = lastEBR.Serialize(fdisk.path, int64(lastEBR.Part_start))
		if err != nil {
			return "", fmt.Errorf("error actualizando el EBR anterior: %v", err)
		}
	}

	// Guardar el nuevo EBR en el disco
	err = newEBR.Serialize(fdisk.path, newEBRStart)
	if err != nil {
		return "", fmt.Errorf("error serializando el nuevo EBR: %v", err)
	}

	fmt.Println("Nueva partición lógica creada:")
	newEBR.Print()

	return "Lógica", nil
}
*/
