# Especificação Técnica e Análise Cirúrgica de APIs: Google Docs (v1), Google Drive (v3) e Servidor MCP em Go

**Data de Elaboração**: 2026-10-29  
**Autor**: Agente Minerador de Especificação (`spec_miner_api_1` / `teamwork_spec_miner`)  
**Contexto**: Servidor Go `mcp-gdocs-latex` — Ponte entre LaTeX acadêmico (TCC ABNT) e Google Docs / Google Drive.

---

## 1. Visão Geral e Fontes de Autoridade

Esta análise consolida as especificações definitivas das APIs do Google Workspace e as práticas do Model Context Protocol (MCP) em Go, com o objetivo de solucionar as falhas crônicas documentadas no histórico do projeto (`resources/outras infromações.md` e `resources/mcp_sync_history_and_improvements.md`):
1. **Efeito Cascata de Estilos ("Bold Bleeding")**: Propagação acidental de negrito e formatações de títulos para parágrafos inseridos via `ReplaceAllTextRequest`.
2. **Efeito Cascata de Listas ("Bullet Bleeding")**: Criação indevida de dezenas de marcadores (*bullets*) vazios ao injetar textos com quebras de linha (`\n`) sobre blocos que continham listas.
3. **Instabilidade de Índices UTF-16**: Deslocamento de offsets de caracteres (*Index Shifting*) e violação de limites de parágrafo.
4. **Ausência de Formatação ABNT Nativa**: Falta de padronização de recuo de primeira linha (1,25 cm), entrelinhas 1,5, alinhamento justificado e margens 3/3/2/2 cm.
5. **Dificuldade na Renderização de Tabelas e Diagramas Mermaid**: Destruição de arte ASCII por fontes proporcionais do Docs e ausência de tabelas estruturadas nativas.
6. **Falsa Sensação de Sincronia do PDF**: Ausência de verificação criptográfica entre o PDF compilado localmente e o arquivo armazenado no Google Drive.

### Fontes Oficiais Consultadas
- **Google Docs API v1 Reference**: `https://developers.google.com/docs/api/reference/rest/v1`
- **Google Drive API v3 Reference**: `https://developers.google.com/drive/api/v3/reference`
- **Model Context Protocol Go SDK (`mark3labs/mcp-go`)**: `https://github.com/mark3labs/mcp-go`
- **Código Fonte do Projeto**: `c:\Projetos\golang\mcp-gdocs-latex` (Go 1.25.5, `google.golang.org/api` v0.280.0)
- **Histórico e Requisitos Locais**: `resources/outras infromações.md` e `resources/tcc`

---

## 2. Tabelas Formais de Especificação

### Features Discovered

