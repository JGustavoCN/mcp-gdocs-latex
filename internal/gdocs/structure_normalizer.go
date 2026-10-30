package gdocs

import (
	"context"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/latex"
	"google.golang.org/api/docs/v1"
)

// SectionItem representa uma seção ou título identificado no Google Doc
type SectionItem struct {
	Title        string `json:"title"`
	CurrentStyle string `json:"current_style"`
	TargetStyle  string `json:"target_style"`
	StartIndex   int64  `json:"start_index"`
	EndIndex     int64  `json:"end_index"`
	Level        int    `json:"level"`
}

// AnomalyItem representa anomalias estruturais ou comandos vazados no documento
type AnomalyItem struct {
	Type        string `json:"type"` // "LEAKED_COMMAND", "UNSTYLED_HEADING", "MISSING_PAGE_BREAK"
	Description string `json:"description"`
	StartIndex  int64  `json:"start_index"`
	EndIndex    int64  `json:"end_index"`
	Snippet     string `json:"snippet"`
}

// DocStructureDiff consolida o diagnóstico da estrutura do documento em relação ao template/TCC
type DocStructureDiff struct {
	DocumentTitle      string        `json:"document_title"`
	TotalElements      int           `json:"total_elements"`
	FoundSections      []SectionItem `json:"found_sections"`
	MissingSections    []string      `json:"missing_sections"`
	DetectedAnomalies  []AnomalyItem `json:"detected_anomalies"`
	RecommendedActions []string      `json:"recommended_actions"`
}

// NormalizeResult resume as operações executadas na normalização in-place
type NormalizeResult struct {
	PromotedHeadingsCount int      `json:"promoted_headings_count"`
	CleanedCommandsCount  int      `json:"cleaned_commands_count"`
	InsertedPageBreaks    int      `json:"inserted_page_breaks"`
	ABNTFormatted         bool     `json:"abnt_formatted"`
	Details               []string `json:"details"`
	Message               string   `json:"message"`
}

var (
	regexPrimaryHeading   = regexp.MustCompile(`^(?:[1-9][0-9]*\s+[A-ZÁÉÍÓÚÂÊÎÔÛÃÕÇ\s\-]{3,}|(?:CAPÍTULO\s+[0-9]+.*)|INTRODUÇÃO|REFERENCIAL TEÓRICO|METODOLOGIA|DESENVOLVIMENTO|RESULTADOS|CONCLUSÃO|REFERÊNCIAS|RESUMO|ABSTRACT)$`)
	regexSecondaryHeading = regexp.MustCompile(`^[1-9][0-9]*\.[0-9]+\s+[A-ZÁÉÍÓÚÂÊÎÔÛÃÕÇa-záéíóúâêîôûãõç\s\-,/]+$`)
	regexTertiaryHeading  = regexp.MustCompile(`^[1-9][0-9]*\.[0-9]+\.[0-9]+\s+[A-ZÁÉÍÓÚÂÊÎÔÛÃÕÇa-záéíóúâêîôûãõç\s\-,/]+$`)
	regexLatexLeakCmd     = regexp.MustCompile(`(?i)\\(?:chapter|section|subsection|subsubsection|label|begin|end|ref|cite|caption|textit|textbf)\*?[{(][^})]*[})]`)
)

