package gdocs

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/googleapi"
)

// ReviewHighlightItem detalha um trecho com marca-texto / grifo de revisão deixado pelo orientador/revisor
type ReviewHighlightItem struct {
	StartIndex        int64  `json:"start_index"`
	EndIndex          int64  `json:"end_index"`
	ColorCategory     string `json:"color_category"`
	RgbHex            string `json:"rgb_hex"`
	Text              string `json:"text"`
	AssociatedComment string `json:"associated_comment,omitempty"`
}

// FontSizeAnomalyItem detalha trechos com tamanho de fonte fora do padrão ABNT
type FontSizeAnomalyItem struct {
	StartIndex   int64   `json:"start_index"`
	EndIndex     int64   `json:"end_index"`
	CurrentSize  float64 `json:"current_size_pt"`
	ExpectedSize float64 `json:"expected_size_pt"`
	ContextType  string  `json:"context_type"`
	Text         string  `json:"text"`
}

// TextColorAnomalyItem detalha trechos com cor de texto anômala fora do preto padrão
type TextColorAnomalyItem struct {
	StartIndex int64  `json:"start_index"`
	EndIndex   int64  `json:"end_index"`
	RgbHex     string `json:"rgb_hex"`
	Text       string `json:"text"`
}

// LintReport contém os resultados da análise estrutural, tipográfica e de revisões do Google Doc
type LintReport struct {
	TotalElements         int                    `json:"total_elements"`
	TotalTables           int                    `json:"total_tables"`
	EmptyTables           []string               `json:"empty_tables,omitempty"`
	FilledTables          []string               `json:"filled_tables,omitempty"`
	BoldCascades          []string               `json:"bold_cascades,omitempty"`
	ReviewHighlights      []ReviewHighlightItem  `json:"review_highlights,omitempty"`
	FontSizeAnomalies     []FontSizeAnomalyItem  `json:"font_size_anomalies,omitempty"`
	TextColorAnomalies    []TextColorAnomalyItem `json:"text_color_anomalies,omitempty"`
	ConsecutiveBlankLines []string               `json:"consecutive_blank_lines,omitempty"`
	SoftLineBreaks        []string               `json:"soft_line_breaks,omitempty"`
	DanglingMarkers       []string               `json:"dangling_markers,omitempty"`
	FixedAnomaliesCount   int                    `json:"fixed_anomalies_count,omitempty"`
	Message               string                 `json:"message"`
}

// classifyColor mapeia valores RGB (0.0 a 1.0) para categoria semântica e código HEX
func classifyColor(r, g, b float64) (string, string) {
	hex := fmt.Sprintf("#%02X%02X%02X", int(r*255), int(g*255), int(b*255))
	if r > 0.8 && g < 0.35 && b < 0.35 {
		return "Vermelho (Alerta / Correção Crítica)", hex
	}
	if r > 0.8 && g > 0.75 && b < 0.4 {
		return "Amarelo (Atenção / Dúvida / Comentário)", hex
	}
	if g > 0.7 && r < 0.5 {
		return "Verde (Aprovação / Observação)", hex
	}
	if b > 0.7 && r < 0.4 {
		return "Azul (Referência / Link / Nota)", hex
	}
	if r > 0.8 && g < 0.5 && b > 0.5 {
		return "Magenta / Rosa (Revisão)", hex
	}
	return fmt.Sprintf("RGB(%.2f, %.2f, %.2f)", r, g, b), hex
}

// getDocCommentsSummary coleta os comentários ativos no Google Drive para cruzamento
func getDocCommentsSummary(ctx context.Context, fileID string) map[string]string {
	commentMap := make(map[string]string)
	driveSvc, err := GetDriveService(ctx)
	if err != nil {
		return commentMap
	}

	fields := "comments(id,content,author/displayName,quotedFileContent/value,resolved)"
	res, err := driveSvc.Comments.List(fileID).Fields(googleapi.Field(fields)).IncludeDeleted(false).PageSize(100).Do()
	if err != nil || res == nil {
		return commentMap
	}

	for _, c := range res.Comments {
		if c.QuotedFileContent != nil && strings.TrimSpace(c.QuotedFileContent.Value) != nilText(c.QuotedFileContent.Value) {
			cleanQuoted := strings.ToLower(strings.TrimSpace(c.QuotedFileContent.Value))
			author := "Autor"
			if c.Author != nil && c.Author.DisplayName != "" {
				author = c.Author.DisplayName
			}
			commentMap[cleanQuoted] = fmt.Sprintf("[%s]: %s", author, c.Content)
		}
	}

	return commentMap
}

func nilText(s string) string {
	return ""
}

