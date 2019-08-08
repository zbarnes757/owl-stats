package main

import (
	"owl-stats/app"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	router := mux.NewRouter()
	router.Use(app.JwtAuthentication)

	// add routes
	// TODO: think of routes to add

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(router)

	n.Run()
}
