# Análise Técnica Profunda da Base de Código Go (`mcp-gdocs-latex`)

**Data de Elaboração:** 29 de outubro de 2026  
**Investigador:** `explorer_codebase_1` (`teamwork_explorer`)  
**Workspace:** `c:\Projetos\golang\mcp-gdocs-latex`  
**Escopo:** Mapeamento arquitetural completo, diagnóstico técnico de falhas históricas de integridade/estilo e avaliação da infraestrutura de testes e qualidade.

---

## 1. Sumário Executivo

A investigação detalhada da base de código Go e dos artefatos históricos do projeto `mcp-gdocs-latex` revelou que o servidor MCP atual implementa uma ponte funcional básica entre o protocolo MCP e as APIs do Google Workspace (Docs v1 e Drive v3). No entanto, a base de código sofre de **fragilidades arquiteturais severas**, que explicam diretamente os problemas documentados de *bold bleeding* (efeito cascata de negrito), geração de bullets vazios, desalinhamento de listas, perda de espaçamentos e ausência de sincronização idempotente de PDF.

### Principais Constatações:
1. **Ausência Total de Testes:** Não existe sequer um arquivo `*_test.go` no projeto (0% de cobertura de testes). O código acopla-se a singletons globais de structs concretas do SDK do Google (`*docs.Service` e `*drive.Service`), dificultando mocks sem refatoração.
2. **Bug Crítico na Normalização de Estilo:** Em `internal/gdocs/writer.go:292`, existe uma cláusula de guarda `if bloco.ParagraphStyle != "NORMAL_TEXT" && bloco.ParagraphStyle != ""` que **impede explicitamente** a aplicação de `NORMAL_TEXT` nos parágrafos, fazendo com que textos injetados continuem retendo estilos de cabeçalho (`HEADING_1`, etc.).
3. **Mecanismo Cego de Substituição:** O `search_and_replace_text` invoca a API nativa `ReplaceAllTextRequest`, a qual herda passivamente a formatação de caractere do primeiro token substituído, propagando negrito indevidamente para parágrafos inteiros.
4. **Desconhecimento de Estrutura de Listas:** Nem o leitor (`reader.go`) nem o escritor (`writer.go`) tratam a propriedade `paragraph.Bullet`. Ao sobrescrever blocos contidos em itens de lista, o Docs interpreta quebras de linha (`\n`) como novos itens, gerando bullets vazios e listas infinitas. A remoção de listas requer `DeleteParagraphBulletsRequest`, requisição inexistente no código.
5. **Upload de PDF sem Validação de Integridade:** A ferramenta `sync_pdf_to_drive` possui apenas 31 linhas e executa um `svc.Files.Update` cego de mídia, sem calcular hash MD5 local, sem checar o `md5Checksum` do Drive v3 e sem garantir idempotência.
6. **Suporte Nulo a Tabelas e Imagens:** Não há implementação de `InsertTableRequest` nem de `InsertInlineImageRequest`. As tabelas e os diagramas Mermaid do TCC foram convertidos paliativamente em deep links ou texto cru ilegível.

---

## 2. Mapeamento Arquitetural Atual do Projeto Go

### 2.1 Estrutura de Diretórios e Arquivos

O projeto segue a convenção de layout padrão de Go para CLI/daemons, concentrando a lógica de negócio em pacotes sob `internal/` e o executável principal em `cmd/mcp-server/`. Não há pacote público `pkg/`.

