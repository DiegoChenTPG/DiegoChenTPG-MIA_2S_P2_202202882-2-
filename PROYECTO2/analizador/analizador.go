package analizador

import (
	comandos "PROYECTO2/comandos"
	"errors"  // Importa el paquete "errors" para manejar errores
	"fmt"     // Importa el paquete "fmt" para formatear e imprimir texto
	"os"      // Importa el paquete "os" para interactuar con el sistema operativo
	"os/exec" // Importa el paquete "os/exec" para ejecutar comandos del sistema
	"strings" // Importa el paquete "strings" para manipulación de cadenas
)

// Analyzer analiza el comando de entrada y ejecuta la acción correspondiente
func Analizador(input string) (string, interface{}, error) {

	tokens := strings.Fields(input)

	// Manejo de los espacios en el archivo
	if len(tokens) == 0 {
		return "", nil, nil
	}

	//Manejo de comentarios
	if strings.HasPrefix(tokens[0], "#") {
		var consola_retorno string
		for i := 0; i < len(tokens); i++ {
			//fmt.Println(tokens[i])
			consola_retorno += tokens[i] + " "
		}
		return consola_retorno, nil, nil
	}

	switch strings.ToLower(tokens[0]) {
	case "mkdisk":
		//fmt.Println("estamos en mkdisk")
		mkdisk, consola_retorno, error := comandos.ParserMkdisk(tokens[1:])

		if error != nil {
			return consola_retorno, nil, error

		}

		return consola_retorno, mkdisk, nil
	case "rmdisk":
		//fmt.Println("estamos en rmdisk")
		consola_retorno, error := comandos.ParserRmdisk(tokens[1:])

		if error != nil {
			return consola_retorno, nil, error
		}

		return consola_retorno, nil, nil
	case "fdisk":
		fmt.Println("estamos en fdisk")
		fdisk, consola_retorno, error := comandos.ParserFdisk(tokens[1:])
		if error != nil {
			return consola_retorno, nil, error
		}

		return consola_retorno, fdisk, nil
	case "mount":
		//fmt.Println("estamos en mount")
		mount, consola_retorno, error := comandos.ParserMount(tokens[1:])

		if error != nil {
			return consola_retorno, nil, error
		}

		return consola_retorno, mount, nil
	case "mkfs":
		mkfs, consola_retorno, error := comandos.ParserMkfs(tokens[1:])
		if error != nil {

			return consola_retorno, nil, error

		}
		return consola_retorno, mkfs, nil
	case "cat":
		cat, consola_retorno, error := comandos.ParserCat(tokens[1:])

		if error != nil {

			return consola_retorno, nil, error

		}
		return consola_retorno, cat, nil

	case "login":
		login, consola_retorno, error := comandos.ParserLogin(tokens[1:])

		if error != nil {

			return consola_retorno, nil, error

		}
		return consola_retorno, login, nil
	case "logout":
		consola_retorno, error := comandos.ParserLogout()

		if error != nil {

			return consola_retorno, nil, error

		}

		return consola_retorno, nil, nil
	case "mkgrp":
		mkgrp, consola_retorno, error := comandos.ParserMkgrp(tokens[1:])
		if error != nil {

			return consola_retorno, nil, error

		}
		return consola_retorno, mkgrp, nil
	case "rmgrp":
		rmgrp, consola_retorno, error := comandos.ParserRmgrp(tokens[1:])
		if error != nil {

			return consola_retorno, nil, error

		}
		return consola_retorno, rmgrp, nil
	case "mkusr":
		mkuser, consola_retorno, error := comandos.ParserMkusr(tokens[1:])
		if error != nil {

			return consola_retorno, nil, error

		}
		return consola_retorno, mkuser, nil
	case "rmusr":
		rmuser, consola_retorno, error := comandos.ParserRmusr(tokens[1:])
		if error != nil {

			return consola_retorno, nil, error

		}
		return consola_retorno, rmuser, nil
	case "chgrp":
		chgrp, consola_retorno, error := comandos.ParserChgrp(tokens[1:])
		if error != nil {

			return consola_retorno, nil, error

		}
		return consola_retorno, chgrp, nil
	case "mkdir":
		consola_retorno, error := comandos.ParserMkdir(tokens[1:])
		if error != nil {

			return consola_retorno, nil, error

		}
		return consola_retorno, nil, nil /*
			En este caso retorno dos nil, a diferencia de los demas comandos ya que
			en las funciones no retornamos el type MKDIR
		*/
	case "mkfile":
		consola_retorno, error := comandos.ParserMkfile(tokens[1:])
		if error != nil {

			return consola_retorno, nil, error

		}
		return consola_retorno, nil, nil /*
			En este caso retorno dos nil, a diferencia de los demas comandos ya que
			en las funciones no retornamos el type MKFILE
		*/

	//NUEVOS COMANDOS ====================
	case "unmount":
		unmount, consola_retorno, error := comandos.ParserUnmount(tokens[1:])
		if error != nil {

			return consola_retorno, nil, error

		}
		return consola_retorno, unmount, nil
	case "remove":
		fmt.Println("remove")
		return "", nil, nil
	case "edit":
		fmt.Println("edit")
		return "", nil, nil
	case "rename":
		fmt.Println("rename")
		return "", nil, nil
	case "copy":
		fmt.Println("copy")
		return "", nil, nil
	case "move":
		fmt.Println("move")
		return "", nil, nil
	case "find":
		fmt.Println("find")
		return "", nil, nil
	case "chown":
		fmt.Println("chown")
		return "", nil, nil
	case "chmod":
		fmt.Println("chmod")
		return "", nil, nil
	//====================================
	case "rep":
		//fmt.Println("estamos en rep")
		rep, consola_retorno, error := comandos.ParserRep(tokens[1:])
		if error != nil {

			return "", nil, error

		}
		return consola_retorno, rep, nil
	case " ":
		return "salto de linea", nil, nil

	case "clear":
		// Crea un comando para limpiar la terminal
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout // Redirige la salida del comando a la salida estándar
		err := cmd.Run()       // Ejecuta el comando
		if err != nil {
			// Si hay un error al ejecutar el comando, devuelve un error
			return "", nil, errors.New("no se pudo limpiar la terminal")
		}
		return "", nil, nil // Devuelve nil si el comando se ejecutó correctamente
	default:
		// Si el comando no es reconocido, devuelve un error
		fmt.Println()
		return "comando desconocido: " + tokens[0], nil, fmt.Errorf("comando desconocido: %s", tokens[0])

	}

}
