package main

import (
	"fmt"
	"github.com/grahms/html2pdf"
	"log"
	"os"
)

func main() {
	html := `<html>
<body>
<div>text</div>
<img src="https://pkg.go.dev/static/shared/gopher/package-search-700x300.jpeg"/>
<img src="https://go.dev/images/gophers/motorcycle.svg"/>
<img src="https://go.dev/images/go_google_case_study_carousel.png" />
</body>
</html>`

	// Create a new PDFGenerator instance
	pdfGen, err := html2pdf.NewWithPool(5)
	if err != nil {
		log.Fatalf("Failed to initialize PDF generator: %v", err)
	}
	defer pdfGen.Close()

	// Convert HTML to PDF
	pdfBuffer, err := pdfGen.GeneratePDF(&html)
	if err != nil {
		log.Fatalf("Failed to convert HTML to PDF: %v", err)
	}

	// Save the PDF buffer to a file
	if err := saveToFile("output.pdf", pdfBuffer); err != nil {
		log.Fatalf("Failed to save PDF to file: %v", err)
	}

	fmt.Println("PDF successfully created: output.pdf")
}

// saveToFile writes the PDF buffer to a file.
func saveToFile(fileName string, data []byte) error {
	return os.WriteFile(fileName, data, 0644)
}
