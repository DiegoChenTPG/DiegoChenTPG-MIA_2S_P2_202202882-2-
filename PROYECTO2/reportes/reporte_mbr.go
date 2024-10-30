package reportes

import (
	estructuras "PROYECTO2/estructuras"
	utilidades "PROYECTO2/utilidades"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ReportMBR genera un reporte del MBR y lo guarda en la ruta especificada
func ReportMBR(mbr *estructuras.MBR, path string, mountedDiskPath string) error {
	// Crear las carpetas padre si no existen
	err := utilidades.CreateParentDirs(path)
	if err != nil {
		return err
	}

	// Obtener el nombre base del archivo sin la extensión
	dotFileName, outputImage := utilidades.GetFileNames(path)

	// Definir el contenido DOT con una tabla
	dotContent := fmt.Sprintf(`digraph G {
        node [shape=plaintext]
        tabla [label=<
            <table border="0" cellborder="1" cellspacing="0">
                <tr><td colspan="2" bgcolor="purple"> REPORTE MBR </td></tr>
                <tr><td>mbr_tamano</td><td>%d</td></tr>
                <tr><td bgcolor = "#E0B0FF">mrb_fecha_creacion</td><td bgcolor = "#E0B0FF">%s</td></tr>
                <tr><td>mbr_disk_signature</td><td>%d</td></tr>
            `, mbr.Mbr_size, time.Unix(int64(mbr.Mbr_creation_date), 0), mbr.Mbr_disk_signature)

	// Agregar las particiones a la tabla
	for i, part := range mbr.Mbr_partitions {
		/*
			// Continuar si el tamaño de la partición es -1 (o sea, no está asignada)
			if part.Part_size == -1 {
				continue
			}
		*/
		fmt.Print(i)
		// Convertir Part_name a string y eliminar los caracteres nulos
		partName := strings.TrimRight(string(part.Part_name[:]), "\x00")
		// Convertir Part_status, Part_type y Part_fit a char
		partStatus := rune(part.Part_status[0])
		partType := rune(part.Part_type[0])
		partFit := rune(part.Part_fit[0])

		// Agregar la partición a la tabla
		dotContent += fmt.Sprintf(`

		
				<tr><td colspan="2" bgcolor = "purple"> PARTICIÓN </td></tr>
				<tr><td>part_status</td><td>%c</td></tr>
				<tr><td bgcolor = "#E0B0FF">part_type</td><td bgcolor = "#E0B0FF">%c</td></tr>
				<tr><td>part_fit</td><td>%c</td></tr>
				<tr><td bgcolor = "#E0B0FF">part_start</td><td bgcolor = "#E0B0FF">%d</td></tr>
				<tr><td>part_size</td><td>%d</td></tr>
				<tr><td bgcolor = "#E0B0FF">part_name</td><td bgcolor = "#E0B0FF">%s</td></tr>
			`, partStatus, partType, partFit, part.Part_start, part.Part_size, partName)
		if string(part.Part_type[0]) == "E" {
			offset := int64(part.Part_start)
			for j := 0; j < 100; j++ {

				var ebr estructuras.EBR2

				err := ebr.Deserialize2(mountedDiskPath, offset)
				if err != nil {
					fmt.Println("Error deserializando el MBR:", err)
					return err
				}

				partName := strings.TrimRight(string(ebr.Ebr_name[:]), "\x00")
				// Convertir Part_status, Part_type y Part_fit a char
				partStatus := rune(ebr.Ebr_mount[0])
				partFit := rune(part.Part_fit[0])

				dotContent += fmt.Sprintf(`
				<tr><td colspan="2" bgcolor = "purple"> PARTICIÓN LOGICA</td></tr>
				<tr><td>EBR_mount</td><td>%c</td></tr>
				<tr><td bgcolor = "#E0B0FF">EBR_fit</td><td bgcolor = "#E0B0FF">%c</td></tr>
				<tr><td>EBR_start</td><td>%d</td></tr>
				<tr><td bgcolor = "#E0B0FF">EBR_size</td><td bgcolor = "#E0B0FF">%d</td></tr>
				<tr><td>EBR_next</td><td>%d</td></tr>
				<tr><td bgcolor = "#E0B0FF">EBR_name</td><td bgcolor = "#E0B0FF">%s</td></tr>
				`, partStatus, partFit, ebr.Ebr_start, ebr.Ebr_size, ebr.Ebr_next, partName)

				if ebr.Ebr_next == -1 {
					break
				}

				offset = offset + 30 + int64(ebr.Ebr_size)

			}

		}
	}

	// Cerrar la tabla y el contenido DOT
	dotContent += "</table>>] }"

	// Guardar el contenido DOT en un archivo
	file, err := os.Create(dotFileName)
	if err != nil {
		return fmt.Errorf("error al crear el archivo: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(dotContent)
	if err != nil {
		return fmt.Errorf("error al escribir en el archivo: %v", err)
	}

	// Ejecutar el comando Graphviz para generar la imagen
	cmd := exec.Command("dot", "-Tpng", dotFileName, "-o", outputImage)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error al ejecutar el comando Graphviz: %v", err)
	}

	return nil
}
