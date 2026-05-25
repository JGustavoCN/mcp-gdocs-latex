package gdocs

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"unicode/utf16"

	"google.golang.org/api/docs/v1"
)

// ReplaceTextInDoc substitui ocorrências de oldText por newText usando índices precisos.
// occurrenceIndex define qual ocorrência substituir (0, 1, 2...). Se for -1, substitui TODAS.
func ReplaceTextInDoc(ctx context.Context, documentID, oldText, newText string, occurrenceIndex int) (string, error) {
	// Busca as ocorrências (busca exata)
	matches, err := searchOccurrences(ctx, documentID, oldText, false)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("texto original '%s' não encontrado no documento", oldText)
	}

	// Filtra as ocorrências desejadas
	var targetMatches []Occurrence
	if occurrenceIndex == -1 {
		targetMatches = matches
	} else {
		if occurrenceIndex < 0 || occurrenceIndex >= len(matches) {
			return "", fmt.Errorf(
				"índice de ocorrência inválido: %d. O documento possui apenas %d ocorrência(s) desse texto (índices válidos: 0 a %d)",
				occurrenceIndex, len(matches), len(matches)-1,
			)
		}
		targetMatches = append(targetMatches, matches[occurrenceIndex])
	}

	// REGRA DE OURO (OFFSET SHIFT):
	// Ordena as substituições em ordem DECRESCENTE pelo StartIndex
	sort.Slice(targetMatches, func(i, j int) bool {
		return targetMatches[i].StartIndex > targetMatches[j].StartIndex
	})

	svc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	var requests []*docs.Request
	for _, m := range targetMatches {
		// Primeiro apaga
		requests = append(requests, &docs.Request{
			DeleteContentRange: &docs.DeleteContentRangeRequest{
				Range: &docs.Range{
					StartIndex: m.StartIndex,
					EndIndex:   m.EndIndex,
				},
			},
		})
		// Depois insere
		requests = append(requests, &docs.Request{
			InsertText: &docs.InsertTextRequest{
				Location: &docs.Location{
					Index: m.StartIndex,
				},
				Text: newText,
			},
		})
	}

	req := &docs.BatchUpdateDocumentRequest{Requests: requests}
	_, err = svc.Documents.BatchUpdate(documentID, req).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao executar batchUpdate: %w", err)
	}

	return fmt.Sprintf("SUCESSO: %d ocorrência(s) de '%s' substituída(s) com sucesso.", len(targetMatches), oldText), nil
}

// Chunk define um bloco absoluto para substituição no Docs
type Chunk struct {
	StartIndex      int64  `json:"start_index"`
	EndIndex        int64  `json:"end_index"`
	ExpectedText    string `json:"expected_text"`
	ReplacementText string `json:"replacement_text"`
}

// MultiReplaceDocBlock faz substituições absolutas por índices.
// Aborta se ExpectedText não bater perfeitamente com o texto no Docs (validação atômica simulada).
func MultiReplaceDocBlock(ctx context.Context, documentID string, chunks []Chunk) (string, error) {
	if len(chunks) == 0 {
		return "Nenhum bloco fornecido.", nil
	}

	svc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	doc, err := svc.Documents.Get(documentID).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao obter documento '%s': %w", documentID, err)
	}

	// 1. Validação de ExpectedText
	// O Docs retorna o documento como uma árvore.
	// Precisamos extrair o texto exato para validar os chunks.
	for i, chunk := range chunks {
		extracted, err := extractTextRange(doc.Body.Content, chunk.StartIndex, chunk.EndIndex)
		if err != nil {
			return "", fmt.Errorf("chunk [%d]: falha ao extrair texto - %w", i, err)
		}

		// Remove possíveis quebras de linha estranhas que o Google Docs embute na estrutura se necessário,
		// mas como queremos validação estrita, compararemos cru.
		if extracted != chunk.ExpectedText {
			return "", fmt.Errorf(
				"CONDIÇÃO DE CORRIDA / VALIDAÇÃO FALHOU no Chunk [%d]:\n"+
					"Texto esperado: '%s'\n"+
					"Texto no Google Docs: '%s'\n"+
					"A operação foi abortada para proteger a integridade do documento.",
				i, chunk.ExpectedText, extracted,
			)
		}
	}

	// 2. Ordenação Decrescente (Offset Shift Rule)
	sort.Slice(chunks, func(i, j int) bool {
		return chunks[i].StartIndex > chunks[j].StartIndex
	})

	var requests []*docs.Request
	for _, ch := range chunks {
		requests = append(requests, &docs.Request{
			DeleteContentRange: &docs.DeleteContentRangeRequest{
				Range: &docs.Range{
					StartIndex: ch.StartIndex,
					EndIndex:   ch.EndIndex,
				},
			},
		})
		requests = append(requests, &docs.Request{
			InsertText: &docs.InsertTextRequest{
				Location: &docs.Location{
					Index: ch.StartIndex,
				},
				Text: ch.ReplacementText,
			},
		})
	}

	req := &docs.BatchUpdateDocumentRequest{Requests: requests}
	_, err = svc.Documents.BatchUpdate(documentID, req).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao executar batchUpdate atômico: %w", err)
	}

	return fmt.Sprintf("SUCESSO: %d blocos processados cirurgicamente.", len(chunks)), nil
}

// extractTextRange extrai o texto do Google Docs para uma faixa específica convertendo para UTF-16
func extractTextRange(elements []*docs.StructuralElement, start, end int64) (string, error) {
	var sb strings.Builder
	var traverse func(els []*docs.StructuralElement)

	traverse = func(els []*docs.StructuralElement) {
		for _, el := range els {
			if el.Paragraph != nil {
				for _, pEl := range el.Paragraph.Elements {
					if pEl.TextRun != nil {
						// Verifica interseção
						if start < pEl.EndIndex && end > pEl.StartIndex {
							// Codifica o conteúdo em UTF-16
							u16 := utf16.Encode([]rune(pEl.TextRun.Content))

							// Calcula limites locais
							localStart := int64(0)
							if start > pEl.StartIndex {
								localStart = start - pEl.StartIndex
							}

							localEnd := int64(len(u16))
							if end < pEl.EndIndex {
								localEnd = end - pEl.StartIndex
							}

							if localStart >= 0 && localEnd <= int64(len(u16)) && localStart <= localEnd {
								sb.WriteString(string(utf16.Decode(u16[localStart:localEnd])))
							}
						}
					}
				}
			} else if el.Table != nil {
				for _, row := range el.Table.TableRows {
					for _, cell := range row.TableCells {
						traverse(cell.Content)
					}
				}
			}
		}
	}

	traverse(elements)
	return sb.String(), nil
}
