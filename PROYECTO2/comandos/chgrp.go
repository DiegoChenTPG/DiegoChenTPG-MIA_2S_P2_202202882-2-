package comandos

import (
	estructuras "PROYECTO2/estructuras"
	global "PROYECTO2/global"
	"bytes"
	"errors" // Paquete para manejar errores y crear nuevos errores con mensajes personalizados
	"fmt"    // Paquete para formatear cadenas y realizar operaciones de entrada/salida

	"strings" // Paquete para manipular cadenas, como unir, dividir, y modificar contenido de cadenas
)

type CHGRP struct {
	user [10]byte
	grp  [10]byte
}

/*
	mkgrp -name=usuarios
	mkgrp -name=usuarios2
*/

func ParserChgrp(tokens []string) (*CHGRP, string, error) {

	//Primera Verificacion de sesion activa
	if len(global.UserSessions) == 0 {
		return nil, "No hay una sesión activa para ejecutar esta accion", errors.New("no hay una sesión activa")
	}

	cmd := &CHGRP{} // Crea una nueva instancia de MOUNT

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
		case "-user":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return nil, "el nombre del user no puede estar vacio", errors.New("el nombre no puede estar vacío")
			}
			copy(cmd.user[:], []byte(value))
		case "-grp":
			// Verifica que el nombre no esté vacío
			if value == "" {
				return nil, "El grupo a asignar no debe estar vacio", errors.New("el nombre no puede estar vacío")
			}
			copy(cmd.grp[:], []byte(value))
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, "parametro desconocido: " + key, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que el parámetro -name haya sido proporcionados
	if len(bytes.TrimSpace(cmd.user[:])) == 0 {
		return nil, "faltan parametros requeridos: -user", errors.New("faltan parámetros requeridos: -name")
	}

	if len(bytes.TrimSpace(cmd.grp[:])) == 0 {
		return nil, "faltan parametros requeridos: -grp", errors.New("faltan parámetros requeridos: -name")
	}

	posible_error, err := commandChgrp(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, posible_error, err
	}

	//return cmd, "Se monto la particion: " + cmd.name + " del disco ubicado en: " + cmd.path + " en el sistema", nil // Devuelve el comando MOUNT creado
	return cmd, posible_error, nil
}

func commandChgrp(chgrp *CHGRP) (string, error) {

	// Verificar si el usuario actual es root
	if !global.VerifacionRoot() {
		return "", fmt.Errorf("solo el usuario root puede ejecutar este comando")
	}

	// Obtener el ID del usuario root
	IDUsuarioRoot := global.ObtenerIDRoot()

	// Obtener la partición montada
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

	// Verificar si el usuario existe y si el nuevo grupo existe
	userExists := false
	groupExists := false
	var updatedLines []string

	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 5 && parts[1] == "U" {
			if string(bytes.Trim(chgrp.user[:], "\x00")) == strings.TrimSpace(parts[3]) {
				userExists = true
				// Cambiar el grupo del usuario
				parts[2] = strings.TrimSpace(string(bytes.Trim(chgrp.grp[:], "\x00")))
				line = strings.Join(parts, ",")
			}
		} else if len(parts) == 3 && parts[1] == "G" {
			if string(bytes.Trim(chgrp.grp[:], "\x00")) == strings.TrimSpace(parts[2]) {
				groupExists = true
			}
		}
		updatedLines = append(updatedLines, line)
	}

	if !userExists {
		return fmt.Sprintf("El usuario %s no existe", string(bytes.Trim(chgrp.user[:], "\x00"))), nil
	}

	if !groupExists {
		return fmt.Sprintf("El grupo %s no existe", string(bytes.Trim(chgrp.grp[:], "\x00"))), nil
	}

	// Actualizar el contenido del archivo
	newContent := strings.Join(updatedLines, "\n")

	if len(newContent) > len(fileBlock.B_content) {
		return "", fmt.Errorf("el nuevo contenido excede el tamaño del bloque")
	}

	copy(fileBlock.B_content[:], newContent)

	// Escribir el bloque actualizado
	err = fileBlock.Serialize(partitionPath, int64(superBlock.S_block_start+usersInode.I_block[0]*superBlock.S_block_size))
	if err != nil {
		return "", fmt.Errorf("error al escribir el bloque de archivo: %v", err)
	}

	// Actualizar el tamaño del archivo en el inodo si es necesario
	if usersInode.I_size != int32(len(newContent)) {
		usersInode.I_size = int32(len(newContent))
		err = usersInode.Serialize(partitionPath, int64(superBlock.S_inode_start+usersInodeIndex*superBlock.S_inode_size))
		if err != nil {
			return "", fmt.Errorf("error al actualizar el inodo de users.txt: %v", err)
		}
	}

	verificacion_usuarios, err := superBlock.VerifyUsersFile(usersInode, partitionPath)
	if err != nil {
		return "", fmt.Errorf("error al verificar el archivo users.txt: %v", err)
	}

	usuarios_act := fmt.Sprintf("Grupo del usuario %s cambiado exitosamente a %s\n",
		string(bytes.Trim(chgrp.user[:], "\x00")),
		string(bytes.Trim(chgrp.grp[:], "\x00")))
	usuarios_act += verificacion_usuarios

	fmt.Println(usuarios_act)

	return fmt.Sprintf("Grupo del usuario %s cambiado exitosamente a %s",
		string(bytes.Trim(chgrp.user[:], "\x00")),
		string(bytes.Trim(chgrp.grp[:], "\x00"))), nil
}
