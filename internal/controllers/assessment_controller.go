package controllers

import (
	"backend/internal/models"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type StaticAssessmentRequest struct {
	Responses []int `json:"responses"`
}

type DynamicAssessmentRequest struct {
	Subject string `json:"subject"`
	Scores  []int  `json:"scores"`
}

type LearningStyle struct {
	Style string `json:"style"`
}

func GetStaticAssessment() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := http.Get("https://ml-service-m5is.onrender.com/assessment/static")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch assessment"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}

func PostStaticAssessment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var requestBody StaticAssessmentRequest
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		requestJSON, err := json.Marshal(requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", "https://ml-service-m5is.onrender.com/assessment/static", bytes.NewBuffer(requestJSON))
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

		var currentStudent models.Student

		student_id := c.Keys["email"]
		colerr := studentCollection.FindOne(ctx, bson.M{"student_id": student_id}).Decode(&currentStudent)
		if colerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find student"})
			return
		}

		var styleData LearningStyle
		err = json.Unmarshal([]byte(body), &styleData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			return
		}

		currentStudent.LearningStyle = append(currentStudent.LearningStyle, styleData.Style)

		_, err = studentCollection.UpdateOne(ctx, bson.M{"student_id": student_id}, bson.M{"$set": bson.M{"learning_style": currentStudent.LearningStyle}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
			return
		}

		c.Data(resp.StatusCode, "application/json", body)
	}
}

func PostDynamicAssessment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var requestBody DynamicAssessmentRequest
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		requestJSON, err := json.Marshal(requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", "https://ml-service-m5is.onrender.com/assessment/dynamic", bytes.NewBuffer(requestJSON))
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