```text
c:\Projetos\golang\mcp-gdocs-latex\
├── cmd/
│   └── mcp-server/
│       └── main.go                 # Entrypoint MCP (stdio, init de logs e server)
├── internal/
│   ├── gdocs/                      # Integração com APIs Google Workspace (Docs v1 e Drive v3)
│   │   ├── auth.go                 # Carregamento de Service Account e singletons do client
│   │   ├── comments.go             # Gestão de comentários no Drive (listar, responder, resolver, colisão)
│   │   ├── listing.go              # Varredura de Docs e PDFs no Drive
│   │   ├── reader.go               # Leitura de esqueleto e extração estrutural em Markdown
│   │   ├── search.go               # Busca de texto/regex com cálculo de offsets UTF-16
│   │   ├── upload.go               # Upload de mídia PDF no Drive
│   │   ├── validator.go            # Extração de docID por Regex, allowlist e health-check
│   │   └── writer.go               # Escrita, substituição de texto e formatação de blocos
│   ├── latex/
│   │   └── local_file.go           # Escrita segura de .tex local com backup (.bak)
│   └── mcp/
│       └── tools.go                # Catálogo e registro de tools no MCP Server
├── resources/
│   ├── Avaliação de uma Arquitetura... (2).md  # Export real do GDocs demonstrando falhas
│   ├── mcp_sync_history_and_improvements.md   # Histórico das dores e melhorias
│   ├── outras infromações.md                  # Requisitos anotados do TCC
│   └── tcc/                                  # Repositório LaTeX de referência
│       ├── capitulos/ (01 a 06 .tex)         # Capítulos acadêmicos com tabelas e figuras
│       ├── figuras/ (.mmd e .pdf)            # Diagramas Mermaid e PDFs compilados
│       ├── config.tex, main.tex, etc.
├── .env.example                    # Template de variáveis de ambiente
├── go.mod / go.sum                 # Módulos e dependências Go
├── Makefile                        # Automação de compilação multi-plataforma
└── README.md                       # Documentação do projeto
```

### 2.2 Dependências e Bibliotecas Externas (`go.mod`)

O arquivo `go.mod` declara dependências modernas:

```text
module github.com/JGustavoCN/mcp-gdocs-latex
go 1.25.5

require (
    github.com/mark3labs/mcp-go v0.54.0     // SDK Go oficial/comunitário do protocolo MCP
    golang.org/x/oauth2 v0.36.0             // Autenticação OAuth2 / Service Account
    google.golang.org/api v0.280.0          // SDK oficial da Google API (Docs v1 e Drive v3)
)
```

- **`github.com/mark3labs/mcp-go (v0.54.0)`:** Biblioteca central para comunicação JSON-RPC 2.0 via `stdio`. Suporta definição declarativa de ferramentas (`mcplib.NewTool`), tipagem de parâmetros (`WithString`, `WithInteger`, `WithBoolean`) e controle de ciclo de vida (`server.NewMCPServer`, `server.ServeStdio`).
- **`google.golang.org/api (v0.280.0)`:** Fornece os pacotes `google.golang.org/api/docs/v1` e `google.golang.org/api/drive/v3`. Os endpoints utilizados realizam operações de batch update, leitura de documento, gestão de comentários e upload de mídia.
- **`golang.org/x/oauth2 (v0.36.0)`:** Utilizado em conjunto com `google.golang.org/api/option` para autenticar requisições através de credenciais JSON de Service Account (`google.CredentialsFromJSON`).

### 2.3 Inicialização, Configuração e Autenticação

#### Entrypoint (`cmd/mcp-server/main.go`)
1. **Redirecionamento de Logs:** `log.SetOutput(os.Stderr)`. Esta linha é mandatória porque o protocolo MCP utiliza o `stdout` estritamente para o transporte JSON-RPC. Qualquer saída acidental para o `stdout` corromperia a comunicação com o cliente MCP.
2. **Criação do Servidor:** `s := server.NewMCPServer("TCC_DocsLaTeX_Bridge", "1.0.0", server.WithToolCapabilities(false))`.
3. **Registro das Ferramentas:** Chamada a `mcptools.RegisterTools(s)`.
4. **Transporte Stdio:** `server.ServeStdio(s)` bloqueia e processa mensagens JSON-RPC até que o cliente encerre o fluxo de entrada.

#### Configuração via Variáveis de Ambiente
- `ALLOWED_DOC_IDS`: Lista separada por vírgulas de IDs ou URLs de documentos permitidos. Se configurada, atua como allowlist preventiva em `internal/gdocs/validator.go:isDocAllowed`.
- `GOOGLE_APPLICATION_CREDENTIALS`: Caminho para o arquivo `credentials.json`. Se omitido, o sistema busca o arquivo no mesmo diretório do binário compilado (`filepath.Dir(exe)`).

#### Mecânica de Autenticação (`internal/gdocs/auth.go`)
- **Tipo de Credencial:** Google Cloud Service Account (chave privada JSON).
- **Escopos Utilizados:**
  - `https://www.googleapis.com/auth/documents` (`docs.DocumentsScope`): Permissão total de leitura e escrita em Google Docs.
  - `https://www.googleapis.com/auth/drive` (`drive.DriveScope`): Permissão total de leitura, comentários e upload no Drive.
