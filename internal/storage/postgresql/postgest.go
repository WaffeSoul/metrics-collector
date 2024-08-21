package postgresql

import (
	"context"
	"errors"
	"strconv"

	"github.com/WaffeSoul/metrics-collector/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(addressDB string) *Repository {
	return &Repository{db: InitDB("addressDB")}
}

func InitDB(addr string) *pgxpool.Pool {
	poolConfig, err := pgxpool.ParseConfig(addr)
	if err != nil {
		return nil
		// log.Fatalln("Unable to parse DATABASE_URL:", err)
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil
	}
	// defer conn.Close()
	err = migrateTables(conn)
	if err != nil {
		return nil
	}
	return conn
}

func migrateTables(pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `
	CREATE TABLE IF NOT EXISTS gauges (
		name VARCHAR(255) PRIMARY KEY,
		value DOUBLE PRECISION
	);
	CREATE TABLE IF NOT EXISTS counters (
		name VARCHAR(255) PRIMARY KEY,
		value INTEGER
	);`)
	return err
}

func (p *Repository) Delete(typeMetric string, key string) error {
	if typeMetric == "gauge" {

	} else if typeMetric == "counter" {

	} else {
		return errors.New("NotFound")
	}
	return nil
}

func (p *Repository) Add(typeMetric string, key string, value string) error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	switch typeMetric {
	case "gauge":
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		if _, err := conn.Exec(context.Background(), `insert into gauges(name, key) values ($1, $2)
		on conflict (name) do update set value=value + $2`, key, value); err == nil {
			return nil
		}
	case "counter":
		value, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		if _, err := conn.Exec(context.Background(), `insert into counters(name, key) values ($1, $2)
		on conflict (name) do update set value=$2`, key, value); err == nil {
			return nil
		}
	default:
		return errors.New("NotFound")
	}
	return errors.New("NotFound")
}

func (p *Repository) AddJSON(data model.Metrics) error{
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	switch data.MType {
	case "gauge":
		if _, err := conn.Exec(context.Background(), `insert into gauges(name, key) values ($1, $2)
		on conflict (name) do update set value=$2`, data.ID, data.Value); err == nil {
			return nil
		}
	case "counter":
		if _, err := conn.Exec(context.Background(), `insert into counters(name, key) values ($1, $2)
		on conflict (name) do update set value=value + $2`, data.ID, data.Delta); err == nil {
			return nil
		}
	default:
		return errors.New("NotFound")
	}
	return errors.New("NotFound")
}

func (p *Repository) GetJSON(data model.Metrics) (model.Metrics, error) {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return data, err
	}
	defer conn.Release()
	switch data.MType {
	case "gauge":
	err := conn.QueryRow(context.Background(), "select * from gauges where name=$1", data.ID).Scan(data.Value)
		if err != nil {
			return data, errors.New("NotFound")
		}
		return data, nil
	case "counter":
		err := conn.QueryRow(context.Background(), "select * from counters where name=$1", data.ID).Scan(data.Delta)
		if err != nil {
			return data, errors.New("NotFound")
		}
		return data, nil
	default:
		return data, errors.New("NotFound")
	}
}

func (p *Repository) Get(typeMetric string, key string) (interface{}, error) {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	switch typeMetric {
	case "gauge":
		var value float64
		err := conn.QueryRow(context.Background(), "select * from gauges where name=$1", key).Scan(&value)
		if err != nil {
			return nil, errors.New("NotFound")
		}
		return value, nil
	case "counter":
		var value int
		err := conn.QueryRow(context.Background(), "select * from counters where name=$1", key).Scan(&value)
		if err != nil {
			return nil, errors.New("NotFound")
		}
		return value, nil
	default:
		return nil, errors.New("NotFound")
	}
}

func (p *Repository) GetAll() []byte {
	return nil
}

func (p *Repository) AutoSaveStorage() {

}

func (p *Repository) Ping() error {
	if p.db == nil {
		return errors.New("db is nil")
	}
	return p.db.Ping(context.Background())
}