package controllers

import (
	"backend/internal/config"
	"backend/internal/models"
	"backend/internal/utils"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetStaticAssessment() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		resp, err := http.Get(config.BaseURL + "/assessment/static")
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

		var requestBody utils.StaticAssessmentInput
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		requestJSON, err := json.Marshal(requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", config.BaseURL+"/assessment/static", bytes.NewBuffer(requestJSON))
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

		studentIDIfc, exists := c.Get("student_id")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Student ID not found in context"})
			return
		}
		studentID, ok := studentIDIfc.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid student ID format"})
			return
		}

		var currentStudent models.Student
		colerr := studentCollection.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&currentStudent)
		if colerr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find student"})
			return
		}

		var styleData utils.LearningStyle
		err = json.Unmarshal([]byte(body), &styleData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			return
		}

		currentStudent.LearningStyle = styleData.Style

		_, err = studentCollection.UpdateOne(ctx, bson.M{"student_id": studentID}, bson.M{"$set": bson.M{"learning_style": currentStudent.LearningStyle}})
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

		var requestBody utils.DynamicAssessmentInput
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		requestJSON, err := json.Marshal(requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", config.BaseURL+"/assessment/dynamic", bytes.NewBuffer(requestJSON))
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
