# Relatório de Análise e Crítica Arquitetural — reviewer_spec_2

**Data:** 29 de Outubro de 2026  
**Revisor:** `reviewer_spec_2` (`teamwork_reviewer` / `critic`)  
**Objeto de Revisão:**
1. `c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md` (v1.1.0)
2. `c:\Projetos\golang\mcp-gdocs-latex\docs\architecture_specification.md`
3. `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_spec_1\handoff.md`

---

## 1. Resumo Executivo e Veredicto

**Veredicto:** **REQUEST_CHANGES**

O trabalho realizado por `worker_spec_1` apresenta uma qualidade técnica e fundamentação teórica excepcionais:
- Identificou com precisão cirúrgica a causa-raiz do *Bold Bleeding* (`ReplaceAllTextRequest`) e do bloqueio de normalização de parágrafo (`internal/gdocs/writer.go:292`).
- Mapeou as discrepâncias métricas ABNT (recuo de 1,25 cm = 35,433 pt; margens 3/3/2/2 cm; espaçamento 150; justificado).
- Mapeou com exatidão as 22 tabelas e 4 diagramas Mermaid do TCC acadêmico.
- Concebeu o Two-Step Approach para tabelas nativas e a sincronização idempotente de PDF com verificação MD5 por streaming.

Contudo, sob a ótica de **revisão de engenharia e crítica adversarial**, foram identificadas **omissões de contratos de interface, armadilhas de corrupção de offsets e fragilidades técnicas** que causariam bloqueios ou retrabalho durante as fases de implementação (M1 a M5).

---

## 2. Análise Detalhada dos 3 Eixos de Robustez

### 2.1 Eixo 1: Cobertura dos Contratos de Interface em Go (Docs, Drive, Tabelas e PDF)

#### Achado 1.1 (CRÍTICO): Ausência do método `ListFiles` na interface `DriveService`
- **Onde:** `PROJECT.md:173-181` e `docs/architecture_specification.md:576-584`.
- **Problema:** A interface `DriveService` foi projetada para abstrair o SDK do Google Drive v3 e permitir desacoplamento e mocks com `httptest.Server` (M4). Ela define 7 métodos:
  - `GetFileMetadata`
  - `UpdateFileMedia`
  - `CreateFile`
  - `CreatePermission`
  - `ListComments`
  - `CreateReply`
  - `ResolveComment`
  Entretanto, o arquivo `internal/gdocs/listing.go` contém a função de produção `ListAvailableDocuments`, que consome extensivamente `svc.Files.List()` com filtros, paginação (`nextPageToken`) e ordenação.
- **Impacto:** Durante a refatoração do Milestone 4, `listing.go` não poderá utilizar `DriveService` sem quebrar a assinatura da interface, comprometendo a meta de 100% de desacoplamento e testabilidade.
- **Correção Necessária:** Incluir na interface `DriveService`:
  ```go
  ListFiles(ctx context.Context, query string, fields string, pageSize int64, pageToken string) (*drive.FileList, error)
  ```

#### Achado 1.2 (MAJOR): Struct `TableInsertionResult` Não Declarada
- **Onde:** `PROJECT.md:185`.
- **Problema:** O método `InsertNativeTable` da interface `TableService` declara retornar `(*TableInsertionResult, error)`. Porém, a struct `TableInsertionResult` **não foi declarada nem documentada** em nenhum arquivo da especificação ou projeto.
- **Correção Necessária:** Declarar formalmente no contrato Go:
  ```go
  type TableInsertionResult struct {
      TableStartIndex int64  `json:"table_start_index"`
      TableEndIndex   int64  `json:"table_end_index"`
      RowsCount       int64  `json:"rows_count"`
      ColsCount       int64  `json:"cols_count"`
      Message         string `json:"message"`
  }
  ```

#### Achado 1.3 (MENOR): Menção a `MediaService` no Diagrama sem Interface Declarada
- **Onde:** `PROJECT.md:31` vs `PROJECT.md:154-193`.
- **Problema:** O diagrama de arquitetura cita `MediaService`, mas a Seção 4.1 não define essa interface. O pipeline de diagramas Mermaid utiliza métodos diretos de `DriveService` e `DocsService`.
- **Correção Necessária:** Remover a menção a `MediaService` no diagrama ou declará-la formalmente.

---

### 2.2 Eixo 2: Ordenação Atômica do `batchUpdate` e Integridade de Offsets

