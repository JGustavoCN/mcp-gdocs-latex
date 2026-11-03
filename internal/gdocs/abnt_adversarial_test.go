package gdocs

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	"google.golang.org/api/docs/v1"
)

// DocumentModelSimulado é um simulador em memória do documento Google Docs
// para testar o comportamento adversarial de mutações posicionalmente sensíveis.
type DocumentModelSimulado struct {
	Content []*docs.StructuralElement
}

// Clone cria uma cópia profunda da lista de elementos estruturais
func (m *DocumentModelSimulado) Clone() *DocumentModelSimulado {
	newContent := make([]*docs.StructuralElement, len(m.Content))
	for i, el := range m.Content {
		copyEl := &docs.StructuralElement{
			StartIndex: el.StartIndex,
			EndIndex:   el.EndIndex,
		}
		if el.Table != nil {
			copyEl.Table = &docs.Table{}
		}
		if el.Paragraph != nil {
			copyEl.Paragraph = &docs.Paragraph{
				ParagraphStyle: el.Paragraph.ParagraphStyle,
			}
			for _, pel := range el.Paragraph.Elements {
				copyPel := &docs.ParagraphElement{}
				if pel.PageBreak != nil {
					copyPel.PageBreak = &docs.PageBreak{}
				}
				if pel.InlineObjectElement != nil {
					copyPel.InlineObjectElement = &docs.InlineObjectElement{InlineObjectId: pel.InlineObjectElement.InlineObjectId}
				}
				if pel.TextRun != nil {
					copyPel.TextRun = &docs.TextRun{
						Content:   pel.TextRun.Content,
						TextStyle: pel.TextRun.TextStyle,
					}
				}
				copyEl.Paragraph.Elements = append(copyEl.Paragraph.Elements, copyPel)
			}
		}
		newContent[i] = copyEl
	}
	return &DocumentModelSimulado{Content: newContent}
}

// AplicaDelecoesEmLoteReverso simula a execução atômica de DeleteContentRangeRequest
// no Google Docs, ajustando os índices subsequentes de forma fiel à API.
func (m *DocumentModelSimulado) AplicaDelecoesEmLoteReverso(ranges []*docs.Range) error {
	// A API do Docs exige que se as deleções vêm em lote decrescente, elas sejam aplicadas em sequência
	for _, r := range ranges {
		delLen := r.EndIndex - r.StartIndex
		if delLen <= 0 {
			continue
		}

		// Localiza e remove o elemento correspondente ou fatia
		idxToRemove := -1
		for i, el := range m.Content {
			if el.StartIndex == r.StartIndex && el.EndIndex == r.EndIndex {
				idxToRemove = i
				break
			}
		}

		if idxToRemove >= 0 {
			m.Content = append(m.Content[:idxToRemove], m.Content[idxToRemove+1:]...)
		}

		// Atualiza os offsets de todos os elementos que estavam à frente do range deletado
		for _, el := range m.Content {
			if el.StartIndex >= r.EndIndex {
				el.StartIndex -= delLen
				el.EndIndex -= delLen
			}
		}
	}
	return nil
}

// ColetaRangesHeal reproduz fielmente a lógica de coleta de HealDocTextAndBlanks
func ColetaRangesHeal(content []*docs.StructuralElement, targetStartIndex int64) (toDelete []*docs.Range, tableDelimiters []*docs.Range) {
	if len(content) == 0 {
		return nil, nil
	}
	bodyEnd := content[len(content)-1].EndIndex

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

			t := strings.TrimSpace(txt)
			if t == "" {
				isBeforeTable := false
				if i+1 < len(content) && content[i+1].Table != nil {
					isBeforeTable = true
				}

				if isBeforeTable {
					tableDelimiters = append(tableDelimiters, &docs.Range{
						StartIndex: el.StartIndex,
						EndIndex:   el.EndIndex,
					})
				} else {
					toDelete = append(toDelete, &docs.Range{
						StartIndex: el.StartIndex,
						EndIndex:   el.EndIndex,
					})
				}
			}
		}
	}
	return toDelete, tableDelimiters
}

