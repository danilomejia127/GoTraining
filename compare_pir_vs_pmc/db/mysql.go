package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/newrelic/go-agent/v3/integrations/nrmysql"
	"log"
	"os"
	"sync"
)

var (
	dbConn *sqlx.DB
	once   sync.Once
)

// GetConnection establece una conexión a la base de datos MySQL y la devuelve.
func GetConnection() (*sqlx.DB, error) {
	// Lee las credenciales de la base de datos desde el archivo de propiedades
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	urlAndPort := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	// Configura la cadena de conexión a la base de datos
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", dbUser, dbPassword, urlAndPort, dbName)

	// Abre una conexión a la base de datos MySQL
	db, err := sqlx.Open("nrmysql", connectionString)
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)

		return nil, err
	}

	// Intenta conectar a la base de datos
	err = db.Ping()
	if err != nil {
		log.Fatal("Error al conectarse a la base de datos:", err)

		return nil, err
	}

	return db, nil
}

func InitDB() (*sqlx.DB, error) {
	var err error

	once.Do(func() {
		dbConn, err = GetConnection()
	})
	if err != nil {
		log.Fatal("Error al obtener la conexión a la base de datos:", err)
	}

	return dbConn, err
}