| # | Categoria | Feature | Descrição | Entradas (Inputs) | Saídas (Outputs) | Comportamento de Erro | Descoberto Via |
|---|-----------|---------|-----------|-------------------|------------------|----------------------|----------------|
| 1 | Docs API (v1) | `documents.batchUpdate` | Executa lote atômico e ordenado de requisições de mutação em um documento Docs | `documentId` (string), `requests` ([]*Request), `writeControl` (opcional) | `BatchUpdateDocumentResponse` contendo array de `Response` | Falha atômica completa: se uma requisição falhar, nenhuma alteração é persistida (HTTP 400/404/500) | Docs API v1 Spec |
| 2 | Docs API (v1) | `UpdateTextStyleRequest` | Atualiza formatação de caracteres (negrito, itálico, sublinhado, links, cor) em um intervalo UTF-16 | `range` (Range), `textStyle` (TextStyle), `fields` (FieldMask string) | Resposta vazia (`{}`) | HTTP 400 se `fields` estiver vazio, se `startIndex > endIndex`, ou se o range cruzar limites inválidos | Docs API v1 Spec |
| 3 | Docs API (v1) | `UpdateParagraphStyleRequest` | Modifica layout e tipografia de parágrafo (alinhamento, entrelinhas, recuo, `namedStyleType`) | `range` (Range), `paragraphStyle` (ParagraphStyle), `fields` (FieldMask string) | Resposta vazia (`{}`) | HTTP 400 se `range` for inválido ou se campos no mask não forem reconhecidos | Docs API v1 Spec |
| 4 | Docs API (v1) | `DeleteParagraphBulletsRequest` | Remove marcadores de lista (*bullets* ou números) de todos os parágrafos que cruzam a faixa especificada | `range` (Range) | Resposta vazia (`{}`) | HTTP 400 se o range estiver fora dos limites do documento | Docs API v1 Spec |
| 5 | Docs API (v1) | `CreateParagraphBulletsRequest` | Converte parágrafos dentro de uma faixa em itens de lista com marcadores ou numeração | `range` (Range), `bulletPreset` (string enum) | Resposta vazia (`{}`) | HTTP 400 se o preset não existir ou faixa for inválida; remove tabs iniciais automaticamente | Docs API v1 Spec |
| 6 | Docs API (v1) | `UpdateDocumentStyleRequest` | Altera propriedades globais de documento, incluindo margens (superior, inferior, esquerda, direita) e tamanho da página | `documentStyle` (DocumentStyle), `fields` (FieldMask string) | Resposta vazia (`{}`) | HTTP 400 se `fields` for omitido; margens são ignoradas em modo `PAGELESS` | Docs API v1 Spec |
| 7 | Docs API (v1) | `InsertTableRequest` | Cria uma tabela com número definido de linhas e colunas em um índice de inserção | `rows` (int), `columns` (int), `location` (Location) ou `endOfSegment` | Resposta vazia (`{}`) | HTTP 400 se `rows < 1` ou `columns < 1` ou inserção ocorrer fora de um parágrafo existente | Docs API v1 Spec |
| 8 | Docs API (v1) | `UpdateTableCellStyleRequest` | Customiza bordas, cor de fundo, alinhamento vertical e padding de células de uma tabela | `tableStartLocation` (Location) ou `tableRange` (TableRange), `tableCellStyle` (TableCellStyle), `fields` (FieldMask) | Resposta vazia (`{}`) | HTTP 400 se `tableStartLocation` não apontar exatamente para o início da tabela | Docs API v1 Spec |
| 9 | Docs API (v1) | `UpdateTableColumnPropertiesRequest` | Configura larguras de colunas específicas em uma tabela | `tableStartLocation` (Location), `columnIndices` ([]int), `tableColumnProperties` (TableColumnProperties), `fields` | Resposta vazia (`{}`) | HTTP 400 se o índice da coluna exceder o número de colunas da tabela | Docs API v1 Spec |
| 10 | Docs API (v1) | `InsertInlineImageRequest` | Insere imagem rasterizada (PNG/JPEG) no fluxo de texto com tamanho definido | `uri` (string URL pública), `location` (Location), `objectSize` (Size) | Resposta contendo `objectId` da imagem inserida | HTTP 400 se a URI for inacessível pelos servidores do Google (requer permissão pública ou link direto) | Docs API v1 Spec |
| 11 | Docs API (v1) | `DeleteContentRangeRequest` | Remove um intervalo contíguo de caracteres e elementos estruturais do documento | `range` (Range com `startIndex` e `endIndex`) | Resposta vazia (`{}`) | HTTP 400 se tentar deletar o caractere `\n` final obrigatório do Body | Docs API v1 Spec |
| 12 | Docs API (v1) | `InsertTextRequest` | Insere texto bruto em uma posição específica do documento | `location` (Location com `index`), `text` (string) | Resposta vazia (`{}`) | HTTP 400 se o índice não estiver dentro dos limites de um parágrafo existente | Docs API v1 Spec |
| 13 | Drive API (v3) | `files.get` (Integridade) | Obtém metadados de arquivos binários e de Workspace no Drive | `fileId` (string), `fields` ("id, name, md5Checksum, size, modifiedTime, webViewLink, webContentLink") | `drive.File` com campos populados | HTTP 404 se o ID não existir; `md5Checksum` é omitido para arquivos nativos do Workspace | Drive API v3 Spec |
| 14 | Drive API (v3) | `files.update` (Mídia) | Atualiza os bytes binários de um arquivo existente (como o PDF de TCC compilado) | `fileId` (string), `media` (io.Reader), `fields` (string) | `drive.File` atualizado | HTTP 403 se o usuário/robô não tiver permissão de edição; recalcula `md5Checksum` | Drive API v3 Spec |
| 15 | Drive API (v3) | `files.create` (Mídia) | Cria novo arquivo binário (ex: imagem PNG gerada de Mermaid) | `file` (File metadata com `name`, `parents`), `media` (io.Reader), `fields` | `drive.File` com ID gerado | HTTP 400 se pasta pai for inválida ou quota estourada | Drive API v3 Spec |
| 16 | Drive API (v3) | `permissions.create` | Concede acesso de leitura para terceiros (torna link público temporariamente) | `fileId` (string), `permission` (`{role: "reader", type: "anyone"}`) | `drive.Permission` com permissão concedida | HTTP 403 se a política da organização Workspace proibir compartilhamento externo público | Drive API v3 Spec |
| 17 | MCP Go SDK | `NewToolResultText` | Retorna resultado bem-sucedido de ferramenta em texto simples/Markdown | `text` (string) | `*mcplib.CallToolResult` com `Content` do tipo texto | N/A | MCP Go SDK Spec |
| 18 | MCP Go SDK | `NewToolResultError` | Retorna falha de execução de ferramenta com `isError: true` no protocolo MCP | `text` (string de erro explicativo) | `*mcplib.CallToolResult` com flag de erro | Permite que o LLM raciocine sobre o erro e se autocorrige em tempo de execução | MCP Go SDK Spec |
| 19 | MCP Go SDK | `NewToolResultImage` | Retorna imagem codificada em Base64 para o LLM visualizar na conversa | `text` (legenda), `imageData` (base64 string), `mimeType` (ex: "image/png") | `*mcplib.CallToolResult` com `ImageContent` | Erro se base64 for malformado | MCP Go SDK Spec |

