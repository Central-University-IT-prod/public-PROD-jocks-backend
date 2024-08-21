package api_customer

import (
	"net/http"
	"server/service"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GetPromocodesRespItem struct {
	Token       string `json:"token"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Activated   bool   `json:"activated"`
}

func HandleGetPromocodes(s *service.Services) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		customerId, err := strconv.Atoi(ctx.Param("customerId"))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат пути"})
			return
		}

		promocodes, err := s.PromocodeService.GetForUser(customerId)

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, utils.Map(promocodes, func(item service.GeneratedPromocodeFull) GetPromocodesRespItem {
			return GetPromocodesRespItem{
				Token:       item.GeneratedPromocode.Token,
				Name:        item.Promocode.Name,
				Description: item.Promocode.Description,
				Activated:   item.GeneratedPromocode.Activated,
			}
		}))
	}
}
