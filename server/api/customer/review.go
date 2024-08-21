package api_customer

import (
	"net/http"
	"server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type baseReview struct {
	Form_id int `json:"id"`
	Value   int `json:"value"`
}
type productReview struct {
	Product_id int `json:"id"`
	Value      int `json:"value"`
}
type reviewRequest struct {
	CustomerId    int             `json:"customer_id"`
	BusinessId    int             `json:"business_id"`
	CheckId       int             `json:"check_id"`
	BaseReview    []baseReview    `json:"base_review_fields"`
	ProductReview []productReview `json:"products_review_fields"`
}

func Review(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request reviewRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат"})
			return
		}
		review_id, err := s.Rs.Review(request.CustomerId, request.BusinessId, request.CheckId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for _, v := range request.BaseReview {
			_, err := s.Rs.ReviewProducer(v.Form_id, v.Value, review_id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		for _, v := range request.ProductReview {
			_, err := s.Rs.ReviewProduct(v.Product_id, v.Value, review_id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		promocode, token, err := s.PromocodeService.GeneratePromocode(request.BusinessId, request.CustomerId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"reason": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"id": review_id,
			"promocode": gin.H{
				"name":        promocode.Name,
				"description": promocode.Description,
				"token":       token,
			},
		})
	}
}

func GetUserReviews(s *service.Services)gin.HandlerFunc{
	return func(c *gin.Context) {
		id,err := strconv.Atoi(c.Param("user_id"))
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err})
			return
		}
		reviews,err := s.ClientService.GetReviews(id)
		if err != nil{
			c.JSON(http.StatusInternalServerError,gin.H{"error":err})
			return
		}
		c.JSON(http.StatusOK,reviews)
	}
}