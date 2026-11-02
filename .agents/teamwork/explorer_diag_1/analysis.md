# Relatório de Diagnóstico Técnico do Codebase Go (mcp-gdocs-latex)
**Marco**: R2 — Diagnóstico Técnico de Falhas de Entrelinhas, Quebras Espúrias e Ferramentas de Leitura/Inspeção  
**Data**: 2026-11-01  
**Autor**: explorer_diag_1 (Teamwork Explorer / Investigator)  
**Documento de Teste Auditado**: `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8`

---

## 1. Sumário Executivo

A investigação minuciosa no código-fonte Go do pacote `internal/gdocs` revelou as causas raízes exatas pelas quais o sistema não identificou nem corrigiu variações indesejadas de entrelinhas (ex.: 1,15 vs 1,5), manteve quebras de linha espúrias após a folha de rosto e deixou modelos de IA incapazes de diagnosticar visualmente esses defeitos:

1. **Causa Raiz de Entrelinhas (1,15 vs 1,5)**: O linter (`linter.go`) possui uma lacuna total: **não existe nenhuma linha de código que leia ou valide `LineSpacing`, `SpaceAbove` ou `SpaceBelow`**. O formatador (`abnt_writer.go`) possui a lógica de aplicação de 150%, mas **ignora completamente parágrafos dentro de tabelas nativas**, pula parágrafos vazios pós-introdução e possui heurística frágil para citações longas. Além disso, nenhuma ferramenta de leitura/inspeção expõe o campo `LineSpacing`, gerando cegueira estrutural completa.
2. **Causa Raiz de Quebras Espúrias após a Folha de Rosto**: O `cover_generator.go` insere dezenas de `\n` artificiais (18 na Capa, 7+5+14 na Folha de Rosto) para emular o `\vfill` do LaTeX. A lógica de saneamento em `abnt_writer.go` (`collapseBlankLines`) exige **duas ou mais linhas consecutivas (`consecutiveBlanks >= 2`)**, preservando sempre a primeira quebra após a folha de rosto. Além disso, a rotina de limpeza reversa em `cover_generator.go` (linhas 891–907) falha ao tentar deletar o parágrafo que abriga a própria `PageBreak`, suprimindo o erro com `_, _ = batchUpdate`.
3. **Causa Raiz de Carência nas Ferramentas de Leitura/Inspeção**: `read_doc_content` (em `reader.go`) converte o documento para Markdown plano e **omite completamente todos os estilos de parágrafo e quebras de página** (`pEl.PageBreak` é descartado porque não tem `TextRun`). `inspect_doc_layout` (em `inspector.go`) lista propriedades geométricas, mas **omitiu o campo `LineSpacing` da struct `ElementLayoutDetail`**, além de não inspecionar margens de página nem propriedades de cabeçalho/seção.

---

## 2. Diagnóstico Detalhado por Questão

---

### Questão 1: Variações Indesejadas de Entrelinhas (1,15 vs 1,5) e Espaçamento Vertical

#### 1.1 Análise do Linter (`internal/gdocs/linter.go`)
- **Ausência Absoluta de Auditoria de Entrelinhas**:  
  No arquivo `internal/gdocs/linter.go` (1.257 linhas), uma busca global por `LineSpacing` ou `line_spacing` retorna **zero correspondências**.  
  O linter examina `doc.Body.Content` nas linhas 259–768, inspecionando exclusivamente:
  - Tabelas (linhas 261–309): células vazias e localização após referências;
  - Elementos de texto `TextRun` (linhas 324–458): cor de fundo (`BackgroundColor` / marca-texto), cor de frente (`ForegroundColor` != preto), família da fonte (`WeightedFontFamily`), tamanho da fonte (`FontSize.Magnitude` com 10pt ou 12pt);
  - Ordem estrutural de seções (linhas 464–550);
  - Vazamentos de sintaxe LaTeX e matemática (linhas 552–610);
  - Espaços duplos e antes de pontuação (linhas 612–651);
  - Qualidade de parágrafos (linhas 653–727): parágrafos < 100 ou > 1500 caracteres, listas com > 4 bullets;
  - Parágrafos vazios consecutivos no corpo (linhas 733–744);
  - Regras de capa e folha de rosto (linhas 771–1093).

  **Conclusão no Linter**: O linter **nunca verifica** `p.ParagraphStyle.LineSpacing`. Não existe limiar, tolerância ou comparação. Se o documento contiver parágrafos inteiros formatados com entrelinhas padrão do Google Docs (1,15 / `LineSpacing: 115`) ou simples (1,0 / `LineSpacing: 100`), o linter reporta que o documento está em plena conformidade.

