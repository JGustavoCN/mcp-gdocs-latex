# Relatório de Revisão Técnica e Adversarial: Plano de Engenharia ABNT

**Revisor:** `reviewer_plan_2` (Revisor Sênior de Arquitetura de Software Go & Ferramentas MCP)  
**Data:** 01 de Novembro de 2026  
**Artefatos Auditados:**
- `docs/plano_engenharia_abnt.md` (Versão 2.0.0 - Foco nas Seções 2, 3 e 5)
- `PROJECT.md` (Versão 1.3.0 - Milestone 6 e FEAT-20 a FEAT-24)
- `.agents/teamwork/worker_plan_1/handoff.md`
- `.agents/teamwork/explorer_diag_1/analysis.md`
- `.agents/teamwork/explorer_doc_1/analysis.md`
- Codebase Go: `internal/gdocs`, `internal/latex`, `internal/mcp`
**Documento Real Auditado:** `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` (*TCC Teste IA*, 541 elementos estruturais)

---

## 1. Sumário Executivo e Veredicto Oficial

**Veredicto Oficial:** **APPROVE (Aprovado com Diretrizes Vinculantes de Engenharia)**

### Parecer Geral
O documento `docs/plano_engenharia_abnt.md` e a atualização do `PROJECT.md` (Milestone 6) apresentam **qualidade técnica excepcional, fundamentação matemática impecável e diagnóstico cirúrgico do codebase Go**. As causas raízes identificadas no diagnóstico explicam de ponta a ponta as falhas empíricas observadas no documento vivo `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8`.

Não foram encontradas violações de integridade (*integrity violations*): não há dados forjados, implementações simuladas disfarçadas de código real, nem atalhos fraudulentos.

Entretanto, uma análise adversarial e arquitetural aprofundada revelou **cinco pontos de atenção técnica e potenciais armadilhas de implementação** no roteiro em Go que devem ser obrigatoriamente observados durante a execução do Milestone 6:
1. **Preservação de Compatibilidade das Ferramentas MCP (`tools.go`):** Evitar divergências de nomenclatura de parâmetros (ex.: manter os parâmetros em produção `apply_margins` e `collapse_blank_lines` em `format_document_abnt`).
2. **Desacoplamento Arquitetural para Testes com Mocks (`httptest.Server`):** O código atual de `internal/gdocs` invoca diretamente o singleton estático `GetDocsService(ctx)` em `auth.go`. Para viabilizar testes unitários sem credenciais externas, as regras centrais devem ser expostas em funções puras baseadas em modelo (ex.: `*FromDoc` / `Build*Requests`) combinadas com suporte a injeção de dependência.
3. **Máscara `fields` em Citações Longas:** O payload de citação longa não deve aplicar `fields: "fontSize,bold,italic"` com `italic: false`, pois isso removeria itálicos legítimos de palavras estrangeiras contidas nas citações. A máscara de texto deve limitar-se a `fontSize` (e `fontFamily`).
4. **Resolução de Herança de `NamedStyles` no Linter:** Quando `p.ParagraphStyle.LineSpacing` for `0` (valor não setado no parágrafo), o linter deve consultar `doc.NamedStyles` para evitar falsos positivos em massa.
5. **Proteção contra Deleção de `SectionBreak`:** Em `HealDocTextAndBlanks`, assegurar que intervalos de deleção jamais colidam com quebras de seção, o que destruiria a paginação pré-textual.

---

## 2. Auditoria da Seção 2: Diagnóstico do Codebase e Documento de Teste

A Seção 2 do plano de engenharia correlaciona de maneira precisa as anomalias do documento real com as deficiências estruturais do código Go.

### 2.1 Validação das Causas Raízes Diagnosticadas