- **Ciclo de Vida dos Clientes:**
  - Implementados via padrão **Singleton** protegido por `sync.Once` (`docsOnce` e `driveOnce`).
  - Funções de acesso: `GetDocsService(ctx)` e `GetDriveService(ctx)`.
  - **Limitação de Testabilidade:** Como os singletons retornam structs concretas ponteiro (`*docs.Service` e `*drive.Service`), não é possível injetar mocks sem instanciar servidores HTTP ou alterar a assinatura das funções para aceitar interfaces de domínio.

### 2.4 Catálogo de Ferramentas MCP Implementadas e Registradas

Apesar de comentários no código mencionarem 9 ferramentas e o `README.md` citar 10 ou 12, a inspeção em `internal/mcp/tools.go` confirma o registro exato de **12 ferramentas MCP**:

| # | Nome da Ferramenta | Parâmetros de Entrada (Schema) | Finalidade / Descrição | Implementação Interna |
|---|---|---|---|---|
| 1 | `list_available_documents` | Nenhum | Lista documentos GDocs e arquivos PDF compartilhados com a Service Account no Google Drive. | `gdocs.ListAvailableDocuments` em `listing.go` |
| 2 | `get_doc_skeleton` | `document_url_or_id` (string, req) | Retorna a estrutura hierárquica (títulos e blocos com highlight) com índices UTF-16 (`start_index`, `end_index`). | `gdocs.GetDocSkeleton` em `reader.go` |
| 3 | `read_doc_content` | `document_url_or_id` (string, req), `start_index` (int, opt), `end_index` (int, opt) | Lê o texto do documento e o converte para Markdown. Permite paginação por índices. | `gdocs.ReadDocContent` em `reader.go` |
| 4 | `search_in_doc` | `document_url_or_id` (string, req), `query` (string, req), `is_regex` (bool, opt) | Busca ocorrências no documento retornando índices absolutos e contexto do parágrafo. | `gdocs.SearchInDoc` em `search.go` |
| 5 | `list_doc_comments` | `document_url_or_id` (string, req) | Lista comentários no Drive, incluindo texto original ancorado, autor, data e respostas. | `gdocs.ListDocComments` em `comments.go` |
| 6 | `reply_to_comment` | `document_url_or_id` (string, req), `comment_id` (string, req), `reply_text` (string, req) | Insere resposta textual em um comentário existente via Drive API. | `gdocs.ReplyToComment` em `comments.go` |
| 7 | `resolve_comment` | `document_url_or_id` (string, req), `comment_id` (string, req) | Marca um comentário como resolvido (`action: "resolve"`). | `gdocs.ResolveComment` em `comments.go` |
| 8 | `apply_text_style` | `document_url_or_id` (string, req), `start_index` (int, req), `end_index` (int, req), `style` (string, req: BOLD/ITALIC/NONE) | Aplica ou remove negrito e itálico em uma faixa de índices sem apagar o texto. | `gdocs.ApplyTextStyle` em `writer.go` |
| 9 | `search_and_replace_text` | `document_url_or_id` (string, req), `old_text` (string, req), `new_text` (string, req) | Substitui globalmente ocorrências exatas de texto via `ReplaceAllTextRequest`. | `gdocs.SearchAndReplaceText` em `writer.go` |
| 10 | `update_block_by_index` | `document_url_or_id` (string, req), `payload_json` (string, req: `{start_index, end_index, expected_text, content: [...]}`) | Substituição atômica de intervalo com validação prévia de texto e injeção de trechos estilizados. | `gdocs.UpdateBlockByIndex` em `writer.go` |
| 11 | `sync_pdf_to_drive` | `local_pdf_path` (string, req), `target_pdf_file_id` (string, req) | Realiza upload direto do binário PDF local para o arquivo existente no Drive. | `gdocs.SyncPDFToDrive` em `upload.go` |
| 12 | `update_local_latex` | `filepath` (string, req), `new_content` (string, req) | Escreve conteúdo em arquivo `.tex` no disco local, gerando backup `.bak` automático. | `latex.UpdateLocalLatex` em `local_file.go` |

*Nota Arquitetural Adicional:* A função `MultiReplaceDocBlock` (em `internal/gdocs/writer.go:103-182`), desenvolvida para aplicar ordenação decrescente de chunks (*Offset Shift Rule*) e prevenir corrupção de índices em substituições múltiplas, **não está registrada no MCP** (`internal/mcp/tools.go`), constituindo código morto.

