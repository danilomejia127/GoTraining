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
	http.HandleFunc("/validate_last_update", validateLastKVSUpdate)

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

	err := validateToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	body, err := validateBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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

func validateLastKVSUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Se esperaba un método POST", http.StatusMethodNotAllowed)
		return
	}

	err := validateToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	body, err := validateBody(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Decodificar el cuerpo del request JSON
	var inputData services.InputSellers

	err = json.Unmarshal(body, &inputData)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)
		return
	}

	dataResponse := services.ValidateLastUpdateIntoKVS(inputData)

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

func validateBody(r *http.Request) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("error al leer el cuerpo del request")
	}

	// Verificar que el cuerpo no esté vacío
	if len(body) == 0 {
		return []byte{}, fmt.Errorf("el cuerpo del request no puede estar vacío")
	}

	return body, nil
}

func validateToken(r *http.Request) error {
	token := r.Header.Get("X-Tiger-Token")
	if token == "" {
		return fmt.Errorf("se esperaba un token en la cabecera X-Tiger-Token")
	}

	services.SetToken(token)

	return nil
}
