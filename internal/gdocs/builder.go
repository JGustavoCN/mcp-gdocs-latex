package gdocs

import (
	"google.golang.org/api/docs/v1"
)

// Constantes métricas tipográficas ABNT (DTP Points: 1 cm = 72 / 2.54 pt ≈ 28.346 pt)
const (
	ABNTDefaultFontFamily      = "Arial"           // Família primária padrão ABNT NBR 14724
	ABNTTimesFontFamily        = "Times New Roman" // Família secundária permitida ABNT
	ABNTFontSizeBodyPt         = 12.0              // Corpo e títulos de seções (NBR 14724 / NBR 6024)
	ABNTFontSizeCaptionPt      = 10.0              // Legendas, fontes, notas de rodapé e citações longas
	ABNTFontSizeHeadingPt      = 12.0              // Títulos de seções na ABNT são sempre 12 pt!
	ABNTIndentFirstLinePt      = 35.433            // 1.25 cm (recuo padrão de parágrafos de corpo)
	ABNTIndentLongQuotePt      = 113.386           // 4.00 cm (recuo para citações diretas longas)
	ABNTIndentPreTextualNotePt = 226.772           // 8.00 cm (recuo para nota institucional na folha de rosto, metade da página A4)
	ABNTLineSpacing            = 150               // 1.5 entrelinhas (150%)
	ABNTLineSpacingSingle      = 100               // 1.0 entrelinhas simples (100%)
	ABNTMarginTopPt            = 85.039            // 3.00 cm
	ABNTMarginLeftPt           = 85.039            // 3.00 cm
	ABNTMarginBottomPt         = 56.693            // 2.00 cm
	ABNTMarginRightPt          = 56.693            // 2.00 cm
	ABNTPageWidthPt            = 595.276           // 21.00 cm (A4)
	ABNTPageHeightPt           = 841.890           // 29.70 cm (A4)
)

// BatchUpdateBuilder auxilia na montagem atômica e segura de lotes de requisições para a Google Docs API
type BatchUpdateBuilder struct {
	requests     []*docs.Request
	writeControl *docs.WriteControl
}

// NewBatchUpdateBuilder instancia um novo construtor
func NewBatchUpdateBuilder() *BatchUpdateBuilder {
	return &BatchUpdateBuilder{
		requests: make([]*docs.Request, 0),
	}
}

// SetWriteControl define a trava de concorrência por revisão
func (b *BatchUpdateBuilder) SetWriteControl(revisionID string) *BatchUpdateBuilder {
	if revisionID != "" {
		b.writeControl = &docs.WriteControl{
			RequiredRevisionId: revisionID,
		}
	}
	return b
}

// AddRequest adiciona uma requisição genérica ao lote
func (b *BatchUpdateBuilder) AddRequest(req *docs.Request) *BatchUpdateBuilder {
	if req != nil {
		b.requests = append(b.requests, req)
	}
	return b
}

// AddRequests adiciona múltiplas requisições ao lote
func (b *BatchUpdateBuilder) AddRequests(reqs []*docs.Request) *BatchUpdateBuilder {
	for _, req := range reqs {
		b.AddRequest(req)
	}
	return b
}

// Requests retorna a lista de requisições acumuladas
func (b *BatchUpdateBuilder) Requests() []*docs.Request {
	return b.requests
}

// Build gera o BatchUpdateDocumentRequest final
func (b *BatchUpdateBuilder) Build() *docs.BatchUpdateDocumentRequest {
	req := &docs.BatchUpdateDocumentRequest{
		Requests: b.requests,
	}
	if b.writeControl != nil {
		req.WriteControl = b.writeControl
	}
	return req
}

// ClampRange impõe truncamento defensivo automático para nunca ultrapassar ou deletar o caractere final \n
func ClampRange(r *docs.Range, bodyEndIndex int64) {
	if r == nil || bodyEndIndex <= 1 {
		return
	}
	maxAllowed := bodyEndIndex - 1
	if r.EndIndex > maxAllowed {
		r.EndIndex = maxAllowed
	}
	if r.StartIndex > r.EndIndex {
		r.StartIndex = r.EndIndex
	}
}

