package main

import (
	"net/http"
	"os"
	"strconv"

	gonnect "github.com/craftamap/atlas-gonnect"
	"github.com/craftamap/atlas-gonnect/middleware"
	gonnectRoutes "github.com/craftamap/atlas-gonnect/routes"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.Use(func(handler http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, handler)
	})
	router.Use(handlers.ProxyHeaders)

	configReader, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	descriptorReader, err := os.Open("atlassian-connect.json")
	if err != nil {
		panic(err)
	}

	addon, err := gonnect.NewAddon(configReader, descriptorReader)
	if err != nil {
		panic(err)
	}

	router.Use(middleware.NewRequestMiddleware(addon, make(map[string]string)))

	gonnectRoutes.RegisterRoutes(addon, router)
	RegisterRoutes(router, addon)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("public")))
	http.ListenAndServe(":"+strconv.Itoa(addon.Config.Port), router)
}
