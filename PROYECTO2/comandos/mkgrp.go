package comandos

import (
	estructuras "PROYECTO2/estructuras"
	global "PROYECTO2/global"
	"bytes"
	"errors" // Paquete para manejar errores y crear nuevos errores con mensajes personalizados
	"fmt"    // Paquete para formatear cadenas y realizar operaciones de entrada/salida
	"strconv"

	"strings" // Paquete para manipular cadenas, como unir, dividir, y modificar contenido de cadenas
)

// MOUNT estructura que representa el comando mount con sus parámetros
type MKGRP struct {
	name string // Nombre de la partición
}

/*
	mkgrp -name=usuarios
	mkgrp -name=usuarios2
*/

func ParserMkgrp(tokens []string) (*MKGRP, string, error) {
	cmd := &MKGRP{} // Crea una nueva instancia de MOUNT

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

	// Verifica que el parámetro -name haya sido proporcionados
	if cmd.name == "" {
		return nil, "faltan parametros requeridos: -name", errors.New("faltan parámetros requeridos: -name")
	}

	posible_error, err := commandMkgrp(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, posible_error, err
	}

	//return cmd, "Se monto la particion: " + cmd.name + " del disco ubicado en: " + cmd.path + " en el sistema", nil // Devuelve el comando MOUNT creado
	return cmd, posible_error, nil
}

func commandMkgrp(mkgrp *MKGRP) (string, error) {
	// Obtener la partición montada (asumiendo que tienes una función para esto)
	mountedPartition, partitionPath, err := global.GetMountedPartition("821A") // Reemplaza "vda1" con el ID correcto
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
	var usersInodeIndex int32
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
				usersInodeIndex = content.B_inodo
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

	content := string(bytes.Trim(fileBlock.B_content[:], "\x00"))
	lines := strings.Split(content, "\n")

	// Verificar si el grupo ya existe y encontrar el máximo ID de grupo
	maxGroupID := 0
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 3 && parts[1] == "G" {
			if parts[2] == mkgrp.name {
				return fmt.Sprintf("El grupo %s ya existe", mkgrp.name), nil
			}
			id, _ := strconv.Atoi(parts[0])
			if id > maxGroupID {
				maxGroupID = id
			}
		}
	}

	// Agregar el nuevo grupo, al ID obtenido le sumamos 1
	newGroupID := maxGroupID + 1
	newLine := fmt.Sprintf("%d,G,%s", newGroupID, mkgrp.name)
	lines = append(lines, newLine)

	// Actualizar el contenido del archivo
	newContent := strings.Join(lines, "\n")
	fmt.Println("IMPRIMIENDO LEN NEW CONTENT")
	fmt.Println(len(newContent))
	fmt.Println("===============================")
	fmt.Println("IMPRIMIENDO fileBlock.B_content")
	fmt.Println(len(fileBlock.B_content))
	fmt.Println("===============================")
	if len(newContent) > len(fileBlock.B_content) {
		return "el nuevo contenido excede el tamaño del bloque", fmt.Errorf("el nuevo contenido excede el tamaño del bloque")
	}

	copy(fileBlock.B_content[:], newContent)

	// Escribir el bloque actualizado
	err = fileBlock.Serialize(partitionPath, int64(superBlock.S_block_start+usersInode.I_block[0]*superBlock.S_block_size))
	if err != nil {
		return "error al escribir el bloque de archivo", fmt.Errorf("error al escribir el bloque de archivo: %v", err)
	}

	// Actualizar el tamaño del archivo en el inodo
	usersInode.I_size = int32(len(newContent))
	err = usersInode.Serialize(partitionPath, int64(superBlock.S_inode_start+usersInodeIndex*superBlock.S_inode_size))
	if err != nil {
		return "error al actualizar el inodo de users.txt", fmt.Errorf("error al actualizar el inodo de users.txt: %v", err)
	}

	verificacion_usuarios, err := superBlock.VerifyUsersFile(usersInode, partitionPath)
	if err != nil {
		return "error al verificar el archivo users.txt", fmt.Errorf("error al verificar el archivo users.txt: %v", err)
	}

	usuarios_act := fmt.Sprintf("Grupo %s creado exitosamente\n", mkgrp.name)
	usuarios_act += verificacion_usuarios

	fmt.Println(usuarios_act)

	return fmt.Sprintf("Grupo %s creado exitosamente", mkgrp.name), nil

}
