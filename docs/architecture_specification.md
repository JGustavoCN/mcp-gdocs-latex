# Especificação Técnica e Arquitetural do Servidor Go: `mcp-gdocs-latex`

**Documento:** Especificação Técnica Formal de Engenharia  
**Referência do Projeto:** `mcp-gdocs-latex`  
**Autor:** `worker_spec_3` (`teamwork_worker`) — Revisão e Polimento Final (Iteração 3)  
**Data:** 29 de Outubro de 2026  
**Status:** Definitivo e Homologado para Implementação (Iteração 3)  
**Ambiente:** Go 1.25.5 | `mark3labs/mcp-go` v0.54.0 | `google.golang.org/api` v0.280.0  

---

## Sumário
1. [Requisito 1 (R1): Especificação das APIs Google Workspace e Servidor MCP](#1-requisito-1-r1-especificação-das-apis-google-workspace-e-servidor-mcp)
   - 1.1 Endpoints, Métodos HTTP e Restrições de Quota
   - 1.2 Semântica de Transação Atômica do `documents.batchUpdate`
   - 1.3 Aritmética de Índices e Unidades UTF-16 no Google Docs
   - 1.4 Padrões de Servidor MCP com `mark3labs/mcp-go`
2. [Requisito 2 (R2): Arquitetura de Formatação, Estilos ABNT e Anti-Bleeding](#2-requisito-2-r2-arquitetura-de-formatação-estilos-abnt-e-anti-bleeding)
   - 2.1 Causa Raiz e Eliminação Definitiva do Bold Bleeding
   - 2.2 Causa Raiz e Eliminação Definitiva do List / Bullet Bleeding
   - 2.3 Correção Cirúrgica do Defeito em `internal/gdocs/writer.go:292`
   - 2.4 Sequência Canônica de Requisições do Lote Atômico
   - 2.5 Parametrização Métrica das Normas ABNT (NBR 14724 e NBR 6023)
3. [Requisito 3 (R3): Arquitetura para Tabelas Nativas e Diagramas Mermaid](#3-requisito-3-r3-arquitetura-para-tabelas-nativas-e-diagramas-mermaid)
   - 3.1 O Padrão Two-Step Approach para Criação e Preenchimento de Tabelas
   - 3.2 Estilização ABNT/IBGE sem Bordas Verticais
   - 3.3 Mapeamento das 22 Tabelas LaTeX do TCC para a Estrutura Go
   - 3.4 Pipeline de Renderização e Inserção de Diagramas Mermaid
4. [Requisito 4 (R4): Verificação Criptográfica e Sincronização de PDF](#4-requisito-4-r4-verificação-criptográfica-e-sincronização-de-pdf)
   - 4.1 Comportamento do Metadado `md5Checksum` no Google Drive v3
   - 4.2 Especificação do Arquivo Canônico do TCC (`main.pdf`)
   - 4.3 Especificação Algorítmica da Sincronização Idempotente
5. [Requisito 5 (R5): Plano de Implementação, Mocks e Estratégia de Testes](#5-requisito-5-r5-plano-de-implementação-mocks-e-estratégia-de-testes)
   - 5.1 Arquitetura Desacoplada: Adapter Pattern e Interfaces de Serviço
   - 5.2 Framework de Testes Automatizados com `httptest.Server`
   - 5.3 Matriz de Cobertura de Testes Unitários e de Integração
   - 5.4 Procedimento de Validação E2E com Artefatos Reais do TCC

---

## 1. Requisito 1 (R1): Especificação das APIs Google Workspace e Servidor MCP

### 1.1 Endpoints, Métodos HTTP e Restrições de Quota

A integração do servidor `mcp-gdocs-latex` opera sobre duas APIs REST do ecossistema Google Workspace, autenticadas via Google Cloud Service Account com credenciais em formato JSON:

| API | Endpoint Base | Escopos OAuth 2.0 Requeridos |
|---|---|---|
| **Google Docs API (v1)** | `https://docs.googleapis.com/v1/documents` | `https://www.googleapis.com/auth/documents` |
| **Google Drive API (v3)** | `https://www.googleapis.com/drive/v3` | `https://www.googleapis.com/auth/drive` |

#### Endpoints Críticos Utilizados:
1. `GET /v1/documents/{documentId}`: Obtém a representação estrutural completa da árvore do documento (corpo, parágrafos, tabelas, objetos posicionados e estilos). Usado pelo leitor estrutural (`reader.go`) e no passo intermediário do preenchimento de tabelas.
2. `POST /v1/documents/{documentId}:batchUpdate`: Canal exclusivo de mutação do documento. Recebe um array contíguo de `Request` e retorna um array pareado de `Response`.
3. `GET /drive/v3/files/{fileId}`: Recupera metadados essenciais de arquivos, especificamente `fields="id, name, md5Checksum, size, modifiedTime, webContentLink, webViewLink"`.
4. `PATCH /upload/drive/v3/files/{fileId}?uploadType=multipart`: Atualiza os bytes binários de um arquivo existente (como o PDF de TCC compilado).
5. `POST /upload/drive/v3/files?uploadType=multipart`: Cria novos arquivos de mídia binária (imagens PNG geradas a partir de diagramas Mermaid).
6. `POST /drive/v3/files/{fileId}/permissions`: Gerencia compartilhamento temporário de recursos do Drive.
7. `GET /drive/v3/files/{fileId}/comments` e `POST .../replies`: Gestão de comentários colaborativos e revisão textual.

#### Quotas e Limites de Vazão:
- Google Docs API v1: Limite padrão de 300 requisições de leitura por minuto por projeto e 300 requisições de escrita por minuto por projeto. O agrupamento de até centenas de operações atômicas em uma única chamada `batchUpdate` é mandatório para evitar estouro de taxa HTTP 429 (`RESOURCE_EXHAUSTED`).
- Google Drive API v3: Limite padrão de 1.200 requisições por minuto por projeto.

---

### 1.2 Semântica de Transação Atômica do `documents.batchUpdate`

O endpoint `batchUpdate` possui garantias formais de transação ACID (especificamente Atomicidade e Consistência) no backend do Google:

1. **Rollback Total em Caso de Falha:** Todas as operações presentes no array `requests: []*docs.Request` são processadas na ordem exata de seus índices de array. Caso qualquer operação viole uma restrição (por exemplo, um índice fora dos limites, exclusão do caractere `\n` final do corpo ou má formação de máscara de campos), o lote inteiro é abortado imediatamente e nenhuma mutação é persistida no histórico do documento.
2. **Dependência Posicional Progressiva:** Cada operação muta o documento e afeta as coordenadas UTF-16 das operações subsequentes no mesmo lote. Por essa razão, operações concorrentes sobre diferentes partes do texto devem ser submetidas em **ordem decrescente de índice de início** (*Reverse Index Order Rule*).
3. **Máscara de Campos (`FieldMask`):** Operações de atualização de estilo (`UpdateTextStyleRequest`, `UpdateParagraphStyleRequest`, `UpdateDocumentStyleRequest`, `UpdateTableCellStyleRequest`) **exigem obrigatoriamente** o preenchimento do campo `fields`. Se `fields` for omitido ou enviado vazio (`""`), a Google Docs API rejeita a requisição com erro HTTP 400 ou ignora a modificação silenciosamente.

---

### 1.3 Aritmética de Índices e Unidades UTF-16 no Google Docs

O modelo de dados do Google Docs é estritamente indexado em **unidades de código UTF-16 (UTF-16 Code Units)**, começando no índice `0` (delimitador inicial reservado) até `body.EndIndex`.

#### Divergência Crítica entre Go e Google Docs:
Em Go, uma `string` é uma sequência nativa de bytes codificados em UTF-8. A função padrão `len(str)` retorna o número de **bytes**, não caracteres ou code units UTF-16:
- Caracteres ASCII (`a`, `Z`, `1`): 1 byte UTF-8 = 1 unidade UTF-16.
- Acentos e caracteres latinos (`ã`, `é`, `ç`): 2 bytes UTF-8 = **1 unidade UTF-16**.
- Caracteres matemáticos e símbolos especiais (`≠`, `≤`, `→`): 3 bytes UTF-8 = **1 unidade UTF-16**.
- Caracteres suplementares e Emojis (`🚀`, `📘`): 4 bytes UTF-8 = **2 unidades UTF-16** (par substituto / *surrogate pair*).

```
Exemplo Prático: String "Atenção 🚀"
- Go len("Atenção 🚀"): 13 bytes
- Runas len([]rune("Atenção 🚀")): 9 runas
- Google Docs API UTF-16: 10 code units ('A', 't', 'e', 'n', 'ç', 'ã', 'o', ' ', 0xD83D, 0xDE80)
```

#### Implementação Obrigatória em Go:
Para garantir que os índices enviados nas requisições `Range` e `Location` coincidam com os offsets internos do Google Docs, deve-se utilizar exclusivamente a biblioteca `unicode/utf16`:

```go
package gdocs

import "unicode/utf16"

// UTF16Length calcula com exatidão matemática o comprimento da string em code units UTF-16
func UTF16Length(s string) int64 {
	return int64(len(utf16.Encode([]rune(s))))
}
```

#### Restrição do Delimitador Final de Parágrafo e Clamping Defensivo:
O último elemento estrutural do `body.Content` de qualquer Google Doc é delimitado por um parágrafo que termina com um caractere `\n` inamovível (`body.EndIndex - 1`). Qualquer tentativa de submeter uma `DeleteContentRangeRequest` cujo `endIndex` alcance ou ultrapasse `body.EndIndex` é sumariamente rejeitada pela API com a mensagem:
`400 Bad Request: The content range must not include the final newline.`

**Regra de Clamping Defensivo Obrigatório no Backend:**
Para blindar o servidor contra essa falha quando ferramentas como `replace_block_abnt` ou clientes LLM utilizarem faixas derivadas de `get_doc_skeleton`, o construtor de requisições (`BatchUpdateBuilder`) deve impor um truncamento defensivo automático:
```go
// Clamping defensivo do delimitador final inamovível
if req.EndIndex >= bodyEndIndex {
    req.EndIndex = bodyEndIndex - 1
}
```

---

### 1.4 Padrões de Servidor MCP com `mark3labs/mcp-go` e Catálogo de Ferramentas

O servidor adota o SDK `github.com/mark3labs/mcp-go` v0.54.0. As seguintes regras arquiteturais devem ser cumpridas:

1. **Isolamento de Stdio:** Todo log de diagnóstico, depuração ou erro deve ser direcionado estritamente para `os.Stderr` (`log.SetOutput(os.Stderr)`). A saída padrão `os.Stdout` é de uso exclusivo do protocolo de transporte JSON-RPC 2.0.
2. **Tratamento Semântico de Falhas nas Ferramentas:**
   - Em caso de falhas de negócio ou validações de API, o handler da tool **não deve retornar erro Go** `(nil, err)`, pois isso encerra a sessão MCP com o cliente.
   - O handler deve retornar `mcplib.NewToolResultError(mensagem), nil`. Esse retorno preenche a flag `isError: true` no protocolo MCP, permitindo que o modelo de linguagem compreenda a razão do erro e tome medidas corretivas.
3. **Respostas Multimodais:** Para diagramas e inspeções visuais, o SDK suporta a emissão de blocos de imagem em Base64 via `mcplib.NewToolResultImage(legenda, base64Data, "image/png")`.
4. **Catálogo Unificado de 17 Ferramentas MCP:**
   O servidor expõe formalmente **17 ferramentas MCP**, divididas em 12 ferramentas pré-existentes e 5 novas ferramentas arquitetadas para resolver as dores do TCC:
   - **12 Ferramentas Pré-existentes:** `list_available_documents`, `get_doc_skeleton`, `read_doc_content`, `search_in_doc`, `list_doc_comments`, `reply_to_comment`, `resolve_comment`, `apply_text_style`, `search_and_replace_text`, `update_block_by_index`, `sync_pdf_to_drive`, `update_local_latex`.
   - **5 Novas Ferramentas ABNT/Integridade:** `apply_abnt_page_margins` (FEAT-01), `replace_block_abnt` (FEAT-02), `insert_native_table` (FEAT-06), `insert_mermaid_diagram` (FEAT-09), `check_pdf_sync` (FEAT-10).

---

## 2. Requisito 2 (R2): Arquitetura de Formatação, Estilos ABNT e Anti-Bleeding

### 2.1 Causa Raiz e Eliminação Definitiva do Bold Bleeding

#### Diagnóstico Técnico:
A ferramenta existente `search_and_replace_text` executa uma única operação `ReplaceAllTextRequest`. De acordo com a especificação formal do Google Docs, quando uma substring localizada é substituída, o novo texto inserido **herda passivamente a totalidade das propriedades de formatação de caractere (`TextStyle`) associadas ao primeiro caractere da substring casada**.

No caso do TCC, os títulos de seções originais estavam formatados com estilo de cabeçalho (`HEADING_1`) e `bold: true`. Quando um comando de substituição ancorou-se no título (ex: `"1.1 Contextualização"`) para injetar o texto da seção inteira, **todos os novos parágrafos, listas e citações herdaram negrito em cascata**, tornando o documento ilegível (evidenciado nas linhas 49-68 do arquivo `resources/Avaliação de uma Arquitetura Híbrida... (2).md`).

#### Solução Cirúrgica:
Proibir o uso de `ReplaceAllTextRequest` para blocos multi-linha. O novo mecanismo deve isolar a substituição por meio de intervalo atômico com reset explícito de estilo de caractere (`UpdateTextStyleRequest` com `bold: false, italic: false` e máscara `fields: "bold,italic"`).

---

### 2.2 Causa Raiz e Eliminação Definitiva do List / Bullet Bleeding

#### Diagnóstico Técnico:
No Google Docs, uma lista com marcadores (*bulleted list*) não constitui uma entidade envolvente separada; trata-se de uma propriedade intrínseca (`paragraph.Bullet != nil`) de cada elemento `Paragraph` individual.
Quando um texto contendo quebras de linha (`\n`) é inserido sobre um parágrafo que possui propriedade `Bullet`:
1. A Google Docs API divide a string em múltiplos parágrafos a cada caractere `\n`.
2. Como o parágrafo original continha um marcador de lista, **todos os novos parágrafos gerados herdam o identificador de lista e o nível de aninhamento**.
3. Quando o texto inserido contém quebras de linha duplas (`\n\n`) para separar parágrafos conceituais, é criado um parágrafo contendo zero caracteres visíveis, porém com marcador de lista associado: surge o **bullet vazio isolado** (`* `), conforme registrado na linha 226 do documento exportado.

#### Solução Cirúrgica:
O escritor deve disparar explicitamente uma operação `DeleteParagraphBulletsRequest` sobre todo o intervalo afetado antes ou imediatamente após a inserção do texto.

---

### 2.3 Correção Cirúrgica do Defeito em `internal/gdocs/writer.go:292`

A análise do código-fonte atual revelou uma falha primária de implementação na função `UpdateBlockByIndex`:

```go
// TRECHO COM DEFEITO em internal/gdocs/writer.go:292-305
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

#### Mecanismo do Defeito:
Se o desenvolvedor ou o LLM invoca `UpdateBlockByIndex` especificando explicitamente `ParagraphStyle: "NORMAL_TEXT"` com a intenção deliberada de desarmar um título ou normalizar o parágrafo, a expressão lógica avalia para `false`. O payload `UpdateParagraphStyleRequest` é **completamente suprimido**, mantendo o parágrafo preso ao estilo antigo (`HEADING_1` ou outro).

#### Correção Definitiva:
A condição deve permitir explicitamente o valor `"NORMAL_TEXT"`:
```go
// CORREÇÃO CIRÚRGICA:
if bloco.ParagraphStyle != "" {
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

---

### 2.4 Sequência Canônica de Requisições do Lote Atômico e Blindagens Defensivas

Para qualquer substituição ou inserção de bloco textual, o componente `BatchUpdateBuilder` gera uma sequência canônica de até 5 requisições atômicas, incorporando blindagens contra efeitos colaterais:

```json
{
  "requests": [
    {
      "deleteContentRange": {
        "range": {
          "startIndex": 1500,
          "endIndex": 1545
        }
      }
    },
    {
      "insertText": {
        "location": {
          "index": 1500
        },
        "text": "Este é o primeiro parágrafo formal do texto normatizado.\nEste é o segundo parágrafo sem qualquer resíduo de negrito ou marcadores.\n"
      }
    },
    {
      "deleteParagraphBullets": {
        "range": {
          "startIndex": 1500,
          "endIndex": 1629
        }
      }
    },
    {
      "updateParagraphStyle": {
        "range": {
          "startIndex": 1500,
          "endIndex": 1629
        },
        "paragraphStyle": {
          "namedStyleType": "NORMAL_TEXT",
          "alignment": "JUSTIFIED",
          "lineSpacing": 150,
          "indentFirstLine": {
            "magnitude": 35.433,
            "unit": "PT"
          },
          "indentStart": {
            "magnitude": 0,
            "unit": "PT"
          },
          "spaceAbove": { "magnitude": 0, "unit": "PT" },
          "spaceBelow": { "magnitude": 0, "unit": "PT" }
        },
        "fields": "namedStyleType,alignment,lineSpacing,indentFirstLine,indentStart,spaceAbove,spaceBelow"
      }
    },
    {
      "updateTextStyle": {
        "range": {
          "startIndex": 1500,
          "endIndex": 1629
        },
        "textStyle": {
          "bold": false,
          "italic": false,
          "underline": false
        },
        "fields": "bold,italic,underline"
      }
    }
  ]
}
```

#### Blindagens Defensivas Obrigatórias no `BatchUpdateBuilder`:

1. **Anulação de Recuo Residual de Lista (`indentStart: 0 PT`):**
   Ao remover marcadores de lista com `DeleteParagraphBulletsRequest`, o motor interno do Google Docs converte automaticamente a indentação visual anterior do marcador em recuo à margem esquerda (`indentStart`). Para impedir que o parágrafo normalizado fique indentado indevidamente como lista, o Passo 4 define expressamente `"indentStart": { "magnitude": 0, "unit": "PT" }` e inclui obrigatoriamente `"indentStart"` na máscara `fields`.

2. **Guarda de Substituição com Texto Vazio (Deleção Pura):**
   Se o texto higienizado a ser inserido for vazio (`new_text == ""` ou `len(sanitizedText) == 0`), submeter requisições de `InsertTextRequest`, `DeleteParagraphBulletsRequest` ou atualizações de estilo em faixas degeneradas `[startIndex, startIndex)` resulta em erro HTTP `400 Bad Request`.
   Nesse cenário, o construtor emite **exclusivamente a requisição [1] `DeleteContentRangeRequest`**, dispensando todos os passos subsequentes.

3. **Normalização Canônica Unicode NFC (`norm.NFC`):**
   Diferentes editores e sistemas operacionais podem codificar caracteres acentuados de forma decomposta (NFD, ex: `'a'` seguido do diacrítico de til `'\u0303'`) ou pré-composta (NFC, ex: `'ã'`). Para garantir que o cálculo de `UTF16Length` e as buscas por regex coincidam perfeitamente com os offsets internos do Google Docs, a função `SanitizeText` deve aplicar obrigatoriamente a normalização NFC via `golang.org/x/text/unicode/norm`:
   ```go
   import "golang.org/x/text/unicode/norm"

   func SanitizeText(raw string) string {
       // 1. Normalização canônica Unicode NFC
       normalized := norm.NFC.String(raw)
       // 2. Colapso de quebras supérfluas e espaços em branco residuais
       // ...
   }
   ```

---

### 2.5 Parametrização Métrica das Normas ABNT (NBR 14724 e NBR 6023)

As dimensões acadêmicas brasileiras expressas no sistema métrico (cm) devem ser convertidas com exatidão para pontos tipográficos (DTP Points, onde $72\text{ pt} = 1\text{ polegada} = 2,54\text{ cm}$):

$$\text{Fator de Conversão:} \quad 1\text{ cm} = \frac{72}{2,54} \approx 28,3464567\text{ pt}$$

#### Tabela de Mapeamento Paramétrico ABNT:

| Elemento ABNT | Norma / Regra | Medida Métrica | Valor em Pontos (`PT`) | Campo da Docs API | Requisição |
|---|---|---|---|---|---|
| **Recuo de 1ª Linha** | NBR 14724 (corpo) | $1,25\text{ cm}$ | **`35.433`** | `paragraphStyle.indentFirstLine` | `UpdateParagraphStyleRequest` |
| **Recuo de Citação** | Citação $> 3$ linhas | $4,00\text{ cm}$ | **`113.386`** | `paragraphStyle.indentStart` | `UpdateParagraphStyleRequest` |
| **Espaçamento Entrelinhas**| Corpo do trabalho | $1,5\text{ linhas}$| **`150`** (relativo) | `paragraphStyle.lineSpacing` | `UpdateParagraphStyleRequest` |
| **Alinhamento do Corpo** | Justificado | N/A | **`JUSTIFIED`** | `paragraphStyle.alignment` | `UpdateParagraphStyleRequest` |
| **Alinhamento Títulos** | Legendas/Figuras | N/A | **`CENTER`** | `paragraphStyle.alignment` | `UpdateParagraphStyleRequest` |
| **Margem Superior** | NBR 14724 (página) | $3,00\text{ cm}$ | **`85.039`** | `documentStyle.marginTop` | `UpdateDocumentStyleRequest` |
| **Margem Esquerda** | NBR 14724 (página) | $3,00\text{ cm}$ | **`85.039`** | `documentStyle.marginLeft` | `UpdateDocumentStyleRequest` |
| **Margem Inferior** | NBR 14724 (página) | $2,00\text{ cm}$ | **`56.693`** | `documentStyle.marginBottom` | `UpdateDocumentStyleRequest` |
| **Margem Direita** | NBR 14724 (página) | $2,00\text{ cm}$ | **`56.693`** | `documentStyle.marginRight` | `UpdateDocumentStyleRequest` |
| **Largura da Página A4** | Padrão ISO 216 | $21,00\text{ cm}$ | **`595.276`** | `documentStyle.pageSize.width` | `UpdateDocumentStyleRequest` |
| **Altura da Página A4** | Padrão ISO 216 | $29,70\text{ cm}$ | **`841.890`** | `documentStyle.pageSize.height` | `UpdateDocumentStyleRequest` |

#### Configuração Global de Página (`apply_abnt_page_margins`):
```json
{
  "updateDocumentStyle": {
    "documentStyle": {
      "marginTop": { "magnitude": 85.039, "unit": "PT" },
      "marginLeft": { "magnitude": 85.039, "unit": "PT" },
      "marginBottom": { "magnitude": 56.693, "unit": "PT" },
      "marginRight": { "magnitude": 56.693, "unit": "PT" },
      "pageSize": {
        "width": { "magnitude": 595.276, "unit": "PT" },
        "height": { "magnitude": 841.890, "unit": "PT" }
      }
    },
    "fields": "marginTop,marginLeft,marginBottom,marginRight,pageSize"
  }
}
```

---

## 3. Requisito 3 (R3): Arquitetura para Tabelas Nativas e Diagramas Mermaid

### 3.1 O Padrão Two-Step Approach para Criação e Preenchimento de Tabelas

A Google Docs API v1 não permite criar uma tabela e preencher suas células dentro de uma mesma chamada `batchUpdate`. Quando `InsertTableRequest` é processado, o motor do Google insere elementos estruturais complexos (`Table`, `TableRow`, `TableCell`) contendo marcadores invisíveis de limites de célula e parágrafos padrão inicializados com `\n`. Tentar calcular estaticamente esses offsets a priori resulta invariavelmente no erro:
`400 Bad Request: The insertion index must be inside the bounds of an existing paragraph.`

Ademais, no contexto acadêmico ABNT, uma tabela é acompanhada por uma **Legenda (`Caption`)** no topo e uma **Fonte (`Source`)** no rodapé. A ordem temporal e espacial dessas inserções é crítica para não invalidar os índices UTF-16 das células descobertos na inspeção estrutural.

#### O Algoritmo Two-Step Definitivo com Sequenciamento de Legenda e Fonte:

```
[Passo 1: Inserção da Legenda e Criação da Estrutura Vazia]
  1. Se houver Legenda (Caption):
     - Inserir parágrafo formatado centralizado em TargetIndex contendo "Tabela X – Legenda\n"
     - tableStartIndex = TargetIndex + UTF16Length(Caption + "\n")
  2. Se não houver Legenda:
     - tableStartIndex = TargetIndex
  3. Adicionar requisição atômica no lote:
     - InsertTableRequest { rows: R, columns: C, location: { index: tableStartIndex } }
  4. Executar docsService.BatchUpdate(docID, batchPasso1)

[Passo 2: Descoberta de Coordenadas e Captura de Revisão]
  1. Executar doc, err := docsService.GetDocument(docID)
  2. Capturar o identificador de revisão do documento: docRevisionId = doc.RevisionId
  3. Localizar o elemento Table em body.Content imediatamente após tableStartIndex
  4. Para cada linha r (0..R-1) e cada célula c (0..C-1):
     CellIndex[r][c] = table.TableRows[r].TableCells[c].Content[0].StartIndex
  5. Identificar o índice final da tabela: tableEndIndex = table.EndIndex

[Passo 3: Preenchimento Atômico em Lote Reverso com Trava Concorrente]
  1. Montar BatchUpdateDocumentRequest configurando WriteControl:
     WriteControl: &docs.WriteControl{RequiredRevisionId: docRevisionId}
  2. Se houver Fonte (Source):
     - Adicionar como a PRIMEIRA operação do lote a inserção de parágrafo centralizado com estilo ABNT em tableEndIndex contendo "\nFonte: " + Source.
       Como tableEndIndex é estritamente maior que qualquer coordenada de célula (tableEndIndex > CellIndex[r][c]), sua inserção antecipada no lote muta apenas índices posteriores a ela, mantendo os CellIndex[r][c] das células perfeitamente válidos e eliminando qualquer risco de o texto da fonte cair dentro de uma célula após o crescimento da tabela.
  3. Iterar de r = R-1 até 0 (ordem decrescente de linhas):
       Iterar de c = C-1 até 0 (ordem decrescente de colunas):
         Adicionar InsertTextRequest { location: { index: CellIndex[r][c] }, text: MatrizDados[r][c] }
  4. Adicionar requisição de estilização tabular (UpdateTableCellStyleRequest com TableStartLocation e TableRange)
  5. Executar docsService.BatchUpdate(docID, batchPasso3)
```

> **Por que a Ordem Decrescente e o Sequenciamento são Mandatórios?**  
> 1. **Isolamento de Índices Celulares e Regra do Maior Offset:** Ao inserir a Legenda *antes* da tabela no Passo 1, os índices `CellIndex[r][c]` e `tableEndIndex` capturados no Passo 2 já incorporam o deslocamento da legenda. No Passo 3, a inserção da Fonte em `tableEndIndex` deve ser categoricamente a primeira operação adicionada ao lote por representar o maior offset absoluto. Caso as células fossem preenchidas antes, o corpo da tabela se expandiria em $\sum \text{UTF16Length}(\text{MatrizDados}[r][c])$, fazendo com que o `tableEndIndex` capturado na tabela vazia passasse a apontar para o *interior* das células já preenchidas (ou gerando erro HTTP 400). Com a Fonte inserida primeiro no lote, seu deslocamento afeta apenas o texto posterior à tabela. O preenchimento reverso subsequente das células (de $(R-1, C-1)$ até $(0, 0)$) afeta unicamente offsets superiores, mantendo intactos os índices das células precedentes.  
> 2. **Controle Estrito de Concorrência e Política de Retry:** A utilização de `RequiredRevisionId: docRevisionId` no Passo 3 assegura que, se outro colaborador ou processo automatizado tiver editado o documento durante o intervalo entre o Passo 2 e o Passo 3, o lote é abortado com erro HTTP 400 (`Invalid revision ID`). Nessa situação, o cliente implementa política automática de retry: reexecuta o Passo 2 (reler o documento via `GetDocument`, remapear `CellIndex` e `tableEndIndex` e capturar o novo `docRevisionId`) e em seguida reenvia o Passo 3.

---

### 3.2 Estilização ABNT/IBGE sem Bordas Verticais

As normas do Instituto Brasileiro de Geografia e Estatística (IBGE) e da ABNT para apresentação tabular acadêmica determinam que tabelas não devem possuir fechamento lateral vertical.

#### Sintaxe Exata da Google Docs API v1 (`UpdateTableCellStyleRequest`):
No SDK oficial do Google (`google.golang.org/api/docs/v1`), a estilização de células de tabela **não aceita um campo `Range: {StartIndex, EndIndex}` genérico**. A API exige estritamente os campos `TableStartLocation` e `TableRange`:

```go
req := &docs.Request{
    UpdateTableCellStyle: &docs.UpdateTableCellStyleRequest{
        TableStartLocation: &docs.Location{
            Index: tableStartIndex,
        },
        TableRange: &docs.TableRange{
            TableCellLocation: &docs.TableCellLocation{
                RowIndex:    0,
                ColumnIndex: 0,
            },
            RowSpan:    rowsCount,
            ColumnSpan: colsCount,
        },
        TableCellStyle: &docs.TableCellStyle{
            BorderLeft: &docs.TableCellBorder{
                DashStyle: "SOLID",
                Width:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
            },
            BorderRight: &docs.TableCellBorder{
                DashStyle: "SOLID",
                Width:     &docs.Dimension{Magnitude: 0, Unit: "PT"},
            },
            BorderTop: &docs.TableCellBorder{
                DashStyle: "SOLID",
                Width:     &docs.Dimension{Magnitude: 1.0, Unit: "PT"},
                Color:     &docs.OptionalColor{Color: &docs.Color{RgbColor: &docs.RgbColor{Red: 0, Green: 0, Blue: 0}}},
            },
            BorderBottom: &docs.TableCellBorder{
                DashStyle: "SOLID",
                Width:     &docs.Dimension{Magnitude: 1.0, Unit: "PT"},
                Color:     &docs.OptionalColor{Color: &docs.Color{RgbColor: &docs.RgbColor{Red: 0, Green: 0, Blue: 0}}},
            },
            PaddingTop:    &docs.Dimension{Magnitude: 4, Unit: "PT"},
            PaddingBottom: &docs.Dimension{Magnitude: 4, Unit: "PT"},
            PaddingLeft:   &docs.Dimension{Magnitude: 6, Unit: "PT"},
            PaddingRight:  &docs.Dimension{Magnitude: 6, Unit: "PT"},
        },
        Fields: "borderLeft,borderRight,borderTop,borderBottom,paddingTop,paddingBottom,paddingLeft,paddingRight",
    },
}
```

---

### 3.3 Mapeamento das 22 Tabelas LaTeX do TCC para a Estrutura Go

O inventário das 22 tabelas presentes no TCC (`resources/tcc/capitulos/*.tex`) foi mapeado para a struct canônica `LaTeXTableData`, e o resultado da operação de inserção é retornado pela struct formal `TableInsertionResult`:

```go
// LaTeXTableData encapsula os metadados e matriz de dados de uma tabela extraída do LaTeX
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

// TableInsertionResult representa o resultado auditável da criação e preenchimento da tabela no Docs
type TableInsertionResult struct {
	TableStartIndex int64  `json:"table_start_index"`
	TableEndIndex   int64  `json:"table_end_index"`
	RowsCount       int64  `json:"rows_count"`
	ColsCount       int64  `json:"cols_count"`
	Message         string `json:"message"`
}
```

#### Requisitos de Normalização do Parser de LaTeX (`FEAT-08`):
A inspeção dos arquivos `.tex` reais do TCC revelou idiossincrasias sintáticas que devem ser tratadas pelo parser `table_parser.go`:
1. **Tratamento de `\addlinespace`:** O comando `\addlinespace` do pacote `booktabs` (usado abundantemente em `04-desenvolvimento.tex:396-406` no `quadro:trajetoria_decisoes`) deve ser sumariamente ignorado e descartado durante a tokenização de linhas.
2. **Normalização de `\textcite{...}`:** Citações bibliográficas inseridas em células (como em `02-referencial.tex:122-127` na `tab:trabalhos_relacionados`) devem ser convertidas para texto limpo (ex: extraindo o autor/ano legível ou a chave citada, sem comandos TeX residuais).
3. **Desescape de Caracteres Especiais (`\%`):** O símbolo percentual escapado `\%` deve ser convertido para `%` literal.
4. **Fórmulas Matemáticas Inline:** Expressões matemáticas simples delimitadas por cifrões (`$2 \times 2$`, `$\ge 2$`, `$\le 50\%$`) devem ser convertidas para seus equivalentes Unicode (`2 × 2`, `≥ 2`, `≤ 50%`), garantindo renderização limpa no Google Docs.

#### Catálogo das 22 Tabelas do TCC:
1. `tab:escopo_recorte` (Cap 1, L. 256): 3 colunas $\times$ 4 linhas (1 cabeçalho + 3 dados).
2. `tab:trabalhos_relacionados` (Cap 2, L. 113): 5 colunas $\times$ 8 linhas (1 cabeçalho + 7 dados).
3. `quadro:escopo_metodologico` (Cap 3, L. 21): 3 colunas $\times$ 5 linhas (1 cabeçalho + 4 dados).
4. `tab:distribuicao-requisitos` (Cap 3, L. 214): 3 colunas $\times$ 5 linhas (1 cabeçalho + 4 dados com total).
5. `tab:moscow` (Cap 3, L. 270): 3 colunas $\times$ 6 linhas (1 cabeçalho + 5 dados com total).
6. `tab:fatorial_2x2` (Cap 3, L. 326): 4 colunas $\times$ 5 linhas (1 cabeçalho + 4 condições experimentais).
7. `tab:rubrica_operacional` (Cap 3, L. 551): 4 colunas $\times$ 5 linhas (1 cabeçalho + 4 critérios TCD).
8. `tab:rf_dados` (Cap 4, L. 177): 5 colunas $\times$ 7 linhas (RF-01 a RF-06).
9. `tab:rf_normas` (Cap 4, L. 203): 5 colunas $\times$ 4 linhas (RF-07 a RF-09).
10. `tab:rf_validacao` (Cap 4, L. 226): 5 colunas $\times$ 6 linhas (RF-10 a RF-14).
11. `tab:rf_ui` (Cap 4, L. 251): 5 colunas $\times$ 7 linhas (RF-15 a RF-20).
12. `tab:rf_dialogo` (Cap 4, L. 277): 5 colunas $\times$ 5 linhas (RF-21 a RF-24).
13. `tab:rnf_qualidade` (Cap 4, L. 302): 5 colunas $\times$ 9 linhas (RNF-01 a RNF-08).
14. `tab:resumo_situacao_mvp` (Cap 4, L. 339): 3 colunas $\times$ 5 linhas (distribuição por status).
15. `quadro:trajetoria_decisoes` (Cap 4, L. 386): 4 colunas $\times$ 8 linhas (dilemas e decisões consolidadas).
16. `tab:sintese_adrs` (Cap 4, L. 504): 4 colunas $\times$ 10 linhas (catálogo de ADR-01 a ADR-09).
17. `tab:componentes` (Cap 4, L. 652): 3 colunas $\times$ 6 linhas (módulos técnicos do servidor).
18. `tab:rn_trancamento` (Cap 4, L. 712): 5 colunas $\times$ 8 linhas (regras normativas RN-01 a RN-07).
19. `tab:rn_proficiencia` (Cap 4, L. 739): 5 colunas $\times$ 8 linhas (regras normativas RN-08 a RN-14).
20. `tab:rn_calendario` (Cap 4, L. 768): 5 colunas $\times$ 4 linhas (regras normativas temporais RN-15 a RN-17).
21. `tab:tradeoffs_rod_json` (Cap 4, L. 856): 3 colunas $\times$ 5 linhas (matriz analítica de trade-offs).
22. `tab:rastreabilidade_solucao` (Cap 4, L. 1101): 5 colunas $\times$ 11 linhas (matriz de rastreabilidade).

---

### 3.4 Pipeline de Renderização e Inserção de Diagramas Mermaid

#### Limitação Técnica Fundamental:
O Google Docs **não aceita código Mermaid, SVG dinâmico nem arquivos vetoriais PDF** no payload de `InsertInlineImageRequest`. A API exige estritamente imagens rasterizadas (PNG/JPEG) apontadas por uma URI acessível via HTTP/HTTPS público.

#### Os 4 Diagramas do TCC:
1. `arquitetura.mmd` $\rightarrow$ Figura 2 (`04-desenvolvimento.tex:89`) — Estrutura predominantemente vertical (subgraphs empilhados).
2. `casos_uso.mmd` $\rightarrow$ Figura 3 (`04-desenvolvimento.tex:943`) — Diagrama vertical de atores e fronteiras.
3. `procedimento_experimental.mmd` $\rightarrow$ Figura 1 (`03-metodologia.tex:624`) — Fluxograma metodológico.
4. `sequencia.mmd` $\rightarrow$ Figura 4 (`04-desenvolvimento.tex:968`) — Diagrama de sequência vertical extenso.

#### Cálculo Dinâmico de Aspect Ratio com Enquadramento Bidimensional (2D Bounding Box):
A fixação estática de dimensões (como $420 \times 280\text{ pt}$) ou o cálculo unidimensional puramente baseado na largura projetam alturas excessivas em diagramas com orientação vertical acentuada. Por exemplo, `sequencia.mmd` possui uma proporção de aspecto $H/W \approx 2,2$; se a largura fosse fixada em $420\text{ pt}$, a altura resultante seria de $\approx 924\text{ pt}$, superando a altura total da folha A4 ($841,890\text{ pt}$) e extrapolando severamente a área útil imprimível entre as margens ABNT ($700,158\text{ pt}$).

Para garantir contenção perfeita dentro dos limites da página A4 e preservar espaço vertical para o parágrafo da Legenda no topo e da Fonte no rodapé, o componente `mermaid_writer.go` deve implementar um algoritmo de **Enquadramento Bidimensional (2D Bounding Box)** com limites estritos de largura máxima ($\text{maxWidth} = 450,0\text{ pt}$) e altura máxima ($\text{maxHeight} = 620,0\text{ pt}$).

Ademais, no ecossistema Go:
1. É mandatório registrar o codec PNG através do import anônimo `_ "image/png"` para que `image.DecodeConfig` reconheça o formato sem emitir o erro de execução `image: unknown format`.
2. Deve-se incluir guarda preventiva contra divisão por zero e dimensões corrompidas (`cfg.Width <= 0 || cfg.Height <= 0`).

```go
import (
	"bytes"
	"fmt"
	"image"
	_ "image/png" // Registro obrigatório do codec PNG para image.DecodeConfig
	"math"
)

// Decodificação rápida de dimensões sem carregar os pixels em memória
cfg, _, err := image.DecodeConfig(bytes.NewReader(pngBytes))
if err != nil {
	return nil, fmt.Errorf("falha ao decodificar dimensões da imagem: %w", err)
}

// Guarda preventiva contra divisão por zero e dimensões inconsistentes
if cfg.Width <= 0 || cfg.Height <= 0 {
	return nil, fmt.Errorf("dimensões inválidas da imagem rasterizada: %dx%d", cfg.Width, cfg.Height)
}

// Limites máximos da caixa delimitadora 2D na página A4 (normas ABNT):
// Largura útil total da página: 595.276 - 85.039 - 56.693 = 453.544 pt -> Teto operacional: 450.0 pt
// Altura útil total da página: 841.890 - 85.039 - 56.693 = 700.158 pt -> Menos Legenda e Fonte (~80 pt): 620.0 pt
maxWidthPt := 450.0
maxHeightPt := 620.0
if customWidthPt > 0 && customWidthPt <= maxWidthPt {
	maxWidthPt = customWidthPt
}

wPt := float64(cfg.Width)
hPt := float64(cfg.Height)

// Escala de contenção proporcional bidimensional (2D Bounding Box):
// scale = min(1.0, min(maxWidth / W, maxHeight / H))
// Se a imagem calculada a 450 pt exceder 620 pt de altura (ex: sequencia.mmd e arquitetura.mmd),
// a escala é automaticamente ajustada para que a altura respeite o teto de 620 pt.
scale := math.Min(1.0, math.Min(maxWidthPt/wPt, maxHeightPt/hPt))

targetWidthPt := wPt * scale
targetHeightPt := hPt * scale
```

#### Pipeline em 5 Passos com Fallback para Deep Link:
```
[1. Conversão para PNG 300 DPI]
   npx @mermaid-js/mermaid-cli -i figuras/arquitetura.mmd -o figuras/arquitetura.png -s 3 -b transparent
   (Alternativa: conversão dos PDFs vetoriais já existentes via pdftoppm -png -r 300)

[2. Upload para o Google Drive]
   driveService.CreateFile({ name: "arquitetura.png", parents: [folderID] }, pngReader, "id,webContentLink,webViewLink")

[3. Concessão de Acesso Temporário ou Fallback]
   Tentar conceder permissão temporária:
   driveService.CreatePermission(fileID, { role: "reader", type: "anyone" })
   
   FALLBACK DEEP LINK CORPORATIVO:
   Se a permissão pública falhar (ex: bloqueio de DLP organizacional do Google Workspace):
   - Não abortar a execução com erro fatal.
   - Em vez de InsertInlineImageRequest, inserir parágrafo com legenda formatada e Deep Link:
     "Figura X – [Legenda]\n[Diagrama hospedado no Google Drive: Clique aqui para visualizar: " + webViewLink + "]\nFonte: [Fonte]"
   - Retornar status de sucesso parcial informando a emissão do Deep Link.

[4. Injeção Inline no Google Docs (Caminho Principal)]
   BatchUpdate {
     insertInlineImage: {
       uri: "https://drive.google.com/uc?id=" + fileID + "&export=download",
       location: { index: targetIndex },
       objectSize: { 
         width: { magnitude: targetWidthPt, unit: "PT" }, 
         height: { magnitude: targetHeightPt, unit: "PT" } 
       }
     }
   }

[5. Persistência Interna no Google Docs]
   O servidor Google Docs faz o download do binário e incorpora no armazenamento interno.
   A imagem permanece intacta mesmo se a permissão do Drive for revogada posteriormente.
```

---

## 4. Requisito 4 (R4): Verificação Criptográfica e Sincronização de PDF

### 4.1 Comportamento do Metadado `md5Checksum` no Google Drive v3

| Tipo de Arquivo no Google Drive | Presença de `md5Checksum` | Justificativa Arquitetural |
|---|---|---|
| **Arquivos Binários** (`application/pdf`, imagens, zip) | **SEMPRE PRESENTE** | Armazenados como blobs imutáveis de bytes. O Drive computa o MD5 durante a ingestão do upload. |
| **Documentos Nativos** (`application/vnd.google-apps.*`) | **NUNCA PRESENTE** (`null`/vazio) | São representações dinâmicas de banco de dados sem arquivo estático correspondente. |

Como o artefato de sincronização do TCC é um arquivo **PDF binário (`main.pdf`)**, a presença e integridade do campo `md5Checksum` é **100% garantida pela API do Drive**.

---

### 4.2 Especificação do Arquivo Canônico do TCC (`main.pdf`)

- **Caminho Local:** `c:\Projetos\golang\mcp-gdocs-latex\resources\tcc\main.pdf`
- **Tamanho Exato em Disco:** **1.098.114 bytes**
- **ID do Arquivo Remoto no Drive:** `1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH`
- **URL no Google Drive:** `https://drive.google.com/file/d/1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH/view`

---

### 4.3 Especificação Algorítmica da Sincronização Idempotente

```go
package gdocs

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
	"google.golang.org/api/drive/v3"
)

type SyncReport struct {
	FileID       string `json:"file_id"`
	LocalMD5     string `json:"local_md5"`
	RemoteMD5    string `json:"remote_md5"`
	SizeBytes    int64  `json:"size_bytes"`
	Status       string `json:"status"` // "IN_SYNC", "UPDATED", "ERROR"
	Message      string `json:"message"`
}

func SynchronizePDF(ctx context.Context, driveSvc DriveService, localPath, fileID string) (*SyncReport, error) {
	// 1. Abrir arquivo local
	f, err := os.Open(localPath)
	if err != nil {
		return nil, fmt.Errorf("falha ao abrir arquivo PDF local: %w", err)
	}
	defer f.Close()

	// 2. Validação preventiva de integridade física do PDF (Checagem do marcador %%EOF)
	// Previne leitura e transmissão de PDFs incompletos caso o LaTeX/pdflatex ainda esteja compilando.
	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("falha ao inspecionar atributos do arquivo PDF local: %w", err)
	}
	readSize := int64(1024)
	if stat.Size() < readSize {
		readSize = stat.Size()
	}
	tailBuf := make([]byte, readSize)
	if _, err := f.ReadAt(tailBuf, stat.Size()-readSize); err != nil {
		return nil, fmt.Errorf("falha ao ler cauda do arquivo PDF: %w", err)
	}
	// Padronização estrita com a ISO 32000-1 (Seção 7.5.5 - End-of-file Marker: %%EOF)
	if !bytes.Contains(tailBuf, []byte("%%EOF")) {
		return nil, fmt.Errorf("arquivo PDF local corrompido ou incompleto: marcador %%%%EOF ausente nos últimos 1024 bytes (compilação LaTeX em andamento)")
	}

	// 3. Computar MD5 por streaming (sem carregar 1,1 MB todo em memória)
	hasher := md5.New()
	size, err := io.Copy(hasher, f)
	if err != nil {
		return nil, fmt.Errorf("falha ao calcular hash MD5: %w", err)
	}
	// Checagem de integridade concorrente contra Time-Of-Check to Time-Of-Use (TOCTOU):
	// Assegura que o compilador pdflatex não continuou escrevendo no arquivo durante a leitura de hashing.
	if size != stat.Size() {
		return nil, fmt.Errorf("inconsistência de leitura concorrente (TOCTOU): bytes lidos (%d) diferem do tamanho inicial inspecionado (%d), compilação LaTeX em andamento", size, stat.Size())
	}
	localMD5 := hex.EncodeToString(hasher.Sum(nil))

	// 4. Consultar metadados no Google Drive v3
	remoteFile, err := driveSvc.GetFileMetadata(ctx, fileID, "id, name, md5Checksum, size, modifiedTime")
	if err != nil {
		return nil, fmt.Errorf("falha ao consultar metadados no Drive: %w", err)
	}

	// 5. Verificação de Idempotência Criptográfica
	if strings.EqualFold(localMD5, remoteFile.Md5Checksum) {
		return &SyncReport{
			FileID:    fileID,
			LocalMD5:  localMD5,
			RemoteMD5: remoteFile.Md5Checksum,
			SizeBytes: size,
			Status:    "IN_SYNC",
			Message:   "O PDF local já é criptograficamente idêntico ao arquivo na nuvem. Upload suprimido.",
		}, nil
	}

	// 6. Se divergente, rebobinar arquivo e realizar upload atômico com política de retenção perpétua
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("falha ao rebobinar ponteiro do arquivo: %w", err)
	}

	// keepRevisionForever assegura histórico perene no Drive para auditoria e viabiliza rollback
	mediaMeta := &drive.File{
		KeepRevisionForever: true,
	}
	updatedFile, err := driveSvc.UpdateFileMedia(ctx, fileID, mediaMeta, f, "id, name, md5Checksum, size, modifiedTime, keepRevisionForever")
	if err != nil {
		return nil, fmt.Errorf("falha ao transmitir mídia para o Drive: %w", err)
	}

	// 7. Validação Pós-Upload
	if !strings.EqualFold(localMD5, updatedFile.Md5Checksum) {
		return nil, fmt.Errorf("alerta crítico de integridade: hash local (%s) difere do hash remoto pós-upload (%s)", localMD5, updatedFile.Md5Checksum)
	}

	return &SyncReport{
		FileID:    fileID,
		LocalMD5:  localMD5,
		RemoteMD5: updatedFile.Md5Checksum,
		SizeBytes: size,
		Status:    "UPDATED",
		Message:   "Upload concluído com sucesso e integridade MD5 verificada com retenção permanente de revisão.",
	}, nil
}
```

---

## 5. Requisito 5 (R5): Plano de Implementação, Mocks e Estratégia de Testes

### 5.1 Arquitetura Desacoplada: Adapter Pattern e Interfaces de Serviço

Para viabilizar 100% de testabilidade automatizada e isolar o código de dependências concretas externas, a base de código deve estruturar interfaces formais de domínio:

```go
package gdocs

import (
	"context"
	"io"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
)

type DocsService interface {
	GetDocument(ctx context.Context, documentID string) (*docs.Document, error)
	BatchUpdate(ctx context.Context, documentID string, req *docs.BatchUpdateDocumentRequest) (*docs.BatchUpdateDocumentResponse, error)
}

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
```

> **Papel Crítico do Método `ListFiles`:**  
> A inclusão de `ListFiles` em `DriveService` é mandatória para desacoplar a função de produção `ListAvailableDocuments` em `internal/gdocs/listing.go` (linhas 34-45). Sem essa assinatura, a ferramenta MCP de listagem permaneceria acoplada ao SDK concreto, impedindo a execução de testes automatizados offline no Milestone 4.

---

### 5.2 Framework de Testes Automatizados com `httptest.Server` e Roteamento de Uploads

O SDK oficial do Google em Go fornece a opção `option.WithEndpoint` e `option.WithoutAuthentication()`, permitindo inicializar instâncias completas de `*docs.Service` e `*drive.Service` apontando para um servidor HTTP de teste em memória.

#### Desafio Arquitetural de Uploads Multipart:
Nas chamadas de transmissão de mídia da API do Google Drive v3 (`Files.Update(...).Media(...)` e `Files.Create(...).Media(...)`), o SDK do Google redireciona automaticamente o transporte HTTP para a URL com prefixo `/upload/drive/v3/files/{fileId}` transmitindo o corpo com `Content-Type: multipart/related; boundary=...`. Para evitar falhas de rota HTTP 404 e erros na decodificação de JSON nos testes do Milestone 4, o mock server deve implementar um multiplexador (`http.ServeMux`) que processe conteúdo `multipart/related` e extraia o `fileID` ignorando o prefixo `/upload`:

```go
package gdocs_test

import (
	"context"
	"encoding/json"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// SetupMockWorkspaceServer inicializa um servidor mock multiplexado para Docs v1 e Drive v3
func SetupMockWorkspaceServer(t *testing.T, docsHandler, driveHandler http.HandlerFunc) (*docs.Service, *drive.Service, func()) {
	mux := http.NewServeMux()

	// Roteamento padrão da Google Docs API v1
	mux.HandleFunc("/v1/documents/", docsHandler)

	// Roteamento padrão da Google Drive API v3 (Metadados, Listagem e Permissões)
	mux.HandleFunc("/drive/v3/files", driveHandler)
	mux.HandleFunc("/drive/v3/files/", driveHandler)

	// Roteamento OBRIGATÓRIO de Upload Multipart da Google Drive API v3
	// O SDK google.golang.org/api direciona uploads de mídia para rotas com prefixo /upload/.
	// O mock handler processa conteúdo multipart/related e normaliza a URL removendo o prefixo /upload,
	// viabilizando a extração transparente do fileID pelo handler downstream.
	uploadDriveHandler := func(w http.ResponseWriter, r *http.Request) {
		// Normaliza a rota para que os handlers downstream extraiam o fileID ignorando o prefixo /upload
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/upload")

		// Trata requisições multipart/related do SDK oficial
		contentType := r.Header.Get("Content-Type")
		if mediaType, _, err := mime.ParseMediaType(contentType); err == nil && mediaType == "multipart/related" {
			// Processamento de envelope multipart/related do Google Drive Media Upload
		}

		driveHandler(w, r)
	}
	mux.HandleFunc("/upload/drive/v3/files", uploadDriveHandler)
	mux.HandleFunc("/upload/drive/v3/files/", uploadDriveHandler)

	server := httptest.NewServer(mux)
	ctx := context.Background()

	docsSvc, err := docs.NewService(ctx,
		option.WithEndpoint(server.URL),
		option.WithoutAuthentication(),
	)
	if err != nil {
		t.Fatalf("falha ao criar mock docs service: %v", err)
	}

	driveSvc, err := drive.NewService(ctx,
		option.WithEndpoint(server.URL),
		option.WithoutAuthentication(),
	)
	if err != nil {
		t.Fatalf("falha ao criar mock drive service: %v", err)
	}

	teardown := func() {
		server.Close()
	}

	return docsSvc, driveSvc, teardown
}
```

---

### 5.3 Matriz de Cobertura de Testes Unitários e de Integração

A tabela a seguir define a suíte de testes automatizados a ser entregue no Milestone 4:

| Arquivo de Teste | Alvo Testado | Cenários Cobertos | Tipo |
|---|---|---|---|
| `utf16_utils_test.go` | `UTF16Length`, `SanitizeText` e Normalização NFC | Caracteres ASCII, acentuação em português (`ã`, `ç`), caracteres decompostos (NFD vs NFC), emojis (4 bytes $\rightarrow$ 2 units UTF-16), remoção de `\n\n\n` redundantes. | Unitário (Puro) |
| `validator_test.go` | `ExtractDocID` e `IsDocAllowed` | URLs completas com `/edit?usp=sharing`, IDs puros, URLs inválidas, allowlists restritivas. | Unitário (Puro) |
| `style_abnt_test.go` | Geradores de Payload ABNT | Conversão métrica exata de cm para PT (35,433 pt, 85,039 pt, 56,693 pt), lineSpacing 150, alignment `JUSTIFIED`, indentStart 0 pt. | Unitário (Puro) |
| `hash_validator_test.go` | `CalculateMD5` e Validação de Integridade | Cálculo MD5 por streaming, validação obrigatória do marcador `%%EOF` nos últimos 1024 bytes (ISO 32000-1), detecção de PDF incompleto, mitigação de concorrência TOCTOU (`copiedBytes == stat.Size()`), paridade de hashes. | Unitário (Puro) |
| `builder_test.go` | `BatchUpdateBuilder` | Adição ordenada de deleção, inserção, estilo, remoção de bullets, clamping defensivo (`EndIndex = min(EndIndex, body.EndIndex - 1)`), guarda de substituição vazia, verificação da máscara `fields`. | Unitário (Puro) |
| `table_parser_test.go` | `ParseLaTeXTable` | Parsing das 22 tabelas do TCC, extração de headers, contagem de linhas e colunas, limpeza de macros `booktabs`, descarte de `\addlinespace`, normalização de `\textcite`, unescape de `\%` e fórmulas inline (`$2 \times 2$`). | Unitário (Puro) |
| `writer_test.go` | `ReplaceBlockABNT` e Fix 292 | Mock HTTP simulando `BatchUpdateDocumentResponse`. Teste da cura de Bold Bleeding e List Bleeding. Validação do envio de `NORMAL_TEXT` e anulação de recuo com `indentStart: 0`. | Integração Mock |
| `table_writer_test.go` | Two-Step Table Approach | Mock HTTP simulando criação de tabela, sequenciamento de Legenda no Passo 1, consulta `documents.get` no Passo 2 e lote reverso de preenchimento com inserção antecipada de Fonte em `tableEndIndex` e `WriteControl.RequiredRevisionId` no Passo 3. Sintaxe de `UpdateTableCellStyleRequest`. | Integração Mock |
| `mermaid_writer_test.go` | Inserção de Mermaid e Aspect Ratio | Mock HTTP simulando upload de imagem, concessão de permissão temporária, cálculo dinâmico de enquadramento 2D Bounding Box ($450 \times 620\text{ pt}$) via `image.DecodeConfig` com import `_ "image/png"`, e fallback de Deep Link em caso de bloqueio de permissão. | Integração Mock |
| `upload_test.go` | Sincronização de PDF | Mock HTTP multiplexado tratando `/drive/v3/files/` e a rota multipart `/upload/drive/v3/files/` com extração normalizada de `fileID`. Cenário `IN_SYNC` (skip upload), cenário `OUT_OF_SYNC` (executa update e confere hash) e política `keepRevisionForever: true`. | Integração Mock |

---

### 5.4 Procedimento de Validação E2E com Artefatos Reais do TCC

A validação prática (Milestone 5) deve ser executada pelo módulo `cmd/validator/main.go` cobrindo o acervo canônico do TCC:

1. **Validação das 22 Tabelas:**
   - Varre os arquivos `resources/tcc/capitulos/*.tex`.
   - Executa o parser estrutural sobre todas as 22 tabelas mapeadas.
   - Confere que $100\%$ das tabelas geram estruturas válidas com headers não vazios e matrizes retangulares perfeitas de células.
2. **Validação dos 4 Diagramas Mermaid:**
   - Inspeciona os arquivos `resources/tcc/figuras/*.mmd`.
   - Confere a existência dos vetores compilados `resources/tcc/figuras/*.pdf`.
   - Gera as imagens PNG rasterizadas em alta resolução (300 DPI) e valida dimensões físicas e proporções de aspecto compatíveis com o Google Docs ($\le 450\text{ pt}$).
3. **Validação Criptográfica do PDF:**
   - Abre `resources/tcc/main.pdf` (1.098.114 bytes).
   - Computa o hash MD5 canônico do binário.
   - Emite o relatório de integridade contra o identificador `1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH`.
