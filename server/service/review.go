package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewService struct {
	pool *pgxpool.Pool
}

type Reviews struct {
	CreatedAt time.Time `json:"created_at"`
	Rating    float64   `json:"rating"`
}
type Review struct {
	Id         int
	CustomerId int
	BusinessId int
	Created_at time.Time
	CheckId int
}
type Count struct{
	User int
	Anon int
}
func (rs *ReviewService) init() error {
	query := `
        CREATE TABLE IF NOT EXISTS reviews_producer (
            id SERIAL PRIMARY KEY,
            form_id INTEGER REFERENCES form_parameters(id),
            value INTEGER,
            review_id INTEGER DEFAULT NULL
        );
		CREATE TABLE IF NOT EXISTS reviews (
            id SERIAL PRIMARY KEY,
            customer_id INTEGER DEFAULT NULL,
			business_id INTEGER REFERENCES businesses(id),
			created_at TIMESTAMPTZ,
			check_id INTEGER UNIQUE REFERENCES checks
        );
		CREATE TABLE IF NOT EXISTS reviews_product (
            id SERIAL PRIMARY KEY,
            product_id INTEGER REFERENCES products(id),
            value INTEGER,
            review_id INTEGER REFERENCES reviews(id)
        );
    `
	return rs.pool.AcquireFunc(context.TODO(), func(c *pgxpool.Conn) error {
		_, err := c.Exec(context.TODO(), query)
		return err
	})
}

func (rs *ReviewService) Review(customer_id int, business_id int, check_id int) (int, error) {
	query := `INSERT INTO reviews (customer_id,business_id,created_at,check_id) VALUES ($1, $2,$3,$4) RETURNING id`

	conn, err := rs.pool.Acquire(context.TODO())
	if err != nil {
		return 0, err
	}
	defer conn.Release()
	var id int
	if err := conn.QueryRow(context.TODO(), query, customer_id, business_id, time.Now(), check_id).Scan(&id); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return 0, fmt.Errorf("извините, на данный чек уже написан отзыв")
		}
		return 0, err
	}
	return id, nil
}

