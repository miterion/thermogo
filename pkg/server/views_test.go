package server

import (
	"html/template"
	"testing"
)

var testTemplate = ThermoTemplate{
	"Test",
	map[string]string{"testvar": "exampleval"},
	"<h1>Testing</h1><br>{{.testvar}}",
}

func TestHighlightVars(t *testing.T) {
	res := vueTemplate(&testTemplate.Variables)
	correct := map[string]string{"testvar": "$testvar$"}
	if res["testvar"] != correct["testvar"] {
		t.Error("expected", correct["testvar"], "got", res["testvar"])
	}
}

func TestRenderThermoTemplate(t *testing.T) {
	brokenTemplate := ThermoTemplate{"", nil, "{{define test}}"}
	if res := renderThermoTemplate(&brokenTemplate, false); res != "" {
		t.Error("expected empty string, got", res)
	}

	resNormal := renderThermoTemplate(&testTemplate, false)
	correctNormal := template.HTML("<h1>Testing</h1><br>exampleval")
	if resNormal != correctNormal {
		t.Error("expected ", correctNormal, "got ", resNormal)
	}
	resPreview := renderThermoTemplate(&testTemplate, true)
	correctPreview := template.HTML("<h1>Testing</h1><br><strong>exampleval</strong>")
	if resNormal != correctNormal {
		t.Error("expected ", correctPreview, "got ", resPreview)
	}
}
