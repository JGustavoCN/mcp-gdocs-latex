package gdocs

import (
	"sort"
	"strings"
	"testing"

	"google.golang.org/api/docs/v1"
)

func TestShouldMergeParagraphs(t *testing.T) {
	tests := []struct {
		name        string
		curr        string
		next        string
		isSpecial1  bool
		isSpecial2  bool
		expected    bool
	}{
		{
			name:       "Hard wrap quebrado no meio da frase (caso clássico LaTeX)",
			curr:       "que retornam dados estruturados, permitindo ao modelo de linguagem recuperar os",
			next:       "elementos pertinentes à solicitação e utilizá-los na elaboração do documento.",
			isSpecial1: false,
			isSpecial2: false,
			expected:   true,
		},
		{
			name:       "Frase com vírgula no final que continua em minúscula",
			curr:       "Além da aplicação específica à elaboração de requerimentos, a disponibilização",
			next:       "estruturada de informações e funcionalidades institucionais pode favorecer o",
			isSpecial1: false,
			isSpecial2: false,
			expected:   true,
		},
		{
			name:       "Parágrafo com ponto final seguido de nova frase iniciada em maiúscula",
			curr:       "Essa abordagem possibilita investigar como a ampliação do contexto informacional disponível ao modelo influencia a conformidade das minutas geradas.",
			next:       "Além da aplicação específica à elaboração de requerimentos, a disponibilização estruturada de informações pode favorecer o desenvolvimento.",
			isSpecial1: false,
			isSpecial2: false,
			expected:   false,
		},
		{
			name:       "Parágrafo terminando em dois pontos (introdução de lista/alínea)",
			curr:       "Os seguintes requisitos foram levantados:",
			next:       "a) Autenticação via OAuth2 com PKCE;",
			isSpecial1: false,
			isSpecial2: false,
			expected:   false,
		},
		{
			name:       "Título primário seguido de texto de abertura",
			curr:       "1 INTRODUÇÃO",
			next:       "Este trabalho investiga a integração de IA na redação acadêmica.",
			isSpecial1: true,
			isSpecial2: false,
			expected:   false,
		},
		{
			name:       "Título secundário",
			curr:       "2.1 CONTEXTUALIZAÇÃO DO PROBLEMA",
			next:       "O contexto atual apresenta desafios de conformidade.",
			isSpecial1: true,
			isSpecial2: false,
			expected:   false,
		},
		{
			name:       "Legenda seguida de fonte",
			curr:       "Figura 1 – Diagrama de arquitetura da solução",
			next:       "Fonte: Autor (2026)",
			isSpecial1: true,
			isSpecial2: true,
			expected:   false,
		},
		{
			name:       "Hifenização de quebra de linha (sis- / tema)",
			curr:       "A arquitetura proposta para o sis-",
			next:       "tema de geração de documentos utiliza o protocolo MCP.",
			isSpecial1: false,
			isSpecial2: false,
			expected:   true,
		},
		{
			name:       "Item com marcador seguido de outro item",
			curr:       "• Primeiro requisito funcional da plataforma.",
			next:       "• Segundo requisito funcional da plataforma.",
			isSpecial1: false,
			isSpecial2: false,
			expected:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := shouldMergeParagraphs(tc.curr, tc.next, tc.isSpecial1, tc.isSpecial2)
			if got != tc.expected {
				t.Errorf("shouldMergeParagraphs(%q, %q, %v, %v) = %v; esperado %v",
					tc.curr, tc.next, tc.isSpecial1, tc.isSpecial2, got, tc.expected)
			}
		})
	}
}

