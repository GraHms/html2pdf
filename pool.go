package html2pdf

import (
	"sync"
)

type PDFWorkerPool struct {
	pdfGen *PDFGenerator
	tasks  chan func() ([]byte, error)
	wg     sync.WaitGroup
}

// NewWithPool creates a PDFWorkerPool with the specified number of worker tabs.
func NewWithPool(size int) (*PDFWorkerPool, error) {
	pdfGen, err := New()
	if err != nil {
		return nil, err
	}
	return newPDFWorkerPool(pdfGen, size), nil
}

func newPDFWorkerPool(pdfGen *PDFGenerator, size int) *PDFWorkerPool {
	pool := &PDFWorkerPool{
		pdfGen: pdfGen,
		tasks:  make(chan func() ([]byte, error), size),
	}
	for i := 0; i < size; i++ {
		pool.wg.Add(1)
		go pool.worker()
	}
	return pool
}

func (p *PDFWorkerPool) worker() {
	defer p.wg.Done()
	for task := range p.tasks {
		task()
	}
}

func (p *PDFWorkerPool) GeneratePDF(html *string) ([]byte, error) {
	resultChan := make(chan struct {
		buf []byte
		err error
	}, 1)

	task := func() ([]byte, error) {
		buf, err := p.pdfGen.GeneratePDF(html)
		resultChan <- struct {
			buf []byte
			err error
		}{buf, err}
		return buf, err
	}

	p.tasks <- task

	result := <-resultChan
	return result.buf, result.err
}

func (p *PDFWorkerPool) Close() {
	close(p.tasks)
	p.wg.Wait()
	p.pdfGen.Close()
}
