package gdocs

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf16"

	"google.golang.org/api/docs/v1"
)

// Occurrence representa uma ocorrência exata de um texto procurado.
type Occurrence struct {
	Index       int    `json:"index"`        // Posição global na lista de achados (0, 1, 2...)
	StartIndex  int64  `json:"start_index"`  // Índice exato no Google Docs
	EndIndex    int64  `json:"end_index"`    // Índice exato no Google Docs
	MatchedText string `json:"matched_text"` // Texto exato (pode ter variação de case)
	Context     string `json:"context"`      // Parágrafo inteiro onde foi achado
}

// searchOccurrences varre o documento buscando todas as ocorrências exatas de 'query'.
// Se isRegex for true, compila a query como expressão regular.
func searchOccurrences(ctx context.Context, documentID, query string, isRegex bool) ([]Occurrence, error) {
	svc, err := GetDocsService(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	doc, err := svc.Documents.Get(documentID).Do()
	if err != nil {
		return nil, fmt.Errorf("erro ao obter documento '%s': %w", documentID, err)
	}

	if doc.Body == nil {
		return nil, nil
	}

	var occurrences []Occurrence

	var re *regexp.Regexp
	queryLower := ""
	queryUTF16Len := int64(0)

	if isRegex {
		re, err = regexp.Compile(query)
		if err != nil {
			return nil, fmt.Errorf("erro ao compilar regex '%s': %w", query, err)
		}
	} else {
		queryLower = strings.ToLower(query)
		queryUTF16Len = int64(len(utf16.Encode([]rune(query))))
	}

	// Função recursiva para varrer elementos
	var searchElements func(elements []*docs.StructuralElement)
	searchElements = func(elements []*docs.StructuralElement) {
		for _, el := range elements {
			if el.Paragraph != nil {
				// Junta o texto do parágrafo inteiro para ignorar quebras de TextRun
				var paraTextBuilder strings.Builder
				var baseStartIndex int64 = -1

				for _, pEl := range el.Paragraph.Elements {
					if pEl.TextRun != nil {
						if baseStartIndex == -1 {
							baseStartIndex = pEl.StartIndex
						}
						paraTextBuilder.WriteString(pEl.TextRun.Content)
					}
				}

				paraText := paraTextBuilder.String()
				paraLower := strings.ToLower(paraText)

				// Busca a query dentro deste parágrafo
				if isRegex {
					// Encontra todos os matches no parágrafo usando a string original (pois regex considera case se não tiver flag)
					matches := re.FindAllStringIndex(paraText, -1)
					for _, m := range matches {
						matchStartStr := m[0]
						matchEndStr := m[1]

						prefixUTF16 := len(utf16.Encode([]rune(paraText[:matchStartStr])))
						matchUTF16Len := int64(len(utf16.Encode([]rune(paraText[matchStartStr:matchEndStr]))))

						matchStartIndex := baseStartIndex + int64(prefixUTF16)
						matchEndIndex := matchStartIndex + matchUTF16Len

						matchedStr := paraText[matchStartStr:matchEndStr]

						occ := Occurrence{
							Index:       len(occurrences),
							StartIndex:  matchStartIndex,
							EndIndex:    matchEndIndex,
							MatchedText: matchedStr,
							Context:     strings.TrimSpace(paraText),
						}
						occurrences = append(occurrences, occ)
					}
				} else {
					offsetStr := 0
					for {
						idx := strings.Index(paraLower[offsetStr:], queryLower)
						if idx == -1 {
							break
						}

						// Posição em UTF-8 na string do parágrafo
						matchStartStr := offsetStr + idx

						// Precisamos converter a substring anterior para UTF-16 para saber o offset real
						prefixUTF16 := len(utf16.Encode([]rune(paraText[:matchStartStr])))

						// Calcula Start e End absolutos no Google Docs
						matchStartIndex := baseStartIndex + int64(prefixUTF16)
						matchEndIndex := matchStartIndex + queryUTF16Len

						// Extrai exatamente o texto casado (para caso a query não considere case)
						matchedStr := paraText[matchStartStr : matchStartStr+len(query)]

						occ := Occurrence{
							Index:       len(occurrences),
							StartIndex:  matchStartIndex,
							EndIndex:    matchEndIndex,
							MatchedText: matchedStr,
							Context:     strings.TrimSpace(paraText),
						}
						occurrences = append(occurrences, occ)

						// Avança o offset
						offsetStr = matchStartStr + len(query)
					}
				}
			} else if el.Table != nil {
				for _, row := range el.Table.TableRows {
					for _, cell := range row.TableCells {
						searchElements(cell.Content)
					}
				}
			}
		}
	}

	searchElements(doc.Body.Content)

	return occurrences, nil
}

// SearchInDoc varre o documento buscando todas as ocorrências exatas de 'query'.
// Retorna uma lista em formato JSON amigável para o LLM.
func SearchInDoc(ctx context.Context, documentID, query string, isRegex bool) (string, error) {
	occurrences, err := searchOccurrences(ctx, documentID, query, isRegex)
	if err != nil {
		return "", err
	}

	if len(occurrences) == 0 {
		return fmt.Sprintf("Nenhuma ocorrência de '%s' foi encontrada.", query), nil
	}

	// Formata em JSON bonitinho
	b, err := json.MarshalIndent(occurrences, "", "  ")
	if err != nil {
		return "", fmt.Errorf("erro ao formatar JSON: %w", err)
	}

	return string(b), nil
}
