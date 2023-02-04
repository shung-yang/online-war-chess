package controller

import (
  "online_chess/util"
  "online_chess/model"
  "golang.org/x/crypto/bcrypt"
  "github.com/gin-gonic/gin"
  "net/http"
  "log"
  "strings"
)

type player struct {
  Email  string  `json:"email"`
  Password string  `json:"password"`
}

func TestAuth(c *gin.Context) {
	result, _ := util.VerifyToken(strings.Split(c.GetHeader("Authorization"), " ")[1])
	if result {
		c.IndentedJSON(http.StatusOK, result)
	} else {
		c.IndentedJSON(http.StatusUnauthorized, result)
	}
}

func ChangeAccountInfo(c *gin.Context) {  //just for test, will discard after build reset password func
	type account_info struct {
		Password string `json:"password"`
	}
	var new_info account_info
	c.BindJSON(&new_info)
	hash_password, _ := bcrypt.GenerateFromPassword([]byte(new_info.Password), 10)
	model.SetPlayerPassword(hash_password)
  c.IndentedJSON(http.StatusOK, new_info.Password)
}

func Login(c *gin.Context) {
	var inputs player
	err := c.BindJSON(&inputs)
	if err != nil {
		log.Fatal("bindjson fail:", err)
	}
	log.Printf("login inputs: %v email: %s password: %s\n", inputs, inputs.Email, inputs.Password)
	
  if hash_password, err := model.GetPlayerPassword(inputs.Email); err != nil {
		c.IndentedJSON(http.StatusUnauthorized, "user login fail")
	} else {
		valid_password_result := bcrypt.CompareHashAndPassword(hash_password, []byte(inputs.Password))
		if valid_password_result == nil {
			type token_response struct {
				Token string
			}
			c.IndentedJSON(http.StatusOK, token_response{ Token: util.GenerateToken()})
		} else {
			c.IndentedJSON(http.StatusUnauthorized, "user login fail")
		}
	}
}