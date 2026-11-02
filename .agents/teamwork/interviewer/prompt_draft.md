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
