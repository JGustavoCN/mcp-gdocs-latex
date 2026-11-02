# Plano de Projeto e Arquitetura de Sistema: `mcp-gdocs-latex`

**Versão:** 1.3.0  
**Data:** 01 de Novembro de 2026  
**Status:** Especificação Definitiva Aprovada para Implementação (Iteração 4 - Conformidade Estrita ABNT)  
**Escopo:** Servidor MCP em Go para Interoperabilidade Bidirecional entre LaTeX Acadêmico (TCC ABNT) e Google Workspace (Docs v1 e Drive v3) com Conformidade ABNT NBR 14724 Estrita  

---

## 1. Visão Geral da Arquitetura do Servidor Go e Ecossistema MCP

O servidor `mcp-gdocs-latex` atua como uma ponte de alta fidelidade entre documentos acadêmicos escritos em LaTeX (normas ABNT NBR 14724 / NBR 6023) e o ambiente colaborativo do Google Workspace (Google Docs e Google Drive), utilizando o protocolo padronizado **Model Context Protocol (MCP)** via transporte `stdio` JSON-RPC 2.0.

```
┌────────────────────────────────────────────────────────────────────────┐
│                          Host MCP / Cliente LLM                         │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │ JSON-RPC 2.0 via Stdio
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│                   Servidor Go: mcp-gdocs-latex                         │
│                                                                        │
│  ┌───────────────────────┐         ┌────────────────────────────────┐  │
│  │   Catálogo de Tools   │         │    Camada de Sanitização &     │  │
│  │   (mark3labs/mcp-go)  │ ◄─────► │    UTF-16 Offset Shifting      │  │
│  └───────────┬───────────┘         └───────────────┬────────────────┘  │
│              │                                     │                   │
│              ▼                                     ▼                   │
│  ┌──────────────────────────────────────────────────────────────────┐  │
│  │           BatchUpdateBuilder & Adapter Interfaces                │  │
│  │      (DocsService, DriveService, TableService, MediaService)     │  │
│  └───────────┬─────────────────────────────────────┬────────────────┘  │
└──────────────┼─────────────────────────────────────┼───────────────────┘
               │ HTTPS (OAuth2 / Service Account)    │ Leitura / Escrita
               ▼                                     ▼
┌───────────────────────────────┐     ┌──────────────────────────────────┐
│ Google Workspace APIs         │     │ Sistema de Arquivos Local        │
│ ├─ Google Docs API (v1)       │     │ ├─ resources/tcc/ (LaTeX Fonte)  │
│ └─ Google Drive API (v3)      │     │ └─ resources/tcc/main.pdf (Hash) │
└───────────────────────────────┘     └──────────────────────────────────┘
```

### 1.1 Princípios Fundamentais de Engenharia
1. **Atomicidade Estrita de Mutação:** Todas as alterações estruturais em documentos Google Docs devem ser executadas em transações atômicas únicas via `documents.batchUpdate`. Nenhuma mutação intermediária deve deixar o documento em estado inconsistente.
2. **Isolamento de Estilos (Zero-Bleeding):** Eliminação categórica do efeito cascata de negrito (*Bold Bleeding*) e de listas (*List/Bullet Bleeding*) através do encadeamento atômico de deleção, inserção, remoção explícita de bullets e aplicação forçada de `NORMAL_TEXT` e estilos ABNT.
3. **Métrica Canônica ABNT:** Parametrização precisa em pontos tipográficos ($1\text{ cm} = 28,3464567\text{ pt}$) para recuo de primeira linha ($1,25\text{ cm} = 35,433\text{ pt}$), entrelinhas 1,5 ($150\%$), alinhamento justificado e margens 3/3/2/2 cm.
4. **Idempotência Criptográfica:** Toda sincronização de mídia (especialmente o PDF compilado do TCC) é precedida pelo cálculo e validação do hash MD5 local contra o campo `md5Checksum` da API Google Drive v3, abortando uploads desnecessários.
5. **Two-Step Approach para Elementos Estruturais:** Tabelas nativas e imagens inline operam em fluxos de duas etapas para garantir a correta indexação de células e referências públicas de mídia.
6. **Desacoplamento e Testabilidade Total:** Isolamento do SDK do Google através do padrão Adapter, viabilizando 100% de cobertura de testes automatizados com `httptest.Server` e dispensando credenciais reais durante a suíte de testes.

---

## 2. Feature Inventory Exaustiva

A tabela a seguir consolida todas as funcionalidades existentes, refatoradas e novas do servidor `mcp-gdocs-latex`, associando cada uma ao seu Milestone de entrega correspondente.

