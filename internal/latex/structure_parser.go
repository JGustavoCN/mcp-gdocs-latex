package latex

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// StructuralElementType define a tipagem de elementos no LaTeX
type StructuralElementType string

const (
	ElementChapter       StructuralElementType = "CHAPTER"
	ElementSection       StructuralElementType = "SECTION"
	ElementSubsection    StructuralElementType = "SUBSECTION"
	ElementSubsubsection StructuralElementType = "SUBSUBSECTION"
	ElementTable         StructuralElementType = "TABLE"
	ElementFigure        StructuralElementType = "FIGURE"
	ElementListing       StructuralElementType = "LISTING"
)

// LatexElement representa um nó estrutural localizado em um arquivo .tex
type LatexElement struct {
	Type      StructuralElementType `json:"type"`
	Title     string                `json:"title"`
	Label     string                `json:"label,omitempty"`
	File      string                `json:"file"`
	StartLine int                   `json:"start_line"`
	EndLine   int                   `json:"end_line,omitempty"`
}

var (
	regexChapterCmd       = regexp.MustCompile(`^\\chapter\*?\{([^}]+)\}`)
	regexSectionCmd       = regexp.MustCompile(`^\\section\*?\{([^}]+)\}`)
	regexSubsectionCmd    = regexp.MustCompile(`^\\subsection\*?\{([^}]+)\}`)
	regexSubsubsectionCmd = regexp.MustCompile(`^\\subsubsection\*?\{([^}]+)\}`)
	regexLineLabel        = regexp.MustCompile(`\\label\{([^}]+)\}`)
	regexBeginTable       = regexp.MustCompile(`\\begin\{(?:table|quadro)\}`)
	regexEndTable         = regexp.MustCompile(`\\end\{(?:table|quadro)\}`)
	regexBeginFigure      = regexp.MustCompile(`\\begin\{figure\}`)
	regexEndFigure        = regexp.MustCompile(`\\end\{figure\}`)
	regexBeginListing     = regexp.MustCompile(`\\begin\{lstlisting\}`)
	regexEndListing       = regexp.MustCompile(`\\end\{lstlisting\}`)
)

// ParseLatexFileStructure analisa linha a linha um arquivo .tex e mapeia sua estrutura
func ParseLatexFileStructure(filePath string) ([]LatexElement, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo %s: %w", filePath, err)
	}
	defer file.Close()

	var elements []LatexElement
	scanner := bufio.NewScanner(file)
	lineNum := 0

	var currentContainer *LatexElement

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Ignora comentários puros de linha inteira
		if strings.HasPrefix(line, "%") {
			continue
		}

		// Rastreia início de ambientes flutuantes (Tabelas, Figuras, Códigos)
		if regexBeginTable.MatchString(line) {
			currentContainer = &LatexElement{
				Type:      ElementTable,
				File:      filepath.Base(filePath),
				StartLine: lineNum,
			}
			continue
		}
		if regexBeginFigure.MatchString(line) {
			currentContainer = &LatexElement{
				Type:      ElementFigure,
				File:      filepath.Base(filePath),
				StartLine: lineNum,
			}
			continue
		}
		if regexBeginListing.MatchString(line) {
			currentContainer = &LatexElement{
				Type:      ElementListing,
				File:      filepath.Base(filePath),
				StartLine: lineNum,
			}
			continue
		}

		// Se estamos dentro de um ambiente flutuante, busca caption e label
		if currentContainer != nil {
			if capMatch := regexCaption.FindStringSubmatch(line); len(capMatch) > 1 {
				currentContainer.Title = CleanLaTeXText(capMatch[1])
			}
			if lblMatch := regexLineLabel.FindStringSubmatch(line); len(lblMatch) > 1 {
				currentContainer.Label = strings.TrimSpace(lblMatch[1])
			}

			if regexEndTable.MatchString(line) || regexEndFigure.MatchString(line) || regexEndListing.MatchString(line) {
				currentContainer.EndLine = lineNum
				elements = append(elements, *currentContainer)
				currentContainer = nil
			}
			continue
		}

		// Títulos hierárquicos
		var el *LatexElement
		if m := regexChapterCmd.FindStringSubmatch(line); len(m) > 1 {
			el = &LatexElement{Type: ElementChapter, Title: CleanLaTeXText(m[1]), File: filepath.Base(filePath), StartLine: lineNum}
		} else if m := regexSectionCmd.FindStringSubmatch(line); len(m) > 1 {
			el = &LatexElement{Type: ElementSection, Title: CleanLaTeXText(m[1]), File: filepath.Base(filePath), StartLine: lineNum}
		} else if m := regexSubsectionCmd.FindStringSubmatch(line); len(m) > 1 {
			el = &LatexElement{Type: ElementSubsection, Title: CleanLaTeXText(m[1]), File: filepath.Base(filePath), StartLine: lineNum}
		} else if m := regexSubsubsectionCmd.FindStringSubmatch(line); len(m) > 1 {
			el = &LatexElement{Type: ElementSubsubsection, Title: CleanLaTeXText(m[1]), File: filepath.Base(filePath), StartLine: lineNum}
		}

		if el != nil {
			// Tenta capturar label na mesma linha se houver
			if lblMatch := regexLineLabel.FindStringSubmatch(line); len(lblMatch) > 1 {
				el.Label = strings.TrimSpace(lblMatch[1])
			}
			elements = append(elements, *el)
		} else if len(elements) > 0 && elements[len(elements)-1].Label == "" {
			// Se a linha anterior foi um título e esta linha contém apenas o \label
			if lblMatch := regexLineLabel.FindStringSubmatch(line); len(lblMatch) > 1 {
				elements[len(elements)-1].Label = strings.TrimSpace(lblMatch[1])
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a leitura do arquivo %s: %w", filePath, err)
	}

	return elements, nil
}

