# Relatório de Handoff — Diagnóstico da Base de Código Go (`mcp-gdocs-latex`)

**Tipo de Handoff:** Hard (Tarefa Concluída)  
**Data:** 29 de outubro de 2026  
**Agente:** `explorer_codebase_1` (`teamwork_explorer`)  
**Destinatário:** Orquestrador / Equipe de Arquitetura e Engenharia  
**Arquivo de Análise Completa:** `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\explorer_codebase_1\analysis.md`

---

## 1. Observation

Durante a investigação da base de código Go e dos documentos de requisitos e histórico do projeto, foram diretamente observados os seguintes fatos:

1. **Estrutura e Inventário de Código:**
   - A base de código contém exatamente 11 arquivos `.go`:
     - `cmd/mcp-server/main.go`
     - `internal/gdocs/auth.go`
     - `internal/gdocs/comments.go`
     - `internal/gdocs/listing.go`
     - `internal/gdocs/reader.go`
     - `internal/gdocs/search.go`
     - `internal/gdocs/upload.go`
     - `internal/gdocs/validator.go`
     - `internal/gdocs/writer.go`
     - `internal/latex/local_file.go`
     - `internal/mcp/tools.go`
   - Não há **nenhum arquivo de teste unitário** (`*_test.go`) em nenhum diretório. A cobertura de testes automatizados é de 0%.

2. **Dependências em `go.mod`:**
   - Go versão: `1.25.5`.
   - SDK MCP: `github.com/mark3labs/mcp-go v0.54.0`.
   - Google APIs: `google.golang.org/api v0.280.0` (incluindo `docs/v1` e `drive/v3`).
   - Autenticação: `golang.org/x/oauth2 v0.36.0`.

3. **Catálogo de Ferramentas Registradas (`internal/mcp/tools.go:20-31`):**
   - 12 ferramentas registradas: `list_available_documents`, `get_doc_skeleton`, `read_doc_content`, `search_in_doc`, `list_doc_comments`, `reply_to_comment`, `resolve_comment`, `apply_text_style`, `search_and_replace_text`, `update_block_by_index`, `sync_pdf_to_drive`, `update_local_latex`.
   - A função `MultiReplaceDocBlock` está implementada em `internal/gdocs/writer.go:103-182`, porém **não está registrada** no `internal/mcp/tools.go`.

4. **Bug Crítico na Aplicação de Estilo de Parágrafo (`internal/gdocs/writer.go:292`):**
   - O código contém a seguinte condição:
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
   - Se o chamador solicitar `ParagraphStyle: "NORMAL_TEXT"`, a instrução é descartada e o parágrafo retém o estilo de cabeçalho anterior.

5. **Herança Passiva de Estilo em Substituição Global (`internal/gdocs/writer.go:29-41`):**
   - `SearchAndReplaceText` executa exclusivamente um `ReplaceAllTextRequest`. A API do Google Docs herda a formatação do primeiro caractere substituído para todo o novo texto inserido.
   - No documento exportado real (`resources/Avaliação de uma Arquitetura Híbrida... (2).md`, linhas 49-68), blocos inteiros de texto como `"1. Realizar consultas..."` até `"7. Avaliar..."` ficaram travados em negrito (`**...**`).

6. **Efeito Cascata de Listas e Bullets Vazios:**
   - Em `internal/gdocs/reader.go:171-240`, `processParagraph` avalia apenas `NamedStyleType` e **ignora completamente** o campo `paragraph.Bullet`.
   - No escritor `internal/gdocs/writer.go`, não há nenhuma chamada a `DeleteParagraphBulletsRequest`.
   - Na linha 226 do documento exportado (`resources/Avaliação de uma Arquitetura Híbrida... (2).md`), observou-se uma linha com bullet vazio isolado: `* `.

7. **Upload de PDF sem Hash e sem Idempotência (`internal/gdocs/upload.go:11-30`):**
   - O arquivo possui 31 linhas. Executa apenas `svc.Files.Update(fileID, &drive.File{}).Media(file).Do()`.
   - Não calcula hash MD5 local do arquivo `resources/tcc/main.pdf`.
   - Não lê o campo `md5Checksum` da API do Drive v3 (`svc.Files.Get`).
   - Retorna a string fixa `"SUCESSO: PDF atualizado silenciosamente no ID %s"`.

8. **Inexistência de Suporte a Tabelas e Imagens:**
   - Não há ocorrências de `InsertTableRequest` nem de `InsertInlineImageRequest` na base Go.
   - No LaTeX do TCC existem dezenas de tabelas `tabularx` (ex: `01-introducao.tex:256`, `03-metodologia.tex:21`) e 4 diagramas Mermaid (`arquitetura.mmd`, `casos_uso.mmd`, `procedimento_experimental.mmd`, `sequencia.mmd`).
   - No GDocs exportado, todas as tabelas e figuras foram substituídas por links de texto do tipo: `[VER TABELA NO PDF COMPILADO — ... : https://drive.google.com/...]`.

---

## 2. Logic Chain

1. **Origem do Bold Bleeding:**
   - *Premissa:* `SearchAndReplaceText` utiliza `ReplaceAllTextRequest` sem enviar requisições adicionais de estilo.
   - *Comportamento da API:* A API Google Docs propaga o `TextStyle` da string antiga para o texto novo.
   - *Constatação:* Quando o Agente substituiu títulos em negrito por parágrafos completos, a formatação de negrito foi herdada passivamente.
   - *Tentativa de Correção Frustrada:* Quando se tentou usar `update_block_by_index` com `ParagraphStyle: "NORMAL_TEXT"`, a linha 292 de `writer.go` ignorou o comando devido à condição `bloco.ParagraphStyle != "NORMAL_TEXT"`. Consequentemente, o negrito e a formatação de parágrafo permaneceram corrompidos.

