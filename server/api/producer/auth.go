package api_producer

import (
	"fmt"
	"net/http"
	"os"
	"server/service"
	"server/utils"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type ClaimsProducer struct {
	ProducerId int `json:"producer_id"`
	jwt.StandardClaims
}

type RegisterBody struct {
	service.Business
	Password  string   `json:"password"`
	Addresses []string `json:"addresses"`
}

type SignInBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func RedisterProducer(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body RegisterBody

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат"})
			return
		}

		if err := utils.Validate.Struct(body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат"})
			return
		}

		hashed_password, err := utils.GenerateHashPassword(body.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "ошибка генерации пароля"})
			return
		}

		producer := body.Business
		producer.HashedPassword = hashed_password

		id, err := s.Bs.Create(producer, body.Addresses)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
			return
		}

		producer.Id = id
		c.JSON(http.StatusOK, producer)
	}
}

func SigninProducer(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body SignInBody

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат"})
			return
		}

		existingProducer, err := s.Bs.GetWithoutAddresses(body.Email)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"reason": "бизнес не найден"})
			return
		}

		if !utils.CompareHashPassword(body.Password, existingProducer.HashedPassword) {
			c.JSON(http.StatusUnauthorized, gin.H{"reason": "неверный пароль"})
			return
		}

		expirationTime := time.Now().Add(60 * time.Minute)
		claims := ClaimsProducer{
			ProducerId:     existingProducer.Id,
			StandardClaims: jwt.StandardClaims{Subject: fmt.Sprint(existingProducer.Id), ExpiresAt: expirationTime.Unix()},
		}

		token, err := utils.GenerateJWT(claims, os.Getenv("JWTKEY"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"reason": "ошибка генерации токена"})
			return
		}
		c.SetCookie("token", token, int(expirationTime.Unix()), "/", strings.Split(os.Getenv("ADDRESS"), ":")[0], false, true)

		c.JSON(http.StatusOK, gin.H{"access_token": token, "business_id": existingProducer.Id})
	}
}
