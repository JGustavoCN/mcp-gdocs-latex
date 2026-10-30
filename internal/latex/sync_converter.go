package latex

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	regexDoubleQuotes = regexp.MustCompile(`[“”]`)
	regexSingleQuotes = regexp.MustCompile(`[‘’]`)
	regexMarkdownBold = regexp.MustCompile(`\*\*(.+?)\*\*`)
	regexMarkdownItal = regexp.MustCompile(`\*(.+?)\*`)
	regexBracketCites = regexp.MustCompile(`\[@([a-zA-Z0-9_\-:]+)\]`)
)

// ConvertGDocsTextToLatex converte texto formatado do Google Docs / Markdown para sintaxe LaTeX compatível com ABNT
func ConvertGDocsTextToLatex(text string) string {
	if text == "" {
		return ""
	}

	result := text

	// 1. Normaliza quebras de linha Windows CRLF -> LF
	result = strings.ReplaceAll(result, "\r\n", "\n")
	result = strings.ReplaceAll(result, "\r", "\n")

	// 2. Converte aspas curvas para aspas tipográficas do TeX
	// Abre aspa dupla “ -> `` e Fecha aspa dupla ” -> ''
	result = regexp.MustCompile(`“([^”]+)”`).ReplaceAllString(result, "``$1''")
	result = regexp.MustCompile(`"([^"]+)"`).ReplaceAllString(result, "``$1''")
	// Aspas simples ‘...’ -> `...'
	result = regexp.MustCompile(`‘([^’]+)’`).ReplaceAllString(result, "`$1'")

	// 3. Converte citações do formato [@chave] ou @chave para \cite{chave}
	result = regexBracketCites.ReplaceAllString(result, `\cite{$1}`)

	// 4. Converte estilos Markdown para comandos LaTeX
	// Negrito **texto** -> \textbf{texto}
	result = regexMarkdownBold.ReplaceAllString(result, `\textbf{$1}`)
	// Itálico *texto* -> \textit{texto}
	result = regexMarkdownItal.ReplaceAllString(result, `\textit{$1}`)

	// 5. Escapa caracteres reservados do LaTeX que ocorrem em texto livre
	// Cuidado para não escapar barras invertidas que já formam comandos (\cite, \textbf, \textit, etc.)
	// Escapa % que não seja precedido de \
	result = regexp.MustCompile(`([^\\])%`).ReplaceAllString(result, `$1\%`)
	// Escapa & que não seja precedido de \
	result = regexp.MustCompile(`([^\\])&`).ReplaceAllString(result, `$1\&`)
	// Escapa _ que não seja precedido de \ e não esteja dentro de comandos de citação
	// Ex: sub_rotina -> sub\_rotina, preservando \cite{chave_1}
	lines := strings.Split(result, "\n")
	for i, line := range lines {
		if !strings.Contains(line, "\\") {
			lines[i] = strings.ReplaceAll(line, "_", `\_`)
		}
	}
	result = strings.Join(lines, "\n")

	return result
}

// UpdateLatexSectionInFile localiza uma seção ou subseção pelo título ou label em um arquivo .tex
// e substitui o conteúdo dela preservando comandos de cabeçalho e gerando backup automático .bak
func UpdateLatexSectionInFile(filePath, targetTitleOrLabel, newContent string) (string, error) {
	if !strings.HasSuffix(strings.ToLower(filePath), ".tex") {
		return "", fmt.Errorf("por segurança, apenas arquivos com extensão .tex podem ser editados")
	}

	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("erro ao ler arquivo %s: %w", filePath, err)
	}

	content := string(contentBytes)
	lines := strings.Split(content, "\n")

	normalizedTarget := strings.ToLower(strings.TrimSpace(targetTitleOrLabel))

	sectionStartIdx := -1
	sectionEndIdx := -1
	inTargetSection := false

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Verifica se a linha inicia uma seção/capítulo
		isHeader := regexSectionCmd.MatchString(trimmed) ||
			regexSubsectionCmd.MatchString(trimmed) ||
			regexSubsubsectionCmd.MatchString(trimmed) ||
			regexChapterCmd.MatchString(trimmed)

		if isHeader {
			var title string
			if m := regexChapterCmd.FindStringSubmatch(trimmed); len(m) > 1 {
				title = CleanLaTeXText(m[1])
			} else if m := regexSectionCmd.FindStringSubmatch(trimmed); len(m) > 1 {
				title = CleanLaTeXText(m[1])
			} else if m := regexSubsectionCmd.FindStringSubmatch(trimmed); len(m) > 1 {
				title = CleanLaTeXText(m[1])
			} else if m := regexSubsubsectionCmd.FindStringSubmatch(trimmed); len(m) > 1 {
				title = CleanLaTeXText(m[1])
			}

			// Se já estávamos na seção alvo, o próximo cabeçalho marca o fim dela
			if inTargetSection {
				sectionEndIdx = i
				break
			}

			// Verifica se coincide com o título procurado
			if strings.EqualFold(strings.TrimSpace(title), normalizedTarget) {
				inTargetSection = true
				sectionStartIdx = i
			}
		}

		// Também verifica se a linha possui o \label correspondente ao target procurado
		if !inTargetSection {
			if lblMatch := regexLineLabel.FindStringSubmatch(trimmed); len(lblMatch) > 1 {
				if strings.EqualFold(strings.TrimSpace(lblMatch[1]), normalizedTarget) {
					inTargetSection = true
					// O início do conteúdo é a partir da linha seguinte ao label
					sectionStartIdx = i
				}
			}
		}
	}

	if sectionStartIdx == -1 {
		return "", fmt.Errorf("seção ou label '%s' não foi encontrada no arquivo %s", targetTitleOrLabel, filePath)
	}

	// Se não encontrou outro cabeçalho depois, vai até o fim do arquivo
	if sectionEndIdx == -1 {
		sectionEndIdx = len(lines)
	}

	// Avança o start index se houver um \label na linha seguinte ao cabeçalho
	contentStartIdx := sectionStartIdx + 1
	if contentStartIdx < len(lines) && regexLineLabel.MatchString(strings.TrimSpace(lines[contentStartIdx])) {
		contentStartIdx++
	}

	// Converte o texto novo para formato LaTeX seguro
	convertedLatex := ConvertGDocsTextToLatex(newContent)

	// Monta o novo conteúdo do arquivo
	var newLines []string
	newLines = append(newLines, lines[:contentStartIdx]...)
	newLines = append(newLines, "")
	newLines = append(newLines, convertedLatex)
	newLines = append(newLines, "")
	if sectionEndIdx < len(lines) {
		newLines = append(newLines, lines[sectionEndIdx:]...)
	}

	finalContent := strings.Join(newLines, "\n")

	// Cria backup automático
	backupPath := filePath + ".bak"
	if err := os.WriteFile(backupPath, contentBytes, 0644); err != nil {
		return "", fmt.Errorf("erro ao gerar backup em %s: %w", backupPath, err)
	}

	// Salva o arquivo atualizado
	if err := os.WriteFile(filePath, []byte(finalContent), 0644); err != nil {
		return "", fmt.Errorf("erro ao escrever arquivo atualizado em %s: %w", filePath, err)
	}

	return fmt.Sprintf("SUCESSO: Seção '%s' no arquivo '%s' atualizada com sucesso!\nBackup salvo em: %s",
		targetTitleOrLabel, filepath.Base(filePath), backupPath), nil
}

