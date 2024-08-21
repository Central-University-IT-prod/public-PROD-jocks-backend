package api_producer

import (
	"net/http"
	"server/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PutPromocodeReqResp struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func HandlePutPromocode(s *service.Services) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		businessId, err := strconv.Atoi(ctx.Param("businessId"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат пути"})
			return
		}

		var body PutPromocodeReqResp
		err = ctx.Bind(&body)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
			return
		}

		if err = validator.New().Struct(body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
			return
		}

		err = s.PromocodeService.Put(service.Promocode{BusinessId: businessId, Name: body.Name, Description: body.Description})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, body)
	}
}
