package html2pdf

import (
	"context"
	"errors"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// PDFGeneratorInterface defines the methods for generating PDFs.
type PDFGeneratorInterface interface {
	GeneratePDF(html *string) ([]byte, error)
	Close()
}

// PDFGenerator holds the browser context and manages PDF generation.
type PDFGenerator struct {
	browserCtx    context.Context
	browserCancel context.CancelFunc
	mu            sync.Mutex
}

// New creates a new PDFGenerator and initializes the browser context.
func New() (*PDFGenerator, error) {
	ctx, cancel := chromedp.NewContext(context.Background())

	// Start the browser
	log.Println("starting the chromium on the background for PDF generation")
	if err := chromedp.Run(ctx); err != nil {
		cancel()
		return nil, err
	}
	instance := &PDFGenerator{
		browserCtx:    ctx,
		browserCancel: cancel,
	}
	closeHandler(instance)
	return instance, nil
}

// Close shuts down the browser context.
func (p *PDFGenerator) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.browserCancel != nil {
		p.browserCancel()
		p.browserCancel = nil
		log.Println("Chromium instance closed")
	}
}

// GeneratePDF converts HTML content to a PDF buffer.
func (p *PDFGenerator) GeneratePDF(html *string) ([]byte, error) {
	if html == nil {
		return nil, errors.New("html content is required")
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.browserCtx == nil || p.browserCtx.Err() != nil {
		return nil, errors.New("browser context is not initialized or has been closed")
	}

	ctx, cancel := chromedp.NewContext(p.browserCtx)
	defer cancel()

	var pdfBuffer []byte
	if err := chromedp.Run(ctx, generatePDFTasks(html, &pdfBuffer)...); err != nil {
		return nil, err
	}

	return pdfBuffer, nil
}

// generatePDFTasks returns a sequence of tasks to generate a PDF from HTML content.
func generatePDFTasks(html *string, pdfBuffer *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		setDocumentContent(*html),
		printToPDF(pdfBuffer),
	}
}

// setDocumentContent sets the document content of the current page.
func setDocumentContent(html string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		frameTree, err := page.GetFrameTree().Do(ctx)
		if err != nil {
			return err
		}
		return page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx)
	}
}

// printToPDF prints the current page to PDF.
func printToPDF(pdfBuffer *[]byte) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		buf, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
		if err != nil {
			return err
		}
		*pdfBuffer = buf
		return nil
	}
}

// CloseHandler sets up a signal handler to gracefully shut down the PDFGenerator on interrupt signals.
func closeHandler(pdfGen PDFGeneratorInterface) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("gracefully shutting down the chromium")
		pdfGen.Close()
		os.Exit(0)
	}()
}
