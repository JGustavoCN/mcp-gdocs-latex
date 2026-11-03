# Relatório de Handoff: Investigação Técnica de Quebras de Linha Espúrias e Títulos de Seção

**Data**: 02 de Novembro de 2026  
**Investigador**: Explorer 1 (`teamwork_explorer`)  
**Documento Auditado**: Google Docs ID `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` ("TCC Teste IA")  
**Arquivos de Código Analisados**: `internal/gdocs/linter.go`, `internal/gdocs/abnt_writer.go`, `internal/gdocs/inspector.go`, `internal/gdocs/reader.go`, `internal/gdocs/cover_generator.go`, `internal/gdocs/writer.go`, `internal/mcp/tools.go`.

---

## 1. Observation (Observações Diretas e Evidências Concretas)

### 1.1 Discrepância de Versão e Binário do Servidor MCP
- O arquivo binário do servidor MCP em execução (`bin\mcp-server-v2.exe`) registrado no Antigravity (`C:\Users\joseg\.gemini\antigravity\mcp_config.json:53`) possui timestamp de compilação `02/11/2026 18:28:44`.
- Os arquivos-fonte `internal/gdocs/linter.go` e `internal/gdocs/abnt_writer.go` possuem modificações não comitadas (*unstaged*) salvas às `18:54:05` e `18:53:31`.
- A chamada de `lint_document_abnt` executada pelo servidor MCP ativo retornou **falso negativo completo**:
  ```json
  {
    "total_elements": 537,
    "total_tables": 7,
    "cover_status": "Capa Oficial IFS com Logo (3.2cm centralizada) e Folha de Rosto NBR 14724 em plena conformidade.",
    "margins_conformity": { "is_conform": true },
    "message": "Auditoria concluída. 0 anomalias/destaques identificados no documento (autoFix=false)."
  }
  ```

### 1.2 Radiografia Geométrica no Documento Real (`inspect_doc_layout`)
Ao inspecionar diretamente os `StructuralElements` na seção de Metodologia (faixa de índices 59122 a 63300), constatamos a presença maciça de **parágrafos vazios espúrios intercalados e estilizados anomalamente como `HEADING_1`**:

