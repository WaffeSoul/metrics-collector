package storage

import (
	"errors"

	"github.com/WaffeSoul/metrics-collector/internal/model"
	"github.com/WaffeSoul/metrics-collector/internal/storage/mem"
	"github.com/WaffeSoul/metrics-collector/internal/storage/postgresql"
)

type Database struct {
	DB Store
}

func New(typeDB string, interlval int, path string, addrDB string) (*Database, error) {
	if typeDB == "mem" {
		return &Database{
			DB: mem.InitMem(interlval, path),
		}, nil
	} else if typeDB == "postgresql" {
		return &Database{
			DB: postgresql.NewRepository(addrDB),
		}, nil
	}
	return nil, errors.New("type metric error")
}

type Store interface {
	Delete(typeMetric string, key string) error
	Add(typeMetric string, key string, value string) error
	AddJSON(data model.Metrics) error
	Get(typeMetric string, key string) (interface{}, error)
	GetJSON(data model.Metrics) (model.Metrics, error)
	GetAll() []byte
	AutoSaveStorage()
	Ping() error
}
