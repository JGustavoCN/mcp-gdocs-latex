# Relatório de Handoff Técnico: Desafio Adversarial de API e Viabilidade em Go

**Agente:** `challenger_plan_2` (Engenheiro de Testes Adversariais e Integração de APIs)  
**Data:** 01 de Novembro de 2026  
**Destinatário:** Orquestrador (`parent` - ID `8522a030-5a02-4454-bc22-525b48f7de69`)  
**Missão:** Submeter o plano de engenharia ABNT (`docs/plano_engenharia_abnt.md`, `PROJECT.md`, `handoff.md` do `worker_plan_1`) a estresse de integração de API e viabilidade prática em Go.  
**Veredicto Oficial:** **REQUEST_CHANGES** (Aprovação Condicionada à Inclusão de Salvaguardas Defensivas de API)

---

## 1. Observation (Observações Diretas)

1. **Integridade da Suíte de Testes Existente:**
   - Comando executado no terminal Windows: `go test -count=1 ./...`
   - Saída obtida:
     ```
     ok  github.com/JGustavoCN/mcp-gdocs-latex/internal/gdocs   0.538s
     ok  github.com/JGustavoCN/mcp-gdocs-latex/internal/latex   1.232s
     ```
   - Código de retorno 0. O repositório encontra-se 100% verde, estável e livre de erros de compilação ou regressões.

2. **Auditoria da Especificação de Lotes e Índices no Plano:**
   - Em `docs/plano_engenharia_abnt.md:323-328` (Seção 1.8), o plano define a *Reverse Index Order Rule* para exclusões de parágrafos vazios com guarda pré-tabela (`content[i+1].Table != nil`) e clamping defensivo (`bodyEndIndex - 1`).
   - Em `docs/plano_engenharia_abnt.md:462-476` (Seção 3.1) e `docs/plano_engenharia_abnt.md:657-663` (Seção 4.3), o plano descreve os passos de `format_document_abnt`:
     - Passo 3.1: Margens de página (`UpdateDocumentStyleRequest`);
     - Passo 3.2: Particionamento de seção (`InsertSectionBreakRequest`);
     - Passo 3.3: Saneamento reverso de quebras de parágrafo (`DeleteContentRangeRequest`);
     - Passo 3.4: Formatação de parágrafos textuais e preservação de citações (`UpdateParagraphStyleRequest`);
     - Passo 3.5: Formatação de células de tabelas (`UpdateTableCellStyleRequest`).
   - O plano não especifica barreira de sincronização (recarga de snapshot com `GetDocument`) entre os passos estruturais (3.2 e 3.3) e os passos de estilização (3.4 e 3.5), o que permitiria que uma implementação ingênua enviasse todas as operações em um único lote atômico.

3. **Auditoria da Dependência de Template para Paginação:**
   - Em `docs/plano_engenharia_abnt.md:286-290` (Seção 1.6), é confirmada a impossibilidade de criar elementos `AutoText` (`PAGE_NUMBER`) dinâmicos via chamadas REST na Google Docs API v1.
   - O plano confia na clonagem de um *Template de Referência* no Google Drive (`resources/template/` ou ID canônico).
   - O plano não define o comportamento de `format_document_abnt` ou de `lint_document_abnt` quando o usuário executa as ferramentas em um documento avulso em branco que não foi originado do template e não contém cabeçalho prévio com `AutoText`.

4. **Auditoria das Interfaces de Serviço (Adapter Pattern):**
   - Em `PROJECT.md:203-207`, a interface proposta para o Docs é:
     ```go
     type DocsService interface {
         GetDocument(ctx context.Context, documentID string) (*docs.Document, error)
         BatchUpdate(ctx context.Context, documentID string, req *docs.BatchUpdateDocumentRequest) (*docs.BatchUpdateDocumentResponse, error)
     }
     ```
   - No SDK oficial da Google Docs API v1 (`google.golang.org/api/docs/v1`), as operações `UpdateDocumentStyleRequest`, `UpdateParagraphStyleRequest`, `InsertSectionBreakRequest` e `CreateHeaderRequest` são todas subcampos da struct polimórfica `docs.Request` e trafegam exclusivamente dentro de `docs.BatchUpdateDocumentRequest`.

---

## 2. Logic Chain (Cadeia de Raciocínio Lógico)

1. **Do Risco de Deriva Posicional (*Index Drift*) ao Erro HTTP 400:**
   - *Observação:* No Google Docs API, as requisições em um `batchUpdate` são aplicadas sequencialmente na ordem em que aparecem no array. Cada inserção ou exclusão altera imediatamente as coordenadas UTF-16 de todos os caracteres subsequentes.
   - *Dedução:* Se o Passo 3.2 insere uma quebra de seção (`\u000b`), todos os índices seguintes avançam +1. Se o Passo 3.3 deleta 10 quebras de linha espúrias, todos os índices seguintes retrocedem 10 posições.
   - *Consequência:* Se o Passo 3.4 aplicar `UpdateParagraphStyleRequest` utilizando os índices obtidos da leitura inicial antes das mutações do Passo 3.2 e 3.3, os parágrafos serão formatados fora de suas fronteiras reais. Ao atingir o final do documento encolhido, o `EndIndex` excederá o novo limite, disparando `HTTP 400 Bad Request: Invalid range`.
   - *Solução:* Exigir no plano a arquitetura de **Two-Phase Commit com Recarregamento de Snapshot** em `format_document_abnt`: Fase 1 (estrutural) -> Recarrega documento via `GetDocument` -> Fase 2 (estilização pura com coordenadas físicas estabilizadas).

