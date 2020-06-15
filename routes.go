package main

import (
	"html/template"
	"net/http"

	gonnect "github.com/craftamap/atlas-gonnect"
	"github.com/craftamap/atlas-gonnect/middleware"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

var templates *template.Template

var iconPath = "images/icons/sport/"
var sports = []struct {
	Id   string
	Name string
	Icon string
}{
	{Id: "nfl", Name: "American Football", Icon: iconPath + "american_football.png"},
	{Id: "baseball", Name: "Baseball", Icon: iconPath + "baseball.png"},
	{Id: "basketball", Name: "Basketball", Icon: iconPath + "basketball.png"},
	{Id: "football", Name: "Football", Icon: iconPath + "football.png"},
	{Id: "golf", Name: "Golf", Icon: iconPath + "golf.png"},
	{Id: "tennis", Name: "Tennis", Icon: iconPath + "tennis.png"},
}

func helloWorldHandleFunc(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "hello-world.html", map[string]interface{}{
		"hostScriptUrl":     context.Get(r, "hostScriptUrl"),
		"hostStylesheetUrl": context.Get(r, "hostStylesheetUrl"),
		"localBaseUrl":      context.Get(r, "localBaseUrl"),
		"hostBaseUrl":       context.Get(r, "hostBaseUrl"),
	})
}

func RenderMacro(w http.ResponseWriter, r *http.Request) {
	sportName := r.FormValue("sport")
	sport := sports[0]
	for _, value := range sports {
		if value.Id == sportName {
			sport = value
			break
		}
	}

	templates.ExecuteTemplate(w, "macro-view.html", map[string]interface{}{
		"sport":             sport,
		"hostScriptUrl":     context.Get(r, "hostScriptUrl"),
		"hostStylesheetUrl": context.Get(r, "hostStylesheetUrl"),
		"localBaseUrl":      context.Get(r, "localBaseUrl"),
		"hostBaseUrl":       context.Get(r, "hostBaseUrl"),
	})
}

func RenderMacroPage(w http.ResponseWriter, r *http.Request) {
	RenderMacro(w, r)
}

func MacroHandleFunc(w http.ResponseWriter, r *http.Request) {
	RenderMacro(w, r)
}

func MacroEditorFunc(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "macro-editor.html", map[string]interface{}{
		"sports":        sports,
		"localBaseUrl":  context.Get(r, "localBaseUrl"),
		"hostBaseUrl":   context.Get(r, "hostBaseUrl"),
		"hostScriptUrl": context.Get(r, "hostScriptUrl"),
	})
}

func RegisterRoutes(router *mux.Router, addon *gonnect.Addon) error {
	tmpl, err := template.ParseGlob("templates/*.html")
	templates = tmpl
	if err != nil {
		panic(err)
	}

	router.Handle("/hello-world", middleware.NewAuthenticationMiddleware(addon, false)(http.HandlerFunc(helloWorldHandleFunc)))
	router.Handle("/macro", middleware.NewAuthenticationMiddleware(addon, false)(http.HandlerFunc(MacroHandleFunc))).Methods("POST")
	router.Handle("/macro-page", middleware.NewAuthenticationMiddleware(addon, false)(http.HandlerFunc(RenderMacroPage)))
	router.Handle("/editor", middleware.NewAuthenticationMiddleware(addon, false)(http.HandlerFunc(MacroEditorFunc)))
	return nil
}