---

### Edge Cases

| # | Feature | Entrada (Input) | Comportamento Observado | Causa Raiz e Solução Cirúrgica |
|---|---------|-----------------|-------------------------|--------------------------------|
| 1 | `ReplaceAllTextRequest` | Substituir string localizada dentro de um título em negrito por parágrafos comuns | **Bold Bleeding**: Todo o texto novo herdou negrito e estilo `HEADING_1` em cascata | A Docs API replica os atributos de caractere do primeiro ponto casado. **Solução**: Abandonar `ReplaceAllTextRequest` para blocos multi-linha; utilizar `DeleteContentRangeRequest` + `InsertTextRequest` + `UpdateParagraphStyleRequest(NORMAL_TEXT)` + `UpdateTextStyleRequest(bold:false)` no mesmo `batchUpdate`. |
| 2 | Inserção de texto em linha com bullet | `InsertTextRequest` contendo quebras de linha (`\n`) sobre posição que tinha marcador de lista | **List Cascading**: Cada `\n` inserido foi convertido pelo Docs em um novo bullet point (dezenas de bullets vazios) | O parágrafo hospedeiro continha estrutura `bullet != nil`. **Solução**: Chamar explicitamente `DeleteParagraphBulletsRequest` no range afetado e redefinir o estilo do parágrafo para `NORMAL_TEXT`. |
| 3 | `DeleteContentRangeRequest` | Range que atinge o final absoluto do documento (`startIndex < endIndex == body.EndIndex`) | Erro da API: `400 Bad Request: The content range must not include the final newline` | O Google Docs exige que todo documento possua um caractere final inamovível `\n`. **Solução**: Garantir que `endIndex <= body.EndIndex - 1`. |
| 4 | Contagem de caracteres vs UTF-16 | Texto com caracteres especiais, acentos (`á`, `ç`) ou emojis | Índices calculados em Go por `len(s)` (UTF-8 bytes) diferem dos índices do Docs (UTF-16 code units) | Acentos ocupam 2 bytes em UTF-8 mas 1 unidade em UTF-16; emojis ocupam 4 bytes em UTF-8 e 2 unidades em UTF-16. **Solução**: Usar estritamente `int64(len(utf16.Encode([]rune(text))))`. |
| 5 | Múltiplas substituições no mesmo lote | Múltiplos chunks aplicados em ordem crescente de índices (`0 -> 100`, depois `200 -> 300`) | Corrupção total: a primeira mutação alterou o comprimento do documento e invalidou os índices subsequentes | Deslocamento de índice (*Index Shifting*). **Solução**: Ordenar todas as deleções e inserções em ordem estritamente decrescente de `startIndex` (Reverse Order). |
| 6 | `InsertTableRequest` no mesmo batch de preenchimento | Inserir tabela e tentar injetar texto nas células no mesmo payload sem consultar a árvore | Rejeição da API: `The insertion index must be inside the bounds of an existing paragraph` | As células de uma tabela recém-criada contêm novos marcadores estruturais invisíveis e parágrafos próprios. **Solução**: Usar o padrão em 2 etapas: (1) `InsertTable`, (2) `documents.get` para capturar os `cell.Content[0].StartIndex`, (3) `InsertTextRequest` em lote reverso. |
| 7 | `InsertInlineImageRequest` com link do Google Drive | URI apontando para `https://drive.google.com/file/d/{ID}/view` ou URI sem acesso público | A imagem não carrega ou a API retorna erro `400 Bad Request: Invalid image uri` | O Google Docs requer download direto sem cookies de autenticação. **Solução**: Compartilhar o arquivo como `anyone/reader` e usar a URI `webContentLink` (`https://drive.google.com/uc?id={ID}&export=download`). |
| 8 | Margens em documento `PAGELESS` | `UpdateDocumentStyleRequest` com `marginTop`, etc. em documento configurado sem páginas | As propriedades de margem são ignoradas silenciosamente ou rejeitadas | Margens físicas só existem no modo `PAGED`. **Solução**: Verificar `doc.DocumentStyle.DocumentMode` e alternar para `PAGED` se necessário. |
| 9 | Leitura de hash em Google Doc nativo | `files.get` com `fields="md5Checksum"` em arquivo do tipo `application/vnd.google-apps.document` | `md5Checksum` retorna vazio (`""` ou `null`) | O Google Drive só calcula MD5 para arquivos binários (PDF, imagens, zip). **Solução**: Para verificação de integridade, comparar o MD5 do `main.pdf` local com o `md5Checksum` do arquivo PDF correspondente na nuvem. |

---

## 3. Google Docs API (v1) — Análise Cirúrgica de Endpoints e Mecanismos

