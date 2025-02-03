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

type TutorBotInput struct {
	Subject struct {
		Subject          string `json:"subject"`
		Chapter          string `json:"chapter"`
		TopicDescription string `json:"topic_description"`
	} `json:"subject"`
	Student struct {
		StudentClass                 int    `json:"student_class"`
		StudentPerformanceFrom1To100 int    `json:"student_performance_from_1_to_100"`
		StudentLearningStyle         string `json:"student_learning_style"`
		StudentPerformanceLevel      string `json:"student_performance_level"`
		StudyPace                    string `json:"study_pace"`
	} `json:"student"`
	ChatHistory []struct {
		Content string `json:"content"`
		Sender  string `json:"sender"`
	} `json:"chat_history"`
	NewMessage string `json:"new_message"`
}

func PostTutorBot() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var requestBody TutorBotInput
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		requestJSON, err := json.Marshal(requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", "https://ml-service-m5is.onrender.com/tutor/session", bytes.NewBuffer(requestJSON))
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
