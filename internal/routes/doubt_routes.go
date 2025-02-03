package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func DoubtRoutes(r *gin.Engine) {

	r.Use(middleware.AuthenticationStudent())
	r.POST("/doubt/ask", controllers.PostDoubtBot())

}
