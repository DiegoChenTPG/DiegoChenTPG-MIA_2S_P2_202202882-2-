package estructuras

import (
	"bytes"           // Paquete para manipulación de buffers
	"encoding/binary" // Paquete para codificación y decodificación de datos binarios
	"fmt"             // Paquete para formateo de E/S
	"os"              // Paquete para funciones del sistema operativo
	// Paquete para manipulación de tiempo
)

type EBR2 struct {
	Ebr_mount [1]byte
	Ebr_fit   [1]byte
	Ebr_start int32
	Ebr_size  int32
	Ebr_next  int32
	Ebr_name  [16]byte
}

// Serialize escribe la estructura SuperBlock en un archivo binario en la posición especificada
func (ebr *EBR2) Serialize2(path string, offset int64) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Serializar la estructura SuperBlock directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, ebr)
	if err != nil {
		return err
	}

	return nil
}

// Deserialize lee la estructura SuperBlock desde un archivo binario en la posición especificada
func (ebr *EBR2) Deserialize2(path string, offset int64) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Mover el puntero del archivo a la posición especificada
	_, err = file.Seek(offset, 0)
	if err != nil {
		return err
	}

	// Obtener el tamaño de la estructura SuperBlock
	sbSize := binary.Size(ebr)
	if sbSize <= 0 {
		return fmt.Errorf("invalid SuperBlock size: %d", sbSize)
	}

	// Leer solo la cantidad de bytes que corresponden al tamaño de la estructura SuperBlock
	buffer := make([]byte, sbSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserializar los bytes leídos en la estructura SuperBlock
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, ebr)
	if err != nil {
		return err
	}

	return nil
}
