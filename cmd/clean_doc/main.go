package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/gdocs"
	"google.golang.org/api/docs/v1"
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

	fmt.Printf("Inspecionando '%s' (Revisão: %s)...\n", doc.Title, doc.RevisionId)

	type RangeToDelete struct {
		Start int64
		End   int64
		Desc  string
	}

	var toDelete []RangeToDelete

	for i, el := range doc.Body.Content {
		if el.Table != nil {
			t := el.Table
			rows := len(t.TableRows)
			cols := 0
			if rows > 0 {
				cols = len(t.TableRows[0].TableCells)
			}
			emptyCells := 0
			totalCells := rows * cols

			for _, r := range t.TableRows {
				for _, c := range r.TableCells {
					cellText := ""
					for _, cEl := range c.Content {
						if cEl.Paragraph != nil {
							for _, pEl := range cEl.Paragraph.Elements {
								if pEl.TextRun != nil {
									cellText += pEl.TextRun.Content
								}
							}
						}
					}
					if strings.TrimSpace(cellText) == "" {
						emptyCells++
					}
				}
			}

			if emptyCells == totalCells {
				fmt.Printf("🗑️ Marcando para exclusão Tabela Vazia: Elemento %d [%d - %d] (%dx%d)\n",
					i, el.StartIndex, el.EndIndex, rows, cols)
				toDelete = append(toDelete, RangeToDelete{
					Start: el.StartIndex,
					End:   el.EndIndex,
					Desc:  fmt.Sprintf("Tabela Vazia [%d-%d]", el.StartIndex, el.EndIndex),
				})
			}
		}
	}

	if len(toDelete) == 0 {
		fmt.Println("Nenhuma tabela vazia encontrada para exclusão.")
		return
	}

	// Deleta de trás para frente para não alterar os índices anteriores
	for i := len(toDelete) - 1; i >= 0; i-- {
		item := toDelete[i]
		fmt.Printf("Excluindo %s...\n", item.Desc)
		req := &docs.BatchUpdateDocumentRequest{
			Requests: []*docs.Request{
				{
					DeleteContentRange: &docs.DeleteContentRangeRequest{
						Range: &docs.Range{
							StartIndex: item.Start,
							EndIndex:   item.End,
						},
					},
				},
			},
		}
		_, err := docsSvc.Documents.BatchUpdate(testDocID, req).Do()
		if err != nil {
			fmt.Printf("❌ Falha ao excluir %s: %v\n", item.Desc, err)
		} else {
			fmt.Printf("✅ Excluída com sucesso: %s\n", item.Desc)
		}
	}

	fmt.Println("🎉 Limpeza concluída com sucesso!")
}