```json
[
  { "index": 217, "start_index": 59122, "end_index": 59136, "type": "PARAGRAPH", "style_name": "HEADING_1", "text_snippet": "3 METODOLOGIA" },
  { "index": 218, "start_index": 59136, "end_index": 59137, "type": "PARAGRAPH", "style_name": "HEADING_1", "space_above_pt": 18, "space_below_pt": 12, "bold": true },
  { "index": 219, "start_index": 59137, "end_index": 59510, "type": "PARAGRAPH", "style_name": "NORMAL_TEXT", "text_snippet": "Este capítulo apresenta os procedimentos metodológicos..." },
  { "index": 220, "start_index": 59510, "end_index": 59511, "type": "PARAGRAPH", "style_name": "HEADING_1", "space_above_pt": 18, "space_below_pt": 12, "bold": true },
  { "index": 221, "start_index": 59511, "end_index": 59908, "type": "PARAGRAPH", "style_name": "NORMAL_TEXT", "text_snippet": "Para assegurar clareza e rigor metodológico..." },
  { "index": 222, "start_index": 59908, "end_index": 59990, "type": "PARAGRAPH", "style_name": "NORMAL_TEXT", "text_snippet": "Tabela 3 – Distinção entre a solução tecnológica ampla e o objeto de investigação" },
  { "index": 223, "start_index": 59990, "end_index": 59991, "type": "PARAGRAPH", "style_name": "HEADING_1", "space_above_pt": 18, "space_below_pt": 12, "bold": true },
  { "index": 224, "start_index": 59991, "end_index": 60764, "type": "TABLE", "text_snippet": "Tabela Nativa 5x3 (vazias: 0/15)" },
  { "index": 225, "start_index": 60764, "end_index": 60765, "type": "PARAGRAPH", "style_name": "NORMAL_TEXT" },
  { "index": 226, "start_index": 60765, "end_index": 60801, "type": "PARAGRAPH", "style_name": "NORMAL_TEXT", "text_snippet": "Fonte: Elaborado pelo autor (2026)." },
  { "index": 227, "start_index": 60801, "end_index": 60802, "type": "PARAGRAPH", "style_name": "HEADING_1", "space_above_pt": 18, "space_below_pt": 12, "bold": true },
  { "index": 228, "start_index": 60802, "end_index": 60833, "type": "PARAGRAPH", "style_name": "HEADING_1", "text_snippet": "3.1 CARACTERIZAÇÃO DA PESQUISA" },
  { "index": 229, "start_index": 60833, "end_index": 60834, "type": "PARAGRAPH", "style_name": "HEADING_1", "space_above_pt": 18, "space_below_pt": 12, "bold": true },
  { "index": 230, "start_index": 60834, "end_index": 61724, "type": "PARAGRAPH", "style_name": "NORMAL_TEXT", "text_snippet": "A partir dessa delimitação epistemológica, o capítulo organiza-se em dez seções articuladas: a..." },
  { "index": 231, "start_index": 61724, "end_index": 61725, "type": "PARAGRAPH", "style_name": "HEADING_1", "space_above_pt": 18, "space_below_pt": 12, "bold": true },
  { "index": 232, "start_index": 61725, "end_index": 61756, "type": "PARAGRAPH", "style_name": "HEADING_1", "text_snippet": "3.1 Caracterização da Pesquisa" },
  { "index": 233, "start_index": 61756, "end_index": 61757, "type": "PARAGRAPH", "style_name": "HEADING_1", "space_above_pt": 18, "space_below_pt": 12, "bold": true },
  { "index": 234, "start_index": 61757, "end_index": 62068, "type": "PARAGRAPH", "style_name": "NORMAL_TEXT", "text_snippet": "A fundamentação do método científico..." },
  { "index": 235, "start_index": 62068, "end_index": 62069, "type": "PARAGRAPH", "style_name": "HEADING_1", "space_above_pt": 18, "space_below_pt": 12, "bold": true },
  { "index": 236, "start_index": 62069, "end_index": 62629, "type": "PARAGRAPH", "style_name": "NORMAL_TEXT", "text_snippet": "Quanto à Natureza: Classifica-se como Pesquisa Aplicada..." },
  { "index": 237, "start_index": 62629, "end_index": 62630, "type": "PARAGRAPH", "style_name": "HEADING_1", "space_above_pt": 18, "space_below_pt": 12, "bold": true },
  { "index": 238, "start_index": 62630, "end_index": 63295, "type": "PARAGRAPH", "style_name": "NORMAL_TEXT", "text_snippet": "Quanto aos Objetivos: O trabalho assume caráter exploratório..." },
  { "index": 239, "start_index": 63295, "end_index": 63296, "type": "PARAGRAPH", "style_name": "HEADING_1", "space_above_pt": 18, "space_below_pt": 12, "bold": true }
]
```

O mesmo padrão repete-se na Introdução (elementos 76 a 80):
- Elemento 76: `1 INTRODUÇÃO` (`HEADING_1`)
- Elemento 77: Parágrafo vazio `[2974-2975]` (`NORMAL_TEXT`)
- Elemento 78: `Contextualização` (`HEADING_2`)
- Elemento 79: Parágrafo vazio `[2992-2993]` (`NORMAL_TEXT`)
- Elemento 80: Texto do parágrafo `[2993-3645]`

### 1.3 Código Atual de Auditoria no Linter (`internal/gdocs/linter.go`)
No commit base (`bin\mcp-server-v2.exe`), as linhas 945-950 de `linter.go` continham:
```go
if trimmed == "" {
    if currentState != statePreTextual {
        consecutiveEmptyParas++
        if consecutiveEmptyParas >= 2 {
            report.ConsecutiveBlankLines = append(report.ConsecutiveBlankLines, ...)
        }
    }
} else {
    consecutiveEmptyParas = 0
}
```