func TestPageBreakAndTableDelimiterProtectionInHeal(t *testing.T) {
	// Cria uma estrutura simulada de elementos do Google Docs
	// 0: Parágrafo com PageBreak (Folha de Rosto -> Resumo)
	// 1: Parágrafo com texto ("1 INTRODUÇÃO")
	// 2: Parágrafo vazio espúrio (deve ser deletado)
	// 3: Parágrafo vazio delimitador pré-tabela (deve ser PRESERVADO e neutralizado)
	// 4: Tabela
	// 5: Parágrafo com InlineObjectElement (Logo)
	// 6: Parágrafo vazio espúrio (deve ser deletado)

	content := []*docs.StructuralElement{
		{
			StartIndex: 100,
			EndIndex:   101,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{PageBreak: &docs.PageBreak{}},
				},
			},
		},
		{
			StartIndex: 101,
			EndIndex:   115,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{TextRun: &docs.TextRun{Content: "1 INTRODUÇÃO\n"}},
				},
			},
		},
		{
			StartIndex: 115,
			EndIndex:   116,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{TextRun: &docs.TextRun{Content: "\n"}},
				},
			},
		},
		{
			StartIndex: 116,
			EndIndex:   117,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{TextRun: &docs.TextRun{Content: "\n"}},
				},
			},
		},
		{
			StartIndex: 117,
			EndIndex:   200,
			Table:      &docs.Table{},
		},
		{
			StartIndex: 200,
			EndIndex:   201,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{InlineObjectElement: &docs.InlineObjectElement{InlineObjectId: "kix.abc"}},
				},
			},
		},
		{
			StartIndex: 201,
			EndIndex:   202,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{
					{TextRun: &docs.TextRun{Content: "\n"}},
				},
			},
		},
	}

	targetStartIndex := int64(100)
	bodyEnd := int64(300)

	var toDeleteRanges []*docs.Range
	var tableDelimiterRanges []*docs.Range

	for i := 0; i < len(content); i++ {
		el := content[i]
		if el.StartIndex < targetStartIndex || el.EndIndex >= bodyEnd {
			continue
		}
		if el.Table != nil {
			continue
		}
		if el.Paragraph != nil {
			hasPageBreak := false
			hasInlineObject := false
			txt := ""
			for _, pel := range el.Paragraph.Elements {
				if pel.PageBreak != nil {
					hasPageBreak = true
				}
				if pel.InlineObjectElement != nil {
					hasInlineObject = true
				}
				if pel.TextRun != nil {
					txt += pel.TextRun.Content
				}
			}

			if hasPageBreak || hasInlineObject {
				continue
			}

			tStr := strings.TrimSpace(txt)
			if tStr == "" {
				isBeforeTable := false
				if i+1 < len(content) && content[i+1].Table != nil {
					isBeforeTable = true
				}

				if isBeforeTable {
					tableDelimiterRanges = append(tableDelimiterRanges, &docs.Range{
						StartIndex: el.StartIndex,
						EndIndex:   el.EndIndex,
					})
				} else {
					toDeleteRanges = append(toDeleteRanges, &docs.Range{
						StartIndex: el.StartIndex,
						EndIndex:   el.EndIndex,
					})
				}
			}
		}
	}

	// 1. Verifica que o parágrafo com PageBreak (index 0) NÃO foi coletado para deleção
	for _, r := range toDeleteRanges {
		if r.StartIndex == 100 {
			t.Errorf("Parágrafo com PageBreak no offset 100 JAMAIS deve ser deletado!")
		}
	}

	// 2. Verifica que o parágrafo com InlineObject (index 5) NÃO foi coletado para deleção
	for _, r := range toDeleteRanges {
		if r.StartIndex == 200 {
			t.Errorf("Parágrafo com InlineObject no offset 200 JAMAIS deve ser deletado!")
		}
	}

	// 3. Verifica que o delimitador pré-tabela (offset 116) está em tableDelimiterRanges e NÃO em toDeleteRanges
	if len(tableDelimiterRanges) != 1 || tableDelimiterRanges[0].StartIndex != 116 {
		t.Errorf("tableDelimiterRanges deveria conter o offset 116, obteve: %v", tableDelimiterRanges)
	}
	for _, r := range toDeleteRanges {
		if r.StartIndex == 116 {
			t.Errorf("Parágrafo pré-tabela no offset 116 NÃO deve estar em toDeleteRanges!")
		}
	}

	// 4. Verifica que as quebras puramente vazias (offset 115 e 201) estão em toDeleteRanges
	if len(toDeleteRanges) != 2 {
		t.Fatalf("toDeleteRanges deveria conter 2 faixas vazias, obteve %d", len(toDeleteRanges))
	}
	if toDeleteRanges[0].StartIndex != 115 || toDeleteRanges[1].StartIndex != 201 {
		t.Errorf("Offsets esperados [115, 201], obteve [%d, %d]", toDeleteRanges[0].StartIndex, toDeleteRanges[1].StartIndex)
	}

	// 5. Verifica ordenação reversa (do maior para o menor offset)
	sort.Slice(toDeleteRanges, func(i, j int) bool {
		return toDeleteRanges[i].StartIndex > toDeleteRanges[j].StartIndex
	})

	if toDeleteRanges[0].StartIndex != 201 || toDeleteRanges[1].StartIndex != 115 {
		t.Errorf("Após ordenação reversa esperava [201, 115], obteve [%d, %d]",
			toDeleteRanges[0].StartIndex, toDeleteRanges[1].StartIndex)
	}
}

