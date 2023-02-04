package util

import (
  "github.com/golang-jwt/jwt"
  "time"
  "reflect"
  "fmt"
)

func GenerateToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(10 * time.Minute).Unix(),
		//"player": player id
	})
	tokenString, _ := token.SignedString([]byte("lakgfnlawng"))
	return tokenString
}

func VerifyToken(request_token string) (bool, error) {
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