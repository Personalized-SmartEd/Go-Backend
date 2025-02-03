package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type RecommentdationInput struct {
	LearningStyle      string   `json:"learning_style"`
	CurrentLevel       string   `json:"current_level"`
	WeakAreas          []string `json:"weak_areas"`
	PerformanceHistory []int    `json:"performance_history"`
	PreferredPace      string   `json:"preferred_pace"`
	AvailableHours     int      `json:"available_hours"`
}

func PostRecommendation() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var requestBody RecommentdationInput
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		requestJSON, err := json.Marshal(requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", "https://ml-service-m5is.onrender.com/reccomend/generate_study_plan", bytes.NewBuffer(requestJSON))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		c.Data(resp.StatusCode, "application/json", body)
	}
}
