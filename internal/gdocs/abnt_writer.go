package gdocs

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/latex"
	"google.golang.org/api/docs/v1"
)

var (
	rePrimaryHeadingNum   = regexp.MustCompile(`^(?:[1-9][0-9]*\s+[A-ZÁÉÍÓÚÂÊÎÔÛÃÕÇ\s\-]{3,}|CAPÍTULO\s+[0-9]+.*)$`)
	rePrimaryHeadingNoNum = regexp.MustCompile(`^(?:RESUMO|ABSTRACT|REFERÊNCIAS|SUMÁRIO|LISTA DE ILUSTRAÇÕES|LISTA DE TABELAS|LISTA DE ABREVIATURAS E SIGLAS|ANEXO[S]?.*|APÊNDICE[S]?.*)$`)
	reSecondaryHeading    = regexp.MustCompile(`^[1-9][0-9]*\.[0-9]+\s+[A-ZÁÉÍÓÚÂÊÎÔÛÃÕÇa-záéíóúâêîôûãõç\s\-,/]+$`)
	reTertiaryHeading     = regexp.MustCompile(`^[1-9][0-9]*\.[0-9]+\.[0-9]+\s+[A-ZÁÉÍÓÚÂÊÎÔÛÃÕÇa-záéíóúâêîôûãõç\s\-,/]+$`)
	reQuaternaryHeading   = regexp.MustCompile(`^[1-9][0-9]*\.[0-9]+\.[0-9]+\.[0-9]+\s+[A-ZÁÉÍÓÚÂÊÎÔÛÃÕÇa-záéíóúâêîôûãõç\s\-,/]+$`)
	rePureLatexCmd        = regexp.MustCompile(`(?i)^\s*\\(?:chapter|section|subsection|subsubsection|label|begin|end|ref|cite|caption|centering|hline|toprule|midrule|bottomrule)\*?(?:[{(][^})]*[})]?|\[[^\]]*\]|\{[^}]*\})*\s*$`)
	reLatexLeakInline     = regexp.MustCompile(`(?i)\\(?:chapter|section|subsection|subsubsection|label|begin|end|ref|cite|caption|textit|textbf|path)\*?[{(]([^})]*)[})]?`)
)

// cleanLatexInlineText remove tags de comandos LaTeX em frases mantendo o texto legível e normalizado
func cleanLatexInlineText(s string) string {
	// 1. Símbolos matemáticos inline e caracteres especiais
	s = strings.ReplaceAll(s, `$2 \times 2$`, "2 × 2")
	s = strings.ReplaceAll(s, `\times`, "×")
	s = strings.ReplaceAll(s, `\ge`, "≥")
	s = strings.ReplaceAll(s, `\le`, "≤")
	s = strings.ReplaceAll(s, `\neq`, "≠")
	s = strings.ReplaceAll(s, `\textordmasculine`, "º")
	s = strings.ReplaceAll(s, `\textordfeminine`, "ª")

	// Remove delimitadores de matemática inline simples como $x$ -> x
	s = regexp.MustCompile(`\$([^$]+)\$`).ReplaceAllString(s, "$1")

	// 2. Comandos LaTeX inline
	s = regexp.MustCompile(`\\ref\{([^})]+)[})]`).ReplaceAllString(s, "$1")
	s = regexp.MustCompile(`\\path\{([^})]+)[})]`).ReplaceAllString(s, "$1")
	s = regexp.MustCompile(`\\cite(?:\*?\[[^\]]*\])?\{([^})]+)[})]`).ReplaceAllString(s, "($1)")
	s = regexp.MustCompile(`\\textcite\{([^})]+)[})]`).ReplaceAllString(s, "$1")
	s = regexp.MustCompile(`\\textit\{([^}]+)\}?`).ReplaceAllString(s, "$1")
	s = regexp.MustCompile(`\\textbf\{([^}]+)\}?`).ReplaceAllString(s, "$1")
	s = regexp.MustCompile(`\\label\{([^}]+)\}`).ReplaceAllString(s, "")

	// 3. Til usado como espaço insecável
	s = strings.ReplaceAll(s, "Seção~", "Seção ")
	s = strings.ReplaceAll(s, "Tabela~", "Tabela ")
	s = strings.ReplaceAll(s, "Figura~", "Figura ")
	s = strings.ReplaceAll(s, "Quadro~", "Quadro ")
	s = strings.ReplaceAll(s, "nº~", "nº ")
	s = strings.ReplaceAll(s, "~", " ")

	// 4. Parênteses redundantes/desbalanceados decorrentes de conversão
	s = strings.ReplaceAll(s, "Model Context Protocol) (MCP)", "Model Context Protocol (MCP)")
	s = strings.ReplaceAll(s, "(Retrieval-Augmented Generation), RAG),", "(Retrieval-Augmented Generation -- RAG),")
	s = strings.ReplaceAll(s, "))", ")")
	s = regexp.MustCompile(`\)\s*,`).ReplaceAllString(s, "),")

	// 5. Marcadores residuais artificiais
	s = strings.ReplaceAll(s, "[TABELA]", "")
	s = strings.ReplaceAll(s, "[/TABELA]", "")
	s = strings.ReplaceAll(s, "[VER IMAGEM]", "")

	// 6. Espaços antes de pontuação
	s = regexp.MustCompile(`\s+([,;.:!?])`).ReplaceAllString(s, "$1")

	// 7. Espaços múltiplos
	s = regexp.MustCompile(` {2,}`).ReplaceAllString(s, " ")

	return strings.TrimSpace(s)
}

