package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RecommendationRoutes(r *gin.Engine) {

	r.Use(middleware.AuthenticationStudent())
	r.POST("/recommend/generate_study_plan", controllers.PostRecommendation())

}
