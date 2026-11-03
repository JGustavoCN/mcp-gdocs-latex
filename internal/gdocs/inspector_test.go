package gdocs

import (
	"testing"
)

func TestElementLayoutDetailModel(t *testing.T) {
	d := ElementLayoutDetail{
		Index:              1,
		StartIndex:         10,
		EndIndex:           50,
		Type:               "PARAGRAPH",
		StyleName:          "NORMAL_TEXT",
		Alignment:          "JUSTIFIED",
		LineSpacingPercent: 150.0,
		IndentStartPt:      226.772,
		IndentFirstLinePt:  0.0,
		KeepWithNext:       true,
		AvoidWidowOrphan:   true,
		FontSizePt:         10.0,
		FontFamily:         "Arial",
		HasPageBreak:       false,
		TextSnippet:        "Exemplo de preâmbulo",
	}

	if d.IndentStartPt != 226.772 || d.IndentFirstLinePt != 0.0 {
		t.Errorf("Recuo incorreto: Start=%.3f, First=%.3f", d.IndentStartPt, d.IndentFirstLinePt)
	}

	if d.Alignment != "JUSTIFIED" {
		t.Errorf("Alinhamento esperado JUSTIFIED, obtido %s", d.Alignment)
	}

	if d.LineSpacingPercent != 150.0 {
		t.Errorf("Entrelinhas esperado 150.0, obtido %.1f", d.LineSpacingPercent)
	}

	geom := DocumentGeometryDetail{
		PageWidthPt:  595.276,
		PageHeightPt: 841.890,
		MarginTopPt:  85.039,
		MarginLeftPt: 85.039,
	}
	if geom.PageWidthPt < 590 || geom.MarginTopPt < 84 {
		t.Errorf("Dimensões de geometria incorretas: %+v", geom)
	}
}