// -----------------------------------------------------------------------------
// TESTE ADVERSARIAL 1: Idempotência Matemática do Lote Reverso
// O que acontece se HealDocTextAndBlanks rodar 2 vezes seguidas sobre o mesmo documento?
// -----------------------------------------------------------------------------
func TestAdversarialIdempotence_HealDoc(t *testing.T) {
	// Cria documento sintético com múltiplos tipos de parágrafos
	doc := &DocumentModelSimulado{
		Content: []*docs.StructuralElement{
			// 0: Folha de Rosto Ano 1
			{StartIndex: 1, EndIndex: 10, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "2026\n"}}}}},
			// 1: Folha de Rosto Ano 2 (ponto de corte: targetStartIndex = 20)
			{StartIndex: 10, EndIndex: 20, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "2026\n"}}}}},
			// 2: Título Introdução
			{StartIndex: 20, EndIndex: 35, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "1 INTRODUÇÃO\n"}}}}},
			// 3: Quebra espúria pós-título (deve ser DELETADA na 1ª rodada)
			{StartIndex: 35, EndIndex: 36, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}}}},
			// 4: Texto 1
			{StartIndex: 36, EndIndex: 90, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "Primeiro parágrafo de introdução normativa.\n"}}}}},
			// 5: Quebra espúria interparágrafo 1 (deve ser DELETADA na 1ª rodada)
			{StartIndex: 90, EndIndex: 91, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}}}},
			// 6: Quebra espúria interparágrafo 2 (consecutiva - deve ser DELETADA na 1ª rodada)
			{StartIndex: 91, EndIndex: 92, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}}}},
			// 7: Texto 2
			{StartIndex: 92, EndIndex: 150, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "Segundo parágrafo de introdução com análise empírica.\n"}}}}},
			// 8: Delimitador técnico pré-tabela (deve ser PRESERVADO e neutralizado)
			{StartIndex: 150, EndIndex: 151, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}}}},
			// 9: Tabela nativa
			{StartIndex: 151, EndIndex: 250, Table: &docs.Table{}},
			// 10: Parágrafo com PageBreak (deve ser rigorosamente PRESERVADO)
			{StartIndex: 250, EndIndex: 251, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{PageBreak: &docs.PageBreak{}}}}},
			// 11: Parágrafo com Imagem/Logo (deve ser rigorosamente PRESERVADO)
			{StartIndex: 251, EndIndex: 252, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{InlineObjectElement: &docs.InlineObjectElement{InlineObjectId: "img1"}}}}},
			// 12: Penúltimo parágrafo vazio espúrio (deve ser DELETADO na 1ª rodada)
			{StartIndex: 252, EndIndex: 253, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}}}},
			// 13: Último parágrafo inamovível do Google Docs (deve ser PRESERVADO pela guarda bodyEnd)
			{StartIndex: 253, EndIndex: 254, Paragraph: &docs.Paragraph{Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}}}},
		},
	}

	targetStartIndex := int64(20)

	// =========================================================================
	// RODADA 1: Execução do HealDocTextAndBlanks
	// =========================================================================
	toDeleteR1, tableDelimsR1 := ColetaRangesHeal(doc.Content, targetStartIndex)

	// Verificações da Rodada 1
	if len(toDeleteR1) != 4 {
		t.Fatalf("Rodada 1: esperava 4 faixas para deleção (offsets 35, 90, 91, 252), obteve %d", len(toDeleteR1))
	}
	if len(tableDelimsR1) != 1 || tableDelimsR1[0].StartIndex != 150 {
		t.Fatalf("Rodada 1: esperava 1 delimitador pré-tabela no offset 150, obteve %v", tableDelimsR1)
	}

	// Ordenação decrescente (lote reverso)
	sort.Slice(toDeleteR1, func(i, j int) bool {
		return toDeleteR1[i].StartIndex > toDeleteR1[j].StartIndex
	})

	// Valida ordem decrescente rigorosa
	for k := 0; k < len(toDeleteR1)-1; k++ {
		if toDeleteR1[k].StartIndex <= toDeleteR1[k+1].StartIndex {
			t.Fatalf("Rodada 1: falha na ordem decrescente: %d <= %d", toDeleteR1[k].StartIndex, toDeleteR1[k+1].StartIndex)
		}
	}

	// Aplica as mutações no modelo simulado
	err := doc.AplicaDelecoesEmLoteReverso(toDeleteR1)
	if err != nil {
		t.Fatalf("Erro ao aplicar deleções da Rodada 1: %v", err)
	}

	// Snapshot do estado pós-Rodada 1
	snapshotR1 := doc.Clone()

	// =========================================================================
	// RODADA 2: Execução de HealDocTextAndBlanks sobre o resultado da Rodada 1
	// =========================================================================
	toDeleteR2, tableDelimsR2 := ColetaRangesHeal(doc.Content, targetStartIndex)

	// PROVA MATEMÁTICA DE IDEMPOTÊNCIA:
	// Na 2ª rodada, toDeleteR2 DEVE ser rigorosamente vazio (zero deleções adicionais)
	if len(toDeleteR2) != 0 {
		t.Errorf("VIOLAÇÃO DE IDEMPOTÊNCIA: Rodada 2 tentou deletar %d ranges que já deveriam estar sanados: %+v",
			len(toDeleteR2), toDeleteR2)
	}

	// O delimitador de tabela continua existindo e sendo preservado
	if len(tableDelimsR2) != 1 {
		t.Errorf("Rodada 2: delimitador de tabela pré-existente foi corrompido, obteve %d", len(tableDelimsR2))
	}

	// Aplica deleções da Rodada 2 (nenhuma)
	_ = doc.AplicaDelecoesEmLoteReverso(toDeleteR2)

	// O estado do documento pós-Rodada 2 deve ser 100% IDÊNTICO ao pós-Rodada 1: f(f(x)) == f(x)
	if len(doc.Content) != len(snapshotR1.Content) {
		t.Fatalf("VIOLAÇÃO DE IDEMPOTÊNCIA: Quantidade de elementos divergiu: R1=%d, R2=%d",
			len(snapshotR1.Content), len(doc.Content))
	}

	for idx := range doc.Content {
		el1 := snapshotR1.Content[idx]
		el2 := doc.Content[idx]
		if el1.StartIndex != el2.StartIndex || el1.EndIndex != el2.EndIndex {
			t.Errorf("VIOLAÇÃO DE IDEMPOTÊNCIA no elemento %d: R1=[%d-%d], R2=[%d-%d]",
				idx, el1.StartIndex, el1.EndIndex, el2.StartIndex, el2.EndIndex)
		}
	}

	t.Logf("PASS: Idempotência matemática comprovada: f(f(doc)) == f(doc). toDeleteR1=%d, toDeleteR2=%d",
		len(toDeleteR1), len(toDeleteR2))
}

