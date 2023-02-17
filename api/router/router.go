package router

import (
	"log"
	"net/http"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/kumareswaramoorthi/companies/api/controller"
	"github.com/kumareswaramoorthi/companies/api/database"
	"github.com/kumareswaramoorthi/companies/api/logging"
	"github.com/kumareswaramoorthi/companies/api/middleware"
	"github.com/kumareswaramoorthi/companies/api/repository"
	"github.com/kumareswaramoorthi/companies/api/service"
	docs "github.com/kumareswaramoorthi/companies/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	//swagger init
	docs.SwaggerInfo.Title = "COMPANIES API"
	docs.SwaggerInfo.Description = "This lists down the endpoints that are part of COMPANIES API server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

func SetupRouter() *gin.Engine {
	//Get default router from gin
	router := gin.Default()

	//create a global logger for the server
	apiLoggerEntry := logging.NewLoggerEntry()

	//router use the global logger
	router.Use(logging.LoggingMiddleware(apiLoggerEntry))
	router.Use(requestid.New())

	dbConn, err := database.NewPostgresDB(database.GetDBConfig())
	if err != nil {
		log.Fatal(err)
	}

	companyRepo := repository.NewRepository(dbConn)
	companySvc := service.NewService(companyRepo)
	companyCtrl := controller.NewController(companySvc)

	loginService := service.StaticLoginService()
	jwtService := service.JWTAuthService()
	loginCtrl := controller.NewLoginController(loginService, jwtService)

	v1 := router.Group("/api/v1")

	//health check API
	v1.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})

	v1.POST("/login", loginCtrl.Login)
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	v1.GET("/company/:id", companyCtrl.GetCompany)
	v1.POST("/company", middleware.AuthorizeJWT(), companyCtrl.CreateCompany)
	v1.PATCH("/company/:id", middleware.AuthorizeJWT(), companyCtrl.UpdateCompany)
	v1.DELETE("/company/:id", middleware.AuthorizeJWT(), companyCtrl.DeleteCompany)

	return router
}
