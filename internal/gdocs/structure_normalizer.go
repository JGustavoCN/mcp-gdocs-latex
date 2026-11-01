package gdocs

import (
	"context"
	"fmt"
	"log"
	"os"
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
	msg, err := FormatDocumentABNTExtended(ctx, documentID, -1, -1, applyABNT, cleanLatexCommands, ABNTDefaultFontFamily, true, true, false)
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
	matchIdx := -1
	targetLevel := 1

	for i, el := range doc.Body.Content {
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
				matchIdx = i
				if el.Paragraph.ParagraphStyle != nil {
					switch el.Paragraph.ParagraphStyle.NamedStyleType {
					case "HEADING_1":
						targetLevel = 1
					case "HEADING_2":
						targetLevel = 2
					case "HEADING_3":
						targetLevel = 3
					default:
						if regexPrimaryHeading.MatchString(pText) {
							targetLevel = 1
						} else if regexSecondaryHeading.MatchString(pText) {
							targetLevel = 2
						} else if regexTertiaryHeading.MatchString(pText) {
							targetLevel = 3
						}
					}
				}
				break
			}
		}
	}

	if foundStart != -1 {
		// Encontra onde a seção termina: início do próximo título de nível igual ou superior
		bodyEnd := doc.Body.Content[len(doc.Body.Content)-1].EndIndex - 1
		foundEnd = bodyEnd
		for j := matchIdx + 1; j < len(doc.Body.Content); j++ {
			nextEl := doc.Body.Content[j]
			if nextEl.Paragraph != nil {
				lvl := 0
				if nextEl.Paragraph.ParagraphStyle != nil {
					switch nextEl.Paragraph.ParagraphStyle.NamedStyleType {
					case "HEADING_1":
						lvl = 1
					case "HEADING_2":
						lvl = 2
					case "HEADING_3":
						lvl = 3
					}
				}
				txt := ""
				for _, pEl := range nextEl.Paragraph.Elements {
					if pEl.TextRun != nil {
						txt += pEl.TextRun.Content
					}
				}
				trimTxt := strings.TrimSpace(txt)
				if lvl == 0 {
					if regexPrimaryHeading.MatchString(trimTxt) {
						lvl = 1
					} else if regexSecondaryHeading.MatchString(trimTxt) {
						lvl = 2
					} else if regexTertiaryHeading.MatchString(trimTxt) {
						lvl = 3
					}
				}
				if lvl > 0 && lvl <= targetLevel {
					foundEnd = nextEl.StartIndex
					break
				}
			}
		}

		// Substituição controlada in-place do bloco da seção
		res, err := ReplaceBlockABNT(ctx, documentID, foundStart, foundEnd, sectionContent)
		if err != nil {
			return "", err
		}
		_, _ = FormatDocumentABNT(ctx, documentID, 0, 0, false, true)
		return res, nil
	}

	// Se não encontrada, anexa antes das Referências (se existirem) ou no final do documento
	if len(doc.Body.Content) > 0 {
		targetIndex = doc.Body.Content[len(doc.Body.Content)-1].EndIndex - 1
	}
	refIndex := int64(-1)
	for _, el := range doc.Body.Content {
		if el.Paragraph != nil {
			txt := ""
			for _, pEl := range el.Paragraph.Elements {
				if pEl.TextRun != nil {
					txt += pEl.TextRun.Content
				}
			}
			t := strings.ToUpper(strings.TrimSpace(txt))
			if strings.Contains(t, "REFERÊNCIAS") || strings.Contains(t, "REFERENCIAS") {
				refIndex = el.StartIndex
				break
			}
		}
	}
	if refIndex > 0 && !strings.Contains(cleanTarget, "referência") && !strings.Contains(cleanTarget, "referencia") && !strings.Contains(cleanTarget, "apêndice") && !strings.Contains(cleanTarget, "anexo") {
		targetIndex = refIndex
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

	_, _ = FormatDocumentABNT(ctx, documentID, 0, 0, false, true)
	return fmt.Sprintf("SUCESSO: Seção '%s' anexada ao documento no índice %d.", targetSectionTitle, targetIndex), nil
}

// ResetDocToTemplate executa a inicialização ou reset do documento com o Template Oficial ABNT/IFS,
// incluindo a Capa canônica com a logo oficial do IFS (3.2cm), Folha de Rosto e capítulos canônicos.
// Se documentID for vazio ou "new", cria um novo Google Doc no Drive do usuário (ou orienta o compartilhamento se cota da Service Account for 0).
func ResetDocToTemplate(ctx context.Context, documentID, templateDir string, metadataPath string, customLogoPath ...string) (string, error) {
	docsSvc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter docs service: %w", err)
	}

	metaFile := metadataPath
	if metaFile == "" {
		metaFile = filepath.Join("resources", "tcc", "metadados.tex")
	}
	meta, errMeta := latex.ParseTCCMetadata(metaFile)
	if errMeta != nil {
		log.Printf("[WARN] Falha ao ler metadados.tex (%v), utilizando configurações padrão", errMeta)
	}

	isNew := documentID == "" || strings.ToLower(strings.TrimSpace(documentID)) == "new"
	if isNew {
		docTitle := "TCC - " + meta.Titulo
		if len(docTitle) > 100 {
			docTitle = docTitle[:100]
		}
		newDoc, errCreate := docsSvc.Documents.Create(&docs.Document{Title: docTitle}).Do()
		if errCreate != nil {
			return "", fmt.Errorf("não foi possível criar um documento vazio autonomamente via Service Account (%v).\n"+
				"Solução: Crie um documento no seu Google Docs e compartilhe-o com o e-mail da sua Service Account como Editor, fornecendo o ID ou URL para inicialização do template", errCreate)
		}
		documentID = newDoc.DocumentId
		log.Printf("[INFO] Novo documento Google Docs criado com sucesso: ID %s", documentID)
	}

	logoPath := ""
	if len(customLogoPath) > 0 && customLogoPath[0] != "" {
		logoPath = customLogoPath[0]
	} else if templateDir != "" {
		possibleLogo := filepath.Join(templateDir, "logo_ifs.png")
		if _, errStat := os.Stat(possibleLogo); errStat == nil {
			logoPath = possibleLogo
		}
	}
	if logoPath == "" {
		if _, errAssets := os.Stat(filepath.Join("assets", "logo_ifs.png")); errAssets == nil {
			logoPath = filepath.Join("assets", "logo_ifs.png")
		} else if _, errRoot := os.Stat("logo_ifs.png"); errRoot == nil {
			logoPath = "logo_ifs.png"
		} else {
			logoPath = filepath.Join("resources", "tcc", "logo_ifs.png")
		}
	}

	errGen := GenerateABNTIFSCoverAndTitlePage(ctx, documentID, meta, logoPath, false)
	if errGen != nil {
		return "", fmt.Errorf("erro ao gerar template oficial ABNT/IFS: %w", errGen)
	}

	action := "resetado e inicializado"
	if isNew {
		action = "criado do zero"
	}

	return fmt.Sprintf("SUCESSO: Documento %s com o Template Oficial ABNT/IFS (com Logo do IFS 3.2cm centralizada, Capa, Folha de Rosto e Capítulos NBR 14724/6024).\nID: %s\nURL: https://docs.google.com/document/d/%s/edit", action, documentID, documentID), nil
}