| ID | Funcionalidade | Tipo | Milestone | Pacote Responsável | Descrição Técnica |
|---|---|---|---|---|---|
| **FEAT-01** | `apply_abnt_page_margins` | Nova Tool | M1 | `internal/gdocs` | Configura margens ABNT de página (3cm sup/esq, 2cm inf/dir) via `UpdateDocumentStyleRequest` em pontos tipográficos. |
| **FEAT-02** | `replace_block_abnt` | Nova Tool | M1 | `internal/gdocs` | Substitui trechos com isolamento anti-bleeding, remoção de bullets, `indentStart: 0 pt`, clamping de fim de documento e normalização NFC. |
| **FEAT-03** | Correção de `NORMAL_TEXT` em `update_block_by_index` | Fix / Refactor | M1 | `internal/gdocs` | Remove a guarda de bloqueio `bloco.ParagraphStyle != "NORMAL_TEXT"` em `writer.go:292`, forçando a normalização de parágrafo. |
| **FEAT-04** | Higienizador de Texto UTF-16 (`SanitizeText`) | Nova Função | M1 | `internal/gdocs` | Normalização canônica Unicode NFC (`norm.NFC`), colapso de `\n{3,}`, espaços à direita e cálculo seguro em UTF-16 code units. |
| **FEAT-05** | `BatchUpdateBuilder` | Novo Core | M1 | `internal/gdocs` | Componente Go fluente para lotes atômicos com clamping defensivo de fim de doc, ordenação reversa e guarda de texto vazio. |
| **FEAT-06** | `insert_native_table` | Nova Tool | M2 | `internal/gdocs` | Cria tabelas ABNT com legenda no Passo 1, inserção da Fonte em `tableEndIndex` como primeira operação do lote reverso com `WriteControl` (Passo 3) e retorno `TableInsertionResult`. |
| **FEAT-07** | Estilização ABNT/IBGE de Tabelas | Nova Função | M2 | `internal/gdocs` | Aplica `UpdateTableCellStyleRequest` com `TableStartLocation` e `TableRange` para zerar bordas laterais e padronizar réguas. |
| **FEAT-08** | Conversor LaTeX Tabular -> Estrutura Go | Novo Parser | M2 | `internal/latex` | Extrai matrizes de dados tratando descarte de `\addlinespace`, normalização de `\textcite`, unescape de `\%` e fórmulas inline (`$2 \times 2$`). |
| **FEAT-09** | Pipeline de Imagens Mermaid (`insert_mermaid_diagram`) | Nova Tool | M2 | `internal/gdocs` | Converte diagramas em PNG 300 DPI, enquadramento 2D Bounding Box ($450 \times 620\text{ pt}$) via `image.DecodeConfig` com import `_ "image/png"` e fallback de Deep Link. |
| **FEAT-10** | `check_pdf_sync` | Nova Tool | M3 | `internal/gdocs` | Valida integridade física prévia (`%%EOF` nos últimos 1024 bytes, ISO 32000-1) e compara hash MD5 com `md5Checksum` do Drive v3. |
| **FEAT-11** | `sync_pdf_to_drive` Idempotente | Refactor Tool | M3 | `internal/gdocs` | Upload condicional com verificação prévia de `%%EOF`, prevenção TOCTOU (`copiedBytes == stat.Size()`), política `keepRevisionForever: true` e validação estrita. |
| **FEAT-12** | Adapter Pattern para Google Docs (`DocsService`) | Refactor Core | M4 | `internal/gdocs` | Cria interface desacoplada para isolar chamadas diretas ao `*docs.Service`. |
| **FEAT-13** | Adapter Pattern para Google Drive (`DriveService`) | Refactor Core | M4 | `internal/gdocs` | Interface desacoplada com suporte a `ListFiles` para desacoplar e testar `internal/gdocs/listing.go`. |
| **FEAT-14** | Infraestrutura de Mocks com `httptest.Server` | Nova Infra | M4 | `internal/gdocs` | Servidor HTTP mock multiplexado cobrindo `/v1/documents/`, `/drive/v3/files/` e rota multipart `/upload/drive/v3/files/` com extração normalizada de `fileID`. |
| **FEAT-15** | Suíte de Testes Unitários de Funções Puras | Novos Testes | M4 | `internal/*` | Cobre `ExtractDocID`, cálculo UTF-16, `SanitizeText` com NFC, validação `%%EOF` (ISO 32000-1) e TOCTOU, clamping e parser de tabelas com macros LaTeX. |
| **FEAT-16** | Suíte de Testes de Integração com Mocks | Novos Testes | M4 | `internal/gdocs` | Testa fluxos de `replace_block_abnt`, Two-Step com `WriteControl` e retry, uploads multipart com remoção do prefixo `/upload` e fallback de Deep Link. |
| **FEAT-17** | Validação E2E com as 22 Tabelas do TCC | Validação | M5 | `cmd/validator` | Script/teste E2E validando a conversão das 22 tabelas reais de `resources/tcc/capitulos/*.tex`. |
| **FEAT-18** | Validação E2E com os 4 Diagramas Mermaid | Validação | M5 | `cmd/validator` | Validação do pipeline de renderização e aspect ratio dos 4 diagramas reais de `resources/tcc/figuras/*.mmd`. |
| **FEAT-19** | Validação E2E de Sincronia do `main.pdf` | Validação | M5 | `cmd/validator` | Validação do hash MD5 real do arquivo de 1.098.114 bytes contra o ID `1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH`. |
| **FEAT-20** | Auditoria de Entrelinhas, Margens e Seções (`lint_document_abnt`) | Refactor Tool | M6 | `internal/gdocs` | Audita `LineSpacing` (150 corpo, 100 tabelas/citações), `spaceAbove/Below`, margens de página A4 e quebra de seção/cabeçalho. |
| **FEAT-21** | Enriquecimento Geométrico de `inspect_doc_layout` | Refactor Tool | M6 | `internal/gdocs` | Expõe `LineSpacingPercent`, `DocumentGeometry` (margens 3/3/2/2 cm) e segmentação de `Sections` para modelos de IA. |
| **FEAT-22** | Preservação Estrutural em `read_doc_content` | Refactor Tool | M6 | `internal/gdocs` | Preserva e exibe marcadores `[PAGE_BREAK]`, `[SECTION_BREAK]` e modo opcional com anotações de estilo para agentes LLM. |
| **FEAT-23** | Formatação de Tabelas e Preservação de Citações Longas | Refactor Tool | M6 | `internal/gdocs` | Formata células de tabela para 1,0 simples / 10pt e preserva blocos de citação direta longa (4,0 cm / 10pt / 1,0) em `format_document_abnt`. |
| **FEAT-24** | Saneamento Inteligente de Quebras e Calibração de Capa | Refactor Core | M6 | `internal/gdocs` | Elimina quebras espúrias em transições de página em `abnt_writer.go`, substitui dezenas de `\n` por `SpaceAbove` e trata erros em `cover_generator.go`. |

