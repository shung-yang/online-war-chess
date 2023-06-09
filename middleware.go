package main

import (
	"net/http"
	"online_chess/modules/game"
	"online_chess/modules/player"
	"online_chess/util"
	"strings"

	"github.com/gin-gonic/gin"
)

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := strings.Split(c.GetHeader("Authorization"), " ")[1]
		is_valid, err := util.VerifyToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		} else if !is_valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is not valid"})
		} else {
			_, err := player.QueryIdByToken(token)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			}
		}
		c.Next()
	}
}

func attachWsServer(wsServer *game.WsServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("wsServer", wsServer)
		c.Next()
	}
}
