package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/apicalls"

	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/services"
)

func main() {
	http.Handle("/pir_vs_pmc", validationMiddleware(http.HandlerFunc(handlePostRequest)))
	http.Handle("/validate_last_update", validationMiddleware(http.HandlerFunc(validateLastKVSUpdate)))
	http.Handle("/sellers_and_site_report", validationMiddleware(http.HandlerFunc(sellersAndSiteReport)))
	http.Handle("/update_special_owners_kvs", validationMiddleware(http.HandlerFunc(updateSpecialOwnersKVS)))

	// Iniciar el servidor en el puerto 8080
	log.Println("Servidor escuchando en http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	}
}

func handlePostRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Se esperaba un método POST", http.StatusMethodNotAllowed)

		return
	}

	body := r.Context().Value("body").([]byte)
	// Decodificar el cuerpo del request JSON
	var inputData services.InputData

	err := json.Unmarshal(body, &inputData)
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
		log.Fatal("Error al escribir la respuesta:", err)

		return
	}
}

func validateLastKVSUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Se esperaba un método POST", http.StatusMethodNotAllowed)

		return
	}

	body := r.Context().Value("body").([]byte)

	// Decodificar el cuerpo del request JSON
	var inputData services.InputSellers

	err := json.Unmarshal(body, &inputData)
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
		log.Fatal("Error al escribir la respuesta:", err)

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

	apicalls.SetToken(token)

	return nil
}

func sellersAndSiteReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Se esperaba un método POST", http.StatusMethodNotAllowed)

		return
	}

	body := r.Context().Value("body").([]byte)

	// Decodificar el cuerpo del request JSON
	var inputData services.SellerSiteReport

	err := json.Unmarshal(body, &inputData)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)

		return
	}

	services.GetSellerSite(inputData)
}

func updateSpecialOwnersKVS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Se esperaba un método POST", http.StatusMethodNotAllowed)

		return
	}

	body := r.Context().Value("body").([]byte)

	// Decodificar el cuerpo del request JSON
	var inputData []string

	err := json.Unmarshal(body, &inputData)
	if err != nil {
		http.Error(w, "Error al decodificar el JSON", http.StatusBadRequest)

		return
	}

	services.UpdateSpecialOwnersKVS(inputData)

	w.WriteHeader(http.StatusOK)
}

func validationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := validateBody(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		err = validateToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)

			return
		}

		// Store the body in the request context
		ctx := context.WithValue(r.Context(), "body", body)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
