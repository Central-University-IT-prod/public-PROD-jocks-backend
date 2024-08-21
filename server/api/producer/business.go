package api_producer

import (
	"net/http"
	"server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetReviews(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("buisnessId"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		reviews,_, err := s.Rs.GetReviews(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(reviews) == 0 {
			c.JSON(http.StatusOK, make([]string, 0))
			return
		}

		c.JSON(http.StatusOK, reviews)
	}
}
