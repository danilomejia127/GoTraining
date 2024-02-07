package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/services"
)

func main() {
	http.HandleFunc("/pir_vs_pmc", handlePostRequest)

	// Iniciar el servidor en el puerto 8080
	fmt.Println("Servidor escuchando en http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error al iniciar el servidor:", err)
	}
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Se esperaba un método POST", http.StatusMethodNotAllowed)
		return
	}

	token := r.Header.Get("X-Tiger-Token")
	if token == "" {
		http.Error(w, "Se esperaba un token en la cabecera X-Tiger-Token", http.StatusUnauthorized)
		return
	}

	services.SetToken(token)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error al leer el cuerpo del request", http.StatusInternalServerError)
		return
	}

	// Verificar que el cuerpo no esté vacío
	if len(body) == 0 {
		http.Error(w, "El cuerpo del request no puede estar vacío", http.StatusBadRequest)
		return
	}
	// Decodificar el cuerpo del request JSON
	var inputData services.InputData
	err = json.Unmarshal(body, &inputData)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	dataProcess := services.CompareData(inputData)

	dataResponse := services.HomologateCustomData(dataProcess)

	// Codificar la respuesta JSON
	response, err := json.Marshal(dataResponse)
	if err != nil {
		http.Error(w, "Error al codificar la respuesta JSON", http.StatusInternalServerError)
		return
	}

	// Establecer el encabezado de tipo de contenido y enviar la respuesta
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	if _, err := w.Write(response); err != nil {
		fmt.Println("Error al escribir la respuesta:", err)
		return
	}
}