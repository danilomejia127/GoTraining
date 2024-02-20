package services

import (
	"fmt"
	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/apicalls"
	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/dao"
	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/db"
	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/entity"
	"time"

	"log"
	"strconv"

	"github.com/joho/godotenv"
)

// GetSellerSite obtiene los seller con su site correspondiente para aquellos casos que tienen customizaciones activas en el KVS Custom
func GetSellerSite() {
	// Carga las variables de entorno desde el archivo .env
	if err := godotenv.Load("compare_pir_vs_pmc/config/local.properties"); err != nil {
		log.Fatal("Error al cargar el archivo .env:", err)
	}

	dbConn, err := db.GetConnection()
	if err != nil {
		log.Fatal("Error al obtener la conexión a la base de datos:", err)
	}
	defer dbConn.Close()

	// Crea un nuevo DAO para la tabla custom_pm_event
	customPMEventDAO := dao.NewCustomPMEventDAO(dbConn)

	// Define el tamaño de página y el desplazamiento inicial
	pageSize := 10000 // Tamaño de página
	offset := 0       // Desplazamiento inicial

	// Obtén eventos personalizados creados después de una fecha específica con paginación
	date := time.Date(2024, time.February, 20, 1, 0, 28, 0, time.UTC)

	events := []entity.CustomPMEvent{}

	for {
		// Obtiene la siguiente página de eventos
		eventsParc, err := customPMEventDAO.GetEventsSinceDate(date, offset, pageSize)
		if err != nil {
			log.Fatal("Error al obtener eventos:", err)
		}

		// Si no hay más eventos, sal del bucle
		if len(eventsParc) == 0 {
			break
		}

		events = append(events, eventsParc...)
		fmt.Println("Total parcial de eventos: " + strconv.Itoa(len(events)))

		// Incrementa el desplazamiento para la siguiente página
		offset += pageSize
	}

	fmt.Println("Total eventos: " + strconv.Itoa(len(events)))

	sellers := apicalls.GetSellerCustomFromReadV2()

	fmt.Println("Total Sellers Customizados: " + strconv.Itoa(len(sellers)))
}