// ScanTCCStructure varre a pasta resources/tcc ou subpastas e retorna o sumário estrutural de todos os capítulos
func ScanTCCStructure(tccBaseDir string) ([]LatexElement, error) {
	capitulosDir := filepath.Join(tccBaseDir, "capitulos")
	if _, err := os.Stat(capitulosDir); err != nil {
		capitulosDir = tccBaseDir
	}

	pattern := filepath.Join(capitulosDir, "*.tex")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar capítulos em %s: %w", pattern, err)
	}

	var allElements []LatexElement
	for _, f := range files {
		elems, err := ParseLatexFileStructure(f)
		if err != nil {
			return nil, err
		}
		allElements = append(allElements, elems...)
	}

	return allElements, nil
}

// KnownSectionInfo detalha uma seção mapeada do TCC com numeração ABNT e nível hierárquico
type KnownSectionInfo struct {
	Level         int                   `json:"level"`          // 1 (Chapter), 2 (Section), 3 (Subsection), 4 (Subsubsection)
	TargetHeading string                `json:"target_heading"` // HEADING_1, HEADING_2, HEADING_3, HEADING_4
	Number        string                `json:"number"`         // "1", "1.1", "1.4.1"
	Title         string                `json:"title"`          // "Contextualização"
	FullABNTTitle string                `json:"full_abnt_title"`// "1.1 CONTEXTUALIZAÇÃO" ou "1.4.1 Objetivo Geral"
	Type          StructuralElementType `json:"type"`
	File          string                `json:"file"`
}

// GetNumberedTCCStructure calcula a numeração progressiva ABNT (NBR 6024) para todas as seções mapeadas do TCC
func GetNumberedTCCStructure(tccBaseDir string) ([]KnownSectionInfo, error) {
	elements, err := ScanTCCStructure(tccBaseDir)
	if err != nil {
		return nil, err
	}

	var result []KnownSectionInfo
	chapNum := 0
	secNum := 0
	subsecNum := 0
	subsubsecNum := 0

	for _, el := range elements {
		switch el.Type {
		case ElementChapter:
			chapNum++
			secNum = 0
			subsecNum = 0
			subsubsecNum = 0
			numStr := fmt.Sprintf("%d", chapNum)
			result = append(result, KnownSectionInfo{
				Level:         1,
				TargetHeading: "HEADING_1",
				Number:        numStr,
				Title:         el.Title,
				FullABNTTitle: fmt.Sprintf("%s %s", numStr, strings.ToUpper(el.Title)),
				Type:          el.Type,
				File:          el.File,
			})
		case ElementSection:
			secNum++
			subsecNum = 0
			subsubsecNum = 0
			numStr := fmt.Sprintf("%d.%d", chapNum, secNum)
			result = append(result, KnownSectionInfo{
				Level:         2,
				TargetHeading: "HEADING_2",
				Number:        numStr,
				Title:         el.Title,
				FullABNTTitle: fmt.Sprintf("%s %s", numStr, strings.ToUpper(el.Title)),
				Type:          el.Type,
				File:          el.File,
			})
		case ElementSubsection:
			subsecNum++
			subsubsecNum = 0
			numStr := fmt.Sprintf("%d.%d.%d", chapNum, secNum, subsecNum)
			result = append(result, KnownSectionInfo{
				Level:         3,
				TargetHeading: "HEADING_3",
				Number:        numStr,
				Title:         el.Title,
				FullABNTTitle: fmt.Sprintf("%s %s", numStr, el.Title),
				Type:          el.Type,
				File:          el.File,
			})
		case ElementSubsubsection:
			subsubsecNum++
			numStr := fmt.Sprintf("%d.%d.%d.%d", chapNum, secNum, subsecNum, subsubsecNum)
			result = append(result, KnownSectionInfo{
				Level:         4,
				TargetHeading: "HEADING_4",
				Number:        numStr,
				Title:         el.Title,
				FullABNTTitle: fmt.Sprintf("%s %s", numStr, el.Title),
				Type:          el.Type,
				File:          el.File,
			})
		}
	}

	return result, nil
}
