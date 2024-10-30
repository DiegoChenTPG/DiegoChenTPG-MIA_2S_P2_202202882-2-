package comandos

import (
	estructuras "PROYECTO2/estructuras"
	global "PROYECTO2/global"
	utilidades "PROYECTO2/utilidades"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// MKFILE estructura que representa el comando mkfile con sus parámetros
type MKFILE struct {
	path string // Ruta del archivo
	r    bool   // Opción recursiva
	size int    // Tamaño del archivo
	cont string // Contenido del archivo
}

// ParserMkfile parsea el comando mkfile y devuelve una instancia de MKFILE
func ParserMkfile(tokens []string) (string, error) {

	//Primera Verificacion de sesion activa, ya que este comando se debe hacer con una sesion activa
	if len(global.UserSessions) == 0 {
		return "No hay una sesión activa para ejecutar esta accion", errors.New("no hay una sesión activa")
	}

	cmd := &MKFILE{} // Crea una nueva instancia de MKFILE

	// Itera sobre cada coincidencia encontrada
	for _, match := range tokens {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		key := strings.ToLower(kv[0])
		var value string
		if len(kv) == 2 {
			value = kv[1]
		}

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
			cmd.path = value
		case "-r":
			// Establece el valor de r a true
			cmd.r = true
		case "-size":
			// Convierte el valor del tamaño a un entero
			size, err := strconv.Atoi(value)
			if err != nil || size < 0 {
				return "el tamaño debe ser un número entero no negativo", errors.New("el tamaño debe ser un número entero no negativo")
			}
			cmd.size = size
		case "-cont":
			// Verifica que el contenido no esté vacío
			if value == "" {
				return "el contenido no puede estar vacío", errors.New("el contenido no puede estar vacío")
			}
			cmd.cont = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "parámetro desconocido: " + key, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -path haya sido proporcionado
	if cmd.path == "" {
		return "faltan parámetros requeridos: -path", errors.New("faltan parámetros requeridos: -path")
	}

	// Si no se proporcionó el tamaño, se establece por defecto a 0
	if cmd.size == 0 {
		cmd.size = 0
	}

	// Si no se proporcionó el contenido, se establece por defecto a ""
	if cmd.cont == "" {
		cmd.cont = ""
	}

	// Crear el archivo con los parámetros proporcionados
	retorno_consola, err := commandMkfile(cmd)
	if err != nil {
		return retorno_consola, err
	}

	return fmt.Sprintf("MKFILE: Archivo %s creado correctamente.", cmd.path), nil // Devuelve el comando MKFILE creado
}

// Función ficticia para crear el archivo (debe ser implementada)
func commandMkfile(mkfile *MKFILE) (string, error) {

	// Obtener el ID del usuario logeado
	IDUsuario := global.ObtenerIDUsuarioLogueado()
	fmt.Println(IDUsuario)

	// Obtener la partición montada
	partitionSuperblock, mountedPartition, partitionPath, err := global.GetMountedPartitionSuperblock(IDUsuario)
	if err != nil {
		return "error al obtener la particion montada", fmt.Errorf("error al obtener la partición montada: %w", err)
	}

	// Generar el contenido del archivo si no se proporcionó
	if mkfile.cont == "" {
		mkfile.cont = generateContent(mkfile.size)
	}

	// Crear el archivo
	consola_retorno, err := createFile(mkfile.path, mkfile.size, mkfile.cont, partitionSuperblock, partitionPath, mountedPartition)
	if err != nil {
		err = fmt.Errorf("error al crear el archivo: %w", err)
		return "error al crear el archivo", err
	}

	return consola_retorno, err
}

// generateContent genera una cadena de números del 0 al 9 hasta cumplir el tamaño ingresado
func generateContent(size int) string {
	content := ""
	for len(content) < size {
		content += "0123456789"
	}
	return content[:size] // Recorta la cadena al tamaño exacto
}

// Funcion para crear un archivo
func createFile(filePath string, size int, content string, sb *estructuras.SuperBlock, partitionPath string, mountedPartition *estructuras.Partition) (string, error) {
	fmt.Println("\nCreando archivo:", filePath)

	parentDirs, destDir := utilidades.GetParentDirectories(filePath)
	fmt.Println("\nDirectorios padres:", parentDirs)
	fmt.Println("Directorio destino:", destDir)

	// Obtener contenido por chunks
	chunks := utilidades.SplitStringIntoChunks(content)
	fmt.Println("\nChunks del contenido:", chunks)

	// Crear el archivo
	err := sb.CreateFile(partitionPath, parentDirs, destDir, size, chunks)
	if err != nil {
		return "error al crear el archivo", fmt.Errorf("error al crear el archivo: %w", err)
	}

	// Imprimir inodos y bloques
	sb.PrintInodes(partitionPath)
	sb.PrintBlocks(partitionPath)

	// Serializar el superbloque
	err = sb.Serialize(partitionPath, int64(mountedPartition.Part_start))
	if err != nil {
		return "error al serializar el superbloque", fmt.Errorf("error al serializar el superbloque: %w", err)
	}

	return "", nil
}
