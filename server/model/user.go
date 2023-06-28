package model

import (

	"github.com/golang-jwt/jwt/v5"
)


type User struct {
	ID       	int    		`json:"id"`
	Name 	 	string 		`json:"name"`
	Email 		string		`json:"email" binding:"required"`
	Password 	string 		`json:"password"  binding:"required"`
} 

type AuthRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type AuthResponse struct {
	JwtToken string `json:"jwtToken"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
