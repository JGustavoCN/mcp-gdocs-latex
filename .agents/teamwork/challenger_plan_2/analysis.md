# Relatório de Análise Adversarial e Testes de Estresse de API
## Projeto: `mcp-gdocs-latex` — Conformidade Estrita ABNT (Milestone 6)

**Challenger:** `challenger_plan_2` (Engenheiro de Testes Adversariais e Integração de APIs)  
**Data:** 01 de Novembro de 2026  
**Status da Auditoria:** Concluído com Apontamento de Vulnerabilidades Críticas  
**Veredicto:** **REQUEST_CHANGES** (Solicitação de Refinamento Defensivo no Plano de Engenharia)  

---

## 1. Sumário Executivo do Desafio

Esta auditoria adversarial submeteu o plano de engenharia (`docs/plano_engenharia_abnt.md`), a arquitetura do sistema (`PROJECT.md`) e os diagnósticos da equipe (`worker_plan_1`) a testes rigorosos de viabilidade computacional no ecossistema Go e de integração com a Google Docs API v1 e Google Drive API v3.

Foram avaliados com profundidade quatro pilares essenciais:
1. **Dinâmica de Lotes Atômicos (`batchUpdate`) e Regra de Ordem Reversa:** Risco de corrupção ou erro HTTP 400 por deriva de índices (*Index Drift*) durante deleções e inserções concorrentes.
2. **Separação de Fases Estrutural vs. Estilização:** Avaliação da execução do `format_document_abnt` e o perigo de mesclar mutações de topologia (quebras de seção e remoção de parágrafos) com estilização pontual de parágrafos.
3. **Robustez da Numeração de Página e Dependência de Template:** Viabilidade do `AutoText: PAGE_NUMBER` em documentos avulsos vs. documentos clonados de template base.
4. **Completude das Interfaces do Adapter Pattern e Mocks:** Avaliação das interfaces `DocsService` e `DriveService` e execução empírica da suíte existente (`go test ./...`).

---

## 2. Análise Detalhada dos Desafios Técnicos

### 2.1 Desafio 1: Deriva Posicional de Índices (*Index Drift*) em Mutações Mistas

#### Cenário de Ataque / Falha
No Google Docs API v1, um lote de requisições enviado via `POST /v1/documents/{documentId}:batchUpdate` processa os elementos na ordem estrita do array `requests`:
> *"Each request is applied in the order it appears in the list. The updates are atomic: either all updates are applied, or none are. Note that the index of any content in the document can shift as previous requests in the batch are applied."*

O plano de engenharia (`docs/plano_engenharia_abnt.md`, Seção 4.3) descreve a ordem do pipeline em `format_document_abnt`:
- Passo 3.1: Ajuste de margens de página (`UpdateDocumentStyleRequest`).
- Passo 3.2: Particionamento de seção antes da Introdução (`InsertSectionBreakRequest`).
- Passo 3.3: Saneamento reverso de quebras de parágrafo espúrias (`DeleteContentRangeRequest`).
- Passo 3.4: Formatação dos parágrafos textuais e preservação de citações longas (`UpdateParagraphStyleRequest`).
- Passo 3.5: Formatação de células de tabelas.

E na Seção 1.1 do `PROJECT.md`, é postulado:
> *"1. Atomicidade Estrita de Mutação: Todas as alterações estruturais em documentos Google Docs devem ser executadas em transações atômicas únicas via documents.batchUpdate."*

#### O que quebra se implementado em lote único:
1. O documento é lido na entrada da função via `GetDocument`. Suponha que o parágrafo `"1 INTRODUÇÃO"` esteja no offset 2962 e o parágrafo textual seguinte ocupe `[2994 - 3646]`.
2. O Passo 3.2 insere uma quebra de seção em 2962 (`InsertSectionBreakRequest`). O caractere `\u000b` (Section Break) é inserido, empurrando todo o conteúdo seguinte em **+1 code unit**. O parágrafo textual agora está fisicamente em `[2995 - 3647]`.
3. O Passo 3.3 deleta 14 linhas em branco espúrias na zona pré-textual (`DeleteContentRangeRequest`). Cada deleção remove 1 caractere (`\n`). O deslocamento acumulado agora é de $1 - 14 = -13$ code units! O parágrafo textual agora reside fisicamente em `[2981 - 3633]`.
4. O Passo 3.4 tenta executar `UpdateParagraphStyleRequest` utilizando o range original `[2994 - 3646]`.
5. **Resultado Catastrófico:**
   - O estilo de corpo (recuo de 1,25 cm, 1,5 entrelinhas) é aplicado no meio de outro parágrafo ou cruza a fronteira de um parágrafo vizinho!
   - Se o documento tiver encolhido e o último parágrafo for atualizado com seu `EndIndex` antigo, o índice excederá o novo fim do documento (`bodyEndIndex`), resultando em:
     `HTTP 400 Bad Request: Invalid range [start, end)`.