---

## 3. Milestones Detalhados

### Milestone 1 (M1): Correções de Formatação e Estilos ABNT
- **Objetivo Central:** Eliminar 100% dos efeitos colaterais de Bold Bleeding, List/Bullet Bleeding e inconsistências de parágrafo, padronizando o motor de mutação segundo a ABNT NBR 14724.
- **Entregáveis Técnicos:**
  1. Criação do componente `internal/gdocs/builder.go` (`BatchUpdateBuilder`) com suporte à ordenação reversa.
  2. Implementação da função pura `SanitizeText(raw string) string` com normalização canônica Unicode NFC (`norm.NFC`) e cálculo seguro `UTF16Length(s string) int64`.
  3. Correção cirúrgica de `internal/gdocs/writer.go:292` permitindo aplicação explícita de `NORMAL_TEXT`.
  4. Implementação da nova ferramenta MCP `replace_block_abnt` incorporando:
     - Clamping defensivo de fim de documento (`EndIndex = min(EndIndex, body.EndIndex - 1)`).
     - Anulação do recuo residual de lista via `indentStart: { magnitude: 0, unit: "PT" }` incluído na máscara `fields`.
     - Guarda preventiva para substituição com texto vazio (`new_text == ""`), emitindo unicamente `DeleteContentRangeRequest`.
  5. Implementação da ferramenta `apply_abnt_page_margins` via `UpdateDocumentStyleRequest` com margens 3/3/2/2 cm ($85,039\text{ pt}$ / $56,693\text{ pt}$).
- **Critérios de Aceite:**
  - Nenhuma substituição de título por texto de parágrafo gera texto em negrito indesejado.
  - A substituição sobre itens de lista remove marcadores e anula recuos residuais de lista sem gerar bullets vazios (`* `).
  - Parágrafos formatados recebem recuo de primeira linha de 35,433 pt, entrelinhas 150 e alinhamento `JUSTIFIED`.
  - Operações de substituição no último parágrafo do documento não disparam erro HTTP 400.

### Milestone 2 (M2): Tabelas Nativas e Diagramas Mermaid
- **Objetivo Central:** Substituir marcadores de texto provisórios por elementos estruturais nativos do Google Docs, renderizando tabelas IBGE/ABNT e diagramas conceituais com alta definição visual.
- **Entregáveis Técnicos:**
  1. Implementação do parser de sintaxe LaTeX para matriz de dados (`internal/latex/table_parser.go`), normalizando `\addlinespace`, macros `\textcite`, caracteres escapados `\%` e fórmulas inline (`$2 \times 2$`).
  2. Implementação do padrão Two-Step para tabelas em `internal/gdocs/table_writer.go`:
     - Passo 1: Inserção da Legenda (`Caption`) no topo em `TargetIndex` e criação da tabela vazia via `InsertTableRequest` no offset ajustado.
     - Passo 2: Leitura de coordenadas via `GetDocument` capturando `CellIndex[r][c]`, `tableEndIndex` e a revisão `doc.RevisionId`.
     - Passo 3: Lote reverso de preenchimento com inserção da Fonte (`Source`) em `tableEndIndex` como a PRIMEIRA operação adicionada ao lote (maior offset), seguida do preenchimento decrescente das células $(R-1, C-1)$ a $(0, 0)$, estilização tabular via `UpdateTableCellStyleRequest` (`TableStartLocation` e `TableRange`), trava de concorrência com `WriteControl.RequiredRevisionId` e retry automático (Passo 2 + Passo 3) em caso de conflito.
     - Emissão formal da struct de retorno `TableInsertionResult`.
  3. Implementação da ferramenta MCP `insert_native_table`.
  4. Implementação do pipeline de conversão de Mermaid para PNG 300 DPI e hospedagem no Drive (`internal/gdocs/mermaid_writer.go`), com decodificação de dimensões via `image.DecodeConfig`, import anônimo `_ "image/png"`, guarda preventiva contra divisão por zero (`cfg.Width <= 0 || cfg.Height <= 0`), enquadramento bidimensional (2D Bounding Box: $\text{maxWidth} = 450.0\text{ pt}$, $\text{maxHeight} = 620.0\text{ pt}$ via escala proporcional $\min(1.0, \min(450.0/W, 620.0/H))$) e fallback para Deep Link textual caso a permissão pública seja bloqueada por DLP corporativo.
  5. Implementação da ferramenta MCP `insert_mermaid_diagram` via `InsertInlineImageRequest`.
- **Critérios de Aceite:**
  - A tabela criada no Docs não possui bordas verticais nas extremidades (norma IBGE).
  - O preenchimento das células ocorre em lote atômico reverso com inserção antecipada da Fonte em `tableEndIndex`, sem corrupção posicional e com trava de concorrência.
  - Os diagramas Mermaid preservam sua proporção visual nativa e são enquadrados na caixa 2D de $450 \times 620\text{ pt}$ sem transbordar a página A4.