- **Auditoria de `SpaceAbove` e `SpaceBelow`**:  
  Assim como ocorre com `LineSpacing`, os campos `p.ParagraphStyle.SpaceAbove` e `p.ParagraphStyle.SpaceBelow` **não são lidos nem auditados pelo linter**.  
  A ABNT NBR 14724 determina que o texto normal deve ter espaçamento 1,5 entrelinhas, com **espaçamento antes e depois nulo (`spaceAbove = 0 pt` e `spaceBelow = 0 pt`)**, sem linhas em branco extras entre parágrafos consecutivos. O linter falha em identificar documentos onde o usuário ou template introduziu espaçamentos de 6 pt ou 8 pt entre parágrafos.

#### 1.2 Análise do Formatador (`internal/gdocs/abnt_writer.go`)
- **Como `abnt_writer.go` aplica 150%**:  
  A constante `ABNTLineSpacing = 150` (definida em `internal/gdocs/builder.go:17`) é aplicada no Caso J (Corpo do texto, linhas 1260–1289):
  ```go
  // internal/gdocs/abnt_writer.go:1260-1277
  // ── CASO J: CORPO PADRÃO DE TEXTO (NORMAL_TEXT) ──
  formattedParas++
  requests = append(requests, &docs.Request{
      UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
          Range: pRange,
          ParagraphStyle: &docs.ParagraphStyle{
              NamedStyleType:      "NORMAL_TEXT",
              Alignment:           "JUSTIFIED",
              LineSpacing:         ABNTLineSpacing, // 150
              IndentFirstLine:     &docs.Dimension{Magnitude: ABNTIndentFirstLinePt, Unit: "PT"}, // 1.25 cm
              IndentStart:         &docs.Dimension{Magnitude: 0, Unit: "PT"},
              SpaceAbove:          &docs.Dimension{Magnitude: 0, Unit: "PT"},
              SpaceBelow:          &docs.Dimension{Magnitude: 0, Unit: "PT"},
              AvoidWidowAndOrphan: true,
          },
          Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow,avoidWidowAndOrphan",
      },
  })
  ```
  A máscara de campos (`fields`) está sintaticamente correta perante a Google Docs API v1.

- **Por que parágrafos escapam da formatação (parágrafos pulados)**:
  1. **Tabelas Nativas Blindadas sem Formatação Interna**:
     Na linha 700:
     ```go
     // Tabelas nativas são blindadas: não aplica recuo nas células
     if el.Table != nil {
         continue
     }
     ```
     O comando `continue` ignora todo o elemento `Table`. As células de tabelas contêm seus próprios parágrafos (`cell.Content[*].Paragraph`). Esses parágrafos mantêm qualquer espaçamento pré-existente (ex.: 1,15), violando a ABNT NBR 14724, que prescreve **espaçamento simples (100%) e fonte 10pt** para o interior de tabelas e ilustrações.
  2. **Parágrafos Vazios do Corpo Pulados**:
     Na linha 746:
     ```go
     if trimmed == "" {
         if beforeFirstChapter {
             // Formatação da Capa
         }
         continue // linha 781
     }
     ```
     Parágrafos vazios pós-capa são imediatamente pulados sem atualização de `LineSpacing` ou `SpaceAbove`/`SpaceBelow`.
  3. **Heurística Frágil de Detecção de Citações Longas**:
     Na linha 819:
     ```go
     isLongQuote := existingIndent >= 90.0 && !isHeading1Num && !isHeading1NoNum && !isHeading2 && !isPreTextualNote(trimmed)
     ```
     - Se um parágrafo que deveria ser citação longa (recuo de 4,0 cm = 113,386 pt) foi digitado pelo usuário sem recuo prévio ou com recuo de primeira linha apenas, `existingIndent` é `< 90.0 pt`. Ele é classificado erroneamente como Corpo do Texto e recebe 150% em vez de 100%.
     - Inversamente, se um parágrafo de texto corrido foi recuado manualmente em mais de 90 pt, ele é rebaixado para citação longa, recebendo fonte 10pt e entrelinhas simples.
  4. **Fatiamento por Range (`startIndex` e `endIndex`)**:
     Quando ferramentas MCP como `update_block_by_index` ou `format_document_abnt` recebem fatias parciais, apenas os parágrafos contidos na fatia são atualizados, deixando o restante do documento com 1,15.

