package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"meisa_xyz/pkg/config"
	"meisa_xyz/pkg/models"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(data *models.TemplateData) *models.TemplateData {

	return data
}

func Template(w http.ResponseWriter, tmpl string, data *models.TemplateData) {

	var templateCache map[string]*template.Template

	if app.UseCache {
		// Get the page cache from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	page, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("page is not ok")
	}

	buf := new(bytes.Buffer)

	data = AddDefaultData(data)

	_ = page.Execute(buf, data)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing page to browser:", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}

	return myCache, nil
}
