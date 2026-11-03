package gdocs

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf16"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/latex"
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

// SkeletonItem representa um título ou tópico identificado no documento Google Docs
type SkeletonItem struct {
	Text         string `json:"text"`
	CurrentStyle string `json:"current_style"`
	TargetStyle  string `json:"target_style,omitempty"`
	Level        int    `json:"level,omitempty"`
	StartIndex   int64  `json:"start_index"`
	EndIndex     int64  `json:"end_index"`
}

// StructuralComparison condensa o confronto entre o Docs e o Template/LaTeX
type StructuralComparison struct {
	MatchedSections []string      `json:"matched_sections"`
	MissingInDoc    []string      `json:"missing_in_doc"`
	Anomalies       []AnomalyItem `json:"anomalies,omitempty"`
	Recommendations []string      `json:"recommendations,omitempty"`
}

// DocSkeletonReport representa o retorno unificado de raio-X e comparação estrutural
type DocSkeletonReport struct {
	DocumentTitle        string                `json:"document_title"`
	TotalElements        int                   `json:"total_elements"`
	DocSkeleton          []SkeletonItem        `json:"doc_skeleton"`
	TemplateSkeleton     []latex.LatexElement  `json:"template_skeleton,omitempty"`
	StructuralComparison *StructuralComparison `json:"structural_comparison,omitempty"`
}