### Milestone 3 (M3): Sincronização de PDF e Integridade MD5
- **Objetivo Central:** Estabelecer uma camada criptográfica e idempotente de sincronização entre o binário local do TCC e a nuvem do Google Drive.
- **Entregáveis Técnicos:**
  1. Implementação do módulo `internal/gdocs/hash_validator.go` com verificação prévia do marcador `%%EOF` nos últimos 1024 bytes (ISO 32000-1 Seção 7.5.5, blindagem contra compilação incompleta), cálculo de MD5 por streaming de bytes (`crypto/md5`) e validação de integridade concorrente contra TOCTOU (`copiedBytes == stat.Size()`).
  2. Implementação da ferramenta MCP `check_pdf_sync` para diagnóstico instantâneo de paridade.
  3. Refatoração completa de `internal/gdocs/upload.go` (`SyncPDFToDrive`) incorporando verificação de `%%EOF`, prevenção contra TOCTOU, verificação prévia do `remoteFile.Md5Checksum` e política `keepRevisionForever: true` para garantia de histórico de versões e viabilização de rollback seguro.
  4. Alerta crítico e abortamento caso o hash pós-upload divirja do hash local computado.
- **Critérios de Aceite:**
  - Se o PDF local for idêntico ao remoto, o upload é suprimido e a ferramenta retorna status `IN_SYNC`.
  - Tentativas de sincronizar PDFs em processo de compilação sem marcador `%%EOF` são abortadas preventivamente com mensagem de erro descritiva.
  - Se os hashes divergirem, o upload é realizado via `files.update` com política perene de revisão e confirmação pós-upload.

### Milestone 4 (M4): Framework de Testes com Mocks e Refatoração de Interfaces
- **Objetivo Central:** Eliminar o acoplamento rígido com o SDK do Google, estabelecendo 100% de cobertura de testes automatizados executáveis sem credenciais reais e sem acesso à internet.
- **Entregáveis Técnicos:**
  1. Criação das interfaces de serviço `DocsService` e `DriveService` em `internal/gdocs/interfaces.go`, incluindo explicitamente o método `ListFiles` para desacoplar a ferramenta de listagem `internal/gdocs/listing.go`.
  2. Implementação do adapter concreto `googleDocsAdapter` e `googleDriveAdapter` preservando compatibilidade de produção.
  3. Criação de infraestrutura de testes mock em `internal/gdocs/testing_mock.go` utilizando `httptest.Server` com multiplexador de rotas atendendo `/v1/documents/`, `/drive/v3/files/` e a rota multipart `/upload/drive/v3/files/` tratando `multipart/related` e normalizando o path sem o prefixo `/upload`.
  4. Criação de testes unitários para cálculo UTF-16 com NFC, parsing de tabelas com macros ABNT, detecção de colisões, checagem `%%EOF` (ISO 32000-1) e TOCTOU, e validação de URLs (`reader_test.go`, `writer_test.go`, `upload_test.go`, `table_test.go`, `mermaid_writer_test.go`).
  5. Configuração da suíte de execução em lote via `go test ./...`.
- **Critérios de Aceite:**
  - `go test ./...` executa com sucesso em tempo $< 5\text{ segundos}$ sem exigir `credentials.json`.
  - Suíte de mocks cobre operações de upload de mídia multipart sem erros 404 de roteamento.

### Milestone 5 (M5): Validação E2E com Artefatos do TCC
- **Objetivo Central:** Auditar e comprovar a aderência do servidor Go aos 6 capítulos, 22 tabelas e 4 diagramas reais do TCC presentes em `resources/tcc/`.
- **Entregáveis Técnicos:**
  1. Desenvolvimento de script de validação E2E (`cmd/validator/main.go` ou suíte de testes de validação).
  2. Teste de parsing e formatação das 22 tabelas mapeadas (Tabela 1 a Tabela 22 do TCC).
  3. Teste de renderização, dimensões e aspect ratio dos 4 diagramas Mermaid (`arquitetura.mmd`, `casos_uso.mmd`, `procedimento_experimental.mmd`, `sequencia.mmd`).
  4. Teste de conferência do PDF real compilado (`resources/tcc/main.pdf`).
  5. Emissão do relatório formal de encerramento e validação de aceitação.
- **Critérios de Aceite:**
  - Todas as 22 tabelas são analisadas e convertidas com dimensões exatas de linhas e colunas.
  - Todos os 4 diagramas geram metadados compatíveis com o Google Docs e preservação de proporção.
  - Hash MD5 do PDF compilado local é confrontado com a referência do projeto.