2. **Origem dos Bullets Vazios:**
   - *Premissa:* No Google Docs, uma lista é uma propriedade intrínseca do parágrafo (`paragraph.Bullet`).
   - *Constatação:* `reader.go` não exibe a existência de listas e `writer.go` não remove listas.
   - *Dedução:* Ao substituir um parágrafo que era item de lista por um texto contendo múltiplas quebras de linha (`\n`), cada `\n` gerou um novo item da lista. Onde havia linhas em branco (`\n\n`), surgiu um parágrafo sem texto porém com marcador de lista: o bullet vazio (`* `). Para desvincular o parágrafo da lista, a API do Google exige obrigatoriamente `DeleteParagraphBulletsRequest`.

3. **Origem das Quebras de Linha Excessivas:**
   - *Premissa:* O Google Docs armazena delimitadores de parágrafo no final de cada `Paragraph`.
   - *Constatação:* `reader.go` adiciona `\n` ao concatenar parágrafos após o `strings.TrimSpace`, enquanto strings injetadas com `\n\n` geram novos parágrafos estruturais vazios.

4. **Insegurança e Desperdício na Sincronização de PDF:**
   - *Premissa:* A integridade da sincronia requer comprovação de que o PDF na nuvem é idêntico ao compilado localmente.
   - *Constatação:* `upload.go` não consulta o `md5Checksum` da Google Drive API v3 e não calcula MD5 do arquivo local. Faz upload cego de mídia a cada execução, sem garantia de integridade nem idempotência.

5. **Inviabilidade Atual de Testes Automatizados:**
   - *Premissa:* Para testar de forma isolada, os serviços devem permitir substituição de chamadas externas por mocks.
   - *Constatação:* `auth.go` instancia structs concretas do Google SDK em variáveis globais com `sync.Once`.
   - *Dedução:* Mocks sem internet só são possíveis por meio de refatoração para interfaces de domínio (Adapter Pattern) ou inicialização do client apontando para `httptest.Server` via `option.WithEndpoint`.

---

## 3. Caveats

1. **Credenciais em Produção:** A investigação foi conduzida em modo estritamente somente-leitura e offline em relação às credenciais reais do Google Cloud. As conclusões foram validadas pelo código-fonte de `google.golang.org/api`, pelo código Go do projeto e pelos arquivos exportados em `resources/`.
2. **MultiReplaceDocBlock:** A função `MultiReplaceDocBlock` em `writer.go` implementa ordenação decrescente de blocos, mas não está registrada em `internal/mcp/tools.go`. Não foram encontrados registros históricos de testes reais dela em produção.

---

## 4. Conclusion

A base de código Go atual possui uma arquitetura limpa em termos de divisão de pacotes (`internal/gdocs`, `internal/latex`, `internal/mcp`), porém sofre de **seis lacunas críticas de engenharia**:
1. Defeito no tratamento de `NORMAL_TEXT` no `update_block_by_index` (linha 292 de `writer.go`).
2. Ausência de ferramenta atômica `search_and_replace_clean` para aplicar texto e redefinir estilos no mesmo payload.
3. Ausência de manipulação de listas via `DeleteParagraphBulletsRequest`.
4. Ausência de cálculo de hash MD5 e checagem de `md5Checksum` em `sync_pdf_to_drive`.
5. Ausência de ferramentas de inserção de tabelas nativas (`InsertTableRequest`) e imagens/diagramas Mermaid (`InsertInlineImageRequest`).
6. Ausência total de suíte de testes unitários (`*_test.go`).

A implementação de uma arquitetura robusta requer a correção dessas seis lacunas, estruturando interfaces para permitir 100% de cobertura com testes automatizados simulados via `httptest.Server`.

---

## 5. Verification Method

Para verificar independentemente todas as afirmações e achados deste relatório:

1. **Verificação do Bug do `NORMAL_TEXT`:**
   - Inspecionar `c:\Projetos\golang\mcp-gdocs-latex\internal\gdocs\writer.go` nas linhas 292 a 305. Constatar que a condição `bloco.ParagraphStyle != "NORMAL_TEXT"` exclui a criação do `UpdateParagraphStyleRequest`.

2. **Verificação da Ausência de Testes:**
   - Executar no terminal:
     ```powershell
     cd c:\Projetos\golang\mcp-gdocs-latex
     go test ./...
     ```
   - O comando retornará `? [no test files]` para todos os pacotes.

3. **Verificação da Ausência de Hash no PDF:**
   - Inspecionar `c:\Projetos\golang\mcp-gdocs-latex\internal\gdocs\upload.go` nas linhas 11 a 30. Confirmar que não há chamadas a `crypto/md5`, `md5Checksum` ou `svc.Files.Get`.

4. **Verificação da Evidência Real de Bold Bleeding e Bullet Vazio:**
   - Inspecionar `c:\Projetos\golang\mcp-gdocs-latex\resources\Avaliação de uma Arquitetura Híbrida RAG-MCP na Validação Normativa de Requerimentos Acadêmicos (2).md`:
     - Linhas 49 a 68: Parágrafos inteiros prefixados e sufixados com `**` (negrito em cascata).
     - Linha 226: Presença de bullet vazio isolado (`* `).
     - Linhas 87, 124, 182, 206, 236: Presença de marcadores textuais `[VER TABELA NO PDF COMPILADO — ...]`.