---

### Questão 2: Quebras de Linha Espúrias após a Folha de Rosto

#### 2.1 O que `cover_generator.go` faz ao gerar a Capa e a Folha de Rosto
O arquivo `internal/gdocs/cover_generator.go` emula comandos do LaTeX (`\vfill`, `\vspace`) injetando dezenas de quebras de linha (`\n`) repetidas diretamente no texto inserido:
- **Capa** (linhas 308–325):
  - Início: `\n\n` (2 quebras)
  - Pós-instituição: `\n\n\n` (3 quebras)
  - Pós-autor: `\n\n\n` (3 quebras)
  - Pós-título: `strings.Repeat("\n", 18)` (18 quebras de linha consecutivas para empurrar Local/Ano para o rodapé)
- **Folha de Rosto** (linhas 334–356):
  - Início: `\n\n` (2 quebras)
  - Entre Autor e Título: `strings.Repeat("\n", 7)` (7 quebras regulamentares)
  - Entre Título e Preâmbulo: `strings.Repeat("\n", 5)` (5 quebras regulamentares)
  - Entre Preâmbulo e Local/Ano: `strings.Repeat("\n", 14)` (14 quebras para calibração de rodapé)
- **Concatenação e Quebra de Página**:
  Na linha 386, concatena:
  `fullPreText.WriteString(capaText) + "\n" + folhaRostoText`
  E insere quebras de página explícitas via `InsertPageBreakRequest`:
  - Após a Capa (`posCapaBreak`)
  - Após a Folha de Rosto (`posFRBreak`)

#### 2.2 Por que as Quebras Espúrias não foram Saneadas?
Existem 3 causas interdependentes no código que impedem o saneamento:

1. **A regra do `consecutiveBlanks >= 2` em `HealDocTextAndBlanks`**:
   Em `internal/gdocs/abnt_writer.go:343-356`:
   ```go
   t := strings.TrimSpace(txt)
   if t == "" {
       consecutiveBlanks++
       if consecutiveBlanks >= 2 {
           toDeleteRanges = append(toDeleteRanges, &docs.Range{
               StartIndex: el.StartIndex,
               EndIndex:   el.EndIndex,
           })
       }
   } else {
       consecutiveBlanks = 0
   }
   ```
   **Mecanismo**: A primeira linha em branco após a quebra de página da Folha de Rosto tem `consecutiveBlanks == 1`. Portanto, ela **nunca é adicionada a `toDeleteRanges`**. O algoritmo foi desenhado para colapsar múltiplas quebras para uma quebra simples, mas na transição entre páginas pré-textuais e o primeiro título ("RESUMO" ou "1 INTRODUÇÃO"), **uma quebra de página já inicia uma nova página limpa no topo**; a presença de qualquer `\n` antes do título empurra o título para baixo, criando uma linha em branco indesejada no cabeçalho da página.

2. **Corte Frágil de `targetStartIndex` por `anoCount == 2`**:
   Em `internal/gdocs/abnt_writer.go:191-196`:
   ```go
   trimmed := strings.TrimSpace(pText)
   if trimmed == "2026" {
       anoCount++
       if anoCount == 2 {
           targetStartIndex = el.EndIndex
       }
   }
   ```
   - O código utiliza a string exata `"2026"` hardcoded. Se o documento contiver `"Lagarto -- SE / 2026"` ou outro ano, `anoCount` não atinge 2.
   - Quando `targetStartIndex = el.EndIndex`, ele marca o fim do parágrafo do Ano da Folha de Rosto.
   - Se a quebra de página (`PageBreak`) estiver no parágrafo do Ano ou imediatamente no início do próximo elemento, os parágrafos vazios residuais entre o Ano e o Resumo começam a contagem de `consecutiveBlanks` do zero, permitindo que a primeira quebra passe intacta.

