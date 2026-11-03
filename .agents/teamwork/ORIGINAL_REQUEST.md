# Original User Request

## Initial Request — 2026-11-01T18:15:41Z

Pesquisa de viabilidade técnica nas APIs do Google Docs e Google Drive e elaboração de plano detalhado de engenharia (/plan) para as ferramentas do MCP mcp-gdocs-latex, visando conformidade estrita com normas ABNT.

Working directory: c:\Projetos\golang\mcp-gdocs-latex
Integrity mode: development

## Requirements

### R1. Pesquisa de Viabilidade Técnica na API do Google Docs e Drive
Investigar as capacidades e restrições da API do Google Docs e do Google Drive para cada uma das normas ABNT do Guia de Referência:
- Ajuste de margens NBR 14724 (Superior 3cm, Esquerda 3cm, Inferior 2cm, Direita 2cm).
- Alinhamento justificado e espaçamento entrelinhas (1,5 no corpo com spaceAbove=0/spaceBelow=0; 1,0 em citações longas, notas, referências, legendas e tabelas).
- Recuo de primeira linha de 1,25 cm no corpo de texto.
- Citação direta longa (> 3 linhas): recuo de 4,0 cm (113.386 pt), tamanho 10 pt, espaçamento 1,0 simples, sem aspas.
- Capa institucional e Folha de rosto NBR 14724 com preâmbulo recuado a 8,0 cm (226.77 pt), tamanho 10 pt e espaçamento simples.
- Numeração de páginas no canto superior direito a partir da introdução, suprimindo o número nas páginas pré-textuais (investigar SectionBreak, Header/Footer e FirstPageDifferent).
- Sumário automático (TableOfContents nativo via Docs API vs. gerador estruturado ABNT).
- Remoção programática de quebras de linha em excesso após a folha de rosto e entre parágrafos.

### R2. Diagnóstico das Ferramentas MCP Atuais
Examinar a implementação atual dos pacotes Go (`internal/gdocs/abnt_writer.go`, `internal/gdocs/linter.go`, `internal/gdocs/inspector.go`, `internal/gdocs/cover_generator.go`, etc.) e avaliar o documento de teste `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` para identificar:
- Por que o linter e o formatador não identificaram nem corrigiram variações indesejadas de entrelinhas (ex.: 1,15 vs 1,5).
- Por que as quebras de linha espúrias após a folha de rosto não foram saneadas.
- O que falta nas ferramentas de leitura e inspeção (`read_doc_content`, `inspect_doc_layout`) para expor metadados visuais e estruturais suficientes para que modelos de IA consigam enxergar e diagnosticar esses problemas.

### R3. Plano Detalhado de Melhoria (/plan) por Ferramenta
Elaborar um plano detalhado de engenharia (/plan) cobrindo cada ferramenta MCP (`format_document_abnt`, `lint_document_abnt`, `inspect_doc_layout`, `read_doc_content`, `cover_generator`):
- Detalhar as assinaturas, payloads, regras de validação e endpoints de batchUpdate da API do Google Docs necessários para cada ferramenta.
- Definir como as ferramentas devem cooperar para garantir que o formatador não desfaça formatações de citações diretas nem quebre a estrutura do documento.

## Verification Resources
- Documento de teste: `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8`
- Código-fonte: repositório Go em `c:\Projetos\golang\mcp-gdocs-latex`
- Guia de Referência de Normas ABNT (NBR 14724) fornecido na solicitação.

## Acceptance Criteria

### Diagnóstico e Viabilidade Técnica
- [ ] Relatório técnico documentando a viabilidade de implementação na Google Docs API para todos os 8 pontos do Guia de Referência ABNT.
- [ ] Identificação fundamentada no código atual e no documento `1IVQ271HdTfXywUhDBeXTUDDcZFP-KRvNEtJvxRy8Dp8` de por que variações de espaçamento (1,15 vs 1,5) e quebras de linha excessivas não foram sinalizadas ou corrigidas.

### Plano de Engenharia (/plan)
- [ ] Plano estruturado por ferramenta MCP com seções dedicadas para `format_document_abnt`, `lint_document_abnt`, `inspect_doc_layout`, `read_doc_content` e geradores de capa/folha de rosto.
- [ ] Especificação dos mecanismos para detecção de citações longas e proteção de sua formatação durante a formatação geral do documento.
- [ ] Diretrizes técnicas claras e acionáveis para guiar a implementação posterior no código Go sem ambiguidades.


## Follow-up — 2026-11-02T21:50:07Z

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