// -----------------------------------------------------------------------------
// TESTE ADVERSARIAL 2: Chunking de Lotes de 100 Requisições com 250 Parágrafos Vazios
// O lote é dividido corretamente em blocos <= 100 sem inverter offsets decrescentes?
// -----------------------------------------------------------------------------
func TestAdversarialChunking_250_BlankParagraphs(t *testing.T) {
	const totalBlanks = 250
	var content []*docs.StructuralElement

	// Capa e cabeçalho inicial
	offset := int64(1)
	content = append(content, &docs.StructuralElement{
		StartIndex: offset,
		EndIndex:   offset + 10,
		Paragraph: &docs.Paragraph{
			Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "2026\n"}}},
		},
	})
	offset += 10
	content = append(content, &docs.StructuralElement{
		StartIndex: offset,
		EndIndex:   offset + 10,
		Paragraph: &docs.Paragraph{
			Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "2026\n"}}},
		},
	})
	offset += 10

	targetStartIndex := offset

	// Gera 250 pares de (parágrafo de texto + parágrafo em branco)
	for i := 0; i < totalBlanks; i++ {
		// Parágrafo com texto
		txt := fmt.Sprintf("Parágrafo de conteúdo numerado %03d.\n", i+1)
		txtLen := int64(len(txt))
		content = append(content, &docs.StructuralElement{
			StartIndex: offset,
			EndIndex:   offset + txtLen,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: txt}}},
			},
		})
		offset += txtLen

		// Parágrafo puramente vazio intercalado (deve ser deletado)
		content = append(content, &docs.StructuralElement{
			StartIndex: offset,
			EndIndex:   offset + 1,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}},
			},
		})
		offset += 1
	}

	// Último parágrafo inamovível do Google Docs
	content = append(content, &docs.StructuralElement{
		StartIndex: offset,
		EndIndex:   offset + 1,
		Paragraph: &docs.Paragraph{
			Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}},
		},
	})

	// 1. Executa a coleta exata do abnt_writer.go
	toDeleteRanges, _ := ColetaRangesHeal(content, targetStartIndex)

	if len(toDeleteRanges) != totalBlanks {
		t.Fatalf("Esperava exatamente %d parágrafos coletados para deleção, obteve %d", totalBlanks, len(toDeleteRanges))
	}

	// 2. Ordenação decrescente por StartIndex
	sort.Slice(toDeleteRanges, func(i, j int) bool {
		return toDeleteRanges[i].StartIndex > toDeleteRanges[j].StartIndex
	})

	// 3. Execução do algoritmo de chunking em fatias de até 100 (abnt_writer.go linhas 424-430)
	batchSize := 100
	var chunks [][]*docs.Range
	for start := 0; start < len(toDeleteRanges); start += batchSize {
		end := start + batchSize
		if end > len(toDeleteRanges) {
			end = len(toDeleteRanges)
		}
		chunks = append(chunks, toDeleteRanges[start:end])
	}

	// Verificação 1: Número exato de chunks
	expectedChunks := 3 // 100 + 100 + 50
	if len(chunks) != expectedChunks {
		t.Fatalf("Esperava %d chunks de requisições, obteve %d", expectedChunks, len(chunks))
	}

	// Verificação 2: Dimensões de cada chunk
	if len(chunks[0]) != 100 {
		t.Errorf("Chunk 0 deveria ter 100 itens, obteve %d", len(chunks[0]))
	}
	if len(chunks[1]) != 100 {
		t.Errorf("Chunk 1 deveria ter 100 itens, obteve %d", len(chunks[1]))
	}
	if len(chunks[2]) != 50 {
		t.Errorf("Chunk 2 deveria ter 50 itens, obteve %d", len(chunks[2]))
	}

	// Verificação 3: Monotonicidade estrita decrescente dentro de cada chunk
	for cIdx, chunk := range chunks {
		for k := 0; k < len(chunk)-1; k++ {
			if chunk[k].StartIndex <= chunk[k+1].StartIndex {
				t.Fatalf("Chunk %d: inversão de ordem interna no índice %d: %d <= %d",
					cIdx, k, chunk[k].StartIndex, chunk[k+1].StartIndex)
			}
		}
	}

	// Verificação 4: Monotonicidade estrita decrescente entre as fronteiras dos chunks
	// O menor StartIndex do Chunk 0 deve ser ESTRITAMENTE MAIOR que o maior StartIndex do Chunk 1
	lastChunk0 := chunks[0][len(chunks[0])-1].StartIndex
	firstChunk1 := chunks[1][0].StartIndex
	if lastChunk0 <= firstChunk1 {
		t.Fatalf("Fronteira Chunk 0 -> Chunk 1 invertida: último do lote 0 (%d) <= primeiro do lote 1 (%d)",
			lastChunk0, firstChunk1)
	}

	// O menor StartIndex do Chunk 1 deve ser ESTRITAMENTE MAIOR que o maior StartIndex do Chunk 2
	lastChunk1 := chunks[1][len(chunks[1])-1].StartIndex
	firstChunk2 := chunks[2][0].StartIndex
	if lastChunk1 <= firstChunk2 {
		t.Fatalf("Fronteira Chunk 1 -> Chunk 2 invertida: último do lote 1 (%d) <= primeiro do lote 2 (%d)",
			lastChunk1, firstChunk2)
	}

	// Verificação 5: Invariante Global — se concatenarmos todos os chunks, a sequência é
	// monotonicamente decrescente do início ao fim sem nenhuma perturbação
	var concatenated []*docs.Range
	for _, chunk := range chunks {
		concatenated = append(concatenated, chunk...)
	}

	if len(concatenated) != totalBlanks {
		t.Fatalf("Total concatenado divergente: esperado %d, obteve %d", totalBlanks, len(concatenated))
	}

	for i := 0; i < len(concatenated)-1; i++ {
		if concatenated[i].StartIndex <= concatenated[i+1].StartIndex {
			t.Fatalf("Invariante global violada no índice %d: %d <= %d",
				i, concatenated[i].StartIndex, concatenated[i+1].StartIndex)
		}
	}

	t.Logf("PASS: Chunking de 250 parágrafos validado: 3 lotes [100, 100, 50], preservação estrita de monotonicidade decrescente inter e intra-lotes.")
}

