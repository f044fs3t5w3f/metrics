package db

import (
	"context"
	"database/sql"

	"github.com/f044fs3t5w3f/metrics/internal/models"
)

type dbStorage struct {
	db *sql.DB
}

func (d *dbStorage) AddCounter(metricName string, value int64) error {
	ctx := context.Background()
	row := d.db.QueryRowContext(ctx, `
SELECT count(*) 
FROM metric 
WHERE name = $1 AND type = 'counter'`, metricName)
	err := row.Err()
	if err != nil {
		return err
	}
	var count int8
	err = row.Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		_, err := d.db.ExecContext(ctx, `
			UPDATE metric
			SET delta = delta + $1
			WHERE name = $2 and type = 'counter'`,
			value, metricName)
		return err
	} else {
		_, err := d.db.ExecContext(ctx, `
			INSERT INTO metric (name, delta, type)
			VALUES ($1, $2, 'counter')`,
			metricName, value)
		return err
	}
}

func (d *dbStorage) GetCounter(metricName string) (int64, error) {
	ctx := context.Background()
	row := d.db.QueryRowContext(ctx, `
SELECT delta 
FROM metric
WHERE name = $1 and type = 'counter'`, metricName)
	err := row.Err()
	if err != nil {
		return 0, err
	}
	var value int64
	err = row.Scan(&value)
	return value, err
}

func (d *dbStorage) GetGauge(metricName string) (float64, error) {
	ctx := context.Background()
	row := d.db.QueryRowContext(ctx, `
SELECT value 
FROM metric
WHERE name = $1 and type = 'gauge'`, metricName)
	err := row.Err()
	if err != nil {
		return 0, err
	}
	var value float64
	err = row.Scan(&value)
	return value, err
}

func (d *dbStorage) GetValuesList() ([]models.Metrics, error) {
	ctx := context.Background()
	rows, err := d.db.QueryContext(ctx, "SELECT name, type, delta, value FROM metric")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	metrics := make([]models.Metrics, 0)

	for rows.Next() {
		metric := models.Metrics{}
		var (
			delta *int64
			value *float64
		)
		err := rows.Scan(&metric.ID, &metric.MType, &delta, &value)
		if err != nil {
			continue
		}
		switch metric.MType {
		case models.Gauge:
			metric.Value = value
		case models.Counter:
			metric.Delta = delta
		default:
			continue
		}
		metrics = append(metrics, metric)
	}
	return metrics, nil
}

func (d *dbStorage) SetGauge(metricName string, value float64) error {
	ctx := context.Background()
	row := d.db.QueryRowContext(ctx, `
SELECT count(*) 
FROM metric 
WHERE name = $1 AND type = 'gauge'`, metricName)
	err := row.Err()
	if err != nil {
		return err
	}
	var count int8
	err = row.Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		_, err := d.db.ExecContext(ctx, `
			UPDATE metric
			SET value = $1
			WHERE name = $2 and type = 'gauge'`,
			value, metricName)
		return err
	} else {
		_, err := d.db.ExecContext(ctx, `
			INSERT INTO metric (name, value, type)
			VALUES ($1, $2, 'gauge')`,
			metricName, value)
		return err
	}
}

func NewDBStorage(db *sql.DB) *dbStorage {
	return &dbStorage{
		db: db,
	}
}