// GetDocSkeleton retorna o Raio-X hierárquico do Google Docs unificado com a comparação do Template/LaTeX
func GetDocSkeleton(ctx context.Context, documentID string, referenceTemplateDir ...string) (string, error) {
	svc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	doc, err := svc.Documents.Get(documentID).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao obter documento '%s': %w", documentID, err)
	}

	if doc.Body == nil {
		return "{}", nil
	}

	report := &DocSkeletonReport{
		DocumentTitle: doc.Title,
		TotalElements: len(doc.Body.Content),
	}

	var anomalies []AnomalyItem
	var headingRecords []HeadingRecord

	// 1. Mapeia títulos existentes e detecta anomalias estruturais
	for _, el := range doc.Body.Content {
		if el.Paragraph != nil {
			p := el.Paragraph
			styleName := "NORMAL_TEXT"
			if p.ParagraphStyle != nil && p.ParagraphStyle.NamedStyleType != "" {
				styleName = p.ParagraphStyle.NamedStyleType
			}

			var textBuilder strings.Builder
			for _, pEl := range p.Elements {
				if pEl.TextRun != nil {
					textBuilder.WriteString(pEl.TextRun.Content)
				}
			}

			trimmed := strings.TrimSpace(textBuilder.String())
			if trimmed == "" {
				continue
			}

			// Detecção de comandos LaTeX vazados
			if matches := regexLatexLeakCmd.FindAllStringIndex(trimmed, -1); len(matches) > 0 {
				anomalies = append(anomalies, AnomalyItem{
					Type:        "LEAKED_COMMAND",
					Description: "Comando LaTeX vazado em texto puro",
					StartIndex:  el.StartIndex,
					EndIndex:    el.EndIndex,
					Snippet:     trimmed,
				})
			}

			isHeading := strings.Contains(styleName, "HEADING") || styleName == "TITLE"

			targetStyle := styleName
			level := 0
			if isHeading {
				if styleName == "TITLE" {
					level = 1
				} else if strings.HasPrefix(styleName, "HEADING_") {
					fmt.Sscanf(styleName, "HEADING_%d", &level)
				}
			} else {
				// Detecta se é um título não estilizado (em NORMAL_TEXT)
				if regexPrimaryHeading.MatchString(trimmed) {
					targetStyle = "HEADING_1"
					level = 1
				} else if regexSecondaryHeading.MatchString(trimmed) {
					targetStyle = "HEADING_2"
					level = 2
				} else if regexTertiaryHeading.MatchString(trimmed) {
					targetStyle = "HEADING_3"
					level = 3
				}

				if targetStyle != styleName {
					anomalies = append(anomalies, AnomalyItem{
						Type:        "UNSTYLED_HEADING",
						Description: fmt.Sprintf("Título '%s' está com estilo %s em vez de %s", trimmed, styleName, targetStyle),
						StartIndex:  el.StartIndex,
						EndIndex:    el.EndIndex,
						Snippet:     trimmed,
					})
				}
			}

			if isHeading || targetStyle != styleName {
				report.DocSkeleton = append(report.DocSkeleton, SkeletonItem{
					Text:         trimmed,
					CurrentStyle: styleName,
					TargetStyle:  targetStyle,
					Level:        level,
					StartIndex:   el.StartIndex,
					EndIndex:     el.EndIndex,
				})

				// Coleta para auditoria de continuidade da numeração progressiva
				cleanForNum := regexp.MustCompile(`(?i)\\(?:textit|textbf|texttt|text|textsc|emph)\b\{?`).ReplaceAllString(trimmed, "")
				cleanForNum = strings.ReplaceAll(cleanForNum, "}", "")
				cleanForNum = strings.TrimSpace(cleanForNum)
				if m := reSectionNumberExtractor.FindStringSubmatch(cleanForNum); len(m) > 1 {
					headingRecords = append(headingRecords, HeadingRecord{
						Text:        trimmed,
						Number:      m[1],
						Parts:       ParseSectionNumber(m[1]),
						StartIndex:  el.StartIndex,
						EndIndex:    el.EndIndex,
						StyleName:   styleName,
					})
				}
			}
		}
	}

	// Auditoria de quebras de sequência numérica progressiva (NBR 6024)
	headingAnomalies := AuditHeadingNumberingContinuity(headingRecords)
	for _, ha := range headingAnomalies {
		anomalies = append(anomalies, AnomalyItem{
			Type:        "HEADING_SEQUENCE_BREAK",
			Description: ha.Details,
			Snippet:     ha.Title,
		})
	}

	// 2. Confronto opcional com o Template / TCC de referência
	refDir := "resources/tcc"
	if len(referenceTemplateDir) > 0 && referenceTemplateDir[0] != "" {
		refDir = referenceTemplateDir[0]
	}

	refElements, errRef := latex.ScanTCCStructure(refDir)
	if errRef == nil && len(refElements) > 0 {
		report.TemplateSkeleton = refElements

		comp := &StructuralComparison{
			Anomalies: anomalies,
		}

		foundTitlesMap := make(map[string]bool)
		for _, fs := range report.DocSkeleton {
			cleanFound := strings.ToLower(strings.TrimSpace(fs.Text))
			foundTitlesMap[cleanFound] = true
		}

		for _, refEl := range refElements {
			if refEl.Type == latex.ElementChapter || refEl.Type == latex.ElementSection {
				cleanRef := strings.ToLower(strings.TrimSpace(refEl.Title))
				found := false
				for foundKey := range foundTitlesMap {
					if strings.Contains(foundKey, cleanRef) || strings.Contains(cleanRef, foundKey) {
						found = true
						comp.MatchedSections = append(comp.MatchedSections, refEl.Title)
						break
					}
				}
				if !found {
					comp.MissingInDoc = append(comp.MissingInDoc, fmt.Sprintf("%s (%s no arquivo %s)", refEl.Title, refEl.Type, refEl.File))
				}
			}
		}

		// Recomendações inteligentes
		if len(anomalies) > 0 {
			comp.Recommendations = append(comp.Recommendations,
				fmt.Sprintf("Executar 'normalize_doc_structure' para promover títulos e limpar %d anomalias.", len(anomalies)))
		}
		if len(comp.MissingInDoc) > 0 {
			comp.Recommendations = append(comp.Recommendations,
				fmt.Sprintf("Usar 'sync_section_to_doc' para injetar seções faltantes (%d pendentes).", len(comp.MissingInDoc)))
		}

		report.StructuralComparison = comp
	}

	b, err := json.MarshalIndent(report, "", "  ")
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
			sb.WriteString("\n\n[SECTION_BREAK: NEXT_PAGE]\n\n")
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
		if elem.PageBreak != nil {
			parts = append(parts, "\n[PAGE_BREAK]\n")
			continue
		}
		if elem.InlineObjectElement != nil {
			parts = append(parts, "[INLINE_OBJECT]")
			continue
		}
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