func (rs *ReviewService) ReviewProducer(form_id int, value int, review_id int) (int, error) {
	query := `INSERT INTO reviews_producer (form_id, value, review_id) VALUES ($1, $2, $3) RETURNING id`

	conn, err := rs.pool.Acquire(context.TODO())
	if err != nil {
		return 0, err
	}
	defer conn.Release()

	var id int
	err = conn.QueryRow(context.TODO(), query, form_id, value, review_id).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (rs *ReviewService) ReviewProduct(product_id int, value int, review_id int) (int, error) {
	query := `INSERT INTO reviews_product (product_id,value,review_id) VALUES ($1,$2,$3) RETURNING id`
	conn, err := rs.pool.Acquire(context.TODO())
	if err != nil {
		return 0, err
	}
	defer conn.Release()
	var id int
	if err := conn.QueryRow(context.TODO(), query, product_id, value, review_id).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (rs *ReviewService) GetRatingForBusiness(businessId int) (float64, int, error) {
	query := `SELECT avg(value), count(value) FROM reviews_producer WHERE form_id in (SELECT id FROM form_parameters WHERE business_id=$1)`
	var (
		ratingAny any
		countAny  any
	)
	err := rs.pool.AcquireFunc(context.TODO(), func(c *pgxpool.Conn) error {
		return c.QueryRow(context.TODO(), query, businessId).Scan(&ratingAny, &countAny)
	})

	if err != nil {
		return 0, 0, err
	}

	ratingPgNumeric, ok := ratingAny.(pgtype.Numeric)
	if !ok {
		return 0, 0, nil
	}

	count, ok := countAny.(int64)
	if !ok {
		return 0, 0, nil
	}

	raitingFloat64, _ := ratingPgNumeric.Float64Value()

	return raitingFloat64.Float64, int(count), nil
}

type ProductsWithRating struct {
	Product
	Rating      float64 `json:"rating"`
	RatingCount int     `json:"rating_count"`
}

func (rs *ReviewService) GetRatingForProductsOf(businessId int) ([]ProductsWithRating, error) {
	query := `SELECT products.id, business_id, name, avg(value), count(value) FROM products LEFT JOIN reviews_product ON products.id=reviews_product.product_id
	WHERE business_id=$1
	GROUP BY (products.id)`
	conn, err := rs.pool.Acquire(context.TODO())

	if err != nil {
		return []ProductsWithRating{}, err
	}

	defer conn.Release()
	rows, err := conn.Query(context.TODO(), query, businessId)

	if err != nil {
		return []ProductsWithRating{}, err
	}

	defer rows.Close()
	res := []ProductsWithRating{}
	for rows.Next() {
		var (
			pwr                 ProductsWithRating
			ratingAny, countAny any
		)
		err := rows.Scan(&pwr.Product.Id, &pwr.Product.BusinessId, &pwr.Product.Name, &ratingAny, &countAny)

		if err != nil {
			fmt.Println(err)
			continue
		}

		ratingPgNumeric, ok := ratingAny.(pgtype.Numeric)

		if !ok {
			pwr.Rating = 0
			pwr.RatingCount = 0
			res = append(res, pwr)
			continue
		}

		count, ok := countAny.(int64)

		if !ok {
			pwr.Rating = 0
			pwr.RatingCount = 0
			res = append(res, pwr)
			continue
		}

		raitingFloat64, _ := ratingPgNumeric.Float64Value()
		pwr.Rating = raitingFloat64.Float64
		pwr.RatingCount = int(count)
		res = append(res, pwr)
	}

	return res, nil
}

type ParamWithRating struct {
	Parameter
	Rating      float64 `json:"rating"`
	RatingCount int     `json:"rating_count"`
}

func (rs *ReviewService) GetBusinessFormsRating(businessId int) ([]ParamWithRating, error) {
	query := `SELECT form_parameters.id, business_id, name, type, avg(value), count(value) FROM form_parameters LEFT JOIN reviews_producer ON form_parameters.id=reviews_producer.form_id
	WHERE business_id=$1
	GROUP BY (form_parameters.id, name)`
	fmt.Println(query, businessId)
	conn, err := rs.pool.Acquire(context.TODO())

	if err != nil {
		return []ParamWithRating{}, err
	}

	defer conn.Release()
	rows, err := conn.Query(context.TODO(), query, businessId)

	if err != nil {
		return []ParamWithRating{}, err
	}

	defer rows.Close()
	res := []ParamWithRating{}
	for rows.Next() {
		var (
			pwr                 ParamWithRating
			ratingAny, countAny any
		)
		err := rows.Scan(&pwr.Parameter.Id, &pwr.Parameter.BusinessId, &pwr.Parameter.Name, &pwr.Parameter.Type, &ratingAny, &countAny)

		if err != nil {
			fmt.Println(err)
			continue
		}

		ratingPgNumeric, ok := ratingAny.(pgtype.Numeric)

		if !ok {
			pwr.Rating = 0
			pwr.RatingCount = 0
			res = append(res, pwr)
			continue
		}

		count, ok := countAny.(int64)

		if !ok {
			pwr.Rating = 0
			pwr.RatingCount = 0
			res = append(res, pwr)
			continue
		}

		raitingFloat64, _ := ratingPgNumeric.Float64Value()
		pwr.Rating = raitingFloat64.Float64
		pwr.RatingCount = int(count)
		res = append(res, pwr)
	}

	return res, nil
}

func (rs *ReviewService) GetReviews(buisness_id int) ([]Reviews,Count,error) {
	var reviews []Review

	query := `SELECT * FROM reviews WHERE business_id =$1`
	conn, err := rs.pool.Acquire(context.TODO())
	if err != nil {
		return []Reviews{},Count{}, err
	}
	defer conn.Release()

	rows, err := conn.Query(context.TODO(), query, buisness_id)
	if err != nil {
		return []Reviews{},Count{}, err
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
	var anonCount int
	var userCount int
	for _, review := range reviews {
		if review.CustomerId == 0{
			anonCount ++
		}else{
			userCount ++
		}

		producerQuery := `SELECT value FROM reviews_producer WHERE review_id = $1`
		var producerValue float64
		if err := conn.QueryRow(context.TODO(), producerQuery, review.Id).Scan(&producerValue); err != nil {
			return []Reviews{},Count{},err
		}

		productQuery := `SELECT value FROM reviews_product WHERE review_id = $1`
		var productValue float64
		if err := conn.QueryRow(context.TODO(), productQuery, review.Id).Scan(&productValue); err != nil {
			return []Reviews{},Count{}, err
		}

		average := (producerValue + productValue) / 2

		reviewObject := Reviews{
			CreatedAt: review.Created_at,
			Rating:    average,
		}
		reviewsGeneral = append(reviewsGeneral, reviewObject)
	}
	counts := Count{User: userCount,Anon: anonCount}
	return reviewsGeneral, counts, nil

}

func NewReviewService(pool *pgxpool.Pool) (*ReviewService, error) {
	rs := ReviewService{pool: pool}
	err := rs.init()
	if err != nil {
		return nil, err
	}
	return &rs, nil
}
