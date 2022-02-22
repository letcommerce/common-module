package db

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

// ConnectAndMigrate connects to PostgresDB
func ConnectAndMigrate(dsn string, schema string, dst ...interface{}) *gorm.DB {
	connectToPublicSchema(dst, dsn)
	db := connectToServiceSchema(dst, schema, dsn)

	log.Info("Postgres connected successfully.")

	DB = db
	return db
}

func connectToServiceSchema(dst []interface{}, schemaName string, dsn string) *gorm.DB {
	log.Infof("Start Connecting to Postgres DB, schema: %v", schemaName)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{NamingStrategy: schema.NamingStrategy{
		TablePrefix:   schemaName + ".",
		SingularTable: false,
	}})
	if err != nil {
		log.Panicf("Can't connect to postgres service scehma: %v", err.Error())
	}

	db.Exec("CREATE SCHEMA IF NOT EXISTS " + schemaName)

	log.Info("Start Auto Migrating on service schema")
	err = db.AutoMigrate(dst...)
	if err != nil {
		log.Panicf("Failed to auto migrate: %v", err.Error())
	}
	return db
}

func connectToPublicSchema(dst []interface{}, dsn string) {
	log.Infof("Start Connecting to Postgres DB, public schema")
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}))
	if err != nil {
		log.Panicf("Can't connect to postgres public schema: %v", err.Error())
	}

	log.Info("Start Auto Migrating to public schema")
	err = db.AutoMigrate(dst...)
	if err != nil {
		log.Panicf("Failed to auto migrate on public schema: %v", err.Error())
	}
}