### Milestone 6 (M6): Conformidade Estrita ABNT (Entrelinhas, Margens, Quebras e Auditoria)
- **Objetivo Central:** Sanar em definitivo as causas raízes de falsos negativos de auditoria e falhas de diagramação ABNT (NBR 14724, NBR 6023, NBR 6024, NBR 6027, NBR 10520) no servidor MCP e no documento de teste.
- **Entregáveis Técnicos:**
  1. **Auditoria no Linter (`lint_document_abnt`):**
     - Criação do checker `checkLineSpacing` comparando `LineSpacing` (150 corpo, 100 tabelas/citações/legendas) e populando `LineSpacingAnomalies`.
     - Criação do checker `checkParagraphSpacing` para validação de `SpaceAbove=0` e `SpaceBelow=0` no corpo textual.
     - Criação do checker `checkPageMargins` auditando `DocumentStyle` (3cm sup/esq, 2cm inf/dir).
     - Remoção da blindagem `statePreTextual` para parágrafos vazios espúrios em transições de página.
  2. **Enriquecimento Estrutural do Inspector (`inspect_doc_layout`):**
     - Inclusão do campo `LineSpacingPercent` na struct `ElementLayoutDetail`.
     - Inclusão das structs `DocumentGeometryDetail` (margens e dimensões da folha A4) e `SectionLayoutDetail` no relatório `InspectLayoutReport`.
  3. **Preservação e Visibilidade no Reader (`read_doc_content`):**
     - Inserção dos marcadores explícitos `[PAGE_BREAK]` e `[SECTION_BREAK]` no Markdown de saída.
     - Suporte ao parâmetro `include_metadata: true` com tags estruturais para inspeção por agentes de IA.
  4. **Formatação ABNT Refatorada (`format_document_abnt`):**
     - Aplicação de margens ABNT via `UpdateDocumentStyleRequest` em pontos DTP ($85,039\text{ pt}$ / $56,693\text{ pt}$).
     - Varredura de células de tabelas nativas aplicando espaçamento 1,0 simples e fonte 10pt.
     - Mecanismo determinístico de detecção e preservação de citações diretas longas ($4,0\text{ cm} = 113,386\text{ pt}$, 10pt, 1,0 simples).
  5. **Saneamento e Calibração de Quebras (`cover_generator.go` e `abnt_writer.go`):**
     - Refatoração de `HealDocTextAndBlanks` para sanear parágrafos vazios isolados pós-`PageBreak` (`consecutiveBlanks >= 1`).
     - Substituição de dezenas de `\n` mecânicos por espaçamento vertical paramétrico (`SpaceAbove`) ou quebras de página diretas.
     - Correção do loop reverso de transição para o Resumo com tratamento explícito de erros na Docs API.
- **Critérios de Aceite:**
  - `lint_document_abnt` detecta com precisão parágrafos com entrelinhas 1,15 e quebras de linha espúrias após a folha de rosto.
  - `format_document_abnt` normaliza todo o documento para margens 3/3/2/2 cm, corpo 1,5 e células de tabelas para 1,0 sem desconfigurar citações longas.
  - `inspect_doc_layout` e `read_doc_content` fornecem visão 100% fiel da geometria e das quebras de página para modelos de IA.
  - 100% de aprovação na suíte de testes automatizados com mocks (`go test ./...`).

---

## 4. Interface Contracts Completos

### 4.1 Assinaturas em Go (Adapter Pattern)

```go
package gdocs

import (
	"context"
	"io"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
)

// DocsService abstrai a Google Docs API v1 para permitir injeção de dependências e mocks
type DocsService interface {
	GetDocument(ctx context.Context, documentID string) (*docs.Document, error)
	BatchUpdate(ctx context.Context, documentID string, req *docs.BatchUpdateDocumentRequest) (*docs.BatchUpdateDocumentResponse, error)
}

// DriveService abstrai a Google Drive API v3 para gestão de arquivos, mídia e comentários
type DriveService interface {
	GetFileMetadata(ctx context.Context, fileID string, fields string) (*drive.File, error)
	ListFiles(ctx context.Context, query string, fields string, pageSize int64, pageToken string) (*drive.FileList, error)
	UpdateFileMedia(ctx context.Context, fileID string, metadata *drive.File, media io.Reader, fields string) (*drive.File, error)
	CreateFile(ctx context.Context, metadata *drive.File, media io.Reader, fields string) (*drive.File, error)
	CreatePermission(ctx context.Context, fileID string, permission *drive.Permission) (*drive.Permission, error)
	ListComments(ctx context.Context, fileID string) ([]*drive.Comment, error)
	CreateReply(ctx context.Context, fileID string, commentID string, reply *drive.Reply) (*drive.Reply, error)
	ResolveComment(ctx context.Context, fileID string, commentID string) error
}

// TableInsertionResult representa o resultado auditável da criação e preenchimento de tabela
type TableInsertionResult struct {
	TableStartIndex int64  `json:"table_start_index"`
	TableEndIndex   int64  `json:"table_end_index"`
	RowsCount       int64  `json:"rows_count"`
	ColsCount       int64  `json:"cols_count"`
	Message         string `json:"message"`
}

// LaTeXTableData encapsula os metadados e matriz de dados de tabela para inserção
type LaTeXTableData struct {
	ID        string     `json:"id"`
	Label     string     `json:"label"`
	Caption   string     `json:"caption"`
	Source    string     `json:"source"`
	RowsCount int64      `json:"rows_count"`
	ColsCount int64      `json:"cols_count"`
	Headers   []string   `json:"headers"`
	Data      [][]string `json:"data"`
}

// TableService gerencia a orquestração Two-Step de tabelas nativas
type TableService interface {
	InsertNativeTable(ctx context.Context, docID string, index int64, table LaTeXTableData) (*TableInsertionResult, error)
}

// PDFSyncService orquestra a validação criptográfica e sincronização idempotente
type PDFSyncService interface {
	CheckSync(ctx context.Context, localPath string, driveFileID string) (*SyncReport, error)
	Sync(ctx context.Context, localPath string, driveFileID string) (*SyncReport, error)
}
```

---

### 4.2 Structs de Payload e Schemas das Ferramentas MCP