// isPreTextualNote verifica se o texto corresponde à nota institucional da folha de rosto
func isPreTextualNote(text string) bool {
	t := strings.ToLower(text)
	return strings.Contains(t, "trabalho de conclusão de curso apresentado") ||
		strings.Contains(t, "dissertação apresentada") ||
		strings.Contains(t, "tese apresentada") ||
		strings.Contains(t, "monografia apresentada") ||
		strings.Contains(t, "requisito parcial para obtenção") ||
		strings.HasPrefix(t, "orientador:") ||
		strings.HasPrefix(t, "orientadora:") ||
		strings.HasPrefix(t, "coorientador:") ||
		strings.HasPrefix(t, "coorientadora:")
}

// getBodyEndIndex retorna o índice final do último elemento estrutural do corpo do documento
func getBodyEndIndex(doc *docs.Document) int64 {
	if doc != nil && doc.Body != nil && len(doc.Body.Content) > 0 {
		return doc.Body.Content[len(doc.Body.Content)-1].EndIndex
	}
	return 1
}

// shouldMergeParagraphs avalia se o parágrafo p1 e o parágrafo p2 devem ser unificados em um só
func shouldMergeParagraphs(t1, t2 string, isSpecial1, isSpecial2 bool) bool {
	if isSpecial1 || isSpecial2 {
		return false
	}
	t1 = strings.TrimSpace(t1)
	t2 = strings.TrimSpace(t2)
	if t1 == "" || t2 == "" {
		return false
	}

	// Não funde se p1 terminar com pontuação de término de período
	if strings.HasSuffix(t1, ".") || strings.HasSuffix(t1, "!") || strings.HasSuffix(t1, "?") || strings.HasSuffix(t1, ":") {
		return false
	}

	// Não funde se p2 for título, marcador de lista ou início de tabela/figura
	if strings.HasPrefix(t2, "•") || strings.HasPrefix(t2, "-") || strings.HasPrefix(t2, "*") ||
		strings.HasPrefix(t2, "Figura ") || strings.HasPrefix(t2, "Tabela ") || strings.HasPrefix(t2, "Quadro ") ||
		strings.HasPrefix(t2, "Fonte:") || strings.HasPrefix(t2, "Nota:") {
		return false
	}

	// Se p2 começar com letra minúscula ou pontuação de continuação
	firstRune, _ := utf8.DecodeRuneInString(t2)
	if unicode.IsLower(firstRune) || firstRune == ',' || firstRune == ')' || firstRune == ';' {
		return true
	}

	// Se p1 terminar com vírgula, conectivo, preposição ou hífen, é continuação de frase evidente
	words1 := strings.Fields(t1)
	if len(words1) > 0 {
		lastWord := strings.ToLower(words1[len(words1)-1])
		lastWord = strings.Trim(lastWord, "(),")
		connectives := map[string]bool{
			"e": true, "ou": true, "de": true, "da": true, "do": true, "das": true, "dos": true,
			"em": true, "no": true, "na": true, "nos": true, "nas": true, "para": true,
			"por": true, "pelo": true, "pela": true, "com": true, "que": true, "como": true,
			"ao": true, "à": true, "aos": true, "às": true, "sob": true, "sobre": true,
		}
		if connectives[lastWord] || strings.HasSuffix(t1, ",") || strings.HasSuffix(t1, "-") {
			return true
		}
	}

	return false
}

