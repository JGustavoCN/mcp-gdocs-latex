package gdocs

import (
	"context"
	"fmt"
	"strings"
)

// ElementLayoutDetail descreve em detalhes a geometria e os estilos de um elemento do documento
type ElementLayoutDetail struct {
	Index              int     `json:"index"`
	StartIndex         int64   `json:"start_index"`
	EndIndex           int64   `json:"end_index"`
	Type               string  `json:"type"` // "PARAGRAPH", "TABLE", "SECTION_BREAK"
	StyleName          string  `json:"style_name,omitempty"`
	Alignment          string  `json:"alignment,omitempty"`
	LineSpacingPercent float64 `json:"line_spacing_percent,omitempty"` // 150 = 1,5; 100 = 1,0; 115 = 1,15
	IndentStartPt      float64 `json:"indent_start_pt,omitempty"`
	IndentFirstLinePt  float64 `json:"indent_first_line_pt,omitempty"`
	SpaceAbovePt       float64 `json:"space_above_pt,omitempty"`
	SpaceBelowPt       float64 `json:"space_below_pt,omitempty"`
	KeepWithNext       bool    `json:"keep_with_next,omitempty"`
	AvoidWidowOrphan   bool    `json:"avoid_widow_orphan,omitempty"`
	FontSizePt         float64 `json:"font_size_pt,omitempty"`
	FontFamily         string  `json:"font_family,omitempty"`
	Bold               bool    `json:"bold,omitempty"`
	Italic             bool    `json:"italic,omitempty"`
	HasPageBreak       bool    `json:"has_page_break,omitempty"`
	HasInlineObject    bool    `json:"has_inline_object,omitempty"`
	TextSnippet        string  `json:"text_snippet,omitempty"`
	TableRows          int     `json:"table_rows,omitempty"`
	TableCols          int     `json:"table_cols,omitempty"`
	EmptyCells         int     `json:"empty_cells,omitempty"`
}

// DocumentGeometryDetail armazena as dimensões da página e margens
type DocumentGeometryDetail struct {
	PageWidthPt    float64 `json:"page_width_pt"`
	PageHeightPt   float64 `json:"page_height_pt"`
	MarginTopPt    float64 `json:"margin_top_pt"`
	MarginTopCm    float64 `json:"margin_top_cm"`
	MarginLeftPt   float64 `json:"margin_left_pt"`
	MarginLeftCm   float64 `json:"margin_left_cm"`
	MarginBottomPt float64 `json:"margin_bottom_pt"`
	MarginBottomCm float64 `json:"margin_bottom_cm"`
	MarginRightPt  float64 `json:"margin_right_pt"`
	MarginRightCm  float64 `json:"margin_right_cm"`
}

// SectionLayoutDetail descreve seções e cabeçalhos
type SectionLayoutDetail struct {
	SectionIndex int    `json:"section_index"`
	StartIndex   int64  `json:"start_index"`
	EndIndex     int64  `json:"end_index"`
	HasHeader    bool   `json:"has_header"`
	HeaderID     string `json:"header_id,omitempty"`
	SectionType  string `json:"section_type,omitempty"`
}

// InspectLayoutReport consolida o Raio-X geométrico e estrutural dos elementos do Google Docs
type InspectLayoutReport struct {
	DocumentID       string                  `json:"document_id"`
	DocumentTitle    string                  `json:"document_title"`
	TotalElements    int                     `json:"total_elements"`
	StartIndex       int64                   `json:"filter_start_index"`
	EndIndex         int64                   `json:"filter_end_index"`
	DocumentGeometry *DocumentGeometryDetail `json:"document_geometry,omitempty"`
	Sections         []SectionLayoutDetail   `json:"sections,omitempty"`
	Elements         []ElementLayoutDetail   `json:"elements"`
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
		Sections:      []SectionLayoutDetail{},
		Elements:      []ElementLayoutDetail{},
	}

	const ptToCm = 2.54 / 72.0
	if doc.DocumentStyle != nil {
		geom := &DocumentGeometryDetail{}
		if doc.DocumentStyle.PageSize != nil {
			if doc.DocumentStyle.PageSize.Width != nil {
				geom.PageWidthPt = doc.DocumentStyle.PageSize.Width.Magnitude
			}
			if doc.DocumentStyle.PageSize.Height != nil {
				geom.PageHeightPt = doc.DocumentStyle.PageSize.Height.Magnitude
			}
		}
		if doc.DocumentStyle.MarginTop != nil {
			geom.MarginTopPt = doc.DocumentStyle.MarginTop.Magnitude
			geom.MarginTopCm = geom.MarginTopPt * ptToCm
		}
		if doc.DocumentStyle.MarginLeft != nil {
			geom.MarginLeftPt = doc.DocumentStyle.MarginLeft.Magnitude
			geom.MarginLeftCm = geom.MarginLeftPt * ptToCm
		}
		if doc.DocumentStyle.MarginBottom != nil {
			geom.MarginBottomPt = doc.DocumentStyle.MarginBottom.Magnitude
			geom.MarginBottomCm = geom.MarginBottomPt * ptToCm
		}
		if doc.DocumentStyle.MarginRight != nil {
			geom.MarginRightPt = doc.DocumentStyle.MarginRight.Magnitude
			geom.MarginRightCm = geom.MarginRightPt * ptToCm
		}
		report.DocumentGeometry = geom
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
			sb := el.SectionBreak
			sectionType := ""
			headerID := ""
			hasHeader := false
			if sb.SectionStyle != nil {
				sectionType = sb.SectionStyle.SectionType
				if sb.SectionStyle.DefaultHeaderId != "" {
					headerID = sb.SectionStyle.DefaultHeaderId
					hasHeader = true
				}
			}
			if sectionType == "" {
				sectionType = "NEXT_PAGE"
			}
			detail.TextSnippet = fmt.Sprintf("Quebra de Seção (%s)", sectionType)
			report.Sections = append(report.Sections, SectionLayoutDetail{
				SectionIndex: len(report.Sections) + 1,
				StartIndex:   el.StartIndex,
				EndIndex:     el.EndIndex,
				SectionType:  sectionType,
				HasHeader:    hasHeader,
				HeaderID:     headerID,
			})
		} else if el.Paragraph != nil {
			detail.Type = "PARAGRAPH"
			p := el.Paragraph
			if p.ParagraphStyle != nil {
				detail.StyleName = p.ParagraphStyle.NamedStyleType
				detail.Alignment = p.ParagraphStyle.Alignment
				detail.LineSpacingPercent = p.ParagraphStyle.LineSpacing
				detail.KeepWithNext = p.ParagraphStyle.KeepWithNext
				detail.AvoidWidowOrphan = p.ParagraphStyle.AvoidWidowAndOrphan
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
