package gdocs

import (
	"google.golang.org/api/docs/v1"
	"testing"
)

func TestClampRange(t *testing.T) {
	bodyEnd := int64(100)

	r1 := &docs.Range{StartIndex: 10, EndIndex: 120}
	ClampRange(r1, bodyEnd)
	if r1.EndIndex != 99 {
		t.Errorf("ClampRange esperava EndIndex 99, obteve %d", r1.EndIndex)
	}

	r2 := &docs.Range{StartIndex: 10, EndIndex: 50}
	ClampRange(r2, bodyEnd)
	if r2.EndIndex != 50 {
		t.Errorf("ClampRange não deveria alterar range válido, obteve %d", r2.EndIndex)
	}
}

func TestBuildABNTReplaceSequence(t *testing.T) {
	t.Run("Texto comum - 5 passos", func(t *testing.T) {
		reqs := BuildABNTReplaceSequence(100, 150, 500, "Parágrafo formatado.")
		if len(reqs) != 5 {
			t.Fatalf("Esperava 5 requisições canônicas, obteve %d", len(reqs))
		}

		if reqs[0].DeleteContentRange == nil {
			t.Errorf("Passo 1 deve ser DeleteContentRange")
		}
		if reqs[1].InsertText == nil {
			t.Errorf("Passo 2 deve ser InsertText")
		}
		if reqs[2].DeleteParagraphBullets == nil {
			t.Errorf("Passo 3 deve ser DeleteParagraphBullets")
		}
		if reqs[3].UpdateParagraphStyle == nil {
			t.Errorf("Passo 4 deve ser UpdateParagraphStyle")
		}
		if reqs[4].UpdateTextStyle == nil {
			t.Errorf("Passo 5 deve ser UpdateTextStyle")
		}
	})

	t.Run("Texto vazio - guarda de deleção pura (1 passo)", func(t *testing.T) {
		reqs := BuildABNTReplaceSequence(100, 150, 500, "")
		if len(reqs) != 1 {
			t.Fatalf("Esperava exatamente 1 requisição para deleção pura, obteve %d", len(reqs))
		}
		if reqs[0].DeleteContentRange == nil {
			t.Errorf("Passo 1 deve ser DeleteContentRange")
		}
	})
}
