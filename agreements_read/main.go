package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Estructuras para mapear la respuesta JSON
type Agreement struct {
	ID                           string               `json:"id"`
	Issuer                       Issuer               `json:"issuer"`
	MaxInstallments              string               `json:"max_installments"`
	Marketplace                  string               `json:"marketplace"`
	PaymentMethods               []PaymentMethod      `json:"payment_methods"`
	Legals                       string               `json:"legals"`
	InterestDeductionByCollector bool                 `json:"interest_deduction_by_collector"`
	AllTotalFinancialCost        []int                `json:"all_total_financial_cost"`
	StartDate                    time.Time            `json:"start_date"`
	ExpirationDate               time.Time            `json:"expiration_date"`
	TotalFinancialCost           int                  `json:"total_financial_cost"`
	Site                         string               `json:"site"`
	LogoURLs                     []string             `json:"logo_urls"`
	FinancingConditions          []FinancingCondition `json:"financing_conditions"`
	Owner                        string               `json:"owner"`
	ApplicationDays              []string             `json:"application_days"`
	FinancingGroups              []interface{}        `json:"financing_groups"`
}

type Issuer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

type PaymentMethod struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Brand           string  `json:"brand"`
	Product         string  `json:"product"`
	SecureThumbnail *string `json:"secure_thumbnail"`
	Thumbnail       *string `json:"thumbnail"`
}

type FinancingCondition struct {
	Installments string  `json:"installments"`
	InterestRate float64 `json:"interest_rate"`
	InterestType string  `json:"interest_type"`
}

// Función para consumir la API de agreements
func getAgreements() ([]Agreement, error) {
	url := "https://agreements-read.melioffice.com/agreements-read/MLA/credit_card_promotions"

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
	var agreements []Agreement
	err = json.Unmarshal(body, &agreements)
	if err != nil {
		return nil, fmt.Errorf("error parseando JSON: %v", err)
	}

	return agreements, nil
}

// Función para guardar los agreements en un archivo CSV
func saveToCSV(agreements []Agreement, filename string) error {
	// Crear la carpeta si no existe
	dir := "data"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("error creando directorio: %v", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creando archivo CSV: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir encabezados
	headers := []string{"site_id", "issuer_id"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("error escribiendo encabezados: %v", err)
	}

	// Escribir datos
	for _, agreement := range agreements {
		record := []string{
			agreement.Site,                         // site_id (usando el campo site)
			fmt.Sprintf("%d", agreement.Issuer.ID), // issuer.id
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error escribiendo registro: %v", err)
		}
	}

	return nil
}

func main() {
	// Llamar a la función para obtener los agreements
	agreements, err := getAgreements()
	if err != nil {
		fmt.Printf("Error obteniendo agreements: %v\n", err)
		return
	}

	// Guardar en CSV
	err = saveToCSV(agreements, "data/agreements_data.csv")
	if err != nil {
		fmt.Printf("Error guardando CSV: %v\n", err)
		return
	}

	fmt.Printf("Datos guardados exitosamente en data/agreements_data.csv\n")
	fmt.Printf("Total de registros: %d\n", len(agreements))
}
