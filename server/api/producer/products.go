package api_producer

import (
	"net/http"
	"server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddProduct(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		var product service.Product
		if err := c.ShouldBindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат"})
			return
		}

		id, err := s.Ps.Create(&product)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "бизнес не найден"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}

func RemoveProduct(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("product_id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "неверный формат пути запроса"})
			return
		}
		removed, err := s.Ps.Remove(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "продукт не найден"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"deleted": removed})

	}
}
