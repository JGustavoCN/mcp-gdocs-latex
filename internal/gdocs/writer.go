package gdocs

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"unicode/utf16"

	"google.golang.org/api/docs/v1"
)

// SearchAndReplaceText substitui ocorrências de oldText por newText globalmente no documento.
func SearchAndReplaceText(ctx context.Context, documentID, oldText, newText string) (string, error) {
	oldText = strings.ReplaceAll(oldText, "\r\n", "\n")
	newText = strings.ReplaceAll(newText, "\r\n", "\n")

	// Validação de colisão de comentários
	err := CheckCommentCollision(ctx, documentID, oldText)
	if err != nil {
		return "", err
	}

	svc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	req := &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{
			{
				ReplaceAllText: &docs.ReplaceAllTextRequest{
					ContainsText: &docs.SubstringMatchCriteria{
						Text:      oldText,
						MatchCase: true,
					},
					ReplaceText: newText,
				},
			},
		},
	}

	result, err := svc.Documents.BatchUpdate(documentID, req).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao executar batchUpdate (ReplaceAllText): %w", err)
	}

	return fmt.Sprintf("SUCESSO: %d ocorrência(s) substituída(s) globalmente.", result.Replies[0].ReplaceAllText.OccurrencesChanged), nil
}

// ApplyTextStyle aplica cirurgicamente o estilo (BOLD, ITALIC, NONE) em uma faixa de índices.
// Não remove nem modifica o texto em si, preservando os comentários fixados.
func ApplyTextStyle(ctx context.Context, documentID string, startIndex, endIndex int64, textStyle string) (string, error) {
	svc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	style := &docs.TextStyle{
		Bold:   false,
		Italic: false,
	}

	if textStyle == "BOLD" {
		style.Bold = true
	} else if textStyle == "ITALIC" {
		style.Italic = true
	}

	req := &docs.BatchUpdateDocumentRequest{
		Requests: []*docs.Request{
			{
				UpdateTextStyle: &docs.UpdateTextStyleRequest{
					Range: &docs.Range{
						StartIndex: startIndex,
						EndIndex:   endIndex,
					},
					TextStyle: style,
					Fields:    "bold,italic",
				},
			},
		},
	}

	_, err = svc.Documents.BatchUpdate(documentID, req).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao aplicar estilo (UpdateTextStyleRequest): %w", err)
	}

	return fmt.Sprintf("SUCESSO: Estilo '%s' aplicado com sucesso no range [%d - %d].", textStyle, startIndex, endIndex), nil
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

	for i := range chunks {
		chunks[i].ExpectedText = strings.ReplaceAll(chunks[i].ExpectedText, "\r\n", "\n")
		chunks[i].ReplacementText = strings.ReplaceAll(chunks[i].ReplacementText, "\r\n", "\n")
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

		err = CheckCommentCollision(ctx, documentID, chunk.ExpectedText)
		if err != nil {
			return "", err
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

// FormattedContent define o estilo do texto a ser inserido.
type FormattedContent struct {
	Text           string `json:"text"`
	ParagraphStyle string `json:"paragraph_style"`
	TextStyle      string `json:"text_style"`
	LinkUrl        string `json:"link_url,omitempty"` // Novo campo
}

// ReplaceRequest encapsula a requisição da nova ferramenta determinística de edição estrutural.
type ReplaceRequest struct {
	StartIndex   int64              `json:"start_index"`
	EndIndex     int64              `json:"end_index"`
	ExpectedText string             `json:"expected_text"`
	Content      []FormattedContent `json:"content"`
}

// UpdateBlockByIndex aplica as formatações baseadas em JSON na range determinística.
func UpdateBlockByIndex(ctx context.Context, documentID string, req ReplaceRequest) (string, error) {
	req.ExpectedText = strings.ReplaceAll(req.ExpectedText, "\r\n", "\n")
	for i := range req.Content {
		req.Content[i].Text = strings.ReplaceAll(req.Content[i].Text, "\r\n", "\n")
	}

	svc, err := GetDocsService(ctx)
	if err != nil {
		return "", fmt.Errorf("erro ao obter serviço do Docs: %w", err)
	}

	doc, err := svc.Documents.Get(documentID).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao obter documento '%s': %w", documentID, err)
	}

	extracted, err := extractTextRange(doc.Body.Content, req.StartIndex, req.EndIndex)
	if err != nil {
		return "", fmt.Errorf("falha ao extrair texto do range: %w", err)
	}

	if extracted != req.ExpectedText {
		return "", fmt.Errorf(
			"CONDIÇÃO DE CORRIDA / VALIDAÇÃO FALHOU:\n"+
				"Texto esperado: '%s'\n"+
				"Texto no Docs: '%s'",
			req.ExpectedText, extracted,
		)
	}

	// Trava de segurança
	err = CheckCommentCollision(ctx, documentID, req.ExpectedText)
	if err != nil {
		return "", err
	}

	var requests []*docs.Request

	// Deleta a range
	requests = append(requests, &docs.Request{
		DeleteContentRange: &docs.DeleteContentRangeRequest{
			Range: &docs.Range{
				StartIndex: req.StartIndex,
				EndIndex:   req.EndIndex,
			},
		},
	})

	// Index Shifting Tracking com regra de UTF-16
	currentIndex := req.StartIndex

	for _, bloco := range req.Content {
		requests = append(requests, &docs.Request{
			InsertText: &docs.InsertTextRequest{
				Location: &docs.Location{
					Index: currentIndex,
				},
				Text: bloco.Text,
			},
		})

		utf16Length := int64(len(utf16.Encode([]rune(bloco.Text))))

		style := &docs.TextStyle{
			Bold:   false,
			Italic: false,
		}

		if bloco.TextStyle == "BOLD" {
			style.Bold = true
		} else if bloco.TextStyle == "ITALIC" {
			style.Italic = true
		}

		if bloco.LinkUrl != "" {
			style.Link = &docs.Link{Url: bloco.LinkUrl}
		} else {
			style.Link = nil // Garante remoção de link em texto comum
		}

		requests = append(requests, &docs.Request{
			UpdateTextStyle: &docs.UpdateTextStyleRequest{
				Range: &docs.Range{
					StartIndex: currentIndex,
					EndIndex:   currentIndex + utf16Length,
				},
				TextStyle: style,
				Fields:    "bold,italic,link",
			},
		})

		if bloco.ParagraphStyle != "NORMAL_TEXT" && bloco.ParagraphStyle != "" {
			requests = append(requests, &docs.Request{
				UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
					Range: &docs.Range{
						StartIndex: currentIndex,
						EndIndex:   currentIndex + utf16Length,
					},
					ParagraphStyle: &docs.ParagraphStyle{
						NamedStyleType: bloco.ParagraphStyle,
					},
					Fields: "namedStyleType",
				},
			})
		}

		currentIndex += utf16Length
	}

	updateReq := &docs.BatchUpdateDocumentRequest{Requests: requests}
	_, err = svc.Documents.BatchUpdate(documentID, updateReq).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao executar batchUpdate (UpdateBlockByIndex): %w", err)
	}

	return "SUCESSO: Bloco de texto atualizado e formatado via JSON com Index Shifting.", nil
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
