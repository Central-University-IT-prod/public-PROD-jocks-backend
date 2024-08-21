package api_producer

import (
	"fmt"
	"net/http"
	"server/service"
	"server/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

type GetCheckResp struct {
	BusinessId           int           `json:"business_id"`
	BaseReviewFields     []ReviewField `json:"base_review_fields"`
	ProductsReviewFields []ReviewField `json:"products_review_fields"`
}

type ReviewField struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func GenerateCheck(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		var check service.Check
		if err := c.ShouldBindJSON(&check); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"reason": "неверный формат"})
			return
		}
		id, err := s.Cs.Create(&check)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		url := fmt.Sprintf("http://158.160.122.246/check/%d/", id)

		qrCode, err := qrcode.Encode(url, qrcode.Medium, 256)
		if err != nil {
			c.String(http.StatusInternalServerError, "Ошибка при создании QR-кода: %v", err)
			return
		}

		c.Data(http.StatusOK, "image/png", qrCode)

	}
}

func GetCheck(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат пути запроса"})
			return
		}

		check, err := s.Cs.GetById(id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		products, err := s.Cs.GetProductsById(id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		params, err := s.Fs.GetParams(check.BusinessId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(params) > 5 {
			params = utils.RandomChoises(params, 5)
		}

		resp := GetCheckResp{
			BusinessId: check.BusinessId,
			BaseReviewFields: utils.Map(params, func(item service.Parameter) ReviewField {
				return ReviewField{
					Id:   item.Id,
					Name: item.Name,
					Type: item.Type,
				}
			}),
			ProductsReviewFields: utils.Map(products, func(item service.Product) ReviewField {
				return ReviewField{
					Id:   item.Id,
					Name: item.Name,
					Type: "STAR",
				}
			}),
		}

		c.JSON(http.StatusOK, resp)
	}
}
