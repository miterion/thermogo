package server

import (
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	thermotempl, err := GetTemplates()
	if err != nil {
		log.Println(err)
	}
	var examples []renderedTemplates
	for _, templ := range thermotempl {
		renderedThermoTemplate := templ.renderThermoTemplate(false)
		url, _ := router.Get("detail").URL("name", templ.Name)
		examples = append(examples, renderedTemplates{
			templ.Name, renderedThermoTemplate, url.String()})
	}
	sort.Slice(examples, func(i, j int) bool { return examples[i].Name < examples[j].Name })
	templates["template_grid"].Execute(w, examples)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	templates["404"].Execute(w, nil)
	return
}

func detailHandler(w http.ResponseWriter, r *http.Request) {
	thermotempl, err := GetTemplates()
	if err != nil {
		log.Println(err)
		errorHandler(w, r)
	}
	vars := mux.Vars(r)

	if templ, ok := thermotempl[vars["name"]]; ok {
		if r.Method == "POST" {
			r.ParseForm()
			newValues := make(map[string]string)
			for key, value := range templ.Variables {
				tempVal := r.PostForm.Get(key)
				if tempVal == "" {
					newValues[key] = value
				} else {
					newValues[key] = tempVal
				}
			}
			templ.Variables = newValues
		}
		if r.PostForm.Get("action") == "pdf" {
			w.Header().Set("Content-type", "application/pdf")
			GeneratePDF(string(templ.renderThermoTemplate(false)), w)
		} else if r.PostForm.Get("action") == "print" {
			copies, err := strconv.Atoi(r.PostForm.Get("copies"))
			if err != nil {
				copies = 1
			}
			Print(string(templ.renderThermoTemplate(false)), copies)
		}
		templates["detail"].Execute(w, struct {
		Template ThermoTemplate
		Rendered template.HTML
		}{templ, templ.renderThermoTemplate(true)})

	} else {
		errorHandler(w, r)
	}

}
