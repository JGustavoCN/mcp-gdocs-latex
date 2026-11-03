# BRIEFING — 2026-11-02T22:35:00Z

## Mission
Implementar correções profundas no linter e formatador ABNT do mcp-gdocs-latex, cobrir com testes unitários, compilar binários e validar no documento Google Docs de produção.

## 🔒 My Identity
- Archetype: teamwork_worker
- Roles: implementer, qa, specialist
- Working directory: c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_r1_1
- Original parent: c8a0340f-c506-4816-9ee0-56615cfe9c68
- Milestone: r1_implementation

## 🔒 Key Constraints
- Comunicação estritamente em Português do Brasil (pt-BR).
- Titularidade de escrita estrita: `internal/gdocs/linter.go`, `internal/gdocs/abnt_writer.go`, `internal/gdocs/structure_normalizer.go`, `internal/gdocs/*_test.go`, `bin/`.
- Proibido trapacear: implementações reais e funcionais, sem hardcoding ou fachadas.
- Todas as alterações no documento real do Google Docs devem ser feitas exclusivamente via ferramentas MCP (Restrição R3).
- Testes unitários com 100% de sucesso.

## Current Parent
- Conversation ID: c8a0340f-c506-4816-9ee0-56615cfe9c68
- Updated: 2026-11-02T22:35:00Z

## Task Summary
- **What to build**:
  1. `internal/gdocs/linter.go`: auditoria de quebras em branco unitárias no corpo textual (`consecutiveEmptyParas >= 1`), proteção de delimitador de tabela, proteção de `PageBreak` e `InlineObjectElement`, auditoria de headings fantasmas, falsos títulos com `:` e hierarquia NBR 6024.
  2. `internal/gdocs/abnt_writer.go`: remoção atômica em lote reverso de parágrafos vazios em `HealDocTextAndBlanks`, proteção de `PageBreak`/`InlineObject`, proteção e neutralização de estilo de parágrafo pré-tabela (1pt, lineSpacing 100, space 0), precedência de regexes hierárquicas sobre `HEADING_1`, rebaixamento de falsos títulos com `:` para `NORMAL_TEXT`, e inserção determinística de `PageBreak` antes de capítulos primários e REFERÊNCIAS em `EnsureChapterPageBreaks`.
  3. Testes unitários em `internal/gdocs/linter_rules_test.go` e `internal/gdocs/abnt_heal_test.go`.
  4. Compilação de `bin/mcp-server-v2.exe` e `bin/mcp-server.exe`.
  5. Validação com o documento real e registro de tentativa de execução via MCP.
- **Success criteria**: 100% de aprovação na suíte de testes unitários (`go test -v ./...`), binários atualizados em `bin/`, relatório de handoff estruturado.

## Change Tracker
- **Files modified**:
  - `internal/gdocs/linter.go`: Adicionada detecção de headings fantasmas, hierarquia NBR 6024, quebras unitárias e proteção de PageBreak/InlineObject/isBeforeTable.
  - `internal/gdocs/abnt_writer.go`: Saneamento atômico em lote reverso (até 100 reqs por batch), proteção de PageBreak e neutralização pré-tabela (1pt, lineSpacing 100, space 0), precedência de regexes sobre HEADING_1, rebaixamento de orações com `:` e inserção determinística de PageBreak em capítulos primários e referências (`EnsureChapterPageBreaks`).
  - `internal/gdocs/linter_rules_test.go`: Novos testes unitários para detecção de headings fantasmas, falsos títulos com `:`, hierarquia NBR 6024 e regras de quebra unitária no corpo.
  - `internal/gdocs/abnt_heal_test.go`: Novos testes unitários para proteção de PageBreak/InlineObject, neutralização pré-tabela e ordenação reversa de lotes de deleção.
  - `bin/mcp-server-v2.exe`: Binário oficial do MCP compilado e atualizado.
  - `bin/mcp-server.exe`: Binário oficial secundário compilado e atualizado.
- **Build status**: PASS (100% dos testes unitários passando em `go test -v ./...`).
- **Pending issues**: Nenhuma pendência de implementação.

## Quality Status
- **Build/test result**: PASS (`go test -v ./...` executou com sucesso em todos os pacotes).
- **Lint status**: 0 violações introduzidas; código formatado segundo padrão Go.
- **Tests added/modified**: `TestLinterGhostHeadingsAndHierarchy`, `TestHeadingPrecedenceRule`, `TestLinterSpuriousBlankLinePreservationRules`, `TestPageBreakAndTableDelimiterProtectionInHeal`.

## Key Decisions Made
- Neutralização do delimitador pré-tabela é executada preventivamente antes do lote de deleção, evitando deslocamento de índices.
- Lotes reversos de deleção agrupados em fatias de até 100 requisições para evitar exceder cotas ou limites de payload da Docs API.
- Proteção categórica de `PageBreak` e `InlineObjectElement` em todas as fases do linter e do saneador.
- Obediência estrita à recusa de permissão para ferramentas mutativas MCP interativas no ambiente de teste, documentando o status de forma transparente.

## Artifact Index
- `.agents/teamwork/worker_r1_1/DISPATCH.md` — Assignment do orquestrador
- `.agents/teamwork/worker_r1_1/BRIEFING.md` — Memória persistente
- `.agents/teamwork/worker_r1_1/progress.md` — Heartbeat de progresso
- `.agents/teamwork/worker_r1_1/handoff.md` — Relatório final
