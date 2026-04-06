package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const baseURL = "https://testing-payment-methods.melioffice.com"

type RequestRecord struct {
	Path        string
	Application string
}

type PagingResponse struct {
	Paging struct {
		Total  int `json:"total"`
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"paging"`
}

func main() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error al obtener directorio actual: %v", err)
	}

	csvPath := filepath.Join(dir, "CannotInferPM", "logs", "excludes_by_rule.csv")
	fmt.Printf("Buscando archivo en: %s\n", csvPath)

	csvFile, err := os.Open(csvPath)
	if err != nil {
		log.Fatalf("Error al abrir el archivo CSV: %v", err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error al leer el CSV: %v", err)
	}

	fmt.Printf("Total de registros leídos: %d\n", len(records))

	requests := extractRequests(records, dir)
	callAPI(requests)
}

func extractRequests(records [][]string, dir string) []RequestRecord {
	if len(records) < 2 {
		fmt.Println("No hay suficientes registros para procesar")
		return nil
	}

	// Buscar índices de columnas dinámicamente desde el header
	header := records[0]
	colMessage := -1
	colApplication := -1
	for i, col := range header {
		switch col {
		case "message":
			colMessage = i
		case "tags.api-name":
			colApplication = i
		}
	}
	if colMessage == -1 || colApplication == -1 {
		log.Fatalf("No se encontraron las columnas requeridas (message, tags.api-name)")
	}

	urlRegex := regexp.MustCompile(`\[url:(/v1/payment_methods/search[^\]]*)\]`)

	dataDir := filepath.Join(dir, "CannotInferPM", "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Error al crear directorio: %v", err)
	}

	outPath := filepath.Join(dataDir, "clean_request.txt")
	outFile, err := os.Create(outPath)
	if err != nil {
		log.Fatalf("Error al crear archivo de salida: %v", err)
	}
	defer outFile.Close()

	writer := bufio.NewWriter(outFile)
	defer writer.Flush()

	var requests []RequestRecord
	for i := 1; i < len(records); i++ {
		record := records[i]
		if len(record) <= colApplication {
			continue
		}

		message := record[colMessage]
		application := record[colApplication]

		matches := urlRegex.FindStringSubmatch(message)
		if len(matches) < 2 {
			continue
		}
		path := matches[1]

		fmt.Fprintf(writer, "%s\t%s\n", path, application)
		requests = append(requests, RequestRecord{Path: path, Application: application})
	}

	fmt.Printf("Archivo guardado en: %s\n", outPath)
	fmt.Printf("Total de requests extraídos: %d\n", len(requests))
	return requests
}

func callAPI(requests []RequestRecord) {
	client := &http.Client{Timeout: 30 * time.Second}

	fmt.Printf("\n--- Iniciando llamadas a la API (%d requests) ---\n\n", len(requests))

	for _, r := range requests {
		fullURL := baseURL + r.Path

		req, err := http.NewRequest("GET", fullURL, nil)
		if err != nil {
			fmt.Printf("Error creando request: %v\n", err)
			continue
		}
		req.Header.Set("x-cannot-infer-response", "true")
		req.Header.Set("X-Api-Client-Application", r.Application)
		req.Header.Set("x-core-flow-type", "pm_offer")

		start := time.Now()
		resp, err := client.Do(req)
		elapsed := time.Since(start).Milliseconds()

		if err != nil {
			fmt.Printf("Error en request %s — %v\n", r.Path, err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Printf("Error leyendo respuesta: %v\n", err)
			continue
		}

		// Normalizar path para imprimir (quitar query string largo)
		printPath := r.Path
		if idx := strings.Index(printPath, "?"); idx != -1 {
			printPath = printPath[:idx] + "?" + r.Path[idx+1:]
		}

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("%s — HTTP %d — Tiempo: %d ms\n", printPath, resp.StatusCode, elapsed)
			continue
		}

		var result PagingResponse
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Printf("%s — Error parseando JSON: %v — Tiempo: %d ms\n", printPath, err, elapsed)
			continue
		}

		fmt.Printf("%s  Total=%d  Tiempo: %d ms\n", r.Path, result.Paging.Total, elapsed)
	}
}