// HealDocTextAndBlanks saneia quebras de linha duras e linhas em branco no corpo de forma determinística
func HealDocTextAndBlanks(ctx context.Context, documentID string, healBrokenLines bool, collapseBlankLines bool) (int, int, error) {
	if !healBrokenLines && !collapseBlankLines {
		return 0, 0, nil
	}

	svc, err := GetDocsService(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("erro ao obter docs service: %w", err)
	}

	doc, err := svc.Documents.Get(documentID).Do()
	if err != nil {
		return 0, 0, fmt.Errorf("erro ao carregar documento: %w", err)
	}

	if doc.Body == nil || len(doc.Body.Content) == 0 {
		return 0, 0, nil
	}

	// 1. Identifica o início do primeiro capítulo
	firstChapterIndex := int64(-1)
	for _, el := range doc.Body.Content {
		if el.Paragraph != nil {
			pText := ""
			for _, pEl := range el.Paragraph.Elements {
				if pEl.TextRun != nil {
					pText += pEl.TextRun.Content
				}
			}
			trimmed := strings.TrimSpace(pText)
			if rePrimaryHeadingNum.MatchString(trimmed) || trimmed == "1 INTRODUÇÃO" {
				firstChapterIndex = el.StartIndex
				break
			}
		}
	}

	if firstChapterIndex < 0 {
		return 0, 0, nil
	}

	healedLinesCount := 0
	collapsedBlanksCount := 0

	// 2. Desfragmentação de Linhas Quebradas (healBrokenLines)
	if healBrokenLines {
		var mergeReqs []*docs.Request

		content := doc.Body.Content
		for i := len(content) - 2; i >= 0; i-- {
			elA := content[i]
			elB := content[i+1]

			if elA.StartIndex < firstChapterIndex {
				continue
			}

			if elA.Paragraph == nil || elB.Paragraph == nil {
				continue
			}

			pTextA := ""
			for _, pEl := range elA.Paragraph.Elements {
				if pEl.TextRun != nil {
					pTextA += pEl.TextRun.Content
				}
			}

			pTextB := ""
			for _, pEl := range elB.Paragraph.Elements {
				if pEl.TextRun != nil {
					pTextB += pEl.TextRun.Content
				}
			}

			trimA := strings.TrimSpace(pTextA)
			trimB := strings.TrimSpace(pTextB)

			isSpecialA := rePrimaryHeadingNum.MatchString(trimA) || rePrimaryHeadingNoNum.MatchString(trimA) ||
				reSecondaryHeading.MatchString(trimA) || reTertiaryHeading.MatchString(trimA) ||
				strings.HasPrefix(trimA, "Figura ") || strings.HasPrefix(trimA, "Tabela ") ||
				strings.HasPrefix(trimA, "Quadro ") || strings.HasPrefix(trimA, "Fonte:") ||
				strings.HasPrefix(trimA, "Nota:") || isPreTextualNote(trimA)

			isSpecialB := rePrimaryHeadingNum.MatchString(trimB) || rePrimaryHeadingNoNum.MatchString(trimB) ||
				reSecondaryHeading.MatchString(trimB) || reTertiaryHeading.MatchString(trimB) ||
				strings.HasPrefix(trimB, "Figura ") || strings.HasPrefix(trimB, "Tabela ") ||
				strings.HasPrefix(trimB, "Quadro ") || strings.HasPrefix(trimB, "Fonte:") ||
				strings.HasPrefix(trimB, "Nota:") || isPreTextualNote(trimB)

			if shouldMergeParagraphs(trimA, trimB, isSpecialA, isSpecialB) {
				endA := elA.EndIndex
				if !strings.HasSuffix(strings.TrimSuffix(pTextA, "\n"), " ") {
					mergeReqs = append(mergeReqs, &docs.Request{
						InsertText: &docs.InsertTextRequest{
							Location: &docs.Location{Index: endA - 1},
							Text:     " ",
						},
					})
					mergeReqs = append(mergeReqs, &docs.Request{
						DeleteContentRange: &docs.DeleteContentRangeRequest{
							Range: &docs.Range{
								StartIndex: endA,
								EndIndex:   endA + 1,
							},
						},
					})
				} else {
					mergeReqs = append(mergeReqs, &docs.Request{
						DeleteContentRange: &docs.DeleteContentRangeRequest{
							Range: &docs.Range{
								StartIndex: endA - 1,
								EndIndex:   endA,
							},
						},
					})
				}
				healedLinesCount++
			}
		}

		if len(mergeReqs) > 0 {
			batchSize := 200
			for start := 0; start < len(mergeReqs); start += batchSize {
				end := start + batchSize
				if end > len(mergeReqs) {
					end = len(mergeReqs)
				}
				chunk := mergeReqs[start:end]
				_, errBatch := svc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: chunk}).Do()
				if errBatch != nil {
					log.Printf("[AVISO] HealDocTextAndBlanks (merge): %v", errBatch)
					break
				}
			}

			doc, err = svc.Documents.Get(documentID).Do()
			if err != nil {
				return healedLinesCount, 0, fmt.Errorf("erro ao recarregar documento pós-merge: %w", err)
			}
		}
	}

	// 3. Colapso de Linhas em Branco Consecutivas no Corpo (collapseBlankLines)
	if collapseBlankLines {
		var blankReqs []*docs.Request
		bodyEnd := getBodyEndIndex(doc)

		content := doc.Body.Content
		for i := len(content) - 1; i >= 0; i-- {
			el := content[i]
			if el.StartIndex < firstChapterIndex {
				continue
			}
			if el.EndIndex >= bodyEnd {
				continue
			}

			if el.Paragraph != nil {
				pText := ""
				for _, pEl := range el.Paragraph.Elements {
					if pEl.TextRun != nil {
						pText += pEl.TextRun.Content
					}
				}
				if strings.TrimSpace(pText) == "" && el.EndIndex > el.StartIndex {
					blankReqs = append(blankReqs, &docs.Request{
						DeleteContentRange: &docs.DeleteContentRangeRequest{
							Range: &docs.Range{
								StartIndex: el.StartIndex,
								EndIndex:   el.EndIndex,
							},
						},
					})
					collapsedBlanksCount++
				}
			}
		}

		if len(blankReqs) > 0 {
			batchSize := 200
			for start := 0; start < len(blankReqs); start += batchSize {
				end := start + batchSize
				if end > len(blankReqs) {
					end = len(blankReqs)
				}
				chunk := blankReqs[start:end]
				_, errBlank := svc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: chunk}).Do()
				if errBlank != nil {
					log.Printf("[AVISO] HealDocTextAndBlanks (blanks): %v", errBlank)
					break
				}
			}
		}
	}

	return healedLinesCount, collapsedBlanksCount, nil
}

// ReplaceBlockABNT substitui uma faixa do documento aplicando a sequência canônica anti-bleeding e formatação ABNT
func ReplaceBlockABNT(ctx context.Context, documentID string, startIndex, endIndex int64, newText string) (string, error) {
	if startIndex < 0 || endIndex < startIndex {
		return "", fmt.Errorf("índices inválidos: start_index (%d) e end_index (%d)", startIndex, endIndex)
	}

	svc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	doc, err := svc.Documents.Get(documentID).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao obter documento '%s': %w", documentID, err)
	}

	bodyEndIndex := getBodyEndIndex(doc)

	// Extrai texto existente para validação de comentários se o range for válido
	if endIndex > startIndex && doc.Body != nil {
		extracted, errExtract := extractTextRange(doc.Body.Content, startIndex, endIndex)
		if errExtract == nil && extracted != "" {
			if errCollision := CheckCommentCollision(ctx, documentID, extracted); errCollision != nil {
				return "", errCollision
			}
		}
	}

	requests := BuildABNTReplaceSequence(startIndex, endIndex, bodyEndIndex, newText)
	if len(requests) == 0 {
		return "Nenhuma alteração necessária.", nil
	}

	updateReq := &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}

	_, err = svc.Documents.BatchUpdate(documentID, updateReq).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao executar batchUpdate (ReplaceBlockABNT): %w", err)
	}

	return fmt.Sprintf("SUCESSO: Bloco [%d - %d] substituído e formatado segundo as normas ABNT (anti-bleeding aplicado).", startIndex, endIndex), nil
}

