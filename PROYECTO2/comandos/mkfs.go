package comandos

import (
	estructuras "PROYECTO2/estructuras"
	global "PROYECTO2/global"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

type MKFS struct {
	id  string
	typ string
	fs  string
}

func ParserMkfs(tokens []string) (*MKFS, string, error) {
	cmd := &MKFS{} // Crea una nueva instancia de MOUNT

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
			// Verifica que el path no esté vacío
			if value == "" {
				return nil, "el id no puede estar vacio", errors.New("el id no puede estar vacío")
			}
			cmd.id = value
		case "-type":
			// Verifica que el nombre no esté vacío
			value = strings.ToLower(value)
			if value != "full" {

				return nil, "el tipo debe ser full", errors.New("el tipo debe ser full")

			}
			cmd.typ = strings.ToLower(value)
		case "-fs":
			fmt.Println("implementar -fs")
			// Verifica que el sistema de archivos sea "2fs" o "3fs"
			if value != "2fs" && value != "3fs" {
				return nil, "", errors.New("el sistema de archivos debe ser 2fs o 3fs")
			}
			cmd.fs = value
		default:
			// Si el parámetro no es reconocido, devuelve un error
			return nil, "parametro desconocido: " + key, fmt.Errorf("parámetro desconocido: %s", key)
		}
	}

	// Verifica que los parámetros -path y -name hayan sido proporcionados
	if cmd.id == "" {
		return nil, "faltan parametros requeridos: -id", errors.New("faltan parámetros requeridos: -id")
	}
	if cmd.typ == "" {
		cmd.typ = "full"
	}

	if cmd.fs == "" {
		cmd.fs = "2fs"
	}

	err := commandMkfs(cmd)
	if err != nil {

		fmt.Println("Error: ", err)

	}

	return cmd, "Se realizo un formateo completo de la particion con el ID: " + cmd.id + " y con el sistema de archivos " + cmd.fs, nil

}

func commandMkfs(mkfs *MKFS) error {
	// Obtener la partición montada
	mountedPartition, partitionPath, err := global.GetMountedPartition(mkfs.id)
	if err != nil {
		return err
	}

	// Calcular el valor de n
	n := calculateN(mountedPartition, mkfs.fs)

	fmt.Printf("Valor de N: %d\n", n)

	// Inicializar un nuevo superbloque
	superBlock := createSuperBlock(mountedPartition, n, mkfs.fs)

	// Crear los bitmaps
	err = superBlock.CreateBitMaps(partitionPath)
	if err != nil {
		return err
	}
	/*
		// Crear archivo users.txt
		err = superBlock.CreateUsersFile(partitionPath)
		if err != nil {
			return err
		}
	*/

	// Validar que sistema de archivos es
	if superBlock.S_filesystem_type == 3 {
		// Crear archivo users.txt ext3
		err = superBlock.CreateUsersFileExt3(partitionPath, int64(mountedPartition.Part_start+int32(binary.Size(estructuras.SuperBlock{}))))
		if err != nil {
			return err
		}
	} else {
		// Crear archivo users.txt ext2
		err = superBlock.CreateUsersFileExt2(partitionPath)
		if err != nil {
			return err
		}
	}

	// Serializar el superbloque
	err = superBlock.Serialize(partitionPath, int64(mountedPartition.Part_start))
	if err != nil {
		return err
	}

	return nil
}

func calculateN(partition *estructuras.Partition, fs string) int32 {
	// Numerador: tamaño de la partición menos el tamaño del superblock
	numerator := int(partition.Part_size) - binary.Size(estructuras.SuperBlock{})

	// Denominador base: 4 + tamaño de inodos + 3 * tamaño de bloques de archivo
	baseDenominator := 4 + binary.Size(estructuras.Inode{}) + 3*binary.Size(estructuras.FileBlock{})

	// Si el sistema de archivos es "3fs", se añade el tamaño del journaling al denominador
	temp := 0
	if fs == "3fs" {
		temp = binary.Size(estructuras.Journal{})
	}

	// Denominador final
	denominator := baseDenominator + temp

	// Calcular n
	n := math.Floor(float64(numerator) / float64(denominator))

	return int32(n)
}

func createSuperBlock(partition *estructuras.Partition, n int32, fs string) *estructuras.SuperBlock {
	// Calcular punteros de las estructuras
	journal_start, bm_inode_start, bm_block_start, inode_start, block_start := calculateStartPositions(partition, fs, n)

	fmt.Printf("Journal Start: %d\n", journal_start)
	fmt.Printf("Bitmap Inode Start: %d\n", bm_inode_start)
	fmt.Printf("Bitmap Block Start: %d\n", bm_block_start)
	fmt.Printf("Inode Start: %d\n", inode_start)
	fmt.Printf("Block Start: %d\n", block_start)

	// Tipo de sistema de archivos
	var fsType int32

	if fs == "2fs" {
		fsType = 2
	} else {
		fsType = 3
	}

	// Crear un nuevo superbloque
	superBlock := &estructuras.SuperBlock{
		S_filesystem_type:   fsType,
		S_inodes_count:      0,
		S_blocks_count:      0,
		S_free_inodes_count: int32(n),
		S_free_blocks_count: int32(n * 3),
		S_mtime:             float32(time.Now().Unix()),
		S_umtime:            float32(time.Now().Unix()),
		S_mnt_count:         1,
		S_magic:             0xEF53,
		S_inode_size:        int32(binary.Size(estructuras.Inode{})),
		S_block_size:        int32(binary.Size(estructuras.FileBlock{})),
		S_first_ino:         inode_start,
		S_first_blo:         block_start,
		S_bm_inode_start:    bm_inode_start,
		S_bm_block_start:    bm_block_start,
		S_inode_start:       inode_start,
		S_block_start:       block_start,
	}
	return superBlock
}

func calculateStartPositions(partition *estructuras.Partition, fs string, n int32) (int32, int32, int32, int32, int32) {
	superblockSize := int32(binary.Size(estructuras.SuperBlock{}))
	journalSize := int32(binary.Size(estructuras.Journal{}))
	inodeSize := int32(binary.Size(estructuras.Inode{}))

	// Inicializar posiciones
	// EXT2
	journalStart := int32(0)
	bmInodeStart := partition.Part_start + superblockSize
	bmBlockStart := bmInodeStart + n
	inodeStart := bmBlockStart + (3 * n)
	blockStart := inodeStart + (inodeSize * n)

	// Ajustar para EXT3
	if fs == "3fs" {
		journalStart = partition.Part_start + superblockSize
		bmInodeStart = journalStart + (journalSize * n)
		bmBlockStart = bmInodeStart + n
		inodeStart = bmBlockStart + (3 * n)
		blockStart = inodeStart + (inodeSize * n)
	}

	return journalStart, bmInodeStart, bmBlockStart, inodeStart, blockStart
}