3. **Falha na Rotina de Limpeza Reversa em `cover_generator.go`**:
   Em `internal/gdocs/cover_generator.go:891-912`:
   ```go
   // Detecta transição da Folha de Rosto para o RESUMO
   if upperT == "RESUMO" || upperT == "# RESUMO" || strings.HasPrefix(upperT, "RESUMO") {
       inResumeTransition = true
       continue
   }
   if inResumeTransition {
       if t == "" {
           cleanReqs = append(cleanReqs, &docs.Request{
               DeleteContentRange: &docs.DeleteContentRangeRequest{
                   Range: &docs.Range{StartIndex: el.StartIndex, EndIndex: el.EndIndex},
               },
           })
       } else {
           inResumeTransition = false
       }
   }
   ```
   - O loop varre de trás para frente (`for i := len(doc.Body.Content) - 1; i >= 0; i--`).
   - Ao encontrar `"RESUMO"`, ativa `inResumeTransition = true`.
   - O elemento imediatamente anterior pode ser o parágrafo que abriga a `PageBreak` da Folha de Rosto (`pEl.PageBreak != nil` e `TextRun == nil`).
   - Como `t == ""`, esse elemento é incluído em `cleanReqs` para ser deletado!
   - Deletar o elemento que contém a `PageBreak` causaria a perda da quebra de página da Folha de Rosto; além disso, a Google Docs API rejeita a exclusão do range se houver conflito de índice de quebra de seção ou início de documento.
   - Na linha 911: `_, _ = docsSvc.Documents.BatchUpdate(documentID, &docs.BatchUpdateDocumentRequest{Requests: cleanReqs}).Do()`. **O erro é completamente silenciado**, e nenhuma linha espúria é removida.

---

### Questão 3: Carências em `read_doc_content` e `inspect_doc_layout`

#### 3.1 O que `read_doc_content` (`reader.go`) retorna hoje vs o que falta
- **Comportamento Atual**:
  - Implementado em `internal/gdocs/reader.go:217-338`.
  - Processa os parágrafos iterando em `paragraph.Elements` e filtrando estritamente `if elem.TextRun == nil { continue }` (linha 280).
  - Mapeia estilos nomeados para cabeçalhos Markdown (`# `, `## `, `### `).
  - Aplica tags `**` (negrito) e `*` (itálico) em `TextRun`.
  - Converte tabelas para formato texto plano `[TABELA] col | col [/TABELA]`.
- **Limitações e Carências**:
  1. **Invisibilidade das Quebras de Página**: Quando `elem.PageBreak != nil`, `elem.TextRun` é `nil`. O código faz `continue`. O retorno Markdown não indica onde ocorrem quebras de página. Uma IA lendo o texto não consegue saber se "RESUMO" está em uma nova página ou na mesma folha do preâmbulo.
  2. **Invisibilidade de Metadados de Parágrafo**: Nenhum dado de `LineSpacing` (1,15 vs 1,5), `IndentFirstLine` (1,25 cm), `IndentStart` (4,0 cm ou 8,0 cm), `SpaceAbove` ou `SpaceBelow` é retornado.
  3. **Invisibilidade de Objetos Inline**: `InlineObjectElement` (logo do IFS ou ilustrações) é ignorado.
  4. **Ausência de Offsets Estruturais**: O texto retornado não indica os índices de início e fim (`startIndex`/`endIndex`) dos parágrafos, exigindo chamadas adicionais a `search_in_doc` para qualquer edição cirúrgica.

#### 3.2 O que `inspect_doc_layout` (`inspector.go`) retorna hoje vs o que falta
- **Comportamento Atual**:
  - Implementado em `internal/gdocs/inspector.go:44-175`.
  - Popula a struct `ElementLayoutDetail`:
    ```go
    type ElementLayoutDetail struct {
        Index             int     `json:"index"`
        StartIndex        int64   `json:"start_index"`
        EndIndex          int64   `json:"end_index"`
        Type              string  `json:"type"`
        StyleName         string  `json:"style_name,omitempty"`
        Alignment         string  `json:"alignment,omitempty"`
        IndentStartPt     float64 `json:"indent_start_pt,omitempty"`
        IndentFirstLinePt float64 `json:"indent_first_line_pt,omitempty"`
        SpaceAbovePt      float64 `json:"space_above_pt,omitempty"`
        SpaceBelowPt      float64 `json:"space_below_pt,omitempty"`
        FontSizePt        float64 `json:"font_size_pt,omitempty"`
        FontFamily        string  `json:"font_family,omitempty"`
        Bold              bool    `json:"bold,omitempty"`
        Italic            bool    `json:"italic,omitempty"`
        HasPageBreak      bool    `json:"has_page_break,omitempty"`
        HasInlineObject   bool    `json:"has_inline_object,omitempty"`
        TextSnippet       string  `json:"text_snippet,omitempty"`
        TableRows         int     `json:"table_rows,omitempty"`
        TableCols         int     `json:"table_cols,omitempty"`
        EmptyCells        int     `json:"empty_cells,omitempty"`
    }
    ```
