package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	gonnect "github.com/craftamap/atlas-gonnect"
	"github.com/craftamap/atlas-gonnect/hostrequest"
	"github.com/craftamap/atlas-gonnect/middleware"
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
	templates.ExecuteTemplate(w, "hello-world.html.tmpl", map[string]interface{}{
		"hostScriptUrl":     r.Context().Value("hostScriptUrl"),
		"hostStylesheetUrl": r.Context().Value("hostStylesheetUrl"),
		"localBaseUrl":      r.Context().Value("localBaseUrl"),
		"hostBaseUrl":       r.Context().Value("hostBaseUrl"),
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

	templates.ExecuteTemplate(w, "macro-view.html.tmpl", map[string]interface{}{
		"sport":             sport,
		"hostScriptUrl":     r.Context().Value("hostScriptUrl"),
		"hostStylesheetUrl": r.Context().Value("hostStylesheetUrl"),
		"localBaseUrl":      r.Context().Value("localBaseUrl"),
		"hostBaseUrl":       r.Context().Value("hostBaseUrl"),
	})
}

func RenderMacroPage(w http.ResponseWriter, r *http.Request) {
	RenderMacro(w, r)
}

func MacroHandleFunc(w http.ResponseWriter, r *http.Request) {
	RenderMacro(w, r)
}

func MacroEditorFunc(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "macro-editor.html.tmpl", map[string]interface{}{
		"sports":        sports,
		"localBaseUrl":  r.Context().Value("localBaseUrl"),
		"hostBaseUrl":   r.Context().Value("hostBaseUrl"),
		"hostScriptUrl": r.Context().Value("hostScriptUrl"),
	})
}

func RegisterRoutes(router *mux.Router, addon *gonnect.Addon) error {
	tmpl, err := template.ParseGlob("templates/*.tmpl")
	templates = tmpl
	if err != nil {
		panic(err)
	}
	fmt.Println(templates.DefinedTemplates())

	router.Handle("/hello-world", middleware.NewAuthenticationMiddleware(addon, false)(http.HandlerFunc(helloWorldHandleFunc)))
	router.Handle("/macro", middleware.NewAuthenticationMiddleware(addon, false)(http.HandlerFunc(MacroHandleFunc))).Methods("POST")
	router.Handle("/macro-page", middleware.NewAuthenticationMiddleware(addon, false)(http.HandlerFunc(RenderMacroPage)))
	router.Handle("/editor", middleware.NewAuthenticationMiddleware(addon, false)(http.HandlerFunc(MacroEditorFunc)))
	router.Handle("/asUser", middleware.NewAuthenticationMiddleware(addon, false)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "asUser.html.tmpl", map[string]interface{}{
			"localBaseUrl":  r.Context().Value("localBaseUrl"),
			"hostBaseUrl":   r.Context().Value("hostBaseUrl"),
			"hostScriptUrl": r.Context().Value("hostScriptUrl"),
		})

		if err != nil {
			panic(err)
		}
	})))
	router.Handle("/api/asUser", middleware.NewTokenMiddleware(addon)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpClient, _ := hostrequest.FromRequest(r)
		request, _ := http.NewRequest("GET", "/rest/api/content", http.NoBody)
		_, err := httpClient.AsUser(request, r.Context().Value("userAccountId").(string))
		if err == nil {
			response, _ := http.DefaultClient.Do(request)
			rBody, _ := ioutil.ReadAll(response.Body)
			w.Write(rBody)
			w.Header().Set("Content-Type", "application/json")
		}
	})))
	return nil
}
