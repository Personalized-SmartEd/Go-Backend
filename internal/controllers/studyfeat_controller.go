package controllers

import (
	"backend/internal/config"
	"backend/internal/utils"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func PostPdfSummariser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var requestBody utils.PdfSummariserInput
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		pdfurl := requestBody.PdfUrl
		req, err := http.NewRequestWithContext(ctx, "POST", config.BaseURL+"/v2/summarise/?pdf_url="+string(pdfurl), nil)
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

func PostYoutubeVideoSummariser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 200*time.Second)
		defer cancel()

		var requestBody utils.YoutubeVideoSummariserInput
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		youtubeUrl := requestBody.YoutubeUrl
		req, err := http.NewRequestWithContext(ctx, "POST", config.BaseURL+"/v2/summarise-youtube/?youtube_url="+string(youtubeUrl), nil)
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

func PostGenerateMindMap() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var requestBody utils.MindMapInput
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		topic := requestBody.Topic
		req, err := http.NewRequestWithContext(ctx, "POST", config.BaseURL+"/v2/generate_mindmap/"+string(topic), nil)
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
