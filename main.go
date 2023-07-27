package main

import (
	_ "online_chess/docs"
	"online_chess/model"
	"online_chess/modules/game"
	"online_chess/modules/player"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	is_success := model.GetDatabaseHandle()
	if is_success {
		router := gin.Default()
		router.GET("/testauth", player.TestAuth)
		router.POST("/account", player.ChangeAccountInfo)
		router.POST("/login", player.Login)
		router.POST("/register", player.Register)

		wsServer := game.NewWebSocketServer()
		go wsServer.Run()
		router.GET("/ws", attachWsServer(wsServer), game.RunClient)

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		router.Run("0.0.0.0:8080")
	}
}
