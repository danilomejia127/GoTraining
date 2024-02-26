package services

import (
	"encoding/binary"
	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/apicalls"
	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/dao"
	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/db"
	"os"
	"time"

	"log"
	"strconv"

	"github.com/joho/godotenv"
)

type SellerSiteReport struct {
	WithLocalSellersFile bool      `json:"with_local_sellers_file"`
	QueryStartDate       time.Time `json:"query_start_date"`
	QueryEndDate         string    `json:"query_end_date"`
}

// GetSellerSite obtiene los seller con su site correspondiente para aquellos casos que tienen customizaciones activas en el KVS Custom
func GetSellerSite(sellerSiteReport SellerSiteReport) {
	start := time.Now()
	defer func() { logTimeExecution(start, time.Now(), "GetSellerSite") }()

	// Carga las variables de entorno desde el archivo .env
	if err := godotenv.Load("compare_pir_vs_pmc/config/local.properties"); err != nil {
		log.Panic("Error al cargar el archivo .env:", err)
	}

	dbConn, err := db.GetConnection()
	if err != nil {
		log.Panic("Error al obtener la conexión a la base de datos:", err)
	}
	defer dbConn.Close()

	// Crea un nuevo DAO para la tabla custom_pm_event
	customPMEventDAO := dao.NewCustomPMEventDAO(dbConn)

	// Define el tamaño de página y el desplazamiento inicial
	pageSize := 10000 // Tamaño de página
	offset := 0       // Desplazamiento inicial

	eventsMap := make(map[int]string)
	sellersMLA := make(map[int]string)
	sellersMLB := make(map[int]string)
	sellersMCO := make(map[int]string)
	sellersMLM := make(map[int]string)
	sellersMLC := make(map[int]string)
	sellersMEC := make(map[int]string)
	sellersMLU := make(map[int]string)
	sellersMPE := make(map[int]string)
	sellersMLV := make(map[int]string)
	sellersNoSite := make(map[int]string)

	eventsProc := 0

	for {
		// Obtiene la siguiente página de eventos
		eventsParc, err := customPMEventDAO.GetEventsSinceDate(sellerSiteReport.QueryStartDate, offset, pageSize)
		if err != nil {
			log.Panic("Error al obtener eventos:", err)
		}

		// Si no hay más eventos, sal del bucle
		if len(eventsParc) == 0 {
			break
		}

		for _, event := range eventsParc {
			eventsMap[int(event.SellerID)] = event.SiteID
		}

		eventsProc += len(eventsParc)

		log.Println("Total parcial de eventos: " + strconv.Itoa(eventsProc))

		// Incrementa el desplazamiento para la siguiente página
		offset += pageSize
	}

	log.Println("Total eventos: " + strconv.Itoa(len(eventsMap)))

	sellers, err := getSellerCustom(sellerSiteReport)
	if err != nil {
		log.Panic("Error al obtener sellers customizados:", err)
	}

	log.Println("Total Sellers Customizados: " + strconv.Itoa(len(sellers)))

	for _, seller := range sellers {
		// buscar el seller dentro del map
		if _, ok := eventsMap[seller]; ok {
			// guardar el seller en el map correspondiente
			switch eventsMap[seller] {
			case "MLA":
				sellersMLA[seller] = eventsMap[seller]
			case "MLB":
				sellersMLB[seller] = eventsMap[seller]
			case "MCO":
				sellersMCO[seller] = eventsMap[seller]
			case "MLM":
				sellersMLM[seller] = eventsMap[seller]
			case "MLC":
				sellersMLC[seller] = eventsMap[seller]
			case "MEC":
				sellersMEC[seller] = eventsMap[seller]
			case "MLU":
				sellersMLU[seller] = eventsMap[seller]
			case "MPE":
				sellersMPE[seller] = eventsMap[seller]
			case "MLV":
				sellersMLV[seller] = eventsMap[seller]
			default:
				sellersNoSite[seller] = eventsMap[seller]
			}
		} else {
			sellersNoSite[seller] = eventsMap[seller]
		}
	}

	// crear funcion que guarde el resultado de cada map en un archivo csv
	saveSellers(sellersMLA, "sellersMLA.csv")
	saveSellers(sellersMLB, "sellersMLB.csv")
	saveSellers(sellersMCO, "sellersMCO.csv")
	saveSellers(sellersMLM, "sellersMLM.csv")
	saveSellers(sellersMLC, "sellersMLC.csv")
	saveSellers(sellersMEC, "sellersMEC.csv")
	saveSellers(sellersMLU, "sellersMLU.csv")
	saveSellers(sellersMPE, "sellersMPE.csv")
	saveSellers(sellersMLV, "sellersMLV.csv")
	saveSellers(sellersNoSite, "sellersNoSite.csv")

	log.Println("Total sellers MLA: " + strconv.Itoa(len(sellersMLA)))
	log.Println("Total sellers MLB: " + strconv.Itoa(len(sellersMLB)))
	log.Println("Total sellers MCO: " + strconv.Itoa(len(sellersMCO)))
	log.Println("Total sellers MLM: " + strconv.Itoa(len(sellersMLM)))
	log.Println("Total sellers MLC: " + strconv.Itoa(len(sellersMLC)))
	log.Println("Total sellers MEC: " + strconv.Itoa(len(sellersMEC)))
	log.Println("Total sellers MLU: " + strconv.Itoa(len(sellersMLU)))
	log.Println("Total sellers MPE: " + strconv.Itoa(len(sellersMPE)))
	log.Println("Total sellers MLV: " + strconv.Itoa(len(sellersMLV)))

	log.Println("Total sellers: " + strconv.Itoa(len(sellersMLA)+len(sellersMLB)+len(sellersMCO)+len(sellersMLM)+len(sellersMLC)+len(sellersMEC)+len(sellersMLU)+len(sellersMPE)+len(sellersMLV)))
}