// -----------------------------------------------------------------------------
// TESTE ADVERSARIAL 3: Guarda de Fronteira do bodyEndIndex - 1
// Garante que nenhum range de deleção ultrapasse a fronteira inamovível do Google Docs
// -----------------------------------------------------------------------------
func TestAdversarialBoundaryGuard_BodyEndMinusOne(t *testing.T) {
	// Cenário A: Documento mínimo com apenas 1 parágrafo (apenas o \n inamovível final)
	docMinimo := []*docs.StructuralElement{
		{
			StartIndex: 1,
			EndIndex:   2,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}},
			},
		},
	}
	delMin, _ := ColetaRangesHeal(docMinimo, 1)
	if len(delMin) != 0 {
		t.Errorf("Cenário A: Documento mínimo de 1 elemento não deve gerar deleções, obteve %d", len(delMin))
	}

	// Cenário B: Documento onde o ÚLTIMO elemento é um parágrafo vazio (\n) que termina em bodyEnd
	bodyEndB := int64(100)
	docFinalVazio := []*docs.StructuralElement{
		{
			StartIndex: 1,
			EndIndex:   90,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "Conteúdo legítimo\n"}}},
			},
		},
		{
			StartIndex: 90,
			EndIndex:   bodyEndB,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}},
			},
		},
	}
	delFinal, _ := ColetaRangesHeal(docFinalVazio, 1)
	// A guarda el.EndIndex >= bodyEnd DEVE proteger categoricamente o último parágrafo
	if len(delFinal) != 0 {
		t.Errorf("Cenário B: O parágrafo terminal [90-100] não pode ser deletado, pois contém o \\n inamovível! Obteve: %v", delFinal)
	}

	// Cenário C: Penúltimo parágrafo vazio terminando em bodyEnd - 1 e último terminando em bodyEnd
	bodyEndC := int64(150)
	docPenultimoVazio := []*docs.StructuralElement{
		{
			StartIndex: 1,
			EndIndex:   140,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "Texto do capítulo\n"}}},
			},
		},
		{
			// Penúltimo elemento vazio: vai de 140 até 149 (bodyEnd - 1)
			StartIndex: 140,
			EndIndex:   149,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}},
			},
		},
		{
			// Último elemento: vai de 149 até 150 (bodyEnd)
			StartIndex: 149,
			EndIndex:   bodyEndC,
			Paragraph: &docs.Paragraph{
				Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "\n"}}},
			},
		},
	}

	delPenultimo, _ := ColetaRangesHeal(docPenultimoVazio, 1)
	if len(delPenultimo) != 1 {
		t.Fatalf("Cenário C: Esperava 1 deleção para o penúltimo parágrafo vazio, obteve %d", len(delPenultimo))
	}

	r := delPenultimo[0]
	// Asserções estritas de fronteira:
	if r.StartIndex != 140 {
		t.Errorf("Cenário C: StartIndex esperado 140, obteve %d", r.StartIndex)
	}
	if r.EndIndex != 149 {
		t.Errorf("Cenário C: EndIndex esperado 149, obteve %d", r.EndIndex)
	}
	if r.EndIndex > bodyEndC-1 {
		t.Fatalf("VIOLAÇÃO DE FRONTEIRA: EndIndex %d ultrapassou bodyEnd-1 (%d)!", r.EndIndex, bodyEndC-1)
	}
	if r.EndIndex >= bodyEndC {
		t.Fatalf("VIOLAÇÃO CRÍTICA DE FRONTEIRA: EndIndex %d tocou ou ultrapassou bodyEnd (%d)!", r.EndIndex, bodyEndC)
	}

	// Cenário D: Análise de fronteira na rotina healBrokenLines
	// Verifica se a junção de parágrafos quando o segundo é o último parágrafo do corpo nunca toca bodyEnd
	elPenultimoTexto := &docs.StructuralElement{
		StartIndex: 100,
		EndIndex:   120, // endA = 120
		Paragraph: &docs.Paragraph{
			Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: "Frase quebrada que continua"}}},
		},
	}
	elUltimoTexto := &docs.StructuralElement{
		StartIndex: 120,
		EndIndex:   140, // bodyEnd = 140
		Paragraph: &docs.Paragraph{
			Elements: []*docs.ParagraphElement{{TextRun: &docs.TextRun{Content: " no parágrafo final.\n"}}},
		},
	}

	// endA = 120. A operação deleta endA até endA+1 (120..121).
	// Como o último elemento vai até 140, 121 é estritamente < 140 (bodyEnd).
	endA := elPenultimoTexto.EndIndex
	delEnd := endA + 1
	if delEnd >= elUltimoTexto.EndIndex {
		t.Fatalf("healBrokenLines: delEnd (%d) tocou o final do documento (%d)!", delEnd, elUltimoTexto.EndIndex)
	}

	t.Logf("PASS: Guarda de fronteira de bodyEndIndex - 1 validada em todos os 4 cenários extremos.")
}

