package main

import (
	"github.com/kenmobility/lemtest/server"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/user", server.CreateUser)
	r.GET("/users", server.GetAllUsers)
	r.POST("/transaction", server.CreateTransaction)

	go server.UserVerificationProcessor()
	go server.TransactionProcessor()

	r.Run()
}