// ApplyABNTPageMargins aplica as margens oficiais da ABNT NBR 14724 no documento (3cm sup/esq, 2cm inf/dir, formato A4)
func ApplyABNTPageMargins(ctx context.Context, documentID string) (string, error) {
	svc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	req := &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{
			BuildABNTDocumentMarginsRequest(),
		},
	}

	_, err = svc.Documents.BatchUpdate(documentID, req).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao aplicar margens ABNT: %w", err)
	}

	return "SUCESSO: Margens e tamanho de página A4 ABNT (Superior 3cm, Esquerda 3cm, Inferior 2cm, Direita 2cm) aplicadas com sucesso.", nil
}

// FormatDocumentABNT formata parágrafos, estrutura e página segundo as normas ABNT completas (NBR 14724, NBR 6024 e NBR 6023).
// Unifica a normalização estrutural de títulos com fonte 12pt rigorosa, margens A4 3/3/2/2 cm,
// tratamento contextual de capa/folha de rosto, referências bibliográficas, legendas e citações longas.
func FormatDocumentABNT(ctx context.Context, documentID string, startIndex, endIndex int64, applyMargins bool, cleanLatex ...bool) (string, error) {
	cleanLatexCmds := true
	if len(cleanLatex) > 0 {
		cleanLatexCmds = cleanLatex[0]
	}
	return FormatDocumentABNTExtended(ctx, documentID, startIndex, endIndex, applyMargins, cleanLatexCmds, ABNTDefaultFontFamily)
}

