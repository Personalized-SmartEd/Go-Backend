package routes

import (
	"backend/internal/controllers"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func ClassroomRoutesStudent(r *gin.Engine) {

	r.Use(middleware.AuthenticationStudent())
	r.GET("/classroom/join/:classroom_id", controllers.JoinClassroom())
	r.GET("/classroom/leave/:classroom_id", controllers.LeaveClassroom())

}

func ClassroomRoutesTeacher(r *gin.Engine) {

	r.Use(middleware.AuthenticationTeacher())
	r.GET("/classroom/get/teacher/:teacher_id", controllers.GetClassroomByTeacherID())
	r.POST("/classroom/create/:teacher_id", controllers.CreateClassroom())
	r.GET("/classroom/delete/:classroom_id", controllers.DeleteClassroom())
	r.GET("/classroom/get/:classroom_id", controllers.GetClassroomByID())
	r.GET("/classroom/get/students/:classroom_id", controllers.GetStudentsByClassroomID())
	r.GET("/classroom/get/teachers/:classroom_id", controllers.GetTeachersByClassroomID())
	r.GET("/classroom/get/school/:school_code", controllers.GetClassroomsBySchoolCode())
	r.GET("/classroom/get/school/:school_code/:teacher_id", controllers.GetClassroomsBySchoolCodeTeacherID())

}
