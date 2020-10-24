package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/google/martian/log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/gorilla/mux"
)

type App struct {
	DB *gorm.DB
}

type User struct {
	ID uint `json:"id"`
	Name string
	RoleID int
}

type Question struct {
	ID uint
	QuizID int
	Question string
}

type Response struct {
	ID uint
	QuestionID int
	UserID int
	Response string
}

type Quiz struct {
	ID uint
	Questions []Question
	UserID int
}

func (a *App) handler(w http.ResponseWriter, r *http.Request) {
	// Create a test Star.
	a.DB.Create(&Response{Response: "test"})

	// Read from DB.
	var res Response
	a.DB.First(&res, "name = ?", "test")

	// Write to HTTP response.
	w.WriteHeader(200)
	w.Write([]byte(star.Name))

	// Delete.
	a.DB.Delete(&star)
}

func main() {
	a := &App{}
	a.Initialize("sqlite3", "test.db")

	r := mux.NewRouter()

	r.HandleFunc("/create/quiz", a.CreateQuizHandler).Methods("POST")
	r.HandleFunc("/get/{quiz:.+}", a.ViewHandler).Methods("GET")

	r.HandleFunc("/stars", a.ListHandler).Methods("GET")
	r.HandleFunc("/stars/{name:.+}", a.UpdateHandler).Methods("PUT")
	r.HandleFunc("/stars/{name:.+}", a.DeleteHandler).Methods("DELETE")

	http.Handle("/", r)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
	defer a.DB.Close()
}

func (a *App) Initialize(dbDriver string, dbURI string) {
	db, err := gorm.Open(dbDriver, dbURI)
	if err != nil {
		panic("failed to connect database")
	}
	a.DB = db
}

func (a *App) CreateQuizHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the POST body to populate r.PostForm.
	if err := r.ParseForm(); err != nil {
		panic("failed in ParseForm() call")
	}

	// Create a new response from the request body.
	err := r.ParseForm()
	if err != nil {
		log.Errorf("error parsing form")
	}

	var questions []Question
	for i := 0; i <= len(r.Form); i++ {
		t := Question{
			Question: r.FormValue("question"),
		}
		questions = append(questions, t)
	}

	userId, err := strconv.Atoi(r.FormValue("userId"))
	if err != nil {
		log.Errorf("error converting userId")
	}
	quiz := &Quiz{
		Questions: questions,
		UserID: userId,
	}
	a.DB.Create(quiz)

	// Form the URL of the newly created response.
	u, err := url.Parse(fmt.Sprintf("/quiz/%s", quiz.ID))
	if err != nil {
		panic("failed to form new Star URL")
	}
	base, err := url.Parse(r.URL.String())
	if err != nil {
		panic("failed to parse request URL")
	}

	// Write to HTTP response.
	w.Header().Set("Location", base.ResolveReference(u).String())
	w.WriteHeader(201)
}

func (a *App) ListHandler(w http.ResponseWriter, r *http.Request) {
	var responses []Response

	// Select all responses and convert to JSON.
	a.DB.Find(&responses, "quiz_id")
	starsJSON, _ := json.Marshal(responses)

	// Write to HTTP response.
	w.WriteHeader(200)
	w.Write([]byte(starsJSON))
}

func (a *App) ViewHandler(w http.ResponseWriter, r *http.Request) {
	var quiz Quiz
	vars := mux.Vars(r)

	// Select the quiz with the given name, and convert to JSON.
	a.DB.First(&quiz, "id = ?", vars["quiz_id"])
	starJSON, _ := json.Marshal(quiz)

	// Write to HTTP response.
	w.WriteHeader(200)
	w.Write([]byte(starJSON))
}

func (a *App) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Parse the request body to populate r.PostForm.
	if err := r.ParseForm(); err != nil {
		panic("failed in ParseForm() call")
	}

	// Set new star values from the request body.
	star := &Star{
		Name: r.PostFormValue("name"),
		Description: r.PostFormValue("description"),
		URL: r.PostFormValue("url"),
	}

	// Update the star with the given name.
	a.DB.Model(&star).Where("name = ?", vars["name"]).Updates(&star)

	// Write to HTTP response.
	w.WriteHeader(204)
}

func (a *App) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// Delete the star with the given name.
	a.DB.Where("name = ?", vars["name"]).Delete(Star{})

	// Write to HTTP response.
	w.WriteHeader(204)
}