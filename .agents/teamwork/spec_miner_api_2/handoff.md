# Handoff Report — Mineração de Especificações Técnicas das APIs Google Docs v1 e Drive v3 (R1)

**Agente Remetente:** `spec_miner_api_2`  
**Destinatário:** Orquestrador (`parent` / `8522a030-5a02-4454-bc22-525b48f7de69`)  
**Data:** 01 de Novembro de 2026  
**Tipo de Handoff:** Hard (Tarefa Concluída)  
**Artefato Primário:** `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\spec_miner_api_2\analysis.md`

---

## 1. Observation (Observações Diretas)

1. **Definições e Estruturas da Google Docs API v1 (`docs/pesquisa-api.md` e Go SDK `google.golang.org/api/docs/v1`):**
   - `DocumentStyle` (`docs/pesquisa-api.md:2093-2197`): Possui campos `MarginTop`, `MarginBottom`, `MarginLeft`, `MarginRight`, `PageSize`, `PageNumberStart`, `UseFirstPageHeaderFooter`. A documentação oficial explicita que atualizar margens no `DocumentStyle` limpa automaticamente as margens correspondentes em todos os `SectionStyle`.
   - `UpdateDocumentStyleRequest` (`docs/pesquisa-api.md:7014-7045`): Requer `documentStyle` e `fields`.
   - `ParagraphStyle` (`docs/pesquisa-api.md:4416-4550`):
     - `LineSpacing float64` (`linha 4495`): Percentual relativo a 100.0 (onde 100.0 = 1,0 simples, 150.0 = 1,5 ABNT).
     - `Alignment string` (`linha 4429`): `"START"`, `"CENTER"`, `"END"`, `"JUSTIFIED"`.
     - `IndentFirstLine *Dimension` (`linha 4479`): Recuo de primeira linha.
     - `IndentStart *Dimension` (`linha 4483`): Recuo em bloco da margem inicial.
     - `TabStops []*TabStop` (`linha 4537`): Documentado oficialmente como **`This property is read-only`**.
   - `CreateHeaderRequest` (`docs/pesquisa-api.md:1156-1182`): Aceita `type: "DEFAULT" | "FIRST_PAGE"` e `sectionBreakLocation`. Cria um cabeçalho vinculado especificamente ao `SectionStyle` da quebra de seção. Retorna `400 Bad Request` se já existir cabeçalho do tipo especificado.
   - `InsertSectionBreakRequest` (`docs/pesquisa-api.md:3528-3564`): Aceita `sectionType: "NEXT_PAGE" | "CONTINUOUS"` e `location`. Insere quebra de seção que particiona o documento.
   - `DeleteContentRangeRequest` (`docs/pesquisa-api.md:1583-1610`): Documenta formalmente restrições que causam `400 Bad Request`: tentar deletar o último `\n` do corpo do documento (`body.EndIndex - 1`); deletar um code unit isolado de par substituto UTF-16; deletar o parágrafo delimitador imediatamente anterior a uma `Table`, `TableOfContents` ou `SectionBreak`.
   - **Ausência de Requisições:** A enumeração exaustiva de todos os tipos `Request` da Docs API v1 confirmou que **NÃO EXISTEM** `InsertTableOfContentsRequest`, `UpdateTableOfContentsRequest`, `InsertAutoTextRequest` ou `InsertPageNumberRequest` no protocolo `batchUpdate`.

