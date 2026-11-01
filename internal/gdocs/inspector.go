package gdocs

import (
	"context"
	"fmt"
	"strings"
)

// ElementLayoutDetail descreve em detalhes a geometria e os estilos de um elemento do documento
type ElementLayoutDetail struct {
	Index             int     `json:"index"`
	StartIndex        int64   `json:"start_index"`
	EndIndex          int64   `json:"end_index"`
	Type              string  `json:"type"` // "PARAGRAPH", "TABLE", "SECTION_BREAK"
	StyleName         string  `json:"style_name,omitempty"`
	Alignment         string  `json:"alignment,omitempty"`
	IndentStartPt     float64 `json:"indent_start_pt,omitempty"`
	IndentFirstLinePt float64 `json:"indent_first_line_pt,omitempty"`
	SpaceAbovePt      float64 `json:"space_above_pt,omitempty"`
	SpaceBelowPt      float64 `json:"space_below_pt,omitempty"`
	FontSizePt        float64 `json:"font_size_pt,omitempty"`
	FontFamily        string  `json:"font_family,omitempty"`
	Bold              bool    `json:"bold,omitempty"`
	Italic            bool    `json:"italic,omitempty"`
	HasPageBreak      bool    `json:"has_page_break,omitempty"`
	HasInlineObject   bool    `json:"has_inline_object,omitempty"`
	TextSnippet       string  `json:"text_snippet,omitempty"`
	TableRows         int     `json:"table_rows,omitempty"`
	TableCols         int     `json:"table_cols,omitempty"`
	EmptyCells        int     `json:"empty_cells,omitempty"`
}

// InspectLayoutReport consolida o Raio-X geométrico e estrutural dos elementos do Google Docs
type InspectLayoutReport struct {
	DocumentID    string                `json:"document_id"`
	DocumentTitle string                `json:"document_title"`
	TotalElements int                   `json:"total_elements"`
	StartIndex    int64                 `json:"filter_start_index"`
	EndIndex      int64                 `json:"filter_end_index"`
	Elements      []ElementLayoutDetail `json:"elements"`
}

// InspectDocLayout inspeciona as propriedades de baixo nível (recuos, alinhamento, quebras de página, tabelas)
func InspectDocLayout(ctx context.Context, documentID string, startIndex, endIndex int64, limit int) (*InspectLayoutReport, error) {
	svc, err := GetDocsService(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter docs service: %w", err)
	}

	doc, err := svc.Documents.Get(documentID).Do()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar documento: %w", err)
	}

	report := &InspectLayoutReport{
		DocumentID:    documentID,
		DocumentTitle: doc.Title,
		TotalElements: len(doc.Body.Content),
		StartIndex:    startIndex,
		EndIndex:      endIndex,
		Elements:      []ElementLayoutDetail{},
	}

	count := 0
	for i, el := range doc.Body.Content {
		if startIndex >= 0 && el.EndIndex <= startIndex {
			continue
		}
		if endIndex > 0 && el.StartIndex >= endIndex {
			break
		}
		if limit > 0 && count >= limit {
			break
		}

		detail := ElementLayoutDetail{
			Index:      i,
			StartIndex: el.StartIndex,
			EndIndex:   el.EndIndex,
		}

		if el.Table != nil {
			detail.Type = "TABLE"
			rows := len(el.Table.TableRows)
			cols := 0
			if rows > 0 {
				cols = len(el.Table.TableRows[0].TableCells)
			}
			emptyCells := 0
			for _, r := range el.Table.TableRows {
				for _, c := range r.TableCells {
					cellTxt := ""
					for _, cEl := range c.Content {
						if cEl.Paragraph != nil {
							for _, pEl := range cEl.Paragraph.Elements {
								if pEl.TextRun != nil {
									cellTxt += pEl.TextRun.Content
								}
							}
						}
					}
					if strings.TrimSpace(cellTxt) == "" {
						emptyCells++
					}
				}
			}
			detail.TableRows = rows
			detail.TableCols = cols
			detail.EmptyCells = emptyCells
			detail.TextSnippet = fmt.Sprintf("Tabela Nativa %dx%d (vazias: %d/%d)", rows, cols, emptyCells, rows*cols)
		} else if el.SectionBreak != nil {
			detail.Type = "SECTION_BREAK"
			detail.TextSnippet = "Quebra de Seção"
		} else if el.Paragraph != nil {
			detail.Type = "PARAGRAPH"
			p := el.Paragraph
			if p.ParagraphStyle != nil {
				detail.StyleName = p.ParagraphStyle.NamedStyleType
				detail.Alignment = p.ParagraphStyle.Alignment
				if p.ParagraphStyle.IndentStart != nil {
					detail.IndentStartPt = p.ParagraphStyle.IndentStart.Magnitude
				}
				if p.ParagraphStyle.IndentFirstLine != nil {
					detail.IndentFirstLinePt = p.ParagraphStyle.IndentFirstLine.Magnitude
				}
				if p.ParagraphStyle.SpaceAbove != nil {
					detail.SpaceAbovePt = p.ParagraphStyle.SpaceAbove.Magnitude
				}
				if p.ParagraphStyle.SpaceBelow != nil {
					detail.SpaceBelowPt = p.ParagraphStyle.SpaceBelow.Magnitude
				}
			}

			var txtBuilder strings.Builder
			for _, pEl := range p.Elements {
				if pEl.TextRun != nil {
					tr := pEl.TextRun
					txtBuilder.WriteString(tr.Content)
					if tr.TextStyle != nil {
						if tr.TextStyle.FontSize != nil && detail.FontSizePt == 0 {
							detail.FontSizePt = tr.TextStyle.FontSize.Magnitude
						}
						if tr.TextStyle.WeightedFontFamily != nil && detail.FontFamily == "" {
							detail.FontFamily = tr.TextStyle.WeightedFontFamily.FontFamily
						}
						if tr.TextStyle.Bold {
							detail.Bold = true
						}
						if tr.TextStyle.Italic {
							detail.Italic = true
						}
					}
				}
				if pEl.PageBreak != nil {
					detail.HasPageBreak = true
				}
				if pEl.InlineObjectElement != nil {
					detail.HasInlineObject = true
				}
			}

			fullTxt := strings.TrimSpace(txtBuilder.String())
			if len(fullTxt) > 100 {
				detail.TextSnippet = fullTxt[:100] + "..."
			} else {
				detail.TextSnippet = fullTxt
			}
		}

		report.Elements = append(report.Elements, detail)
		count++
	}

	return report, nil
}
