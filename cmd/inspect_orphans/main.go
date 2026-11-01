package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/gdocs"
)

const testDocID = "1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8"

func main() {
	credPath, _ := filepath.Abs("credentials.json")
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credPath)
	ctx := context.Background()

	docsSvc, err := gdocs.GetDocsService(ctx)
	if err != nil {
		fmt.Printf("Erro docs: %v\n", err)
		return
	}

	doc, err := docsSvc.Documents.Get(testDocID).Do()
	if err != nil {
		fmt.Printf("Erro ao buscar doc: %v\n", err)
		return
	}

	targets := []string{
		"Etapa de Engenharia",
		"Etapa de Avaliação",
		"Didática e o Calendário",
		"como atestados médicos",
	}

	for i, el := range doc.Body.Content {
		if el.Paragraph != nil {
			txt := ""
			for _, pel := range el.Paragraph.Elements {
				if pel.TextRun != nil {
					txt += pel.TextRun.Content
				}
			}
			t := strings.TrimSpace(txt)
			for _, tg := range targets {
				if strings.Contains(t, tg) {
					// Imprime contexto: elemento anterior, atual e posterior
					fmt.Printf("\n=== MATCH %q no Elemento %d [%d - %d] ===\n", tg, i, el.StartIndex, el.EndIndex)
					if i > 0 && doc.Body.Content[i-1].Paragraph != nil {
						pPrev := ""
						for _, pel := range doc.Body.Content[i-1].Paragraph.Elements {
							if pel.TextRun != nil {
								pPrev += pel.TextRun.Content
							}
						}
						fmt.Printf("ANTERIOR (Elem %d): %q\n", i-1, strings.TrimSpace(pPrev))
					}
					fmt.Printf("ATUAL:    %q\n", t)
					if i+1 < len(doc.Body.Content) && doc.Body.Content[i+1].Paragraph != nil {
						pNext := ""
						for _, pel := range doc.Body.Content[i+1].Paragraph.Elements {
							if pel.TextRun != nil {
								pNext += pel.TextRun.Content
							}
						}
						fmt.Printf("POSTERIOR (Elem %d): %q\n", i+1, strings.TrimSpace(pNext))
					}
				}
			}
		}
	}
}