// BuildABNTReplaceSequence gera a sequência canônica atômica de 5 passos com blindagens anti-bleeding
func BuildABNTReplaceSequence(startIndex, endIndex, bodyEndIndex int64, newText string) []*docs.Request {
	sanitized := SanitizeText(newText)

	// Clamping defensivo no range a deletar
	targetRange := &docs.Range{
		StartIndex: startIndex,
		EndIndex:   endIndex,
	}
	ClampRange(targetRange, bodyEndIndex)

	var reqs []*docs.Request

	// Passo 1: DeleteContentRangeRequest
	if targetRange.EndIndex > targetRange.StartIndex {
		reqs = append(reqs, &docs.Request{
			DeleteContentRange: &docs.DeleteContentRangeRequest{
				Range: targetRange,
			},
		})
	}

	// Guarda de texto vazio (Deleção pura): se novo texto for vazio, encerra por aqui sem disparar HTTP 400
	if sanitized == "" {
		return reqs
	}

	// Passo 2: InsertTextRequest
	reqs = append(reqs, &docs.Request{
		InsertText: &docs.InsertTextRequest{
			Location: &docs.Location{
				Index: startIndex,
			},
			Text: sanitized,
		},
	})

	u16Len := UTF16Length(sanitized)
	affectedRange := &docs.Range{
		StartIndex: startIndex,
		EndIndex:   startIndex + u16Len,
	}

	// Passo 3: DeleteParagraphBulletsRequest (elimina marcadores residuais de lista / bullets vazios)
	reqs = append(reqs, &docs.Request{
		DeleteParagraphBullets: &docs.DeleteParagraphBulletsRequest{
			Range: affectedRange,
		},
	})

	// Passo 4: UpdateParagraphStyleRequest (ABNT: Normal Text, Justificado, 1.5, recuo 1.25cm e indentStart zerado)
	reqs = append(reqs, &docs.Request{
		UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
			Range: affectedRange,
			ParagraphStyle: &docs.ParagraphStyle{
				NamedStyleType: "NORMAL_TEXT",
				Alignment:      "JUSTIFIED",
				LineSpacing:    ABNTLineSpacing,
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
				AvoidWidowAndOrphan: true,
			},
			Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,avoidWidowAndOrphan",
		},
	})

	// Passo 5: UpdateTextStyleRequest (limpeza estrita de bold e italic para impedir bold bleeding)
	reqs = append(reqs, &docs.Request{
		UpdateTextStyle: &docs.UpdateTextStyleRequest{
			Range: affectedRange,
			TextStyle: &docs.TextStyle{
				Bold:      false,
				Italic:    false,
				Underline: false,
			},
			Fields: "bold,italic,underline",
		},
	})

	return reqs
}

// BuildABNTDocumentMarginsRequest gera a requisição de padronização de margens e folha A4 ABNT
func BuildABNTDocumentMarginsRequest() *docs.Request {
	return &docs.Request{
		UpdateDocumentStyle: &docs.UpdateDocumentStyleRequest{
			DocumentStyle: &docs.DocumentStyle{
				MarginTop: &docs.Dimension{
					Magnitude: ABNTMarginTopPt,
					Unit:      "PT",
				},
				MarginLeft: &docs.Dimension{
					Magnitude: ABNTMarginLeftPt,
					Unit:      "PT",
				},
				MarginBottom: &docs.Dimension{
					Magnitude: ABNTMarginBottomPt,
					Unit:      "PT",
				},
				MarginRight: &docs.Dimension{
					Magnitude: ABNTMarginRightPt,
					Unit:      "PT",
				},
				PageSize: &docs.Size{
					Width: &docs.Dimension{
						Magnitude: ABNTPageWidthPt,
						Unit:      "PT",
					},
					Height: &docs.Dimension{
						Magnitude: ABNTPageHeightPt,
						Unit:      "PT",
					},
				},
			},
			Fields: "marginTop,marginLeft,marginBottom,marginRight,pageSize",
		},
	}
}
