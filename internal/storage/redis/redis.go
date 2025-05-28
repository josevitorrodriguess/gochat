package redis

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
)

type Config struct {
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
	Wait        bool
}

type RedisStorage struct {
	pool *redis.Pool
}

type redisVariables struct {
	host     string
	port     string
	password string
	db       int
}

func DefaultConfig() Config {
	return Config{
		MaxIdle:     10,
		MaxActive:   50,
		IdleTimeout: 5 * time.Minute,
		Wait:        true,
	}
}

func New(host, port, password string, db int, cfg Config) (*RedisStorage, error) {
	pool := &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: cfg.IdleTimeout,
		Wait:        cfg.Wait,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
			if err != nil {
				return nil, err
			}

			// Autenticar se senha foi fornecida
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}

			// Selecionar database
			if db != 0 {
				if _, err := c.Do("SELECT", db); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	// Testar conexão
	conn := pool.Get()
	defer conn.Close()

	if _, err := conn.Do("PING"); err != nil {
		return nil, fmt.Errorf("couldn't connect to Redis: %w", err)
	}

	return &RedisStorage{pool: pool}, nil
}

func NewDefault() (*RedisStorage, error) {
	cfg := DefaultConfig()

	dbStr := os.Getenv("REDIS_DB")
	db := 0
	if dbStr != "" {
		var err error
		db, err = strconv.Atoi(dbStr)
		if err != nil {
			return nil, fmt.Errorf("invalid REDIS_DB value: %w", err)
		}
	}

	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	password := os.Getenv("REDIS_PASSWORD")

	return New(host, port, password, db, cfg)
}

func (r *RedisStorage) GetPool() *redis.Pool {
	return r.pool
}

func (r *RedisStorage) Ping() error {
	conn := r.pool.Get()
	defer conn.Close()
	
	_, err := conn.Do("PING")
	return err
}

func (r *RedisStorage) Close() error {
	return r.pool.Close()
}