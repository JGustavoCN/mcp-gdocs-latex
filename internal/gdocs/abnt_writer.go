package gdocs

import (
	"context"
	"fmt"

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

// FormatDocumentABNT formata parágrafos do documento segundo as normas ABNT (NBR 14724)
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

	targetRange := &docs.Range{
		StartIndex: start,
		EndIndex:   end,
	}

	var requests []*docs.Request

	// 1. Aplica margens se solicitado
	if applyMargins {
		requests = append(requests, BuildABNTDocumentMarginsRequest())
	}

	// 2. Aplica estilo de parágrafo ABNT no range
	requests = append(requests, &docs.Request{
		UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
			Range: targetRange,
			ParagraphStyle: &docs.ParagraphStyle{
				Alignment:   "JUSTIFIED",
				LineSpacing: ABNTLineSpacing,
				IndentFirstLine: &docs.Dimension{
					Magnitude: ABNTIndentFirstLinePt,
					Unit:      "PT",
				},
				IndentStart: &docs.Dimension{
					Magnitude: 0,
					Unit:      "PT",
				},
				SpaceAbove: &docs.Dimension{
					Magnitude: 0,
					Unit:      "PT",
				},
				SpaceBelow: &docs.Dimension{
					Magnitude: 0,
					Unit:      "PT",
				},
			},
			Fields: "alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow",
		},
	})

	updateReq := &docs.BatchUpdateDocumentRequest{Requests: requests}
	_, err = svc.Documents.BatchUpdate(documentID, updateReq).Do()
	if err != nil {
		return "", fmt.Errorf("erro ao formatar documento para ABNT: %w", err)
	}

	return fmt.Sprintf("SUCESSO: Documento formatado em conformidade ABNT no intervalo [%d - %d].", start, end), nil
}
