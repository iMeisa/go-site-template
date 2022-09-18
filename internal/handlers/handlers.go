package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/iMeisa/errortrace"
	"github.com/iMeisa/go-site-template/internal/config"
	"github.com/iMeisa/go-site-template/internal/dbDriver"
	"github.com/iMeisa/go-site-template/internal/models"
	"github.com/iMeisa/go-site-template/internal/render"
	"github.com/iMeisa/go-site-template/internal/repository"
	"github.com/iMeisa/go-site-template/internal/repository/dbrepo"
	"log"
	"net/http"
)

// Main handlers file

// Repo the repository used by the handlers
var Repo *Repository

// Template folder consts
const (
	PublicDir = "public"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *dbDriver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewDBRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Public handler directs you to any public page based on the url
func (m *Repository) Public(w http.ResponseWriter, r *http.Request) {
	page := chi.URLParam(r, "page")

	templatePage := fmt.Sprintf("%s.page.tmpl", page)
	render.Template(w, r, PublicDir, templatePage, &models.TemplateData{})
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, PublicDir, "home.page.tmpl", &models.TemplateData{})
}

func writeResp(w http.ResponseWriter, resp models.JsonResponse) {
	respJSON, err := json.Marshal(resp)
	if err != nil {
		trace := errortrace.NewTrace(err)
		trace.Read()
		log.Println("Error marshalling response struct to json")
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(respJSON)
}
