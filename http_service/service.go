package httpservice

import (
	"github.com/Brandon-lz/myopcua/http_service/routers"
	webhookrouters "github.com/Brandon-lz/myopcua/http_service/routers/webhook_routers"

	"github.com/gin-gonic/gin"

	_ "github.com/Brandon-lz/myopcua/docs" // 引入文档目录

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

	router.GET("/health", healthCheck)

	url := ginSwagger.URL("http://localhost:8080/docs/doc.json") // The URL pointing to API definition

	// swagger:  http://localhost:8080/docs/index.html
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	v1 := router.Group("/api/v1")
	routers.RegisterRoutes(v1)
	webhookrouters.GetAllWebhookConfigFromDB()

	router.Run("0.0.0.0:8080")
}

// healthCheck 路由
// @Summary  healthCheck 路由
// @Description  healthCheck 路由
// @Tags     default
// @Accept   json
// @Produce  json
// @Success  200  {string}  pong  "pong"
// @Router   /health [get]
func healthCheck(c *gin.Context) {
	// c.Header("Content-Type", "charset=utf-8")
	c.String(200, "I am healthy")
}