### 1.4 Código Atual de Saneamento no Formatador (`internal/gdocs/abnt_writer.go`)
No commit base, as linhas 350-355 de `abnt_writer.go` em `HealDocTextAndBlanks` continham:
```go
if (consecutiveBlanks >= 2 || isBeforeHeading) && !isBeforeTable {
    toDeleteRanges = append(toDeleteRanges, &docs.Range{
        StartIndex: el.StartIndex,
        EndIndex:   el.EndIndex,
    })
}
```
E nas linhas 366-380 de `abnt_writer.go`:
```go
for i := len(toDeleteRanges) - 1; i >= 0; i-- {
    r := toDeleteRanges[i]
    req := &docs.BatchUpdateDocumentRequest{
        Requests: []*docs.Request{
            { DeleteContentRange: &docs.DeleteContentRangeRequest{ Range: r } },
        },
    }
    svc.Documents.BatchUpdate(documentID, req).Do() // UMA REQUISIÇÃO HTTP INDIVIDUAL POR PARÁGRAFO!
}
```
E no loop de formatação de parágrafos de `FormatDocumentABNTExtended` (linhas 786-822):
```go
if trimmed == "" {
    if beforeFirstChapter { ... }
    continue // IGNORA COMPLETAMENTE PARÁGRAFOS VAZIOS NA PARTE TEXTUAL!
}
```

### 1.5 Classificação Incorreta de Títulos Secundários (`internal/gdocs/abnt_writer.go:825`)
```go
isHeading1Num := rePrimaryHeadingNum.MatchString(trimmed) || (currentStyle == "HEADING_1" && !rePrimaryHeadingNoNum.MatchString(trimmed) && len(trimmed) < 120)
isHeading2 := reSecondaryHeading.MatchString(trimmed) || (currentStyle == "HEADING_2" && len(trimmed) < 120)
```
Como a verificação de `isHeading1Num` precede `isHeading2` e possui a guarda permissiva `currentStyle == "HEADING_1"`, qualquer subtítulo secundário ("3.1 CARACTERIZAÇÃO DA PESQUISA" ou "3.1 Caracterização da Pesquisa") que já estivesse com o estilo `HEADING_1` permanece classificado como `HEADING_1`, recebendo caixa alta forçada, negrito, `SpaceAbove: 18pt` e `SpaceBelow: 12pt`, violando a ABNT NBR 6024.

---

## 2. Logic Chain (Cadeia Lógica de Causa Raiz e Deduções Técnicas)

1. **Causa Raiz 1 — Falso Negativo no Linter para Parágrafos Vazios Isolados:**
   - *Premissa A:* No documento real, as quebras de linha em excesso não ocorrem em blocos contíguos de 2 ou mais `\n`, mas sim de forma **intercalada** (1 parágrafo de texto, 1 parágrafo vazio, 1 parágrafo de texto, etc.) ou imediatamente após títulos.
   - *Premissa B:* A cada parágrafo de texto, o linter executa `consecutiveEmptyParas = 0`.
   - *Premissa C:* Ao encontrar o parágrafo vazio isolado seguinte, o contador atinge `consecutiveEmptyParas = 1`.
   - *Conclusão:* Como a lógica original exigia `if consecutiveEmptyParas >= 2`, a condição NUNCA foi satisfeita. O linter reportou 0 anomalias de linhas em branco espúrias, gerando um falso negativo completo no checklist normativo ABNT NBR 14724.

2. **Causa Raiz 2 — Falha de Saneamento no `format_document_abnt` (`HealDocTextAndBlanks`):**
   - *Premissa A:* O saneador `HealDocTextAndBlanks` utilizava a mesma guarda restritiva `(consecutiveBlanks >= 2 || isBeforeHeading)`.
   - *Premissa B:* As quebras espúrias após títulos (ex.: logo após `3 METODOLOGIA` ou `3.1 CARACTERIZAÇÃO DA PESQUISA`) ocorrem **DEPOIS** do título (`isBeforeHeading == false`) e possuem `consecutiveBlanks == 1`.
   - *Conclusão:* Essas quebras nunca foram incluídas na lista de deleção `toDeleteRanges`. Elas permaneceram intocadas no documento.

3. **Causa Raiz 3 — Herança de Estilo `HEADING_1` em Parágrafos Vazios com Espaçamento Gigantesco:**
   - *Premissa A:* No Google Docs, quando o autor aperta ENTER após um título, o novo parágrafo herda o `NamedStyleType: HEADING_1`.
   - *Premissa B:* O estilo `HEADING_1` ABNT aplica `SpaceAbove: 18pt`, `SpaceBelow: 12pt`, e `LineSpacing: 150%` (18pt adicionais).
   - *Premissa C:* No loop principal de formatação (`FormatDocumentABNTExtended`), ao encontrar `trimmed == ""` após o primeiro capítulo, o código executa `continue` sem limpar nem redefinir o estilo do elemento.
   - *Conclusão:* Os parágrafos vazios espúrios mantiveram o estilo `HEADING_1`, gerando um vão em branco artificial de aproximadamente **48 pontos tipográficos (~1,7 cm)** entre títulos e parágrafos de texto, distorcendo gravemente a diagramação do trabalho.

