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
	correct := map[string]string{"testvar": "<span class='previewtext' id='prev-testvar'>$testvar$</span>"}
	if res["testvar"] != correct["testvar"] {
		t.Error("expected", correct["testvar"], "got", res["testvar"])
	}
}

func TestRenderThermoTemplate(t *testing.T) {
	brokenTemplate := ThermoTemplate{"", nil, "{{define test}}"}
	if res := brokenTemplate.renderThermoTemplate(false); res != "" {
		t.Error("expected empty string, got", res)
	}

	resNormal := testTemplate.renderThermoTemplate(false)
	correctNormal := template.HTML("<h1>Testing</h1><br>exampleval")
	if resNormal != correctNormal {
		t.Error("expected ", correctNormal, "got ", resNormal)
	}
	resPreview := testTemplate.renderThermoTemplate(false)
	correctPreview := template.HTML("<h1>Testing</h1><br>$exampleval$")
	if resNormal != correctNormal {
		t.Error("expected ", correctPreview, "got ", resPreview)
	}
}
