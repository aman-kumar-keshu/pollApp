package web

import (
	"polling-app/model"
	"encoding/json"
	"log"
	"net/http"
	"polling-app/database"
	"github.com/go-chi/chi"
	"strconv"
)

type App struct {
	d        database.DB
	handlers map[string]http.HandlerFunc
}

func NewApp(d database.DB, cors bool) error {
	app := App{
		d:        d,
	}
	pollHandler := app.GetPolls
	if !cors {
		pollHandler = disableCors(pollHandler)
	}
	r:= chi.NewRouter()
	r.Get("/ping", pingPong)
	r.Get("/polls", app.GetPolls)
	r.Put("/poll/{id}", app.UpdatePoll)
	r.Get("/", http.FileServer(http.Dir("/webapp")).ServeHTTP)
	r.Post("/login", app.LoginUser)
	r.Post("/signup", app.SignUpUser)

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


func (a *App) LoginUser(w http.ResponseWriter, r *http.Request)  {
	// 
}
func (a *App) SignUpUser(w http.ResponseWriter, r *http.Request)  {
	// 
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


func (a *App) UpdatePoll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var poll model.Poll
	index,_ := strconv.Atoi(chi.URLParam(r, "id"))
	json.NewDecoder(r.Body).Decode(&poll)

	// index, _ := strconv.Atoi(c.Param("index"))

	id, err := a.d.UpdatePoll(index, poll.Name, poll.Upvotes, poll.Downvotes)

	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
	} else {
		 
		err := json.NewEncoder(w).Encode(id)
		if err != nil {
			sendErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		log.Println("Updated poll", id)

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
