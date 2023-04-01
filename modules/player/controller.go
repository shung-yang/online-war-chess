package player

import (
  "online_chess/util"
  "golang.org/x/crypto/bcrypt"
  "github.com/gin-gonic/gin"
  "net/http"
  "log"
  "strings"
)

type login_player struct {
  Email  string  `json:"email" binding:"required,email" example:"example@gnka.com"`
  Password string  `json:"password" binding:"min=8" example:"password"`
}

type Register_player struct {
	Name string `json:"name" binding:"required" example:"wilson"`
  Email  string  `json:"email" binding:"required,email" example:"example@gnka.com"`
  Password string  `json:"password" binding:"min=8" example:"asdqwezxc"`
}

type Player struct {
	Id int
	Token string
	Name string
	Email string
	Password string
  Level  int8
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
	SetPlayerPassword(hash_password)
  c.IndentedJSON(http.StatusOK, new_info.Password)
}

// @Summary      player register
// @Description  register player in game
// @Tags         player
// @Accept       json
// @Produce      json
// @Param register_inputs body Register_player true "register player"
// @Success      200  {object}  object{token=string}
// @Failure      500  {object}  object{error=string}
// @Router       /register [post]
func Register(c *gin.Context) {
	var inputs Register_player
	err := c.ShouldBindJSON(&inputs)
	if err != nil {
		log.Println("register bind err:", err)
		c.JSON(http.StatusBadRequest, gin.H{ "error": err.Error() })
	} else {
		token, err := AddNewPlayer(inputs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{ "error": err.Error() })
		} else {
			c.JSON(http.StatusOK, gin.H{ "token": token })
		}
	}
}

// @Summary      player login
// @Description  player login
// @Tags         player
// @Accept       json
// @Produce      json
// @Param login_inputs body login_player true "login player"
// @Success      200  {object}  login_player
// @Failure      401  {object}  object{error=string} "email or password not correct"
// @Failure      500  {object}  object{error=string}
// @Router       /login [post]
func Login(c *gin.Context) {
	var inputs login_player
	err := c.ShouldBindJSON(&inputs)
	if err != nil {
		log.Fatal("bind json fail:", err)
	}
	log.Printf("login inputs: %v email: %s password: %s\n", inputs, inputs.Email, inputs.Password)
	
  if hash_password, err := GetPlayerPassword(inputs.Email); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{ "error": err.Error()})
	} else {
		valid_password_result := bcrypt.CompareHashAndPassword(hash_password, []byte(inputs.Password))
		if valid_password_result == nil {
			token := util.GenerateToken()
			err := UpdatePlayerToken(inputs.Email, token)
			if err != nil { 
				c.IndentedJSON(http.StatusInternalServerError, gin.H{ "error": err.Error() })
			} else {
				c.JSON(http.StatusOK, gin.H{ "token": token })
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{ "error": "Email or Password is not correct" })
		}
	}
}