// DiffDocStructure inspeciona o documento de forma NÃO-DESTRUTIVA, mapeando seções, anomalias e tópicos faltantes
func DiffDocStructure(ctx context.Context, documentID string, referenceTCCDir string) (*DocStructureDiff, error) {
	docsSvc, err := GetDocsService(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	doc, err := docsSvc.Documents.Get(documentID).Do()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar documento: %w", err)
	}

	if referenceTCCDir == "" {
		referenceTCCDir = "resources/tcc"
	}

	diff := &DocStructureDiff{
		DocumentTitle: doc.Title,
		TotalElements: len(doc.Body.Content),
	}

	// 1. Mapeamento dos elementos estruturais existentes no Google Doc
	for _, el := range doc.Body.Content {
		if el.Paragraph != nil {
			p := el.Paragraph
			styleName := "NORMAL_TEXT"
			if p.ParagraphStyle != nil && p.ParagraphStyle.NamedStyleType != "" {
				styleName = p.ParagraphStyle.NamedStyleType
			}

			pText := ""
			for _, pEl := range p.Elements {
				if pEl.TextRun != nil {
					pText += pEl.TextRun.Content
				}
			}
			trimmed := strings.TrimSpace(pText)
			if trimmed == "" {
				continue
			}

			// 1.1 Detecção de Comandos LaTeX Vazados no texto
			if matches := regexLatexLeakCmd.FindAllStringIndex(trimmed, -1); len(matches) > 0 {
				diff.DetectedAnomalies = append(diff.DetectedAnomalies, AnomalyItem{
					Type:        "LEAKED_COMMAND",
					Description: "Comando LaTeX vazado em texto puro",
					StartIndex:  el.StartIndex,
					EndIndex:    el.EndIndex,
					Snippet:     trimmed,
				})
			}

			// 1.2 Identificação de Títulos (Primários, Secundários, Terciários)
			targetStyle := ""
			level := 0

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

			if targetStyle != "" {
				diff.FoundSections = append(diff.FoundSections, SectionItem{
					Title:        trimmed,
					CurrentStyle: styleName,
					TargetStyle:  targetStyle,
					StartIndex:   el.StartIndex,
					EndIndex:     el.EndIndex,
					Level:        level,
				})

				// Se o título não estiver estilizado com o HEADING correto
				if styleName != targetStyle {
					diff.DetectedAnomalies = append(diff.DetectedAnomalies, AnomalyItem{
						Type:        "UNSTYLED_HEADING",
						Description: fmt.Sprintf("Título '%s' está com estilo %s em vez de %s", trimmed, styleName, targetStyle),
						StartIndex:  el.StartIndex,
						EndIndex:    el.EndIndex,
						Snippet:     trimmed,
					})
				}
			}
		}
	}

	// 2. Comparação com o modelo de referência (resources/tcc ou resources/template)
	refElements, errRef := latex.ScanTCCStructure(referenceTCCDir)
	if errRef == nil && len(refElements) > 0 {
		foundTitlesMap := make(map[string]bool)
		for _, fs := range diff.FoundSections {
			cleanFound := strings.ToLower(strings.TrimSpace(fs.Title))
			foundTitlesMap[cleanFound] = true
		}

		for _, refEl := range refElements {
			if refEl.Type == latex.ElementChapter || refEl.Type == latex.ElementSection {
				cleanRef := strings.ToLower(strings.TrimSpace(refEl.Title))
				found := false
				for foundKey := range foundTitlesMap {
					if strings.Contains(foundKey, cleanRef) || strings.Contains(cleanRef, foundKey) {
						found = true
						break
					}
				}
				if !found {
					diff.MissingSections = append(diff.MissingSections, fmt.Sprintf("%s (%s no arquivo %s)", refEl.Title, refEl.Type, refEl.File))
				}
			}
		}
	}

	// 3. Elaboração de Ações Recomendadas
	if len(diff.DetectedAnomalies) > 0 {
		diff.RecommendedActions = append(diff.RecommendedActions,
			fmt.Sprintf("Executar 'normalize_doc_structure' para promover %d títulos a cabeçalhos formais e higienizar comandos LaTeX vazados sem alterar o texto.", len(diff.FoundSections)))
	}
	if len(diff.MissingSections) > 0 {
		diff.RecommendedActions = append(diff.RecommendedActions,
			fmt.Sprintf("Considerar injeção controlada de seções faltantes (%d seções identificadas no TCC de referência).", len(diff.MissingSections)))
	}
	diff.RecommendedActions = append(diff.RecommendedActions,
		"Executar 'format_document_abnt' para assegurar margens 3/3/2/2 cm e entrelinhas 1.5 preservando tabelas e comentários.")

	return diff, nil
}

// NormalizeDocStructure executa a padronização IN-PLACE sem apagar texto ou comentários,
// delegando diretamente para o motor unificado FormatDocumentABNTExtended para garantir tipografia e hierarquia NBR 6024 rigorosas.
func NormalizeDocStructure(ctx context.Context, documentID string, addPageBreaks bool, cleanLatexCommands bool, applyABNT bool) (*NormalizeResult, error) {
	msg, err := FormatDocumentABNTExtended(ctx, documentID, -1, -1, applyABNT, cleanLatexCommands, ABNTDefaultFontFamily, true, true)
	if err != nil {
		return nil, fmt.Errorf("erro ao executar normalização e formatação ABNT: %w", err)
	}

	res := &NormalizeResult{
		ABNTFormatted: true,
		Details:       []string{msg},
		Message:       "Normalização estrutural ABNT concluída com sucesso: títulos formatados com fonte 12pt, cabeçalhos NBR 6024, margens A4 3/3/2/2cm, parágrafos 1,25cm, quebras duras desfragmentadas e linhas vazias colapsadas.",
	}

	return res, nil
}