2. **Código Existente do Servidor MCP:**
   - Em `internal/gdocs/builder.go:7-24`: As constantes métricas já adotam a conversão exata $1\text{ cm} = 28,3464567\text{ pt}$ (`ABNTMarginTopPt = 85.039`, `ABNTMarginLeftPt = 85.039`, `ABNTMarginBottomPt = 56.693`, `ABNTMarginRightPt = 56.693`, `ABNTIndentFirstLinePt = 35.433`, `ABNTIndentLongQuotePt = 113.386`, `ABNTIndentPreTextualNotePt = 226.772`, `ABNTLineSpacing = 150`, `ABNTLineSpacingSingle = 100`).
   - Em `internal/gdocs/linter.go:103-126`: A struct `LintReport` **não contém nenhum campo para `LineSpacing`**. Na função `LintDocument` (`linhas 312-760`), `p.ParagraphStyle.LineSpacing` nunca é lido nem validado.
   - Em `internal/gdocs/linter.go:734`: A checagem de linhas em branco consecutivas contém a guarda `if currentState != statePreTextual`, fazendo com que todas as linhas vazias entre a folha de rosto e o primeiro capítulo numerado sejam ignoradas.
   - Em `internal/gdocs/inspector.go:9-31`: A struct `ElementLayoutDetail` **omite o campo `LineSpacing`**, impedindo que clientes MCP ou LLMs diagnostiquem o espaçamento entrelinhas através de `inspect_doc_layout`.

---

## 2. Logic Chain (Cadeia de Raciocínio Lógico)

1. **Margens e Dimensões (Item 1):**
   - Observação: $1\text{ cm} = \frac{72}{2,54}\text{ pt} \approx 28,3464567\text{ pt}$.
   - Inferência: As margens ABNT NBR 14724 ($3,0\text{ cm}$ superior e esquerda, $2,0\text{ cm}$ inferior e direita) convertem-se estritamente para $85,039\text{ pt}$ e $56,693\text{ pt}$, e a folha A4 para $595,276\text{ pt} \times 841,890\text{ pt}$.
   - Aplicação: A chamada `UpdateDocumentStyleRequest` com máscara `fields: "marginTop,marginLeft,marginBottom,marginRight,pageSize"` padroniza o documento globalmente e remove automaticamente margens anômalas residuais em todas as seções.

2. **Espaçamento e Alinhamento (Item 2) e Causa da Falha no Teste:**
   - Observação: No Docs API, `lineSpacing: 150` equivale a 1,5 entrelinhas e `lineSpacing: 100` a 1,0 simples. O linter atual (`internal/gdocs/linter.go`) não verifica `p.ParagraphStyle.LineSpacing`, e o inspector (`internal/gdocs/inspector.go`) não o expõe.
   - Conclusão: O documento de teste `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` permaneceu com variações de entrelinhas (1,15 vs 1,5) porque o linter possui um ponto cego arquitetural (falso negativo silencioso).

3. **Recuo de 1,25 cm e Blindagem contra Bleeding (Item 3):**
   - Observação: $1,25\text{ cm} = 35,433\text{ pt}$. No Google Docs, `indentStart` desloca o bloco inteiro e `indentFirstLine` a primeira linha.
   - Conclusão: Se `indentStart` contiver valor residual (ex: 36 pt de marcadores anteriores), o parágrafo fica com recuo duplo. O payload deve obrigatoriamente forçar `indentFirstLine: 35.433 pt` e `indentStart: 0 pt` simultaneamente.

4. **Citação Longa e Preâmbulo (Itens 4 e 5):**
   - Citação longa: $4,0\text{ cm} = 113,386\text{ pt}$ em bloco (`indentStart: 113.386`, `indentFirstLine: 0`), 10 pt (`UpdateTextStyleRequest`), espaçamento simples (`lineSpacing: 100`), sem aspas.
   - Preâmbulo da Folha de Rosto: $8,0\text{ cm} = 226,772\text{ pt}$ (`indentStart: 226.772`, `indentFirstLine: 226.772`), justificado, 10 pt, espaçamento simples.

5. **Numeração de Página, Quebra de Seção e Cabeçalhos (Item 6):**
   - Observação: `InsertSectionBreakRequest` com `SECTION_TYPE_NEXT_PAGE` particiona o documento. `CreateHeaderRequest` com `sectionBreakLocation` cria cabeçalho isolado para a Seção 2.
   - Conclusão: Deletando o cabeçalho global do `DocumentStyle` e criando o cabeçalho exclusivamente na Seção 2, as páginas pré-textuais (Seção 1) ficam rigorosamente sem numeração. Como a Seção 2 herda a contagem acumulada, o número da página reflete a contagem desde a Folha de Rosto.
   - Restrição da API: A API v1 não permite inserir o token `AutoText: PAGE_NUMBER` via `batchUpdate`. Ele deve provir de template base do Google Drive ou o layout deve ser pré-configurado.

