package web

import (
 	"polling-app/model"
	"log"
	"net/http"
	"polling-app/database"
	"github.com/gin-gonic/gin"
	"strconv"
)

type App struct {
	d        database.DB
}


func NewApp(d database.DB)  {
	app := App{
		d:  d,
	} 
	app.d.Migrate()

	server := gin.New()
	server.Use(cors())
	publicAPI := server.Group("/")
	
	publicAPI.GET("/ping", pingPong)
	publicAPI.GET("/polls", app.GetPolls)
	publicAPI.PUT("/poll/:id", app.UpdatePoll)
	// r.Get("/", http.FileServer(http.Dir("/webapp")).ServeHTTP)
	// r.Post("/login", app.LoginUser)
	// r.Post("/signup", routes.SignUpUser)
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
