with open('internal/gdocs/cover_generator.go', 'r', encoding='utf-8') as f:
    c = f.read().replace('\r\n', '\n')

# 1. Ajuste em GenerateABNTIFSCoverAndTitlePage
old_gen = '''	// Instituição
	instLines := meta.InstituicaoLines()
	for _, l := range instLines {
		sb.WriteString(strings.ToUpper(l))
		sb.WriteString("\\n")
	}
	sb.WriteString(strings.Repeat("\\n", 6)) // Exatamente 6 quebras de linhas entre Instituição e Autor (NBR 14724)

	// Autor
	sb.WriteString(strings.ToUpper(meta.Autor))
	sb.WriteString(strings.Repeat("\\n", 6)) // Exatamente 6 quebras de linhas entre Autor e Título (NBR 14724)

	// Título do Trabalho
	sb.WriteString(strings.ToUpper(meta.Titulo))
	sb.WriteString(strings.Repeat("\\n", 22)) // Exatamente 22 quebras vazias posicionam Local e Ano no limite inferior da Capa (rodapé)

	// Local e Ano
	sb.WriteString(strings.ToUpper(meta.Local))
	sb.WriteString("\\n")
	sb.WriteString(meta.Data)

	capaText := sb.String()

	// --- PÁGINA 2: FOLHA DE ROSTO (NBR 14724 + IFS config.tex) ---
	var sbFR strings.Builder
	sbFR.WriteString(strings.Repeat("\\n", 5)) // Exatamente 5 quebras simples no topo da Folha de Rosto (NBR 14724)
	sbFR.WriteString(strings.ToUpper(meta.Autor))
	sbFR.WriteString(strings.Repeat("\\n", 8)) // Exatamente 8 quebras de linhas entre Autor e Título
	sbFR.WriteString(strings.ToUpper(meta.Titulo))
	sbFR.WriteString(strings.Repeat("\\n", 7)) // Exatamente 7 quebras de linhas entre Título e Preâmbulo

	// Bloco da Nota Institucional + Orientador (alinhado da metade para a margem direita, recuo 8.0 cm)
	sbFR.WriteString(meta.Preambulo)
	sbFR.WriteString("\\n")
	sbFR.WriteString("Orientador: ")
	sbFR.WriteString(meta.Orientador)
	if meta.Coorientador != "" {
		sbFR.WriteString("\\nCoorientador: ")
		sbFR.WriteString(meta.Coorientador)
	}
	sbFR.WriteString(strings.Repeat("\\n", 19)) // Exatamente 19 quebras vazias posicionam Local e Ano no limite inferior da Folha de Rosto'''

new_gen = '''	// Instituição
	instLines := meta.InstituicaoLines()
	for _, l := range instLines {
		sb.WriteString(strings.ToUpper(l))
		sb.WriteString("\\n")
	}
	sb.WriteString(strings.Repeat("\\n", 3)) // 3 quebras verticais calibradas (NBR 14724 / sem overflow)

	// Autor
	sb.WriteString(strings.ToUpper(meta.Autor))
	sb.WriteString(strings.Repeat("\\n", 3)) // 3 quebras verticais calibradas (NBR 14724 / sem overflow)

	// Título do Trabalho
	sb.WriteString(strings.ToUpper(meta.Titulo))
	sb.WriteString(strings.Repeat("\\n", 15)) // 15 quebras verticais posicionam Local e Ano no rodapé da Capa (Página 1 estrita)

	// Local e Ano
	sb.WriteString(strings.ToUpper(meta.Local))
	sb.WriteString("\\n")
	sb.WriteString(meta.Data)

	capaText := sb.String()

	// --- PÁGINA 2: FOLHA DE ROSTO (NBR 14724 + IFS config.tex) ---
	var sbFR strings.Builder
	sbFR.WriteString(strings.Repeat("\\n", 2)) // 2 quebras verticais no topo (Página 2 estrita)
	sbFR.WriteString(strings.ToUpper(meta.Autor))
	sbFR.WriteString(strings.Repeat("\\n", 4)) // 4 quebras entre Autor e Título
	sbFR.WriteString(strings.ToUpper(meta.Titulo))
	sbFR.WriteString(strings.Repeat("\\n", 3)) // 3 quebras entre Título e Preâmbulo

	// Bloco da Nota Institucional + Orientador (alinhado da metade para a margem direita, recuo 8.0 cm)
	sbFR.WriteString(meta.Preambulo)
	sbFR.WriteString("\\n")
	sbFR.WriteString("Orientador: ")
	sbFR.WriteString(meta.Orientador)
	if meta.Coorientador != "" {
		sbFR.WriteString("\\nCoorientador: ")
		sbFR.WriteString(meta.Coorientador)
	}
	sbFR.WriteString(strings.Repeat("\\n", 12)) // 12 quebras verticais posicionam Local e Ano no rodapé da Folha de Rosto (Página 2 estrita)'''

