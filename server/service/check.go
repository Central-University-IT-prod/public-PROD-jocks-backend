package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Check struct {
	Id         int   `json:"id"`
	Products   []int `json:"products"`
	BusinessId int   `json:"business_id"`
}

type CheckService struct {
	pool *pgxpool.Pool
}

func (cs *CheckService) init() error {
	query := `CREATE TABLE IF NOT EXISTS checks (
        id SERIAL PRIMARY KEY,
		business_id BIGINT REFERENCES businesses(id),
        products INTEGER[]
    )`

	return cs.pool.AcquireFunc(context.TODO(), func(c *pgxpool.Conn) error {
		_, err := c.Exec(context.TODO(), query)
		return err
	})
}

func (cs *CheckService) Create(check *Check) (int, error) {
	var id int
	query := "INSERT INTO checks (business_id,products) VALUES ($1,$2) RETURNING id"

	conn, err := cs.pool.Acquire(context.TODO())
	if err != nil {
		return 0, err
	}
	defer conn.Release()
	err = conn.QueryRow(context.TODO(), query, check.BusinessId, check.Products).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (cs *CheckService) GetProductsById(checkId int) ([]Product, error) {
	check_query := `SELECT products FROM checks WHERE id = $1`

	conn, err := cs.pool.Acquire(context.TODO())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	var productsId []int
	err = conn.QueryRow(context.TODO(), check_query, checkId).Scan(&productsId)
	if err != nil {
		return nil, err
	}
	var products []Product
	for _, productId := range productsId {
		productQuery := `SELECT id, name, business_id FROM products WHERE id = $1`
		var product Product
		err := conn.QueryRow(context.TODO(), productQuery, productId).Scan(&product.Id, &product.Name, &product.BusinessId)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (cs *CheckService) GetById(id int) (Check, error) {
	query := `SELECT * FROM checks WHERE id=$1`
	conn, err := cs.pool.Acquire(context.TODO())

	if err != nil {
		return Check{}, err
	}

	defer conn.Release()
	var check Check
	row := conn.QueryRow(context.TODO(), query, id)
	err = row.Scan(&check.Id, &check.BusinessId, &check.Products)

	if err != nil {
		return Check{}, err
	}

	return check, nil
}

func NewCheckService(pool *pgxpool.Pool) (*CheckService, error) {
	cs := &CheckService{pool: pool}
	err := cs.init()

	if err != nil {
		return nil, err
	}

	return cs, nil
}