2. **Da Fusão de Intervalos na *Reverse Index Order Rule*:**
   - *Observação:* A regra de ordenação decrescente de ranges em `DeleteContentRangeRequest` protege a integridade dos offsets anteriores quando aplicada em deleções puras.
   - *Dedução:* Contudo, se linhas em branco forem consecutivas (ex.: 3 `\n` seguidos), deletá-las como micro-requisições separadas gera tráfego desnecessário e amplia a superfície de falha nas fronteiras.
   - *Solução:* Exigir que intervalos adjacentes contíguos sejam colapsados em um único intervalo antes da submissão reversa.

3. **Do Beco sem Saída em Documentos Avulsos sem Template:**
   - *Observação:* A Google Docs API v1 não permite inserir `AutoText` (`PAGE_NUMBER`) dinâmico.
   - *Dedução:* Se um documento em branco for formatado, o `format_document_abnt` não conseguirá gerar o número dinâmico no cabeçalho. Consequentemente, o novo checker `checkSectionsAndNumbering` do `lint_document_abnt` continuará acusando anomalia de numeração indefinidamente.
   - *Solução:* O plano deve instruir o `format_document_abnt` a retornar uma advertência estruturada (*warning*) orientando a ação manual ou uso de template, e o linter deve tratar esse desvio como advertência consultiva e não como falha bloqueante em documentos sem template.

4. **Da Suficiência da Interface `DocsService`:**
   - *Observação:* Todas as quatro operações destacadas na missão (`UpdateDocumentStyle`, `UpdateParagraphStyle`, `InsertSectionBreak`, `CreateHeader`) são variantes de `docs.Request` processadas por `BatchUpdate`.
   - *Dedução:* A interface `DocsService` é 100% aderente ao padrão canônico do Google e cobre todas as necessidades do Milestone 6 sem exigir métodos REST espúrios.
   - *Melhoria:* Para agilizar testes de unidade determinísticos sem depender do servidor HTTP, recomenda-se disponibilizar uma struct `MockDocsService` (in-memory spy) que registre as requisições enviadas.

---

## 3. Caveats (Ressalvas e Limitações)

1. **Restrições Inerentes da API Google Workspace:**
   A impossibilidade de gerar `AutoText` e `TableOfContents` dinâmicos via REST v1 é uma limitação da infraestrutura da Google, confirmada oficialmente e contornada com sucesso através de templates pré-configurados e tabelas invisíveis ABNT.
2. **Escopo do Veredicto:**
   O veredicto **REQUEST_CHANGES** não invalida a arquitetura ou o trabalho do `worker_plan_1`. Trata-se de um aprimoramento preventivo estrito para blindar a fase de codificação do Milestone 6 contra bugs de concorrência e deslocamento de índices.

---

## 4. Conclusion (Conclusão e Veredicto)

**Veredicto Oficial: REQUEST_CHANGES**

O plano de engenharia em `docs/plano_engenharia_abnt.md` está aprovado em suas bases conceituais, métricas tipográficas e diagnósticos, mas **deve ser refinado antes do início da implementação** para incorporar formalmente os seguintes requisitos defensivos:

1. **Arquitetura Two-Phase Commit em `format_document_abnt`:**
   - Formalizar explicitamente que mutações estruturais (`InsertSectionBreakRequest` e `DeleteContentRangeRequest`) devem ser commitadas em um lote inicial isolado.
   - Inserir a barreira obrigatória de recarregamento do documento (`GetDocument`) para atualizar todos os offsets UTF-16 antes de emitir os lotes de estilização (`UpdateParagraphStyle`, `UpdateTextStyle`, `UpdateDocumentStyle`).
2. **Fusão de Intervalos Consecutivos no Saneamento Reverso:**
   - Incluir a etapa de colapso de intervalos em branco contíguos antes de ordenar decrescentemente.
3. **Contrato de Advertência para Documentos Sem Template Base:**
   - Especificar a mensagem de retorno estruturada de `warnings` para `format_document_abnt` e o nível consultivo para `lint_document_abnt` em documentos avulsos desprovidos de `AutoText`.
4. **Spy Mock Go para o Milestone 6:**
   - Definir a struct `MockDocsService` no pacote `internal/gdocs` para testes unitários com asserção estrita de payloads.

---

## 5. Verification Method (Método de Verificação Independente)

Para auditar e verificar as observações deste relatório:

1. **Executar a Suíte de Testes em Modo Sem Cache:**
   ```powershell
   go test -count=1 ./...
   ```
   *Resultado esperado:* Todos os pacotes (`internal/gdocs` e `internal/latex`) retornam `ok` com tempo de execução inferior a 2 segundos.

2. **Inspecionar o Relatório Completo de Análise Adversarial:**
   ```powershell
   Get-Content "c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\challenger_plan_2\analysis.md"
   ```

3. **Inspecionar as Linhas Auditadas de Deslocamento de Índice no Plano:**
   ```powershell
   Select-String -Path "c:\Projetos\golang\mcp-gdocs-latex\docs\plano_engenharia_abnt.md" -Pattern "Passo 3\.|Reverse Index Order"
   ```