6. **Sumário Automático (Item 7):**
   - Observação: A API Docs v1 não possui `InsertTableOfContentsRequest`, e `ParagraphStyle.TabStops` é `read-only` (além de não suportar leader dots).
   - Conclusão: A abordagem ótima na API v1 é a geração de um Sumário ABNT estruturado via **Tabela Nativa de 2 colunas sem bordas visíveis** (Coluna 1: títulos NBR 6024, Coluna 2: números alinhados à direita).

7. **Remoção de Linhas em Branco (Item 8) e Causa da Falha no Teste:**
   - Observação: O linter ignora parágrafos vazios se `currentState == statePreTextual`. O saneador falha em identificar o marco pós-folha de rosto se o ano não bater exatamente com `"2026"`.
   - Conclusão: As quebras espúrias após a folha de rosto não foram eliminadas por falha na máquina de estados do linter e fragilidade no marco de corte do `abnt_writer.go`.
   - Regra de Deleção: Toda exclusão via `DeleteContentRangeRequest` deve rodar em ordem reversa (*Reverse Index Order Rule*), nunca tocar em `body.EndIndex - 1` e nunca deletar o parágrafo imediatamente anterior a tabelas.

---

## 3. Caveats (Ressalvas e Limitações)

1. **AutoText na API Docs v1:** A criação programática ex-nihilo do campo de numeração de página dinâmica (`PAGE_NUMBER`) não é suportada por requisições de escrita da Docs API v1. É indispensável que documentos em branco utilizem um modelo institucional base no Google Drive com o campo pré-inserido no cabeçalho ou dependam de sincronização de PDF.
2. **TabStops na API Docs v1:** Modificações de tabulações personalizadas via `UpdateParagraphStyleRequest` são rejeitadas pela API porque a propriedade é somente leitura (`read-only`).
3. **Modo Sem Páginas (Pageless Mode):** Se um usuário configurar o Google Doc no modo `PAGELESS`, as margens de página e os cabeçalhos/rodapés deixam de ser renderizados visualmente pela interface do Google Docs, embora continuem armazenados na estrutura de dados da API.

---

## 4. Conclusion (Conclusão)

A viabilidade técnica de adequação às normas ABNT NBR 14724 nas APIs do Google Docs v1 e Google Drive v3 é **plenamente comprovada e matematicamente parametrizável**, com exceção das duas restrições de escrita conhecidas da API v1 (`AutoText` dinâmico e `InsertTableOfContentsRequest`), para as quais existem soluções arquiteturais de engenharia consolidadas (uso de templates para numeração e tabela nativa de 2 colunas para o sumário).

As falhas observadas no documento de teste `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` foram cabalmente explicadas:
- Omissão de auditoria de `LineSpacing` no `linter.go`;
- Supressão de auditoria de linhas em branco no `statePreTextual` do `linter.go`;
- Omissão de `LineSpacing` na struct `ElementLayoutDetail` do `inspector.go`.

---

## 5. Verification Method (Método de Verificação Independente)

1. **Inspeção dos Arquivos Entregues:**
   - Relatório técnico exaustivo: `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\spec_miner_api_2\analysis.md`
   - Relatório de handoff: `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\spec_miner_api_2\handoff.md`
2. **Validação do Código Go e Cobertura de Testes Existentes:**
   - Executar `go test ./internal/gdocs/...` para certificar que os testes unitários atuais continuam passando.
3. **Verificação de Schemas e Tipos:**
   - Conferir as structs em `docs/pesquisa-api.md`:
     - Linhas 4417-4550 (`ParagraphStyle`), confirmando `LineSpacing float64` e `TabStops read-only`.
     - Linhas 2093-2197 (`DocumentStyle`), confirmando margens e dimensões.
     - Linhas 1156-1205 (`CreateHeaderRequest`), confirmando criação de cabeçalho por quebra de seção.
     - Linhas 1583-1610 (`DeleteContentRangeRequest`), confirmando restrições de deleção.
