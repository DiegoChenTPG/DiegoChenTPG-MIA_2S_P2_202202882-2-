package comandos

import (
	// Importa el paquete "structure" desde el directorio "EDD2021/structure"
	estructuras "PROYECTO2/estructuras" // Importa el paquete "estructuras" desde el directorio "EDD2021/estructuras"
	utilidades "PROYECTO2/utilidades"   // Importa el paquete "utilidades" desde el directorio "EDD2021/utilidades"
	"errors"                            // Paquete para manejar errores y crear nuevos errores con mensajes personalizados
	"fmt"                               // Paquete para formatear cadenas y realizar operaciones de entrada/salida
	"math/rand"
	"os"
	"path/filepath" // Paquete para trabajar con expresiones regulares, útil para encontrar y manipular patrones en cadenas
	"strconv"       // Paquete para convertir cadenas a otros tipos de datos, como enteros
	"strings"       // Paquete para manipular cadenas, como unir, dividir, y modificar contenido de cadenas
	"time"
)

// MKDISK estructura que representa el comando mkdisk con sus parámetros
type MKDISK struct {
	size int    // Tamaño del disco
	unit string // Unidad de medida del tamaño (K o M)
	fit  string // Tipo de ajuste (BF, FF, WF)
	path string // Ruta del archivo del disco
}

/*
	mkdisk -size=3000 -unit=K -path=/home/user/Disco1.mia
	mkdisk -size=3000 -path=/home/user/Disco1.mia
	mkdisk -size=5 -unit=M -fit=WF -path="/home/keviin/University/PRACTICAS/MIA_LAB_S2_2024/CLASE03/disks/Disco1.mia"
	mkdisk -size=10 -path="/home/mis discos/Disco4.mia"
*/

func ParserMkdisk(tokens []string) (*MKDISK, string, error) {
	cmd := &MKDISK{} // Crea una nueva instancia de MKDISK

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
			if value != "K" && value != "M" {
				return nil, "la unidad debe ser K o M", errors.New("la unidad debe ser K o M")
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
				return nil, "el path no puede estar vacio", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, "parametro desconocido: " + key, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -size y -path hayan sido proporcionados
	if cmd.size == 0 {
		return nil, "faltan parametros requeridos: -size", errors.New("faltan parámetros requeridos: -size")
	}
	if cmd.path == "" {
		return nil, "faltan parámetros requeridos: -path", errors.New("faltan parámetros requeridos: -path")
	}

	// Si no se proporcionó la unidad, se establece por defecto a "M"
	if cmd.unit == "" {
		cmd.unit = "M"
	}

	// Si no se proporcionó el ajuste, se establece por defecto a "FF"
	if cmd.fit == "" {
		cmd.fit = "FF"
	}

	// Crear el disco con los parámetros proporcionados
	err := commandMkdisk(cmd)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return cmd, "SE CREO EL DISCO CORRECTAMENTE EN LA RUTA: " + cmd.path, nil // Devuelve el comando MKDISK creado
}

func commandMkdisk(mkdisk *MKDISK) error {
	// Convertir el tamaño a bytes
	sizeBytes, err := utilidades.ConvertToBytes(mkdisk.size, mkdisk.unit)
	if err != nil {
		fmt.Println("Error converting size:", err)
		return err
	}

	// Crear el disco con el tamaño proporcionado
	err = createDisk(mkdisk, sizeBytes)
	if err != nil {
		fmt.Println("Error creating disk:", err)
		return err
	}

	// Crear el MBR con el tamaño proporcionado
	err = createMBR(mkdisk, sizeBytes)
	if err != nil {
		fmt.Println("Error creating MBR:", err)
		return err
	}

	return nil
}

func createDisk(mkdisk *MKDISK, sizeBytes int) error {
	// Crear las carpetas necesarias
	err := os.MkdirAll(filepath.Dir(mkdisk.path), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return err
	}

	// Crear el archivo binario
	file, err := os.Create(mkdisk.path)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}
	defer file.Close()

	// Escribir en el archivo usando un buffer de 1 MB
	buffer := make([]byte, 1024*1024) // Crea un buffer de 1 MB
	for sizeBytes > 0 {
		writeSize := len(buffer)
		if sizeBytes < writeSize {
			writeSize = sizeBytes // Ajusta el tamaño de escritura si es menor que el buffer
		}
		if _, err := file.Write(buffer[:writeSize]); err != nil {
			return err // Devuelve un error si la escritura falla
		}
		sizeBytes -= writeSize // Resta el tamaño escrito del tamaño total
	}
	utilidades.GuardarDiscosCreados(mkdisk.path)
	discos := utilidades.ObtenerNombresDiscos()
	fmt.Println("IMPRIENDO DISCOS")
	fmt.Println(discos)
	return nil
}

func createMBR(mkdisk *MKDISK, sizeBytes int) error {

	// Crear el MBR con los valores proporcionados
	mbr := &estructuras.MBR{
		Mbr_size:           int32(sizeBytes),
		Mbr_creation_date:  float32(time.Now().Unix()),
		Mbr_disk_signature: rand.Int31(),
		Mbr_disk_fit:       [1]byte{mkdisk.fit[0]},
		Mbr_partitions: [4]estructuras.Partition{
			{Part_status: [1]byte{'9'}, Part_type: [1]byte{'0'}, Part_fit: [1]byte{'0'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'0'}, Part_correlative: -1, Part_id: [4]byte{'0'}},
			{Part_status: [1]byte{'9'}, Part_type: [1]byte{'0'}, Part_fit: [1]byte{'0'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'0'}, Part_correlative: -1, Part_id: [4]byte{'0'}},
			{Part_status: [1]byte{'9'}, Part_type: [1]byte{'0'}, Part_fit: [1]byte{'0'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'0'}, Part_correlative: -1, Part_id: [4]byte{'0'}},
			{Part_status: [1]byte{'9'}, Part_type: [1]byte{'0'}, Part_fit: [1]byte{'0'}, Part_start: -1, Part_size: -1, Part_name: [16]byte{'0'}, Part_correlative: -1, Part_id: [4]byte{'0'}},
		},
	}

	// Serializar el MBR en el archivo
	err := mbr.Serialize(mkdisk.path)
	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}
