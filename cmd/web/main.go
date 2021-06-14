package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/taldrori/bookings/internal/config"
	"github.com/taldrori/bookings/internal/handlers"
	"github.com/taldrori/bookings/internal/models"
	"github.com/taldrori/bookings/internal/render"
)

const portNumber = ":8080"

var app config.Appconfig
var session *scs.SessionManager

func main() {
	gob.Register(models.Reservation{})
	//change to true in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cant create template cache:", err)
	}

	app.TemplateCache = tc
	app.UseChache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)

	fmt.Printf("Starting at port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
