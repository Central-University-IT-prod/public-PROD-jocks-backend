package main

import (
	"context"
	"fmt"
	"os"
	"server/config"
	"server/router"
	"server/service"
	"server/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectToPostres(cfg *config.Config, maxRetries int) (*pgxpool.Pool, error) {
	var (
		pool       *pgxpool.Pool
		err        error
		retryCount int
	)

	pool, err = pgxpool.New(context.TODO(), cfg.PostgresConn)

	for err != nil && retryCount < maxRetries {
		fmt.Println(err, "retrying to connect")
		time.Sleep(time.Second * 5)
		pool, err = pgxpool.New(context.TODO(), cfg.PostgresConn)
		retryCount++
	}

	if err != nil {
		return nil, err
	}

	retryCount = 0
	err = pool.Ping(context.TODO())

	for err != nil && retryCount < maxRetries {
		fmt.Println(err, "retrying to connect")
		time.Sleep(time.Second * 5)
		err = pool.Ping(context.TODO())
		retryCount++
	}

	return pool, nil
}

func main() {

	cfg, err := config.Get()
	if err != nil {
		fmt.Println(err)
		return
	}

	pool, err := ConnectToPostres(cfg, 5)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pool.Close()

	services, err := service.NewService(pool)
	if err != nil {
		fmt.Println("ошибка в подключении сервисов:",err)
		return
	}

	fmt.Println("Connected to Postgres")

	utils.RegisterAllValidators()

	r := gin.Default()
	router.RouteAll(r, services)
	err = r.Run(os.Getenv("ADDRESS"))

	if err != nil {
		fmt.Println(err)
		return
	}
}