```
Leitura Inicial:        [---- P1 ----][\n][---- P2 ----][\n] ... [---- Pn ----][\n]
                             │                 │                       │
Mutação (Delete/Insert): [-3 chars]            │                       │
                                               ▼                       ▼
Índices Reais no Docs:  [-- P1 --][\n][---- P2 ----][\n] ... [---- Pn ----][\n]
Requisição no Batch:                  [---- P2 (velho) ----]       [---- Pn (velho) ----]
                                             ▲                              ▲
                                    Deslocado / Desalinhado        Cruza o Fim (HTTP 400!)
```

#### Solução Arquitetural Mandatória (Two-Phase Commit com Recarga de Snapshot)
O plano de engenharia deve estipular formalmente que o motor de formatação `format_document_abnt` **não deve misturar mutações de comprimento de texto com estilizações no mesmo lote**. Ele deve operar em duas fases sequenciais com barreira de sincronização:

```
┌────────────────────────────────────────────────────────────────────────┐
│ Fase 1: Mutação Estrutural e Topológica (Batch Atômico 1)              │
│ - InsertSectionBreakRequest (Particionamento pré-textual / textual)    │
│ - DeleteContentRangeRequest (Saneamento reverso de quebras de linha)   │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │ Commit Atômico na API
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│ Barreira de Sincronização: Recarregamento do Snapshot                 │
│ - doc, err = svc.GetDocument(ctx, documentID)                          │
│ - Reindexação completa com offsets UTF-16 físicos exatos               │
└───────────────────────────────────┬────────────────────────────────────┘
                                    │ Offsets Estabilizados
                                    ▼
┌────────────────────────────────────────────────────────────────────────┐
│ Fase 2: Estilização Pura e Geometria (Batch Atômico 2)                 │
│ - UpdateDocumentStyleRequest (Margens 3/3/2/2 cm e Folha A4)           │
│ - UpdateParagraphStyleRequest (Corpo 1,5, Títulos, Citações 4,0 cm)    │
│ - UpdateTextStyleRequest (Fontes 10pt / 12pt)                          │
│ - UpdateTableCellStyleRequest (Células de Tabelas)                     │
│ * Nota: Nenhuma operação nesta fase altera o comprimento do texto!     │
└────────────────────────────────────────────────────────────────────────┘
```

---

### 2.2 Desafio 2: Regra de Ordem Reversa (*Reverse Index Order Rule*) e Fusão de Intervalos

#### Avaliação do Algoritmo da Seção 1.8
O plano define:
1. Varrer `doc.Body.Content` e coletar intervalos vazios (`trimmedText == ""`).
2. Guarda Pré-Tabela: Não deletar parágrafo se `content[i+1].Table != nil`.
3. Execução Reversa: Ordenar intervalos do maior índice para o menor antes do lote.
4. Clamping Defensivo: Garantir que `endIndex <= bodyEndIndex - 1`.

#### Análise Crítica e Falhas de Borda Encontradas
1. **Ranges Adjacentes Contíguos:**
   - Se o documento tiver 3 linhas em branco contíguas (ex.: índices 100-101, 101-102, 102-103):
   - Se deletadas como 3 requests separados:
     - Request 1 deleta 102-103. O documento agora tem tamanho diminuído em 1.
     - Request 2 tenta deletar 101-102. Como 101-102 está antes de 102, ela continua válida.
     - Contudo, emitir centenas de mini-deleções unitárias satura o payload da API.
   - **Recomendação Defensiva:** O algoritmo deve colapsar ranges adjacentes consecutivos em um único intervalo contíguo `[start_min, end_max]` antes da submissão reversa (ex.: deletar diretamente `[100, 103]` em um único request).
2. **Vulnerabilidade em `cover_generator.go`:**
   - Em `cover_generator.go:891-912`, a limpeza tenta remover parágrafos que antecedem `"RESUMO"`. O código atual falha porque tenta deletar o parágrafo que abriga a `PageBreak` da Folha de Rosto e silencia o erro com `_, _ = docsSvc.Documents.BatchUpdate(...)`.
   - O plano exige proteger a `PageBreak`, mas deve enfatizar que a coleta de ranges para limpeza na capa também deve utilizar a mesma ordenação reversa estrita.
