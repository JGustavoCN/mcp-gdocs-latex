## 2026-10-29T19:56:24Z
Você é o worker_spec_3 (teamwork_worker) na Iteração 3 de polimento final da especificação técnica e arquitetural.
Seu diretório de trabalho é: c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_spec_3\
Sua identidade: teamwork_worker
LEIA OBRIGATORIAMENTE:
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\ORIGINAL_REQUEST.md
- c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md
- c:\Projetos\golang\mcp-gdocs-latex\docs\architecture_specification.md
- c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\challenger_spec_4\handoff.md

IMPORTANTE: Todas as suas explicações, relatórios e comunicações DEVEM ser em Português do Brasil (pt-BR).

MANDATORY INTEGRITY WARNING:
DO NOT CHEAT. All implementations must be genuine. DO NOT hardcode test results, create dummy/facade implementations, or circumvent the intended task. A teamwork_auditor will independently verify your work. Integrity violations WILL be detected and your work WILL be rejected.

PROPRIEDADE EXCLUSIVA DE ESCRITA:
Você é o único agente autorizado a escrever em:
- c:\Projetos\golang\mcp-gdocs-latex\PROJECT.md
- c:\Projetos\golang\mcp-gdocs-latex\docs\architecture_specification.md
- Diretório c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_spec_3\

SUA MISSÃO:
Aplicar cirurgicamente as 4 correções finais em `PROJECT.md` e `docs/architecture_specification.md`:
1. Two-Step de Tabelas (Passo 3):
   - Mudar a ordenação do Passo 3 para que a inserção da Fonte no `tableEndIndex` seja a PRIMEIRA operação adicionada ao lote (maior offset), antes do loop de preenchimento reverso das células (índices menores). Isso garante que quando as células forem preenchidas, a fonte já terá sido inserida na sua coordenada final correta, eliminando o risco de que o texto da fonte caia dentro de uma célula da tabela.
2. Diagramas Mermaid (Enquadramento Bidimensional / 2D Bounding Box):
   - Substituir o cálculo 1D por enquadramento 2D: $\text{maxWidth} = 450.0\text{ pt}$ e $\text{maxHeight} = 620.0\text{ pt}$.
   - Definir escala de contenção proporcional: se a altura calculada a 450 pt exceder 620 pt (como no caso de `sequencia.mmd` e `arquitetura.mmd`), calcular $\text{scale} = \min(1.0, \min(450.0 / W_{\text{pt}}, 620.0 / H_{\text{pt}}))$ garantindo que a imagem caiba perfeitamente na página A4 com margens ABNT, legenda e fonte.
   - Adicionar guarda preventiva contra divisão por zero (`cfg.Width <= 0 || cfg.Height <= 0`).
   - Declarar expressamente o import anônimo `_ "image/png"` no Go para registro do codec PNG.
3. Sincronização de PDF (%EOF e TOCTOU):
   - Padronizar o marcador formal de encerramento do PDF para `[]byte("%%EOF")` (dois percentuais conforme ISO 32000-1).
   - Adicionar checagem pós-hashing de integridade concorrente: validar se `copiedBytes == stat.Size()`, prevenindo leitura durante compilações simultâneas do LaTeX.
4. Mock Server Multipart:
   - Explicitar no handler do mock server que a rota de upload `/upload/drive/v3/files/` processa conteúdo `multipart/related` e extrai o `fileID` ignorando o prefixo `/upload`.

Entregáveis:
- Arquivos `PROJECT.md` e `docs/architecture_specification.md` devidamente atualizados.
- Relatório de handoff formal em `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\worker_spec_3\handoff.md`.
- Envie mensagem ao orquestrador ao concluir.
