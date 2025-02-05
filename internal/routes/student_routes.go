package routes

import (
	"backend/internal/controllers"

	"github.com/gin-gonic/gin"
)

func StudentRoutes(r *gin.Engine) {

	r.GET("/students", controllers.GetStudents())
	r.GET("/student/:student_id", controllers.GetStudent())
	r.POST("/login", controllers.Login())
	r.POST("/signup", controllers.SignUp())
	r.PUT("/student/:student_id", controllers.UpdateStudent())
	r.DELETE("/student/:student_id", controllers.DeleteStudent())

}