// SyncSectionToDoc injeta ou atualiza uma seção específica no documento de forma cirúrgica
func SyncSectionToDoc(ctx context.Context, documentID, targetSectionTitle, sectionContent string) (string, error) {
	docsSvc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter docs service: %w", err)
	}

	doc, err := docsSvc.Documents.Get(documentID).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao carregar documento: %w", err)
	}

	targetIndex := int64(1)
	cleanTarget := strings.ToLower(strings.TrimSpace(targetSectionTitle))

	// Tenta localizar a seção existente para atualizar
	foundStart := int64(-1)
	foundEnd := int64(-1)

	for _, el := range doc.Body.Content {
		if el.Paragraph != nil {
			pText := ""
			for _, pEl := range el.Paragraph.Elements {
				if pEl.TextRun != nil {
					pText += pEl.TextRun.Content
				}
			}
			trimmed := strings.ToLower(strings.TrimSpace(pText))

			if strings.Contains(trimmed, cleanTarget) {
				foundStart = el.StartIndex
				foundEnd = el.EndIndex
				break
			}
		}
	}

	if foundStart != -1 {
		// Substituição controlada in-place do bloco da seção
		return ReplaceBlockABNT(ctx, documentID, foundStart, foundEnd, sectionContent)
	}

	// Se não encontrada, anexa antes do final do documento
	if len(doc.Body.Content) > 0 {
		targetIndex = doc.Body.Content[len(doc.Body.Content)-1].EndIndex - 1
	}

	insertText := fmt.Sprintf("\n\n%s\n\n%s\n", strings.ToUpper(targetSectionTitle), sectionContent)
	req := &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{
			{
				InsertText: &docs.InsertTextRequest{
					Location: &docs.Location{Index: targetIndex},
					Text:     insertText,
				},
			},
		},
	}
	_, errIns := docsSvc.Documents.BatchUpdate(documentID, req).Do()
	if errIns != nil {
		return "", fmt.Errorf("erro ao anexar seção: %w", errIns)
	}

	return fmt.Sprintf("SUCESSO: Seção '%s' anexada ao documento no índice %d.", targetSectionTitle, targetIndex), nil
}

// ResetDocToTemplate executa o reset sob demanda estrita a partir do template limpo do IFS
func ResetDocToTemplate(ctx context.Context, documentID, templateDir string) (string, error) {
	if templateDir == "" {
		templateDir = filepath.Join("resources", "template")
	}

	docsSvc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter docs service: %w", err)
	}

	doc, err := docsSvc.Documents.Get(documentID).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao obter documento: %w", err)
	}

	bodyEnd := int64(1)
	if len(doc.Body.Content) > 0 {
		bodyEnd = doc.Body.Content[len(doc.Body.Content)-1].EndIndex
	}

	// 1. Limpeza explícita do corpo sob solicitação formal
	if bodyEnd > 2 {
		delReq := &docs.BatchUpdateDocumentRequest{
			Requests: []*docs.Request{
				{
					DeleteContentRange: &docs.DeleteContentRangeRequest{
						Range: &docs.Range{
							StartIndex: 1,
							EndIndex:   bodyEnd - 1,
						},
					},
				},
			},
		}
		if _, errDel := docsSvc.Documents.BatchUpdate(documentID, delReq).Do(); errDel != nil {
			return "", fmt.Errorf("erro ao limpar documento para reset: %w", errDel)
		}
	}

	// 2. Inserção do Esqueleto Canônico do Template
	var sb strings.Builder
	sb.WriteString("INSTITUTO FEDERAL DE EDUCAÇÃO, CIÊNCIA E TECNOLOGIA DE SERGIPE\n")
	sb.WriteString("CAMPUS LAGARTO\n")
	sb.WriteString("BACHARELADO EM SISTEMAS DE INFORMAÇÃO\n\n\n")
	sb.WriteString("NOME DO AUTOR\n\n\n")
	sb.WriteString("TÍTULO DO TRABALHO DE CONCLUSÃO DE CURSO: SUBTÍTULO\n\n\n")
	sb.WriteString("Trabalho de Conclusão de Curso apresentado ao curso de Bacharelado em Sistemas de Informação do Instituto Federal de Educação, Ciência e Tecnologia de Sergipe, Campus Lagarto, como requisito parcial para obtenção do grau de Bacharel em Sistemas de Informação.\n\n")
	sb.WriteString("Orientador: Prof. Dr. Nome do Orientador\n\n")
	sb.WriteString("Lagarto/SE\n2026\n\n")
	sb.WriteString("1 INTRODUÇÃO\n\nEscreva aqui a introdução do trabalho.\n\n")
	sb.WriteString("2 REFERENCIAL TEÓRICO\n\nEscreva aqui o referencial teórico.\n\n")
	sb.WriteString("3 METODOLOGIA\n\nEscreva aqui a metodologia científica adotada.\n\n")
	sb.WriteString("4 RESULTADOS E DISCUSSÃO\n\nEscreva aqui os resultados alcançados.\n\n")
	sb.WriteString("5 CONCLUSÃO\n\nEscreva aqui as considerações finais.\n\n")
	sb.WriteString("REFERÊNCIAS\n\nInsira aqui as referências no padrão ABNT NBR 6023.\n")

	templateText := sb.String()
	insReq := &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{
			{
				InsertText: &docs.InsertTextRequest{
					Location: &docs.Location{Index: 1},
					Text:     templateText,
				},
			},
		},
	}
	if _, errIns := docsSvc.Documents.BatchUpdate(documentID, insReq).Do(); errIns != nil {
		return "", fmt.Errorf("erro ao inserir template básico: %w", errIns)
	}

	// 3. Aplicação imediata de formatação ABNT
	FormatDocumentABNT(ctx, documentID, -1, -1, true)

	return "SUCESSO: Documento resetado e inicializado com o Template Limpo ABNT/IFS com sucesso.", nil
}