3. **Parágrafo Antecessor de Tabela:**
   - A guarda `doc.Body.Content[i+1].Table != nil` é vital. No Google Docs, uma tabela é um elemento de nível de bloco que exige obrigatoriamente um parágrafo precedente. Deletar esse parágrafo resulta em erro irrecuperável `400 Bad Request`. O plano previu isso adequadamente, mas essa regra deve ser validada por teste unitário específico.

---

### 2.3 Desafio 3: Robustez da Dependência de Template para `AutoText: PAGE_NUMBER`

#### A Restrição Intransponível da Google Docs API v1
A Google Docs API v1 REST **não possui nenhum endpoint ou request para criar elementos `AutoText`** (tokens dinâmicos como `PAGE_NUMBER` ou `PAGE_COUNT`). O elemento `AutoText` existe apenas como leitura na struct `docs.ParagraphElement`.

#### O Ponto de Ruptura (Cenário do Documento Avulso / em Branco)
O que acontece se o usuário submeter ao `format_document_abnt` ou ao `lint_document_abnt` um documento que:
- Não foi clonado a partir de um template institucional;
- Não possui cabeçalho configurado na interface Web do Docs;
- Está completamente em branco ou foi gerado por outra ferramenta externa?

1. **Tentativa de Criação via API:**
   - Se o formatador executar `CreateHeaderRequest` vinculado à Seção 2 (`sectionBreakLocation`), a API criará um cabeçalho contendo apenas um parágrafo vazio com `\n`.
   - Como não é possível injetar `AutoText`, o cabeçalho permanecerá sem numeração de página dinâmica.
   - Se o código tentar inserir texto estático (`"1"`, `"2"`), o número não incrementará nas páginas seguintes.
2. **Ciclo de Falso Erro no Linter:**
   - O `lint_document_abnt` (novo checker `checkSectionsAndNumbering`) identificará a ausência do número de página no cabeçalho e acusará anomalia.
   - O usuário tentará rodar `format_document_abnt` para corrigir a anomalia.
   - O formatador não conseguirá resolver a limitação da API REST.
   - O linter continuará reportando 100% de erro, deixando o usuário e o cliente MCP em um beco sem saída.

#### Mitigação e Requisitos Aditivos para o Plano:
1. **Contrato de Comportamento para Documentos Sem Template:**
   - O `format_document_abnt` deve verificar se o cabeçalho existente contém `AutoText`.
   - Se não contiver, o formatador cria a quebra de seção `SECTION_TYPE_NEXT_PAGE` e configura as margens, mas NÃO injeta strings falsas nem tenta operações inválidas.
   - O formatador deve emitir uma advertência descritiva no payload de retorno:
     ```json
     {
       "status": "FORMATTED_WITH_WARNINGS",
       "warnings": [
         "Numeração de página dinâmica ABNT requer a presença do campo de número de página no template original. Para exibir o número no canto superior direito, insira o campo uma única vez via interface Google Docs (Inserir > Números de página) ou utilize o template canônico IFS."
       ]
     }
     ```
2. **Classificação no Linter:**
   - No `lint_document_abnt`, a ausência de `AutoText` em documento não-template deve ser relatada como `WARNING` consultivo e não como defeito impeditivo do motor.

---

### 2.4 Desafio 4: Resiliência Concorrente, Adapter Pattern e Mocks em Go

#### Análise de Cobertura das Interfaces (`DocsService` e `DriveService`)
No `PROJECT.md` (linhas 203-219), a interface definida é:
```go
type DocsService interface {
    GetDocument(ctx context.Context, documentID string) (*docs.Document, error)
    BatchUpdate(ctx context.Context, documentID string, req *docs.BatchUpdateDocumentRequest) (*docs.BatchUpdateDocumentResponse, error)
}
```

Avaliamos se essa assinatura atende às 4 operações demandadas no Milestone 6:
- `UpdateDocumentStyle`: Enviado dentro de `req.Requests` como `*docs.UpdateDocumentStyleRequest`.
- `UpdateParagraphStyle`: Enviado dentro de `req.Requests` como `*docs.UpdateParagraphStyleRequest`.
- `InsertSectionBreak`: Enviado dentro de `req.Requests` como `*docs.InsertSectionBreakRequest`.
- `CreateHeader`: Enviado dentro de `req.Requests` como `*docs.CreateHeaderRequest`.

