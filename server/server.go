package main

import (
	"polling-app/database"
	"polling-app/web"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	d := initDB()
	// CORS is enabled only in prod profile
	web.NewApp(database.NewDB(d))
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
	pass := ""
	user := "amankumarkeshu"
	if os.Getenv("profile") == "prod" {
		host = "db"
		pass = os.Getenv("db_pass")
	}
	return "postgresql://" + host + ":5433/poll" +
		"?user=" + user + "&sslmode=disable&password=" + pass
}