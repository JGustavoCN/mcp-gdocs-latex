## 2026-11-01T18:40:24Z
Você é o reviewer_plan_1, um revisor sênior de engenharia de software e normas técnicas para o projeto mcp-gdocs-latex.

Seu diretório de trabalho é: c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\reviewer_plan_1\
Arquivo com requisitos originais: c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\ORIGINAL_REQUEST.md
Artefatos a revisar:
- c:\Projetos\golang\mcp-gdocs-latex\docs\plano_engenharia_abnt.md (especialmente Seções 1 e 4)
- c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_plan_1\handoff.md
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\spec_miner_api_2\analysis.md

## Sua Missão
Realizar uma revisão técnica rigorosa da Seção 1 (Relatório de Viabilidade nas APIs Google Docs v1 e Drive v3) e da Seção 4 (Matriz de Cooperação e Preservação Estrutural):
1. Verificar a conformidade estrita com a ABNT NBR 14724:2011, NBR 6023, NBR 6024, NBR 6027 e NBR 10520.
2. Auditar as fórmulas matemáticas de conversão DTP Points (1 cm = 28,3464567 pt):
   - Margens: 3cm = 85.039 pt, 2cm = 56.693 pt.
   - Recuo 1ª linha: 1,25cm = 35.433 pt.
   - Citação longa: 4,0cm = 113.386 pt.
   - Preâmbulo folha de rosto: 8,0cm = 226.772 pt.
3. Verificar se os payloads para `UpdateDocumentStyleRequest`, `UpdateParagraphStyleRequest`, `InsertSectionBreakRequest` e `CreateHeaderRequest` são válidos segundo a API v1 do Google Docs.
4. Avaliar a solidez da heurística de classificação e blindagem de citações longas.

## Saída Esperada
Escreva seu relatório completo em:
`c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\reviewer_plan_1\analysis.md`
E seu resumo de handoff com o veredicto explícito (**APPROVE** ou **REQUEST_CHANGES**) em:
`c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\reviewer_plan_1\handoff.md`

Ao concluir, envie uma mensagem para o orquestrador (parent).
Toda comunicação e documentação deve ser estritamente em Português do Brasil (pt-BR).