#### Tool 1: `replace_block_abnt`
- **Nome:** `replace_block_abnt`
- **Descrição:** Substitui um bloco contíguo de texto no Google Docs garantindo isolamento total de formatação (eliminação de Bold Bleeding, remoção de marcadores de lista/bullet e aplicação estrita do padrão de parágrafo ABNT).
- **Schema de Entrada (JSON Schema):**
```json
{
  "type": "object",
  "properties": {
    "document_url_or_id": {
      "type": "string",
      "description": "ID ou URL completa do documento Google Docs"
    },
    "start_index": {
      "type": "integer",
      "description": "Índice UTF-16 inicial do bloco a ser substituído"
    },
    "end_index": {
      "type": "integer",
      "description": "Índice UTF-16 final do bloco a ser substituído"
    },
    "expected_text": {
      "type": "string",
      "description": "Texto esperado na faixa para validação preventiva de conflito concorrente"
    },
    "new_text": {
      "type": "string",
      "description": "Novo texto a ser inserido no documento"
    },
    "apply_abnt_style": {
      "type": "boolean",
      "description": "Se verdadeiro, aplica recuo de 1,25 cm, entrelinhas 1,5 e alinhamento justificado",
      "default": true
    }
  },
  "required": ["document_url_or_id", "start_index", "end_index", "new_text"]
}
```
- **Struct Go:**
```go
type ReplaceBlockABNTInput struct {
	DocumentURLOrID string `json:"document_url_or_id"`
	StartIndex      int64  `json:"start_index"`
	EndIndex        int64  `json:"end_index"`
	ExpectedText    string `json:"expected_text,omitempty"`
	NewText         string `json:"new_text"`
	ApplyABNTStyle  bool   `json:"apply_abnt_style"`
}
```
> **Blindagens Defensivas Integradas:**
> 1. *Clamping de Fim de Documento:* `EndIndex = min(EndIndex, body.EndIndex - 1)` previne erro HTTP 400 no `\n` final inamovível.
> 2. *Anulação de Recuo Residual de Lista:* Aplicação de `indentStart: 0 pt` remove a indentação deixada por `DeleteParagraphBulletsRequest`.
> 3. *Guarda de Substituição Vazia:* Se `new_text == ""`, emite unicamente `DeleteContentRangeRequest`.
> 4. *Normalização Unicode NFC:* Executa `norm.NFC.String(s)` antes do cômputo de `UTF16Length`.

---

#### Tool 2: `apply_abnt_page_margins`
- **Nome:** `apply_abnt_page_margins`
- **Descrição:** Ajusta as margens globais de página do Google Docs para a norma ABNT (Superior: 3cm, Esquerda: 3cm, Inferior: 2cm, Direita: 2cm) e formato A4 ($595.276 \times 841.890\text{ pt}$).
- **Schema de Entrada:**
```json
{
  "type": "object",
  "properties": {
    "document_url_or_id": {
      "type": "string",
      "description": "ID ou URL do documento Google Docs"
    }
  },
  "required": ["document_url_or_id"]
}
```

---

#### Tool 3: `insert_native_table`
- **Nome:** `insert_native_table`
- **Descrição:** Insere uma tabela nativa estruturada no Google Docs com bordas no padrão acadêmico ABNT/IBGE (sem bordas verticais externas), orquestrando o Two-Step Approach com legenda no Passo 1, preenchimento reverso protegido por `WriteControl` e fonte no Passo 3.
- **Schema de Entrada:**
```json
{
  "type": "object",
  "properties": {
    "document_url_or_id": { "type": "string", "description": "ID ou URL do documento" },
    "insertion_index": { "type": "integer", "description": "Índice de inserção no fluxo do texto" },
    "caption": { "type": "string", "description": "Legenda da tabela posicionada no topo (ABNT)" },
    "source": { "type": "string", "description": "Fonte da tabela posicionada no rodapé (ABNT)" },
    "headers": {
      "type": "array",
      "items": { "type": "string" },
      "description": "Lista de títulos das colunas (cabeçalho em negrito)"
    },
    "rows": {
      "type": "array",
      "items": {
        "type": "array",
        "items": { "type": "string" }
      },
      "description": "Matriz bidimensional contendo os dados das linhas"
    }
  },
  "required": ["document_url_or_id", "insertion_index", "headers", "rows"]
}
```
- **Struct Go:**
```go
type InsertNativeTableInput struct {
	DocumentURLOrID string     `json:"document_url_or_id"`
	InsertionIndex  int64      `json:"insertion_index"`
	Caption         string     `json:"caption,omitempty"`
	Source          string     `json:"source,omitempty"`
	Headers         []string   `json:"headers"`
	Rows            [][]string `json:"rows"`
}
```
- **Struct Go de Retorno:** `TableInsertionResult` (contendo `TableStartIndex`, `TableEndIndex`, `RowsCount`, `ColsCount`, `Message`).

---

#### Tool 4: `insert_mermaid_diagram`
- **Nome:** `insert_mermaid_diagram`
- **Descrição:** Renderiza e injeta um diagrama Mermaid no Google Docs com enquadramento bidimensional (2D Bounding Box: $\text{maxWidth} = 450.0\text{ pt}$ e $\text{maxHeight} = 620.0\text{ pt}$ via escala proporcional $\min(1.0, \min(450.0/W, 620.0/H))$) utilizando `image.DecodeConfig` com import anônimo `_ "image/png"` e guarda preventiva contra divisão por zero, além de fallback para Deep Link textual caso a concessão de permissão pública seja bloqueada por DLP corporativo.
- **Schema de Entrada:**
```json
{
  "type": "object",
  "properties": {
    "document_url_or_id": { "type": "string", "description": "ID ou URL do documento" },
    "insertion_index": { "type": "integer", "description": "Índice de inserção no texto" },
    "diagram_name": { "type": "string", "description": "Nome base do diagrama em resources/tcc/figuras (ex: arquitetura)" },
    "caption": { "type": "string", "description": "Legenda da figura no topo" },
    "source": { "type": "string", "description": "Fonte da figura no rodapé" },
    "width_pt": { "type": "number", "description": "Largura da imagem em pontos (máximo ~450 pt)", "default": 420.0 }
  },
  "required": ["document_url_or_id", "insertion_index", "diagram_name"]
}
```

