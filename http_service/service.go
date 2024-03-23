package httpservice

import (
	"earth/http_service/routers"

	"github.com/gin-gonic/gin"

	_ "earth/docs" // 引入文档目录

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Start() {
	router := gin.Default()
	// CORS for https://foo.com and https://github.com origins, allowing:
	// - PUT and PATCH methods
	// - Origin header
	// - Credentials share
	// - Preflight requests cached for 12 hours
	router.Use(ginCors())
	router.Use(gin.CustomRecovery(ErrorHandler))
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})
	url := ginSwagger.URL("http://localhost:8080/docs/doc.json") // The URL pointing to API definition

	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	v1 := router.Group("/api/v1")
	routers.RegisterRoutes(v1)

	router.Run("0.0.0.0:8080")

}
