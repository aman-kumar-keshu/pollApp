package web

import (
 	"polling-app/model"
	"log"
	"fmt"
	"net/http"
	"polling-app/database"
	"github.com/gin-gonic/gin"
	"strconv"
)

type App struct {
	d        database.PostgresDB
}


func NewApp(d database.PostgresDB)  {
	app := App{
		d:  d,
	} 
	app.d.Migrate()

	server := gin.New()
	server.Use(cors())
	publicAPI := server.Group("/")
	
	publicAPI.GET("/ping", pingPong)
	publicAPI.GET("/polls", app.GetPolls)
	publicAPI.GET("/poll/:id",app.GetPoll)
	publicAPI.PUT("/poll/:id", app.UpdatePoll)
	publicAPI.POST("/poll", app.createPoll)
	publicAPI.DELETE("/poll/:id", app.deletePoll)
	// publicAPI.POST("/login", app.LoginUser)
	// publicAPI.POST("/signup", routes.SignUpUser)
	log.Println("Web server is available on port 8080")
	server.Run(":8080")
	
}

func cors() gin.HandlerFunc{
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

func pingPong(c *gin.Context) {
	data := "Ping Successful"
	c.JSON(http.StatusOK, gin.H{
		data:data,
	})
}


func (a *App) LoginUser(w http.ResponseWriter, r *http.Request)  {
	// 
}
func (a *App) SignUpUser(w http.ResponseWriter, r *http.Request)  {
	// 
}

func (a *App) GetPolls(c *gin.Context) {

	log.Println("Fetching polls from DB")
 
	polls, err := a.d.GetPolls()
	log.Println("We got poll", polls)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"polls": polls,
	})
}

func (a *App) createPoll(c *gin.Context) {
	log.Println("Create POll")
	var poll model.Poll
	if err:= c.ShouldBindJSON(&poll); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	log.Println("New Poll : ",poll)
	err := a.d.CreatePoll(poll.Name, poll.Topic, poll.Src)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}  
	c.Status(http.StatusCreated)

}
func (a *App) UpdatePoll(c *gin.Context) {
	log.Println("updating polls")

	var poll model.Poll
	id,_ := strconv.Atoi(c.Param("id"))
	log.Println("id",id)
	if err:= c.ShouldBindJSON(&poll); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
 	log.Println("update called",poll)
	err := a.d.UpdatePoll(id, poll.Name, poll.Upvotes, poll.Downvotes)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}  
	c.Status(http.StatusOK)
}
func (a *App) GetPoll(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	log.Println("Get Poll Info id:", id);

	poll, err := a.d.GetPoll(id)
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
func (a *App) deletePoll (c *gin.Context) {
	log.Println("Deleting the poll");
	id,_ := strconv.Atoi(c.Param("id"))
	poll,err := a.d.GetPoll(id)
	log.Println("Poll to be deleted", poll)
	if poll.ID != id {
		c.JSON(http.StatusNotFound, gin.H{
			"message": fmt.Sprintf("Id not Found %d!", id),
		})
		return
	}
	err= a.d.DeletePoll(id);
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}  
	c.JSON(http.StatusOK, gin.H{
		"Deleted Poll Id": id,
	})
}