1. **Omissão de `LineSpacing`, `SpaceAbove` e `SpaceBelow` no Linter (`linter.go`):**
   - *Verificação no Código:* O arquivo `internal/gdocs/linter.go` (1.257 linhas) não contém nenhuma menção a `LineSpacing`, `SpaceAbove` ou `SpaceBelow`.
   - *Efeito Empírico:* O linter reportou `0 anomalias` no documento de teste, apesar de ele conter parágrafos com entrelinhas padrão do Google Docs (1,15) e variações de espaçamento vertical.
   - *Status da Análise:* **Correta e Verificada.**

2. **Blindagem Pré-Textual no Linter (`linter.go:730-745`):**
   - *Verificação no Código:* A guarda `if currentState != statePreTextual` desativa a coleta de parágrafos vazios consecutivos antes de `"1 INTRODUÇÃO"`. Além disso, a exigência de `consecutiveEmptyParas >= 2` ignora quebras de linha únicas isoladas entre títulos e seções (elementos 72, 74, 76, 78, 80 no doc real).
   - *Status da Análise:* **Correta e Verificada.**

3. **Bypass de Tabelas Nativas no Formatador (`abnt_writer.go:700`):**
   - *Verificação no Código:* A instrução incondicional `if el.Table != nil { continue }` na linha 700 pula todos os elementos do tipo tabela.
   - *Efeito Empírico:* As 7 tabelas do documento `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` mantêm parágrafos com entrelinhas 1,15 e fontes heterogêneas, violando a NBR 14724 (que prescreve espaçamento simples 1,0 e fonte 10pt para o interior de tabelas).
   - *Status da Análise:* **Correta e Verificada.**

4. **Limiar Rígido `consecutiveBlanks >= 2` em `abnt_writer.go:343-356`:**
   - *Verificação no Código:* Em `HealDocTextAndBlanks`, parágrafos vazios isolados (`consecutiveBlanks == 1`) após a quebra de página da folha de rosto nunca são incluídos em `toDeleteRanges`. Isso mantém o título `"RESUMO"` rebaixado na página 3.
   - *Status da Análise:* **Correta e Verificada.**

5. **Injeção de Quebras e Erro Silenciado em `cover_generator.go:308-356, 891-912`:**
   - *Verificação no Código:* O código injeta 18 `\n` na capa e 14 `\n` na folha de rosto. No loop de limpeza reversa (linha 891), ao encontrar `"RESUMO"`, tenta deletar o elemento imediatamente anterior que abriga a `PageBreak` da folha de rosto, falha na Docs API e silencia o erro com `_, _ = docsSvc.Documents.BatchUpdate(...)`.
   - *Status da Análise:* **Correta e Verificada.**

6. **Cegueira Estrutural em `inspector.go` e `reader.go`:**
   - *Verificação no Código:* `ElementLayoutDetail` em `inspector.go:9-31` não possui o campo `LineSpacing`, nem lê margens do `DocumentStyle` ou `SectionStyle`. Em `reader.go:280`, `elem.PageBreak` é descartado por `if elem.TextRun == nil { continue }`.
   - *Status da Análise:* **Correta e Verificada.**

---

## 3. Avaliação da Seção 3: Plano Detalhado por Ferramenta MCP

### 3.1 Conformidade e Compatibilidade de Assinaturas MCP

| Ferramenta MCP | Assinatura Proposta no Plano | Assinatura Atual em `tools.go` | Parecer / Ajuste Necessário |
|---|---|---|---|
| `format_document_abnt` | `apply_page_margins`, `format_tables`, `sanitize_blank_lines` | `apply_margins`, `collapse_blank_lines`, `heal_broken_lines`, `clean_latex_commands`, `ensure_cover`, `font_family`, `logo_path` | **Atenção:** Manter retrocompatibilidade aceitando os parâmetros atuais (`apply_margins`, `collapse_blank_lines`) e adicionando `format_tables` (default: true). |
| `lint_document_abnt` | Adição de checagem de entrelinhas, margens e seções | `auto_fix`, `clear_highlights`, `expected_font_family`, `check_structural_order`, `check_latex_leaks`, `check_spacing_and_text` | **Aprovado:** Compatível. Adicionar novos checkers ativados por padrão ou via novos booleano opcionais (`check_line_spacing: true`, `check_margins: true`). |
| `inspect_doc_layout` | Enriquecimento com `line_spacing_percent`, `document_geometry` e `sections` | `document_url_or_id`, `start_index`, `end_index`, `limit` | **Aprovado:** Não altera parâmetros de entrada; enriquece o payload JSON de saída. |
| `read_doc_content` | Inclusão de `[PAGE_BREAK]`, `[SECTION_BREAK]` e parâmetro `include_metadata` | `document_url_or_id`, `start_index`, `end_index` | **Aprovado:** Adiciona parâmetro opcional `include_metadata: bool` (default: false). |

