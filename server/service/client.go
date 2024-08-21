package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	Id             int    `json:"id"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"-"`
}

type ClientService struct {
	pool *pgxpool.Pool
}

func (cs *ClientService) init() error {
	query := `CREATE TABLE IF NOT EXISTS clients (
		id SERIAL PRIMARY KEY,
		username varchar(255),
		email varchar(255) UNIQUE,
		hashed_password TEXT
	)`

	return cs.pool.AcquireFunc(context.TODO(), func(c *pgxpool.Conn) error {
		_, err := c.Exec(context.TODO(), query)
		return err
	})
}

func (cs *ClientService) Create(client Client) (int, error) {
	query := `INSERT INTO clients (username, email, hashed_password) VALUES ($1, $2, $3) RETURNING id`
	var id int

	err := cs.pool.AcquireFunc(context.TODO(), func(c *pgxpool.Conn) error {
		return c.QueryRow(context.TODO(), query, client.Username, client.Email, client.HashedPassword).Scan(&id)
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (cs *ClientService) GetByEmail(email string) (Client, error) {
	query := `SELECT * FROM clients WHERE email=$1`

	var client Client
	err := cs.pool.AcquireFunc(context.TODO(), func(c *pgxpool.Conn) error {
		return c.QueryRow(context.TODO(), query, email).Scan(&client.Id, &client.Username, &client.Email, &client.HashedPassword)
	})

	if err != nil {
		return Client{}, err
	}

	return client, nil
}

func (cs *ClientService) GetReviews(customer_id int) ([]Reviews,error) {
	var reviews []Review

	query := `SELECT * FROM reviews WHERE customer_id =$1`
	conn, err := cs.pool.Acquire(context.TODO())
	if err != nil {
		return []Reviews{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.TODO(), query, customer_id)
	if err != nil {
		return []Reviews{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var review Review
		err := rows.Scan(&review.Id, &review.CustomerId, &review.BusinessId, &review.Created_at,&review.CheckId)
		if err != nil {
			fmt.Println(err)
			continue
		}

		reviews = append(reviews, review)
	}

	var reviewsGeneral []Reviews

	for _, review := range reviews {

		producerQuery := `SELECT value FROM reviews_producer WHERE review_id = $1`
		var producerValue float64
		if err := conn.QueryRow(context.TODO(), producerQuery, review.Id).Scan(&producerValue); err != nil {
			return []Reviews{},err
		}

		productQuery := `SELECT value FROM reviews_product WHERE review_id = $1`
		var productValue float64
		if err := conn.QueryRow(context.TODO(), productQuery, review.Id).Scan(&productValue); err != nil {
			return []Reviews{}, err
		}

		average := (producerValue + productValue) / 2

		reviewObject := Reviews{
			CreatedAt: review.Created_at,
			Rating:    average,
		}
		reviewsGeneral = append(reviewsGeneral, reviewObject)
	}

	return reviewsGeneral,  nil

}


func NewClientService(pool *pgxpool.Pool) (*ClientService, error) {
	cs := &ClientService{pool: pool}
	err := cs.init()

	if err != nil {
		return nil, err
	}

	return cs, nil
}
