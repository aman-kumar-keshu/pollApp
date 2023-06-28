package routes

import (
	"github.com/gin-gonic/gin"
	"polling-app/handlers"
	"log"
)



func InitialiseRoutes(h *handlers.Handler) {
	server := gin.New()
	server.Use(handlers.Cors())
	publicAPI := server.Group("/")
	
	publicAPI.GET("/ping", h.PingPong)
	publicAPI.GET("/polls", h.GetPolls)
	publicAPI.GET("/poll/:id",h.GetPoll)
	publicAPI.PUT("/poll/:id", h.UpdatePoll)
	publicAPI.POST("/poll", h.CreatePoll)
	publicAPI.DELETE("/poll/:id", h.DeletePoll)
	publicAPI.POST("users/login", h.LoginUser)
	publicAPI.POST("users/signup", h.SignUpUser)
	log.Println("Web server is available on port 8080")
	server.Run(":8080")
}