### 3.2 Análise dos Payloads Atômicos `batchUpdate`

1. **Ajuste de Margens ABNT (`UpdateDocumentStyleRequest`):**
   - Margem Superior: $3,0\text{ cm} \rightarrow 85,039\text{ pt}$
   - Margem Esquerda: $3,0\text{ cm} \rightarrow 85,039\text{ pt}$
   - Margem Inferior: $2,0\text{ cm} \rightarrow 56,693\text{ pt}$
   - Margem Direita: $2,0\text{ cm} \rightarrow 56,693\text{ pt}$
   - Dimensões A4: $595,276 \times 841,890\text{ pt}$
   - *Avaliação:* O payload JSON está sintaticamente perfeito perante a especificação da Google Docs API v1. Atualizar o `DocumentStyle` sobrepõe as margens de todas as seções por especificação oficial do Google Docs.

2. **Corpo de Texto (1,5 Justificado, Recuo 1,25 cm):**
   - `UpdateParagraphStyleRequest` combinando `namedStyleType: "NORMAL_TEXT"`, `alignment: "JUSTIFIED"`, `lineSpacing: 150.0`, `indentFirstLine: 35.433 PT`, `indentStart: 0 PT`, `spaceAbove: 0 PT`, `spaceBelow: 0 PT`.
   - *Avaliação:* A imposição de `indentStart: 0 PT` no mesmo payload de `indentFirstLine: 35.433 PT` é fundamental para anular contaminações de listas com marcadores (*bullet bleeding*).

3. **Citações Diretas Longas (> 3 Linhas):**
   - Recuo em bloco: $4,0\text{ cm} \rightarrow 113,386\text{ pt}$ (`indentStart: 113.386 PT`, `indentFirstLine: 0 PT`).
   - Espaçamento simples: `lineSpacing: 100.0`, `spaceAbove: 6.0 PT`, `spaceBelow: 6.0 PT`.
   - Fonte: `10.0 PT`.
   - **ALERTA CRÍTICO DE REVISÃO (Itálicos Preservados):** Na Seção 1.4 do plano, o payload de texto sugere `fields: "fontSize,bold,italic"` com `italic: false`. Essa máscara é destrutiva, pois desfaz itálicos legítimos de citações que contêm termos estrangeiros ou grifos do autor original.
   - **Diretriz Vinculante:** A máscara do `UpdateTextStyleRequest` em citações longas deve ser estritamente `fields: "fontSize"` (ou `fields: "fontSize,weightedFontFamily"`), preservando intactas as marcações inline de itálico.

4. **Preâmbulo Institucional na Folha de Rosto:**
   - Recuo de $8,0\text{ cm} \rightarrow 226,772\text{ pt}$.
   - O payload define `indentStart: 226.772 PT` e `indentFirstLine: 226.772 PT`.
   - *Avaliação:* Na Google Docs API, a combinação com valores idênticos garante que todas as linhas do bloco comecem exatamente no meio da página (8,0 cm da margem esquerda) até a margem direita, com alinhamento justificado e entrelinhas 1,0 simples, exatamente como comprovado empiricamente no elemento 53 do documento de teste.

