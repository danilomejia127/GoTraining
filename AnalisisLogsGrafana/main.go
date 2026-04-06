package main

import (
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

type LogEntry struct {
	Timestamp      string
	URL            string
	Bins           string
	ProductID      string
	FinancingGroup string
	CallerID       string
	IssuerName     string
	CardTypeCode   string
}

func main() {
	// Obtener el directorio actual del ejecutable
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error al obtener directorio actual: %v", err)
	}
	// Construir la ruta absoluta al archivo CSV
	csvPath := filepath.Join(dir, "AnalisisLogsGrafana", "logs", "logs-poa-excluded_by_rule.csv")
	fmt.Printf("Buscando archivo en: %s\n", csvPath)

	// Leer el archivo CSV
	csvFile, err := os.Open(csvPath)
	if err != nil {
		fmt.Println("Error al abrir el archivo CSV:", err)
		return
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	// Leer todas las filas
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error al leer el CSV: %v", err)
	}

	fmt.Printf("Total de registros leídos: %d\n", len(records))

	// Procesar los registros y extraer URLs
	processLogsWithURLs(records)
}

func processLogsWithURLs(records [][]string) {
	if len(records) < 2 {
		fmt.Println("No hay suficientes registros para procesar")
		return
	}

	// Regex para extraer la URL del mensaje
	urlRegex := regexp.MustCompile(`\[url:([^\]]+)\]`)

	var logEntries []LogEntry

	// Procesar cada registro (saltando el header)
	for i := 1; i < len(records); i++ {
		fmt.Printf("Procesando registro %d\n", i)
		record := records[i]
		if len(record) < 2 {
			continue
		}

		message := record[1] // El campo message está en la segunda columna

		// Extraer URL usando regex
		urlMatches := urlRegex.FindStringSubmatch(message)
		var url string
		var bins string
		var productID string
		var financingGroup string
		var callerID string

		if len(urlMatches) > 1 {
			url = urlMatches[1]
		}

		// extraer los parametros de la url y guardarlos en un LogEntry
		// Validar que la URL tenga parámetros antes de intentar extraerlos
		if strings.Contains(url, "?") {
			paramsParts := strings.Split(url, "?")
			if len(paramsParts) > 1 {
				paramsFull := paramsParts[1]
				paramsSplit := strings.Split(paramsFull, "&")

				for _, param := range paramsSplit {
					paramParts := strings.Split(param, "=")
					if len(paramParts) > 1 {
						switch paramParts[0] {
						case "bins":
							bins = paramParts[1]
						case "product_id":
							productID = paramParts[1]
						case "financing_group":
							financingGroup = paramParts[1]
						case "caller.id":
							callerID = paramParts[1]
						}
					}
				}
			}
		}

		// Consumir la API de binapi para obtener el issuer name y card type code
		binAPIInfo, err := getBINAPIInfo(bins, "MLB")
		if err != nil {
			fmt.Printf("Error al consumir la API de binapi: %v\n", err)
			continue
		}

		// Usar el primer elemento del array settings si existe
		var issuerName, cardTypeCode string
		if len(binAPIInfo.Settings) > 0 {
			issuerName = binAPIInfo.Settings[0].IssuerName
			cardTypeCode = binAPIInfo.Settings[0].CardTypeCode
		}

		entry := LogEntry{
			Timestamp:      record[0],
			URL:            url,
			Bins:           bins,
			ProductID:      productID,
			FinancingGroup: financingGroup,
			CallerID:       callerID,
			IssuerName:     issuerName,
			CardTypeCode:   cardTypeCode,
		}

		logEntries = append(logEntries, entry)
	}

	saveToCSV(logEntries)
}

func saveToCSV(logEntries []LogEntry) {
	// Crear el directorio si no existe
	dataDir := "AnalisisLogsGrafana/data"
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("Error al crear directorio: %v", err)
	}

	csvFile, err := os.Create(filepath.Join(dataDir, "log_entries.csv"))
	if err != nil {
		log.Fatalf("Error al crear el archivo CSV: %v", err)
	}
	defer csvFile.Close()

	// Crear el writer CSV
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Escribir el header
	header := []string{"Timestamp", "URL", "Bins", "ProductID", "FinancingGroup", "CallerID", "IssuerName", "CardTypeCode"}
	if err := writer.Write(header); err != nil {
		log.Fatalf("Error al escribir header: %v", err)
	}

	// Escribir los datos
	for _, entry := range logEntries {
		record := []string{
			entry.Timestamp,
			entry.URL,
			entry.Bins,
			entry.ProductID,
			entry.FinancingGroup,
			entry.CallerID,
			entry.IssuerName,
			entry.CardTypeCode,
		}
		if err := writer.Write(record); err != nil {
			log.Fatalf("Error al escribir registro: %v", err)
		}
	}

	fmt.Printf("Archivo CSV guardado exitosamente: %s\n", filepath.Join(dataDir, "log_entries.csv"))
	fmt.Printf("Total de registros guardados: %d\n", len(logEntries))
}

// Estructura que coincide con la respuesta JSON de la API
type BinAPIResponse struct {
	Bin          int        `json:"bin"`
	RequestedBin int        `json:"requested_bin"`
	Version      int        `json:"version"`
	Settings     []Settings `json:"settings"`
}

type Settings struct {
	IssuerName   string `json:"issuer_name"`
	CardTypeCode string `json:"card_type_code"`
}

// Función para consumir la API de agreements
func getBINAPIInfo(bin, site_id string) (*BinAPIResponse, error) {
	url := fmt.Sprintf("https://production-binapi.melioffice.com/binapi/v1/search/%s/%s?with=fallback", bin, site_id)

	// Crear cliente HTTP con timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Hacer la llamada GET
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error haciendo la llamada HTTP: %v", err)
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo la respuesta: %v", err)
	}

	// Verificar el status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error en la respuesta: %d - %s", resp.StatusCode, string(body))
	}

	// Parsear el JSON
	var binAPIResponse BinAPIResponse
	err = json.Unmarshal(body, &binAPIResponse)
	if err != nil {
		return nil, fmt.Errorf("error parseando JSON: %v", err)
	}

	return &binAPIResponse, nil
}
