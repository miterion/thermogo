package server

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gobuffalo/packr/v2"
	"github.com/gorilla/mux"
)

type renderedTemplates struct {
	Name      string
	Template  template.HTML
	DetailUrl string
}

var pagesBox *packr.Box

var templates map[string]*template.Template

var router *mux.Router

// Run starts the server on the specified port.
func Run(port int) {
	// initialize content boxes
	bootstrapBox := packr.New("bootstrap", "../../node_modules/bootstrap/dist")
	jqueryBox := packr.New("jquery", "../../node_modules/jquery/dist")
	vueBox := packr.New("vue", "../../node_modules/vue/dist")
	pagesBox = packr.New("pages", "./site_templates")

	initializeHTMLTemplates([]string{"template_grid", "404", "detail"})

	// initialize handle functions
	router = mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.PathPrefix("/static/bootstrap/").Handler(http.StripPrefix("/static/bootstrap/", http.FileServer(bootstrapBox)))
	router.PathPrefix("/static/jquery/").Handler(http.StripPrefix("/static/jquery/", http.FileServer(jqueryBox)))
	router.PathPrefix("/static/vue/").Handler(http.StripPrefix("/static/vue/", http.FileServer(vueBox)))
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
	baseTemplateStr, err := pagesBox.FindString("base.html")
	if err != nil {
		panic(err)
	}
	baseTemplate, _ := template.New("base").Parse(baseTemplateStr)
	for _, t := range templateList {
		templateStr, err := pagesBox.FindString(t + ".html")
		if err != nil {
			panic(err)
		}
		newTemplate, _ := baseTemplate.Clone()
		templates[t], _ = newTemplate.Parse(templateStr)
	}

}
