package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// {
// 	"name": "My Busisnesss",
// 	"description": "This is as test businesss",
// 	"addresses": "chita",
// 	"userContactEmail": "contaacsst@example.com",
// 	"email": "busisnsesss@example.com",
// 	"password": "mysescretpassword"
//   }

type Product struct {
	Id         int    `json:"id"`
	BusinessId int    `json:"business_id"`
	Name       string `json:"name"`
}

type ProductService struct {
	pool *pgxpool.Pool
}

func (ps *ProductService) init() error {
	query := `CREATE TABLE IF NOT EXISTS products (
		id SERIAL PRIMARY KEY,
		business_id INTEGER references businesses(id),
		name varchar(255)
	)`

	return ps.pool.AcquireFunc(context.TODO(), func(c *pgxpool.Conn) error {
		_, err := c.Exec(context.TODO(), query)
		return err
	})
}

func (ps *ProductService) Create(p *Product) (int, error) {
	var id int
	query := "INSERT INTO products (business_id, name) VALUES ($1, $2) RETURNING id"

	conn, err := ps.pool.Acquire(context.TODO())
	if err != nil {
		return 0, err
	}
	defer conn.Release()
	err = conn.QueryRow(context.TODO(), query, p.BusinessId, p.Name).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
func (ps *ProductService) Remove(id int) (bool, error) {
	query := "DELETE FROM products WHERE id = $1"
	conn, err := ps.pool.Acquire(context.TODO())
	if err != nil {
		return false, err
	}
	defer conn.Release()
	c, err := conn.Exec(context.TODO(), query, id)
	if err != nil {
		return false, err
	}

	return c.RowsAffected() != 0, nil
}

func (ps *ProductService) GetForBusiness(businessId int) ([]Product, error) {
	query := `SELECT * FROM products WHERE business_id=$1`
	conn, err := ps.pool.Acquire(context.TODO())
	if err != nil {
		return []Product{}, nil
	}

	defer conn.Release()
	rows, err := conn.Query(context.TODO(), query, businessId)

	if err != nil {
		return []Product{}, err
	}

	res := []Product{}
	for rows.Next() {
		var product Product
		err := rows.Scan(&product.Id, &product.BusinessId, &product.Name)

		if err != nil {
			fmt.Println(err)
			continue
		}

		res = append(res, product)
	}

	return res, nil
}

func NewProductService(pool *pgxpool.Pool) (*ProductService, error) {
	ps := &ProductService{pool: pool}
	err := ps.init()

	if err != nil {
		return nil, err
	}

	return ps, err
}
