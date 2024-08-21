package api_customer

import (
	"fmt"
	"net/http"
	"os"
	"server/service"
	"server/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ClientRegisterReq struct {
	Email    string `json:"email" validate:"required"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ClientRegusterResp struct {
	UserId int `json:"user_id"`
}

type ClientSigninReq struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ClientSigninResp struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
	UserId      int    `json:"user_id"`
}

type ClaimsClient struct {
	ClientId int `json:"client_id"`
	jwt.StandardClaims
}

func HandleClientRegister(s *service.Services) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body ClientRegisterReq
		err := ctx.Bind(&body)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат"})
			return
		}

		if err = validator.New().Struct(body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
			return
		}

		hashedPassword, err := utils.GenerateHashPassword(body.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		id, err := s.ClientService.Create(service.Client{
			Username:       body.Username,
			Email:          body.Email,
			HashedPassword: hashedPassword,
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, ClientRegusterResp{UserId: id})
	}
}

func HandleClientSignin(s *service.Services) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body ClientSigninReq

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат"})
			return
		}

		if err := validator.New().Struct(body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
			return
		}

		client, err := s.ClientService.GetByEmail(body.Email)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"reason": err.Error()})
			return
		}

		if !utils.CompareHashPassword(body.Password, client.HashedPassword) {
			ctx.JSON(http.StatusUnauthorized, gin.H{"reason": "неверный пароль"})
			return
		}

		expirationTime := time.Now().Add(60 * time.Minute)
		claims := ClaimsClient{
			ClientId:       client.Id,
			StandardClaims: jwt.StandardClaims{Subject: fmt.Sprint(client.Id), ExpiresAt: expirationTime.Unix()},
		}

		token, err := utils.GenerateJWT(claims, os.Getenv("JWTKEY"))
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"reason": "ошибка генерации токена"})
			return
		}
		ctx.JSON(http.StatusOK, ClientSigninResp{AccessToken: token, Username: client.Username, UserId: client.Id})
	}
}
