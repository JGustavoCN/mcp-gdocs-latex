package gdocs

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/api/docs/v1"
)

// LintReport contém os resultados da análise estrutural e tipográfica do Google Doc
type LintReport struct {
	TotalElements          int      `json:"total_elements"`
	TotalTables            int      `json:"total_tables"`
	EmptyTables            []string `json:"empty_tables,omitempty"`
	FilledTables           []string `json:"filled_tables,omitempty"`
	BoldCascades           []string `json:"bold_cascades,omitempty"`
	ConsecutiveBlankLines  []string `json:"consecutive_blank_lines,omitempty"`
	SoftLineBreaks         []string `json:"soft_line_breaks,omitempty"`
	DanglingMarkers        []string `json:"dangling_markers,omitempty"`
	FixedAnomaliesCount    int      `json:"fixed_anomalies_count,omitempty"`
	Message                string   `json:"message"`
}

// LintDocumentABNT inspeciona anomalias no documento e opcionalmente aplica autocorreções
func LintDocumentABNT(ctx context.Context, documentID string, autoFix bool) (*LintReport, error) {
	docsSvc, err := GetDocsService(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter docs service: %w", err)
	}

	doc, err := docsSvc.Documents.Get(documentID).Do()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar documento: %w", err)
	}

	report := &LintReport{
		TotalElements: len(doc.Body.Content),
	}

	type TableToDelete struct {
		Start int64
		End   int64
	}
	var emptyTableRanges []TableToDelete

	consecutiveEmptyParas := 0

	for i, el := range doc.Body.Content {
		// 1. Análise de Tabelas
		if el.Table != nil {
			report.TotalTables++
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
				report.EmptyTables = append(report.EmptyTables, desc+" [TOTALMENTE VAZIA]")
				emptyTableRanges = append(emptyTableRanges, TableToDelete{Start: el.StartIndex, End: el.EndIndex})
			} else if emptyCells > 0 {
				report.EmptyTables = append(report.EmptyTables, desc+" [PARCIALMENTE VAZIA]")
			} else {
				report.FilledTables = append(report.FilledTables, desc+" [PREENCHIDA]")
			}
		}

		// 2. Análise de Parágrafos
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

			// Parágrafo vazio ou apenas quebra de linha
			if trimmed == "" {
				consecutiveEmptyParas++
				if consecutiveEmptyParas >= 2 {
					report.ConsecutiveBlankLines = append(report.ConsecutiveBlankLines,
						fmt.Sprintf("Elemento %d [%d-%d]: Quebra múltipla/Parágrafo vazio consecutivo (#%d)",
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
				report.BoldCascades = append(report.BoldCascades,
					fmt.Sprintf("Elemento %d [%d-%d]: Cascata de Negrito (%d caracteres em negrito) -> \"%s\"",
						i, el.StartIndex, el.EndIndex, pBoldLen, snippet))
			}

			// Quebras manuais de linha (\u000b)
			if hasSoftBreak {
				snippet := trimmed
				if len(snippet) > 80 {
					snippet = snippet[:80] + "..."
				}
				report.SoftLineBreaks = append(report.SoftLineBreaks,
					fmt.Sprintf("Elemento %d [%d-%d]: Quebra manual (\\u000b) -> \"%s\"",
						i, el.StartIndex, el.EndIndex, snippet))
			}

			// Marcadores artificiais soltos
			if strings.Contains(pText, "[TABELA]") || strings.Contains(pText, "[/TABELA]") || strings.Contains(pText, "[VER IMAGEM") {
				report.DanglingMarkers = append(report.DanglingMarkers,
					fmt.Sprintf("Elemento %d [%d-%d]: Marcador artificial -> \"%s\"",
						i, el.StartIndex, el.EndIndex, trimmed))
			}
		}
	}

	// 3. Execução de Autocorreção (se autoFix == true)
	if autoFix {
		fixedCount := 0

		// 3.1 Exclui tabelas vazias (de trás para frente para preservar índices)
		for i := len(emptyTableRanges) - 1; i >= 0; i-- {
			tbl := emptyTableRanges[i]
			req := &docs.BatchUpdateDocumentRequest{
				Requests: []*docs.Request{
					{
						DeleteContentRange: &docs.DeleteContentRangeRequest{
							Range: &docs.Range{
								StartIndex: tbl.Start,
								EndIndex:   tbl.End,
							},
						},
					},
				},
			}
			if _, errDel := docsSvc.Documents.BatchUpdate(documentID, req).Do(); errDel == nil {
				fixedCount++
			}
		}

		// 3.2 Substitui quebras manuais \u000b por espaço simples
		if len(report.SoftLineBreaks) > 0 {
			reqReplace := &docs.BatchUpdateDocumentRequest{
				Requests: []*docs.Request{
					{
						ReplaceAllText: &docs.ReplaceAllTextRequest{
							ContainsText: &docs.SubstringMatchCriteria{
								MatchCase: true,
								Text:      "\u000b",
							},
							ReplaceText: " ",
						},
					},
				},
			}
			if _, errSoft := docsSvc.Documents.BatchUpdate(documentID, reqReplace).Do(); errSoft == nil {
				fixedCount += len(report.SoftLineBreaks)
			}
		}

		report.FixedAnomaliesCount = fixedCount
		report.Message = fmt.Sprintf("Auditoria concluída. %d anomalias corrigidas automaticamente.", fixedCount)
	} else {
		totalAnomalies := len(report.EmptyTables) + len(report.BoldCascades) + len(report.ConsecutiveBlankLines) + len(report.SoftLineBreaks) + len(report.DanglingMarkers)
		report.Message = fmt.Sprintf("Auditoria concluída. %d anomalias identificadas no documento (autoFix=false).", totalAnomalies)
	}

	return report, nil
}