4. **Causa Raiz 4 — Precedência Invertida na Detecção Hierárquica de Títulos:**
   - *Premissa A:* O título de seção secundária `3.1 CARACTERIZAÇÃO DA PESQUISA` casa perfeitamente com a regex `reSecondaryHeading`.
   - *Premissa B:* Porém, a verificação de `isHeading1Num` testava primeiro `currentStyle == "HEADING_1" && len(trimmed) < 120`.
   - *Conclusão:* O título secundário foi mantido como `HEADING_1`, perpetuando o erro de hierarquia ABNT NBR 6024 e impedindo a aplicação dos estilos canônicos de `HEADING_2` (alinhado à esquerda, sem negrito, `SpaceAbove: 12pt`, `SpaceBelow: 6pt`).

5. **Causa Raiz 5 — Risco Crítico de Regressão e Deleção de Quebras de Página (`PageBreak`):**
   - *Premissa A:* No Google Docs, um parágrafo que contém uma quebra de página (`PageBreak`) possui texto `txt == ""` (não há `TextRun`, ou há apenas `\n`).
   - *Premissa B:* Se `HealDocTextAndBlanks` simplesmente testar `if strings.TrimSpace(txt) == ""` sem inspecionar os elementos internos (`pel.PageBreak != nil` e `pel.InlineObjectElement != nil`), a quebra de página da Folha de Rosto ou entre seções será considerada "parágrafo vazio" e adicionada para deleção.
   - *Conclusão:* Isso destruiria a separação entre a Folha de Rosto e o Resumo (colapsando a página 2 na página 3), disparando erro crítico no linter pré-textual. A proteção estrutural explícita de `PageBreak` e `InlineObjectElement` é indispensável.

6. **Causa Raiz 6 — Ineficiência e Risco de Quota da Deleção 1 a 1 via HTTP:**
   - *Premissa A:* Executar uma requisição `batchUpdate` HTTP separada para cada parágrafo vazio a deletar faz com que um documento com dezenas de quebras dispare dezenas de chamadas HTTP sequenciais.
   - *Premissa B:* A API do Google Docs possui cota rígida de requisições por minuto por projeto (300 req/min).
   - *Premissa C:* O protocolo `documents.batchUpdate` suporta múltiplos `DeleteContentRangeRequest` na mesma requisição, processados em ordem sequencial.
   - *Conclusão:* Ordenando os ranges em ordem decrescente de `StartIndex` (do maior offset para o menor), a deleção é **100% segura, matematicamente atômica e executada em uma única chamada HTTP**.

---

## 3. Caveats (Pontos de Atenção, Casos Especiais e Restrições)

1. **Parágrafo Delimitador Pré-Tabela Obrigatório (`isBeforeTable`):**
   - A Google Docs API **exige estritamente** que tabelas nativas sejam precedidas por um `StructuralElement` do tipo `Paragraph`. Se tentarmos deletar o parágrafo imediatamente anterior a uma tabela, a API rejeita a operação com erro HTTP 400 (`Invalid requests[X].deleteContentRange: Cannot delete table start boundary`).
   - **Solução necessária**: O parágrafo pré-tabela deve ser **preservado da deleção**, mas seu estilo deve ser cirurgicamente neutralizado: reset para `NORMAL_TEXT`, `lineSpacing: 100`, `fontSize: 1.0pt` (ou mínimo), `spaceAbove: 0` e `spaceBelow: 0`. Isso elimina o espaçamento visual excessivo entre a legenda e o topo da tabela sem violar as restrições da API.

