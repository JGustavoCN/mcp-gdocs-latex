package latex

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// stripLatexComments remove linhas de comentários '%' e comentários inline do LaTeX
func stripLatexComments(input string) string {
	var sb strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "%") {
			continue
		}
		inComment := false
		var cleanLine strings.Builder
		for i := 0; i < len(line); i++ {
			if line[i] == '%' && (i == 0 || line[i-1] != '\\') {
				inComment = true
				break
			}
			cleanLine.WriteByte(line[i])
		}
		if !inComment {
			sb.WriteString(line)
		} else {
			sb.WriteString(cleanLine.String())
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

// TCCMetadata consolida os metadados acadêmicos extraídos de metadados.tex
type TCCMetadata struct {
	Titulo       string `json:"titulo"`
	Autor        string `json:"autor"`
	Local        string `json:"local"`
	Data         string `json:"data"`
	Instituicao  string `json:"instituicao"`
	TipoTrabalho string `json:"tipo_trabalho"`
	Preambulo    string `json:"preambulo"`
	Orientador   string `json:"orientador"`
	Coorientador string `json:"coorientador,omitempty"`
}

// InstituicaoLines retorna as linhas da instituição divididas
func (m *TCCMetadata) InstituicaoLines() []string {
	var lines []string
	raw := strings.ReplaceAll(m.Instituicao, "\r\n", "\n")
	for _, l := range strings.Split(raw, "\n") {
		trimmed := strings.TrimSpace(l)
		if trimmed != "" {
			lines = append(lines, trimmed)
		}
	}
	if len(lines) == 0 {
		return []string{
			"INSTITUTO FEDERAL DE EDUCAÇÃO, CIÊNCIA E TECNOLOGIA DE SERGIPE",
			"CAMPUS LAGARTO",
			"BACHARELADO EM SISTEMAS DE INFORMAÇÃO",
		}
	}
	return lines
}

// ParseTCCMetadata lê e extrai os campos de configuração acadêmica de metadados.tex
func ParseTCCMetadata(filePath string) (*TCCMetadata, error) {
	meta := &TCCMetadata{
		Titulo:       "Avaliação de uma Arquitetura Híbrida RAG-MCP na Validação Normativa de Requerimentos Acadêmicos",
		Autor:        "José Gustavo Correia Nascimento",
		Local:        "Lagarto -- SE",
		Data:         "2026",
		Instituicao:  "Instituto Federal de Educação, Ciência e Tecnologia de Sergipe\nCampus Lagarto\nBacharelado em Sistemas de Informação",
		TipoTrabalho: "Trabalho de Conclusão de Curso",
		Preambulo:    "Trabalho de Conclusão de Curso apresentado ao curso de Bacharelado em Sistemas de Informação do Instituto Federal de Educação, Ciência e Tecnologia de Sergipe, Campus Lagarto, como requisito parcial para obtenção do grau de Bacharel em Sistemas de Informação.",
		Orientador:   "Prof. Diego Armando de Oliveira Meneses",
	}

	if filePath == "" {
		filePath = "resources/tcc/metadados.tex"
	}

	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		// Retorna os dados padrão se o arquivo não existir
		return meta, nil
	}
	content := string(contentBytes)
	contentClean := stripLatexComments(content)

	extractField := func(cmd string) string {
		pattern := `(?s)\\` + cmd + `\s*\{`
		re := regexp.MustCompile(pattern)
		loc := re.FindStringIndex(contentClean)
		if loc == nil {
			return ""
		}
		// Inicia logo após a primeira chave '{'
		startIdx := loc[1]
		depth := 1
		endIdx := -1
		for i := startIdx; i < len(contentClean); i++ {
			if contentClean[i] == '{' && (i == 0 || contentClean[i-1] != '\\') {
				depth++
			} else if contentClean[i] == '}' && (i == 0 || contentClean[i-1] != '\\') {
				depth--
				if depth == 0 {
					endIdx = i
					break
				}
			}
		}
		if endIdx == -1 {
			return ""
		}
		val := contentClean[startIdx:endIdx]
		val = strings.ReplaceAll(val, `\par`, "\n")
		val = strings.ReplaceAll(val, `\textit{`, "")
		val = strings.ReplaceAll(val, `\textbf{`, "")
		val = strings.ReplaceAll(val, `\MakeUppercase{`, "")
		val = strings.ReplaceAll(val, `}`, "")
		val = strings.ReplaceAll(val, `\\`, "\n")
		lines := strings.Split(val, "\n")
		var cleanLines []string
		for _, l := range lines {
			t := strings.TrimSpace(l)
			if t != "" {
				cleanLines = append(cleanLines, t)
			}
		}
		return strings.Join(cleanLines, "\n")
	}

	if v := extractField("titulo"); v != "" {
		meta.Titulo = v
	}
	if v := extractField("autor"); v != "" {
		meta.Autor = v
	}
	if v := extractField("local"); v != "" {
		meta.Local = v
	}
	if v := extractField("data"); v != "" {
		meta.Data = v
	}
	if v := extractField("instituicao"); v != "" {
		meta.Instituicao = v
	}
	if v := extractField("tipotrabalho"); v != "" {
		meta.TipoTrabalho = v
	}
	if v := extractField("preambulo"); v != "" {
		meta.Preambulo = v
	}
	if v := extractField("orientador"); v != "" {
		meta.Orientador = v
	}
	if v := extractField("coorientador"); v != "" {
		// Ignora placeholders genéricos do template caso não tenham sido preenchidos
		low := strings.ToLower(v)
		if !strings.Contains(low, "nome do coorientador") && !strings.Contains(low, "se houver") {
			meta.Coorientador = v
		} else {
			meta.Coorientador = ""
		}
	} else {
		meta.Coorientador = ""
	}

	if meta.Titulo == "" {
		return nil, fmt.Errorf("metadados inválidos: título não localizado em %s", filePath)
	}

	return meta, nil
}

