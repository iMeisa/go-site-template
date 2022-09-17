package models

import (
	"encoding/json"
	"fmt"
	"github.com/iMeisa/errortrace"
	"github.com/iMeisa/go-site-template/internal/config"
	"github.com/justinas/nosurf"
	"net/http"
)

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	User
	CSRFToken string
	Flash     string // Some message
	Warning   string
	Error     string
}

// AddDefaultData adds data for all templates
func (data *TemplateData) AddDefaultData(r *http.Request, a *config.AppConfig) {
	// String map
	//fmt.Println("Adding default data")
	stringMap := make(map[string]string)
	if data.StringMap != nil {
		stringMap = data.StringMap
	}

	data.StringMap = stringMap

	data.CSRFToken = nosurf.Token(r)

	// Separate types
	// User context to user type
	if user := a.Session.Get(r.Context(), "user"); user != nil {
		userString := fmt.Sprintf("%v", user)
		err := json.Unmarshal([]byte(userString), &data.User)

		if err != nil {
			trace := errortrace.NewTrace(err)
			trace.Read()
		}
	}

	if data.Data == nil {
		data.Data = make(map[string]interface{})
	}
}

// GetDefaultData returns a TemplateData object with default data
func GetDefaultData(r *http.Request, a *config.AppConfig) TemplateData {
	data := TemplateData{}
	data.AddDefaultData(r, a)

	return data
}
