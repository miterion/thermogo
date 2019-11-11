package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"
	t "text/template"
)

// ThermoTemplate is the basic struct for saving templates.
type ThermoTemplate struct {
	Name      string
	Variables map[string]string // the used variables with an example value
	Template  string            // the cleaned template without the front matter
}

// GetTemplates returns a map containing the templates in the template directory.
func GetTemplates() (map[string]ThermoTemplate, error) {
	_, err := os.Stat("templates")
	if err != nil {
		log.Println("No template directory, creating one.")
		err = os.Mkdir("templates", 0755)
		if err != nil {
			return nil, err
		}
	}
	return loadTemplates()
}

func loadTemplates() (map[string]ThermoTemplate, error) {
	templates, err := ioutil.ReadDir("templates")
	if err != nil {
		panic(err)
	}
	if len(templates) == 0 {
		return nil, errors.New("No templates found")
	}
	var parsedTemplates = make(map[string]ThermoTemplate)
	for _, file := range templates {
		str, err := ioutil.ReadFile("templates/" + file.Name())
		if err != nil {
			return nil, err
		}
		tmplName := strings.Split(file.Name(), ".")[0]
		template, err := parseMetadata(tmplName, string(str))
		if err != nil {
			return nil, err
		}
		parsedTemplates[tmplName] = template
	}
	return parsedTemplates, nil
}

func parseMetadata(name, templatestring string) (ThermoTemplate, error) {
	split := strings.Split(templatestring, "---")
	if len(split) != 3 {
		return ThermoTemplate{}, fmt.Errorf("Wrongly formatted template, does not contain a block surrounded by two ---. Here is the offending template:\n%s", templatestring)
	}
	// check whether we got valid json
	if !json.Valid([]byte(split[1])) {
		return ThermoTemplate{}, fmt.Errorf("Invalid JSON, got:\n%s", split[1])
	}
	// load and parse json
	var jsonOutput map[string]interface{}
	json.Unmarshal([]byte(split[1]), &jsonOutput)
	variables := make(map[string]string)
	for key := range jsonOutput {
		variables[key] = jsonOutput[key].(string)
	}
	return ThermoTemplate{
		name, variables, split[2],
	}, nil
}

func (templ ThermoTemplate) renderThermoTemplate(preview bool) template.HTML {
	rendTempl, err := t.New(templ.Name).Parse(templ.Template)
	if err != nil {
		log.Println(err)
		return template.HTML("")
	}
	var b strings.Builder
	if preview {
		renderedVariables := vueTemplate(&templ.Variables)
		renderedVariables["media"] = "..//media/"
		rendTempl.Execute(&b, renderedVariables)
	} else {
		templ.Variables["media"] = "http://localhost:3000//media/"
		rendTempl.Execute(&b, templ.Variables)
	}
	return template.HTML(b.String())
}

// helper functions

func vueTemplate(vars *map[string]string) map[string]string {
	ret := make(map[string]string)
	for key := range *vars {
		ret[key] = fmt.Sprintf("<span class='previewtext' id='prev-%s'>$%s$</span>", key, key)
	}
	return ret
}
