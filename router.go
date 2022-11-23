package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func setupRouter() *gin.Engine {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080"}
	// assets, err := fs.Sub(staticAssets, "build")
	// if err != nil {
	// 	fmt.Println("build folder is not readable")
	// 	return nil
	// }
	// assetsFS := http.FS(assets)
	router := gin.Default()
	router.Use(cors.New(config))
	// router.GET("/", func(c *gin.Context) {
	// 	c.Redirect(http.StatusTemporaryRedirect, "/assets/index.html")
	// })
	// router.StaticFS("/assets/", assetsFS)
	router.GET("/prom-target/:id", prometheusHandler)
	router.GET("/health/", healthHandler)
	api := router.Group("/api")
	{
		auth := api.Group("/")
		{
			auth.POST("/target/", createTargetHandler)
			auth.GET("/target/:id", getTargetHandler)
			auth.POST("/target/:id", updateTargetHandler)
			auth.DELETE("/target/:id", removeTargetHandler)
			auth.GET("/targets/", getTargetsHandler)
		}

	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
}
