package server

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	t "text/template"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	thermotempl, err := GetTemplates()
	if err != nil {
		log.Println(err)
	}
	var examples []renderedTemplates
	for _, templ := range thermotempl {
		renderedThermoTemplate := renderThermoTemplate(&templ, false)
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
			GeneratePDF(string(renderThermoTemplate(&templ, false)), w)
		} else if r.PostForm.Get("action") == "print" {
			copies, err := strconv.Atoi(r.PostForm.Get("copies"))
			if err != nil {
				copies = 1
			}
			Print(string(renderThermoTemplate(&templ, false)), copies)
			url, _ := router.Get("detail").URL("name", templ.Name)
			http.Redirect(w, r, url.String(), 301)
		} else {
			templates["detail"].Execute(w, struct {
				Template ThermoTemplate
				Rendered template.HTML
			}{templ, renderThermoTemplate(&templ, true)})
		}

	} else {
		errorHandler(w, r)
	}

}

// helper functions

func renderThermoTemplate(templ *ThermoTemplate, preview bool) template.HTML {
	rendTempl, err := t.New(templ.Name).Parse(templ.Template)
	if err != nil {
		log.Println(err)
		return template.HTML("")
	}
	var b strings.Builder
	if preview {
		rendTempl.Execute(&b, highlightVars(&templ.Variables))
	} else {
		rendTempl.Execute(&b, templ.Variables)
	}
	return template.HTML(b.String())
}

func highlightVars(vars *map[string]string) map[string]string {
	ret := make(map[string]string)
	for key, value := range *vars {
		ret[key] = fmt.Sprintf("<strong>%s</strong>", value)
	}
	return ret
}
