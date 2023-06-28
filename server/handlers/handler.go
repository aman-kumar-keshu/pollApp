
package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"polling-app/model"
	"log"
	"fmt"
	"polling-app/database"
	"strconv"
	"github.com/google/uuid"

	 )

type Handler struct {
	D       database.PostgresDB
}

func (h *Handler) PingPong(c *gin.Context) {
	data := "Ping Successful"
	c.JSON(http.StatusOK, gin.H{
		data:data,
	})
}

func (h *Handler) SignUpUser(c *gin.Context)  {
	log.Println("Create User")

	var user model.User
	if err:= c.ShouldBindJSON(&user); err != nil {
		log.Println("Failed to parse user", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	hashedPassword, err1 := GenerateHashedPassword(user.Password)
	if err1 != nil {
		log.Println("Error in generating hashed Password", err1)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	userCreated, err := h.D.CreateUser(user.Name, user.Email, hashedPassword)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H {
			"User":  "user not created",
		})
		return
	}
	log.Println(" User Created : ",userCreated)

	token := generateToken()
	// save token to db
	err = h.D.SaveToken(token, userCreated.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"message":  "Token could not be created",
		})
		return
	}


	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (h *Handler) LoginUser(c *gin.Context) {
	user := model.User{}
	 if err:= c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	 }
	 dbUser,err := h.D.FetchUser(user.Email)
	 log.Println("Login User fetched from DB:", dbUser)
	 if dbUser.Email != user.Email {
		log.Println("User does not exists with this email")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User not Found with email",
		})
		return

	}

	 if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H {
			"message":  "Cannot fetch user",
		})
		return
	}


	 if CheckPasswordHash(user.Password, dbUser.Password) {
		token,err := h.D.FetchUserToken(dbUser.ID)
		if err != nil {
			log.Println("Cannot fetch Token from DB:")
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		log.Println("Token Fetched from DB:", token)

		 c.JSON(http.StatusAccepted, gin.H {
			"status": 200,
			"message": "Successfully Logged in",
			"token": token,
		})
	 } else {
		 c.JSON(http.StatusBadRequest, gin.H {
			"message": "Incorrect Password",
		})
	 }
}
func (h *Handler) GetPolls(c *gin.Context) {

	log.Println("Fetching polls from DB")
 
	polls, err := h.D.GetPolls()
	log.Println("We got poll", polls)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"polls": polls,
	})
}

func (h *Handler) CreatePoll(c *gin.Context) {
	log.Println("Create POll")
	var poll model.Poll
	if err:= c.ShouldBindJSON(&poll); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Println("New Poll : ",poll)
	err := h.D.CreatePoll(poll.Name, poll.Topic, poll.Src)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}  
	c.Status(http.StatusCreated)

}
func (h *Handler) UpdatePoll(c *gin.Context) {
	log.Println("updating polls")

	var poll model.Poll
	id,_ := strconv.Atoi(c.Param("id"))
	log.Println("id",id)
	if err:= c.ShouldBindJSON(&poll); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
 	log.Println("update called",poll)
	err := h.D.UpdatePoll(id, poll.Name, poll.Upvotes, poll.Downvotes)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}  
	c.Status(http.StatusOK)
}
func (h *Handler) GetPoll(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	log.Println("Get Poll Info id:", id);

	poll, err := h.D.GetPoll(id)
	if poll.ID != id {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Id not Found %d!", id),
		})
		return
	}
 	if err != nil {
		c.AbortWithError(http.StatusNotFound,err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Poll":poll,
	})

}
func (h *Handler) DeletePoll (c *gin.Context) {
	log.Println("Deleting the poll");
	id,_ := strconv.Atoi(c.Param("id"))
	poll,err := h.D.GetPoll(id)
	log.Println("Poll to be deleted", poll)
	if poll.ID != id {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Id not Found %d!", id),
		})
		return
	}
	err= h.D.DeletePoll(id);
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}  
	c.JSON(http.StatusOK, gin.H{
		"Deleted Poll Id": id,
	})
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Authorization")
	if c.Request.Method ==http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent);
		return
	}
	c.Next();
}
}

func generateToken() string{
	token := uuid.New()


	return token.String()
}