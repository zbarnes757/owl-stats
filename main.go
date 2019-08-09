package main

import (
	"owl-stats/app"
	"owl-stats/controllers"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

func main() {
	router := mux.NewRouter()
	router.Use(app.JwtAuthentication)

	// add routes
	router.HandleFunc("/api/v1/user/new", controllers.CreateAccount).Methods("POST")

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(router)

	app.StartScheduledProcesses()

	n.Run()
}
