package main

import (
	"context"
	"fmt"
	"os"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/gdocs"
	"github.com/JGustavoCN/mcp-gdocs-latex/internal/latex"
)

func main() {
	docID := "1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8"
	ctx := context.Background()

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "c:/Projetos/golang/mcp-gdocs-latex/credentials.json")

	docsSvc, err := gdocs.GetDocsService(ctx)
	if err != nil {
		fmt.Printf("GetDocsService err: %v\n", err)
		return
	}

	doc, err := docsSvc.Documents.Get(docID).Do()
	if err != nil {
		fmt.Printf("Get doc err: %v\n", err)
		return
	}

	meta, err := latex.ParseTCCMetadata("")
	if err != nil {
		fmt.Printf("ParseTCCMetadata err: %v\n", err)
	}

	pdfBytes, err := gdocs.ExportDocPDF(ctx, docID)
	if err != nil {
		fmt.Printf("ExportDocPDF err: %v\n", err)
		return
	}

	pages, err := gdocs.ExtractPDFPages(pdfBytes)
	if err != nil {
		fmt.Printf("ExtractPDFPages err: %v\n", err)
		return
	}

	fmt.Printf("Total PDF pages: %d\n", len(pages))
	for i := 0; i < len(pages) && i < 5; i++ {
		fmt.Printf("--- PÁGINA %d (len=%d) ---\n%s\n", i+1, len(pages[i]), pages[i])
	}

	report, err := gdocs.VerifyPreTextualPagination(ctx, docID, meta, doc)
	fmt.Printf("\nVerify report: %+v, err: %v\n", report, err)
}