### 3.1 Anatomia do `documents.batchUpdate` e Atomicidade

O endpoint `POST https://docs.googleapis.com/v1/documents/{documentId}:batchUpdate` é o único canal de modificação estrutural no Google Docs.

#### Propriedades Críticas:
1. **Atomicidade Estrita**: Todas as requisições em `requests: [...]` são executadas sequencialmente na ordem do array. Se qualquer operação falhar (ex: índice inválido, sobreposição proibida, colisão de estilo), a transação inteira sofre rollback imediato pelo backend do Google Docs. Nenhuma alteração intermediária é gravada no histórico.
2. **Semântica de Resposta**: O objeto `BatchUpdateDocumentResponse` retorna um array `replies`, onde cada elemento corresponde à resposta da requisição no mesmo índice da entrada. Operações de mutação de estilo (`UpdateTextStyle`, `UpdateParagraphStyle`, `DeleteParagraphBullets`) retornam respostas vazias (`{}`), enquanto operações de inserção de elementos (como imagens) retornam identificadores (`InsertInlineImage` retorna `objectId`).

---

### 3.2 Isolamento de Formatação e Eliminação do "Bold Bleeding"

#### O Fenômeno do Bold Bleeding
Ao usar `ReplaceAllTextRequest`, a API do Google Docs busca instâncias de texto e aplica uma regra interna: o texto substituto herda a totalidade dos atributos de formatação do **primeiro caractere** do texto substituído. Quando o usuário substitui títulos (ex: `1. INTRODUÇÃO` formatado com `bold: true` e `HEADING_1`) por textos longos, **todas as palavras e parágrafos do texto novo tornam-se negrito e herdam o cabeçalho**, inutilizando a leitura.

#### O Mecanismo Cirúrgico de Substituição com Reset Estrito
Para garantir isolamento absoluto, a mutação deve ser estruturada em 4 passos em lote atômico:

