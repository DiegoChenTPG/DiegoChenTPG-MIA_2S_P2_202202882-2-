package comandos

import (
	estructuras "PROYECTO2/estructuras"
	global "PROYECTO2/global"
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type CAT struct {
	file string
}

func ParserCat(tokens []string) (*CAT, string, error) {

	//Primera Verificacion de sesion activa
	if len(global.UserSessions) == 0 {
		return nil, "No hay una sesión activa para ejecutar esta accion", errors.New("no hay una sesión activa")
	}

	cmd := &CAT{} // Crea una nueva instancia de MOUNT

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
		fmt.Println(key)
		switch key {
		case "-file1":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return nil, "el nombre no puede estar vacio", errors.New("el nombre no puede estar vacío")
			}
			cmd.file = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, "parametro desconocido: " + key, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -name haya sido proporcionados
	if cmd.file == "" {
		return nil, "faltan parametros requeridos: -file", errors.New("faltan parámetros requeridos: -name")
	}

	posible_error, err := commandCat(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, posible_error, err
	}

	//return cmd, "Se monto la particion: " + cmd.name + " del disco ubicado en: " + cmd.path + " en el sistema", nil // Devuelve el comando MOUNT creado
	return cmd, posible_error, nil
}

func commandCat(cat *CAT) (string, error) {
	fmt.Println(cat.file)
	// Verificar si el usuario actual es root
	if !global.VerifacionRoot() {
		return "", fmt.Errorf("solo el usuario root puede ejecutar este comando")
	}

	// Obtener el ID del usuario root
	IDUsuarioRoot := global.ObtenerIDRoot()

	// Obtener la partición montada (asumiendo que tienes una función para esto)
	mountedPartition, partitionPath, err := global.GetMountedPartition(IDUsuarioRoot)
	if err != nil {
		return "error al obtener la particion montada", fmt.Errorf("error al obtener la partición montada: %v", err)
	}

	// Leer el superbloque
	superBlock := &estructuras.SuperBlock{}
	err = superBlock.Deserialize(partitionPath, int64(mountedPartition.Part_start))
	if err != nil {
		return "", fmt.Errorf("error al leer el superbloque: %v", err)
	}

	// Buscar el inodo de users.txt
	rootInode := &estructuras.Inode{}
	err = rootInode.Deserialize(partitionPath, int64(superBlock.S_inode_start))
	if err != nil {
		return "", fmt.Errorf("error al leer el inodo raíz: %v", err)
	}

	var usersInode *estructuras.Inode
	for _, blockIndex := range rootInode.I_block {
		if blockIndex == -1 {
			break
		}
		folderBlock := &estructuras.FolderBlock{}
		err = folderBlock.Deserialize(partitionPath, int64(superBlock.S_block_start+blockIndex*superBlock.S_block_size))
		if err != nil {
			return "", fmt.Errorf("error al leer el bloque de carpeta: %v", err)
		}
		for _, content := range folderBlock.B_content {
			if string(bytes.Trim(content.B_name[:], "\x00")) == "users.txt" {
				usersInode = &estructuras.Inode{}
				err = usersInode.Deserialize(partitionPath, int64(superBlock.S_inode_start+content.B_inodo*superBlock.S_inode_size))
				if err != nil {
					return "", fmt.Errorf("error al leer el inodo de users.txt: %v", err)
				}
				break
			}
		}
		if usersInode != nil {
			break
		}
	}

	if usersInode == nil {
		return "", fmt.Errorf("no se encontró el archivo users.txt")
	}

	// Leer el contenido de users.txt
	fileBlock := &estructuras.FileBlock{}
	err = fileBlock.Deserialize(partitionPath, int64(superBlock.S_block_start+usersInode.I_block[0]*superBlock.S_block_size))
	if err != nil {
		return "", fmt.Errorf("error al leer el bloque de archivo: %v", err)
	}

	// Retornar el contenido del archivo users.txt como una cadena
	content := string(bytes.Trim(fileBlock.B_content[:], "\x00"))
	return content, nil
}
