package server

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func embedhead(html string) string {
	head := "<head><meta charset='UTF-8'> </head><body>%s</body>"
	return fmt.Sprintf(head, html)
}

func GeneratePDF(html string, w io.Writer) {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		log.Fatal(err)
	}
	pdfg.Dpi.Set(300)
	pdfg.Grayscale.Set(true)
	pdfg.PageWidth.Set(72)
	pdfg.PageHeight.Set(200)
	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(embedhead(html))))

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
