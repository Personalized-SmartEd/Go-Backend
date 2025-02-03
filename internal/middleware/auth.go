package middleware

import (
	"backend/internal/helper"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthenticationStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		claims, err := helper.ValidateStudentToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("student_id", claims.StudentID)
		c.Set("name", claims.Name)
		c.Set("email", claims.Email)
		c.Set("class", claims.Class)

		c.Next()
	}
}

func AuthenticationTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("token")
		if clientToken == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("No Authorization header provided")})
			c.Abort()
			return
		}

		claims, err := helper.ValidateTeacherToken(clientToken)
		if err != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			c.Abort()
			return
		}

		c.Set("teacher_id", claims.TeacherID)
		c.Set("name", claims.Name)
		c.Set("email", claims.Email)
		c.Set("class", claims.SchoolCode)

		c.Next()
	}
}
