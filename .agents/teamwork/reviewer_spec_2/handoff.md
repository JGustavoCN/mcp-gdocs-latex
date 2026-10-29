# Relatório de Handoff — reviewer_spec_2

**Data:** 29 de Outubro de 2026  
**De:** `reviewer_spec_2` (`teamwork_reviewer` / `critic`)  
**Para:** Orquestrador (`c37f0df8-888e-4a9d-8960-88e6b44fcd03`)  
**Tipo de Handoff:** Hard (Revisão Arquitetural e Crítica Adversarial Concluídas)  
**Veredicto:** **REQUEST_CHANGES**  
**Artefatos Produzidos:**
1. `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\reviewer_spec_2\analysis.md`
2. `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\reviewer_spec_2\handoff.md`

---

## 1. Observation (Observação)

Durante a auditoria analítica e adversarial dos documentos de especificação e do código-fonte do servidor Go, foram diretamente observados e verificados os seguintes fatos:

1. **Lacuna na Interface `DriveService`:**
   - Em `c:\Projetos\golang\mcp-gdocs-latex\internal\gdocs\listing.go` (linhas 34-45), a função de produção `ListAvailableDocuments` invoca diretamente o SDK do Google Drive via `svc.Files.List().Q(query).Fields(...).PageSize(100).OrderBy("name")`.
   - Em `c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md` (linhas 173-181) e em `c:\Projetos\golang\mcp-gdocs-latex\docs\architecture_specification.md` (linhas 576-584), a interface `DriveService` define exclusivamente 7 métodos (`GetFileMetadata`, `UpdateFileMedia`, `CreateFile`, `CreatePermission`, `ListComments`, `CreateReply`, `ResolveComment`), **omitindo completamente o método `ListFiles`**.

2. **Struct de Retorno Inexistente (`TableInsertionResult`):**
   - Em `c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md` (linha 185), a interface `TableService` declara a assinatura:
     `InsertNativeTable(ctx context.Context, docID string, index int64, table LaTeXTableData) (*TableInsertionResult, error)`.
   - A busca por `TableInsertionResult` em todo o workspace retornou apenas essa linha. Não há definição de struct para `TableInsertionResult` em nenhum ponto de `PROJECT.md` ou `docs/architecture_specification.md`.

3. **Armadilha de Offsets com Legenda/Fonte no Two-Step de Tabelas:**
   - Em `c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md` (linhas 279-280, 303-304), a tool `insert_native_table` declara os campos opcionais `caption` e `source`.
   - Em `docs/architecture_specification.md` (linhas 323-347), o algoritmo Two-Step define:
     Passo 1: `InsertTableRequest` no `TargetIndex`.
     Passo 2: Leitura de `CellIndex[r][c] = table.TableRows[r].TableCells[c].Content[0].StartIndex`.
     Passo 3: Lote reverso de `InsertTextRequest` preenchendo as células.
   - O algoritmo não especifica onde nem em qual etapa a `caption` (legenda superior) e a `source` (fonte inferior) são injetadas. Se a legenda for injetada antes da tabela no Passo 3, todos os `CellIndex[r][c]` calculados no Passo 2 são deslocados pelo tamanho da legenda, corrompendo as posições de preenchimento.

4. **Distorção Geométrica de Imagens Mermaid:**
   - Em `docs/architecture_specification.md` (linha 436), o payload de `InsertInlineImageRequest` fixa estaticamente:
     `objectSize: { width: { magnitude: 420.0, unit: "PT" }, height: { magnitude: 280.0, unit: "PT" } }`.
   - Os 4 diagramas em `resources/tcc/figuras/` possuem aspectos gráficos distintos (especialmente `sequencia.mmd` e `casos_uso.mmd`, que são verticalizados). A fixação incondicional de proporção 3:2 achatará ou esticará diagramas verticais.

5. **Clamping do Delimitador Final de Parágrafo (`body.EndIndex`):**
   - Em `docs/architecture_specification.md` (linhas 110-112), a especificação constata que uma deleção que alcance `body.EndIndex` gera erro HTTP `400 Bad Request: The content range must not include the final newline`.
   - Nem o schema nem o fluxo de `replace_block_abnt` especificam a regra de validação e truncamento preventivo (`endIndex = min(endIndex, body.EndIndex - 1)`).

6. **Divergência na Contagem de Ferramentas:**
   - `PROJECT.md` (linha 426) indica "Catálogo e registro de todas as 15 ferramentas MCP".
   - `internal/mcp/tools.go` registra atualmente 12 ferramentas. Somadas às 5 novas ferramentas projetadas (`apply_abnt_page_margins`, `replace_block_abnt`, `insert_native_table`, `insert_mermaid_diagram`, `check_pdf_sync`), o catálogo totalizará **17 ferramentas**.

---

## 2. Logic Chain (Cadeia Lógica)

1. *Premissa (Obs. 1):* `DriveService` foi projetado para permitir que toda a base Go seja desacoplada e testada com mocks (`httptest.Server`). `listing.go` é parte integrante do pacote `gdocs` e chama `Files.List`.
   *Conclusão 1:* Sem o método `ListFiles` em `DriveService`, a refatoração do Milestone 4 não conseguirá desacoplar `listing.go`, falhando no critério de 100% de testabilidade isolada.

