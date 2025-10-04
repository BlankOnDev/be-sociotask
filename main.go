package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/harundarat/be-socialtask/internal/app"
	"github.com/harundarat/be-socialtask/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "Go backend server port")

	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}
	defer app.DB.Close()

	r := routes.SetupRoutes(app)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("We are running on port %d\n", port)

	err = server.ListenAndServe()
	if err != nil {

		app.Logger.Fatal(err)
	}
}
