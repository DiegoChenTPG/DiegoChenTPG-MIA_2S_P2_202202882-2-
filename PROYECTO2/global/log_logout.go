package global

import (
	estructuras "PROYECTO2/estructuras"
	"bytes"
	"fmt"
	"strings"
)

type SesionActiva struct {
	User     string
	Pass     string
	ID       string
	Loggeado bool
	Root     bool
}

var UserSessions []SesionActiva

// Función para verificar si el usuario ya está logueado
func IsUserLogged(username string) bool {
	for _, session := range UserSessions {
		if session.User == username && session.Loggeado {
			return true
		}
	}
	return false
}

// Función para registrar la sesión del usuario
func RegisterSession(username, password, id string, isRoot bool) {
	for i, session := range UserSessions {
		if session.User == username {
			UserSessions[i].Loggeado = true
			UserSessions[i].ID = id
			UserSessions[i].Pass = password
			UserSessions[i].Root = (username == "root")
			return
		}
	}

	// Si el usuario no existe, agregar uno nuevo
	UserSessions = append(UserSessions, SesionActiva{
		User:     username,
		Pass:     password,
		ID:       id,
		Loggeado: true,
		Root:     (username == "root"),
	})
}

func ObtenerIDRoot() string {
	for _, session := range UserSessions {
		if session.Root {
			return session.ID
		}
	}
	return ""
}

// Función para verificar si el usuario actual es root
func VerifacionRoot() bool {
	for _, session := range UserSessions {
		if session.Root && session.Loggeado {
			return true
		}
	}
	return false
}

func ObtenerIDUsuarioLogueado() string {
	for _, session := range UserSessions {
		if session.Loggeado {
			return session.ID
		}
	}
	return ""
}

func Verificar_login(n_usuario string, contrasena string, id_particion string) (string, bool) {
	//PRIMERO VERIFICAMOS SI EL USUARIO ESTA LOGEADO
	if IsUserLogged(n_usuario) {
		return fmt.Sprintf("Error: La sesión ya está iniciada para el usuario %s", n_usuario), false //, errors.New("sesión ya iniciada")
	}

	// Obtener la partición montada usando el ID proporcionado
	mountedPartition, partitionPath, err := GetMountedPartition(id_particion)
	if err != nil {
		return "error al obtener la particion montada", false //, fmt.Errorf("error al obtener la partición montada: %v", err)
	}

	// Leer el superbloque
	superBlock := &estructuras.SuperBlock{}
	err = superBlock.Deserialize(partitionPath, int64(mountedPartition.Part_start))
	if err != nil {
		return "", false //, fmt.Errorf("error al leer el superbloque: %v", err)
	}

	// Buscar el inodo de users.txt
	rootInode := &estructuras.Inode{}
	err = rootInode.Deserialize(partitionPath, int64(superBlock.S_inode_start))
	if err != nil {
		return "", false //, fmt.Errorf("error al leer el inodo raíz: %v", err)
	}

	var usersInode *estructuras.Inode
	//var usersInodeIndex int32
	for _, blockIndex := range rootInode.I_block {
		if blockIndex == -1 {
			break
		}
		folderBlock := &estructuras.FolderBlock{}
		err = folderBlock.Deserialize(partitionPath, int64(superBlock.S_block_start+blockIndex*superBlock.S_block_size))
		if err != nil {
			return "", false //, fmt.Errorf("error al leer el bloque de carpeta: %v", err)
		}
		for _, content := range folderBlock.B_content {
			if string(bytes.Trim(content.B_name[:], "\x00")) == "users.txt" {
				usersInode = &estructuras.Inode{}
				err = usersInode.Deserialize(partitionPath, int64(superBlock.S_inode_start+content.B_inodo*superBlock.S_inode_size))
				if err != nil {
					return "", false //, fmt.Errorf("error al leer el inodo de users.txt: %v", err)
				}
				//usersInodeIndex = content.B_inodo
				break
			}
		}
		if usersInode != nil {
			break
		}
	}

	if usersInode == nil {
		return "", false //, fmt.Errorf("no se encontró el archivo users.txt")
	}

	// Leer el contenido de users.txt
	fileBlock := &estructuras.FileBlock{}
	err = fileBlock.Deserialize(partitionPath, int64(superBlock.S_block_start+usersInode.I_block[0]*superBlock.S_block_size))
	if err != nil {
		return "", false //, fmt.Errorf("error al leer el bloque de archivo: %v", err)
	}

	content := string(bytes.Trim(fileBlock.B_content[:], "\x00"))
	lines := strings.Split(content, "\n")

	// Verificar si el usuario y la contraseña son correctos
	for _, line := range lines {
		parts := strings.Split(line, ",")
		if len(parts) == 5 && parts[1] == "U" {
			if strings.TrimSpace(parts[3]) == n_usuario {
				if strings.TrimSpace(parts[4]) == contrasena {
					RegisterSession(n_usuario, contrasena, id_particion, n_usuario == "root")
					return fmt.Sprintf("Inicio de sesión exitoso para el usuario %s", n_usuario), true //, nil
				}
				return "contraseña incorrecta para el usuario: " + n_usuario, false //, fmt.Errorf("contraseña incorrecta para el usuario %s", n_usuario)
			}
		}
	}

	return "el usuario no existe: " + n_usuario, false //, fmt.Errorf("el usuario %s no existe", n_usuario)
}
