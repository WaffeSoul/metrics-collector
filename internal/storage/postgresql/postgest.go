package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/WaffeSoul/metrics-collector/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(addressDB string) *Repository {
	return &Repository{db: InitDB(addressDB)}
}

func InitDB(addr string) *pgxpool.Pool {
	fmt.Println("Alee")
	poolConfig, err := pgxpool.ParseConfig(addr)
	if err != nil {
		fmt.Println(err)
		return nil
		// log.Fatalln("Unable to parse DATABASE_URL:", err)
	}
	fmt.Println("Alee")
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil
	}
	fmt.Println("Alee")
	// defer conn.Close()
	err = migrateTables(conn)
	if err != nil {
		fmt.Println(err)
		return conn
	}
	return conn
}

func migrateTables(pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	fmt.Println("Alee")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer conn.Release()
	fmt.Println("Alee")
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS gauges (
		name VARCHAR(255) PRIMARY KEY,
		value DOUBLE PRECISION
	);
	CREATE TABLE IF NOT EXISTS counters (
		name VARCHAR(255) PRIMARY KEY,
		value bigint
	);`)
	fmt.Println("Alee")
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
		if _, err := conn.Exec(context.Background(), `insert into gauges(name, value) values ($1, $2)
		on conflict (name) do update set value=value + $2`, key, value); err == nil {
			return nil
		}
	case "counter":
		value, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		if _, err := conn.Exec(context.Background(), `insert into counters(name, value) values ($1, $2)
		on conflict (name) do update set value=$2`, key, value); err == nil {
			return nil
		}
	default:
		return errors.New("NotFound")
	}
	return errors.New("NotFound")
}

func (p *Repository) AddJSON(data model.Metrics) error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()
	switch data.MType {
	case "gauge":

		_, err := conn.Exec(context.Background(), `insert into gauges(name, value) values ($1, $2)
		on conflict (name) do update set value=$2`, data.ID, data.Value)
		if err == nil {
			return nil
		}
		fmt.Println(err)
	case "counter":
		fmt.Println(*data.Delta)
		_, err := conn.Exec(context.Background(), `insert into counters(name, value) values ($1, $2)
		on conflict (name) do update set value = counters.value + $2`, data.ID, data.Delta)
		if err == nil {
			return nil
		}
		fmt.Println(err)
	default:
		return errors.New("NotFound")
	}
	return errors.New("NotFound")
}

func (p *Repository) AddMuiltJSON(data []model.Metrics) error {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return err
	}
	batch := &pgx.Batch{}
	defer conn.Release()
	for _, i := range data {
		switch i.MType {
		case "gauge":
			batch.Queue(`insert into gauges(name, value) values ($1, $2)
			on conflict (name) do update set value=$2`, i.ID, i.Value)
		case "counter":
			fmt.Println(*i.Delta)
			batch.Queue(`insert into counters(name, value) values ($1, $2)
			on conflict (name) do update set value = counters.value + $2`, i.ID, i.Delta)
		}
	}
	br := conn.SendBatch(context.Background(), batch)
	err = br.Close()
	if err != nil {
		return err
	}
	return nil
}

func (p *Repository) GetJSON(data model.Metrics) (model.Metrics, error) {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return data, err
	}
	defer conn.Release()
	switch data.MType {
	case "gauge":
		err := conn.QueryRow(context.Background(), "select value from gauges where name=$1", data.ID).Scan(&data.Value)
		if err != nil {
			return data, errors.New("NotFound")
		}
		return data, nil
	case "counter":
		err := conn.QueryRow(context.Background(), "select value from counters where name=$1", data.ID).Scan(&data.Delta)
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
		err := conn.QueryRow(context.Background(), "select value from gauges where name=$1", key).Scan(&value)
		if err != nil {
			return nil, errors.New("NotFound")
		}
		return value, nil
	case "counter":
		var value int
		err := conn.QueryRow(context.Background(), "select value from counters where name=$1", key).Scan(&value)
		if err != nil {
			return nil, errors.New("NotFound")
		}
		return value, nil
	default:
		return nil, errors.New("NotFound")
	}
}

func (p *Repository) GetAll() []byte {
	conn, err := p.db.Acquire(context.Background())
	if err != nil {
		return nil
	}
	defer conn.Release()
	resultData := []byte{}
	rows, err := conn.Query(context.Background(), "select * from counters")
	if err != nil {
		fmt.Println(1)
		fmt.Println(err)
		return nil
	}
	for rows.Next() {
		var name string
		var value int32
		err := rows.Scan(&name, &value)
		if err != nil {
			fmt.Println(2)
			fmt.Println(err)
			return nil
		}
		resultData = append(resultData, []byte(fmt.Sprintf("%v: %v\n", name, value))...)
	}
	rows, err = conn.Query(context.Background(), "select * from gauges")
	if err != nil {
		fmt.Println(3)
		fmt.Println(err)
		return nil
	}
	for rows.Next() {
		var name string
		var value float64
		err := rows.Scan(&name, &value)
		if err != nil {
			fmt.Println(4)
			fmt.Println(err)
			return nil
		}
		resultData = append(resultData, []byte(fmt.Sprintf("%v: %v\n", name, value))...)
	}
	return resultData
}

func (p *Repository) AutoSaveStorage() {

}

func (p *Repository) Ping() error {
	if p.db == nil {
		return errors.New("db is nil")
	}
	return p.db.Ping(context.Background())
}
