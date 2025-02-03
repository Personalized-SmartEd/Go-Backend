package routes

import (
	"backend/internal/controllers"

	"github.com/gin-gonic/gin"
)

func StudentRoutes(r *gin.Engine) {

	r.GET("/students", controllers.GetStudents())
	r.GET("/student/:student_id", controllers.GetStudent())
	r.POST("/signup", controllers.SignUp())
	r.POST("/login", controllers.Login())

}
