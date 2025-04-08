package main

import (
	analizador "PROYECTO2/analizador"
	"PROYECTO2/global"
	"PROYECTO2/utilidades"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

// Estructuras para el login
type LoginRequest struct {
	IDParticion   string `json:"ID_particion"`
	NombreUsuario string `json:"nombre_usuario"`
	Password      string `json:"password"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type LogoutResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type DiscosResponse struct {
	Discos []string `json:"discos"`
}

func main() {

	http.HandleFunc("/api/execute", withCORS(ejecutarComando))        // Define el endpoint
	http.HandleFunc("/api/login", withCORS(controlLogin))             // Endpoint para login
	http.HandleFunc("/api/logout", withCORS(controlLogout))           // Endpoint para logout
	http.HandleFunc("/api/discos", withCORS(ObtenerDiscosHandler))    // Endpoint para el visualizador de discos
	http.HandleFunc("/api/particiones", withCORS(ParticionesHandler)) // Endpoint para el visualizador de discos

	fmt.Println("Servidor escuchando en http://3.23.105.151:8080") //es un print, pero igual cambiar a localhost cuando se deje de trabajar con AWS para evitar confusion
	err := http.ListenAndServe(":8080", nil)                       // Inicia el servidor
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}

}

// Middleware para manejar CORS
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Configurar los encabezados CORS
		w.Header().Set("Access-Control-Allow-Origin", "*") // Permite todos los orígenes
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			// Responder a las solicitudes preflight con un código 204 (No Content)
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Llamar al manejador principal
		next(w, r)
	}
}

func ejecutarComando(w http.ResponseWriter, r *http.Request) {

	var req CommandRequest

	// Decodifica la solicitud JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Aquí es donde procesas el comando
	var salida_consola string
	comandos := strings.Split(req.Command, "\n") // SE DIVIDE POR LINEAS LA ENTRADA PARA ANALIZAR 1 A 1 y no solo el ultimo

	// Crea la respuesta
	for _, comando := range comandos {

		//fmt.Println("==================")

		// Ignorar líneas vacías
		/*
			if strings.TrimSpace(comando) == "" {
				continue
			}
		*/

		//Se analiza linea por linea
		salida, _, err := analizador.Analizador(comando)
		if err != nil {
			// Si hay un error, almacenar el mensaje de error en lugar del resultado
			fmt.Printf("Error: %s", err.Error())
			salida_consola = salida_consola + salida + "\n"
			salida_consola += ""
			continue
		}
		//fmt.Println("IMPRIMIENTO LAS SALIDAS")
		//fmt.Println(salida)
		salida_consola = salida_consola + salida + "\n"
		salida_consola += ""
	}
	/*
		if err != nil {
			res := CommandResponse{
				Output: salida_consola,
				Error:  fmt.Sprintf("Error al ejecutar el comando: %v", err),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)
			return
		}
	*/

	//fmt.Println(salida_consola)
	// Preparar la respuesta
	res := CommandResponse{
		Output: salida_consola,
	}

	// Devuelve la respuesta como JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)

}

// Manejador para el login
func controlLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	fmt.Println("ENTRAMOS")
	// Decodificar la solicitud JSON
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mensaje, verificador := global.Verificar_login(req.NombreUsuario, req.Password, req.IDParticion)

	if verificador {
		// Aquí puedes implementar tu lógica real de validación de usuario y contraseña
		// Responder con éxito si las credenciales son correctas
		res := LoginResponse{
			Success: verificador,
			Message: mensaje,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		fmt.Println("VALIDADO")
		return
	}

	// Responder con error si las credenciales son incorrectas
	res := LoginResponse{
		Success: verificador,
		Message: mensaje,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	fmt.Println(mensaje)
}

// Manejador para el logout
func controlLogout(w http.ResponseWriter, r *http.Request) {
	// Aquí puedes limpiar cualquier sesión o token si estás usando uno en el futuro
	salida, _, _ := analizador.Analizador("logout")

	if salida == "Sesión cerrada exitosamente" {
		res := LogoutResponse{
			Success: true,
			Message: "Sesión cerrada exitosamente",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(res)
		fmt.Println("")
		fmt.Println("Sesión cerrada")
		return
	}

	res := LogoutResponse{
		Success: true,
		Message: "Error al cerrar sesion",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	fmt.Println("")
	fmt.Println("Sesión no cerrada")

}

// MANEJO DE DISCOS
// Handler para obtener los nombres de los discos
func ObtenerDiscosHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	discos := utilidades.ObtenerNombresDiscos()

	// Serializa y envía la respuesta en JSON
	if err := json.NewEncoder(w).Encode(discos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// MANEJO DE PARTICIONES
// Handler para obtener particiones del disco específico
func ParticionesHandler(w http.ResponseWriter, r *http.Request) {
	disco := r.URL.Query().Get("disco")
	fmt.Println("IMPRIMIENDO DISCO " + disco)
	if disco == "" {
		http.Error(w, "Disco no especificado", http.StatusBadRequest)
		return
	}
	path_disco, _ := utilidades.ObtenerPathDisco(disco)
	particiones := global.ObtenerParticionesPorDisco(path_disco)

	json.NewEncoder(w).Encode(particiones)

}