#### Achado 2.1 (CRÍTICO): Armadilha de Deslocamento de Offsets com Legenda (`Caption`) no Two-Step de Tabelas
- **Onde:** `PROJECT.md:279, 303` e `docs/architecture_specification.md:323-347`.
- **Problema:** A ferramenta `insert_native_table` recebe os parâmetros `caption` (legenda no topo, ABNT) e `source` (fonte no rodapé, ABNT). O algoritmo Two-Step descrito na seção 3.1 da especificação explica como criar a tabela no `TargetIndex`, ler os índices das células `CellIndex[r][c]` no passo 2, e preenchê-las no passo 3.
  **Armadilha Identificada:** Se a legenda (`Caption`) for inserida *antes* da tabela no Passo 3 do lote, **TODAS as coordenadas de células `CellIndex[r][c]` descobertas no Passo 2 serão invalidadas**, sendo deslocadas exatamente pelo comprimento UTF-16 de `Caption`. O texto destinado à célula $(0, 0)$ será inserido no lugar errado, causando corrupção estrutural imediata!
- **Correção Necessária:** A especificação deve padronizar formalmente a sequência atômica de inserção de elementos complementares da tabela:
  - **Estratégia Recomendada (Passo 1 Antecipado):**
    No Passo 1, insere-se a `Caption` e uma quebra de linha `\n` em `TargetIndex`, e o `InsertTableRequest` é enviado em `TargetIndex + UTF16Length(Caption + "\n")`. Assim, no Passo 2, a leitura do documento via `GetDocument` já reflete a posição definitiva da tabela e suas células!
    No Passo 3, preenchem-se as células em ordem reversa, estilizam-se as bordas IBGE, e insere-se a `Source` após o fim da tabela (`table.EndIndex`).

#### Achado 2.2 (MAJOR): Proteção contra Estouro do Delimitador Final de Parágrafo (`body.EndIndex`)
- **Onde:** `docs/architecture_specification.md:110-112` e `PROJECT.md:200-246`.
- **Problema:** A especificação alerta corretamente que a Google Docs API rejeita requisições `DeleteContentRangeRequest` cujo `endIndex` seja igual a `body.EndIndex` (erro 400: *The content range must not include the final newline*). Contudo, a ferramenta `replace_block_abnt` não define uma regra de saneamento/clamping para impedir que chamadas com `end_index` no final do documento disparem esse erro.
- **Correção Necessária:** O handler de `replace_block_abnt` deve verificar o tamanho do corpo e garantir:
  ```go
  if endIndex >= doc.Body.EndIndex {
      endIndex = doc.Body.EndIndex - 1
  }
  ```

---

### 2.3 Eixo 3: Tratamento de Erros, Fallbacks e Casos de Borda das Ferramentas MCP

#### Achado 3.1 (MAJOR): Distorção de Aspect Ratio na Inserção de Diagramas Mermaid
- **Onde:** `docs/architecture_specification.md:436` e `PROJECT.md:325`.
- **Problema:** O payload de `InsertInlineImageRequest` especifica dimensões fixas:
  `objectSize: { width: { magnitude: 420.0, unit: "PT" }, height: { magnitude: 280.0, unit: "PT" } }` (proporção 3:2).
  Contudo, dos 4 diagramas do TCC, `sequencia.mmd` e `casos_uso.mmd` possuem orientações marcadamente verticais (altura muito superior à largura). Fixar largura 420 pt e altura 280 pt distorcerá as imagens violentamente, tornando-as ilegíveis.
- **Correção Necessária:** O pipeline Go deve decodificar o cabeçalho do PNG gerado (`image.DecodeConfig`) para obter `width_px` e `height_px`, calculando dinamicamente:
  $$\text{height\_pt} = \text{target\_width\_pt} \times \frac{\text{height\_px}}{\text{width\_px}}$$
  Se `height_pt` exceder a altura útil da página A4 ($500\text{ pt}$), o dimensionamento deve limitar pela altura máxima e recalcular a largura.

#### Achado 3.2 (MAJOR): Falta de Fallback e Tratamento de Falha de Execução de CLI Mermaid
- **Onde:** `docs/architecture_specification.md:421-423`.
- **Problema:** O pipeline especifica a execução de `npx @mermaid-js/mermaid-cli` ou `pdftoppm`. Se a máquina host não possuir Node.js/npx instalado ou estiver sem conexão externa, a ferramenta falhará.
- **Correção Necessária:**
  1. Verificar preventivamente a existência de versões compiladas em `resources/tcc/figuras/` ou cache local.
  2. Emitir mensagem com diagnóstico amigável (`mcplib.NewToolResultError`) instruindo sobre a dependência do `mmdc` ou `pdftoppm`.