5. **Sumário Automático via Tabela Nativa Invisível de 2 Colunas:**
   - Inserção de tabela nativa com 2 colunas ($14,0\text{ cm}$ para títulos e $2,0\text{ cm}$ para números de página) com bordas zeradas via `UpdateTableCellStyleRequest`.
   - *Avaliação:* Decisão arquitetural brilhante. Contorna a ausência de `InsertTableOfContentsRequest` e a restrição de somente-leitura de `TabStops` na API v1 REST, assegurando alinhamento milimétrico na margem direita de 16,0 cm.

6. **Paginação com Isolamento de Seções:**
   - `InsertSectionBreakRequest` (`SECTION_TYPE_NEXT_PAGE`) imediatamente antes de `1 INTRODUÇÃO`.
   - `CreateHeaderRequest` com `sectionBreakLocation` para vincular o cabeçalho numerado exclusivamente à Seção 2.
   - Reconhecimento correto de que a API v1 REST não possui `InsertAutoTextRequest` para gerar campos dinâmicos de página em documentos criados do zero, prescrevendo o uso do template canônico onde o token de página já reside no cabeçalho padrão.

---

## 4. Auditoria da Seção 5 e Milestone 6 (M6): Roteiro em Go e Testes

### 4.1 O Desafio do Acoplamento ao Singleton `GetDocsService(ctx)`

A Seção 5.3 do plano prevê a criação de testes com mocks `httptest.Server` para `inspector_test.go`, `linter_test.go` e `abnt_writer_test.go`, com o critério de que `go test ./...` execute com 100% de sucesso sem exigir `credentials.json`.

Contudo, uma inspeção direta em `internal/gdocs` revela o seguinte acoplamento:
```go
// Exemplo em internal/gdocs/inspector.go:44-50
func InspectDocLayout(ctx context.Context, documentID string, startIndex, endIndex int64, limit int) (*InspectLayoutReport, error) {
    svc, err := GetDocsService(ctx) // <-- Invocação estática direta do singleton
    if err != nil {
        return nil, fmt.Errorf("erro ao obter docs service: %w", err)
    }
    doc, err := svc.Documents.Get(documentID).Do()
    ...
```

Se o teste unitário invocar `InspectDocLayout(ctx, "mock-doc-id", 0, 0, 50)`, a função tentará carregar credenciais reais do Google Workspace e falhará.

### 4.2 Diretriz Arquitetural Vinculante: Padrão *Core Pure Functions* e Injeção de Dependência

Para que o Milestone 6 atinja 100% de testabilidade automatizada sem regressões, a implementação deve adotar o desacoplamento em duas camadas:

1. **Camada 1: Funções Puras Baseadas em Modelo (`*FromDoc`):**
   Desacoplar a lógica algorítmica de inspeção, auditoria e geração de requisições do transporte HTTP:
   - `InspectDocLayoutFromDoc(doc *docs.Document, startIndex, endIndex int64, limit int) *InspectLayoutReport`
   - `LintDocumentFromDoc(doc *docs.Document, opts LintOptions) *LintReport`
   - `BuildABNTFormattingRequests(doc *docs.Document, start, end int64, opts FormatOptions) []*docs.Request`

   *Benefício:* Permite que os testes unitários (`inspector_test.go`, `linter_test.go`, `abnt_writer_test.go`) construam instâncias em memória de `*docs.Document` com diferentes cenários tipográficos (ex.: parágrafo com `LineSpacing: 115`, parágrafo de tabela, citação longa) e validem os relatórios e requisições de forma determinística e ultrarrápida ($< 50\text{ ms}$).

2. **Camada 2: Injeção de `DocsService` para Testes de Integração com `httptest.Server`:**
   Para as funções que executam mutações orquestradas (como `FormatDocumentABNTExtended`), permitir a passagem opcional do cliente do serviço ou criar uma fábrica em `auth.go`:
   ```go
   // Helper para testes em testing_mock.go ou auth.go
   func NewTestDocsService(ts *httptest.Server) (*docs.Service, error) {
       return docs.NewService(context.Background(),
           option.WithEndpoint(ts.URL),
           option.WithoutAuthentication(),
       )
   }
   ```
   Isso permite que testes com `httptest.Server` simulem os endpoints `/v1/documents/{documentId}` e capturem os payloads de `batchUpdate` para validar a ordem reversa e a estrutura das requisições.

