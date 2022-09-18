package render

import (
	"bytes"
	"fmt"
	"github.com/iMeisa/errortrace"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/iMeisa/go-site-template/internal/config"
	"github.com/iMeisa/go-site-template/internal/funcs"
	"github.com/iMeisa/go-site-template/internal/models"
)

var functions = funcs.Functions

var app *config.AppConfig

// NewRenderer sets the config for the template package
func NewRenderer(a *config.AppConfig) {
	app = a
}

func Template(w http.ResponseWriter, r *http.Request, dir, tmpl string, data *models.TemplateData) {

	var templateCache map[string]*template.Template

	if app.UseCache {
		// Get the page cache from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	pageName := fmt.Sprintf("%v/%v", dir, tmpl)
	page, ok := templateCache[pageName]
	if !ok {
		log.Println("page is not ok")
		http.Error(w, "Page not found", http.StatusNotFound)
		return
	}

	buf := new(bytes.Buffer)

	data.AddDefaultData(r, app)

	err := page.Execute(buf, data)
	if err != nil {
		trace := errortrace.NewTrace(err)
		trace.Read()
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		trace := errortrace.NewTrace(err)
		trace.Read()
	}
}

func CreateTemplateCache() (map[string]*template.Template, errortrace.ErrorTrace) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*/*.page.tmpl")
	if err != nil {
		return myCache, errortrace.NewTrace(err)
	}

	for _, page := range pages {
		name := filepath.Base(page)
		pageDir := strings.Split(filepath.Dir(page), string(os.PathSeparator))[1]
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, errortrace.NewTrace(err)
		}

		matches, err := filepath.Glob("./templates/layouts/*.layout.tmpl")
		if err != nil {
			return myCache, errortrace.NewTrace(err)
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/layouts/*.layout.tmpl")
			if err != nil {
				return myCache, errortrace.NewTrace(err)
			}
		}

		tmplName := fmt.Sprintf("%v/%v", pageDir, name)
		myCache[tmplName] = ts
	}

	return myCache, errortrace.NilTrace()
}
