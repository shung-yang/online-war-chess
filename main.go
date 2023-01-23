package main

import (
  "fmt"
	"strings"
	"net/http"
  "github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt"
	"time"
	"database/sql"
	"log"
	"github.com/go-sql-driver/mysql"
  "reflect"
  "os"
  "github.com/joho/godotenv"
)

func readEnvVariable(key string) (string, bool) {
  value, not_empty := os.LookupEnv(key)
  if !not_empty {
      fmt.Printf("%s not set\n", key)
      return "", false
  }
  return value, true
}

func generateToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(10 * time.Minute).Unix(),
		//"player": player id
	})
	tokenString, _ := token.SignedString([]byte("lakgfnlawng"))
	return tokenString
}

func verifyToken(request_token string) (bool, error) {
	token, err := jwt.Parse(request_token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("lakgfnlawng"), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid{
		var expiry_time, _ = claims["exp"].(float64)
		fmt.Println("claims exp:", expiry_time, reflect.TypeOf(claims["exp"]), time.Unix(int64(expiry_time), 0))
		return true, err
	} else {
		return false, err
	}
}

func testAuth(c *gin.Context) {
	result, _ := verifyToken(strings.Split(c.GetHeader("Authorization"), " ")[1])
	if result {
		c.IndentedJSON(http.StatusOK, result)
	} else {
		c.IndentedJSON(http.StatusUnauthorized, result)
	}
}

var db *sql.DB
func getDatabaseHandle() bool {
  db_user, user_ok := readEnvVariable("DB_USER")
  db_pwd, pwd_ok := readEnvVariable("DB_PWD")
  if !user_ok || !pwd_ok {
    fmt.Println("db env is empty!!!")
    return false
  }
	cfg := mysql.Config{
    User:   db_user,
    Passwd: db_pwd,
    Net:    "tcp",
    Addr:   "127.0.0.1:3306",
    DBName: "online_war_chess",
  }
  var err error
  db, err = sql.Open("mysql", cfg.FormatDSN())
  if err != nil {
    log.Fatal(err)
    return false
  }
  pingErr := db.Ping()
  if pingErr != nil {
    log.Fatal(pingErr)
    return false
  }
  return true
}

type player struct {
    Email  string  `json:"email"`
    Password string  `json:"password"`
}

func changeAccountInfo(c *gin.Context) {  //just for test, will discard after build reset password func
	type account_info struct {
		Password string `json:"password"`
	}
	var new_info account_info
	c.BindJSON(&new_info)
	hash_password, _ := bcrypt.GenerateFromPassword([]byte(new_info.Password), 10)
	db.Exec("UPDATE player SET password = ? WHERE name = 'wilson'", hash_password)
	c.IndentedJSON(http.StatusOK, new_info.Password)
}

func login(c *gin.Context) {
	var inputs player
	err := c.BindJSON(&inputs)
	if err != nil {
		log.Fatal("bindjson fail:", err)
	}
	log.Printf("login inputs: %v email: %s password: %s\n", inputs, inputs.Email, inputs.Password)
	
	var hash_password []byte
	if err := db.QueryRow("SELECT password FROM player WHERE email = ?", inputs.Email).Scan(&hash_password); err != nil {
		fmt.Errorf("QUERY player fail %v", err)
		c.IndentedJSON(http.StatusUnauthorized, "user login fail")
	} else {
		valid_password_result := bcrypt.CompareHashAndPassword(hash_password, []byte(inputs.Password))
		if valid_password_result == nil {
			type token_response struct {
				Token string
			}
			c.IndentedJSON(http.StatusOK, token_response{ Token: generateToken()})
		} else {
			c.IndentedJSON(http.StatusUnauthorized, "user login fail")
		}
	}
}

func main() {
  err := godotenv.Load() //load .env file
  if err != nil {
    log.Fatal("Error loading .env file", err)
  } else {
    is_success := getDatabaseHandle()
    if is_success {
      router := gin.Default()
      router.GET("/testauth", testAuth)
      router.POST("/account", changeAccountInfo)
      router.POST("/login", login)
      router.Run("localhost:8080")
    }
  }
}