package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/doh-halle/ntarikoon-park/internal/config"
	"github.com/doh-halle/ntarikoon-park/internal/driver"
	"github.com/doh-halle/ntarikoon-park/internal/handlers"
	"github.com/doh-halle/ntarikoon-park/internal/helpers"
	"github.com/doh-halle/ntarikoon-park/internal/models"
	"github.com/doh-halle/ntarikoon-park/internal/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":9000"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

//Main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	defer close(app.MailChan)

	fmt.Println("Starting mail listener...")
	listenForMail()

	fmt.Println(fmt.Printf("Starting Ntarikon Park Web Application on port %s", portNumber))
	//_ = http.ListenAndServe(portNumber, nil)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// what I am going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Apartment{})
	gob.Register(models.Restriction{})
	gob.Register(models.ApartmentRestriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	//Change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//Connect to Database
	log.Println("connecting to np database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=ntarikoon_park user=hallecraft password=manager123")
	if err != nil {
		log.Fatal("cannot connect to np database...")
	}
	log.Println("Connected to np database...")

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Cannot create template Cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
