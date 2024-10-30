package estructuras

import (
	"bytes"           // Paquete para manipulación de buffers
	"encoding/binary" // Paquete para codificación y decodificación de datos binarios
	"errors"
	"fmt" // Paquete para formateo de E/S
	"os"  // Paquete para funciones del sistema operativo
	"strings"
	"time" // Paquete para manipulación de tiempo
)

type MBR struct {
	Mbr_size           int32        // Tamaño del MBR en bytes
	Mbr_creation_date  float32      // Fecha y hora de creación del MBR
	Mbr_disk_signature int32        // Firma del disco
	Mbr_disk_fit       [1]byte      // Tipo de ajuste
	Mbr_partitions     [4]Partition // Particiones del MBR
}

// SerializeMBR escribe la estructura MBR al inicio de un archivo binario
func (mbr *MBR) Serialize(path string) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Serializar la estructura MBR directamente en el archivo
	err = binary.Write(file, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

// DeserializeMBR lee la estructura MBR desde el inicio de un archivo binario
func (mbr *MBR) Deserialize(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Obtener el tamaño de la estructura MBR
	mbrSize := binary.Size(mbr)
	if mbrSize <= 0 {
		return fmt.Errorf("invalid MBR size: %d", mbrSize)
	}

	// Leer solo la cantidad de bytes que corresponden al tamaño de la estructura MBR
	buffer := make([]byte, mbrSize)
	_, err = file.Read(buffer)
	if err != nil {
		return err
	}

	// Deserializar los bytes leídos en la estructura MBR
	reader := bytes.NewReader(buffer)
	err = binary.Read(reader, binary.LittleEndian, mbr)
	if err != nil {
		return err
	}

	return nil
}

func (mbr *MBR) Print() string {
	// Convertir Mbr_creation_date a time.Time
	creationTime := time.Unix(int64(mbr.Mbr_creation_date), 0)

	// Convertir Mbr_disk_fit a char
	diskFit := rune(mbr.Mbr_disk_fit[0])
	/*
		fmt.Printf("MBR Size: %d\n", mbr.Mbr_size)
		fmt.Printf("Creation Date: %s\n", creationTime.Format(time.RFC3339))
		fmt.Printf("Disk Signature: %d\n", mbr.Mbr_disk_signature)
		fmt.Printf("Disk Fit: %c\n", diskFit)
	*/
	mensaje := fmt.Sprint("MBR Size: ", mbr.Mbr_size, "\n")
	mensaje += fmt.Sprint("Creation Date: ", creationTime.Format(time.RFC3339), "\n")
	mensaje += fmt.Sprint("Disk Signature: ", mbr.Mbr_disk_signature, "\n")
	mensaje += fmt.Sprint("Disk Fit: ", diskFit, "\n")

	return mensaje

}

// Método para obtener la primera partición disponible
func (mbr *MBR) GetFirstAvailablePartition() (*Partition, int, int) {
	// Calcular el offset para el start de la partición
	offset := binary.Size(mbr) // Tamaño del MBR en bytes

	// Recorrer las particiones del MBR
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		// Si el start de la partición es -1, entonces está disponible
		if mbr.Mbr_partitions[i].Part_start == -1 {
			// Devolver la partición, el offset y el índice
			return &mbr.Mbr_partitions[i], offset, i
		} else {
			// Calcular el nuevo offset para la siguiente partición, es decir, sumar el tamaño de la partición
			offset += int(mbr.Mbr_partitions[i].Part_size)
		}
	}
	return nil, -1, -1
}

func (mbr *MBR) GetPartitionByName(name string) (*Partition, int, error) {
	// Recorrer las particiones del MBR
	fmt.Println("GET PARTITIONBYNAME")
	for i, partition := range mbr.Mbr_partitions {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partitionName := strings.Trim(string(partition.Part_name[:]), "\x00 ")
		// Convertir el nombre de la partición a string y eliminar los caracteres nulos
		inputName := strings.Trim(name, "\x00 ")
		// Si el nombre de la partición coincide, devolver la partición y el índice
		if strings.EqualFold(partitionName, inputName) {
			return &partition, i, nil
		}
	}
	return nil, -1, errors.New("particion no encontrada")
}

