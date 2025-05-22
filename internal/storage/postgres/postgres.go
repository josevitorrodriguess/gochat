package postgres

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	MaxOpenConn     int16
	MaxIdleConn     int16
	ConnMaxLifeTime time.Duration
	ConnMaxIdleTime time.Duration
}
type PostgresVariables struct {
	DbName     string
	DbUser     string
	DbPassword string
	DbPort     string
	UseSSL     bool
}

func (pv *PostgresVariables) GetConnectionString() string {
	sslMode := "disable"
	if pv.UseSSL {
		sslMode = "require"
	}

	return fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=%s",
		pv.DbUser, pv.DbPassword, pv.DbName, pv.DbPort, sslMode)
}

func DefaultPostgresConfig() Config {
	return Config{
		MaxOpenConn:     10,
		MaxIdleConn:     5,
		ConnMaxLifeTime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	}
}

func NewPostgresConnection(connString string, cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("couldn't connect to database")
	}

	return db, nil
}

func NewDefaultPostgresClient() *sqlx.DB {
	pgConfig := DefaultPostgresConfig()

	pgVars := PostgresVariables{
		DbUser:     os.Getenv("POSTGRES_USER"),
		DbPassword: os.Getenv("POSTGRES_PASSWORD"),
		DbName:     os.Getenv("POSTGRES_DB"),
		DbPort:     os.Getenv("POSTGRES_PORT"),
		UseSSL:     false,
	}

	connStr := pgVars.GetConnectionString()

	db, err := NewPostgresConnection(connStr, pgConfig)
	if err != nil {
		panic(err)
	}
	return db
}
