package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func QuizRoutes(r *gin.Engine) {

	r.Use(middleware.AuthenticationStudent())
	r.POST("/quiz", controllers.PostQuizBot())

}