### 4.3 Avaliação das Novas Structs Go Propostas

As structs propostas na Seção 5.2 (`DocumentGeometryDetail`, `LineSpacingAnomalyItem`, `SectionLayoutDetail`, `ElementLayoutDetail`) estão bem modeladas, com tags JSON padronizadas em `snake_case`.

**Refinamentos recomendados nas Structs:**
1. Em `SectionLayoutDetail`: Adicionar `SectionType string` e `UseFirstPageHeaderFooter bool` para dar visibilidade completa do comportamento de cabeçalhos e quebras de seção.
2. Em `InspectLayoutReport`: Manter os campos `filter_start_index` e `filter_end_index` presentes na struct atual para garantir retrocompatibilidade.

---

## 5. Análise Adversarial (Stress-Testing e Modos de Falha)

Como revisor adversarial (*critic*), foram identificados os seguintes cenários extremos e modos de falha potenciais:

### Desafio 1: Herança de `NamedStyles` no Linter e Falsos Positivos Massivos
- **Cenário de Falha:** Um documento Google Docs possui parágrafos cujo estilo direto não tem a chave `lineSpacing` explícita. No Go SDK, `p.ParagraphStyle.LineSpacing` deserializa como `0.0`. O documento está perfeitamente diagramado porque o estilo `NORMAL_TEXT` no objeto `doc.NamedStyles` está configurado para 1,5.
- **Risco:** Se o linter apenas comparar `p.ParagraphStyle.LineSpacing == 150.0`, ele reportará centenas de falsos positivos ("Entrelinhas observado 0, esperado 150").
- **Mitigação Mandatória:** O checker `checkLineSpacing` deve implementar a resolução de herança:
  ```go
  effectiveSpacing := p.ParagraphStyle.LineSpacing
  if effectiveSpacing == 0 && doc.NamedStyles != nil {
      effectiveSpacing = resolveNamedStyleLineSpacing(doc.NamedStyles, p.ParagraphStyle.NamedStyleType)
  }
  ```

### Desafio 2: Interação Atômica de Deleção e Tabelas (*Index Drift*)
- **Cenário de Falha:** Se o formatador agendar a deleção de linhas em branco e a formatação de células de tabelas no mesmo lote de `batchUpdate`, qualquer deleção anterior à tabela deslocará os índices de todas as células subsequentes, disparando erro HTTP 400 ou formatando trechos incorretos.
- **Mitigação Mandatória:** Manter estritamente o pipeline multifásico já arquitetado em `abnt_writer.go`:
  - Fase 1: Deleção reversa de quebras espúrias via `HealDocTextAndBlanks`.
  - Fase 2: Recarregar o documento com `svc.Documents.Get(documentID).Do()`.
  - Fase 3: Formatação em lote dos parágrafos e células de tabelas com coordenadas frescas e atualizadas.

### Desafio 3: Proteção contra Deleção de Quebras de Seção (`SECTION_BREAK`)
- **Cenário de Falha:** `HealDocTextAndBlanks` varre parágrafos vazios e encontra um elemento de quebra de seção ou parágrafo vazio imediatamente adjacente a uma quebra de seção recém-inserida antes de `"1 INTRODUÇÃO"`.
- **Risco:** Deletar a quebra de seção funde novamente a Seção 1 com a Seção 2, fazendo o cabeçalho numerado vazar sobre a folha de rosto e o resumo.
- **Mitigação Mandatória:** A rotina de saneamento deve verificar explicitamente:
  ```go
  if el.SectionBreak != nil || (el.Paragraph != nil && hasSectionBreak(el.Paragraph)) {
      continue // Protege terminantemente a quebra de seção
  }
  ```