2. **Último Parágrafo Inamovível do Documento (`bodyEndIndex - 1`):**
   - O caractere final do Google Docs (`doc.Body.Content[len-1].EndIndex - 1`) é um newline inamovível mantido pela arquitetura interna do Google Workspace.
   - **Guarda mandatória**: Qualquer range com `EndIndex >= bodyEndIndex` **jamais** pode ser incluído no `DeleteContentRangeRequest`.

3. **Duplicação de Título na Metodologia:**
   - O documento real possui dois elementos com a numeração 3.1:
     - Elemento 228: `3.1 CARACTERIZAÇÃO DA PESQUISA`
     - Elemento 230: Parágrafo de texto explicando a divisão do capítulo em 10 seções.
     - Elemento 232: `3.1 Caracterização da Pesquisa` (início real do conteúdo metodológico).
   - O Worker deve estar ciente de que o elemento 228 foi inserido indevidamente antes do texto de visão geral do capítulo ou trata-se de um título preliminar redundante. O formatador deve normalizar ambos os títulos para `HEADING_2`, e o linter deve auditar a unicidade da numeração de seções.

---

## 4. Conclusion & Engenharia da Solução para o Worker

A solução para sanar as quebras de linha espúrias e desvios de títulos requer alterações coordenadas em **dois arquivos Go**: `internal/gdocs/linter.go` e `internal/gdocs/abnt_writer.go`, seguidas da **recompilação do binário do servidor MCP** (`bin/mcp-server-v2.exe`).

### 4.1 Especificação de Mudanças em `internal/gdocs/linter.go`

1. **Auditoria de Linhas em Branco Espúrias:**
   - No loop principal de elementos (linha 958):
     - Substituir a condição de contagem consecutiva `>= 2` por `>= 1` para toda a parte textual e pós-textual (`currentState != statePreTextual`).
     - Adicionar verificação de proteção: NÃO reportar se o parágrafo contiver `pel.PageBreak != nil` ou `pel.InlineObjectElement != nil`.
     - NÃO reportar se for o parágrafo delimitador pré-tabela (`isBeforeTable == true`).
     - Para todos os demais parágrafos vazios (`trimmed == ""`), registrar em `report.ConsecutiveBlankLines`.
2. **Auditoria de Títulos Vazios (Headings Fantasmas):**
   - Se `isHeading && trimmed == ""`: registrar em `report.StructuralOrderViolations` com severidade `ALERTA`: "Título de seção vazio / Heading sem texto atribuído".
3. **Auditoria de Falsos Títulos terminados em dois pontos (`:`):**
   - Confirmar a regra existente que sinaliza parágrafos com estilo `HEADING_X` terminados em `:`.
4. **Auditoria de Hierarquia de Títulos (NBR 6024):**
   - Se um parágrafo que começa com numeração secundária (ex: `3.1 `) estiver com estilo `HEADING_1`: sinalizar em `report.StructuralOrderViolations`: "Seção secundária formatada indevidamente como título primário (HEADING_1 em vez de HEADING_2)".

### 4.2 Especificação de Mudanças em `internal/gdocs/abnt_writer.go`

1. **Refatoração de `HealDocTextAndBlanks` (Saneamento Atômico em Lote Reverso):**
   - Identificar parágrafos puramente vazios na parte textual (`StartIndex >= targetStartIndex` e `EndIndex < bodyEnd`):
     - Inspecionar os elementos do parágrafo: se possuir `pel.PageBreak != nil` ou `pel.InlineObjectElement != nil`, **PRESERVAR SEMPRE**.
     - Se `isBeforeTable == true`, **PRESERVAR DA DELEÇÃO**.
     - Se `trimmed == ""`, adicionar a `toDeleteRanges`.
   - **Construção do Lote Atômico Reverso:**
     - Ordenar `toDeleteRanges` decrescentemente por `StartIndex`:
       ```go
       sort.Slice(toDeleteRanges, func(i, j int) bool {
           return toDeleteRanges[i].StartIndex > toDeleteRanges[j].StartIndex
       })
       ```
     - Agrupar em lotes de até 100 requisições `DeleteContentRangeRequest` e enviar em transações atômicas via `svc.Documents.BatchUpdate`:
       ```go
       var delRequests []*docs.Request
       for _, r := range toDeleteRanges {
           delRequests = append(delRequests, &docs.Request{
               DeleteContentRange: &docs.DeleteContentRangeRequest{Range: r},
           })
       }
       // Envia em chunks para respeitar limites de payload
       ```
