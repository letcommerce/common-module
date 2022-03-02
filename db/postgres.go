package db

import (
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	DB *gorm.DB
)

// ConnectAndMigrate connects to PostgresDB
func ConnectAndMigrate(dsn string, schema string, useCloudSql bool, dst ...interface{}) *gorm.DB {
	connectToPublicSchema(dst, dsn, useCloudSql)
	db := connectToServiceSchema(dst, schema, dsn, useCloudSql)

	log.Info("Postgres connected successfully.")

	DB = db
	return db
}

func connectToServiceSchema(dst []interface{}, schemaName string, dsn string, useCloudSql bool) *gorm.DB {
	log.Infof("Start Connecting to Postgres DB, schema: %v", schemaName)
	driverName := ""
	if useCloudSql {
		driverName = "cloudsqlpostgres"
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName:           driverName,
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

func connectToPublicSchema(dst []interface{}, dsn string, useCloudSql bool) {
	log.Infof("Start Connecting to Postgres DB, public schema")
	driverName := ""
	if useCloudSql {
		driverName = "cloudsqlpostgres"
	}
	db, err := gorm.Open(postgres.New(postgres.Config{
		DriverName:           driverName,
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
