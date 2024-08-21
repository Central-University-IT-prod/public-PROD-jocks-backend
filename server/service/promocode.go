package service

import (
	"context"
	"errors"
	"fmt"
	"server/utils"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Promocode struct {
	Id          int    `json:"id"`
	BusinessId  int    `json:"business_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GeneratedPromocode struct {
	Token       string `json:"token"`
	CustomerId  int    `json:"customer_Id"`
	PromocodeId int    `json:"promocode_id"`
	Activated   bool   `json:"activated"`
}

type PromocodeService struct {
	pool *pgxpool.Pool
}

func (ps *PromocodeService) init() error {
	query := `CREATE TABLE IF NOT EXISTS promocodes (
		id SERIAL PRIMARY KEY,
		business_id INTEGER references businesses(id),
		name varchar(255),
		description varchar(255)
	);
	CREATE TABLE IF NOT EXISTS generated_promocodes (
		token TEXT UNIQUE,
		customer_id INTEGER references clients(id),
		promocode_id INTEGER references promocodes(id),
		activated BOOl
	)
	`

	return ps.pool.AcquireFunc(context.TODO(), func(c *pgxpool.Conn) error {
		_, err := c.Exec(context.TODO(), query)
		return err
	})
}

func (ps *PromocodeService) Put(promocode Promocode) error {
	deleteQuery := `DELETE FROM promocodes WHERE business_id=$1 RETURNING id`
	conn, err := ps.pool.Acquire(context.TODO())

	if err != nil {
		return err
	}

	defer conn.Release()
	tx, err := conn.Begin(context.TODO())

	if err != nil {
		return err
	}

	var id int
	err = tx.QueryRow(context.TODO(), deleteQuery, promocode.BusinessId).Scan(&id)

	if errors.Is(err, pgx.ErrNoRows) {
		insertQuery := `INSERT INTO promocodes (business_id, name, description) VALUES ($1, $2, $3)`
		_, err = tx.Exec(context.TODO(), insertQuery, promocode.BusinessId, promocode.Name, promocode.Description)
		if err != nil {
			tx.Rollback(context.TODO())
			return err
		}

		return tx.Commit(context.TODO())
	}

	if err != nil {
		tx.Rollback(context.TODO())
		return err
	}

	insertQuery := `INSERT INTO promocodes (id, business_id, name, description) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(context.TODO(), insertQuery, id, promocode.BusinessId, promocode.Name, promocode.Description)

	if err != nil {
		tx.Rollback(context.TODO())
		return err
	}

	return tx.Commit(context.TODO())

}

func (ps *PromocodeService) CheckTokenExists(token string) bool {
	conn, err := ps.pool.Acquire(context.TODO())
	if err != nil {
		return true
	}

	defer conn.Release()
	query := `SELECT token FROM generated_promocodes WHERE token=$1`
	err = conn.QueryRow(context.TODO(), query, token).Scan(&token)

	return err == nil
}

// returns Promocode and token
func (ps *PromocodeService) GeneratePromocode(businessId, customerId int) (Promocode, string, error) {
	conn, err := ps.pool.Acquire(context.TODO())

	if err != nil {
		return Promocode{}, "", err
	}

	defer conn.Release()
	getPromoCodeQuery := `SELECT * FROM promocodes WHERE business_id=$1`
	var promocode Promocode
	err = conn.QueryRow(context.TODO(), getPromoCodeQuery, businessId).Scan(&promocode.Id, &promocode.BusinessId, &promocode.Name, &promocode.Description)

	if err != nil {
		return Promocode{}, "", err
	}

	token := strings.ToUpper(utils.GenerateToken()[:6])
	for ps.CheckTokenExists(token) {
		token = strings.ToUpper(utils.GenerateToken()[:6])
	}

	insertQuery := `INSERT INTO generated_promocodes (token, customer_id, promocode_id, activated) VALUES ($1, $2, $3, $4)`
	_, err = conn.Exec(context.TODO(), insertQuery, token, customerId, promocode.Id, false)

	if err != nil {
		return Promocode{}, "", err
	}

	return promocode, token, nil

}

type GeneratedPromocodeFull struct {
	Promocode
	GeneratedPromocode
}

func (ps *PromocodeService) GetForUser(customerId int) ([]GeneratedPromocodeFull, error) {
	query := `SELECT * FROM generated_promocodes
	JOIN promocodes ON generated_promocodes.promocode_id=promocodes.id
	WHERE customer_id=$1`
	conn, err := ps.pool.Acquire(context.TODO())

	if err != nil {
		return []GeneratedPromocodeFull{}, err
	}

	defer conn.Release()
	rows, err := conn.Query(context.TODO(), query, customerId)

	if err != nil {
		return []GeneratedPromocodeFull{}, err
	}

	defer rows.Close()
	res := []GeneratedPromocodeFull{}
	for rows.Next() {
		var gp GeneratedPromocodeFull
		err := rows.Scan(
			&gp.GeneratedPromocode.Token,
			&gp.GeneratedPromocode.CustomerId,
			&gp.GeneratedPromocode.PromocodeId,
			&gp.GeneratedPromocode.Activated,
			&gp.Promocode.Id,
			&gp.Promocode.BusinessId,
			&gp.Promocode.Name,
			&gp.Promocode.Description,
		)

		if err != nil {
			fmt.Println(err)
			continue
		}

		res = append(res, gp)
	}

	return res, nil
}

func NewPromocodeService(pool *pgxpool.Pool) (*PromocodeService, error) {
	ps := &PromocodeService{pool: pool}
	err := ps.init()

	if err != nil {
		return nil, err
	}

	return ps, nil
}
