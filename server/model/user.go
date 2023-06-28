package model


type User struct {
	ID       	int    		`json:"id"`
	Name 	 	string 		`json:"name"`
	Email 		string		`json:"email" binding:"required"`
	Password 	string 		`json:"password"  binding:"required"`
} 

type TokenToUserId struct {
	ID       	int  
	UserId 		int  	
	Token 		string 
}