#### Achado 3.3 (MAJOR): Omissão de Validação Concorrente (`expected_text` e `CheckCommentCollision`) em `replace_block_abnt`
- **Onde:** `PROJECT.md:219-222` e `internal/gdocs/comments.go:186`.
- **Problema:** A struct `ReplaceBlockABNTInput` possui o campo `expected_text`. No código atual de `writer.go`, as rotinas de escrita chamam `CheckCommentCollision(ctx, documentID, req.ExpectedText)` para evitar sobrescrever trechos com comentários acadêmicos pendentes. A especificação técnica de `replace_block_abnt` não detalhou a execução dessa checagem prévia nem a validação de que o texto existente na faixa `[start_index, end_index]` coincide com `expected_text`.
- **Correção Necessária:** Definir na especificação que `replace_block_abnt` deve:
  1. Extrair o texto da faixa no Google Docs;
  2. Comparar com `expected_text` (se fornecido) e abortar com erro semântico em caso de divergência;
  3. Executar `CheckCommentCollision` para garantir que comentários de orientadores/revisores não sejam apagados silenciosamente.

#### Achado 3.4 (MENOR): Divergência na Contagem do Catálogo de Tools
- **Onde:** `PROJECT.md:426` cita "todas as 15 ferramentas MCP".
- **Problema:** O código existente em `internal/mcp/tools.go` possui 12 ferramentas registradas. Com a adição das 5 novas ferramentas (`apply_abnt_page_margins`, `replace_block_abnt`, `insert_native_table`, `insert_mermaid_diagram`, `check_pdf_sync`), o total é de **17 ferramentas**. A contagem de 15 ferramentas cria ambiguidade sobre quais seriam mantidas ou descontinuadas.
- **Correção Necessária:** Atualizar a documentação para explicitar o total exato de 17 ferramentas no catálogo.

---

## 3. Avaliação de Conformidade com Requisitos (R1 a R5)

| Requisito | Status na Spec | Avaliação Crítica |
|---|---|---|
| **R1 (Pesquisa de APIs Docs/Drive/MCP)** | **CONFORME** | Excelente cobertura de endpoints, UTF-16, stdio e batchUpdate. |
| **R2 (Formatação, ABNT, Anti-Bleeding)** | **PARCIALMENTE CONFORME** | Lógica de 5 passos atômicos e métricas ABNT impecáveis; requer proteção de `body.EndIndex - 1` e validação de `expected_text`/comentários. |
| **R3 (Tabelas Nativas e Mermaid)** | **REQUER AJUSTES** | Mapeamento das 22 tabelas e 4 diagramas excelente; requer ordenação da legenda/fonte no Two-Step e proporção dinâmica de imagens. |
| **R4 (Sincronização e Hash MD5 PDF)** | **CONFORME** | Algoritmo por streaming (`crypto/md5`), comparação contra `md5Checksum` e confirmação pós-gravação 100% sólidos. |
| **R5 (Plano de Testes e Mocks)** | **REQUER AJUSTES** | Estrutura com `httptest.Server` e cobertura por milestones sólida; requer inclusão de `ListFiles` no `DriveService` e definição de `TableInsertionResult`. |

---

## 4. Verificação de Integridade e Adversarial

- **Hardcoded test results:** Nenhuma violação detectada.
- **Implementações dummy/facade:** Nenhuma violação detectada.
- **Bypass de requisitos:** Nenhuma violação detectada. O plano utiliza nativamente as APIs oficiais em Go.
- **Resultados fabricados:** Nenhuma evidência de fabricação. As análises referenciam linhas e comportamentos reais do código.

---

## 5. Recomendações e Plano de Ação para Aprovação

Para que a especificação atinja prontidão completa para implementação, os seguintes itens devem ser incorporados em `PROJECT.md` e `docs/architecture_specification.md`:

1. **Adicionar `ListFiles` à interface `DriveService`** (Seção 4.1 do `PROJECT.md` e Seção 5.1 de `architecture_specification.md`).
2. **Definir a struct `TableInsertionResult`** na Seção 4.1 do `PROJECT.md`.
3. **Especificar a sequência atômica de Legenda e Fonte no Two-Step de Tabelas**, inserindo a legenda no Passo 1 antes do cálculo dos offsets das células.
4. **Substituir dimensões estáticas (420x280 pt) por cálculo proporcional de aspect ratio** no pipeline de imagens Mermaid.
5. **Incluir validações prévias de `expected_text`, clamping de `EndIndex - 1` e `CheckCommentCollision`** na especificação da tool `replace_block_abnt`.
6. **Ajustar a contagem total de ferramentas para 17** na Seção 5 de `PROJECT.md`.
