//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/gdocs"
)

func main() {
	ctx := context.Background()
	docID := "1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8"

	svc, err := gdocs.GetDocsService(ctx)
	if err != nil {
		log.Fatalf("Erro svc: %v", err)
	}

	doc, err := svc.Documents.Get(docID).Do()
	if err != nil {
		log.Fatalf("Erro doc: %v", err)
	}

	indices := []int{186, 187, 190, 200}
	for _, idx := range indices {
		if idx < len(doc.Body.Content) {
			el := doc.Body.Content[idx]
			txt := ""
			if el.Paragraph != nil {
				for _, pel := range el.Paragraph.Elements {
					if pel.TextRun != nil {
						txt += pel.TextRun.Content
					}
				}
			}
			fmt.Printf("=== ELEMENTO %d [%d - %d] (Len: %d) ===\n%s\n\n",
				idx, el.StartIndex, el.EndIndex, len(strings.TrimSpace(txt)), strings.TrimSpace(txt))
		}
	}
}

