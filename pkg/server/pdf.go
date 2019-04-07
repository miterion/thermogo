package server

import (
	"strconv"
	"bytes"
	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"io"
	"log"
	"os/exec"
	"strings"
)

func GeneratePDF(html string, w io.Writer) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}
	pdfg.Dpi.Set(300)
	pdfg.Grayscale.Set(true)
	pdfg.PageWidth.Set(72)
	pdfg.PageHeight.Set(100)
	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(html)))

	pdfg.SetOutput(w)

	err = pdfg.Create()
	if err != nil {
		log.Println(err)
	}
}

func Print(html string, copies int) {
	var pdfBuilder bytes.Buffer
	GeneratePDF(html, &pdfBuilder)
	cmd := exec.Command("lp", "-d", "Thermodrucker", "-n", strconv.Itoa(copies))
	cmd.Stdin = &pdfBuilder
	go cmd.Run()
}
