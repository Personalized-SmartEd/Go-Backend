package routes

import (
	"backend/internal/controllers"

	"github.com/gin-gonic/gin"
)

func TeacherRoutes(r *gin.Engine) {

	r.GET("/teacher/:teacher_id", controllers.GetTeacher())
	r.POST("/teacher/signup", controllers.SignUpTeacher())
	r.POST("/teacher/login", controllers.LoginTeacher())
	r.GET("/teacher/logout", controllers.LogOutTeacher())

}
