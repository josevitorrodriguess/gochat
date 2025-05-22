package postgres

import (
	"context"
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

type PostgresStorage struct {
	db *sqlx.DB
}

type postgresVariables struct {
	dbName     string
	dbUser     string
	dbPassword string
	dbPort     string
	useSSL     bool
}

func (pv *postgresVariables) getConnectionString() string {
	sslMode := "disable"
	if pv.useSSL {
		sslMode = "require"
	}

	return fmt.Sprintf("user=%s password=%s dbname=%s port=%s sslmode=%s",
		pv.dbUser, pv.dbPassword, pv.dbName, pv.dbPort, sslMode)
}

func DefaultConfig() Config {
	return Config{
		MaxOpenConn:     10,
		MaxIdleConn:     5,
		ConnMaxLifeTime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
	}
}

func newConnection(connString string, cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", connString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("couldn't connect to database")
	}

	db.SetMaxOpenConns(int(cfg.MaxOpenConn))
	db.SetMaxIdleConns(int(cfg.MaxIdleConn))
	db.SetConnMaxLifetime(cfg.ConnMaxLifeTime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	return db, nil
}

func New(connString string, cfg Config) (*PostgresStorage, error) {
	db, err := newConnection(connString, cfg)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func NewDefault() (*PostgresStorage, error) {
	cfg := DefaultConfig()

	vars := postgresVariables{
		dbUser:     os.Getenv("POSTGRES_USER"),
		dbPassword: os.Getenv("POSTGRES_PASSWORD"),
		dbName:     os.Getenv("POSTGRES_DB"),
		dbPort:     os.Getenv("POSTGRES_PORT"),
		useSSL:     false,
	}

	connStr := vars.getConnectionString()
	return New(connStr, cfg)
}

func (p *PostgresStorage) GetClient() interface{} {
	return p.db
}

func (p *PostgresStorage) GetDB() *sqlx.DB {
	return p.db
}

func (p *PostgresStorage) Ping(ctx context.Context) error {
	return p.db.PingContext(ctx)
}

func (p *PostgresStorage) Close() error {
	return p.db.Close()
}
