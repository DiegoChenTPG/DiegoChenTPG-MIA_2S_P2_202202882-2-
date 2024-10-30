package comandos

import (
	estructuras "PROYECTO2/estructuras"
	global "PROYECTO2/global"
	utilidades "PROYECTO2/utilidades"
	"errors"
	"fmt"
	"strings"
)

// MKDIR estructura que representa el comando mkdir con sus parámetros
type MKDIR struct {
	path string // Path del directorio
	p    bool   // Opción -p (crea directorios padres si no existen)
}

/*
	mkdir -p -path=/home/user/docs/usac
	mkdir -path="/home/mis documentos/archivos clases"
*/

func ParserMkdir(tokens []string) (string, error) {

	//Primera Verificacion de sesion activa, ya que este comando se debe hacer con una sesion activa
	if len(global.UserSessions) == 0 {
		return "No hay una sesión activa para ejecutar esta accion", errors.New("no hay una sesión activa")
	}

	cmd := &MKDIR{} // Crea una nueva instancia de MKDIR

	// Itera sobre cada coincidencia encontrada
	for _, match := range tokens {
		// Divide cada parte en clave y valor usando "=" como delimitador
		kv := strings.SplitN(match, "=", 2)
		key := strings.ToLower(kv[0])

		// Switch para manejar diferentes parámetros
		switch key {
		case "-path":
			if len(kv) != 2 {
				return "formato de parametro invalido", fmt.Errorf("formato de parámetro inválido: %s", match)
			}
			value := kv[1]
			// Remove quotes from value if present
			if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
				value = strings.Trim(value, "\"")
			}
			cmd.path = value
		case "-p":
			cmd.p = true
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return "parametro desconocido: " + key, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -path haya sido proporcionado
	if cmd.path == "" {
		return "faltan parámetros requeridos: -path", errors.New("faltan parámetros requeridos: -path")
	}

	// Aquí se puede agregar la lógica para ejecutar el comando mkdir con los parámetros proporcionados
	consola_retorno, err := commandMkdir(cmd)
	if err != nil {
		return consola_retorno, err
	}

	return fmt.Sprintf("Directorio %s creado correctamente.", cmd.path), nil // Devuelve el comando MKDIR creado
}

// Aquí debería de estar logeado un usuario, por lo cual el usuario debería tener consido el id de la partición
// En este caso el ID va a estar quemado
//var idPartition = "531A"

func commandMkdir(mkdir *MKDIR) (string, error) {

	// Obtener el ID del usuario logeado
	IDUsuario := global.ObtenerIDUsuarioLogueado()
	fmt.Println(IDUsuario)

	// Obtener la partición montada
	partitionSuperblock, mountedPartition, partitionPath, err := global.GetMountedPartitionSuperblock(IDUsuario)
	if err != nil {
		return "error al obtener la partición montada", fmt.Errorf("error al obtener la partición montada: %w", err)
	}

	// Crear el directorio
	consola_retorno, err := createDirectory(mkdir.path, partitionSuperblock, partitionPath, mountedPartition)
	if err != nil {
		err = fmt.Errorf("error al crear el directorio: %w", err)
		return "error al crear el directorio", err
	}

	return consola_retorno, err
}

func createDirectory(dirPath string, sb *estructuras.SuperBlock, partitionPath string, mountedPartition *estructuras.Partition) (string, error) {
	fmt.Println("\nCreando directorio:", dirPath)

	parentDirs, destDir := utilidades.GetParentDirectories(dirPath)
	fmt.Println("\nDirectorios padres:", parentDirs)
	fmt.Println("Directorio destino:", destDir)

	// Crear el directorio segun el path proporcionado
	err := sb.CreateFolder(partitionPath, parentDirs, destDir)
	if err != nil {
		return "error al crear el directorio", fmt.Errorf("error al crear el directorio: %w", err)
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
