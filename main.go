package main

import (
  "github.com/gin-gonic/gin"
  "log"
  "github.com/joho/godotenv"
  "online_chess/model"
  "online_chess/controller"
)

func main() {
  err := godotenv.Load() //load .env file
  if err != nil {
    log.Fatal("Error loading .env file", err)
  } else {
    is_success := model.GetDatabaseHandle()
    if is_success {
      router := gin.Default()
      router.GET("/testauth", controller.TestAuth)
      router.POST("/account", controller.ChangeAccountInfo)
      router.POST("/login", controller.Login)
      router.Run("localhost:8080")
    }
  }
}