func saveSellers(sellers map[int]string, fileName string) {
	start := time.Now()
	defer func() { logTimeExecution(start, time.Now(), "saveSellers "+fileName) }()

	// crear el archivo csv
	file, err := os.Create("compare_pir_vs_pmc/data/" + fileName)
	if err != nil {
		log.Panic("Error al crear el archivo csv:", err)
	}
	defer file.Close()

	// escribir el encabezado del archivo csv
	_, err = file.WriteString("seller_id,site_id\n")
	if err != nil {
		log.Panic("Error al escribir el encabezado del archivo csv:", err)
	}

	// escribir los sellers en el archivo csv
	for seller, _ := range sellers {
		_, err = file.WriteString(strconv.Itoa(seller) + "\n")
		if err != nil {
			log.Panic("Error al escribir el seller en el archivo csv:", err)
		}
	}
}

func getSellerCustom(sellerSiteReport SellerSiteReport) ([]int, error) {
	start := time.Now()
	defer func() { logTimeExecution(start, time.Now(), "getSellerCustom") }()

	sellers := make([]int, 0)
	if sellerSiteReport.WithLocalSellersFile {
		sellers, err := getSellerCustomFromLocalFile()

		return sellers, err
	}

	sellers = apicalls.GetSellerCustomFromReadV2()
	saveSellerInLocalFile(sellers)

	return sellers, nil
}

func getSellerCustomFromLocalFile() ([]int, error) {
	start := time.Now()
	defer func() { logTimeExecution(start, time.Now(), "getSellerCustomFromLocalFile") }()

	// Abrir el archivo para lectura
	file, err := os.Open("compare_pir_vs_pmc/data/sellerIds.bin")
	if err != nil {
		log.Println("Error al abrir el archivo:", err)

		return nil, err
	}
	defer file.Close()

	// Leer la cantidad de IDs del archivo
	var readCount int32
	if err := binary.Read(file, binary.LittleEndian, &readCount); err != nil {
		log.Println("Error al leer la cantidad de IDs:", err)

		return nil, err
	}

	// Leer cada ID del archivo
	var readIDs []int

	for i := 0; i < int(readCount); i++ {
		var id int32
		if err := binary.Read(file, binary.LittleEndian, &id); err != nil {
			log.Println("Error al leer el ID:", err)

			return nil, err
		}

		readIDs = append(readIDs, int(id))
	}

	log.Println("IDs recuperados del archivo:", len(readIDs))

	return readIDs, nil
}

func saveSellerInLocalFile(sellers []int) {
	start := time.Now()
	defer func() { logTimeExecution(start, time.Now(), "saveSellerInLocalFile") }()

	file, err := os.Create("compare_pir_vs_pmc/data/sellerIds.bin")
	if err != nil {
		log.Println("Error al crear el archivo:", err)

		return
	}
	defer file.Close()

	count := int32(len(sellers))
	if err := binary.Write(file, binary.LittleEndian, count); err != nil {
		log.Println("Error al escribir la cantidad de IDs:", err)

		return
	}

	for _, id := range sellers {
		if err := binary.Write(file, binary.LittleEndian, int32(id)); err != nil {
			log.Println("Error al escribir el ID:", err)

			return
		}
	}

	log.Println("IDs guardados satisfactoriamente.")
}

func logTimeExecution(start time.Time, end time.Time, functionName string) {
	duration := end.Sub(start)
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) % 60
	milliseconds := int(duration.Milliseconds()) % 1000

	log.Println("Tiempo de ejecución de " + functionName + ": " + strconv.Itoa(minutes) + " minutos " + strconv.Itoa(seconds) + " segundos " + strconv.Itoa(milliseconds) + " milisegundos")
}
