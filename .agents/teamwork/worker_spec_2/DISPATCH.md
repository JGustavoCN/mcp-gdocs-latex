## 2026-10-29T19:37:01Z

Você é o worker_spec_2 (teamwork_worker) na Iteração 2 do refinamento arquitetural.
Seu diretório de trabalho é: c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_spec_2\
Sua identidade: teamwork_worker
LEIA OBRIGATORIAMENTE:
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\ORIGINAL_REQUEST.md
- c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md
- c:\Projetos\golang\mcp-gdocs-latex\docs\architecture_specification.md
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\reviewer_spec_2\handoff.md
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\challenger_spec_1\handoff.md
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\challenger_spec_2\handoff.md

IMPORTANTE: Todas as suas explicações, relatórios e comunicações DEVEM ser em Português do Brasil (pt-BR).

MANDATORY INTEGRITY WARNING:
DO NOT CHEAT. All implementations must be genuine. DO NOT hardcode test results, create dummy/facade implementations, or circumvent the intended task. A teamwork_auditor will independently verify your work. Integrity violations WILL be detected and your work WILL be rejected.

PROPRIEDADE EXCLUSIVA DE ESCRITA:
Você é o único agente autorizado a escrever em:
- c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md
- c:\Projetos\golang\mcp-gdocs-latex\docs\architecture_specification.md
- Diretório c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_spec_2\

SUA MISSÃO:
Atualizar e aprimorar cirurgicamente `c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md` e `c:\Projetos\golang\mcp-gdocs-latex\docs\architecture_specification.md`, sanando com precisão os 8 pontos identificados pelos validadores:
1. `DriveService`: Adicionar o método `ListFiles(ctx context.Context, query string, fields string, pageSize int64, pageToken string) (*drive.FileList, error)` na interface em `PROJECT.md` e `docs/architecture_specification.md`.
2. `TableInsertionResult`: Declarar formalmente a struct em `PROJECT.md` e `docs/architecture_specification.md` com campos `TableStartIndex`, `TableEndIndex`, `RowsCount`, `ColsCount`, `Message`.
3. Two-Step Approach de Tabelas:
   - Adicionar `WriteControl: &docs.WriteControl{RequiredRevisionId: doc.RevisionId}` no Passo 3 do preenchimento para travar concorrência.
   - Especificar a sintaxe correta de `UpdateTableCellStyleRequest` contendo `TableStartLocation: &docs.Location{Index: tableStartIndex}` e `TableRange`.
   - Documentar o sequenciamento da Legenda (`Caption`) no Passo 1 antes da criação da tabela e da Fonte (`Source`) no Passo 3 após a tabela.
   - No parser de LaTeX (`FEAT-08`), especificar expressamente o tratamento e normalização de `\addlinespace`, citações `\textcite{...}`, caracteres escapados `\%` e fórmulas inline `$2 \times 2$`.
4. Diagramas Mermaid:
   - Substituir o tamanho estático fixo por cálculo dinâmico de aspect ratio: decodificar dimensões via `image.DecodeConfig` e calcular altura em pontos tipográficos proporcional à largura ($W \le 450.0\text{ pt}$).
   - Incluir fallback em formato de Deep Link textual caso a permissão pública no Drive seja bloqueada por política corporativa.
5. Blindagem Defensiva de Textos e Estilos:
   - Clamping de fim de documento: limitar `EndIndex = min(EndIndex, body.EndIndex - 1)` para impedir erro HTTP 400 no caractere `\n` final inamovível.
   - No `UpdateParagraphStyleRequest`, aplicar `indentStart: { magnitude: 0, unit: "PT" }` e incluir no `fields` mask para anular recuo residual de lista deixado por `DeleteParagraphBulletsRequest`.
   - Adicionar guarda para substituições vazias (`new_text == ""`), emitindo somente `DeleteContentRangeRequest`.
   - Adicionar normalização Unicode NFC (`norm.NFC.String(s)`).
6. Sincronização de PDF e Integridade:
   - Adicionar validação do marcador `%EOF` nos últimos 1024 bytes do PDF local antes do cálculo de hash MD5, prevenindo leitura de PDFs incompletos.
   - Especificar política de segurança de revisões no Drive (`keepRevisionForever: true`) para rollback em caso de falha pós-upload.
7. Mock Server (`httptest.Server`):
   - Mapear e especificar o roteamento de `/upload/drive/v3/files/` além de `/drive/v3/files/` no mock server em Go.
8. Catálogo de Ferramentas:
   - Ajustar a contagem formal para 17 ferramentas MCP (12 pré-existentes + 5 novas).
