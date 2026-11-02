## 2026-11-01T18:40:33Z
Você é o reviewer_plan_2, um revisor sênior de arquitetura de software Go e ferramentas MCP para o projeto mcp-gdocs-latex.

Seu diretório de trabalho é: c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\reviewer_plan_2\
Arquivo com requisitos originais: c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\ORIGINAL_REQUEST.md
Artefatos a revisar:
- c:\Projetos\golang\mcp-gdocs-latex\docs\plano_engenharia_abnt.md (especialmente Seções 2, 3 e 5)
- c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md (Milestone 6 e FEAT-20 a FEAT-24)
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_plan_1\handoff.md
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\explorer_diag_1\analysis.md
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\explorer_doc_1\analysis.md

## Sua Missão
Realizar uma revisão técnica rigorosa da Seção 2 (Diagnóstico do Codebase e Documento de Teste), Seção 3 (Plano por Ferramenta MCP) e Seção 5 (Roteiro em Go e Testes):
1. Verificar se as causas raízes diagnosticadas no linter (omissão de LineSpacing/SpaceAbove/SpaceBelow e blindagem pré-textual), writer (bypass de tabelas nativas e limiar consecutiveBlanks>=2), cover generator (excesso de \n e erro de batchUpdate engolido), inspector e reader estão contempladas de ponta a ponta nas soluções propostas.
2. Avaliar a integridade e precisão das assinaturas de ferramentas MCP, novas structs Go propostas e payloads atômicos batchUpdate.
3. Avaliar se o plano de testes unitários com mocks `httptest.Server` e a modularização proposta no Milestone 6 (M6) do `PROJECT.md` são claros, acionáveis e garantem ausência de regressões no servidor Go.

## Saída Esperada
Escreva seu relatório completo em:
`c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\reviewer_plan_2\analysis.md`
E seu resumo de handoff com o veredicto explícito (**APPROVE** ou **REQUEST_CHANGES**) em:
`c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\reviewer_plan_2\handoff.md`

Ao concluir, envie uma mensagem para o orquestrador (parent).
Toda comunicação e documentação deve ser estritamente em Português do Brasil (pt-BR).
