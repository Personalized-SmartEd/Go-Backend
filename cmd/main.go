package main

import (
	"log"
	"net/http"

	"backend/internal/config"
	"backend/internal/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()

	r := gin.New()
	r.Use(gin.Logger())

	routes.TeacherRoutes(r)
	routes.StudentRoutes(r)
	routes.AssessmentRoutes(r)
	routes.QuizRoutes(r)
	routes.TutorRoutes(r)
	routes.DoubtRoutes(r)
	routes.RecommendationRoutes(r)
	routes.ClassroomRoutesStudent(r)
	routes.ClassroomRoutesTeacher(r)
	routes.StudyfeatRoutes(r)

	port := "8080"

	log.Printf("Starting server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