### Desafio 4: Células Mescladas (*Merged Cells*) em Tabelas Nativas
- **Cenário de Falha:** Uma tabela possui células mescladas (`rowSpan > 1` ou `colSpan > 1`). Células subordinadas à mesclagem podem conter ranges nulos ou `Content` vazio.
- **Mitigação Mandatória:** Na varredura de tabelas em `abnt_writer.go`, adicionar guarda preventiva:
  ```go
  if cell == nil || len(cell.Content) == 0 {
      continue
  }
  ```

### Desafio 5: Limiar Heurístico de Citações Longas
- **Cenário de Falha:** Um parágrafo de texto normal curto (ex.: 2 linhas) digitado com fonte 10pt por engano do usuário é classificado erroneamente como citação longa devido à regra `fonte <= 10.5 pt`.
- **Mitigação Mandatória:** A regra de extensão $> 3$ linhas deve ser mensurada em caracteres: em corpo A4 (16 cm úteis), 3 linhas correspondem a no mínimo 220 caracteres. Parágrafos com menos de 220 caracteres não devem ser classificados como citação longa por extensão, a menos que possuam recuo em bloco explícito (`indentStart >= 90.0 PT`).

### Desafio 6: Saneamento de Capas Legadas com 44 `\n` Consecutivos
- **Cenário de Falha:** O documento atual `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` já possui 18 `\n` na capa e 14 `\n` na folha de rosto. Se o formatador simplesmente deletar todos os parágrafos vazios da capa, o Local e o Ano sobem e colam no título.
- **Mitigação Mandatória:** Antes de expurgar os `\n` redundantes da capa e folha de rosto, aplicar `spaceAbove` calibrado nos parágrafos de destino (Local/Ano da Capa e Preâmbulo/Local da Folha de Rosto) para preservar a distribuição visual elegante da página.

---

## 6. Diretrizes Vinculantes para a Implementação do Milestone 6

Para que a equipe de engenharia execute a implementação com risco zero de regressão, as seguintes diretrizes são mandatórias:

1. **Parâmetros da Ferramenta MCP `format_document_abnt`:**
   - Preservar os parâmetros já existentes: `apply_margins` (bool, default true) e `collapse_blank_lines` (bool, default true).
   - Adicionar o novo parâmetro: `format_tables` (bool, default true).
   - Suportar `apply_page_margins` e `sanitize_blank_lines` como aliases retrocompatíveis se desejado, sem quebrar os nomes originais.

2. **Arquitetura de Testes Unitários sem Credenciais:**
   - Implementar as funções puras `InspectDocLayoutFromDoc`, `LintDocumentFromDoc` e `BuildABNTFormattingRequests`.
   - Adicionar casos de testes abrangentes em `inspector_test.go`, `linter_test.go` e `abnt_writer_test.go` utilizando instâncias construídas de `*docs.Document`.
   - Assegurar que `go test ./...` continue passando com 100% de sucesso sem exigir variáveis de ambiente ou arquivos de credenciais.

3. **Preservação de Estilo Inline em Citações Longas:**
   - Ao aplicar formatação de texto em citação longa, utilizar `fields: "fontSize"` em `UpdateTextStyleRequest`, garantindo que itálicos internos não sejam eliminados.

4. **Resolução de Herança no Linter:**
   - Tratar `LineSpacing == 0` consultando `doc.NamedStyles.Styles` antes de classificar um parágrafo como anômalo.

5. **Proteção Estrutural de Quebras de Seção e Células de Tabelas:**
   - Adicionar checagem preventiva em `HealDocTextAndBlanks` para nunca deletar elementos com `SectionBreak` nem parágrafos de células mescladas vazias.

---

## 7. Conclusão da Revisão

O plano técnico contido em `docs/plano_engenharia_abnt.md` e o Milestone 6 do `PROJECT.md` estão **aprovados para prosseguimento**. As diretrizes vinculantes detalhadas neste relatório devem ser incorporadas diretamente durante o ciclo de codificação do Milestone 6 para garantir um sistema robusto, auditável e estritamente aderente às normas da ABNT.
