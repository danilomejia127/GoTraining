package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/mercadolibre/GoTraining/compare_pir_vs_pmc/entity"
	"log"
	"time"
)

// CustomPMEventDAO proporciona métodos para interactuar con la tabla custom_pm_event
type CustomPMEventDAO struct {
	db *sqlx.DB
}

// NewCustomPMEventDAO crea un nuevo DAO para la tabla custom_pm_event
func NewCustomPMEventDAO(db *sqlx.DB) *CustomPMEventDAO {
	return &CustomPMEventDAO{db}
}

// GetEventsSinceDate obtiene eventos personalizados creados después de una fecha específica con paginación
func (dao *CustomPMEventDAO) GetEventsSinceDate(date time.Time, offset, limit int) ([]entity.CustomPMEvent, error) {
	query := "SELECT DISTINCT site_id, seller_id FROM custom_pm_event WHERE site_id = 'MCO' and date_created > ? LIMIT ?, ?"

	rows, err := dao.db.Query(query, date, offset, limit)
	if err != nil {
		log.Println("Error al ejecutar la consulta:", err)

		return nil, err
	}
	defer rows.Close()

	var events []entity.CustomPMEvent

	for rows.Next() {
		var event entity.CustomPMEvent

		err := rows.Scan(&event.SiteID, &event.SellerID)
		if err != nil {
			log.Println("Error al escanear fila:", err)

			return nil, err
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error durante la iteración de filas:", err)

		return nil, err
	}

	return events, nil
}