// LintDocumentABNT inspeciona o documento, detectando tabelas vazias, cascatas de negrito,
// anomalias de fonte ABNT, marca-texto de revisão e comentários do orientador.
func LintDocumentABNT(ctx context.Context, documentID string, autoFix bool, clearHighlights bool) (*LintReport, error) {
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

	// Carrega comentários do Drive para correlacionar com marca-texto
	commentsSummary := getDocCommentsSummary(ctx, documentID)

	type RangeItem struct {
		Start int64
		End   int64
	}
	var emptyTableRanges []RangeItem
	var fontFixRanges []RangeItem
	var highlightRanges []RangeItem

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
				emptyTableRanges = append(emptyTableRanges, RangeItem{Start: el.StartIndex, End: el.EndIndex})
			} else if emptyCells > 0 {
				report.EmptyTables = append(report.EmptyTables, desc+" [PARCIALMENTE VAZIA]")
			} else {
				report.FilledTables = append(report.FilledTables, desc+" [PREENCHIDA]")
			}
		}

		// 2. Análise de Parágrafos
		if el.Paragraph != nil {
			p := el.Paragraph
			styleName := ""
			if p.ParagraphStyle != nil {
				styleName = p.ParagraphStyle.NamedStyleType
			}
			isHeading := strings.HasPrefix(styleName, "HEADING_") || styleName == "TITLE"

			pText := ""
			pBoldLen := 0
			hasSoftBreak := false

			for _, pEl := range p.Elements {
				if pEl.TextRun != nil {
					tr := pEl.TextRun
					txt := tr.Content
					pText += txt
					ts := tr.TextStyle

					if strings.Contains(txt, "\u000b") || strings.Contains(txt, "\x0b") {
						hasSoftBreak = true
					}

					if ts != nil {
						if ts.Bold {
							pBoldLen += len(txt)
						}

						// 2.1 Detecção de Marca-texto (BackgroundColor)
						if ts.BackgroundColor != nil && ts.BackgroundColor.Color != nil && ts.BackgroundColor.Color.RgbColor != nil {
							rgb := ts.BackgroundColor.Color.RgbColor
							colorCat, hex := classifyColor(rgb.Red, rgb.Green, rgb.Blue)
							cleanTxt := strings.TrimSpace(txt)

							// Tenta correlacionar com comentário existente
							associatedComment := ""
							cleanLower := strings.ToLower(cleanTxt)
							for q, comm := range commentsSummary {
								if strings.Contains(q, cleanLower) || strings.Contains(cleanLower, q) {
									associatedComment = comm
									break
								}
							}

							report.ReviewHighlights = append(report.ReviewHighlights, ReviewHighlightItem{
								StartIndex:        pEl.StartIndex,
								EndIndex:          pEl.EndIndex,
								ColorCategory:     colorCat,
								RgbHex:            hex,
								Text:              cleanTxt,
								AssociatedComment: associatedComment,
							})

							highlightRanges = append(highlightRanges, RangeItem{Start: pEl.StartIndex, End: pEl.EndIndex})
						}

						// 2.2 Detecção de Cores Anômalas de Texto (ForegroundColor != preto)
						if ts.ForegroundColor != nil && ts.ForegroundColor.Color != nil && ts.ForegroundColor.Color.RgbColor != nil {
							rgb := ts.ForegroundColor.Color.RgbColor
							// Ignora se for preto (ou cinza quase preto < 0.05) e ignora quebras puras
							if (rgb.Red > 0.15 || rgb.Green > 0.15 || rgb.Blue > 0.15) && strings.TrimSpace(txt) != "" {
								// Verifica se não é um link azul
								isLink := ts.Link != nil && ts.Link.Url != ""
								if !isLink {
									_, hex := classifyColor(rgb.Red, rgb.Green, rgb.Blue)
									report.TextColorAnomalies = append(report.TextColorAnomalies, TextColorAnomalyItem{
										StartIndex: pEl.StartIndex,
										EndIndex:   pEl.EndIndex,
										RgbHex:     hex,
										Text:       strings.TrimSpace(txt),
									})
								}
							}
						}

						// 2.3 Detecção de Tamanho de Fonte ABNT
						if ts.FontSize != nil && strings.TrimSpace(txt) != "" {
							pt := ts.FontSize.Magnitude

							trimmedP := strings.TrimSpace(pText)
							isCaptionOrSource := strings.HasPrefix(trimmedP, "Figura ") ||
								strings.HasPrefix(trimmedP, "Tabela ") ||
								strings.HasPrefix(trimmedP, "Quadro ") ||
								strings.HasPrefix(trimmedP, "Fonte:") ||
								strings.HasPrefix(trimmedP, "Nota:")

							if isCaptionOrSource {
								// Legenda e Fonte exigem 10 pt
								if pt < 9.5 || pt > 10.5 {
									report.FontSizeAnomalies = append(report.FontSizeAnomalies, FontSizeAnomalyItem{
										StartIndex:   pEl.StartIndex,
										EndIndex:     pEl.EndIndex,
										CurrentSize:  pt,
										ExpectedSize: 10.0,
										ContextType:  "Legenda / Fonte (ABNT)",
										Text:         strings.TrimSpace(txt),
									})
								}
							} else if !isHeading {
								// Corpo comum (NORMAL_TEXT) exige 12 pt
								if pt < 11.5 || pt > 12.5 {
									report.FontSizeAnomalies = append(report.FontSizeAnomalies, FontSizeAnomalyItem{
										StartIndex:   pEl.StartIndex,
										EndIndex:     pEl.EndIndex,
										CurrentSize:  pt,
										ExpectedSize: 12.0,
										ContextType:  "Corpo do Texto (ABNT)",
										Text:         strings.TrimSpace(txt),
									})
									fontFixRanges = append(fontFixRanges, RangeItem{Start: pEl.StartIndex, End: pEl.EndIndex})
								}
							}
						}
					}
				}
			}

			trimmed := strings.TrimSpace(pText)

			// Parágrafo vazio ou quebra múltipla
			if trimmed == "" {
				consecutiveEmptyParas++
				if consecutiveEmptyParas >= 2 {
					report.ConsecutiveBlankLines = append(report.ConsecutiveBlankLines,
						fmt.Sprintf("Elemento %d [%d-%d]: Quebra múltipla / Linha em branco excessiva (#%d)",
							i, el.StartIndex, el.EndIndex, consecutiveEmptyParas))
				}
			} else {
				consecutiveEmptyParas = 0
			}

			// Cascata de negrito (> 120 caracteres contínuos em negrito fora de títulos)
			if !isHeading && pBoldLen > 120 {
				snippet := trimmed
				if len(snippet) > 80 {
					snippet = snippet[:80] + "..."
				}
				report.BoldCascades = append(report.BoldCascades,
					fmt.Sprintf("Elemento %d [%d-%d]: Cascata de Negrito (%d caracteres em negrito) -> %q",
						i, el.StartIndex, el.EndIndex, pBoldLen, snippet))
			}

			// Quebras manuais de linha (\u000b)
			if hasSoftBreak {
				snippet := trimmed
				if len(snippet) > 80 {
					snippet = snippet[:80] + "..."
				}
				report.SoftLineBreaks = append(report.SoftLineBreaks,
					fmt.Sprintf("Elemento %d [%d-%d]: Quebra manual (\\u000b) -> %q",
						i, el.StartIndex, el.EndIndex, snippet))
			}

			// Marcadores artificiais soltos
			if strings.Contains(pText, "[TABELA]") || strings.Contains(pText, "[/TABELA]") || strings.Contains(pText, "[VER IMAGEM") {
				report.DanglingMarkers = append(report.DanglingMarkers,
					fmt.Sprintf("Elemento %d [%d-%d]: Marcador artificial -> %q",
						i, el.StartIndex, el.EndIndex, trimmed))
			}
		}
	}

	// 3. Execução de Autocorreção (se autoFix == true)
	if autoFix {
		fixedCount := 0
		var requests []*docs.Request

		// 3.1 Exclui tabelas vazias (de trás para frente para manter índices)
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

		// 3.3 Normaliza tamanhos de fonte do corpo para 12 pt
		for _, fr := range fontFixRanges {
			requests = append(requests, &docs.Request{
				UpdateTextStyle: &docs.UpdateTextStyleRequest{
					Range: &docs.Range{
						StartIndex: fr.Start,
						EndIndex:   fr.End,
					},
					TextStyle: &docs.TextStyle{
						FontSize: &docs.Dimension{
							Magnitude: 12.0,
							Unit:      "PT",
						},
					},
					Fields: "fontSize",
				},
			})
		}

		// 3.4 Limpeza de marca-texto se solicitado explicitamente
		if clearHighlights && len(highlightRanges) > 0 {
			for _, hr := range highlightRanges {
				requests = append(requests, &docs.Request{
					UpdateTextStyle: &docs.UpdateTextStyleRequest{
						Range: &docs.Range{
							StartIndex: hr.Start,
							EndIndex:   hr.End,
						},
						TextStyle: &docs.TextStyle{
							// Para remover marca-texto no Google Docs, não define BackgroundColor
							// ou define campos vazios
						},
						Fields: "backgroundColor",
					},
				})
			}
		}

		if len(requests) > 0 {
			batchReq := &docs.BatchUpdateDocumentRequest{Requests: requests}
			if _, errBatch := docsSvc.Documents.BatchUpdate(documentID, batchReq).Do(); errBatch == nil {
				fixedCount += len(requests)
			}
		}

		report.FixedAnomaliesCount = fixedCount
		report.Message = fmt.Sprintf("Auditoria concluída com auto_fix=true. %d anomalias corrigidas com sucesso.", fixedCount)
	} else {
		totalAnomalies := len(report.EmptyTables) + len(report.BoldCascades) + len(report.ConsecutiveBlankLines) +
			len(report.SoftLineBreaks) + len(report.DanglingMarkers) + len(report.ReviewHighlights) +
			len(report.FontSizeAnomalies) + len(report.TextColorAnomalies)
		report.Message = fmt.Sprintf("Auditoria concluída. %d anomalias/destaques identificados no documento (autoFix=false).", totalAnomalies)
	}

	return report, nil
}