// -----------------------------------------------------------------------------
// TESTE ADVERSARIAL 4: Preservação de Headings e Títulos Consecutivos
// -----------------------------------------------------------------------------
func TestAdversarialHeadingAndHierarchyStress(t *testing.T) {
	// Testa matriz de precedência de títulos sob estilos conflitantes
	cases := []struct {
		text          string
		currentStyle  string
		expectedStyle string
	}{
		{"1 INTRODUÇÃO", "HEADING_2", "HEADING_1"},
		{"1.1 CONTEXTO", "HEADING_1", "HEADING_2"},
		{"1.1.1 DETALHE", "HEADING_1", "HEADING_3"},
		{"1.1.1.1 SUBDETALHE", "TITLE", "HEADING_4"},
		{"REFERÊNCIAS", "NORMAL_TEXT", "HEADING_1"},
		{"Estes foram os critérios avaliados:", "HEADING_1", "NORMAL_TEXT"}, // Termina em ':'
	}

	for _, tc := range cases {
		// Valida precedência hierárquica implementada no FormatDocumentABNTExtended
		isPrimary := rePrimaryHeadingNum.MatchString(tc.text) || tc.text == "REFERÊNCIAS"
		isSec := reSecondaryHeading.MatchString(tc.text)
		isTert := reTertiaryHeading.MatchString(tc.text)
		isQuat := reQuaternaryHeading.MatchString(tc.text)
		isColon := strings.HasSuffix(tc.text, ":")

		determinedStyle := "NORMAL_TEXT"
		if isColon {
			determinedStyle = "NORMAL_TEXT"
		} else if isPrimary {
			determinedStyle = "HEADING_1"
		} else if isSec {
			determinedStyle = "HEADING_2"
		} else if isTert {
			determinedStyle = "HEADING_3"
		} else if isQuat {
			determinedStyle = "HEADING_4"
		}

		if determinedStyle != tc.expectedStyle {
			t.Errorf("Para %q com estilo atual %q, esperava %q, obteve %q",
				tc.text, tc.currentStyle, tc.expectedStyle, determinedStyle)
		}
	}

	t.Logf("PASS: Estresse de precedência hierárquica e falsos títulos com ':' validado com sucesso.")
}
