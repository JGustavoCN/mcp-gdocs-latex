# Relatório de Handoff Técnico (Diagnóstico do Codebase Go - R2)
**Agente Responsável**: `explorer_diag_1`  
**Destinatário**: Agente Orquestrador (`parent` / `8522a030-5a02-4454-bc22-525b48f7de69`)  
**Data**: 2026-11-01  
**Tipo de Handoff**: Hard (Tarefa Concluída)  
**Arquivos de Entrega**:
- `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\explorer_diag_1\analysis.md` (Relatório Técnico Detalhado)
- `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\explorer_diag_1\handoff.md` (Este documento)

---

## 1. Observation (Observações Diretas)

1. **`internal/gdocs/linter.go`**:
   - Uma busca por `LineSpacing`, `spaceAbove` ou `spaceBelow` retorna 0 ocorrências em todo o arquivo.
   - O loop de parágrafos (linhas 312–768) avalia apenas `ts.BackgroundColor` (linha 341), `ts.ForegroundColor` (linha 368), `ts.WeightedFontFamily` (linha 385) e `ts.FontSize` (linha 400).
   - O objeto `p.ParagraphStyle` é utilizado exclusivamente para extrair `NamedStyleType` (linha 316) e `Alignment` na auditoria de capa (linha 809). **Nenhum teste de entrelinhas é executado.**

2. **`internal/gdocs/abnt_writer.go`**:
   - Linha 700–702:
     ```go
     // Tabelas nativas são blindadas: não aplica recuo nas células
     if el.Table != nil {
         continue
     }
     ```
     Todas as tabelas nativas são puladas; parágrafos dentro de células mantêm qualquer entrelinhas existente (ex.: 1,15), em desconformidade com a ABNT NBR 14724 (que exige espaçamento simples 1,0 e fonte 10pt).
   - Linhas 343–356 (`collapseBlankLines` em `HealDocTextAndBlanks`):
     ```go
     if t == "" {
         consecutiveBlanks++
         if consecutiveBlanks >= 2 {
             toDeleteRanges = append(toDeleteRanges, &docs.Range{StartIndex: el.StartIndex, EndIndex: el.EndIndex})
         }
     } else {
         consecutiveBlanks = 0
     }
     ```
     O limiar `>= 2` preserva intencionalmente a primeira quebra vazia, impedindo a limpeza de linhas em branco residuais no topo de páginas pós-quebra de seção/página.
   - Linhas 191–196:
     ```go
     if trimmed == "2026" {
         anoCount++
         if anoCount == 2 {
             targetStartIndex = el.EndIndex
         }
     }
     ```
     Corte de início textual depende rigidamente da string literal `"2026"`.

3. **`internal/gdocs/cover_generator.go`**:
   - Linha 324: Injeta `strings.Repeat("\n", 18)` na Capa.
   - Linhas 337, 339, 350: Injeta 7, 5 e 14 quebras `\n` consecutivas na Folha de Rosto.
   - Linhas 891–912: A rotina de limpeza reversa entre Folha de Rosto e Resumo inclui na lista de deleção parágrafos vazios anteriores a "RESUMO", colidindo com o parágrafo que abriga o `PageBreak`. O erro retornado pela Google Docs API é descartado na linha 911 via `_, _ = docsSvc.Documents.BatchUpdate(...)`.

4. **`internal/gdocs/inspector.go`**:
   - Struct `ElementLayoutDetail` (linhas 10–31): Contém `IndentStartPt`, `IndentFirstLinePt`, `SpaceAbovePt`, `SpaceBelowPt`, `FontSizePt`, `FontFamily`, mas **NÃO possui o campo `LineSpacing`**.
   - Linhas 117–132: O código lê `IndentStart`, `IndentFirstLine`, `SpaceAbove`, `SpaceBelow`, mas **ignora `p.ParagraphStyle.LineSpacing`**.

5. **`internal/gdocs/reader.go`**:
   - Linhas 278–282:
     ```go
     for _, elem := range paragraph.Elements {
         if elem.TextRun == nil {
             continue
         }
     ```
     Quebras de página explícitas (`elem.PageBreak != nil`) são descartadas porque não possuem `TextRun`. O leitor Markdown não exibe quebras de página, recuos ou entrelinhas.

---

## 2. Logic Chain (Cadeia Lógica de Raciocínio)

1. **Premissa 1**: Um linter só pode reportar anomalias de formatação se inspecionar os campos correspondentes da estrutura de dados retornada pela API.
   - **Fato 1**: `linter.go` não contém nenhuma leitura do campo `LineSpacing` nem dos campos `SpaceAbove` e `SpaceBelow`.
   - **Dedução 1**: É computacionalmente impossível que o linter atual sinalize desvios de entrelinhas (como 1,15 vs 1,5) ou espaços extras entre parágrafos, resultando em falsos negativos sistemáticos.

