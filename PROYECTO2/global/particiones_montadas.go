package global

import (
	//estructuras "PROYECTO1/estructuras"
	estructuras "PROYECTO2/estructuras"
	"errors"
)

// Modificacion del guardado para que sea el nombre y el path
type MountInfo struct {
	Path string
	Name string
}

// Carnet de estudiante
const Carnet string = "82" // 202202882

// Declaración de las particiones montadas
var (
	MountedPartitions map[string]MountInfo = make(map[string]MountInfo)
)

// GetMountedPartition obtiene la partición montada con el id especificado
func GetMountedPartition(id string) (*estructuras.Partition, string, error) {
	// Obtener el path de la partición montada
	mountInfo, exists := MountedPartitions[id]
	if !exists {
		return nil, "La particion no esta montada", errors.New("la partición no está montada")
	}

	// Crear una instancia de MBR
	var mbr estructuras.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(mountInfo.Path)
	if err != nil {
		return nil, "", err
	}

	// Buscar la partición con el id especificado
	partition, err := mbr.GetPartitionByID(id)
	if partition == nil {
		return nil, "", err
	}

	return partition, mountInfo.Path, nil
}

// GetMountedMBR obtiene el MBR de la partición montada con el id especificado
func GetMountedPartitionRep(id string) (*estructuras.MBR, *estructuras.SuperBlock, string, error) {
	// Obtener el path de la partición montada
	mountInfo, exists := MountedPartitions[id]
	if !exists {
		return nil, nil, "La particion no esta montada", errors.New("la partición no está montada")
	}

	// Crear una instancia de MBR
	var mbr estructuras.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(mountInfo.Path)
	if err != nil {
		return nil, nil, "", err
	}

	// Buscar la partición con el id especificado
	partition, err := mbr.GetPartitionByID(id)
	if partition == nil {
		return nil, nil, "", err
	}

	// Crear una instancia de SuperBlock
	var sb estructuras.SuperBlock

	// Deserializar la estructura SuperBlock desde un archivo binario
	err = sb.Deserialize(mountInfo.Path, int64(partition.Part_start))
	if err != nil {
		return nil, nil, "", err
	}

	return &mbr, &sb, mountInfo.Path, nil
}

func GetMountedPartitionSuperblock(id string) (*estructuras.SuperBlock, *estructuras.Partition, string, error) {
	// Obtener el MountInfo de la partición montada
	mountInfo, exists := MountedPartitions[id]
	if !exists {
		return nil, nil, "", errors.New("la partición no está montada")
	}

	// Crear una instancia de MBR
	var mbr estructuras.MBR

	// Deserializar la estructura MBR desde un archivo binario
	err := mbr.Deserialize(mountInfo.Path)
	if err != nil {
		return nil, nil, "", err
	}

	// Buscar la partición con el id especificado
	partition, err := mbr.GetPartitionByID(id)
	if partition == nil {
		return nil, nil, "", err
	}

	// Crear una instancia de SuperBlock
	var sb estructuras.SuperBlock

	// Deserializar la estructura SuperBlock desde un archivo binario, usando el Part_start de la partición
	err = sb.Deserialize(mountInfo.Path, int64(partition.Part_start))
	if err != nil {
		return nil, nil, "", err
	}

	return &sb, partition, mountInfo.Path, nil
}

// DeleteMountedPartition elimina la partición montada con el ID especificado
func DeleteMountedPartition(id string) error {
	// Verificar si la partición está montada en la lista de montajes
	_, exists := MountedPartitions[id]
	if !exists {
		return errors.New("la partición no está montada")
	}

	// Eliminar la partición de la lista de montajes
	delete(MountedPartitions, id)
	return nil
}

// ObtenerParticionesPorDisco recibe un path y devuelve un slice con los nombres de las particiones montadas en ese disco
func ObtenerParticionesPorDisco(path string) []string {
	var partitionNames []string

	// Recorrer las particiones montadas y verificar si coinciden con el path especificado
	for _, mountInfo := range MountedPartitions {
		if mountInfo.Path == path {
			partitionNames = append(partitionNames, mountInfo.Name)
		}
	}

	return partitionNames
}
