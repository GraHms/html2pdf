# html2pdf

`html2pdf` is a Go library that provides an interface for converting HTML content to PDF using Chromium. It interacts with a headless Chromium instance and generate PDFs from HTML content.

## Installation

To install the `html2pdf` package, use `go get`:

```sh
go get github.com/grahms/html2pdf
```

## Prerequisites

Ensure that Chromium or Google Chrome is installed on your machine. You can download it from [here](https://www.chromium.org/getting-involved/download-chromium) or install it using a package manager:

- On macOS:
  ```sh
  brew install chromium
  ```
- On Ubuntu/Debian:
  ```sh
  sudo apt-get install chromium-browser
  ```

## Usage

Here's a basic example of how to use the `html2pdf` library to convert HTML content to a PDF file:

### Example

```go
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
	pdfGen, err := html2pdf.New()
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
```

## Features

- **Easy to use**: Simple API to generate PDFs from HTML content.
- **Graceful shutdown**: Handles interrupt signals to ensure the Chromium instance is properly terminated.
- **Mockable**: Provides an interface for easy testing and mocking.

## API

### `PDFGeneratorInterface`

```go
type PDFGeneratorInterface interface {
	GeneratePDF(html *string) ([]byte, error)
	Close()
}
```

### `PDFGenerator`

#### `New() (*PDFGenerator, error)`

Creates a new `PDFGenerator` and initializes the browser context.

#### `GeneratePDF(html *string) ([]byte, error)`

Converts HTML content to a PDF buffer.

#### `Close()`

Shuts down the browser context.



## Handling Interrupt Signals

The library automatically sets up a signal handler to gracefully shut down the `PDFGenerator` on interrupt signals.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License.