---

## 3. Diagnóstico Técnico dos Problemas Históricos e Recorrentes

### 3.1 O Efeito Cascata de Negrito (Bold Bleeding) e Perda de Estilos

#### Mecanismo em `search_and_replace_text`:
Em `internal/gdocs/writer.go:29-41`, a ferramenta envia a seguinte requisição para a API do Google Docs:

```go
req := &docs.BatchUpdateDocumentRequest{
    Requests: []*docs.Request{
        {
            ReplaceAllText: &docs.ReplaceAllTextRequest{
                ContainsText: &docs.SubstringMatchCriteria{
                    Text:      oldText,
                    MatchCase: true,
                },
                ReplaceText: newText,
            },
        },
    },
}
```

**Diagnóstico:** A especificação da Google Docs API define que `ReplaceAllTextRequest` herda passivamente todas as propriedades de `TextStyle` (como `bold: true`, tamanho de fonte e cor) do primeiro caractere casado no texto original (`oldText`).  
Quando o Agente anterior utilizou um título em negrito (por exemplo, `"1.1 Contextualização"` ou `"3.4 Métricas"`) como texto âncora para injetar uma seção inteira em `newText`, **todos os caracteres e parágrafos do texto injetado herdaram o negrito do título**. O arquivo exportado do Docs (`resources/Avaliação de uma Arquitetura Híbrida... (2).md`, linhas 49-68) comprova essa falha: centenas de linhas de objetivos específicos e listas inteiras ficaram presas em negrito indevido (`**1. Realizar consultas...**`).

#### Bug Crítico de Lógica em `update_block_by_index`:
Em `internal/gdocs/writer.go:292-305`, a ferramenta itera sobre os blocos estruturados e executa a aplicação de estilo de parágrafo:

```go
if bloco.ParagraphStyle != "NORMAL_TEXT" && bloco.ParagraphStyle != "" {
    requests = append(requests, &docs.Request{
        UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
            Range: &docs.Range{
                StartIndex: currentIndex,
                EndIndex:   currentIndex + utf16Length,
            },
            ParagraphStyle: &docs.ParagraphStyle{
                NamedStyleType: bloco.ParagraphStyle,
            },
            Fields: "namedStyleType",
        },
    })
}
```

**Diagnóstico:** Este trecho contém um **defeito de software evidente**: se o chamador passa explicitamente `ParagraphStyle: "NORMAL_TEXT"` com o objetivo de "curar" ou resetar a formatação de um cabeçalho para texto comum, o condicional `bloco.ParagraphStyle != "NORMAL_TEXT"` é falso e a instrução `UpdateParagraphStyleRequest` é **completamente omitida**. O Google Docs continua mantendo o estilo de parágrafo anterior (por exemplo, `HEADING_1`), mantendo fontes grandes, espaçamentos anômalos e herança indesejada.

---

### 3.2 Efeito Cascata de Listas (Bullets Vazios e Desalinhamento)

#### Mecanismo de Herança de Lista na Google Docs API:
No modelo estrutural do Google Docs, uma lista com marcadores ou numerada não é um contêiner separado; é um atributo do próprio parágrafo (`Paragraph.Bullet`).  
Quando um parágrafo que possui marcador de lista é substituído:
1. Ao enviar um texto com quebras de linha (`\n`), a Google Docs API interpreta cada quebra como a criação de um **novo parágrafo pertencente à mesma lista**.
2. Se o texto injetado contiver linhas em branco (duplo `\n\n`) para separar parágrafos ou subtítulos, o Google Docs cria um parágrafo de conteúdo vazio, porém contendo o marcador da lista. Surge então o **bullet vazio** (`* `).
3. Essa falha materializa-se exatamente na linha 226 do documento exportado (`Avaliação... (2).md`):
   ```markdown
   224: A construção do protótipo orienta-se pela arquitetura híbrida de injeção de contexto...
   225: 
   226: * 
   227: 
   228: 3.4.1 Trajetória de desenvolvimento e dilemas de engenharia
   ```

