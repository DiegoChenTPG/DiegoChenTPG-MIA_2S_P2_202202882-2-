package estructuras

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

type EBR struct {
	Part_status [1]byte  // Indica si la partición está montada o no
	Part_fit    [1]byte  // Tipo de ajuste de la partición (B, F, W)
	Part_start  int32    // Byte en el que inicia la partición
	Part_size   int32    // Tamaño de la partición en bytes
	Part_next   int32    // Byte en el que se encuentra el próximo EBR (-1 si no hay)
	Part_name   [16]byte // Nombre de la partición
}

// SerializeEBR escribe la estructura EBR a un archivo binario en la posición exacta
func (ebr *EBR) Serialize(path string, offset int64) error {
	file, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero al offset donde se debe escribir el EBR
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Serializar la estructura EBR directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, ebr)
	if err != nil {
		return err
	}

	return nil
}

// DeserializeEBR lee una estructura EBR desde un archivo binario
func (ebr *EBR) Deserialize(path string, offset int64) error {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("ERROR 1")
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo al offset dado (inicio de este EBR)
	_, err = file.Seek(offset, 0)
	if err != nil {
		fmt.Println("ERROR 2")
		return err
	}

	// Leer la estructura EBR desde el archivo
	err = binary.Read(file, binary.LittleEndian, ebr)
	if err != nil {
		fmt.Println("ERROR 3")
		return err
	}

	return nil
}

func (ebr *EBR) Print() string {
	// Convertir Part_fit a char
	partFit := rune(ebr.Part_fit[0])
	// Convertir Part_status a char
	partStatus := rune(ebr.Part_status[0])
	// Convertir Part_name a string
	partName := strings.Trim(string(ebr.Part_name[:]), "\x00 ")

	mensaje := fmt.Sprint("Part Status: ", partStatus, "\n")
	mensaje += fmt.Sprint("Part Fit: ", partFit, "\n")
	mensaje += fmt.Sprint("Part Start: ", ebr.Part_start, "\n")
	mensaje += fmt.Sprint("Part Size: ", ebr.Part_size, "\n")
	mensaje += fmt.Sprint("Part Next: ", ebr.Part_next, "\n")
	mensaje += fmt.Sprint("Part Name: ", partName, "\n")

	return mensaje
}

func (e *EBR) CreateLogicalPartition(partStart, partSize, nextEBR int, partFit, partName string) {
	// Asignar status de la partición
	e.Part_status[0] = '0' // El valor '0' indica que la partición ha sido creada

	// Asignar el byte de inicio de la partición
	e.Part_start = int32(partStart)

	// Asignar el tamaño de la partición
	e.Part_size = int32(partSize)

	// Asignar el ajuste de la partición
	if len(partFit) > 0 {
		e.Part_fit[0] = partFit[0]
	}

	// Asignar el nombre de la partición
	copy(e.Part_name[:], partName)

	// Asignar el puntero al siguiente EBR
	e.Part_next = int32(nextEBR)
}
