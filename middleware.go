package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"online_chess/modules/player"
	"online_chess/util"
	"strings"
)

func auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		exceptions := []string{"/register", "/login"}
		if !util.SliceContains(exceptions, c.Request.URL.Path) {
			token := strings.Split(c.GetHeader("Authorization"), " ")[1]
			is_valid, err := util.VerifyToken(token)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{ "error": err.Error() })
			} else if !is_valid {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{ "error": "token is not valid" })
			} else {
				_, err := player.QueryIdByToken(token)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{ "error": err.Error() })
				}
			}
		}
		c.Next()
	}
}