#### Cegueira do Servidor MCP:
- **No Leitor (`reader.go:171-240`):** `processParagraph` avalia apenas `paragraph.ParagraphStyle.NamedStyleType`. O campo `paragraph.Bullet` é totalmente ignorado. O servidor converte o texto para Markdown como se fossem parágrafos simples, omitindo o fato de que aquele trecho está dentro de uma lista. A IA opera no "escuro", acreditando estar alterando texto fluido quando está alterando itens de lista.
- **No Escritor (`writer.go`):** Não há suporte para `DeleteParagraphBulletsRequest`. Mesmo que o Agente queira remover o bullet de um bloco, não há chamada na API para desvincular o parágrafo da lista.

---

### 3.3 Quebras de Linhas Duplas e Espaços Supérfluos

#### Mecanismo:
1. No Google Docs, cada `Paragraph` termina obrigatoriamente com um caractere delimitador `\n` implícito em seus índices de caracteres.
2. Na conversão para Markdown em `internal/gdocs/reader.go:236-239`:
   ```go
   line := strings.Join(parts, "")
   if prefix != "" && strings.TrimSpace(line) != "" {
       line = prefix + strings.TrimSpace(line) + "\n"
   }
   sb.WriteString(line)
   ```
   Como os `TextRun` do Google Docs frequentemente já incluem `\n` ao final do texto, a manipulação de strings com concatenações e trims adicionais sem controle estrito de quebras de bloco provoca acúmulo de `\n\n` duplicados ou triplicados.
3. No envio de texto via `InsertTextRequest`, strings que contenham `\n\n` geram novos elementos estruturais de parágrafo vazios. Quando misturadas com estilos em cascata e listas, geram lacunas visuais inaceitáveis segundo as normas ABNT (as quais vedam linhas em branco avulsas entre parágrafos).

---

### 3.4 Sincronização de PDF (`sync_pdf_to_drive`)

A análise de `internal/gdocs/upload.go` revela a total fragilidade da ferramenta atual:

```go
func SyncPDFToDrive(ctx context.Context, fileID string, localPDFPath string) (string, error) {
    svc, err := GetDriveService(ctx)
    if err != nil {
        return "", fmt.Errorf("erro ao obter serviço do Drive: %w", err)
    }

    file, err := os.Open(localPDFPath)
    if err != nil {
        return "", fmt.Errorf("erro ao abrir PDF local: %w", err)
    }
    defer file.Close()

    // Faz APENAS o update da mídia no ID fornecido
    _, err = svc.Files.Update(fileID, &drive.File{}).Media(file).Do()
    if err != nil {
        return "", fmt.Errorf("erro ao atualizar o PDF no Drive: %w", err)
    }

    return fmt.Sprintf("SUCESSO: PDF atualizado silenciosamente no ID %s", fileID), nil
}
```

#### Falhas Críticas Identificadas:
1. **Ausência de Cálculo de Hash Local:** A ferramenta não calcula o hash MD5 (nem SHA256) do binário local compilado em disco (`resources/tcc/main.pdf`).
2. **Ausência de Checagem Remota via Google Drive API:** A Google Drive API v3 fornece nativamente o metadado `md5Checksum` no recurso `Files.Get(fileID).Fields("md5Checksum,size,modifiedTime")`. A ferramenta simplesmente não consulta o Drive antes de enviar o arquivo.
3. **Falta de Idempotência:** A cada chamada de `sync_pdf_to_drive`, o binário inteiro (geralmente superior a 1 MB) é transmitido pela rede via upload HTTP, mesmo que o arquivo local não tenha sofrido qualquer alteração.
4. **Relatório Pobre de Status:** Retorna uma string genérica (`"SUCESSO: PDF atualizado silenciosamente no ID %s"`), sem informar se houve modificação real, qual o hash de integridade verificado e qual o timestamp de atualização.

---

### 3.5 Suporte Atual para Tabelas e Imagens

#### Suporte a Tabelas:
- **Leitura (`reader.go:242-255`):** `processTable` limita-se a extrair o texto de cada célula e concatenar com barras simples (` | `) dentro de delimitadores arbitrários `[TABELA]` e `[/TABELA]`. Não preserva mesclagem de células, larguras, cabeçalhos nem alinhamentos.
- **Escrita (`writer.go`):** **Totalmente inexistente.** Não há nenhum suporte para `InsertTableRequest`, `InsertTableRowRequest` ou preenchimento de células estruturadas na Google Docs API.
- **Impacto no Projeto:** O repositório LaTeX (`resources/tcc/capitulos/`) possui dezenas de tabelas formais ABNT em ambiente `tabularx` (ex: `01-introducao.tex:256`, `03-metodologia.tex:21`, `04-desenvolvimento.tex:177`). Por impossibilidade técnica do servidor atual, os agentes anteriores simplesmente desistiram de renderizar tabelas no Docs e injetaram marcadores de texto:
  `[VER TABELA NO PDF COMPILADO — Tabela 1: Delimitação do escopo... : https://drive.google.com/file/d/1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH/view]`.

