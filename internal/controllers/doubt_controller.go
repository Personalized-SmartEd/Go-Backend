package controllers

import (
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

func PostDoubtBot() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var requestBody utils.DoubtBotInput
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var payload map[string]interface{}

		tempBytes, err := json.Marshal(requestBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
			return
		}

		if err := json.Unmarshal(tempBytes, &payload); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert request body"})
			return
		}

		delete(payload, "newchat")

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
		err = studentCollection.FindOne(ctx, bson.M{"student_id": studentID}).Decode(&currentStudent)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve student information"})
			return
		}

		payload["student"] = map[string]interface{}{
			"student_class":                     currentStudent.Class,
			"student_performance_from_1_to_100": currentStudent.Performance,
			"student_learning_style":            currentStudent.LearningStyle,
			"student_performance_level":         currentStudent.PerformanceLvl,
			"study_pace":                        currentStudent.Pace,
		}

		requestJSON, err := json.Marshal(payload)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request"})
			return
		}

		req, err := http.NewRequestWithContext(ctx, "POST", "https://ml-service-m5is.onrender.com/doubt/ask", bytes.NewBuffer(requestJSON))
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
