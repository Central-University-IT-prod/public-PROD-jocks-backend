package service

import (
	"context"
	"fmt"
	"server/utils"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Business struct {
	Id               int    `json:"id"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	UserContactEmail string `json:"userContactEmail"`
	Email            string `json:"email"`
	HashedPassword   string `json:"-"`
}

type Address struct {
	BusinessId int          `json:"-"`
	Address    string       `json:"address"`
	Coords     utils.Coords `json:"coords"`
}

type BusinessWithAddress struct {
	Business
	Address
}

type BusinessService struct {
	pool *pgxpool.Pool
}

func (bs *BusinessService) init() error {
	query := `CREATE TABLE IF NOT EXISTS businesses (
		id SERIAL PRIMARY KEY,
		name varchar(255),
		description varchar(255),
		userContactEmail varchar(255),
		email varchar(255) UNIQUE,
		hashedPassword TEXT
	);

	CREATE TABLE IF NOT EXISTS addresses (
		businessId INTEGER references businesses(id),
		address varchar(255),
		lat double precision,
		lng double precision
	);
	
	CREATE OR REPLACE FUNCTION haversine(lat1 FLOAT, lon1 FLOAT, lat2 FLOAT, lon2 FLOAT)
RETURNS FLOAT AS $$
DECLARE
    dLat FLOAT;
    dLon FLOAT;
    a FLOAT;
    c FLOAT;
    rad FLOAT := 6371;
BEGIN
    dLat := radians(lat2 - lat1);
    dLon := radians(lon2 - lon1);
    a := sin(dLat / 2) * sin(dLat / 2) +
         cos(radians(lat1)) * cos(radians(lat2)) *
         sin(dLon / 2) * sin(dLon / 2);
    c := 2 * asin(sqrt(a));
    RETURN rad * c;
END;
$$ LANGUAGE plpgsql;`

	return bs.pool.AcquireFunc(context.TODO(), func(c *pgxpool.Conn) error {
		_, err := c.Exec(context.TODO(), query)
		return err
	})
}

func (bs *BusinessService) CreateAddress(address Address) error {
	query := `INSERT INTO addresses (businessId, address, lat, lng) VALUES ($1, $2, $3, $4)`

	return bs.pool.AcquireFunc(context.TODO(), func(c *pgxpool.Conn) error {
		_, err := c.Exec(context.TODO(), query, address.BusinessId, address.Address, address.Coords.Lat, address.Coords.Lng)
		return err
	})
}

func (bs *BusinessService) Create(business Business, addresses []string) (int, error) {

	var (
		id  int
		err error
	)

	addBusinessQuery := `INSERT INTO businesses (name, description, userContactEmail, email, hashedPassword) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	addAddressQuery := `INSERT INTO addresses (businessId, address, lat, lng) VALUES ($1, $2, $3, $4)`

	conn, err := bs.pool.Acquire(context.TODO())
	if err != nil {
		return 0, err
	}
	defer conn.Release()
	err = conn.QueryRow(context.TODO(), addBusinessQuery, business.Name, business.Description, business.UserContactEmail, business.Email, business.HashedPassword).Scan(&id)
	if err != nil {
		return 0, err
	}

	tx, err := conn.Begin(context.TODO())
	if err != nil {
		return 0, err
	}

	for _, addr := range addresses {
		coords, err := utils.GetCoords(addr)
		if err != nil {
			tx.Rollback(context.TODO())
			return 0, err
		}

		_, err = tx.Exec(context.TODO(), addAddressQuery, id, addr, coords.Lat, coords.Lng)
		if err != nil {
			tx.Rollback(context.TODO())
			return 0, err
		}
	}

	err = tx.Commit(context.TODO())
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (bs *BusinessService) GetWithoutAddresses(email string) (Business, error) {
	var res Business
	query := `SELECT * FROM businesses WHERE email=$1`
	err := bs.pool.AcquireFunc(context.TODO(), func(c *pgxpool.Conn) error {
		return c.QueryRow(context.TODO(), query, email).Scan(
			&res.Id,
			&res.Name,
			&res.Description,
			&res.UserContactEmail,
			&res.Email,
			&res.HashedPassword,
		)
	})

	if err != nil {
		return Business{}, err
	}

	return res, nil
}

func (bs *BusinessService) GetById(id int) (Business, error) {
	query := `SELECT * FROM businesses WHERE id=$1`
	conn, err := bs.pool.Acquire(context.TODO())
	if err != nil {
		return Business{}, err
	}

	defer conn.Release()
	var res Business

	row := conn.QueryRow(context.TODO(), query, id)
	err = row.Scan(&res.Id, &res.Name, &res.Description, &res.UserContactEmail, &res.Email, &res.HashedPassword)
	if err != nil {
		return Business{}, err
	}

	return res, nil
}

func (bs *BusinessService) GetAddressesFor(businessId int) ([]Address, error) {
	query := `SELECT * FROM addresses WHERE businessId=$1`

	conn, err := bs.pool.Acquire(context.TODO())
	if err != nil {
		return []Address{}, err
	}

	defer conn.Release()
	rows, err := conn.Query(context.TODO(), query, businessId)
	if err != nil {
		return []Address{}, err
	}

	res := []Address{}
	defer rows.Close()
	for rows.Next() {
		var arrd Address

		err := rows.Scan(&arrd.BusinessId, &arrd.Address, &arrd.Coords.Lat, &arrd.Coords.Lng)
		if err != nil {
			fmt.Println(err)
			continue
		}

		res = append(res, arrd)
	}

	return res, nil
}

type QueryOpts struct {
	SearchOpts *SearchOpts
	Query      string
	Offset     int
	Limit      int
}

type SearchOpts struct {
	Lat    float64
	Lng    float64
	Radius int
}

// radius in kilometers
func (bs *BusinessService) QueryBusinesses(opts QueryOpts) ([]BusinessWithAddress, error) {
	query := `SELECT * FROM businesses LEFT JOIN addresses ON businesses.id=addresses.businessId`
	filters := []string{}
	if opts.SearchOpts != nil {
		filters = append(filters, fmt.Sprintf(`haversine(addresses.lat, addresses.lng, %v, %v) <= %v`, opts.SearchOpts.Lat, opts.SearchOpts.Lng, opts.SearchOpts.Radius))
	}
	if opts.Query != "" {
		filters = append(filters, fmt.Sprintf(`LOWER(businesses.name) LIKE '%%%v%%'`, strings.ToLower(opts.Query)))
	}

	if len(filters) != 0 {
		query += ` WHERE ` + strings.Join(filters, " AND ") + fmt.Sprintf(" LIMIT %v OFFSET %v", opts.Limit, opts.Offset)
	}

	fmt.Println(query)
	conn, err := bs.pool.Acquire(context.TODO())
	if err != nil {
		return []BusinessWithAddress{}, nil
	}
	defer conn.Release()

	rows, err := conn.Query(context.TODO(), query)
	if err != nil {
		return []BusinessWithAddress{}, nil
	}

	res := []BusinessWithAddress{}

	defer rows.Close()

	for rows.Next() {
		var bwa BusinessWithAddress
		values, err := rows.Values()
		if err != nil {
			fmt.Println(err)
			continue
		}

		if values[6] == nil {
			var businessId, address, lat, lng any
			err = rows.Scan(&bwa.Id, &bwa.Name, &bwa.Description, &bwa.UserContactEmail, &bwa.Email, &bwa.HashedPassword, &businessId, &address, &lat, &lng)
			if err != nil {
				fmt.Println(err)
				continue
			}
		} else {
			err = rows.Scan(&bwa.Id, &bwa.Name, &bwa.Description, &bwa.UserContactEmail, &bwa.Email, &bwa.HashedPassword, &bwa.Address.BusinessId, &bwa.Address.Address, &bwa.Coords.Lat, &bwa.Coords.Lng)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
		res = append(res, bwa)

	}

	return res, nil
}

func NewBusinessService(pool *pgxpool.Pool) (*BusinessService, error) {
	bs := &BusinessService{pool: pool}
	err := bs.init()

	if err != nil {
		return nil, err
	}

	return bs, nil
}
