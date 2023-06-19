package web

import (
	"encoding/json"
	"log"
	"net/http"
	"polling-app/database"
	"github.com/go-chi/chi"
)

type App struct {
	d        database.DB
	handlers map[string]http.HandlerFunc
}

func NewApp(d database.DB, cors bool) error {
	app := App{
		d:        d,
		handlers: make(map[string]http.HandlerFunc),
	}
	pollHandler := app.GetPolls
	if !cors {
		pollHandler = disableCors(pollHandler)
	}
	r:= chi.NewRouter()
	r.Get("/ping", pingPong)
	r.Get("/polls", app.GetPolls)
	r.Get("/", http.FileServer(http.Dir("/webapp")).ServeHTTP)
	log.Println("Web server is available on port 8080")

	err:= http.ListenAndServe(":8080",r)
	return err
	
}

func pingPong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data := "Ping Successful"
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	log.Println("Ping");
}


func (a *App) GetPolls(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	polls, err := a.d.GetPolls()
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(polls)
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
	}
}

func sendErr(w http.ResponseWriter, code int, message string) {
	resp, _ := json.Marshal(map[string]string{"error": message})
	http.Error(w, string(resp), code)
}

// Needed in order to disable CORS for local development
func disableCors(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		h(w, r)
	}
}