```json
{
  "requests": [
    {
      "deleteContentRange": {
        "range": {
          "startIndex": 100,
          "endIndex": 125
        }
      }
    },
    {
      "insertText": {
        "location": {
          "index": 100
        },
        "text": "Este é o primeiro parágrafo do texto normal.\nEste é o segundo parágrafo.\n"
      }
    },
    {
      "deleteParagraphBullets": {
        "range": {
          "startIndex": 100,
          "endIndex": 172
        }
      }
    },
    {
      "updateParagraphStyle": {
        "range": {
          "startIndex": 100,
          "endIndex": 172
        },
        "paragraphStyle": {
          "namedStyleType": "NORMAL_TEXT",
          "alignment": "JUSTIFIED",
          "lineSpacing": 150,
          "indentFirstLine": {
            "magnitude": 35.433,
            "unit": "PT"
          },
          "spaceAbove": {
            "magnitude": 0,
            "unit": "PT"
          },
          "spaceBelow": {
            "magnitude": 0,
            "unit": "PT"
          }
        },
        "fields": "namedStyleType,alignment,lineSpacing,indentFirstLine,spaceAbove,spaceBelow"
      }
    },
    {
      "updateTextStyle": {
        "range": {
          "startIndex": 100,
          "endIndex": 172
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

> **Regra Fundamental**: A especificação do `fields` mask é obrigatória. Se o `fields` for omitido ou vazio, o Google Docs não aplica a alteração.

---

### 3.3 Gestão e Higienização de Listas e Bullets

#### Como o Google Docs Trata Listas
No Google Docs, uma lista com marcadores (*bulleted list*) não é uma tag envolvente; ela é uma propriedade atribuída a cada elemento `Paragraph` individual (`paragraph.bullet != nil`).
Quando um parágrafo que possui `bullet` é editado ou sobrescrito com um texto que contém quebras de linha (`\n`), cada quebra de linha gera um **novo parágrafo subordinado ao mesmo identificador de lista**.
Isso gerou o problema histórico do TCC: centenas de marcadores vazios e bullets órfãos.

#### Procedimento de Desativação
Para remover bullets de qualquer trecho:
1. Executar `DeleteParagraphBulletsRequest` cobrindo o intervalo afetado:
```json
{
  "deleteParagraphBullets": {
    "range": {
      "startIndex": startOffset,
      "endIndex": endOffset
    }
  }
}
```
2. O Google Docs remove o marcador visual de lista e preserva o texto intacto.

---

### 3.4 Higienização de Quebras de Linha e Regras de Índices UTF-16

#### Regras de Índices do Google Docs
1. **Unidade de Medida**: Todas as posições (`startIndex`, `endIndex`, `location.index`) são medidas em **unidades de código UTF-16 (code units)**.
2. **Cálculo em Go**:
   - Em Go nativo, strings são arrays de bytes UTF-8 (`len(s)` retorna número de bytes).
   - O caractere `ã` tem 2 bytes em UTF-8, mas apenas **1 unidade em UTF-16**.
   - O emoji `🚀` tem 4 bytes em UTF-8, mas **2 unidades em UTF-16** (par substituto).
   - O cálculo correto para avançar índices no Go deve ser:
   ```go
   import "unicode/utf16"
   
   func UTF16Length(s string) int64 {
       return int64(len(utf16.Encode([]rune(s))))
   }
   ```
3. **Parágrafo Final Imutável**: O `endIndex` do último elemento do corpo do documento contém o `\n` final do documento, que **não pode** ser incluído em nenhuma `DeleteContentRangeRequest`. Caso contrário, a API responde `400 Bad Request`.

#### Algoritmo de Higienização de Texto Pré-Inserção
Antes de despachar o texto para o Google Docs, o servidor deve sanitizar os caracteres:
- Normalização de novas linhas: `\r\n` -> `\n`.
- Remoção de espaços em branco antes de quebras: `[ \t]+\n` -> `\n`.
- Remoção de quebras de linha redundantes: `\n{3,}` -> `\n\n` (no máximo uma linha em branco entre parágrafos conceituais).
- Remoção de múltiplos espaços contíguos: `[ ]{2,}` -> ` `.

---

### 3.5 Especificação de Formatação ABNT via API

A norma ABNT NBR 14724 / NBR 10520 estabelece padrões rigorosos de apresentação gráfica para trabalhos acadêmicos.

#### Conversão Métrica para Pontos Tipográficos
- $1\text{ polegada} = 72\text{ pontos (pt)} = 2,54\text{ cm}$
- $1\text{ cm} = \frac{72}{2,54}\text{ pt} \approx 28,3464567\text{ pt}$

#### Mapeamento Numérico ABNT na Google Docs API:

1. **Recuo de Primeira Linha do Parágrafo**:
   - Norma ABNT: $1,25\text{ cm}$
   - Em pontos: $1,25 \times 28,3464567 \approx \mathbf{35,433\text{ pt}}$
   - Payload `UpdateParagraphStyleRequest`:
   ```json
   "indentFirstLine": {
     "magnitude": 35.433,
     "unit": "PT"
   }
   ```

2. **Espaçamento Entre Linhas (Line Spacing)**:
   - Norma ABNT: $1,5$ linhas (no corpo do texto)
   - Na Docs API, a porcentagem é definida em relação a 100 (normal): $1,5 \times 100 = \mathbf{150}$
   - Payload:
   ```json
   "lineSpacing": 150
   ```

3. **Alinhamento do Texto**:
   - Norma ABNT: Justificado para parágrafos textuais; Centralizado para títulos de figuras e tabelas.
   - Payload:
   ```json
   "alignment": "JUSTIFIED"
   ```

4. **Margens da Página**:
   - Norma ABNT:
     * Superior: $3\text{ cm} \approx \mathbf{85,039\text{ pt}}$
     * Esquerda: $3\text{ cm} \approx \mathbf{85,039\text{ pt}}$
     * Inferior: $2\text{ cm} \approx \mathbf{56,693\text{ pt}}$
     * Direita: $2\text{ cm} \approx \mathbf{56,693\text{ pt}}$
   - Payload `UpdateDocumentStyleRequest`:
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

### 3.6 Tabelas Nativas no Google Docs

#### Estrutura de Dados de uma Tabela
No documento Google Docs, uma tabela é representada por um `StructuralElement` do tipo `Table`:
```
StructuralElement (startIndex, endIndex)
└── Table
    ├── TableRows[]
    │   └── TableCells[]
    │       └── Content[] (Array de StructuralElement: Paragraph, etc.)
