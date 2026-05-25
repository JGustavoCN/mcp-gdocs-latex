package gdocs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf16"

	"google.golang.org/api/docs/v1"
)

// headingMap mapeia os estilos de parágrafo do Google Docs para prefixos Markdown.
var headingMap = map[string]string{
	"HEADING_1": "# ",
	"HEADING_2": "## ",
	"HEADING_3": "### ",
	"HEADING_4": "#### ",
	"HEADING_5": "##### ",
	"HEADING_6": "###### ",
	"TITLE":     "# ",
	"SUBTITLE":  "## ",
}

// SkeletonItem representa um título ou trecho destacado no documento.
type SkeletonItem struct {
	Text       string `json:"text"`
	Style      string `json:"style"`
	Level      int    `json:"level,omitempty"`
	StartIndex int64  `json:"start_index"`
	EndIndex   int64  `json:"end_index"`
}

// GetDocSkeleton retorna um "Raio-X" do documento (apenas títulos e highlights).
func GetDocSkeleton(ctx context.Context, documentID string) (string, error) {
	svc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	doc, err := svc.Documents.Get(documentID).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao obter documento '%s': %w", documentID, err)
	}

	if doc.Body == nil {
		return "[]", nil
	}

	var skeleton []SkeletonItem

	var traverse func(elements []*docs.StructuralElement)
	traverse = func(elements []*docs.StructuralElement) {
		for _, el := range elements {
			if el.Paragraph != nil {
				style := ""
				if el.Paragraph.ParagraphStyle != nil {
					style = el.Paragraph.ParagraphStyle.NamedStyleType
				}

				isHeading := strings.Contains(style, "HEADING") || style == "TITLE"
				hasHighlight := false

				var textBuilder strings.Builder
				for _, pEl := range el.Paragraph.Elements {
					if pEl.TextRun != nil {
						textBuilder.WriteString(pEl.TextRun.Content)
						if pEl.TextRun.TextStyle != nil && pEl.TextRun.TextStyle.BackgroundColor != nil {
							hasHighlight = true
						}
					}
				}

				text := strings.TrimSpace(textBuilder.String())
				if text != "" && (isHeading || hasHighlight) {
					displayStyle := style
					level := 0
					if !isHeading && hasHighlight {
						displayStyle = "HIGHLIGHT"
					} else if isHeading {
						if style == "TITLE" {
							level = 1
						} else if strings.HasPrefix(style, "HEADING_") {
							fmt.Sscanf(style, "HEADING_%d", &level)
						}
					}
					skeleton = append(skeleton, SkeletonItem{
						Text:       text,
						Style:      displayStyle,
						Level:      level,
						StartIndex: el.StartIndex,
						EndIndex:   el.EndIndex,
					})
				}
			} else if el.Table != nil {
				for _, row := range el.Table.TableRows {
					for _, cell := range row.TableCells {
						traverse(cell.Content)
					}
				}
			}
		}
	}

	traverse(doc.Body.Content)

	b, err := json.MarshalIndent(skeleton, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao formatar JSON do skeleton: %w", err)
	}

	return string(b), nil
}

// ReadDocContent lê o conteúdo de um Google Doc.
// Se startIndex e endIndex forem -1, lê o documento inteiro.
// Caso contrário, retorna estritamente o conteúdo dentro desses índices,
// mantendo a estrutura Markdown.
func ReadDocContent(ctx context.Context, documentID string, startIndex, endIndex int64) (string, error) {
	svc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	doc, err := svc.Documents.Get(documentID).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao obter documento '%s': %w", documentID, err)
	}

	title := doc.Title
	if title == "" {
		title = "Sem título"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== DOCUMENTO: %s ===\n", title))
	sb.WriteString(fmt.Sprintf("ID: %s\n", documentID))
	if startIndex != -1 && endIndex != -1 {
		sb.WriteString(fmt.Sprintf("Filtro de Leitura: Índices %d a %d\n", startIndex, endIndex))
	}
	sb.WriteString(strings.Repeat("=", 50) + "\n\n")

	if doc.Body != nil {
		extractStructuralElements(&sb, doc.Body.Content, startIndex, endIndex)
	}

	return sb.String(), nil
}

// extractStructuralElements percorre recursivamente os elementos.
func extractStructuralElements(sb *strings.Builder, elements []*docs.StructuralElement, startIdx, endIdx int64) {
	for _, element := range elements {
		// Pula se estiver completamente fora do range (quando ativado)
		if startIdx != -1 && endIdx != -1 {
			if element.EndIndex <= startIdx || element.StartIndex >= endIdx {
				continue
			}
		}

		switch {
		case element.Paragraph != nil:
			processParagraph(sb, element.Paragraph, startIdx, endIdx)
		case element.Table != nil:
			processTable(sb, element.Table, startIdx, endIdx)
		case element.SectionBreak != nil:
			sb.WriteString("\n---\n")
		}
	}
}

// processParagraph extrai o texto de um parágrafo aplicando formatação Markdown.
func processParagraph(sb *strings.Builder, paragraph *docs.Paragraph, startIdx, endIdx int64) {
	namedStyle := ""
	if paragraph.ParagraphStyle != nil {
		namedStyle = paragraph.ParagraphStyle.NamedStyleType
	}
	prefix := headingMap[namedStyle]

	var parts []string
	for _, elem := range paragraph.Elements {
		if elem.TextRun == nil {
			continue
		}

		// Filtragem por índice
		if startIdx != -1 && endIdx != -1 {
			if elem.EndIndex <= startIdx || elem.StartIndex >= endIdx {
				continue
			}
		}

		content := elem.TextRun.Content

		// Recorte de string caso o elemento cruze as bordas do startIdx/endIdx
		if startIdx != -1 && endIdx != -1 {
			u16 := utf16.Encode([]rune(content))

			localStart := int64(0)
			if startIdx > elem.StartIndex {
				localStart = startIdx - elem.StartIndex
			}

			localEnd := int64(len(u16))
			if endIdx < elem.EndIndex {
				localEnd = endIdx - elem.StartIndex
			}

			if localStart >= 0 && localEnd <= int64(len(u16)) && localStart <= localEnd {
				content = string(utf16.Decode(u16[localStart:localEnd]))
			} else {
				content = ""
			}
		}

		if content == "" {
			continue
		}

		// Aplica formatação
		if elem.TextRun.TextStyle != nil {
			trimmed := strings.TrimSpace(content)
			if elem.TextRun.TextStyle.Bold && trimmed != "" {
				content = strings.Replace(content, trimmed, "**"+trimmed+"**", 1)
			}
			if elem.TextRun.TextStyle.Italic && trimmed != "" {
				content = strings.Replace(content, trimmed, "*"+trimmed+"*", 1)
			}
		}

		parts = append(parts, content)
	}

	line := strings.Join(parts, "")

	if prefix != "" && strings.TrimSpace(line) != "" {
		line = prefix + strings.TrimSpace(line) + "\n"
	}

	sb.WriteString(line)
}

// processTable processa tabela mantendo filtro de índices.
func processTable(sb *strings.Builder, table *docs.Table, startIdx, endIdx int64) {
	sb.WriteString("\n[TABELA]\n")
	for _, row := range table.TableRows {
		var cells []string
		for _, cell := range row.TableCells {
			var cellSb strings.Builder
			extractStructuralElements(&cellSb, cell.Content, startIdx, endIdx)
			cells = append(cells, strings.TrimSpace(cellSb.String()))
		}
		sb.WriteString(strings.Join(cells, " | ") + "\n")
	}
	sb.WriteString("[/TABELA]\n\n")
}
