package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/gdocs"
)

const (
	testDocID   = "1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8"
	targetPdfID = "1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH"
)

func main() {
	credPath, _ := filepath.Abs("credentials.json")
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", credPath)
	ctx := context.Background()

	fmt.Println("=== 1. DIAGNÓSTICO DO GOOGLE DRIVE (COTA E PAIS) ===")
	driveSvc, err := gdocs.GetDriveService(ctx)
	if err != nil {
		fmt.Printf("Erro drive: %v\n", err)
		return
	}

	about, err := driveSvc.About.Get().Fields("user, storageQuota").Do()
	if err == nil {
		aboutJSON, _ := json.MarshalIndent(about, "", "  ")
		fmt.Printf("About Drive:\n%s\n", string(aboutJSON))
	} else {
		fmt.Printf("Erro About: %v\n", err)
	}

	docMeta, err := driveSvc.Files.Get(testDocID).Fields("id, name, parents, owners, capabilities").Do()
	if err == nil {
		metaJSON, _ := json.MarshalIndent(docMeta, "", "  ")
		fmt.Printf("Doc Meta:\n%s\n", string(metaJSON))
	}

	pdfMeta, err := driveSvc.Files.Get(targetPdfID).Fields("id, name, parents, owners, capabilities").Do()
	if err == nil {
		pdfJSON, _ := json.MarshalIndent(pdfMeta, "", "  ")
		fmt.Printf("PDF Meta:\n%s\n", string(pdfJSON))
	}

	fmt.Println("\n=== 2. VARREDURA COMPLETA DE ANOMALIAS NO GOOGLE DOCS ===")
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

	fmt.Printf("Título: %s | Revisão: %s | Total Elementos: %d\n", doc.Title, doc.RevisionId, len(doc.Body.Content))

	var emptyTables []string
	var filledTables []string
	var boldCascades []string
	var emptyParagraphs []string
	var softBreakList []string
	var markerList []string

	consecutiveEmptyParas := 0

	for i, el := range doc.Body.Content {
		// 1. Inspeciona Tabelas
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

			desc := fmt.Sprintf("Elemento %d [%d-%d]: Tabela %dx%d (Células Vazias: %d/%d)",
				i, el.StartIndex, el.EndIndex, rows, cols, emptyCells, totalCells)
			if emptyCells == totalCells {
				emptyTables = append(emptyTables, desc+" [TOTALMENTE VAZIA!]")
			} else if emptyCells > 0 {
				emptyTables = append(emptyTables, desc+" [PARCIALMENTE VAZIA]")
			} else {
				filledTables = append(filledTables, desc+" [PREENCHIDA]")
			}
		}

		// 2. Inspeciona Parágrafos
		if el.Paragraph != nil {
			p := el.Paragraph
			pText := ""
			pBoldLen := 0
			hasSoftBreak := false

			for _, pEl := range p.Elements {
				if pEl.TextRun != nil {
					tr := pEl.TextRun
					pText += tr.Content

					if strings.Contains(tr.Content, "\u000b") || strings.Contains(tr.Content, "\x0b") {
						hasSoftBreak = true
					}

					if tr.TextStyle != nil && tr.TextStyle.Bold {
						pBoldLen += len(tr.Content)
					}
				}
			}

			trimmed := strings.TrimSpace(pText)

			// Parágrafo vazio ou apenas quebra
			if trimmed == "" {
				consecutiveEmptyParas++
				if consecutiveEmptyParas >= 2 {
					emptyParagraphs = append(emptyParagraphs, fmt.Sprintf("Elemento %d [%d-%d]: Quebra múltipla/Parágrafo vazio consecutivo (#%d)",
						i, el.StartIndex, el.EndIndex, consecutiveEmptyParas))
				}
			} else {
				consecutiveEmptyParas = 0
			}

			// Cascata de negrito (> 150 caracteres contínuos em negrito fora de títulos)
			isHeading := strings.HasPrefix(p.ParagraphStyle.NamedStyleType, "HEADING_")
			if !isHeading && pBoldLen > 150 {
				snippet := trimmed
				if len(snippet) > 80 {
					snippet = snippet[:80] + "..."
				}
				boldCascades = append(boldCascades, fmt.Sprintf("Elemento %d [%d-%d]: Cascata de Negrito (%d caracteres em negrito) -> \"%s\"",
					i, el.StartIndex, el.EndIndex, pBoldLen, snippet))
			}

			// Soft breaks (\u000b)
			if hasSoftBreak {
				snippet := trimmed
				if len(snippet) > 80 {
					snippet = snippet[:80] + "..."
				}
				softBreakList = append(softBreakList, fmt.Sprintf("Elemento %d [%d-%d]: Contém quebra de linha manual (\\u000b) -> \"%s\"",
					i, el.StartIndex, el.EndIndex, snippet))
			}

			// Marcadores artificiais soltos ([TABELA], [VER IMAGEM...])
			if strings.Contains(pText, "[TABELA]") || strings.Contains(pText, "[/TABELA]") || strings.Contains(pText, "[VER IMAGEM") {
				markerList = append(markerList, fmt.Sprintf("Elemento %d [%d-%d]: Marcador textual solto -> \"%s\"",
					i, el.StartIndex, el.EndIndex, trimmed))
			}
		}
	}

	fmt.Printf("\n--- 📊 TABELAS IDENTIFICADAS (%d) ---\n", len(emptyTables)+len(filledTables))
	for _, t := range emptyTables {
		fmt.Println("  ⚠️ ", t)
	}
	for _, t := range filledTables {
		fmt.Println("  ✅ ", t)
	}

	fmt.Printf("\n--- ⚠️ CASCATAS DE NEGRITO ANÔMALAS (%d) ---\n", len(boldCascades))
	for _, b := range boldCascades {
		fmt.Println("  ", b)
	}

	fmt.Printf("\n--- ⚠️ PARÁGRAFOS VAZIOS / QUEBRAS REPETIDAS (%d) ---\n", len(emptyParagraphs))
	for _, ep := range emptyParagraphs {
		fmt.Println("  ", ep)
	}

	fmt.Printf("\n--- ⚠️ QUEBRAS MANUAIS DE LINHA (\\u000b) (%d) ---\n", len(softBreakList))
	for _, sb := range softBreakList {
		fmt.Println("  ", sb)
	}

	fmt.Printf("\n--- ⚠️ MARCADORES TEXTUAIS SOLTOS (%d) ---\n", len(markerList))
	for _, m := range markerList {
		fmt.Println("  ", m)
	}
}
