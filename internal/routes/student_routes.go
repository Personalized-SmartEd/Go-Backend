package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func StudentRoutes(r *gin.Engine) {

	r.POST("/login", controllers.Login())

	r.Use(middleware.AuthenticationTeacher())
	r.GET("/students", controllers.GetStudents())
	r.GET("/student/:student_id", controllers.GetStudent())
	r.POST("/student/add", controllers.AddStudent())
	r.PUT("/student/:student_id", controllers.UpdateStudent())
	r.DELETE("/student/:student_id", controllers.DeleteStudent())

}
