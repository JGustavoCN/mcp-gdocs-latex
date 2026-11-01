package gdocs

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/googleapi"

	"github.com/JGustavoCN/mcp-gdocs-latex/internal/latex"
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

// FontFamilyAnomalyItem detalha trechos com família de fonte diferente do padrão ABNT
type FontFamilyAnomalyItem struct {
	StartIndex     int64  `json:"start_index"`
	EndIndex       int64  `json:"end_index"`
	CurrentFamily  string `json:"current_family"`
	ExpectedFamily string `json:"expected_family"`
	Text           string `json:"text"`
}

// TextColorAnomalyItem detalha trechos com cor de texto anômala fora do preto padrão
type TextColorAnomalyItem struct {
	StartIndex int64  `json:"start_index"`
	EndIndex   int64  `json:"end_index"`
	RgbHex     string `json:"rgb_hex"`
	Text       string `json:"text"`
}

// OrderViolationItem detalha anomalias de ordem estrutural segundo ABNT NBR 14724 / NBR 6024
type OrderViolationItem struct {
	ElementIndex int    `json:"element_index"`
	StartIndex   int64  `json:"start_index"`
	EndIndex     int64  `json:"end_index"`
	SectionTitle string `json:"section_title"`
	Severity     string `json:"severity"` // "CRÍTICO" ou "ALERTA"
	RuleViolated string `json:"rule_violated"`
	Details      string `json:"details"`
}

// LatexLeakItem detalha vazamentos de comandos ou fórmulas LaTeX no texto
type LatexLeakItem struct {
	ElementIndex int    `json:"element_index"`
	StartIndex   int64  `json:"start_index"`
	EndIndex     int64  `json:"end_index"`
	LeakType     string `json:"leak_type"`
	Snippet      string `json:"snippet"`
	Suggestion   string `json:"suggestion"`
}

// DoubleSpaceAnomalyItem detalha múltiplos espaços ou espaços antes de pontuação
type DoubleSpaceAnomalyItem struct {
	StartIndex int64  `json:"start_index"`
	EndIndex   int64  `json:"end_index"`
	IssueType  string `json:"issue_type"`
	Snippet    string `json:"snippet"`
}

// ParagraphQualityItem detalha parágrafos curtos demais (órfãos) ou longos demais
type ParagraphQualityItem struct {
	ElementIndex int    `json:"element_index"`
	StartIndex   int64  `json:"start_index"`
	EndIndex     int64  `json:"end_index"`
	IssueType    string `json:"issue_type"`
	Length       int    `json:"length"`
	Snippet      string `json:"snippet"`
	Suggestion   string `json:"suggestion"`
}

// BulletListAnomalyItem detalha listas de marcadores que deveriam ser alíneas ABNT
type BulletListAnomalyItem struct {
	StartIndex       int64  `json:"start_index"`
	EndIndex         int64  `json:"end_index"`
	ConsecutiveCount int    `json:"consecutive_count"`
	Snippet          string `json:"snippet"`
	Suggestion       string `json:"suggestion"`
}

// LintReport contém os resultados da análise estrutural, tipográfica e de revisões do Google Doc
type LintReport struct {
	TotalElements             int                      `json:"total_elements"`
	TotalTables               int                      `json:"total_tables"`
	EmptyTables               []string                 `json:"empty_tables,omitempty"`
	FilledTables              []string                 `json:"filled_tables,omitempty"`
	BoldCascades              []string                 `json:"bold_cascades,omitempty"`
	ReviewHighlights          []ReviewHighlightItem    `json:"review_highlights,omitempty"`
	FontSizeAnomalies         []FontSizeAnomalyItem    `json:"font_size_anomalies,omitempty"`
	FontFamilyAnomalies       []FontFamilyAnomalyItem  `json:"font_family_anomalies,omitempty"`
	TextColorAnomalies        []TextColorAnomalyItem   `json:"text_color_anomalies,omitempty"`
	StructuralOrderViolations []OrderViolationItem     `json:"structural_order_violations,omitempty"`
	MisplacedElements         []string                 `json:"misplaced_elements,omitempty"`
	CoverStatus               string                   `json:"cover_status,omitempty"`
	LatexSyntaxLeaks          []LatexLeakItem          `json:"latex_syntax_leaks,omitempty"`
	DoubleSpacingAnomalies    []DoubleSpaceAnomalyItem `json:"double_spacing_anomalies,omitempty"`
	ParagraphQualityAnomalies []ParagraphQualityItem   `json:"paragraph_quality_anomalies,omitempty"`
	BulletListAnomalies       []BulletListAnomalyItem  `json:"bullet_list_anomalies,omitempty"`
	ConsecutiveBlankLines     []string                 `json:"consecutive_blank_lines,omitempty"`
	SoftLineBreaks            []string                 `json:"soft_line_breaks,omitempty"`
	DanglingMarkers           []string                 `json:"dangling_markers,omitempty"`
	FixedAnomaliesCount       int                      `json:"fixed_anomalies_count,omitempty"`
	Message                   string                   `json:"message"`
}

// LintOptions configura parâmetros opcionais de auditoria e correção
type LintOptions struct {
	AutoFix              bool   `json:"auto_fix"`
	ClearHighlights      bool   `json:"clear_highlights"`
	ExpectedFontFamily   string `json:"expected_font_family"`
	CheckStructuralOrder bool   `json:"check_structural_order"`
	CheckLatexLeaks      bool   `json:"check_latex_leaks"`
	CheckSpacingAndText  bool   `json:"check_spacing_and_text"`
}

