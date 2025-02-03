package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func AssessmentRoutes(r *gin.Engine) {

	r.Use(middleware.AuthenticationStudent())
	r.GET("/assessment/static", controllers.GetStaticAssessment())
	r.POST("/assessment/static", controllers.PostStaticAssessment())
	r.POST("/assessment/dynamic", controllers.PostDynamicAssessment())

}
