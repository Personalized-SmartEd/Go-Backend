package routes

import (
	"backend/internal/controllers"

	"github.com/gin-gonic/gin"
)

func StudyfeatRoutes(r *gin.Engine) {

	r.POST("/studyfeat/summarise/pdf", controllers.PostPdfSummariser())
	r.POST("/studyfeat/summarise/youtube", controllers.PostYoutubeVideoSummariser())
	r.POST("/studyfeat/mindmap", controllers.PostGenerateMindMap())

}