func (mbr *MBR) GetPartitionByID(id string) (*Partition, error) {
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		// Convertir Part_name a string y eliminar los caracteres nulos
		partitionID := strings.Trim(string(mbr.Mbr_partitions[i].Part_id[:]), "\x00 ")
		// Convertir el id a string y eliminar los caracteres nulos
		inputID := strings.Trim(id, "\x00 ")
		// Si el nombre de la partición coincide, devolver la partición
		if strings.EqualFold(partitionID, inputID) {
			return &mbr.Mbr_partitions[i], nil
		}
	}
	return nil, errors.New("partición no encontrada")
}

func (mbr *MBR) GetExtendedPartition() (*Partition, int, error) {
	// Recorrer las particiones del MBR
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		// Verificar si la partición es de tipo 'E' (extendida)
		if mbr.Mbr_partitions[i].Part_type[0] == 'E' {
			return &mbr.Mbr_partitions[i], i, nil
		}
	}
	return nil, -1, fmt.Errorf("no se encontró una partición extendida en el MBR")
}

func (mbr *MBR) PrintPartitions() string {
	var mensaje string
	for i, partition := range mbr.Mbr_partitions {
		// Convertir Part_status, Part_type y Part_fit a char
		partStatus := rune(partition.Part_status[0])
		partType := rune(partition.Part_type[0])
		partFit := rune(partition.Part_fit[0])

		// Convertir Part_name a string
		partName := string(partition.Part_name[:])

		fmt.Printf("Partition %d:\n", i+1)
		//COmentar de nuevo luego
		fmt.Printf("  Status: %c\n", partStatus)
		fmt.Printf("  Type: %c\n", partType)
		fmt.Printf("  Fit: %c\n", partFit)
		fmt.Printf("  Start: %d\n", partition.Part_start)
		fmt.Printf("  Size: %d\n", partition.Part_size)
		fmt.Printf("  Name: %s\n", partName)
		fmt.Printf("  Correlative: %d\n", partition.Part_correlative)
		fmt.Printf("  ID: %d\n", partition.Part_id)

		mensaje = fmt.Sprint("Partition: \n", mbr.Mbr_size, "\n")
		mensaje += fmt.Sprint("Status: \n", partStatus, "\n")
		mensaje += fmt.Sprint("Type: \n", partType, "\n")
		mensaje += fmt.Sprint("Fit: \n", partFit, "\n")
		mensaje += fmt.Sprint("Start: \n", partition.Part_start, "\n")
		mensaje += fmt.Sprint("Size: \n", partition.Part_start, "\n")
		mensaje += fmt.Sprint("Name: \n", partName, "\n")
		mensaje += fmt.Sprint("Name: \n", partition.Part_correlative, "\n")
		mensaje += fmt.Sprint("ID: \n", partition.Part_id, "\n")
	}
	return mensaje
}

// Metodo verificar si existe una partición extendida
func (mbr *MBR) GetExtended() (*Partition, int, int) {
	offset := binary.Size(mbr) // Tamaño del MBR en bytes

	// Recorrer las particiones del MBR
	for i := 0; i < len(mbr.Mbr_partitions); i++ {
		// Si el tipo de partición es E de extendida
		if string(mbr.Mbr_partitions[i].Part_type[0]) == "E" {
			// Devolver la particion extendida, el puntero y el indice
			return &mbr.Mbr_partitions[i], offset, i
		} else {
			// Calcular el nuevo offset para la siguiente partición, es decir, sumar el tamaño de la partición
			offset += int(mbr.Mbr_partitions[i].Part_size)
		}
	}
	//si ninguna partición tiene un tipo E, aun no existe la partición extendida, se devuelve nulo
	return nil, -1, -1
}