```
Cada célula (`TableCell`) é criada automaticamente com um parágrafo inicial vazio contendo um caractere de quebra de linha `\n`.

#### O Padrão em 2 Etapas para Inserção e Preenchimento Seguro
Tentar calcular offsets estáticos de células antes da inserção da tabela no Docs frequentemente falha devido a marcadores de linha e célula que variam com a presença de cabeçalhos e estilos. O padrão oficial e mais resiliente é o **Two-Step Approach**:

1. **Etapa 1**: Enviar `InsertTableRequest`:
   ```json
   {
     "insertTable": {
       "rows": 4,
       "columns": 3,
       "location": { "index": targetInsertionIndex }
     }
   }
   ```
2. **Etapa 2**: Executar `svc.Documents.Get(documentID)` para obter o documento atualizado. Localizar a tabela recém-criada em `body.Content` a partir de `targetInsertionIndex`.
3. **Mapeamento de Células**:
   Para cada célula `cell := row.TableCells[colIndex]`, o índice de inserção exato do texto daquela célula é:
   $$\text{CellInsertionIndex} = \text{cell.Content}[0].\text{StartIndex}$$
4. **Etapa 3**: Gerar requisições `InsertTextRequest` para cada célula e enviar em um único `batchUpdate`.
   > **Atenção**: Como as inserções de texto alteram os índices das células subsequentes, o lote de preenchimento deve ser ordenado em **ordem decrescente de índice** (da última célula da última linha para a primeira célula da primeira linha). Dessa forma, nenhuma inserção desloca os índices das células ainda não preenchidas!

#### Formatação ABNT / Normas de Tabela (IBGE)
As normas acadêmicas brasileiras (IBGE / ABNT) determinam que tabelas:
- Devem ter bordas superior e inferior delimitando o cabeçalho.
- Devem ter uma borda inferior delimitando o final da tabela.
- **Não devem possuir bordas verticais nas extremidades laterais**.

Configuração via `UpdateTableCellStyleRequest`:
```json
{
  "updateTableCellStyle": {
    "tableStartLocation": {
      "index": tableStartIndex
    },
    "tableCellStyle": {
      "borderLeft": { "dashStyle": "NONE", "width": { "magnitude": 0, "unit": "PT" } },
      "borderRight": { "dashStyle": "NONE", "width": { "magnitude": 0, "unit": "PT" } },
      "paddingTop": { "magnitude": 4, "unit": "PT" },
      "paddingBottom": { "magnitude": 4, "unit": "PT" },
      "paddingLeft": { "magnitude": 6, "unit": "PT" },
      "paddingRight": { "magnitude": 6, "unit": "PT" }
    },
    "fields": "borderLeft,borderRight,paddingTop,paddingBottom,paddingLeft,paddingRight"
  }
}
```

---

### 3.7 Imagens Inline e Diagramas Mermaid

#### Requisitos do `InsertInlineImageRequest`
O payload oficial para inserção de imagem é:
```json
{
  "insertInlineImage": {
    "uri": "https://drive.google.com/uc?id={FILE_ID}&export=download",
    "location": {
      "index": insertionIndex
    },
    "objectSize": {
      "width": {
        "magnitude": 420.0,
        "unit": "PT"
      },
      "height": {
        "magnitude": 280.0,
        "unit": "PT"
      }
    }
  }
}
```

#### Restrições Críticas da URI e Ciclo de Vida da Mídia:
1. **Origem da URI**: O Google Docs API não aceita envio de bytes diretos no payload de inserção de imagem nem esquemas `data:image/png;base64,...`. A API exige uma URI externa acessível via HTTP/HTTPS público.
2. **Download Interno**: Quando a API recebe o `InsertInlineImageRequest`, os servidores do Google realizam um download imediato dos bytes da imagem a partir da URI fornecida e copiam o binário para o armazenamento de mídia interno do próprio documento Google Docs.
3. **Autonomia Após Inserção**: Uma vez processada a requisição com sucesso, o documento Google Docs mantém sua própria cópia da imagem. Se a URL original for desativada ou a permissão pública for revogada no Drive, a imagem **permanece perfeitamente visível e preservada no documento**.

#### Pipeline para Diagramas Mermaid (`resources/tcc/figuras/*.mmd` ou `*.pdf`):
1. **Geração do PNG**: Converter o arquivo Mermaid (`.mmd`) em imagem PNG usando `mmdc` (Mermaid CLI) ou extrair do `.pdf` compilado em alta resolução (300 DPI).
2. **Upload para o Google Drive**: Fazer upload do PNG para a pasta de recursos do trabalho via `drive.Files.Create`.
3. **Concessão de Acesso**: Adicionar permissão `anyone` como `reader` temporária via `drive.Permissions.Create`.
4. **Obtenção do Link**: Extrair `webContentLink` (`https://drive.google.com/uc?id={ID}&export=download`).
5. **Injeção no Docs**: Executar `InsertInlineImageRequest` no índice desejado com largura proporcional ajustada à margem útil da página (máximo de ~450 pt).

---

## 4. Google Drive API (v3) — Análise Cirúrgica de Metadados e Sincronização

### 4.1 Metadados de Integridade e Endpoint `files.get`

Para verificar a integridade de arquivos sem realizar o download desnecessário de grandes volumes de dados (como o PDF de dezenas de megabytes), a Google Drive API v3 oferece o endpoint:
```http
GET https://www.googleapis.com/drive/v3/files/{fileId}?fields=id,name,md5Checksum,size,modifiedTime,webViewLink,webContentLink
```

#### Propriedades dos Metadados:
- `id`: Identificador único do arquivo no Drive (string).
- `name`: Nome do arquivo com extensão (ex: `main.pdf`).
- `md5Checksum`: Hash MD5 hexadecimal do conteúdo binário do arquivo (calculado e armazenado pelo Google Drive no momento do upload).
- `size`: Tamanho exato em bytes no disco da nuvem (int64 em formato string).
- `modifiedTime`: Data/hora RFC 3339 da última modificação remota.
- `webContentLink`: Link de download direto utilizável pelo Google Docs para incorporar imagens.
- `webViewLink`: Link de visualização no navegador.

---

### 4.2 Especificação Estrita do Campo `md5Checksum`

| Tipo de Arquivo no Google Drive | Presença de `md5Checksum` | Motivo Arquitetural |
|---------------------------------|---------------------------|---------------------|
| **Arquivos Binários** (PDF, PNG, JPEG, ZIP) | **Sim (sempre populado)** | Armazenados como blobs brutos imutáveis por versão. O Google Drive calcula o hash MD5 durante o streaming de upload. |
| **Google Workspace Docs** (`application/vnd.google-apps.document`) | **Não (vazio ou omitido)** | São estruturas JSON relacionais de bancos de dados dinâmicos do Google, sem arquivo estático correspondente. |

> **Conclusão para o Requisito R4**: O arquivo de sincronização do TCC é o `main.pdf`. Por ser um arquivo binário, o campo `md5Checksum` **é 100% garantido e retornado pela API**.

---

### 4.3 Fluxo de Verificação de Integridade e Upload Atômico do PDF

#### Algoritmo em Go:
1. **Cálculo Local**: O servidor Go abre o arquivo local `resources/tcc/main.pdf` e computa o hash MD5:
   ```go
   f, err := os.Open(localPDFPath)
   if err != nil { return err }
   defer f.Close()

   hasher := md5.New()
   if _, err := io.Copy(hasher, f); err != nil { return err }
   localMD5 := hex.EncodeToString(hasher.Sum(nil))
   ```
2. **Consulta Remota**: O servidor consulta os metadados do PDF no Drive com `fields="id, name, md5Checksum, size, modifiedTime"`:
   ```go
   remoteFile, err := driveSvc.Files.Get(targetPDFFileID).Fields("id, name, md5Checksum, size, modifiedTime").Do()
   ```
3. **Decisão de Sincronia**:
   - Se `strings.EqualFold(localMD5, remoteFile.Md5Checksum)`:
     * O PDF local e o PDF na nuvem são criptograficamente idênticos.
     * Retorna status `IN_SYNC` com metadados (timestamp, tamanho, hash). Nenhuma operação de rede pesada é realizada.
   - Se os hashes divergirem:
     * Status `OUT_OF_SYNC`.
     * O servidor executa `driveSvc.Files.Update(targetPDFFileID, &drive.File{}).Media(f).Fields("id, name, md5Checksum, size, modifiedTime").Do()`.
     * Após o upload, confirma que o novo `md5Checksum` retornado pelo Drive bate com o `localMD5`.

---

## 5. Práticas e Arquitetura de Servidores MCP em Go

### 5.1 Especificação e Schema de Ferramentas no `mark3labs/mcp-go`

O servidor `mcp-gdocs-latex` utiliza o SDK Go oficial da comunidade (`github.com/mark3labs/mcp-go` v0.54.0).

#### Boas Práticas de Definição de Ferramentas:
1. **Descrições Ricas e Racionalizadas**: O modelo de linguagem utiliza a descrição da ferramenta para decidir *quando* e *como* chamá-la. Descrições devem conter restrições de uso, regras de formato e pré-condições.
2. **Validação de Schema Estruturado**:
   - Usar `mcplib.WithString(...)`, `mcplib.WithNumber(...)`, `mcplib.WithBoolean(...)` com `mcplib.Required()` para validação em tempo de decodificação.
   - Para payloads complexos, definir structs Go e realizar `json.Unmarshal`, validando antes de invocar a Google API.
3. **Tratamento de Erros Sem Quebrar a Sessão**:
   - **Nunca** retornar `(nil, err)` no handler do MCP a menos que ocorra uma falha de protocolo de transporte.
   - Em caso de erro de negócio ou falha na API do Google, retornar `mcplib.NewToolResultError(msg), nil`. Isso sinaliza `isError: true` na especificação do MCP, permitindo que o LLM leia a mensagem de erro no fluxo da conversa e tome ações corretivas imediatas.

---

### 5.2 Respostas com Conteúdo Multimodal (`ImageContent`)

O SDK `mark3labs/mcp-go` suporta nativamente respostas com imagens:
```go
func HandleRenderMermaid(ctx context.Context, req mcplib.CallToolRequest) (*mcplib.CallToolResult, error) {
    // 1. Gera ou lê o PNG da figura
    pngBytes, err := os.ReadFile("resources/tcc/figuras/arquitetura.png")
    if err != nil {
        return mcplib.NewToolResultError(fmt.Sprintf("Falha ao ler diagrama: %v", err)), nil
    }
    
    // 2. Codifica em base64
    base64Data := base64.StdEncoding.EncodeToString(pngBytes)
    
    // 3. Retorna imagem renderizada para o LLM
    return mcplib.NewToolResultImage("Diagrama de Arquitetura do TCC", base64Data, "image/png"), nil
}
```

---

### 5.3 Padrão Arquitetural Go `BatchUpdateBuilder`

Para substituir implementações dispersas e garantir atomicidade e isolamento de estilo, recomenda-se a criação do componente `BatchUpdateBuilder` em `internal/gdocs`:

```go
package gdocs

import (
	"context"
	"fmt"
	"sort"
	"google.golang.org/api/docs/v1"
)

// BatchUpdateBuilder gerencia e organiza operações atômicas no Google Docs
type BatchUpdateBuilder struct {
	documentID string
	requests   []*docs.Request
}

func NewBatchUpdateBuilder(documentID string) *BatchUpdateBuilder {
	return &BatchUpdateBuilder{documentID: documentID}
}

// AddDeleteRange adiciona deleção de faixa
func (b *BatchUpdateBuilder) AddDeleteRange(start, end int64) *BatchUpdateBuilder {
	b.requests = append(b.requests, &docs.Request{
		DeleteContentRange: &docs.DeleteContentRangeRequest{
			Range: &docs.Range{StartIndex: start, EndIndex: end},
		},
	})
	return b
}

// AddInsertText adiciona inserção de texto
func (b *BatchUpdateBuilder) AddInsertText(index int64, text string) *BatchUpdateBuilder {
	b.requests = append(b.requests, &docs.Request{
		InsertText: &docs.InsertTextRequest{
			Location: &docs.Location{Index: index},
			Text:     text,
		},
	})
	return b
}

// AddResetBullets remove marcadores de lista no intervalo
func (b *BatchUpdateBuilder) AddResetBullets(start, end int64) *BatchUpdateBuilder {
	b.requests = append(b.requests, &docs.Request{
		DeleteParagraphBullets: &docs.DeleteParagraphBulletsRequest{
			Range: &docs.Range{StartIndex: start, EndIndex: end},
		},
	})
	return b
}

// AddABNTParagraphStyle aplica estilo ABNT (recuo 1,25 cm, entrelinhas 1,5, justificado)
func (b *BatchUpdateBuilder) AddABNTParagraphStyle(start, end int64) *BatchUpdateBuilder {
	b.requests = append(b.requests, &docs.Request{
		UpdateParagraphStyle: &docs.UpdateParagraphStyleRequest{
			Range: &docs.Range{StartIndex: start, EndIndex: end},
			ParagraphStyle: &docs.ParagraphStyle{
				NamedStyleType: "NORMAL_TEXT",
				Alignment:      "JUSTIFIED",
				LineSpacing:    150,
				IndentFirstLine: &docs.Dimension{
					Magnitude: 35.433,
					Unit:      "PT",
				},
				SpaceAbove: &docs.Dimension{Magnitude: 0, Unit: "PT"},
				SpaceBelow: &docs.Dimension{Magnitude: 0, Unit: "PT"},
			},
			Fields: "namedStyleType,alignment,lineSpacing,indentFirstLine,spaceAbove,spaceBelow",
		},
	})
	return b
}

// Execute dispara todas as requisições em transação única atômica
func (b *BatchUpdateBuilder) Execute(ctx context.Context, svc *docs.Service) (*docs.BatchUpdateDocumentResponse, error) {
	if len(b.requests) == 0 {
		return nil, nil
	}
	req := &docs.BatchUpdateDocumentRequest{Requests: b.requests}
	return svc.Documents.BatchUpdate(b.documentID, req).Context(ctx).Do()
}
```

---

## 6. Mapeamento de Ferramentas MCP a Refatorar e Criar

Com base nas especificações mineradas, as seguintes ferramentas devem compor a arquitetura do servidor:

1. `check_pdf_sync` (Nova):
   - Compara hash MD5 do `resources/tcc/main.pdf` local com `md5Checksum` do Google Drive v3 (`files.get`).
   - Retorna relatório detalhado de sincronia (ID, data de compilação, hashes e status).

2. `sync_pdf_to_drive` (Refatorada):
   - Faz upload do PDF local para o Google Drive (`files.update`) e valida o hash MD5 após o upload.

3. `search_and_replace_clean` ou `replace_block_abnt` (Nova/Refatorada):
   - Substitui blocos textuais eliminando completamente o Bold Bleeding e List Bleeding.
   - Aplica `DeleteContentRangeRequest` + `InsertTextRequest` + `DeleteParagraphBulletsRequest` + `UpdateParagraphStyleRequest(NORMAL_TEXT, ABNT)` + `UpdateTextStyleRequest(bold:false)`.

4. `insert_native_table` (Nova):
   - Converte tabelas LaTeX estruturadas em tabelas nativas do Google Docs com bordas ABNT/IBGE via Two-Step Approach.

5. `insert_mermaid_diagram` (Nova):
   - Converte diagramas `.mmd` ou `.pdf` em PNG, realiza upload temporário com permissão pública no Drive e injeta via `InsertInlineImageRequest` no Docs.

6. `apply_abnt_page_margins` (Nova):
   - Aplica margens de página ABNT (3 cm superior/esquerda, 2 cm inferior/direita) via `UpdateDocumentStyleRequest`.
