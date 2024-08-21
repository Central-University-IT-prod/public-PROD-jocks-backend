package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Parameter struct {
	Id         int    `json:"id"`
	BusinessId int    `json:"-"`
	Name       string `json:"name"`
	Type       string `json:"type"`
}

type FormService struct {
	pool *pgxpool.Pool
}

func (fs *FormService) init() error {
	query := `CREATE TABLE IF NOT EXISTS form_parameters (
		id SERIAL PRIMARY KEY,
		business_id INTEGER references businesses(id),
		name varchar(255),
		type varchar(255)
	)`

	return fs.pool.AcquireFunc(context.TODO(), func(c *pgxpool.Conn) error {
		_, err := c.Exec(context.TODO(), query)
		return err
	})
}

func (fs *FormService) PutParams(businessId int, parameters []Parameter) error {
	conn, err := fs.pool.Acquire(context.TODO())
	defer conn.Release()

	if err != nil {
		return err
	}

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		return err
	}

	deleteQuery := `DELETE FROM form_parameters WHERE business_id=$1`
	_, err = tx.Exec(context.TODO(), deleteQuery, businessId)
	if err != nil {
		tx.Rollback(context.TODO())
		return err
	}

	insertQuery := `INSERT INTO form_parameters (business_id, name, type) VALUES ($1, $2, $3)`
	for _, param := range parameters {
		_, err := tx.Exec(context.TODO(), insertQuery, businessId, param.Name, param.Type)
		if err != nil {
			tx.Rollback(context.TODO())
			return err
		}
	}

	return tx.Commit(context.TODO())
}

func (fs *FormService) GetParams(businessId int) ([]Parameter, error) {
	res := []Parameter{}
	query := `SELECT * FROM form_parameters WHERE business_id=$1`
	conn, err := fs.pool.Acquire(context.TODO())
	if err != nil {
		return res, err
	}

	defer conn.Release()
	rows, err := conn.Query(context.TODO(), query, businessId)
	if err != nil {
		return res, err
	}

	defer rows.Close()
	for rows.Next() {
		var param Parameter
		err := rows.Scan(&param.Id, &param.BusinessId, &param.Name, &param.Type)
		if err != nil {
			fmt.Println(err)
			continue
		}
		res = append(res, param)
	}

	return res, err
}

func NewFormService(pool *pgxpool.Pool) (*FormService, error) {
	fs := &FormService{pool: pool}
	err := fs.init()

	if err != nil {
		return nil, err
	}

	return fs, nil
}
