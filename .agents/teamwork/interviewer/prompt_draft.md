Diagnóstico aprofundado e saneamento de quebras de linha espúrias em títulos de seção e entre parágrafos no servidor MCP mcp-gdocs-latex, com implementação de correções em Go, validação no documento de teste do Google Docs e elaboração de checklist de provas de conformidade ABNT.

Working directory: c:\Projetos\golang\mcp-gdocs-latex
Integrity mode: development

## Requirements

### R1. Diagnóstico e Correção de Quebras de Linha Espúrias no Servidor MCP
Investigar no histórico de pesquisas em `.agents/teamwork` e no documento Google Docs `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` a causa raiz pela qual quebras de linha em excesso entre parágrafos e nos títulos de seção (com foco crítico na seção de Metodologia) não são sinalizadas por `lint_document_abnt` nem saneadas por `format_document_abnt`. Implementar no código Go (`internal/gdocs`) as melhorias necessárias no linter e no formatador para detectar e remover essas quebras espúrias de forma robusta e idempotente.

### R2. Auditoria e Conformidade com o Guia ABNT e Plano de Engenharia
Validar e garantir que as ferramentas do servidor MCP atendam plenamente aos padrões do plano de engenharia aprovado e às normas ABNT:
- Margens: Superior 3cm, Esquerda 3cm, Inferior 2cm, Direita 2cm.
- Alinhamento justificado e entrelinhas 1,5 no corpo; espaçamento simples (1,0) em citações longas, notas, referências, tabelas e legendas.
- Recuo de primeira linha de 1,25 cm no corpo.
- Citações diretas longas (> 3 linhas): recuo de 4 cm, tamanho 10 pt e entrelinhas simples.
- Capa institucional e Folha de rosto NBR 14724.
- Numeração de páginas e sumário estruturado.

### R3. Restrição de Infraestrutura e Modificação do Documento
Todas as atualizações e correções no documento de teste do Google Docs (`1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8`) devem ser executadas exclusivamente por meio das ferramentas expostas pelo servidor MCP (`format_document_abnt`, etc.), sendo vedado o uso de scripts externos avulsos para mutação direta do documento.

### R4. Lista de Tarefas de Verificação e Evidências Objetivas de Comprovação
Elaborar uma lista estruturada de tarefas cobrindo todos os requisitos ABNT e os componentes do plano aprovado, acompanhada do procedimento objetivo e reproduzível para comprovar que cada item foi executado com sucesso e está conforme.

## Verification Resources
- Documento Google Docs de teste: `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` ("TCC Teste IA")
- Código-fonte do servidor MCP: `c:\Projetos\golang\mcp-gdocs-latex`
- Plano de engenharia de referência: `C:\Users\joseg\.gemini\antigravity\brain\0c400b2f-aaca-4279-9a2f-8ec020311d37\plano_engenharia_abnt_mcp_gdocs.md`
- Histórico de pesquisas da equipe: `c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork`

## Acceptance Criteria

### Diagnóstico e Implementação no Código MCP
- [ ] Causa técnica da não detecção/não remoção de quebras de linha em títulos e parágrafos contíguos (especialmente na Metodologia) identificada e corrigida nos pacotes Go do MCP.
- [ ] Todos os testes unitários Go do repositório (`go test ./...`) executam e passam com sucesso (sem falhas ou regressões).

### Execução e Validação no Documento Google Docs
- [ ] Ferramenta `format_document_abnt` executada no documento `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8`, eliminando as quebras de linha espúrias em títulos de seção e entre parágrafos.
- [ ] Ferramenta `lint_document_abnt` executada após a formatação, reportando 0 anomalias normativas para margens, entrelinhas, tabelas e quebras espúrias.
- [ ] Nenhuma mutação no documento realizada fora dos endpoints e ferramentas do servidor MCP.

### Checklist e Comprovação de Tarefas
- [ ] Lista de tarefas finalizada contendo 100% dos requisitos do plano e do Guia ABNT, com evidências concretas (saída de comandos, logs do linter e trechos inspecionados) provando a conformidade de cada item.
