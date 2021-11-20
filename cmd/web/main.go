package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/doh-halle/ntarikoon-park/pkg/config"
	"github.com/doh-halle/ntarikoon-park/pkg/handlers"
	"github.com/doh-halle/ntarikoon-park/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":9000"

var app config.AppConfig
var session *scs.SessionManager

//Main is the main application function
func main() {

	//Change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template Cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	//http.HandleFunc("/", handlers.Repo.Home)
	//http.HandleFunc("/about", handlers.Repo.About)

	fmt.Println(fmt.Printf("Starting Ntarikon Park Web Application on port %s", portNumber))
	//_ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
