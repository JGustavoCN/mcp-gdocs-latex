package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/gdocs"
)

func main() {
	docID := "1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8"
	if len(os.Args) > 1 {
		docID = os.Args[1]
	}

	ctx := context.Background()
	svc, err := gdocs.GetDocsService(ctx)
	if err != nil {
		log.Fatalf("erro svc: %v", err)
	}

	doc, err := svc.Documents.Get(docID).Do()
	if err != nil {
		log.Fatalf("erro doc: %v", err)
	}

	for i, el := range doc.Body.Content {
		if el.Paragraph != nil {
			txt := ""
			for _, pel := range el.Paragraph.Elements {
				if pel.TextRun != nil {
					txt += pel.TextRun.Content
				}
			}
			trimmed := strings.TrimSpace(txt)
			if strings.Contains(trimmed, "Trabalho de Conclusão") || strings.Contains(trimmed, "Orientador") || strings.Contains(trimmed, "LAGARTO") || strings.Contains(trimmed, "2026") {
				b, _ := json.Marshal(el.Paragraph.ParagraphStyle)
				fmt.Printf("[%d] Range: %d-%d | Text: %s\n  Style: %s\n", i, el.StartIndex, el.EndIndex, trimmed, string(b))
			}
		}
	}
}

