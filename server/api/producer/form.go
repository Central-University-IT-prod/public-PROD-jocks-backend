package api_producer

import (
	"net/http"
	"server/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PutFormReqParam struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required"`
}

func HandlePutForm(s *service.Services) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		businessIdStr := ctx.Param("businessId")
		businessId, err := strconv.Atoi(businessIdStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат пути"})
			return
		}

		body := []PutFormReqParam{}
		err = ctx.Bind(&body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
			return
		}

		formParams := []service.Parameter{}
		for _, param := range body {
			err := validator.New().Struct(param)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
				return
			}
			formParams = append(formParams, service.Parameter{Name: param.Name, Type: param.Type})
		}

		err = s.Fs.PutParams(businessId, formParams)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func HandleGetForm(s *service.Services) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		businessIdStr := ctx.Param("businessId")
		businessId, err := strconv.Atoi(businessIdStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат пути"})
			return
		}

		params, err := s.Fs.GetParams(businessId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, params)
	}
}