assert old_gen in c, "old_gen not found in cover_generator.go"
c = c.replace(old_gen, new_gen)

# 2. Ajuste em EnsurePreTextualPageBreaksAndLayout (metas de quebras)
c = c.replace(
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, tituloCapa.EndIndex, cidadeCapa.StartIndex, 22)',
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, tituloCapa.EndIndex, cidadeCapa.StartIndex, 15)'
)
c = c.replace(
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, autorCapa.EndIndex, tituloCapa.StartIndex, 6)',
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, autorCapa.EndIndex, tituloCapa.StartIndex, 3)'
)
c = c.replace(
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, instCapa.EndIndex, autorCapa.StartIndex, 6)',
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, instCapa.EndIndex, autorCapa.StartIndex, 3)'
)

c = c.replace(
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, orientadorFR.EndIndex, cidadeFR.StartIndex, 19)',
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, orientadorFR.EndIndex, cidadeFR.StartIndex, 12)'
)
c = c.replace(
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, tituloFR.EndIndex, preambuloFR.StartIndex, 7)',
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, tituloFR.EndIndex, preambuloFR.StartIndex, 3)'
)
c = c.replace(
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, autorFR.EndIndex, tituloFR.StartIndex, 8)',
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, autorFR.EndIndex, tituloFR.StartIndex, 4)'
)
c = c.replace(
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, ano1EndIndex, autorFR.StartIndex, 5)',
    '_ = calibrateEmptyParagraphsBetween(docsSvc, documentID, ano1EndIndex, autorFR.StartIndex, 2)'
)

# 3. Ajuste crucial no recuo de primeira linha da Nota e Orientador em FormatABNTCoverAndTitlePageIsolated:
# IndentFirstLine DEVE SER IGUAL A IndentStart (226.772 pt), para que a 1ª linha NÃO fique na margem esquerda!
old_indent = '''IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}}, // alinhado ao bloco de 8.0 cm (deslocamento zero relativo a indentStart)
						IndentStart:         &docs.Dimension{Magnitude: ABNTIndentPreTextualNotePt, Unit: "PT"},                                 // 8.0 cm (226.772 pt)'''

new_indent = '''IndentFirstLine:     &docs.Dimension{Magnitude: ABNTIndentPreTextualNotePt, Unit: "PT"}, // IGUAL a IndentStart (8.0 cm), formando bloco retilíneo sem recuo na 1ª linha
						IndentStart:         &docs.Dimension{Magnitude: ABNTIndentPreTextualNotePt, Unit: "PT"}, // 8.0 cm (226.772 pt)'''

assert old_indent in c, "old_indent not found in cover_generator.go"
c = c.replace(old_indent, new_indent)

# Também ajusta para parágrafos vazios dentro do bloco
old_empty_note = '''IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT", ForceSendFields: []string{"Magnitude"}},
							IndentStart:         &docs.Dimension{Magnitude: ABNTIndentPreTextualNotePt, Unit: "PT"}, // 8.0 cm'''
new_empty_note = '''IndentFirstLine:     &docs.Dimension{Magnitude: ABNTIndentPreTextualNotePt, Unit: "PT"},
							IndentStart:         &docs.Dimension{Magnitude: ABNTIndentPreTextualNotePt, Unit: "PT"}, // 8.0 cm'''

if old_empty_note in c:
    c = c.replace(old_empty_note, new_empty_note)

with open('internal/gdocs/cover_generator.go', 'w', encoding='utf-8', newline='\n') as f:
    f.write(c)

print("cover_generator.go updated with calibrated heights and aligned indentFirstLine!")