**Conclusão Técnica:** A interface `DocsService` é **100% compatível e suficiente**. Na API v1 do Google Docs, todas as mutações trafegam exclusivamente através de um único endpoint HTTP (`batchUpdate`). A granularidade da interface é canônica e não deve ser fragmentada em métodos específicos.

#### Recomendação para a Camada de Testes Mock em Go
Para que os testes do Milestone 6 garantam conformidade estrita sem depender de conexões de rede ou credenciais reais:
1. O `httptest.Server` previsto no Milestone 4 é essencial para testar tráfego HTTP.
2. Contudo, para testes unitários de regras de negócio (como verificar se `LineSpacing == 150`, se o recuo é `35.433 pt` ou se a ordem de deleção é decrescente), o plano deve prever um **Spy Mock** in-memory:
```go
type MockDocsService struct {
    DocToReturn      *docs.Document
    GetErr           error
    RecordedRequests []*docs.Request
    BatchUpdateErr   error
}

func (m *MockDocsService) GetDocument(ctx context.Context, docID string) (*docs.Document, error) {
    if m.GetErr != nil {
        return nil, m.GetErr
    }
    return m.DocToReturn, nil
}

func (m *MockDocsService) BatchUpdate(ctx context.Context, docID string, req *docs.BatchUpdateDocumentRequest) (*docs.BatchUpdateDocumentResponse, error) {
    if m.BatchUpdateErr != nil {
        return nil, m.BatchUpdateErr
    }
    m.RecordedRequests = append(m.RecordedRequests, req.Requests...)
    return &docs.BatchUpdateDocumentResponse{}, nil
}
```
Esse spy permite asserções determinísticas e instantâneas sobre os payloads gerados por `format_document_abnt` sem necessidade de serialização/deserialização JSON.

#### Verificação Empírica da Suíte Existente
Execução real realizada no terminal powershell:
```powershell
go test -count=1 ./...
```
- Pacote `internal/gdocs`: `ok  0.538s`
- Pacote `internal/latex`: `ok  1.232s`
- Pacote `internal/mcp`: sem arquivos de teste (neutro).
- **Resultado:** 100% verde, compilação limpa, sem testes quebrados ou instáveis.

---

## 3. Matriz de Riscos e Recomendações de Engenharia

| ID | Desafio / Vulnerabilidade | Severidade | Blast Radius | Mitigação Obrigatória no Plano |
|---|---|---|---|---|
| **VULN-01** | Deriva de índices (*Index Drift*) ao misturar mutação estrutural e estilização em lote único | **ALTA** | Corrupção de formatação de parágrafos e erro fatal HTTP 400 | Formalizar arquitetura **Two-Phase Commit**: Fase 1 (estrutural) -> Recarga de Snapshot (`GetDocument`) -> Fase 2 (estilização pura). |
| **VULN-02** | Deleção de intervalos adjacentes fragmentados sem fusão | **MÉDIA** | Ineficiência de payload e instabilidade em micro-deleções | Implementar algoritmo de colapso de intervalos vazios contíguos antes da ordenação reversa. |
| **VULN-03** | Ciclo de erro no linter e formatador em documentos avulsos sem `AutoText` | **MÉDIA** | Falha de conformidade impossível de sanar via REST API | Incluir retorno de `warnings` informativos no formatador e classificar a ausência como advertência no linter. |
| **VULN-04** | Falta de Spy Mock in-memory para asserção precisa de payloads de estilo | **BAIXA** | Testes de integração mais lentos e dependentes de HTTP | Adicionar struct `MockDocsService` para testes unitários dos novos checkers e builders. |

---

## 4. Veredicto Final

**VEREDICTO: REQUEST_CHANGES**

O plano de engenharia elaborado pelo `worker_plan_1` apresenta excelente profundidade teórica e diagnóstica. No entanto, sua aprovação final para codificação depende da incorporação formal das seguintes revisões defensivas:
1. **Inclusão da Separação em Duas Fases com Recarregamento de Snapshot** em `format_document_abnt` para neutralizar o risco de *Index Drift*.
2. **Definição Clara da Tolerância e Mensagem de Advertência para Documentos Sem Template** no que tange ao `AutoText: PAGE_NUMBER`.
3. **Especificação do Algoritmo de Fusão de Intervalos Vazios Consecutivos** antes da execução da *Reverse Index Order Rule*.