2. *Premissa (Obs. 2):* Go é uma linguagem estaticamente tipada. Declarar uma interface que retorna `*TableInsertionResult` sem declarar essa struct impede a compilação do código (`undefined: TableInsertionResult`).
   *Conclusão 2:* A struct `TableInsertionResult` precisa ser explicitamente definida antes da fase de implementação.

3. *Premissa (Obs. 3):* No Google Docs API, a inserção de texto antes de qualquer índice existente desloca todos os índices subsequentes para frente.
   *Conclusão 3:* No Two-Step Approach, a inserção da legenda (que fica acima da tabela) deve ocorrer obrigatoriamente no Passo 1, antes da criação da tabela, ou no Passo 3 como a última operação atômica de um lote reverso; caso contrário, os índices de células obtidos no Passo 2 serão corrompidos.

4. *Premissa (Obs. 4):* Diagramas de software acadêmicos variam entre orientações horizontais e verticais. A imposição de dimensões rígidas ($420 \times 280\text{ pt}$) viola a integridade visual da documentação.
   *Conclusão 4:* O backend deve inspecionar a proporção real dos pixels da imagem gerada (`image.DecodeConfig`) e calcular dinamicamente a altura em pontos tipográficos para preservar o aspecto original.

5. *Premissa (Obs. 5):* Usuários e modelos de linguagem que utilizam `get_doc_skeleton` recebem blocos cujo `end_index` pode alcançar o final do documento.
   *Conclusão 5:* O handler `replace_block_abnt` deve incorporar proteção defensiva para limitar o `end_index` a `body.EndIndex - 1`, prevenindo falhas 400.

---

## 3. Caveats (Ressalvas)

1. Não foram identificadas violações de integridade (como código facade, testes com dados hardcoded ou bypass de lógica real). O trabalho desenvolvido por `worker_spec_1` é legítimo, aprofundado e de altíssimo valor técnico.
2. A suíte de testes automatizados via terminal (`go test`) não pôde ser executada neste ambiente devido a restrições de permissão do sandbox no Windows; entretanto, a análise do repositório confirmou que o projeto atualmente conta com 0 arquivos `*_test.go`, confirmando o diagnóstico de 0% de cobertura prévia relatado em `worker_spec_1/handoff.md`.
3. As alterações solicitadas não desqualificam a arquitetura concebida, mas constituem refinamentos mandatórios de engenharia para assegurar que a implementação ocorra sem bloqueios ou falhas em tempo de execução.

---

## 4. Conclusion (Conclusão e Veredicto)

**Veredicto Oficial:** **REQUEST_CHANGES**

Solicita-se a atualização de `PROJECT.md` e `docs/architecture_specification.md` com as seguintes correções pontuais:

1. **Adicionar `ListFiles` à interface `DriveService`:**
   ```go
   ListFiles(ctx context.Context, query string, fields string, pageSize int64, pageToken string) (*drive.FileList, error)
   ```
2. **Declarar a struct `TableInsertionResult`:**
   ```go
   type TableInsertionResult struct {
       TableStartIndex int64  `json:"table_start_index"`
       TableEndIndex   int64  `json:"table_end_index"`
       RowsCount       int64  `json:"rows_count"`
       ColsCount       int64  `json:"cols_count"`
       Message         string `json:"message"`
   }
   ```
3. **Sequenciar atomicamente Legenda e Fonte no algoritmo Two-Step de Tabelas:**
   - Inserir `Caption` no Passo 1 antes do `InsertTableRequest`, garantindo que a descoberta de coordenadas no Passo 2 já considere o deslocamento da legenda.
   - Inserir `Source` após o fim da tabela.
4. **Adicionar cálculo dinâmico de aspect ratio para imagens Mermaid:**
   - Decodificar dimensões do PNG via `image.DecodeConfig` e calcular $\text{height\_pt} = \text{width\_pt} \times (\text{height\_px} / \text{width\_px})$.
5. **Adicionar regra de Clamping defensivo em `replace_block_abnt`:**
   - Truncar `endIndex` para no máximo `doc.Body.EndIndex - 1` e validar conflito com comentários (`CheckCommentCollision`).
6. **Ajustar o total de ferramentas no catálogo para 17 tools.**

---

## 5. Verification Method (Método de Verificação Independente)

Para auditar e verificar as observações deste relatório:

1. **Verificação da Interface `DriveService` vs `listing.go`:**
   - Inspecionar `c:\Projetos\golang\mcp-gdocs-latex\internal\gdocs\listing.go` nas linhas 34-45 e observar a dependência de `svc.Files.List()`.
   - Comparar com `c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md` linhas 173-181 e constatar a ausência de `ListFiles`.
2. **Verificação de `TableInsertionResult`:**
   - Executar busca por `TableInsertionResult` em todo o workspace (`grep_search`) e constatar que a struct não possui declaração.
3. **Verificação das Dimensões Mermaid:**
   - Inspecionar `docs/architecture_specification.md` linha 436 e verificar os valores fixos de `width: 420.0` e `height: 280.0`.
   - Inspecionar `resources/tcc/figuras/sequencia.mmd` e constatar que se trata de diagrama vertical.
4. **Condição de Invalidação:**
   - Este veredicto será invalidado caso a Google Docs API passe a aceitar índices negativos ou suporte tabelas inline com preenchimento em chamada única, ou caso `listing.go` seja descontinuado do servidor.