---

#### Tool 5: `check_pdf_sync`
- **Nome:** `check_pdf_sync`
- **Descrição:** Compara o hash criptográfico MD5 do PDF compilado localmente com o metadado `md5Checksum` do arquivo no Google Drive, verificando integridade sem tráfego de mídia, validando previamente o marcador final `%%EOF` (ISO 32000-1 Seção 7.5.5) e checando integridade concorrente pós-hashing (`copiedBytes == stat.Size()`) contra TOCTOU.
- **Schema de Entrada:**
```json
{
  "type": "object",
  "properties": {
    "local_pdf_path": { "type": "string", "description": "Caminho do arquivo PDF compilado local", "default": "resources/tcc/main.pdf" },
    "target_pdf_file_id": { "type": "string", "description": "ID do arquivo PDF no Google Drive", "default": "1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH" }
  },
  "required": ["local_pdf_path", "target_pdf_file_id"]
}
```
- **Struct Go de Resposta:** `SyncReport` (com campos `FileID`, `FileName`, `LocalPath`, `LocalSize`, `LocalMD5`, `RemoteSize`, `RemoteMD5`, `RemoteModified`, `Status`, `Message`).

---

#### Catálogo Completo das 17 Ferramentas MCP Registradas

O servidor unifica e disponibiliza 17 ferramentas MCP sob o protocolo JSON-RPC 2.0:
1. `list_available_documents`: Listagem e busca de documentos Google Docs e PDFs acessíveis no Drive.
2. `get_doc_skeleton`: Estrutura hierárquica e índices UTF-16 de cabeçalhos e parágrafos.
3. `read_doc_content`: Leitura estruturada em Markdown de faixas do documento.
4. `search_in_doc`: Busca por texto ou regex retornando offsets absolutos.
5. `list_doc_comments`: Listagem de comentários e revisões pendentes.
6. `reply_to_comment`: Resposta encadeada a comentários.
7. `resolve_comment`: Resolução formal de comentários colaborativos.
8. `apply_text_style`: Aplicação pontual de estilos de caractere (negrito, itálico, cores).
9. `search_and_replace_text`: Substituição global simples de termos.
10. `update_block_by_index`: Atualização de bloco posicional com normalização de estilo.
11. `sync_pdf_to_drive`: Sincronização inteligente de PDF com verificação de `%%EOF` (ISO 32000-1), mitigação TOCTOU e retenção de revisões.
12. `update_local_latex`: Atualização bidirecional de arquivos LaTeX com backup `.bak`.
13. `apply_abnt_page_margins` *(Nova - FEAT-01)*: Margens ABNT 3/3/2/2 cm em formato A4.
14. `replace_block_abnt` *(Nova - FEAT-02)*: Substituição anti-bleeding com reset para `NORMAL_TEXT`, `indentStart: 0 pt` e clamping.
15. `insert_native_table` *(Nova - FEAT-06)*: Inserção de tabela nativa IBGE/ABNT com Two-Step e `TableInsertionResult`.
16. `insert_mermaid_diagram` *(Nova - FEAT-09)*: Diagrama Mermaid rasterizado com enquadramento 2D Bounding Box ($450 \times 620\text{ pt}$) e fallback Deep Link.
17. `check_pdf_sync` *(Nova - FEAT-10)*: Diagnóstico criptográfico de paridade local/remota via MD5, marcador `%%EOF` e guarda TOCTOU.

---

### 4.3 Contratos de Lote `batchUpdate` e Atomicidade

#### Contrato 1: Sequência Canônica de Substituição Limpa (`replace_block_abnt`)
O payload atômico deve ser emitido segundo as seguintes regras determinísticas:
- **Clamping de Segurança:** `endIndex = min(endIndex, body.EndIndex - 1)`.
- **Caso Texto Vazio:** Se `len(sanitizedText) == 0`, o lote contém exclusivamente a requisição `[1] DeleteContentRangeRequest`.
- **Caso Padrão (5 Operações Atômicas Sequenciais):**
```
[1] DeleteContentRangeRequest (startIndex, endIndex)
[2] InsertTextRequest (location.index = startIndex, text = sanitizedText)
[3] DeleteParagraphBulletsRequest (startIndex, startIndex + lenUTF16)
[4] UpdateParagraphStyleRequest (startIndex, startIndex + lenUTF16, NORMAL_TEXT, lineSpacing=150, indentFirstLine=35.433 pt, indentStart=0 pt, fields="namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow")
[5] UpdateTextStyleRequest (startIndex, startIndex + lenUTF16, bold=false, italic=false, underline=false, fields="bold,italic,underline")
```

#### Contrato 2: Two-Step Table Approach e Sequenciamento Estrutural
1. **Passo 1 (Legenda + Estrutura Vazia):**
   - Se `caption != ""`: insere parágrafo centralizado de legenda em `TargetIndex` (`tableStartIndex = TargetIndex + UTF16Length(caption + "\n")`).
   - Insere a tabela via `InsertTableRequest` no `tableStartIndex`.
2. **Passo 2 (Inspeção e Travamento de Revisão):**
   - Executa `GetDocument(docID)` para capturar `CellIndex[r][c]`, `tableEndIndex` e `doc.RevisionId`.
