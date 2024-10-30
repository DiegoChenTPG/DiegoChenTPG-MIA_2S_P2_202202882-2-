package comandos

import (
	// Importa el paquete "estructuras" desde el directorio "EDD2021/estructuras"
	global "PROYECTO2/global"
	reportes "PROYECTO2/reportes"
	"errors"
	"fmt"
	"strings"
)

// MKDISK estructura que representa el comando mkdisk con sus parámetros
type REP struct {
	id           string
	name         string
	path_file_ls string
	path         string // Ruta del archivo del disco
}

// CommandRep parsea el comando rep y devuelve una instancia de REP
func ParserRep(tokens []string) (*REP, string, error) {
	cmd := &REP{} // Crea una nueva instancia de REP
	var retorno string

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
		case "-id":
			// Verifica que el id no esté vacío
			if value == "" {
				return nil, "el id no puede estar vacio", errors.New("el id no puede estar vacío")
			}
			cmd.id = value
		case "-path":
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, "el path no puede estar vacio", errors.New("el path no puede estar vacío")
			}
			cmd.path = value
		case "-name":
			// Verifica que el nombre sea uno de los valores permitidos
			validNames := []string{"mbr", "disk", "inode", "block", "bm_inode", "bm_block", "sb", "file", "ls"}
			if !contains(validNames, value) {
				return nil, "nombre inválido, debe ser uno de los siguientes: mbr, disk, inode, block, bm_inode, bm_block, sb, file, ls", errors.New("nombre inválido, debe ser uno de los siguientes: mbr, disk, inode, block, bm_inode, bm_block, sb, file, ls")
			}
			cmd.name = value
		case "-path_file_ls":
			cmd.path_file_ls = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, "parametro desconocido: " + key, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros obligatorios hayan sido proporcionados
	if cmd.id == "" || cmd.path == "" || cmd.name == "" {
		return nil, "faltan parametros requeridos: -id, -path, -name", errors.New("faltan parámetros requeridos: -id, -path, -name")
	}

	// Crear el disco con los parámetros proporcionados
	retorno, err := commandRep(cmd)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, "Error al crear el archivo", nil
	}

	return cmd, retorno, nil // Devuelve el comando FDISK creado
}

func contains(list []string, value string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func commandRep(rep *REP) (string, error) {
	// Obtener la partición montada
	var mensaje_retorno string
	mountedMbr, mountedSb, mountedDiskPath, err := global.GetMountedPartitionRep(rep.id)
	if err != nil {
		return "", err
	}

	// Switch para manejar diferentes tipos de reportes
	switch rep.name {
	case "mbr":
		err = reportes.ReportMBR(mountedMbr, rep.path, mountedDiskPath)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		mensaje_retorno = "Reporte MBR creado correctamente"

	case "inode":
		err = reportes.ReportInode(mountedSb, mountedDiskPath, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		mensaje_retorno = "Reporte Inode creado correctamente"
	case "bm_inode":
		err = reportes.ReportBMInode(mountedSb, mountedDiskPath, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		mensaje_retorno = "Reporte bm_inode creado correctamente"
	case "bm_block":
		err = reportes.ReportBMBlock(mountedSb, mountedDiskPath, rep.path)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
		mensaje_retorno = "Reporte bm_inode creado correctamente"

	}

	return mensaje_retorno, nil
}
