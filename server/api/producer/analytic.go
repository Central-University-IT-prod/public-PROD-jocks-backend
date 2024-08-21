package api_producer

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"server/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewResp struct {
	Stat    service.Count
	Reviews []service.Reviews
}

func ConvertToCSV(reviews *ReviewResp) ([]byte, error) {
	var csvData [][]string
	general := []string{fmt.Sprintf("|GENERAL| anon:%d,user:%d",reviews.Stat.Anon,reviews.Stat.User)}
	csvData = append(csvData, general)

	header := []string{"created_at", "rating"}
	csvData = append(csvData, header)

	for _, s := range reviews.Reviews {
		record := []string{
			s.CreatedAt.Format("2006-01-02 15:04:05"),
			fmt.Sprint(s.Rating),
		}
		csvData = append(csvData, record)
	}

	buf := bytes.NewBuffer([]byte{})
	writer := csv.NewWriter(buf)

	err := writer.WriteAll(csvData)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(buf)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

func ExportData(s *service.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		reviews, count, err := s.Rs.GetReviews(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		resp := ReviewResp{
			Stat:    count,
			Reviews: reviews,
		}

		csvBytes, err := ConvertToCSV(&resp)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename=reviews.csv")
		c.Data(http.StatusOK, "text/csv", csvBytes)
	}
}