2. **Premissa 2**: O formatador `FormatDocumentABNTExtended` aplica regras seletivamente baseando-se no tipo do elemento e na máquina de estados contextual.
   - **Fato 2**: Células de tabela são filtradas pelo comando `continue` na linha 701; parágrafos vazios pós-capa são pulados na linha 781; citações longas dependem de recuo pré-existente `>= 90 pt`.
   - **Dedução 2**: Qualquer texto em tabela ou citação longa não pré-formatada não recebe a calibração de entrelinhas da ABNT.

3. **Premissa 3**: A remoção de quebras vazias baseada em `consecutiveBlanks >= 2` preserva exatamente um parágrafo vazio a cada ocorrência.
   - **Fato 3**: Após uma `PageBreak`, um novo topo de página já é criado. A preservação de 1 parágrafo vazio resulta visualmente em uma linha em branco no cabeçalho da página 3 (antes de RESUMO). Além disso, a rotina de exclusão em `cover_generator.go` falha silenciosamente ao colidir com o parágrafo da quebra de página.
   - **Dedução 3**: As quebras espúrias após a folha de rosto permanecem no documento devido à combinação de um limiar excessivamente tolerante (`>= 2`) com a falha silenciada da limpeza de transição.

4. **Premissa 4**: Agentes de IA que operam ferramentas MCP só enxergam as propriedades expostas nos payloads JSON ou Markdown retornados.
   - **Fato 4**: `inspect_doc_layout` não expõe `LineSpacing` e `read_doc_content` não exibe `[PAGE_BREAK]`, nem recuos, nem entrelinhas.
   - **Dedução 4**: Modelos de IA são incapazes de diagnosticar ou auditar essas falhas de forma autônoma sem a adição explícita desses campos aos payloads.

---

## 3. Caveats (Ressalvas e Limitações da Investigação)

1. A investigação foi estritamente em modo de leitura (read-only), sem aplicação de patches nos arquivos do repositório nesta etapa.
2. O documento de teste `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` foi inspecionado via ferramentas MCP locais (`read_doc_content` e `inspect_doc_layout`), confirmando a presença de múltiplos parágrafos vazios nas transições de página e a ausência do metadado de entrelinhas no output.
3. Não foram identificadas restrições na Google Docs API que impeçam a leitura ou escrita de `LineSpacing`, `SpaceAbove` ou `SpaceBelow`. A API do Google Docs suporta integralmente `paragraphStyle.lineSpacing` (em porcentagem: 100, 115, 150, 200).

---

## 4. Conclusion (Conclusão Técnica)

As falhas diagnosticadas no marco R2 decorrem de **lacunas explícitas de implementação e modelagem de dados** no codebase Go, e não de limitações da API do Google Docs:
- O **linter precisa ser expandido** com novos checkers para `LineSpacing` (comparando 150% no corpo vs 100% em notas/legendas/tabelas) e `SpaceAbove`/`SpaceBelow` (garantindo 0 pt no corpo).
- O **formatador deve varrer células de tabelas** (aplicando 100% e 10pt) e adotar heurística aprimorada para citações longas.
- O **saneamento de quebras (`HealDocTextAndBlanks`) deve eliminar 100% dos parágrafos vazios em topos de página** imediatamente precedidos por `PageBreak` (`consecutiveBlanks >= 1`).
- As ferramentas **`inspect_doc_layout` e `read_doc_content` devem ser enriquecidas** adicionando `line_spacing_percent`, margens da página (`document_geometry`), propriedades de seção e marcadores visuais explícitos de quebra de página (`[PAGE_BREAK]`).

---

## 5. Verification Method (Método de Verificação Independente)

Para reproduzir e verificar independentemente as observações apontadas:

1. **Inspeção de Código do Linter**:
   - Abrir `internal/gdocs/linter.go` e verificar a ausência do campo `LineSpacing` em todo o arquivo.
   - Linhas 312–460: verificar que apenas `TextStyle` (cores, fontes) e texto são inspecionados.

2. **Inspeção de Código do Inspetor**:
   - Abrir `internal/gdocs/inspector.go` nas linhas 10–31 (`ElementLayoutDetail`) e 117–132: constatar que `p.ParagraphStyle.LineSpacing` não é extraído.

3. **Inspeção de Código do Formatador e Limpeza**:
   - Abrir `internal/gdocs/abnt_writer.go` na linha 701 (tabelas puladas com `continue`) e linhas 343–356 (limiar `consecutiveBlanks >= 2`).
   - Abrir `internal/gdocs/cover_generator.go` na linha 911 (chamada `batchUpdate` com erro ignorado `_, _ = ...`).

4. **Comando de Testes Unitários Go**:
   - Executar no terminal:
     ```powershell
     go test ./internal/gdocs/... -v
     ```
     Verifica que os testes existentes cobrem apenas regras regex de texto e não cobrem estilos de parágrafo.
