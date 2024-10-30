package comandos

import (
	global "PROYECTO2/global"
	"errors"
)

func ParserLogout() (string, error) {
	//YA QUE SOLO PUEDE VER UNA SESION ACTIVA A LA VES
	// Verificar si hay alguna sesión activa
	if len(global.UserSessions) == 0 {
		return "No hay sesiones activas para cerrar", errors.New("no hay sesiones activas")
	}

	// si no hay sesiones activas eliminamos la primera (y única) sesión activa
	global.UserSessions = global.UserSessions[:0]

	return "Sesión cerrada exitosamente", nil

}
