## 2026-10-29T19:16:10Z
Você é o worker_spec_1 (teamwork_worker).
Seu diretório de trabalho é: c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_spec_1\
Sua identidade: teamwork_worker
LEIA OBRIGATORIAMENTE O ARQUIVO: c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\ORIGINAL_REQUEST.md
LEIA TAMBÉM:
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\spec_miner_api_1\analysis.md e handoff.md
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\explorer_codebase_1\analysis.md e handoff.md
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\explorer_tcc_1\analysis.md e handoff.md

IMPORTANTE: Todas as suas explicações, relatórios e comunicações DEVEM ser em Português do Brasil (pt-BR).

MANDATORY INTEGRITY WARNING:
DO NOT CHEAT. All implementations must be genuine. DO NOT hardcode test results, create dummy/facade implementations, or circumvent the intended task. A teamwork_auditor will independently verify your work. Integrity violations WILL be detected and your work WILL be rejected.

PROPRIEDADE EXCLUSIVA DE ARQUIVOS (Write ownership):
Você é o único agente autorizado a escrever nos seguintes arquivos:
- c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md
- c:\Projetos\golang\mcp-gdocs-latex\docs\architecture_specification.md
- Arquivos no seu diretório de trabalho c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_spec_1\

SUA MISSÃO:
Elaborar a especificação arquitetural exaustiva e o plano detalhado de implementação para o servidor Go `mcp-gdocs-latex` cobrindo integralmente todos os requisitos (R1, R2, R3, R4, R5) e critérios de aceite da solicitação original:

1. Criar o arquivo `c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md` estruturado segundo o padrão Project Pattern:
   - Arquitetura geral do servidor Go e ecossistema MCP / Google Workspace.
   - Feature Inventory exaustiva: mapear todas as funcionalidades necessárias, associando cada uma ao seu Milestone correspondente.
   - Milestones detalhados (M1: Correções de Formatação e Estilos ABNT; M2: Tabelas Nativas e Diagramas Mermaid; M3: Sincronização de PDF e Integridade MD5; M4: Framework de Testes com Mocks e Refatoração de Interfaces; M5: Validação E2E com Artefatos do TCC).
   - Interface Contracts completos (assinaturas em Go, structs de payload de tools MCP, contratos de lote `batchUpdate`).
   - Code Layout detalhado para o projeto Go.

2. Criar o documento `c:\Projetos\golang\mcp-gdocs-latex\docs\architecture_specification.md` com a especificação técnica formal e cirúrgica:
   - R1: Especificação das APIs Google Docs v1 e Drive v3, endpoints, regras de lote `batchUpdate`, offsets UTF-16, atomicidade, e práticas de servidores MCP em Go (`mark3labs/mcp-go`).
   - R2: Solução definitiva para Bold Bleeding e List Bleeding. Sequência atômica exata de requests (`DeleteContentRangeRequest`, `InsertTextRequest`, `DeleteParagraphBulletsRequest`, `UpdateParagraphStyleRequest(NORMAL_TEXT)`, `UpdateTextStyleRequest(bold:false)`). Correção do bug de `internal/gdocs/writer.go:292`. Especificação métrica ABNT: recuo 1,25 cm (35,433 pt), espaçamento 1,5 (150%), justificado (`JUSTIFIED`), margens 3/3/2/2 cm (85,039 pt / 56,693 pt via `UpdateDocumentStyleRequest`).
   - R3: Arquitetura de inserção de tabelas nativas: conversão das 22 tabelas LaTeX mapeadas em `resources/tcc`, Two-Step approach (`InsertTableRequest` -> `documents.get` para descoberta de índices de células -> preenchimento via `InsertTextRequest` em lote reverso), estilização ABNT/IBGE sem bordas verticais (`UpdateTableCellStyleRequest`). Pipeline de conversão e inserção dos 4 diagramas Mermaid (`.mmd`/`.pdf` -> PNG 300 DPI -> upload Drive com permissão pública -> link direto `webContentLink` -> `InsertInlineImageRequest`).
   - R4: Especificação algorítmica da verificação de integridade por hash MD5 do PDF compilado (`resources/tcc/main.pdf`, 1.098.114 bytes) contra o campo `md5Checksum` do arquivo Google Drive (`1mbBbWNz-N1P4n7SYnj28i9yHCFz9PqZH`). Idempotência e eliminação de uploads cegos.
   - R5: Plano de implementação em etapas para o código Go, detalhando o Adapter Pattern / Interfaces (`DocsService`, `DriveService`), mocks com `httptest.Server` para 100% de cobertura nos testes unitários e de integração, e procedimento de validação prática com o TCC de `resources/tcc`.

3. Gerar o relatório de handoff formal em `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_spec_1\handoff.md` contendo:
   - Observation
   - Logic Chain
   - Caveats
   - Conclusion
   - Verification Method

Envie uma mensagem ao orquestrador ao concluir.
