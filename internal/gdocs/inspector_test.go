package gdocs

import (
	"testing"
)

func TestElementLayoutDetailModel(t *testing.T) {
	d := ElementLayoutDetail{
		Index:             1,
		StartIndex:        10,
		EndIndex:          50,
		Type:              "PARAGRAPH",
		StyleName:         "NORMAL_TEXT",
		Alignment:         "JUSTIFIED",
		IndentStartPt:     226.772,
		IndentFirstLinePt: 226.772,
		FontSizePt:        10.0,
		FontFamily:        "Arial",
		HasPageBreak:      false,
		TextSnippet:       "Exemplo de preâmbulo",
	}

	if d.IndentStartPt != 226.772 || d.IndentFirstLinePt != 226.772 {
		t.Errorf("Recuo incorreto: Start=%.3f, First=%.3f", d.IndentStartPt, d.IndentFirstLinePt)
	}

	if d.Alignment != "JUSTIFIED" {
		t.Errorf("Alinhamento esperado JUSTIFIED, obtido %s", d.Alignment)
	}
}

