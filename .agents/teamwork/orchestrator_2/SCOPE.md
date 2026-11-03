# Scope: Diagnóstico e Saneamento de Quebras de Linha Espúrias e Conformidade ABNT

## Architecture
- Servidor MCP Go (`cmd/mcp-gdocs-latex`, `internal/gdocs`)
- Linter ABNT (`internal/gdocs/linter.go`)
- Formatador ABNT (`internal/gdocs/abnt_writer.go`, `writer.go`, `cover_generator.go`)
- Inspetor e Leitor (`internal/gdocs/inspector.go`, `reader.go`)
- Documento Google Docs de teste (`1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8`)

## Feature Inventory
| # | Feature | Description | Milestone | Source |
|---|---------|-------------|-----------|--------|
| 1 | Detecção de quebras espúrias em títulos de seção | Identificar quebras de linha/parágrafos vazios em títulos de seção (especialmente Metodologia) | M1 | R1 |
| 2 | Detecção de quebras espúrias entre parágrafos contíguos | Identificar quebras de linha excessivas entre parágrafos de texto no linter | M1 | R1 |
| 3 | Saneamento idempotente de quebras espúrias no formatador | Remover parágrafos vazios espúrios em títulos e parágrafos sem quebrar preâmbulo/capa/seções | M2 | R1 |
| 4 | Testes unitários para linter e formatador de quebras | Garantir que go test ./... cubra as novas regras e passe com 100% de sucesso | M3 | R1 |
| 5 | Execução e validação no Google Docs de teste | Rodar format_document_abnt e lint_document_abnt no documento 1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8 | M4 | R3 |
| 6 | Auditoria de conformidade e checklist ABNT | Checklist exaustivo com evidências de margens, entrelinhas, recuo, citações, capa/folha de rosto e sumário | M5 | R2, R4 |

## Milestones
| # | Name | Scope | Dependencies | Status |
|---|------|-------|-------------|--------|
| 1 | Diagnóstico e Investigação | Causa raiz em linter/formatador e documento de teste | none | DONE |
| 2 | Correção no Código Go | Ajustes em internal/gdocs/linter.go e abnt_writer.go | M1 | DONE |
| 3 | Testes Unitários Go | Testes e validação de regressão com go test ./... | M2 | DONE |
| 4 | Execução no Google Docs | Formatação e linting no documento real via MCP | M3 | IN_PROGRESS |
| 5 | Checklist de Conformidade ABNT | Provas e evidências de conformidade ABNT e plano | M4 | IN_PROGRESS |

## Interface Contracts
- `LinterService` / `LintDocumentABNT`: detecção precisa de quebras espúrias (`SpuriousLineBreakAnomaly`).
- `ABNTFormatterService` / `FormatDocumentABNT`: remoção atômica via `DeleteContentRangeRequest` calculando índices corretos em ordem reversa ou lote idempotente.
