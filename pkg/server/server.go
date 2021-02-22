package server

import (
	"os"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"embed"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type renderedTemplates struct {
	Name      string
	Template  template.HTML
	DetailURL string
}

//go:embed site_templates/*.html
var pages embed.FS

var templates map[string]*template.Template

var router *mux.Router

var sessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

// Run starts the server on the specified port.
func Run(port int, static *embed.FS) {
	initializeHTMLTemplates([]string{"template_grid", "404", "detail"})

	// initialize handle functions
	router = mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(static))))
	router.HandleFunc("/template/{name}", detailHandler).Name("detail")
	router.PathPrefix("/media/").Handler(http.StripPrefix("/media/", http.FileServer(http.Dir("media"))))
	router.NotFoundHandler = http.HandlerFunc(errorHandler)
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + strconv.Itoa(port)}

	log.Fatal(srv.ListenAndServe())
}

func initializeHTMLTemplates(templateList []string) {
	templates = make(map[string]*template.Template)
	baseTemplateStr, err := pages.ReadFile("site_templates/base.html")
	if err != nil {
		panic(err)
	}
	notifyTmpl, err := pages.ReadFile("site_templates/notify.html")
	if err != nil {
		panic(err)
	}
	baseTemplate, _ := template.New("base").Parse(string(baseTemplateStr))
	baseTemplate, _ = baseTemplate.Parse(string(notifyTmpl))
	for _, t := range templateList {
		templateStr, err := pages.ReadFile("site_templates/"+ t + ".html")
		if err != nil {
			panic(err)
		}
		newTemplate, _ := baseTemplate.Clone()
		templates[t], _ = newTemplate.Parse(string(templateStr))
	}

}