- **Limitações e Carências**:
  1. **OMISSÃO CRÍTICA: Ausência do campo `LineSpacing`**:
     Nas linhas 117–132 de `inspector.go`:
     ```go
     if p.ParagraphStyle != nil {
         detail.StyleName = p.ParagraphStyle.NamedStyleType
         detail.Alignment = p.ParagraphStyle.Alignment
         if p.ParagraphStyle.IndentStart != nil { detail.IndentStartPt = p.ParagraphStyle.IndentStart.Magnitude }
         if p.ParagraphStyle.IndentFirstLine != nil { detail.IndentFirstLinePt = p.ParagraphStyle.IndentFirstLine.Magnitude }
         if p.ParagraphStyle.SpaceAbove != nil { detail.SpaceAbovePt = p.ParagraphStyle.SpaceAbove.Magnitude }
         if p.ParagraphStyle.SpaceBelow != nil { detail.SpaceBelowPt = p.ParagraphStyle.SpaceBelow.Magnitude }
     }
     ```
     O campo `p.ParagraphStyle.LineSpacing` **NÃO É EXTRAÍDO**! Ele sequer existe na struct `ElementLayoutDetail`. Portanto, uma IA ou cliente que chame `inspect_doc_layout` recebe recuos e tamanhos de fonte, mas fica 100% cega quanto ao espaçamento entrelinhas (100 vs 115 vs 150).
  2. **Ausência de Metadados Globais do Documento e Seções**:
     Não expõe `DocumentStyle` (margens gerais: topo 3cm, esquerda 3cm, rodapé 2cm, direita 2cm; dimensões A4: 595,28 x 841,89 pt).
  3. **Ausência de Propriedades de Seção (`SectionStyle`)**:
     Não expõe `SectionBreak`, `PageNumbering` (estilo e numeração de página), `HeaderId`, `FooterId`, `UseFirstPageHeaderFooter`. Isso impede diagnosticar a exigência ABNT de numeração no canto superior direito apenas a partir da Introdução.
  4. **Ausência de Flags Tipográficas ABNT**:
     `KeepWithNext` (obrigatório em títulos e legendas para não deixar cabeçalhos órfãos no fim da página) e `AvoidWidowAndOrphan` não são extraídos.

#### 3.3 Especificação Exata dos Metadados a Serem Adicionados aos Payloads JSON

##### Payload Enriquecido de `inspect_doc_layout`:
```json
{
  "document_id": "1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8",
  "document_title": "TCC Teste IA",
  "total_elements": 541,
  "document_geometry": {
    "page_width_pt": 595.28,
    "page_height_pt": 841.89,
    "margin_top_cm": 3.0,
    "margin_bottom_cm": 2.0,
    "margin_left_cm": 3.0,
    "margin_right_cm": 2.0
  },
  "sections": [
    {
      "section_index": 0,
      "start_index": 0,
      "end_index": 1250,
      "has_header": false,
      "has_footer": false,
      "page_numbering_start": null,
      "use_first_page_header_footer": true
    }
  ],
  "elements": [
    {
      "index": 47,
      "start_index": 336,
      "end_index": 432,
      "type": "PARAGRAPH",
      "style_name": "NORMAL_TEXT",
      "alignment": "CENTER",
      "line_spacing_percent": 150.0,
      "indent_start_cm": 0.0,
      "indent_first_line_cm": 0.0,
      "space_above_pt": 6.0,
      "space_below_pt": 6.0,
      "keep_with_next": false,
      "avoid_widow_and_orphan": true,
      "font_size_pt": 12.0,
      "font_family": "Arial",
      "bold": true,
      "italic": false,
      "text_snippet": "AVALIAÇÃO DE UMA ARQUITETURA HÍBRIDA RAG-MCP..."
    }
  ]
}
```

##### Payload Enriquecido de `read_doc_content`:
Adicionar parâmetro opcional `include_metadata: bool` (default: false). Quando ativado:
1. Emissão explícita de marcadores:
   - `\n[PAGE_BREAK (Página %d)]\n`
   - `\n[SECTION_BREAK (Seção %d)]\n`
   - `[INLINE_OBJECT: Logo IFS (3.2cm centralizada)]`
