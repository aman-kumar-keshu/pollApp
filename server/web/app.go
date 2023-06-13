package web

import (
	"encoding/json"
	"log"
	"net/http"
	"polling-app/db"
	"github.com/go-chi/chi"
)

type App struct {
	d        db.DB
	handlers map[string]http.HandlerFunc
}

func NewApp(d db.DB, cors bool) error {
	app := App{
		d:        d,
		handlers: make(map[string]http.HandlerFunc),
	}
	techHandler := app.GetTechnologies
	if !cors {
		techHandler = disableCors(techHandler)
	}
	r:= chi.NewRouter()
	r.Get("/ping", pingPong)
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

func (a *App) GetTechnologies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	technologies, err := a.d.GetTechnologies()
	if err != nil {
		sendErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	err = json.NewEncoder(w).Encode(technologies)
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