2. **Normalização do Parágrafo Delimitador Pré-Tabela em `FormatDocumentABNTExtended`:**
   - No loop de formatação (linha 786), quando `trimmed == "" && !beforeFirstChapter`:
     - Verificar se `isBeforeTable == true`:
       - Aplicar `UpdateParagraphStyleRequest` configurando `NamedStyleType: "NORMAL_TEXT"`, `Alignment: "START"`, `LineSpacing: 100`, `SpaceAbove: 0`, `SpaceBelow: 0`.
       - Aplicar `UpdateTextStyleRequest` configurando `FontSize: 1.0 pt` (ou `6.0 pt`) para neutralizar o impacto visual.
3. **Correção da Precedência de Títulos em `FormatDocumentABNTExtended`:**
   - As regexes de numeração hierárquica **DEVEM TER PRECEDÊNCIA** sobre o estilo existente `currentStyle`:
     ```go
     // 1. Títulos Primários Canônicos
     isHeading1Num := rePrimaryHeadingNum.MatchString(trimmed)
     isHeading1NoNum := rePrimaryHeadingNoNum.MatchString(trimmed)

     // 2. Títulos Secundários, Terciários e Quaternários por REGEX
     isHeading2 := reSecondaryHeading.MatchString(trimmed)
     isHeading3 := reTertiaryHeading.MatchString(trimmed)
     isHeading4 := reQuaternaryHeading.MatchString(trimmed)

     // 3. Fallback para estilos existentes apenas se NÃO for número de outro nível:
     if !isHeading1Num && !isHeading1NoNum && !isHeading2 && !isHeading3 && !isHeading4 {
         if currentStyle == "HEADING_1" && len(trimmed) < 120 && !reSecondaryHeading.MatchString(trimmed) {
             isHeading1Num = true
         } else if currentStyle == "HEADING_2" && len(trimmed) < 120 {
             isHeading2 = true
         } else if currentStyle == "HEADING_3" && len(trimmed) < 120 {
             isHeading3 = true
         }
     }

     // 4. Orações com ':' nunca são títulos
     if strings.HasSuffix(trimmed, ":") {
         isHeading1Num, isHeading1NoNum, isHeading2, isHeading3, isHeading4 = false, false, false, false, false
     }
     ```

### 4.3 Recompilação do Binário MCP
Após os ajustes no código:
- Executar `go build -o bin/mcp-server-v2.exe ./cmd/mcp-server`
- Verificar se o processo em segundo plano do MCP server é atualizado para refletir o novo binário.

---

## 5. Verification Method (Como Verificar e Provar a Solução)

1. **Testes Unitários Automatizados Go:**
   - Executar a suíte completa de testes:
     ```powershell
     go test -v ./...
     ```
   - Validar que todos os testes existentes (`linter_rules_test.go`, `abnt_heal_test.go`, etc.) passam sem regressão.
   - Adicionar novo teste unitário em `abnt_heal_test.go` verificando que parágrafos com `PageBreak` **nunca** são elegíveis para deleção.

2. **Execução no Documento Real do Google Docs:**
   - Executar `format_document_abnt` no documento `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` via tool MCP.
   - Inspecionar a saída da ferramenta, confirmando o número de quebras eliminadas.

3. **Auditoria com `inspect_doc_layout` na Seção de Metodologia:**
   - Executar `inspect_doc_layout` na faixa de índices 59000 a 63000.
   - Provar que os parágrafos vazios espúrios (elementos 218, 220, 227, 229, 231, 233, 235, 237, 239) deixaram de existir.
   - Provar que a Tabela 3 mantém seu parágrafo pré-tabela delimitador com estilo `NORMAL_TEXT` e espaçamento zerado.
   - Provar que `3.1 CARACTERIZAÇÃO DA PESQUISA` está com `style_name: "HEADING_2"` (sem negrito, 12pt).

4. **Auditoria Oficial com `lint_document_abnt`:**
   - Executar `lint_document_abnt` com `auto_fix: false`.
   - Verificar relatório com **0 anomalias**, confirmando integridade da Capa, Folha de Rosto (com quebra de página intacta), margens 3/3/2/2 cm e ausência de linhas em branco espúrias.