2. Anotação estrutural de cabeçalho de bloco com metadados ABNT:
   ```markdown
   <!-- block_start: 336-432, style: NORMAL_TEXT, line_spacing: 1.5, align: CENTER -->
   **AVALIAÇÃO DE UMA ARQUITETURA HÍBRIDA RAG-MCP NA VALIDAÇÃO NORMATIVA DE REQUERIMENTOS ACADÊMICOS**
   ```

---

## 3. Síntese Comparativa dos Módulos Auditados

| Módulo | Componente Analisado | Comportamento Esperado (ABNT) | Comportamento Atual no Código | Causa Raiz Técnica |
| :--- | :--- | :--- | :--- | :--- |
| **`linter.go`** | Auditoria de Entrelinhas | Detectar desvios de 1,5 no corpo e 1,0 em citações/tabelas | Não valida entrelinhas | Ausência completa de leitura de `LineSpacing` |
| **`linter.go`** | Espaçamento Vertical | Validar `SpaceAbove=0` e `SpaceBelow=0` no corpo | Não valida espaçamento vertical | Ausência completa de leitura de `SpaceAbove` e `SpaceBelow` |
| **`abnt_writer.go`** | Formatação de Tabelas | Aplicar entrelinhas simples (1,0) e 10pt nas células | Pula tabelas inteiras | `if el.Table != nil { continue }` na linha 701 |
| **`abnt_writer.go`** | Citação Longa | Recuo de 4cm, 10pt e 1,0 simples | Heurística baseada em recuo prévio >= 90pt | Se o texto não estiver pré-recuado, recebe 1,5 e 12pt |
| **`abnt_writer.go`** | Remoção de Linhas Vazias | Eliminar quebras espúrias em transições de página | Mantém a primeira quebra vazia | Limiar rígido `consecutiveBlanks >= 2` em `HealDocTextAndBlanks` |
| **`cover_generator.go`** | Quebras da Folha de Rosto | Calibração de rodapé sem quebras residuais pós-página | Insere dezenas de `\n` e falha ao limpar | Concatena `\n` redundantes e silencia erro na limpeza reversa |
| **`inspector.go`** | Inspeção Estrutural | Expor geometria completa para diagnóstico por IA | Omite `LineSpacing` e margens de página | Struct `ElementLayoutDetail` não contém o campo `LineSpacing` |
| **`reader.go`** | Leitura de Conteúdo | Expor estrutura fiel incluindo quebras de página | Markdown simplificado sem quebras de página | Ignora `elem.PageBreak` porque `elem.TextRun == nil` |

---

## 4. Recomendações Técnicas para o Plano de Engenharia (R3)

1. **Atualizar `inspector.go`**:
   - Incluir `LineSpacing float64` na struct `ElementLayoutDetail` e extrair `p.ParagraphStyle.LineSpacing`.
   - Incluir propriedades de `DocumentStyle` (margens ABNT 3/3/2/2 cm) e `SectionStyle` no relatório.
2. **Atualizar `linter.go`**:
   - Criar regra `LineSpacingAnomalyItem` comparando `LineSpacing` de cada parágrafo com o esperado pelo contexto (150 para corpo e títulos; 100 para notas, citações, tabelas, legendas e referências).
   - Criar regra `ParagraphMarginAnomalyItem` para sinalizar `SpaceAbove > 0` ou `SpaceBelow > 0` em parágrafos de corpo normal.
3. **Refatorar `abnt_writer.go`**:
   - Iterar sobre as células de tabelas (`el.Table.TableRows[*].TableCells[*]`) e aplicar `LineSpacing: 100` e fonte 10pt.
   - Refatorar `HealDocTextAndBlanks`: após uma quebra de página (`PageBreak`), qualquer parágrafo vazio imediatamente anterior ao próximo título textual deve ser removido (`consecutiveBlanks >= 1`).
4. **Refatorar `cover_generator.go`**:
   - Substituir quebras múltiplas de `\n` para calibração de rodapé por espaçamento vertical de parágrafo (`SpaceAbove`) ou quebra de página direta, evitando acúmulo de parágrafos vazios.
   - Corrigir o loop de limpeza reversa para não tentar deletar o parágrafo da `PageBreak`.
5. **Enriquecer `reader.go`**:
   - Inserir marcador textual `[PAGE_BREAK]` ao encontrar `elem.PageBreak != nil`.
   - Adicionar modo com anotações de estilo para consumo por modelos de IA.
