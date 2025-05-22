package storage

import "context"

type Storage interface {
	GetClient() interface{}
	Ping(ctx context.Context) error
	Close() error
}