3. **Passo 3 (Preenchimento Reverso, Estilização e Fonte com Trava de Revisão):**
   - Define `WriteControl: &docs.WriteControl{RequiredRevisionId: doc.RevisionId}` no `BatchUpdateDocumentRequest`.
   - Se `source != ""`: adiciona como **primeira operação** do lote atômico a inserção do parágrafo centralizado de fonte em `tableEndIndex` (maior offset absoluto), blindando as coordenadas das células contra o deslocamento provocado pelo preenchimento de seus textos.
   - Preenche células de $(R-1, C-1)$ até $(0, 0)$ em ordem decrescente estrita de offsets.
   - Aplica `UpdateTableCellStyleRequest` com `TableStartLocation: &docs.Location{Index: tableStartIndex}` e `TableRange`.
   - Política de Retry: Se a requisição falhar com erro HTTP 400 (`Invalid revision ID`), reexecuta o Passo 2 e reemite o Passo 3.

---

## 5. Code Layout Detalhado para o Projeto Go

A estrutura modular a seguir organiza as responsabilidades de forma desacoplada, facilitando a implementação por múltiplos agentes e garantindo 100% de testabilidade:

```text
c:\Projetos\golang\mcp-gdocs-latex\
├── cmd/
│   ├── mcp-server/
│   │   └── main.go                 # Entrypoint MCP (stdio, init de logs stderr, adapters e server)
│   └── validator/
│       └── main.go                 # Utilitário de validação E2E com artefatos do TCC
├── internal/
│   ├── gdocs/                      # Integração com Google Workspace (Docs v1 e Drive v3)
│   │   ├── auth.go                 # Carregamento de Service Account e credenciais
│   │   ├── builder.go              # BatchUpdateBuilder fluente para requisições atômicas
│   │   ├── comments.go             # Gestão de comentários no Drive (listar, responder, resolver)
│   │   ├── hash_validator.go       # Cálculo MD5 streaming, validação de %%EOF (ISO 32000-1) e mitigação TOCTOU
│   │   ├── interfaces.go           # Definição de interfaces de serviço (DocsService, DriveService)
│   │   ├── listing.go              # Listagem de arquivos Docs e PDFs
│   │   ├── mermaid_writer.go       # Upload de imagem de diagrama, enquadramento 2D (450x620 pt) e injeção inline no Docs
│   │   ├── reader.go               # Leitura estrutural de Docs e esqueleto em Markdown
│   │   ├── search.go               # Busca de padrões e regex com offsets UTF-16
│   │   ├── style_abnt.go           # Constantes e geradores de payload ABNT (recuo, margens)
│   │   ├── table_writer.go         # Two-Step table writer (criação, leitura, preenchimento reverso com Fonte em lote e WriteControl)
│   │   ├── upload.go               # Sincronização inteligente e idempotente de PDF
│   │   ├── utf16_utils.go          # Utilitários seguros de cálculo UTF-16 e sanitização de texto
│   │   ├── validator.go            # Extração de docID por Regex, allowlist e health-check
│   │   ├── writer.go               # Escrita, substituição de texto com reset de estilo
│   │   ├── builder_test.go         # Testes unitários do BatchUpdateBuilder
│   │   ├── hash_validator_test.go  # Testes de hash MD5, %%EOF e TOCTOU com arquivos sintéticos
│   │   ├── mock_server_test.go     # Servidor HTTP mock com httptest.Server multiplexado e rotas multipart para Docs e Drive
│   │   ├── style_abnt_test.go      # Testes de integridade das dimensões métricas ABNT
│   │   ├── table_writer_test.go    # Testes unitários e simulados do Two-Step table approach
│   │   ├── utf16_utils_test.go     # Testes de cálculo UTF-16 com acentos e emojis
│   │   └── writer_test.go          # Testes de cura de Bold Bleeding e List Bleeding
│   ├── latex/
│   │   ├── local_file.go           # Leitura e escrita de .tex com backup .bak automático
│   │   ├── table_parser.go         # Parser de ambientes tabularx e tabular do LaTeX
│   │   ├── local_file_test.go      # Testes unitários de manipulação de arquivos LaTeX
│   │   └── table_parser_test.go    # Testes do parser com as 22 tabelas do TCC
│   └── mcp/
│       ├── tools.go                # Catálogo e registro de todas as 17 ferramentas MCP (12 pré-existentes + 5 novas)
│       └── tools_test.go           # Testes de decodificação de schemas e validação de parâmetros
├── docs/
│   └── architecture_specification.md # Especificação técnica formal e cirúrgica (R1 a R5)
├── resources/
│   └── tcc/                        # Repositório LaTeX canônico de referência
├── .env.example                    # Template de variáveis de ambiente
├── go.mod / go.sum                 # Módulos Go 1.25.5
├── Makefile                        # Automação de compilação, testes e validações
└── PROJECT.md                      # Este documento (Project Pattern)
```

---

## 6. Governança de Dependências e Build

- **Linguagem:** Go `1.25.5`
- **Módulos Principais:**
  - `github.com/mark3labs/mcp-go v0.54.0` (Camada MCP Server)
  - `google.golang.org/api v0.280.0` (SDK Google Docs v1 e Google Drive v3)
  - `golang.org/x/oauth2 v0.36.0` (Autenticação Google Service Account)
- **Comandos de Automação (`Makefile`):**
  - `make build`: Compila o binário `bin/mcp-server.exe`.
  - `make test`: Executa todos os testes unitários e com mocks (`go test -v -race ./...`).
  - `make validate-tcc`: Executa a validação cruzada com os artefatos de `resources/tcc`.
  - `make fmt`: Formata o código Go segundo as convenções oficiais (`gofmt -s -w .`).