#### Suporte a Imagens e Diagramas Mermaid:
- **Leitura (`reader.go`):** Elementos do tipo `InlineObjectElement` ou `PositionedObject` são completamente ignorados no percurso estrutural de `doc.Body.Content`.
- **Escrita (`writer.go`):** **Totalmente inexistente.** Não existe suporte a `InsertInlineImageRequest`.
- **Diagramas Mermaid (`resources/tcc/figuras/*.mmd`):** O projeto dispõe de 4 diagramas conceituais cruciais:
  - `arquitetura.mmd` (compilado para `arquitetura.pdf`)
  - `casos_uso.mmd` (compilado para `casos_uso.pdf`)
  - `procedimento_experimental.mmd` (compilado para `procedimento_experimental.pdf`)
  - `sequencia.mmd` (compilado para `sequencia.pdf`)
- Atualmente, não há ferramenta para converter Mermaid (`.mmd`) em imagens bitmap (PNG/JPEG) nem para inseri-las de forma embutida no documento Google Docs. A tentativa inicial de colar arte ASCII provocou destruição de alinhamento pela proporcionalidade de fontes do GDocs, forçando a substituição por links externos.

---

## 4. Estrutura de Testes e Qualidade Existente

### 4.1 Diagnóstico da Suíte de Testes
- **Arquivos `*_test.go`:** **Nenhum.**
- **Cobertura de Testes Automatizados:** **0.0%.**
- **Testabilidade Atual:** A validação hoje depende integralmente de execução manual via cliente MCP em documentos Google Docs reais, exigindo chave de Service Account configurada na máquina de quem testa.

### 4.2 Análise de Acoplamento do Código Go
O código atual apresenta acoplamento direto com a biblioteca cliente oficial do Google:
1. `GetDocsService(ctx)` e `GetDriveService(ctx)` retornam os ponteiros concretos `*docs.Service` e `*drive.Service`.
2. As chamadas como `svc.Documents.BatchUpdate(documentID, req).Do()` dependem de structs internas do Google API que não implementam interfaces simples de Go e utilizam chamadas encadeadas (`.Documents.Get()`, `.Files.Update()`, `.Do()`).

### 4.3 Viabilidade e Facilidade de Implementação de Mocks em Go

Existem duas abordagens consolidadas e altamente viáveis para introduzir 100% de testabilidade na base de código Go:

#### Abordagem 1: Camada de Adaptação com Interfaces de Domínio (Recomendada)
Desacoplar a lógica das ferramentas MCP do SDK concreto através de interfaces de serviço:

```go
type DocsService interface {
    GetDocument(ctx context.Context, docID string) (*docs.Document, error)
    BatchUpdate(ctx context.Context, docID string, req *docs.BatchUpdateDocumentRequest) (*docs.BatchUpdateDocumentResponse, error)
}

type DriveService interface {
    GetFileMetadata(ctx context.Context, fileID string, fields string) (*drive.File, error)
    UpdateFileMedia(ctx context.Context, fileID string, body io.Reader) (*drive.File, error)
    ListComments(ctx context.Context, fileID string) ([]*drive.Comment, error)
    CreateReply(ctx context.Context, fileID, commentID string, reply *drive.Reply) (*drive.Reply, error)
}
```
*Vantagens:* Permite injetar mocks simples em testes unitários sem abrir portas de rede, testando validações de colisão de comentários, ordenação decrescente de chunks, lógica de substituição e controle de índices com precisão cirúrgica.

#### Abordagem 2: Servidor Mock HTTP com `httptest.Server`
O SDK oficial do Google em Go (`google.golang.org/api/option`) oferece suporte nativo para redirecionar endpoints e dispensar autenticação:

```go
ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // Retorna payloads JSON simulados de doc ou drive
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"documentId": "teste-123", "title": "Meu TCC"}`))
}))
defer ts.Close()