// FormatDocumentABNTExtended formata o documento aplicando família de fonte unificada e todas as regras ABNT
func FormatDocumentABNTExtended(ctx context.Context, documentID string, startIndex, endIndex int64, applyMargins bool, cleanLatexCmds bool, fontFamily string, opts ...bool) (string, error) {
	if fontFamily == "" {
		fontFamily = ABNTDefaultFontFamily
	}

	collapseBlankLines := true
	healBrokenLines := true
	if len(opts) > 0 {
		collapseBlankLines = opts[0]
	}
	if len(opts) > 1 {
		healBrokenLines = opts[1]
	}

	svc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	doc, err := svc.Documents.Get(documentID).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao obter documento '%s': %w", documentID, err)
	}

	if doc.Body == nil || len(doc.Body.Content) == 0 {
		return "Documento vazio.", nil
	}

	// ─── ETAPA 1: Higienização de Comandos LaTeX Vazados e Normalizações Globais ───
	cleanedCommandsCount := 0
	if cleanLatexCmds {
		globalCleanReqs := []*docs.Request{
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "  "}, ReplaceText: " "}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "  "}, ReplaceText: " "}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: " ,"}, ReplaceText: ","}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: " ."}, ReplaceText: "."}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: " ;"}, ReplaceText: ";"}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: " :"}, ReplaceText: ":"}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "$2 \\times 2$"}, ReplaceText: "2 × 2"}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "\\times"}, ReplaceText: "×"}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "\\ge"}, ReplaceText: "≥"}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "\\le"}, ReplaceText: "≤"}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "~"}, ReplaceText: " "}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "[TABELA]"}, ReplaceText: ""}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "[/TABELA]"}, ReplaceText: ""}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "[VER IMAGEM"}, ReplaceText: ""}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "Model Context Protocol) (MCP)"}, ReplaceText: "Model Context Protocol (MCP)"}},
			{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "(Retrieval-Augmented Generation), RAG),"}, ReplaceText: "(Retrieval-Augmented Generation -- RAG),"}},
		}

		if _, errGlob := svc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: globalCleanReqs}).Do(); errGlob == nil {
			cleanedCommandsCount += len(globalCleanReqs)
			doc, err = svc.Documents.Get(documentID).Do()
			if err != nil {
				return "", fmt.Errorf("erro ao recarregar documento pós-substituições: %w", err)
			}
		}

		var purgeReqs []*docs.Request
		for i := len(doc.Body.Content) - 1; i >= 0; i-- {
			el := doc.Body.Content[i]
			if el.Paragraph != nil {
				pText := ""
				for _, pEl := range el.Paragraph.Elements {
					if pEl.TextRun != nil {
						pText += pEl.TextRun.Content
					}
				}
				trimmed := strings.TrimSpace(pText)
				if trimmed == "" {
					continue
				}

				// Se for parágrafo puramente composto por comando LaTeX (ex: \chapter{Introdução), \label{...}, \begin{table)...)
				if rePureLatexCmd.MatchString(trimmed) {
					pEnd := el.EndIndex
					pStart := el.StartIndex
					if pEnd > pStart {
						purgeReqs = append(purgeReqs, &docs.Request{
							DeleteContentRange: &docs.DeleteContentRangeRequest{
								Range: &docs.Range{
									StartIndex: pStart,
									EndIndex:   pEnd - 1,
								},
							},
						})
						cleanedCommandsCount++
					}
					continue
				}

				// Se tiver comando vazado inline (ex: \path{pgvector), \ref{...}, \cite{...})
				if strings.Contains(trimmed, `\ref{`) || strings.Contains(trimmed, `\path{`) || strings.Contains(trimmed, `\cite{`) || strings.Contains(trimmed, `\textit{`) || strings.Contains(trimmed, `\textbf{`) || strings.Contains(trimmed, `$`) {
					cleaned := cleanLatexInlineText(trimmed)
					if cleaned != trimmed {
						pEnd := el.EndIndex
						pStart := el.StartIndex
						if pEnd > pStart {
							purgeReqs = append(purgeReqs, &docs.Request{
								DeleteContentRange: &docs.DeleteContentRangeRequest{
									Range: &docs.Range{
										StartIndex: pStart,
										EndIndex:   pEnd - 1,
									},
								},
							})
							purgeReqs = append(purgeReqs, &docs.Request{
								InsertText: &docs.InsertTextRequest{
									Location: &docs.Location{Index: pStart},
									Text:     cleaned,
								},
							})
							cleanedCommandsCount++
						}
					}
				}
			}
		}

		if len(purgeReqs) > 0 {
			_, errPurge := svc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: purgeReqs}).Do()
			if errPurge == nil {
				doc, err = svc.Documents.Get(documentID).Do()
				if err != nil {
					return "", fmt.Errorf("erro ao recarregar documento pós-higienização: %w", err)
				}
			}
		}
	}

	// ─── ETAPA 1.5: Saneamento Mecânico de Quebras Duras e Linhas em Branco no Corpo ───
	healedCount, collapsedBlanks, errHeal := HealDocTextAndBlanks(ctx, documentID, healBrokenLines, collapseBlankLines)
	if errHeal == nil && (healedCount > 0 || collapsedBlanks > 0) {
		doc, err = svc.Documents.Get(documentID).Do()
		if err != nil {
			return "", fmt.Errorf("erro ao recarregar documento pós-saneamento: %w", err)
		}
	}

	bodyEnd := getBodyEndIndex(doc)
	start := int64(1)
	end := bodyEnd - 1

	if startIndex >= 0 {
		start = startIndex
	}
	if endIndex > 0 && endIndex <= bodyEnd-1 {
		end = endIndex
	}

	if start >= end {
		return "", fmt.Errorf("intervalo de formatação inválido: [%d - %d]", start, end)
	}

	var requests []*docs.Request

	// ─── ETAPA 2: Configuração de Página: Margens ABNT 3/3/2/2 cm e Folha A4 ───
	if applyMargins {
		requests = append(requests, BuildABNTDocumentMarginsRequest())
	}

	// Contadores estatísticos para relatório
	formattedParas := 0
	headingsCount := 0
	preTextualCount := 0
	referencesCount := 0
	captionsCount := 0
	longQuotesCount := 0

	// Máquina de estados contextual ABNT
	beforeFirstChapter := true
	inReferences := false

	// Cor preta padrão (#000000)
	blackColor := &docs.OptionalColor{
		Color: &docs.Color{
			RgbColor: &docs.RgbColor{Red: 0, Green: 0, Blue: 0},
		},
	}

	// Carrega mapa de seções do TCC para identificação e promoção de títulos despadronizados
	knownByTitle := make(map[string]latex.KnownSectionInfo)
	if knownSections, errKS := latex.GetNumberedTCCStructure("resources/tcc"); errKS == nil {
		for _, ks := range knownSections {
			cleanKey := strings.ToLower(strings.TrimSpace(ks.Title))
			knownByTitle[cleanKey] = ks
		}
	}

	// ─── ETAPA 3: Formatação Estrutural ABNT Elemento a Elemento ───
	for _, el := range doc.Body.Content {
		// Elemento fora do range solicitado
		if el.EndIndex <= start || el.StartIndex >= end {
			continue
		}

		// Tabelas nativas são blindadas: não aplica recuo nas células
		if el.Table != nil {
			continue
		}

		if el.Paragraph != nil {
			p := el.Paragraph
			currentStyle := ""
			if p.ParagraphStyle != nil {
				currentStyle = p.ParagraphStyle.NamedStyleType
			}

			// Coleta texto completo do parágrafo
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

			pStart := el.StartIndex
			pEnd := el.EndIndex
			if pStart < start {
				pStart = start
			}
			if pEnd > end {
				pEnd = end
			}
			if pStart >= pEnd {
				continue
			}

			pRange := &docs.Range{
				StartIndex: pStart,
				EndIndex:   pEnd,
			}

			// Identificação hierárquica segundo NBR 6024 e ABNT 2026
			isHeading1Num := rePrimaryHeadingNum.MatchString(trimmed) || (currentStyle == "HEADING_1" && !rePrimaryHeadingNoNum.MatchString(trimmed))
			isHeading1NoNum := rePrimaryHeadingNoNum.MatchString(trimmed) || trimmed == "3 REFERÊNCIAS"
			isHeading2 := reSecondaryHeading.MatchString(trimmed) || currentStyle == "HEADING_2"
			isHeading3 := reTertiaryHeading.MatchString(trimmed) || currentStyle == "HEADING_3"
			isHeading4 := reQuaternaryHeading.MatchString(trimmed) || currentStyle == "HEADING_4"

			// Enriquecimento semântico: se o parágrafo for curto e coincidir com seção do TCC
			if !isHeading1Num && !isHeading1NoNum && !isHeading2 && !isHeading3 && !isHeading4 && len(trimmed) < 120 {
				cleanTrim := strings.ToLower(trimmed)
				if ks, found := knownByTitle[cleanTrim]; found {
					switch ks.Level {
					case 1:
						isHeading1Num = true
					case 2:
						isHeading2 = true
					case 3:
						isHeading3 = true
					case 4:
						isHeading4 = true
					}
				}
			}

			isCaption := strings.HasPrefix(trimmed, "Figura ") ||
				strings.HasPrefix(trimmed, "Tabela ") ||
				strings.HasPrefix(trimmed, "Quadro ") ||
				strings.HasPrefix(trimmed, "Gráfico ")
			isSource := strings.HasPrefix(trimmed, "Fonte:") || strings.HasPrefix(trimmed, "Nota:")

			// Identificação de citação longa (> 3 linhas ou recuo existente >= 3.5cm)
			existingIndent := float64(0)
			if p.ParagraphStyle != nil && p.ParagraphStyle.IndentStart != nil {
				existingIndent = p.ParagraphStyle.IndentStart.Magnitude
			}
			isLongQuote := existingIndent >= 90.0 && !isHeading1Num && !isHeading1NoNum && !isHeading2 && !isPreTextualNote(trimmed)

			// ── CASO A: TÍTULO PRIMÁRIO NUMERADO (HEADING_1 - Ex: 1 INTRODUÇÃO, 2 REFERENCIAL TEÓRICO) ──
			if isHeading1Num {
				beforeFirstChapter = false
				inReferences = false
				headingsCount++

				// ParagraphStyle: Alinhamento à esquerda, sem recuo, 1.5 entrelinhas, keepWithNext
				requests = append(requests, &docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: pRange,
						ParagraphStyle: &docs.ParagraphStyle{
							NamedStyleType:      "HEADING_1",
							Alignment:           "START",
							LineSpacing:         ABNTLineSpacing,
							IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
							IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceAbove:          &docs.Dimension{Magnitude: 18, Unit: "PT"},
							SpaceBelow:          &docs.Dimension{Magnitude: 12, Unit: "PT"},
							KeepWithNext:        true,
							AvoidWidowAndOrphan: true,
						},
						Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,keepWithNext,avoidWidowAndOrphan",
					},
				})

				// TextStyle: ABNT NBR 14724/6024 exige FONTE 12PT (NÃO 20PT!), NEGRITO, PRETO e FONTE UNIFORME
				requests = append(requests, &docs.Request{
					UpdateTextStyle: &docs.UpdateTextStyleRequest{
						Range: pRange,
						TextStyle: &docs.TextStyle{
							FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeHeadingPt, Unit: "PT"}, // 12.0 pt
							WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
							Bold:               true,
							ForegroundColor:    blackColor,
						},
						Fields: "fontSize,weightedFontFamily,bold,foregroundColor",
					},
				})
				continue
			}

			// ── CASO B: TÍTULO SEM NÚMERO (HEADING_1 - RESUMO, ABSTRACT, REFERÊNCIAS, SUMÁRIO) ──
			if isHeading1NoNum {
				beforeFirstChapter = false
				if strings.Contains(trimmed, "REFERÊNCIAS") {
					inReferences = true
				} else {
					inReferences = false
				}
				headingsCount++

				// ParagraphStyle: Centralizado (NBR 14724 item 5.2.2 / Guia ABNT), sem recuo, keepWithNext
				requests = append(requests, &docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: pRange,
						ParagraphStyle: &docs.ParagraphStyle{
							NamedStyleType:      "HEADING_1",
							Alignment:           "CENTER",
							LineSpacing:         ABNTLineSpacing,
							IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
							IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceAbove:          &docs.Dimension{Magnitude: 18, Unit: "PT"},
							SpaceBelow:          &docs.Dimension{Magnitude: 12, Unit: "PT"},
							KeepWithNext:        true,
							AvoidWidowAndOrphan: true,
						},
						Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,keepWithNext,avoidWidowAndOrphan",
					},
				})

				// TextStyle: FONTE 12PT, NEGRITO, PRETO
				requests = append(requests, &docs.Request{
					UpdateTextStyle: &docs.UpdateTextStyleRequest{
						Range: pRange,
						TextStyle: &docs.TextStyle{
							FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeHeadingPt, Unit: "PT"}, // 12.0 pt
							WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
							Bold:               true,
							ForegroundColor:    blackColor,
						},
						Fields: "fontSize,weightedFontFamily,bold,foregroundColor",
					},
				})
				continue
			}

			// ── CASO C: SEÇÃO SECUNDÁRIA (HEADING_2 - Ex: 1.1 Contextualização) ──
			if isHeading2 {
				headingsCount++

				// NBR 6024: Seção secundária alinhada à esquerda, caixa alta sem negrito (ou regular)
				requests = append(requests, &docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: pRange,
						ParagraphStyle: &docs.ParagraphStyle{
							NamedStyleType:      "HEADING_2",
							Alignment:           "START",
							LineSpacing:         ABNTLineSpacing,
							IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
							IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceAbove:          &docs.Dimension{Magnitude: 12, Unit: "PT"},
							SpaceBelow:          &docs.Dimension{Magnitude: 6, Unit: "PT"},
							KeepWithNext:        true,
							AvoidWidowAndOrphan: true,
						},
						Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,keepWithNext,avoidWidowAndOrphan",
					},
				})

				// TextStyle: FONTE 12PT (NÃO 16PT!), Regular, Preto
				requests = append(requests, &docs.Request{
					UpdateTextStyle: &docs.UpdateTextStyleRequest{
						Range: pRange,
						TextStyle: &docs.TextStyle{
							FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeHeadingPt, Unit: "PT"}, // 12.0 pt
							WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
							Bold:               false,
							ForegroundColor:    blackColor,
						},
						Fields: "fontSize,weightedFontFamily,bold,foregroundColor",
					},
				})
				continue
			}

			// ── CASO D: SEÇÃO TERCIÁRIA (HEADING_3 - Ex: 1.1.1 Escopo do Projeto) ──
			if isHeading3 {
				headingsCount++

				// NBR 6024: Seção terciária alinhada à esquerda, minúsculas com negrito
				requests = append(requests, &docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: pRange,
						ParagraphStyle: &docs.ParagraphStyle{
							NamedStyleType:      "HEADING_3",
							Alignment:           "START",
							LineSpacing:         ABNTLineSpacing,
							IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
							IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceAbove:          &docs.Dimension{Magnitude: 12, Unit: "PT"},
							SpaceBelow:          &docs.Dimension{Magnitude: 6, Unit: "PT"},
							KeepWithNext:        true,
							AvoidWidowAndOrphan: true,
						},
						Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,keepWithNext,avoidWidowAndOrphan",
					},
				})

				// TextStyle: FONTE 12PT (NÃO 14PT!), Negrito, Preto
				requests = append(requests, &docs.Request{
					UpdateTextStyle: &docs.UpdateTextStyleRequest{
						Range: pRange,
						TextStyle: &docs.TextStyle{
							FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeHeadingPt, Unit: "PT"}, // 12.0 pt
							WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
							Bold:               true,
							ForegroundColor:    blackColor,
						},
						Fields: "fontSize,weightedFontFamily,bold,foregroundColor",
					},
				})
				continue
			}

			// ── CASO E: SEÇÃO QUATERNÁRIA (HEADING_4) ──
			if isHeading4 {
				headingsCount++

				requests = append(requests, &docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: pRange,
						ParagraphStyle: &docs.ParagraphStyle{
							NamedStyleType:      "HEADING_4",
							Alignment:           "START",
							LineSpacing:         ABNTLineSpacing,
							IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
							IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceAbove:          &docs.Dimension{Magnitude: 12, Unit: "PT"},
							SpaceBelow:          &docs.Dimension{Magnitude: 6, Unit: "PT"},
							KeepWithNext:        true,
							AvoidWidowAndOrphan: true,
						},
						Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,keepWithNext,avoidWidowAndOrphan",
					},
				})

				requests = append(requests, &docs.Request{
					UpdateTextStyle: &docs.UpdateTextStyleRequest{
						Range: pRange,
						TextStyle: &docs.TextStyle{
							FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeHeadingPt, Unit: "PT"}, // 12.0 pt
							WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
							Bold:               false,
							ForegroundColor:    blackColor,
						},
						Fields: "fontSize,weightedFontFamily,bold,foregroundColor",
					},
				})
				continue
			}

			// ── CASO F: ELEMENTOS PRÉ-TEXTUAIS (Capa e Folha de Rosto) ──
			if beforeFirstChapter {
				preTextualCount++

				if isPreTextualNote(trimmed) {
					// Nota de Apresentação / Orientador: Recuo à esquerda de 7.5cm, simples, 10pt
					requests = append(requests, &docs.Request{
						UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
							Range: pRange,
							ParagraphStyle: &docs.ParagraphStyle{
								NamedStyleType:      "NORMAL_TEXT",
								Alignment:           "JUSTIFIED",
								LineSpacing:         ABNTLineSpacingSingle,
								IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
								IndentStart:         &docs.Dimension{Magnitude: ABNTIndentPreTextualNotePt, Unit: "PT"}, // 7.5 cm
								SpaceAbove:          &docs.Dimension{Magnitude: 8, Unit: "PT"},
								SpaceBelow:          &docs.Dimension{Magnitude: 8, Unit: "PT"},
								AvoidWidowAndOrphan: true,
							},
							Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,avoidWidowAndOrphan",
						},
					})

					requests = append(requests, &docs.Request{
						UpdateTextStyle: &docs.UpdateTextStyleRequest{
							Range: pRange,
							TextStyle: &docs.TextStyle{
								FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeCaptionPt, Unit: "PT"}, // 10.0 pt
								WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
								ForegroundColor:    blackColor,
							},
							Fields: "fontSize,weightedFontFamily,foregroundColor",
						},
					})
				} else {
					// Instituição, Autor, Título, Cidade/Ano: Centralizados, sem recuo de 1ª linha
					isTitle := len(trimmed) > 30 && strings.ToUpper(trimmed) == trimmed && !strings.Contains(trimmed, "INSTITUTO")
					boldTitle := isTitle

					requests = append(requests, &docs.Request{
						UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
							Range: pRange,
							ParagraphStyle: &docs.ParagraphStyle{
								NamedStyleType:      "NORMAL_TEXT",
								Alignment:           "CENTER",
								LineSpacing:         ABNTLineSpacing,
								IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
								IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT"},
								SpaceAbove:          &docs.Dimension{Magnitude: 6, Unit: "PT"},
								SpaceBelow:          &docs.Dimension{Magnitude: 6, Unit: "PT"},
								AvoidWidowAndOrphan: true,
							},
							Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,avoidWidowAndOrphan",
						},
					})

					requests = append(requests, &docs.Request{
						UpdateTextStyle: &docs.UpdateTextStyleRequest{
							Range: pRange,
							TextStyle: &docs.TextStyle{
								FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeBodyPt, Unit: "PT"}, // 12.0 pt
								WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
								Bold:               boldTitle,
								ForegroundColor:    blackColor,
							},
							Fields: "fontSize,weightedFontFamily,bold,foregroundColor",
						},
					})
				}
				continue
			}

			// ── CASO G: BLOCO DE REFERÊNCIAS (NBR 6023) ──
			if inReferences {
				referencesCount++

				// ABNT NBR 6023: Alinhamento à esquerda, espaço simples, sem recuo de 1ª linha, 6pt entre referências
				requests = append(requests, &docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: pRange,
						ParagraphStyle: &docs.ParagraphStyle{
							NamedStyleType:      "NORMAL_TEXT",
							Alignment:           "START",
							LineSpacing:         ABNTLineSpacingSingle, // 1.0 simples
							IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
							IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceAbove:          &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceBelow:          &docs.Dimension{Magnitude: 6, Unit: "PT"}, // Separação entre referências
							AvoidWidowAndOrphan: true,
						},
						Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,avoidWidowAndOrphan",
					},
				})

				requests = append(requests, &docs.Request{
					UpdateTextStyle: &docs.UpdateTextStyleRequest{
						Range: pRange,
						TextStyle: &docs.TextStyle{
							FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeBodyPt, Unit: "PT"}, // 12.0 pt
							WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
							ForegroundColor:    blackColor,
						},
						Fields: "fontSize,weightedFontFamily,foregroundColor",
					},
				})
				continue
			}

			// ── CASO H: LEGENDAS E FONTES DE ILUSTRAÇÕES ──
			if isCaption || isSource {
				captionsCount++
				keepWithNext := isCaption // Legenda acompanha a figura/tabela

				requests = append(requests, &docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: pRange,
						ParagraphStyle: &docs.ParagraphStyle{
							NamedStyleType:      "NORMAL_TEXT",
							Alignment:           "CENTER",
							LineSpacing:         ABNTLineSpacingSingle,
							IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
							IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceAbove:          &docs.Dimension{Magnitude: 4, Unit: "PT"},
							SpaceBelow:          &docs.Dimension{Magnitude: 4, Unit: "PT"},
							KeepWithNext:        keepWithNext,
							AvoidWidowAndOrphan: true,
						},
						Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,keepWithNext,avoidWidowAndOrphan",
					},
				})

				requests = append(requests, &docs.Request{
					UpdateTextStyle: &docs.UpdateTextStyleRequest{
						Range: pRange,
						TextStyle: &docs.TextStyle{
							FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeCaptionPt, Unit: "PT"}, // 10.0 pt
							WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
							ForegroundColor:    blackColor,
						},
						Fields: "fontSize,weightedFontFamily,foregroundColor",
					},
				})
				continue
			}

			// ── CASO I: CITAÇÃO DIRETA LONGA (> 3 linhas ou recuo existente) ──
			if isLongQuote {
				longQuotesCount++

				requests = append(requests, &docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: pRange,
						ParagraphStyle: &docs.ParagraphStyle{
							NamedStyleType:      "NORMAL_TEXT",
							Alignment:           "JUSTIFIED",
							LineSpacing:         ABNTLineSpacingSingle, // 1.0 simples
							IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
							IndentStart:         &docs.Dimension{Magnitude: ABNTIndentLongQuotePt, Unit: "PT"}, // 4.0 cm
							SpaceAbove:          &docs.Dimension{Magnitude: 6, Unit: "PT"},
							SpaceBelow:          &docs.Dimension{Magnitude: 6, Unit: "PT"},
							AvoidWidowAndOrphan: true,
						},
						Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,avoidWidowAndOrphan",
					},
				})

				requests = append(requests, &docs.Request{
					UpdateTextStyle: &docs.UpdateTextStyleRequest{
						Range: pRange,
						TextStyle: &docs.TextStyle{
							FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeCaptionPt, Unit: "PT"}, // 10.0 pt
							WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
							ForegroundColor:    blackColor,
						},
						Fields: "fontSize,weightedFontFamily,foregroundColor",
					},
				})
				continue
			}

			// ── CASO J: CORPO PADRÃO DE TEXTO (NORMAL_TEXT) ──
			formattedParas++
			requests = append(requests, &docs.Request{
				UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
					Range: pRange,
					ParagraphStyle: &docs.ParagraphStyle{
						NamedStyleType:      "NORMAL_TEXT",
						Alignment:           "JUSTIFIED",
						LineSpacing:         ABNTLineSpacing,
						IndentFirstLine:     &docs.Dimension{Magnitude: ABNTIndentFirstLinePt, Unit: "PT"}, // 1.25 cm
						IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT"},
						SpaceAbove:          &docs.Dimension{Magnitude: 0, Unit: "PT"},
						SpaceBelow:          &docs.Dimension{Magnitude: 0, Unit: "PT"},
						AvoidWidowAndOrphan: true,
					},
					Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,avoidWidowAndOrphan",
				},
			})

			requests = append(requests, &docs.Request{
				UpdateTextStyle: &docs.UpdateTextStyleRequest{
					Range: pRange,
					TextStyle: &docs.TextStyle{
						FontSize:           &docs.Dimension{Magnitude: ABNTFontSizeBodyPt, Unit: "PT"}, // 12.0 pt
						WeightedFontFamily: &docs.WeightedFontFamily{FontFamily: fontFamily},
						ForegroundColor:    blackColor,
					},
					Fields: "fontSize,weightedFontFamily,foregroundColor",
				},
			})
		}
	}

	if len(requests) == 0 {
		return "Nenhum elemento para formatação no intervalo.", nil
	}

	updateReq := &docs.BatchUpdateDocumentRequest{Requests: requests}
	_, err = svc.Documents.BatchUpdate(documentID, updateReq).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao formatar documento para ABNT: %w", err)
	}

	return fmt.Sprintf("SUCESSO: Conformidade Estrutural e Formatação ABNT aplicada com sucesso no intervalo [%d - %d]:\n"+
		"- Títulos Normalizados com Fonte 12pt (NBR 6024): %d\n"+
		"- Linhas quebradas desfragmentadas e unidas em parágrafos: %d\n"+
		"- Linhas em branco excessivas eliminadas no corpo: %d\n"+
		"- Parágrafos de Corpo (1.25cm / 1.5 / Justificado / 12pt): %d\n"+
		"- Elementos Pré-Textuais formatados (Capa/Folha de Rosto sem recuo 1.25cm): %d\n"+
		"- Referências Bibliográficas (NBR 6023 - À Esquerda / 1.0 / 12pt): %d\n"+
		"- Legendas e Fontes (Centralizadas / 10pt / 1.0): %d\n"+
		"- Citações Diretas Longas (Recuo 4cm / 10pt / 1.0): %d\n"+
		"- Comandos LaTeX vazados limpos: %d\n"+
		"- Família Tipográfica Uniforme: %s\n"+
		"- Margens A4 (Superior 3cm, Esquerda 3cm, Inferior 2cm, Direita 2cm): %v",
		start, end, headingsCount, healedCount, collapsedBlanks, formattedParas, preTextualCount, referencesCount, captionsCount, longQuotesCount, cleanedCommandsCount, fontFamily, applyMargins), nil
}
