package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func TutorRoutes(r *gin.Engine) {

	r.Use(middleware.AuthenticationStudent())
	r.GET("/tutor/classes", controllers.GetTutorClasses())
	r.POST("/tutor/session", controllers.PostTutorBot())

}