var (
	reLatexCmdInText    = regexp.MustCompile(`(?i)\\(?:chapter|section|subsection|subsubsection|label|begin|end|ref|cite|textcite|caption|centering|hline|toprule|midrule|bottomrule|textit|textbf|texttt|allowbreak|url|path|verb|footnote|hfill|newline|times|ge|le|neq|textordmasculine|textordfeminine|item)\*?(?:\{[^}\n]*\}|\{[^}\n]*\)|\[[^\]\n]*\])?`)
	reMathSyntaxInText  = regexp.MustCompile(`\$[^$]+\$`)
	reDoubleSpace       = regexp.MustCompile(` {2,}`)
	reSpaceBeforePunct  = regexp.MustCompile(`\s+([,;.:!?])`)
	reBrokenParenthesis = regexp.MustCompile(`\)\s*\)|\(\s*\)`)
	reChapterNumber     = regexp.MustCompile(`^([1-9][0-9]*)\s+`)
)

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
		if c.QuotedFileContent != nil && strings.TrimSpace(c.QuotedFileContent.Value) != "" {
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

// LintDocumentABNT executa a auditoria completa de conformidade ABNT no Google Docs
func LintDocumentABNT(ctx context.Context, documentID string, autoFix bool, clearHighlights bool, extraOpts ...LintOptions) (*LintReport, error) {
	opts := LintOptions{
		AutoFix:              autoFix,
		ClearHighlights:      clearHighlights,
		ExpectedFontFamily:   ABNTDefaultFontFamily, // "Arial"
		CheckStructuralOrder: true,
		CheckLatexLeaks:      true,
		CheckSpacingAndText:  true,
	}

	if len(extraOpts) > 0 {
		userOpts := extraOpts[0]
		opts.AutoFix = autoFix || userOpts.AutoFix
		opts.ClearHighlights = clearHighlights || userOpts.ClearHighlights
		if userOpts.ExpectedFontFamily != "" {
			opts.ExpectedFontFamily = userOpts.ExpectedFontFamily
		}
		opts.CheckStructuralOrder = userOpts.CheckStructuralOrder
		opts.CheckLatexLeaks = userOpts.CheckLatexLeaks
		opts.CheckSpacingAndText = userOpts.CheckSpacingAndText
	}

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

	commentsSummary := getDocCommentsSummary(ctx, documentID)

	type RangeItem struct {
		Start int64
		End   int64
	}
	var emptyTableRanges []RangeItem
	var fontFixRanges []RangeItem
	var fontFamilyFixRanges []RangeItem
	var highlightRanges []RangeItem

	consecutiveEmptyParas := 0
	consecutiveBullets := 0

	// Máquina de estados para rastreamento da ordem canônica ABNT NBR 14724
	type docSectionState int
	const (
		statePreTextual docSectionState = iota
		stateTextual
		statePostTextual
	)
	currentState := statePreTextual
	highestChapterNum := 0
	firstChapterIndex := -1
	referencesIndex := -1

	for i, el := range doc.Body.Content {
		// 1. Análise de Tabelas
		if el.Table != nil {
			consecutiveEmptyParas = 0 // CRÍTICO: Tabelas nativas zeram a contagem de parágrafos vazios para evitar falsos positivos
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

			// Se a tabela estiver posicionada dentro ou após a seção de REFERÊNCIAS
			if opts.CheckStructuralOrder && currentState == statePostTextual {
				report.MisplacedElements = append(report.MisplacedElements,
					fmt.Sprintf("Elemento %d [%d-%d]: Tabela nativa posicionada indevidamente dentro da seção de REFERÊNCIAS (NBR 6023). Tabelas e ilustrações devem integrar os capítulos de desenvolvimento ou anexos/apêndices.",
						i, el.StartIndex, el.EndIndex))
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
							if (rgb.Red > 0.15 || rgb.Green > 0.15 || rgb.Blue > 0.15) && strings.TrimSpace(txt) != "" {
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

						// 2.3 Detecção de Família de Fonte ABNT (Arial ou Times New Roman uniforme)
						if ts.WeightedFontFamily != nil && ts.WeightedFontFamily.FontFamily != "" {
							fam := ts.WeightedFontFamily.FontFamily
							if !strings.EqualFold(fam, opts.ExpectedFontFamily) && strings.TrimSpace(txt) != "" {
								report.FontFamilyAnomalies = append(report.FontFamilyAnomalies, FontFamilyAnomalyItem{
									StartIndex:     pEl.StartIndex,
									EndIndex:       pEl.EndIndex,
									CurrentFamily:  fam,
									ExpectedFamily: opts.ExpectedFontFamily,
									Text:           strings.TrimSpace(txt),
								})
								fontFamilyFixRanges = append(fontFamilyFixRanges, RangeItem{Start: pEl.StartIndex, End: pEl.EndIndex})
							}
						}

						// 2.4 Detecção de Tamanho de Fonte ABNT
						if ts.FontSize != nil && strings.TrimSpace(txt) != "" {
							pt := ts.FontSize.Magnitude

							trimmedP := strings.TrimSpace(pText)
							isCaptionOrSource := strings.HasPrefix(trimmedP, "Figura ") ||
								strings.HasPrefix(trimmedP, "Tabela ") ||
								strings.HasPrefix(trimmedP, "Quadro ") ||
								strings.HasPrefix(trimmedP, "Fonte:") ||
								strings.HasPrefix(trimmedP, "Nota:")

							if isCaptionOrSource {
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
							} else if isPreTextualNote(trimmedP) {
								if pt < 9.5 || pt > 12.5 {
									report.FontSizeAnomalies = append(report.FontSizeAnomalies, FontSizeAnomalyItem{
										StartIndex:   pEl.StartIndex,
										EndIndex:     pEl.EndIndex,
										CurrentSize:  pt,
										ExpectedSize: 10.0,
										ContextType:  "Nota Institucional Pré-Textual (Folha de Rosto)",
										Text:         strings.TrimSpace(txt),
									})
								}
							} else if isHeading {
								if pt < 11.5 || pt > 12.5 {
									report.FontSizeAnomalies = append(report.FontSizeAnomalies, FontSizeAnomalyItem{
										StartIndex:   pEl.StartIndex,
										EndIndex:     pEl.EndIndex,
										CurrentSize:  pt,
										ExpectedSize: 12.0,
										ContextType:  "Título de Seção (ABNT NBR 14724/6024)",
										Text:         strings.TrimSpace(txt),
									})
									fontFixRanges = append(fontFixRanges, RangeItem{Start: pEl.StartIndex, End: pEl.EndIndex})
								}
							} else {
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

			// 2.5 Auditoria de Ordem Estrutural (NBR 14724)
			if opts.CheckStructuralOrder {
				upperTrim := strings.ToUpper(trimmed)

				// Verifica se é título pré-textual importante
				isPreTextualTitle := upperTrim == "RESUMO" || upperTrim == "ABSTRACT" ||
					upperTrim == "SUMÁRIO" || upperTrim == "SUMARIO" ||
					strings.HasPrefix(upperTrim, "LISTA DE ")

				if isPreTextualTitle {
					if currentState == stateTextual {
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: i,
							StartIndex:   el.StartIndex,
							EndIndex:     el.EndIndex,
							SectionTitle: trimmed,
							Severity:     "CRÍTICO",
							RuleViolated: "ABNT NBR 14724: Elemento pré-textual obrigatório localizado dentro da parte textual",
							Details:      fmt.Sprintf("O título '%s' é pré-textual e deve anteceder o primeiro capítulo numerado (iniciado no elemento %d).", trimmed, firstChapterIndex),
						})
					} else if currentState == statePostTextual {
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: i,
							StartIndex:   el.StartIndex,
							EndIndex:     el.EndIndex,
							SectionTitle: trimmed,
							Severity:     "CRÍTICO",
							RuleViolated: "ABNT NBR 14724: Elemento pré-textual obrigatório localizado na parte pós-textual",
							Details:      fmt.Sprintf("O título '%s' é pré-textual e não pode estar localizado após a seção de REFERÊNCIAS (elemento %d). Reorganize-o para antes da Introdução.", trimmed, referencesIndex),
						})
					}
				}

				// Verifica início ou progressão de capítulos numerados
				if rePrimaryHeadingNum.MatchString(trimmed) || (isHeading && reChapterNumber.MatchString(trimmed)) {
					if currentState == statePostTextual {
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: i,
							StartIndex:   el.StartIndex,
							EndIndex:     el.EndIndex,
							SectionTitle: trimmed,
							Severity:     "CRÍTICO",
							RuleViolated: "ABNT NBR 14724: Capítulo textual numerado posicionado após REFERÊNCIAS",
							Details:      fmt.Sprintf("O capítulo '%s' está após a seção de REFERÊNCIAS (elemento %d). Capítulos devem anteceder as referências bibliográficas.", trimmed, referencesIndex),
						})
					}

					currentState = stateTextual
					if firstChapterIndex < 0 {
						firstChapterIndex = i
					}

					if m := reChapterNumber.FindStringSubmatch(trimmed); len(m) > 1 {
						if num, errNum := strconv.Atoi(m[1]); errNum == nil {
							if num <= highestChapterNum {
								report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
									ElementIndex: i,
									StartIndex:   el.StartIndex,
									EndIndex:     el.EndIndex,
									SectionTitle: trimmed,
									Severity:     "ALERTA",
									RuleViolated: "ABNT NBR 6024: Numeração progressiva de capítulos em ordem não crescente ou duplicada",
									Details:      fmt.Sprintf("Capítulo %d encontrado após capítulo %d.", num, highestChapterNum),
								})
							} else {
								highestChapterNum = num
							}
						}
					}
				}

				// Verifica seção pós-textual de REFERÊNCIAS
				if strings.Contains(upperTrim, "REFERÊNCIAS") && (isHeading || strings.HasPrefix(trimmed, "3 ") || strings.HasPrefix(trimmed, "7 ") || trimmed == "REFERÊNCIAS") {
					currentState = statePostTextual
					if referencesIndex < 0 {
						referencesIndex = i
					}
				}

				// Alerta sobre ilustrações/figuras dentro de referências
				isFigureOrTableLabel := strings.HasPrefix(trimmed, "Figura ") || strings.HasPrefix(trimmed, "Tabela ") || strings.HasPrefix(trimmed, "Quadro ")
				if isFigureOrTableLabel && currentState == statePostTextual {
					report.MisplacedElements = append(report.MisplacedElements,
						fmt.Sprintf("Elemento %d [%d-%d]: Legenda de ilustração '%s' posicionada indevidamente dentro da seção de REFERÊNCIAS (NBR 6023).",
							i, el.StartIndex, el.EndIndex, trimmed))
				}
			}

			// 2.6 Auditoria de Vazamentos LaTeX (comandos, matemática, til, tags)
			if opts.CheckLatexLeaks && trimmed != "" {
				// Comandos LaTeX residuais (incluindo macros inline e fechamentos mistos)
				if matches := reLatexCmdInText.FindAllString(pText, -1); len(matches) > 0 {
					for _, m := range matches {
						report.LatexSyntaxLeaks = append(report.LatexSyntaxLeaks, LatexLeakItem{
							ElementIndex: i,
							StartIndex:   el.StartIndex,
							EndIndex:     el.EndIndex,
							LeakType:     "Comando LaTeX residual",
							Snippet:      m,
							Suggestion:   "Remover comando ou converter para texto/estilo nativo do Google Docs.",
						})
					}
				}

				// Sintaxe Matemática $...$
				if m := reMathSyntaxInText.FindString(pText); m != "" {
					report.LatexSyntaxLeaks = append(report.LatexSyntaxLeaks, LatexLeakItem{
						ElementIndex: i,
						StartIndex:   el.StartIndex,
						EndIndex:     el.EndIndex,
						LeakType:     "Sintaxe Matemática ($...$)",
						Snippet:      m,
						Suggestion:   "Substituir delimitação matemática por caracteres Unicode diretos (ex: 2 × 2 em vez de $2 \\times 2$).",
					})
				}

				// Til de espaço insecável
				if strings.Contains(pText, "~") && !strings.Contains(pText, "http") {
					report.LatexSyntaxLeaks = append(report.LatexSyntaxLeaks, LatexLeakItem{
						ElementIndex: i,
						StartIndex:   el.StartIndex,
						EndIndex:     el.EndIndex,
						LeakType:     "Espaço Insecável LaTeX (~)",
						Snippet:      "~",
						Suggestion:   "Substituir o caractere til (~) por espaço simples.",
					})
				}

				// Parênteses quebrados/duplicados
				if m := reBrokenParenthesis.FindString(pText); m != "" {
					report.LatexSyntaxLeaks = append(report.LatexSyntaxLeaks, LatexLeakItem{
						ElementIndex: i,
						StartIndex:   el.StartIndex,
						EndIndex:     el.EndIndex,
						LeakType:     "Parênteses redundantes/desbalanceados",
						Snippet:      m,
						Suggestion:   "Ajustar pontuação ou fecho de citações (ex: substituir '))' por ')').",
					})
				}

				// Marcadores artificiais soltos
				if strings.Contains(pText, "[TABELA]") || strings.Contains(pText, "[/TABELA]") || strings.Contains(pText, "[VER IMAGEM") {
					report.DanglingMarkers = append(report.DanglingMarkers,
						fmt.Sprintf("Elemento %d [%d-%d]: Marcador artificial -> %q",
							i, el.StartIndex, el.EndIndex, trimmed))
				}
			}

			// 2.7 Auditoria de Espaços Duplos e Espaços Antes de Pontuação
			if opts.CheckSpacingAndText && trimmed != "" {
				if matches := reDoubleSpace.FindAllStringIndex(pText, -1); len(matches) > 0 {
					for _, match := range matches {
						snippetStart := match[0] - 10
						if snippetStart < 0 {
							snippetStart = 0
						}
						snippetEnd := match[1] + 10
						if snippetEnd > len(pText) {
							snippetEnd = len(pText)
						}
						report.DoubleSpacingAnomalies = append(report.DoubleSpacingAnomalies, DoubleSpaceAnomalyItem{
							StartIndex: el.StartIndex + int64(match[0]),
							EndIndex:   el.StartIndex + int64(match[1]),
							IssueType:  "Espaço Duplo / Múltiplo",
							Snippet:    pText[snippetStart:snippetEnd],
						})
					}
				}

				if matches := reSpaceBeforePunct.FindAllStringSubmatchIndex(pText, -1); len(matches) > 0 {
					for _, match := range matches {
						snippetStart := match[0] - 10
						if snippetStart < 0 {
							snippetStart = 0
						}
						snippetEnd := match[1] + 10
						if snippetEnd > len(pText) {
							snippetEnd = len(pText)
						}
						report.DoubleSpacingAnomalies = append(report.DoubleSpacingAnomalies, DoubleSpaceAnomalyItem{
							StartIndex: el.StartIndex + int64(match[0]),
							EndIndex:   el.StartIndex + int64(match[1]),
							IssueType:  "Espaço Antes de Pontuação",
							Snippet:    pText[snippetStart:snippetEnd],
						})
					}
				}
			}

			// 2.8 Qualidade do Parágrafo e Alíneas ABNT
			if opts.CheckSpacingAndText && currentState == stateTextual && !isHeading {
				isCaptionOrSource := strings.HasPrefix(trimmed, "Figura ") ||
					strings.HasPrefix(trimmed, "Tabela ") ||
					strings.HasPrefix(trimmed, "Quadro ") ||
					strings.HasPrefix(trimmed, "Fonte:") ||
					strings.HasPrefix(trimmed, "Nota:")

				isBullet := p.Bullet != nil || strings.HasPrefix(trimmed, "•") || strings.HasPrefix(trimmed, "- ")

				if isBullet {
					consecutiveBullets++
					if consecutiveBullets == 4 {
						report.BulletListAnomalies = append(report.BulletListAnomalies, BulletListAnomalyItem{
							StartIndex:       el.StartIndex,
							EndIndex:         el.EndIndex,
							ConsecutiveCount: consecutiveBullets,
							Snippet:          trimmed,
							Suggestion:       "ABNT NBR 6024/14724: Enumerações sequenciais longas no corpo devem preferencialmente seguir a formatação de alíneas (a, b, c...) com recuo de 1,25 cm, finalizadas com ponto e vírgula (;).",
						})
					}
				} else {
					consecutiveBullets = 0
				}

				isHeadingItem := isHeading ||
					rePrimaryHeadingNum.MatchString(trimmed) ||
					rePrimaryHeadingNoNum.MatchString(trimmed) ||
					reSecondaryHeading.MatchString(trimmed) ||
					reTertiaryHeading.MatchString(trimmed) ||
					reQuaternaryHeading.MatchString(trimmed)

				isAlineaOrListItem := regexp.MustCompile(`^[a-z]\)\s+`).MatchString(trimmed) ||
					regexp.MustCompile(`^[0-9]+[\.\)]\s+`).MatchString(trimmed) ||
					regexp.MustCompile(`^(?:RF|RN|RNF|RE|C)[0-9\-]+:`).MatchString(trimmed) ||
					strings.HasSuffix(trimmed, ";")

				isFormulaOrMath := strings.HasPrefix(trimmed, "H₀:") || strings.HasPrefix(trimmed, "H₁:") ||
					strings.HasPrefix(trimmed, "H0:") || strings.HasPrefix(trimmed, "H1:") ||
					strings.HasPrefix(trimmed, "TCD_") || strings.Contains(trimmed, " = ") ||
					strings.Contains(trimmed, " > ") || strings.Contains(trimmed, " < ")

				isListOrTableIntro := strings.HasSuffix(trimmed, ":") ||
					strings.HasPrefix(trimmed, "A Tabela ") || strings.HasPrefix(trimmed, "O Quadro ") || strings.HasPrefix(trimmed, "A Figura ")

				if !isCaptionOrSource && !isBullet && !isPreTextualNote(trimmed) && !isHeadingItem && !isAlineaOrListItem && !isFormulaOrMath && !isListOrTableIntro && trimmed != "" {
					runeLen := utf8.RuneCountInString(trimmed)
					isSubheadingOrStage := strings.HasPrefix(trimmed, "Etapa de ") || strings.HasPrefix(trimmed, "Fase de ")
					if runeLen < 100 && !strings.Contains(trimmed, "[TABELA]") && !isSubheadingOrStage {
						report.ParagraphQualityAnomalies = append(report.ParagraphQualityAnomalies, ParagraphQualityItem{
							ElementIndex: i,
							StartIndex:   el.StartIndex,
							EndIndex:     el.EndIndex,
							IssueType:    "Parágrafo Órfão / Frase Curta (< 100 caracteres)",
							Length:       runeLen,
							Snippet:      trimmed,
							Suggestion:   "ABNT NBR 14724: Parágrafos acadêmicos devem conter introdução, desenvolvimento e conclusão da ideia. Evite períodos soltos de uma só linha no corpo do trabalho.",
						})
					} else if runeLen > 1500 {
						snippet := trimmed
						if len(snippet) > 80 {
							snippet = snippet[:80] + "..."
						}
						report.ParagraphQualityAnomalies = append(report.ParagraphQualityAnomalies, ParagraphQualityItem{
							ElementIndex: i,
							StartIndex:   el.StartIndex,
							EndIndex:     el.EndIndex,
							IssueType:    "Parágrafo Excessivamente Longo (> 1500 caracteres)",
							Length:       runeLen,
							Snippet:      snippet,
							Suggestion:   "Recomenda-se dividir parágrafos com mais de 1500 caracteres para favorecer a clareza e fluidez da leitura acadêmica.",
						})
					}
				}
			}

			// Parágrafo vazio ou quebra múltipla
			// ABNT NBR 14724: Espaçamentos verticais na Capa e Folha de Rosto (elementos pré-textuais)
			// são canônicos e regulamentares para distribuição dos blocos. A auditoria de quebras
			// excessivas atua estritamente na parte textual e pós-textual para evitar falsos positivos.
			if trimmed == "" {
				if currentState != statePreTextual {
					consecutiveEmptyParas++
					if consecutiveEmptyParas >= 2 {
						report.ConsecutiveBlankLines = append(report.ConsecutiveBlankLines,
							fmt.Sprintf("Elemento %d [%d-%d]: Quebra múltipla / Linha em branco excessiva (#%d)",
								i, el.StartIndex, el.EndIndex, consecutiveEmptyParas))
					}
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
		}
	}

	// 2.10 Auditoria da Capa Oficial IFS e Folha de Rosto (NBR 14724 + config.tex)
	if opts.CheckStructuralOrder {
		hasCover, hasLogo, hasTitlePage := HasABNTCoverAndLogo(doc)
		var coverIssues []string

		if !hasCover {
			coverIssues = append(coverIssues, "Capa institucional com dados de identificação ausente")
		}
		if !hasLogo {
			coverIssues = append(coverIssues, "Logo oficial do IFS (3.2cm) não localizada no topo")
		}
		if !hasTitlePage {
			coverIssues = append(coverIssues, "Folha de Rosto com nota institucional ausente")
		}

		// Se os elementos básicos existirem, audita alinhamentos, quebras de página e regras de formatação
		if hasCover && hasLogo && hasTitlePage {
			// Carrega metadados oficiais para comparar se há coorientador
			meta, _ := latex.ParseTCCMetadata("")

			// Máquina de auditoria estrita de páginas pré-textuais NBR 14724
			inCover := true
			coverYearFound := false
			coverYearEndIndex := int64(-1)
			coverHasPageBreak := false
			titlePageYearFound := false
			titlePageHasPageBreak := false
			titlePageStarted := false

			for idx, el := range doc.Body.Content {
				if firstChapterIndex >= 0 && idx >= firstChapterIndex {
					break
				}
				if el.Paragraph != nil {
					p := el.Paragraph
					align := ""
					indentStart := float64(0)
					indentFirstLine := float64(0)
					if p.ParagraphStyle != nil {
						align = p.ParagraphStyle.Alignment
						if p.ParagraphStyle.IndentStart != nil {
							indentStart = p.ParagraphStyle.IndentStart.Magnitude
						}
						if p.ParagraphStyle.IndentFirstLine != nil {
							indentFirstLine = p.ParagraphStyle.IndentFirstLine.Magnitude
						}
					}

					pText := ""
					hasInline := false
					pHasBreak := false
					for _, pEl := range p.Elements {
						if pEl.TextRun != nil {
							pText += pEl.TextRun.Content
						}
						if pEl.InlineObjectElement != nil {
							hasInline = true
						}
						if pEl.PageBreak != nil {
							pHasBreak = true
						}
					}
					trimmed := strings.TrimSpace(pText)

					// Audita alinhamento da Logo no topo da Capa
					if hasInline {
						if align != "CENTER" {
							coverIssues = append(coverIssues, "Logo oficial do IFS não está centralizada no topo da Capa")
							report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
								ElementIndex: idx,
								StartIndex:   el.StartIndex,
								EndIndex:     el.EndIndex,
								SectionTitle: "Logo Oficial IFS",
								Severity:     "ALERTA",
								RuleViolated: "ABNT NBR 14724 / IFS config.tex: Logo oficial deve estar rigorosamente centralizada no topo da Capa",
								Details:      "O parágrafo da logo está com alinhamento '" + align + "'. Execute 'format_document_abnt' para centralizar.",
							})
						}
					}

					if inCover {
						// Estamos na Página 1 (Capa)
						if trimmed == "2026" || (meta != nil && trimmed == strings.TrimSpace(meta.Data)) {
							coverYearFound = true
							coverYearEndIndex = el.EndIndex
							if pHasBreak {
								coverHasPageBreak = true
								inCover = false
							}
						} else if coverYearFound {
							// Se já encontrou o ano da Capa e não houve quebra de página, qualquer texto aqui vazou para a página 1!
							if pHasBreak {
								coverHasPageBreak = true
								inCover = false
							} else if trimmed != "" {
								coverIssues = append(coverIssues, fmt.Sprintf("Elemento indevido vazou para a Capa (Página 1): %q", trimmed))
								report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
									ElementIndex: idx,
									StartIndex:   el.StartIndex,
									EndIndex:     el.EndIndex,
									SectionTitle: "Capa (Página 1)",
									Severity:     "ALERTA",
									RuleViolated: "ABNT NBR 14724: A Capa deve conter exclusivamente Logo, Instituição, Autor, Título, Local e Ano, finalizando com Quebra de Página",
									Details:      fmt.Sprintf("O texto %q apareceu antes da Quebra de Página da Capa. Execute 'format_document_abnt' para inserir a quebra.", trimmed),
								})
								inCover = false // Força transição para auditar o restante como Folha de Rosto
							}
						}
					} else {
						// Estamos na Página 2 (Folha de Rosto)
						if !titlePageStarted && trimmed != "" {
							titlePageStarted = true
							// Primeiro elemento de texto deve ser o Autor
							if meta != nil && !strings.EqualFold(trimmed, strings.TrimSpace(meta.Autor)) && !strings.Contains(strings.ToUpper(trimmed), "JOSÉ GUSTAVO") {
								coverIssues = append(coverIssues, fmt.Sprintf("A Folha de Rosto deve iniciar com o Nome do Autor no topo da Página 2 (iniciou com %q)", trimmed))
								report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
									ElementIndex: idx,
									StartIndex:   el.StartIndex,
									EndIndex:     el.EndIndex,
									SectionTitle: "Folha de Rosto (Página 2)",
									Severity:     "ALERTA",
									RuleViolated: "ABNT NBR 14724: A Folha de Rosto deve iniciar com o Nome do Autor no topo da Página 2",
									Details:      "Verifique se o nome do autor não vazou para o rodapé da Capa por falta de Quebra de Página.",
								})
							}
						}

						// Audita nota institucional da folha de rosto
						if strings.HasPrefix(trimmed, "Trabalho de Conclusão de Curso apresentado") ||
							strings.Contains(trimmed, "requisito parcial para obtenção") {
							if !strings.Contains(trimmed, "obtenção do grau") {
								coverIssues = append(coverIssues, "Nota institucional da Folha de Rosto está incompleta/truncada")
								report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
									ElementIndex: idx,
									StartIndex:   el.StartIndex,
									EndIndex:     el.EndIndex,
									SectionTitle: "Folha de Rosto (Nota de Apresentação Incompleta)",
									Severity:     "ALERTA",
									RuleViolated: "ABNT NBR 14724 / metadados.tex: A nota institucional deve conter o preâmbulo completo com curso, campus e objetivo do grau",
									Details:      "O texto da nota foi cortado antes da conclusão. Execute 'format_document_abnt' para restaurar o preâmbulo completo.",
								})
							}
							if indentStart < 200.0 || indentFirstLine < 200.0 || align != "JUSTIFIED" {
								coverIssues = append(coverIssues, "Nota institucional da Folha de Rosto sem recuo de 8.0cm à direita / justificado (ou 1ª linha desalinhada)")
								report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
									ElementIndex: idx,
									StartIndex:   el.StartIndex,
									EndIndex:     el.EndIndex,
									SectionTitle: "Folha de Rosto (Nota de Apresentação)",
									Severity:     "ALERTA",
									RuleViolated: "ABNT NBR 14724: Nota de apresentação deve possuir recuo de 8.0cm à esquerda (metade da página) com 1ª linha alinhada e texto justificado",
									Details:      fmt.Sprintf("Recuo corpo: %.1f pt, Recuo 1ª linha: %.1f pt (esperado: >= 212.6 pt / 8.0cm).", indentStart, indentFirstLine),
								})
							}
						}

						// Audita alinhamento do Orientador na folha de rosto
						if strings.HasPrefix(trimmed, "Orientador:") || strings.HasPrefix(trimmed, "Orientadora:") {
							if indentStart < 200.0 || indentFirstLine < 200.0 {
								coverIssues = append(coverIssues, "Orientador na Folha de Rosto não está alinhado à direita junto ao bloco da nota institucional")
								report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
									ElementIndex: idx,
									StartIndex:   el.StartIndex,
									EndIndex:     el.EndIndex,
									SectionTitle: "Folha de Rosto (Orientador)",
									Severity:     "ALERTA",
									RuleViolated: "ABNT NBR 14724: O Orientador deve ser formatado no mesmo alinhamento de bloco à direita da nota de apresentação (recuo 8.0cm)",
									Details:      fmt.Sprintf("Recuo corpo: %.1f pt, Recuo 1ª linha: %.1f pt (esperado: >= 212.6 pt / 8.0cm).", indentStart, indentFirstLine),
								})
							}
						}

						// Audita coorientador fake / placeholder
						if strings.Contains(trimmed, "Nome do Coorientador") ||
							(meta != nil && meta.Coorientador == "" && strings.HasPrefix(trimmed, "Coorientador:")) {
							coverIssues = append(coverIssues, "Folha de Rosto contém menção a Coorientador inexistente nos metadados")
							report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
								ElementIndex: idx,
								StartIndex:   el.StartIndex,
								EndIndex:     el.EndIndex,
								SectionTitle: "Folha de Rosto (Coorientador)",
								Severity:     "ALERTA",
								RuleViolated: "Coorientador é opcional e não consta ativo no metadados.tex",
								Details:      fmt.Sprintf("Remova a linha %q ou configure o coorientador real em metadados.tex.", trimmed),
							})
						}

						// Encerramento da Folha de Rosto
						if trimmed == "2026" || (meta != nil && trimmed == strings.TrimSpace(meta.Data)) {
							titlePageYearFound = true
							if pHasBreak {
								titlePageHasPageBreak = true
							}
						} else if titlePageYearFound {
							if pHasBreak {
								titlePageHasPageBreak = true
							} else if (trimmed == "RESUMO" || strings.HasPrefix(trimmed, "RESUMO")) && !titlePageHasPageBreak {
								coverIssues = append(coverIssues, "Ausência de Quebra de Página após a Folha de Rosto (RESUMO está na mesma página)")
								report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
									ElementIndex: idx,
									StartIndex:   el.StartIndex,
									EndIndex:     el.EndIndex,
									SectionTitle: "Folha de Rosto / RESUMO",
									Severity:     "ALERTA",
									RuleViolated: "ABNT NBR 14724: A Folha de Rosto deve encerrar com Quebra de Página, iniciando o RESUMO na Página 3",
									Details:      "O RESUMO apareceu sem que houvesse Quebra de Página após o Ano da Folha de Rosto.",
								})
							}
						}
					}
				}
			}

			// Verificação das quebras regulamentares na Folha de Rosto (7 entre Autor e Título; 5 entre Título e Preâmbulo)
			if coverYearFound && coverYearEndIndex > 0 {
				var autorFREl, tituloFREl, preambuloFREl *docs.StructuralElement
				for _, el := range doc.Body.Content {
					if el.StartIndex <= coverYearEndIndex {
						continue
					}
					if el.Paragraph != nil {
						txt := ""
						for _, pel := range el.Paragraph.Elements {
							if pel.TextRun != nil {
								txt += pel.TextRun.Content
							}
						}
						trim := strings.TrimSpace(txt)
						if trim == "" {
							continue
						}
						if autorFREl == nil && (strings.Contains(strings.ToUpper(trim), "JOSÉ GUSTAVO") || (meta != nil && strings.EqualFold(trim, strings.ToUpper(meta.Autor)))) {
							autorFREl = el
							continue
						}
						if autorFREl != nil && tituloFREl == nil && (strings.Contains(strings.ToUpper(trim), "AVALIAÇÃO DE UMA") || (meta != nil && strings.EqualFold(trim, strings.ToUpper(meta.Titulo)))) {
							tituloFREl = el
							continue
						}
						if tituloFREl != nil && preambuloFREl == nil && (strings.HasPrefix(trim, "Trabalho de Conclusão") || isPreTextualNote(trim)) {
							preambuloFREl = el
							break
						}
					}
				}

				if autorFREl != nil && tituloFREl != nil {
					emptyAT := 0
					for _, el := range doc.Body.Content {
						if el.StartIndex >= autorFREl.EndIndex && el.EndIndex <= tituloFREl.StartIndex {
							if el.Paragraph != nil {
								t := ""
								for _, pel := range el.Paragraph.Elements {
									if pel.TextRun != nil {
										t += pel.TextRun.Content
									}
								}
								if strings.TrimSpace(t) == "" {
									emptyAT++
								}
							}
						}
					}
					if emptyAT != 7 {
						coverIssues = append(coverIssues, fmt.Sprintf("Folha de Rosto: espaçamento entre Autor e Título irregular (%d quebras; esperado: exatamente 7)", emptyAT))
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: 0,
							StartIndex:   autorFREl.StartIndex,
							EndIndex:     tituloFREl.EndIndex,
							SectionTitle: "Folha de Rosto (Espaçamento Autor -> Título)",
							Severity:     "ALERTA",
							RuleViolated: "ABNT NBR 14724: Devem existir exatamente 7 quebras de linhas entre o Nome do Autor e o Título do Trabalho na Folha de Rosto",
							Details:      fmt.Sprintf("Atualmente existem %d quebras vazias. Execute 'format_document_abnt' para calibrar para 7.", emptyAT),
						})
					}
				}

				if tituloFREl != nil && preambuloFREl != nil {
					emptyTP := 0
					for _, el := range doc.Body.Content {
						if el.StartIndex >= tituloFREl.EndIndex && el.EndIndex <= preambuloFREl.StartIndex {
							if el.Paragraph != nil {
								t := ""
								for _, pel := range el.Paragraph.Elements {
									if pel.TextRun != nil {
										t += pel.TextRun.Content
									}
								}
								if strings.TrimSpace(t) == "" {
									emptyTP++
								}
							}
						}
					}
					if emptyTP != 5 {
						coverIssues = append(coverIssues, fmt.Sprintf("Folha de Rosto: espaçamento entre Título e Preâmbulo irregular (%d quebras; esperado: exatamente 5)", emptyTP))
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: 0,
							StartIndex:   tituloFREl.StartIndex,
							EndIndex:     preambuloFREl.EndIndex,
							SectionTitle: "Folha de Rosto (Espaçamento Título -> Preâmbulo)",
							Severity:     "ALERTA",
							RuleViolated: "ABNT NBR 14724: Devem existir exatamente 5 quebras de linhas entre o Título e a Nota de Apresentação / Preâmbulo",
							Details:      fmt.Sprintf("Atualmente existem %d quebras vazias. Execute 'format_document_abnt' para calibrar para 5.", emptyTP),
						})
					}
				}
			}

			// Verificações finais de presença de quebra de página
			if coverYearFound && !coverHasPageBreak {
				coverIssues = append(coverIssues, "Ausência de Quebra de Página obrigatória (NBR 14724) após o Ano da Capa")
			}
			if titlePageYearFound && !titlePageHasPageBreak {
				coverIssues = append(coverIssues, "Ausência de Quebra de Página obrigatória (NBR 14724) após o Ano da Folha de Rosto")
			}
		}

		if len(coverIssues) > 0 {
			report.CoverStatus = "Incompleta ou Desalinhada (" + strings.Join(coverIssues, "; ") + ")"
		} else {
			report.CoverStatus = "Capa Oficial IFS com Logo (3.2cm centralizada) e Folha de Rosto NBR 14724 em plena conformidade."
		}
	}

	// 3. Execução de Autocorreção (se autoFix == true)
	if opts.AutoFix {
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

		// 3.3 Substitui espaços duplos e espaços antes de pontuação
		var textCleanReqs []*docs.Request
		if len(report.DoubleSpacingAnomalies) > 0 {
			// Dois passos para colapsar múltiplos espaços
			textCleanReqs = append(textCleanReqs,
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "  "}, ReplaceText: " "}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "  "}, ReplaceText: " "}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: " ,"}, ReplaceText: ","}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: " ."}, ReplaceText: "."}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: " ;"}, ReplaceText: ";"}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: " :"}, ReplaceText: ":"}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "))"}, ReplaceText: ")"}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "(("}, ReplaceText: "("}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "( "}, ReplaceText: "("}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: " )"}, ReplaceText: ")"}},
			)
		}

		// 3.4 Limpa vazamentos de sintaxe simples de LaTeX e marcadores
		if len(report.LatexSyntaxLeaks) > 0 || len(report.DanglingMarkers) > 0 {
			textCleanReqs = append(textCleanReqs,
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "$2 \\times 2$"}, ReplaceText: "2 × 2"}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "\\times"}, ReplaceText: "×"}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "\\ge"}, ReplaceText: "≥"}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "\\le"}, ReplaceText: "≤"}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "\\textordmasculine"}, ReplaceText: "º"}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "\\textordfeminine"}, ReplaceText: "ª"}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "\\allowbreak"}, ReplaceText: ""}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "~"}, ReplaceText: " "}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "[TABELA]"}, ReplaceText: ""}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "[/TABELA]"}, ReplaceText: ""}},
				&docs.Request{ReplaceAllText: &docs.ReplaceAllTextRequest{ContainsText: &docs.SubstringMatchCriteria{MatchCase: true, Text: "[VER IMAGEM"}, ReplaceText: ""}},
			)
		}

		if len(textCleanReqs) > 0 {
			if _, errTxt := docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: textCleanReqs}).Do(); errTxt == nil {
				fixedCount += len(textCleanReqs)
			}
		}

		// 3.5 Normaliza tamanhos de fonte do corpo para 12 pt
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

		// 3.6 Normaliza família de fonte para o padrão ABNT (Arial / Times)
		for _, ffr := range fontFamilyFixRanges {
			requests = append(requests, &docs.Request{
				UpdateTextStyle: &docs.UpdateTextStyleRequest{
					Range: &docs.Range{
						StartIndex: ffr.Start,
						EndIndex:   ffr.End,
					},
					TextStyle: &docs.TextStyle{
						WeightedFontFamily: &docs.WeightedFontFamily{
							FontFamily: opts.ExpectedFontFamily,
						},
					},
					Fields: "weightedFontFamily",
				},
			})
		}

		// 3.7 Limpeza de marca-texto se solicitado explicitamente
		if opts.ClearHighlights && len(highlightRanges) > 0 {
			for _, hr := range highlightRanges {
				requests = append(requests, &docs.Request{
					UpdateTextStyle: &docs.UpdateTextStyleRequest{
						Range: &docs.Range{
							StartIndex: hr.Start,
							EndIndex:   hr.End,
						},
						TextStyle: &docs.TextStyle{
							// Para remover marca-texto no Google Docs, não define BackgroundColor
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
			len(report.FontSizeAnomalies) + len(report.FontFamilyAnomalies) + len(report.TextColorAnomalies) +
			len(report.StructuralOrderViolations) + len(report.MisplacedElements) + len(report.LatexSyntaxLeaks) +
			len(report.DoubleSpacingAnomalies) + len(report.ParagraphQualityAnomalies) + len(report.BulletListAnomalies)

		report.Message = fmt.Sprintf("Auditoria concluída. %d anomalias/destaques identificados no documento (autoFix=false).", totalAnomalies)
	}

	return report, nil
}
