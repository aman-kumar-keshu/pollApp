package main

import (
	"polling-app/handlers"
	"polling-app/database"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
	"polling-app/routes"
)



func main() {
	d := initDB()
	
	handler:= NewHandler(database.NewDB(d))
	routes.InitialiseRoutes(handler)

}

func NewHandler(d database.PostgresDB) *handlers.Handler {
	handler := &handlers.Handler{
		D:d,
	} 
	handler.D.Migrate()
	return handler
	
}

func initDB() *sql.DB  {

	db, err := sql.Open("postgres", filePath())
	if err != nil {
		log.Fatal(err)
	}
	
	if db == nil {
		log.Fatal("db nil")
	}
	log.Println("Successfully connected to the db")
	return db
}

func filePath() string {
	host := "localhost"
	pass := "Aman@1999"
	user := "amankumarkeshu"
	if os.Getenv("profile") == "prod" {
		host = "db"
		pass = os.Getenv("db_pass")
	}
	return "postgresql://" + host + ":5433/poll" +
		"?user=" + user + "&sslmode=disable&password=" + pass
}