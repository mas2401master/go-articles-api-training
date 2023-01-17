package router

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mas2401master/go-articles-api-training/docs"
	"github.com/mas2401master/go-articles-api-training/internal/api/controller"
	"github.com/mas2401master/go-articles-api-training/internal/api/middlewares"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func Setup() *gin.Engine {
	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server Petstore server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "petstore.swagger.io"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	//
	app := gin.New()

	// Logging to a file.
	f, _ := os.Create("log/api.log")
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(f)

	// Middlewares
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - - [%s] \"%s %s %s %d %s \" \" %s\" \" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("02/Jan/2006:15:04:05 -0700"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// Middlewares
	app.Use(gin.Recovery())
	app.Use(middlewares.CORS())
	app.NoRoute(middlewares.NoRouteHandler())

	//User Login Route
	app.POST("/api/v1/login", controller.Login)

	// Docs Routes
	app.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middlewares for Routes
	app.Use(middlewares.AuthRequired())

	//User Routes
	SetudUsersRoutes(app)
	//Items Routes
	SetupItemRoutes(app)
	//Promotion Routes
	SetupPromotionRoutes(app)
	//Order  Routes
	SetudOrderRoutes(app)
	//Order_item Routes
	SetudOrderItemRoutes(app)

	return app
}
