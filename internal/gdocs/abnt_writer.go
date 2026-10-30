package gdocs

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/api/docs/v1"
)

// getBodyEndIndex retorna o índice final do último elemento estrutural do corpo do documento
func getBodyEndIndex(doc *docs.Document) int64 {
	if doc != nil && doc.Body != nil && len(doc.Body.Content) > 0 {
		return doc.Body.Content[len(doc.Body.Content)-1].EndIndex
	}
	return 1
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

// FormatDocumentABNT formata parágrafos e página segundo as normas ABNT (NBR 14724)
// de forma inteligente e unificada: preserva títulos, legendas, fontes e tabelas nativas,
// e aplica margens oficiais A4 (3/3/2/2 cm) se applyMargins == true.
func FormatDocumentABNT(ctx context.Context, documentID string, startIndex, endIndex int64, applyMargins bool) (string, error) {
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

	// 1. Configuração unificada de página: Margens ABNT 3/3/2/2 cm e Folha A4
	if applyMargins {
		requests = append(requests, BuildABNTDocumentMarginsRequest())
	}

	formattedParas := 0
	preservedHeadings := 0
	preservedCaptions := 0

	// 2. Iteração estrutural elemento por elemento (Preservação de títulos, legendas e tabelas)
	for _, el := range doc.Body.Content {
		// Elemento fora do range solicitado
		if el.EndIndex <= start || el.StartIndex >= end {
			continue
		}

		// Se for tabela nativa: PRESERVA INTACTA (não aplica recuo de parágrafo nas células)
		if el.Table != nil {
			continue
		}

		if el.Paragraph != nil {
			p := el.Paragraph
			styleName := ""
			if p.ParagraphStyle != nil {
				styleName = p.ParagraphStyle.NamedStyleType
			}

			// Coleta texto do parágrafo
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

			// Delimita o range seguro deste parágrafo
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

			isHeading := strings.HasPrefix(styleName, "HEADING_") || styleName == "TITLE"
			isCaption := strings.HasPrefix(trimmed, "Figura ") ||
				strings.HasPrefix(trimmed, "Tabela ") ||
				strings.HasPrefix(trimmed, "Quadro ") ||
				strings.HasPrefix(trimmed, "Gráfico ")
			isSource := strings.HasPrefix(trimmed, "Fonte:") || strings.HasPrefix(trimmed, "Nota:")

			if isHeading {
				// TÍTULOS ABNT: Alinhamento à esquerda (START), sem recuo de 1ª linha, entrelinhas 1.5, keepWithNext
				preservedHeadings++
				requests = append(requests, &docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: pRange,
						ParagraphStyle: &docs.ParagraphStyle{
							Alignment:           "START",
							LineSpacing:         150,
							IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
							IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceAbove:          &docs.Dimension{Magnitude: 12, Unit: "PT"},
							SpaceBelow:          &docs.Dimension{Magnitude: 6, Unit: "PT"},
							KeepWithNext:        true,
							AvoidWidowAndOrphan: true,
						},
						Fields: "alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,keepWithNext,avoidWidowAndOrphan",
					},
				})
			} else if isCaption || isSource {
				// LEGENDAS E FONTES ABNT: Centralizadas ou à esquerda, sem recuo de 1ª linha, entrelinhas simples, fonte 10pt
				preservedCaptions++
				keepWithNext := isCaption // Legenda sempre acompanha a figura/tabela abaixo

				requests = append(requests, &docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: pRange,
						ParagraphStyle: &docs.ParagraphStyle{
							Alignment:           "CENTER",
							LineSpacing:         100, // Espaço simples
							IndentFirstLine:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
							IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceAbove:          &docs.Dimension{Magnitude: 4, Unit: "PT"},
							SpaceBelow:          &docs.Dimension{Magnitude: 4, Unit: "PT"},
							KeepWithNext:        keepWithNext,
							AvoidWidowAndOrphan: true,
						},
						Fields: "alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,keepWithNext,avoidWidowAndOrphan",
					},
				})

				// Aplica fonte 10pt na legenda/fonte
				requests = append(requests, &docs.Request{
					UpdateTextStyle: &docs.UpdateTextStyleRequest{
						Range: pRange,
						TextStyle: &docs.TextStyle{
							FontSize: &docs.Dimension{Magnitude: 10.0, Unit: "PT"},
						},
						Fields: "fontSize",
					},
				})
			} else {
				// CORPO DO TEXTO (ABNT): Justificado, entrelinhas 1.5, recuo de 1ª linha de 1.25 cm (35.433 pt), fonte 12pt
				formattedParas++
				requests = append(requests, &docs.Request{
					UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
						Range: pRange,
						ParagraphStyle: &docs.ParagraphStyle{
							Alignment:           "JUSTIFIED",
							LineSpacing:         ABNTLineSpacing,
							IndentFirstLine:     &docs.Dimension{Magnitude: ABNTIndentFirstLinePt, Unit: "PT"},
							IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceAbove:          &docs.Dimension{Magnitude: 0, Unit: "PT"},
							SpaceBelow:          &docs.Dimension{Magnitude: 0, Unit: "PT"},
							AvoidWidowAndOrphan: true,
						},
						Fields: "alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,avoidWidowAndOrphan",
					},
				})

				// Garante fonte 12pt no corpo comum
				requests = append(requests, &docs.Request{
					UpdateTextStyle: &docs.UpdateTextStyleRequest{
						Range: pRange,
						TextStyle: &docs.TextStyle{
							FontSize: &docs.Dimension{Magnitude: 12.0, Unit: "PT"},
						},
						Fields: "fontSize",
					},
				})
			}
		}
	}

	if len(requests) == 0 {
		return "Nenhum elemento encontrado para formatação no intervalo.", nil
	}

	updateReq := &docs.BatchUpdateDocumentRequest{Requests: requests}
	_, err = svc.Documents.BatchUpdate(documentID, updateReq).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao formatar documento para ABNT: %w", err)
	}

	return fmt.Sprintf("SUCESSO: Documento formatado em conformidade ABNT no intervalo [%d - %d].\n"+
		"- Parágrafos de corpo formatados (1.25cm / 1.5 / Justificado / 12pt): %d\n"+
		"- Títulos preservados com alinhamento e keepWithNext: %d\n"+
		"- Legendas/Fontes preservadas com 10pt e espaço simples: %d\n"+
		"- Margens ABNT A4 aplicadas: %v",
		start, end, formattedParas, preservedHeadings, preservedCaptions, applyMargins), nil
}
