package main

import (
  "github.com/gin-gonic/gin"
  "log"
  "github.com/joho/godotenv"
  "online_chess/model"
	"online_chess/modules/player"
	"online_chess/modules/room"
	_ "online_chess/docs"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
)

// @title           Online War Chess
// @version         1.0
// @description     Online War Chess backend API server.
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:8080
// @securityDefinitions.basic  BasicAuth
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
  err := godotenv.Load() //load .env file
  if err != nil {
    log.Fatal("Error loading .env file", err)
  } else {
    is_success := model.GetDatabaseHandle()
    if is_success {
      router := gin.Default()
			router.Use(auth())
      router.GET("/testauth", player.TestAuth)
      router.POST("/account", player.ChangeAccountInfo)
      router.POST("/login", player.Login)
			router.POST("/register", player.Register)
			router.POST("/room", room.Create)

			router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
      router.Run("localhost:8080")
    }
  }
}