// Inicializa o client oficial apontando para o servidor local mock:
docsSvc, err := docs.NewService(ctx, option.WithEndpoint(ts.URL), option.WithoutAuthentication())
```
*Vantagens:* Testa o fluxo completo das ferramentas MCP contra respostas reais da API (incluindo testes de erro 403, 404, conflitos e validação de payloads `BatchUpdate`) de forma 100% autônoma, sem internet e sem credenciais reais de Service Account.

#### Abordagem 3: Testes de Unidade Imediatos em Funções Puras
Mesmo sem alterar a assinatura do SDK do Google, as seguintes funções da base de código já podem receber testes unitários puros hoje mesmo:
1. `gdocs.ExtractDocID`: Testar extração via URLs completas com parâmetros, IDs puros e strings vazias.
2. `internal/latex.UpdateLocalLatex`: Testar gravação segura, rejeição de extensões não-`.tex`, verificação de diretório e criação de backup `.bak`.
3. Cálculo de comprimentos em UTF-16 (`utf16.Encode([]rune(text))`): Garantir que caracteres multibyte (emojis, acentos da língua portuguesa como `ã`, `ç`, `é`) não desalinhem os offsets do Google Docs.
4. `CheckCommentCollision`: Testar lógica de detecção de colisão por substring ou limite de palavras contra strings normalizadas.

---

## 5. Recomendações Técnicas para o Plano Arquitetural

Com base no diagnóstico da base de código, as seguintes intervenções técnicas são mandatórias na arquitetura de evolução do servidor Go:

1. **Correção Imediata do Bug de Estilo de Parágrafo:**  
   Em `writer.go`, remover a condição excludente `if bloco.ParagraphStyle != "NORMAL_TEXT"` e garantir que quando `bloco.ParagraphStyle` for `NORMAL_TEXT`, o payload `UpdateParagraphStyleRequest` seja efetivamente despachado com `NamedStyleType: "NORMAL_TEXT"`, além de zerar os espaçamentos indevidos.
2. **Nova Tool Atômica `search_and_replace_clean` (ou parâmetro `clear_formatting`):**  
   Encapsular no mesmo `BatchUpdateDocumentRequest` a substituição do texto e a aplicação imediata de `UpdateTextStyleRequest` (`bold: false, italic: false`) e `UpdateParagraphStyleRequest` (`NORMAL_TEXT`) sobre a faixa alterada, eliminando requisições manuais e curando o *bold bleeding* nativamente.
3. **Mecanismo de Desvinculação de Listas (`DeleteParagraphBulletsRequest`):**  
   Implementar suporte a `DeleteParagraphBulletsRequest` para que, ao sobrescrever ou normalizar um parágrafo que era item de lista, o marcador seja removido na API do Google Docs, eliminando os bullets vazios.
4. **Verificação de Integridade por Hash MD5 no `sync_pdf_to_drive`:**  
   Refatorar a ferramenta para:
   - Calcular localmente o hash MD5 dos bytes do PDF compilado.
   - Consultar o metadado `md5Checksum` do arquivo remoto no Google Drive via `svc.Files.Get(fileID).Fields("md5Checksum").Do()`.
   - Comparar os hashes: se idênticos, reportar status `"IDENTICO: O PDF do Google Drive já está sincronizado com a compilação local (MD5: <hash>)"`, abortando o upload desnecessário.
   - Se divergentes, realizar o upload via `svc.Files.Update` e retornar o novo hash confirmado.
5. **Implementação de Ferramentas Nativas para Tabelas (`InsertTableRequest`):**  
   Criar ferramenta para receber especificações de tabelas (linhas, colunas, conteúdo das células e alinhamentos) e gerar as chamadas `InsertTableRequest` e `InsertTextRequest` em coordenadas de células (`TableCellLocation`).
6. **Pipeline de Inserção de Imagens e Diagramas Mermaid:**  
   Especificar ferramenta MCP para carregar diagramas `.mmd` ou imagens compiladas em `resources/tcc/figuras/`, realizar upload temporário no Google Drive para obter URL acessível (ou usar API de imagem) e emitir `InsertInlineImageRequest` no documento Google Docs.
7. **Suíte de Testes Automatizados com Mocks:**  
   Criar estrutura de testes sob `internal/gdocs/*_test.go` utilizando `httptest.Server` para simular as APIs Docs e Drive v1/v3, assegurando cobertura para todos os cenários de manipulação de texto, integridade de hash e validação de erros.
