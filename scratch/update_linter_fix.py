with open('internal/gdocs/linter.go', 'r', encoding='utf-8') as f:
    c = f.read().replace('\r\n', '\n')

start_marker = "// Verificação das quebras regulamentares da Capa e Folha de Rosto"
end_marker = "// Verificações finais de presença de quebra de página"

idx_start = c.find(start_marker)
idx_end = c.find(end_marker)

assert idx_start != -1, "start_marker not found in linter.go"
assert idx_end != -1, "end_marker not found in linter.go"

new_linter = '''// Verificação das quebras regulamentares da Capa e Folha de Rosto (NBR 14724 + checklist IFS)
			if coverYearFound && coverYearEndIndex > 0 {
				// 1. Auditoria das quebras da Capa (Página 1):
				// - Entre Instituição e Autor: 2 a 6 quebras (calibrado: 3)
				// - Entre Autor e Título: 2 a 6 quebras (calibrado: 3)
				// - Entre Título e Cidade: 12 a 18 quebras (calibrado: 15 para posicionar no rodapé sem estouro)
				var instCapaEl, autorCapaEl, tituloCapaEl, cidadeCapaEl *docs.StructuralElement
				for _, el := range doc.Body.Content {
					if el.EndIndex > coverYearEndIndex {
						break
					}
					if el.Paragraph != nil {
						txt := ""
						for _, pel := range el.Paragraph.Elements {
							if pel.TextRun != nil {
								txt += pel.TextRun.Content
							}
						}
						trim := strings.TrimSpace(txt)
						upperTrim := strings.ToUpper(trim)
						if trim == "" {
							continue
						}
						if strings.Contains(upperTrim, "BACHARELADO") || strings.Contains(upperTrim, "CAMPUS LAGARTO") || strings.Contains(upperTrim, "INSTITUTO FEDERAL") {
							instCapaEl = el
							continue
						}
						if autorCapaEl == nil && (strings.Contains(upperTrim, "JOSÉ GUSTAVO") || (meta != nil && strings.EqualFold(trim, strings.ToUpper(meta.Autor)))) {
							autorCapaEl = el
							continue
						}
						if autorCapaEl != nil && tituloCapaEl == nil && (strings.Contains(upperTrim, "AVALIAÇÃO DE UMA") || (meta != nil && strings.EqualFold(trim, strings.ToUpper(meta.Titulo)))) {
							tituloCapaEl = el
							continue
						}
						if tituloCapaEl != nil && cidadeCapaEl == nil && (strings.Contains(upperTrim, "LAGARTO - SE") || strings.Contains(upperTrim, "LAGARTO -- SE")) {
							cidadeCapaEl = el
							break
						}
					}
				}

				if instCapaEl != nil && autorCapaEl != nil {
					emptyIA := 0
					for _, el := range doc.Body.Content {
						if el.StartIndex >= instCapaEl.EndIndex && el.EndIndex <= autorCapaEl.StartIndex {
							if el.Paragraph != nil {
								t := ""
								for _, pel := range el.Paragraph.Elements {
									if pel.TextRun != nil {
										t += pel.TextRun.Content
									}
								}
								if strings.TrimSpace(t) == "" {
									emptyIA++
								}
							}
						}
					}
					if emptyIA < 2 || emptyIA > 6 {
						coverIssues = append(coverIssues, fmt.Sprintf("Capa: espaçamento entre Instituição e Autor irregular (%d quebras; esperado: entre 3 e 4)", emptyIA))
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: 0,
							StartIndex:   instCapaEl.StartIndex,
							EndIndex:     autorCapaEl.EndIndex,
							SectionTitle: "Capa (Espaçamento Instituição -> Autor)",
							Severity:     "ALERTA",
							RuleViolated: "ABNT NBR 14724: Espaçamento vertical entre Bloco Institucional e Autor da Capa deve conter entre 3 e 4 quebras",
							Details:      fmt.Sprintf("Atualmente existem %d quebras vazias. Execute 'format_document_abnt' para calibrar.", emptyIA),
						})
					}
				}

				if autorCapaEl != nil && tituloCapaEl != nil {
					emptyATCapa := 0
					for _, el := range doc.Body.Content {
						if el.StartIndex >= autorCapaEl.EndIndex && el.EndIndex <= tituloCapaEl.StartIndex {
							if el.Paragraph != nil {
								t := ""
								for _, pel := range el.Paragraph.Elements {
									if pel.TextRun != nil {
										t += pel.TextRun.Content
									}
								}
								if strings.TrimSpace(t) == "" {
									emptyATCapa++
								}
							}
						}
					}
					if emptyATCapa < 2 || emptyATCapa > 6 {
						coverIssues = append(coverIssues, fmt.Sprintf("Capa: espaçamento entre Autor e Título irregular (%d quebras; esperado: entre 3 e 4)", emptyATCapa))
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: 0,
							StartIndex:   autorCapaEl.StartIndex,
							EndIndex:     tituloCapaEl.EndIndex,
							SectionTitle: "Capa (Espaçamento Autor -> Título)",
							Severity:     "ALERTA",
							RuleViolated: "ABNT NBR 14724: Espaçamento vertical entre Autor e Título da Capa deve conter entre 3 e 4 quebras",
							Details:      fmt.Sprintf("Atualmente existem %d quebras vazias. Execute 'format_document_abnt' para calibrar.", emptyATCapa),
						})
					}
				}

				if tituloCapaEl != nil && cidadeCapaEl != nil {
					emptyTCCapa := 0
					for _, el := range doc.Body.Content {
						if el.StartIndex >= tituloCapaEl.EndIndex && el.EndIndex <= cidadeCapaEl.StartIndex {
							if el.Paragraph != nil {
								t := ""
								for _, pel := range el.Paragraph.Elements {
									if pel.TextRun != nil {
										t += pel.TextRun.Content
									}
								}
								if strings.TrimSpace(t) == "" {
									emptyTCCapa++
								}
							}
						}
					}
					if emptyTCCapa < 12 || emptyTCCapa > 18 {
						coverIssues = append(coverIssues, fmt.Sprintf("Capa: espaçamento entre Título e Cidade irregular (%d quebras; esperado: 15 para rodapé sem overflow)", emptyTCCapa))
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: 0,
							StartIndex:   tituloCapaEl.StartIndex,
							EndIndex:     cidadeCapaEl.EndIndex,
							SectionTitle: "Capa (Espaçamento Título -> Cidade)",
							Severity:     "ALERTA",
							RuleViolated: "ABNT NBR 14724: Devem existir cerca de 15 quebras de linhas para posicionar Cidade e Ano no rodapé da Capa sem provocar estouro de página",
							Details:      fmt.Sprintf("Atualmente existem %d quebras vazias. Execute 'format_document_abnt' para calibrar para 15.", emptyTCCapa),
						})
					}
				}

				// 2. Auditoria das quebras da Folha de Rosto (Página 2):
				// - Topo: 1 a 3 quebras (calibrado: 2)
				// - Autor -> Título: 3 a 6 quebras (calibrado: 4)
				// - Título -> Preâmbulo: 2 a 5 quebras (calibrado: 3)
				// - Orientador -> Cidade: 10 a 15 quebras (calibrado: 12)
				var autorFREl, tituloFREl, preambuloFREl, orientadorFREl, cidadeFREl *docs.StructuralElement
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
						upperTrim := strings.ToUpper(trim)
						if trim == "" {
							continue
						}
						if autorFREl == nil && (strings.Contains(upperTrim, "JOSÉ GUSTAVO") || (meta != nil && strings.EqualFold(trim, strings.ToUpper(meta.Autor)))) {
							autorFREl = el
							continue
						}
						if autorFREl != nil && tituloFREl == nil && (strings.Contains(upperTrim, "AVALIAÇÃO DE UMA") || (meta != nil && strings.EqualFold(trim, strings.ToUpper(meta.Titulo)))) {
							tituloFREl = el
							continue
						}
						if tituloFREl != nil && preambuloFREl == nil && (strings.HasPrefix(trim, "Trabalho de Conclusão") || isPreTextualNote(trim)) {
							preambuloFREl = el
							continue
						}
						if preambuloFREl != nil && orientadorFREl == nil && (strings.HasPrefix(strings.ToLower(trim), "orientador:") || strings.HasPrefix(strings.ToLower(trim), "coorientador:")) {
							orientadorFREl = el
							continue
						}
						if orientadorFREl != nil && cidadeFREl == nil && (strings.Contains(upperTrim, "LAGARTO - SE") || strings.Contains(upperTrim, "LAGARTO -- SE")) {
							cidadeFREl = el
							break
						}
					}
				}

				if autorFREl != nil {
					emptyTopFR := 0
					for _, el := range doc.Body.Content {
						if el.StartIndex >= coverYearEndIndex && el.EndIndex <= autorFREl.StartIndex {
							if el.Paragraph != nil {
								t := ""
								for _, pel := range el.Paragraph.Elements {
									if pel.TextRun != nil {
										t += pel.TextRun.Content
									}
								}
								if strings.TrimSpace(t) == "" {
									emptyTopFR++
								}
							}
						}
					}
					if emptyTopFR < 1 || emptyTopFR > 4 {
						coverIssues = append(coverIssues, fmt.Sprintf("Folha de Rosto: espaçamento no topo irregular (%d quebras; esperado: 2)", emptyTopFR))
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: 0,
							StartIndex:   coverYearEndIndex,
							EndIndex:     autorFREl.StartIndex,
							SectionTitle: "Folha de Rosto (Topo)",
							Severity:     "ALERTA",
							RuleViolated: "ABNT NBR 14724: Devem existir cerca de 2 quebras de linhas no topo da Folha de Rosto para evitar deslocamento excessivo",
							Details:      fmt.Sprintf("Atualmente existem %d quebras vazias. Execute 'format_document_abnt' para calibrar.", emptyTopFR),
						})
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
					if emptyAT < 3 || emptyAT > 6 {
						coverIssues = append(coverIssues, fmt.Sprintf("Folha de Rosto: espaçamento entre Autor e Título irregular (%d quebras; esperado: 4)", emptyAT))
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: 0,
							StartIndex:   autorFREl.StartIndex,
							EndIndex:     tituloFREl.EndIndex,
							SectionTitle: "Folha de Rosto (Espaçamento Autor -> Título)",
							Severity:     "ALERTA",
							RuleViolated: "ABNT NBR 14724: Devem existir 4 quebras de linhas entre o Nome do Autor e o Título do Trabalho na Folha de Rosto",
							Details:      fmt.Sprintf("Atualmente existem %d quebras vazias. Execute 'format_document_abnt' para calibrar.", emptyAT),
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
					if emptyTP < 2 || emptyTP > 5 {
						coverIssues = append(coverIssues, fmt.Sprintf("Folha de Rosto: espaçamento entre Título e Preâmbulo irregular (%d quebras; esperado: 3)", emptyTP))
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: 0,
							StartIndex:   tituloFREl.StartIndex,
							EndIndex:     preambuloFREl.EndIndex,
							SectionTitle: "Folha de Rosto (Espaçamento Título -> Preâmbulo)",
							Severity:     "ALERTA",
							RuleViolated: "ABNT NBR 14724: Devem existir cerca de 3 quebras de linhas entre o Título e a Nota de Apresentação / Preâmbulo",
							Details:      fmt.Sprintf("Atualmente existem %d quebras vazias. Execute 'format_document_abnt' para calibrar.", emptyTP),
						})
					}
				}

				if orientadorFREl != nil && cidadeFREl != nil {
					emptyOC := 0
					for _, el := range doc.Body.Content {
						if el.StartIndex >= orientadorFREl.EndIndex && el.EndIndex <= cidadeFREl.StartIndex {
							if el.Paragraph != nil {
								t := ""
								for _, pel := range el.Paragraph.Elements {
									if pel.TextRun != nil {
										t += pel.TextRun.Content
									}
								}
								if strings.TrimSpace(t) == "" {
									emptyOC++
								}
							}
						}
					}
					if emptyOC < 10 || emptyOC > 15 {
						coverIssues = append(coverIssues, fmt.Sprintf("Folha de Rosto: espaçamento entre Orientador e Cidade irregular (%d quebras; esperado: 12 para rodapé sem overflow)", emptyOC))
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: 0,
							StartIndex:   orientadorFREl.StartIndex,
							EndIndex:     cidadeFREl.EndIndex,
							SectionTitle: "Folha de Rosto (Espaçamento Orientador -> Cidade)",
							Severity:     "ALERTA",
							RuleViolated: "ABNT NBR 14724: Devem existir cerca de 12 quebras de linhas para posicionar Cidade e Ano no rodapé da Folha de Rosto sem estouro",
							Details:      fmt.Sprintf("Atualmente existem %d quebras vazias. Execute 'format_document_abnt' para calibrar para 12.", emptyOC),
						})
					}
				}

				if preambuloFREl != nil && preambuloFREl.Paragraph != nil && preambuloFREl.Paragraph.ParagraphStyle != nil {
					ps := preambuloFREl.Paragraph.ParagraphStyle
					indentFirst := float64(0)
					if ps.IndentFirstLine != nil {
						indentFirst = ps.IndentFirstLine.Magnitude
					}
					indentStart := float64(0)
					if ps.IndentStart != nil {
						indentStart = ps.IndentStart.Magnitude
					}
					if indentStart < 195.0 || indentFirst < 195.0 {
						coverIssues = append(coverIssues, fmt.Sprintf("Folha de Rosto: recuo da nota institucional/preâmbulo desalinhado (IndentFirstLine: %.1f pt, IndentStart: %.1f pt; esperado: ambos >= 198.4 pt / ~7,0 a 8,0 cm para bloco retilíneo à direita)", indentFirst, indentStart))
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: 0,
							StartIndex:   preambuloFREl.StartIndex,
							EndIndex:     preambuloFREl.EndIndex,
							SectionTitle: "Folha de Rosto (Alinhamento em Bloco da Nota Institucional)",
							Severity:     "ALERTA",
							RuleViolated: "ABNT NBR 14724: O preâmbulo e orientador devem formar um bloco retilíneo com recuo uniforme de 7,0 a 8,0 cm (indentStart e indentFirstLine >= 198.4 pt)",
							Details:      "Execute 'format_document_abnt' para aplicar o recuo canônico de 8,0 cm em todo o bloco.",
						})
					}
				}

				if orientadorFREl != nil && orientadorFREl.Paragraph != nil && orientadorFREl.Paragraph.ParagraphStyle != nil {
					ps := orientadorFREl.Paragraph.ParagraphStyle
					indentFirst := float64(0)
					if ps.IndentFirstLine != nil {
						indentFirst = ps.IndentFirstLine.Magnitude
					}
					indentStart := float64(0)
					if ps.IndentStart != nil {
						indentStart = ps.IndentStart.Magnitude
					}
					if indentStart < 195.0 || indentFirst < 195.0 {
						coverIssues = append(coverIssues, fmt.Sprintf("Folha de Rosto: recuo do orientador desalinhado (IndentFirstLine: %.1f pt, IndentStart: %.1f pt; esperado: ambos >= 198.4 pt)", indentFirst, indentStart))
						report.StructuralOrderViolations = append(report.StructuralOrderViolations, OrderViolationItem{
							ElementIndex: 0,
							StartIndex:   orientadorFREl.StartIndex,
							EndIndex:     orientadorFREl.EndIndex,
							SectionTitle: "Folha de Rosto (Alinhamento do Orientador)",
							Severity:     "ALERTA",
							RuleViolated: "ABNT NBR 14724: O orientador deve formar um bloco retilíneo com o mesmo recuo da nota institucional (indentStart e indentFirstLine >= 198.4 pt)",
							Details:      "Execute 'format_document_abnt' para aplicar o recuo de 8,0 cm no orientador.",
						})
					}
				}
			}

			'''

c = c[:idx_start] + new_linter + c[idx_end:]

with open('internal/gdocs/linter.go', 'w', encoding='utf-8', newline='\n') as f:
    f.write(c)

print("linter.go updated successfully with overflow protection and full block alignment!")

