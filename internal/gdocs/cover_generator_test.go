package gdocs

import (
	"testing"

	"google.golang.org/api/docs/v1"
)

func TestHasABNTCoverAndLogo_Empty(t *testing.T) {
	doc := &docs.Document{
		Body: &docs.Body{
			Content: []*docs.StructuralElement{},
		},
	}
	hasCover, hasLogo, hasTitlePage := HasABNTCoverAndLogo(doc)
	if hasCover || hasLogo || hasTitlePage {
		t.Errorf("Esperava false para todos em documento vazio")
	}
}

func TestHasABNTCoverAndLogo_Complete(t *testing.T) {
	doc := &docs.Document{
		InlineObjects: map[string]docs.InlineObject{
			"kix.logo": {},
		},
		Body: &docs.Body{
			Content: []*docs.StructuralElement{
				{
					StartIndex: 1,
					EndIndex:   100,
					Paragraph: &docs.Paragraph{
						Elements: []*docs.ParagraphElement{
							{
								InlineObjectElement: &docs.InlineObjectElement{
									InlineObjectId: "kix.logo",
								},
							},
							{
								TextRun: &docs.TextRun{
									Content: "INSTITUTO FEDERAL DE EDUCAÇÃO, CIÊNCIA E TECNOLOGIA DE SERGIPE\nCAMPUS LAGARTO\n",
								},
							},
						},
					},
				},
				{
					StartIndex: 101,
					EndIndex:   300,
					Paragraph: &docs.Paragraph{
						Elements: []*docs.ParagraphElement{
							{
								TextRun: &docs.TextRun{
									Content: "Trabalho de Conclusão de Curso apresentado ao curso de Bacharelado em Sistemas de Informação como requisito parcial para obtenção do grau de Bacharel.\nOrientador: Prof. Diego Armando\n",
								},
							},
						},
					},
				},
				{
					StartIndex: 301,
					EndIndex:   320,
					Paragraph: &docs.Paragraph{
						Elements: []*docs.ParagraphElement{
							{
								TextRun: &docs.TextRun{
									Content: "1 INTRODUÇÃO\n",
								},
							},
						},
					},
				},
			},
		},
	}

	hasCover, hasLogo, hasTitlePage := HasABNTCoverAndLogo(doc)
	if !hasCover {
		t.Errorf("Deveria detectar Capa com Instituto Federal")
	}
	if !hasLogo {
		t.Errorf("Deveria detectar Logo em InlineObjects")
	}
	if !hasTitlePage {
		t.Errorf("Deveria detectar Folha de Rosto com nota institucional")
	}
}

