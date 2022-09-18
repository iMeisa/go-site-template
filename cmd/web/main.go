package main

import (
	"database/sql"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/iMeisa/errortrace"
	"github.com/iMeisa/go-site-template/internal/config"
	"github.com/iMeisa/go-site-template/internal/dbDriver"
	"github.com/iMeisa/go-site-template/internal/handlers"
	"github.com/iMeisa/go-site-template/internal/render"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// Load env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// WebApp settings
	app.Prod = os.Getenv("ENV") == "prod"

	session = scs.New()
	session.Lifetime = 24 * time.Hour * 3 // 3 days
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.Prod

	app.Session = session

	// Compile TypeScript
	log.Println("Compiling TypeScript...")
	go func() {
		err := exec.Command("tsc").Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Connect to db
	log.Println("Connecting to DB...")
	db, trace := dbDriver.ConnectSQL(os.Getenv("DBNAME"))
	if trace.HasError() {
		trace.Read()
		log.Fatal()
	}
	//Close connection
	defer func(SQL *sql.DB) {
		err := SQL.Close()
		if err != nil {
			trace = errortrace.NewTrace(err)
			trace.Read()
		}
	}(db.SQL)
	log.Println("Connected to DB")

	// Templates
	tc, trace := render.CreateTemplateCache()
	if trace.HasError() {
		trace.Read()
		log.Fatal("Cannot create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)

	// Serve
	fmt.Println(fmt.Sprintf("Starting %s application on port %s", os.Getenv("ENV"), os.Getenv("SITE_PORT")))

	srv := &http.Server{
		Addr:    os.Getenv("SITE_PORT"),
		Handler: routes(&app),
	}

	err := srv.ListenAndServe()
	if err != nil {
		trace = errortrace.NewTrace(err)
		trace.Read()
		log.Fatal(err)
	}
}
