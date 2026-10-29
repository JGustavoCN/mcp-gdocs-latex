# DISPATCH — Survey Subagent 1 (spec_miner_api_1)
Timestamp: 2026-10-29T19:02:30Z

## Tarefa
Pesquisar formalmente e mapear as especificações das APIs Google Docs v1 e Google Drive v3, além das convenções de servidores MCP (Model Context Protocol), aplicadas ao problema de sincronização e formatação com LaTeX.
Leia C:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\ORIGINAL_REQUEST.md.


## 2026-10-29T19:03:00Z
[Message] sender=c37f0df8-888e-4a9d-8960-88e6b44fcd03 priority=MESSAGE_PRIORITY_HIGH
Você é o spec_miner_api_1 (teamwork_spec_miner).
Seu diretório de trabalho é: c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\spec_miner_api_1\
Sua identidade: teamwork_spec_miner
LEIA OBRIGATORIAMENTE O ARQUIVO: c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\ORIGINAL_REQUEST.md
IMPORTANTE: Todas as suas explicações, relatórios e comunicações DEVEM ser em Português do Brasil (pt-BR).

Sua missão é minerar e documentar com precisão cirúrgica as especificações técnicas das APIs do Google Docs (v1) e Google Drive (v3), além das melhores práticas para servidores MCP (Model Context Protocol) em Go:

1. Google Docs API (v1) batchUpdate Requests e Comportamento:
   - Mecanismos de estilo e controle de formatação: Como isolar a formatação de texto e parágrafos para impedir o efeito cascata de negrito (bold bleeding) e outros estilos ao substituir ou inserir texto. Como funcionam `UpdateTextStyleRequest` (fields, TextStyle), `UpdateParagraphStyleRequest` (fields, ParagraphStyle, NamedStyleType, NORMAL_TEXT), `InsertTextRequest`, `DeleteContentRangeRequest`.
   - Como garantir que novos parágrafos inseridos resetem para NORMAL_TEXT ou herdem as propriedades corretas sem propagar estilos de títulos.
   - Listas e bullets: Como `CreateParagraphBulletsRequest` e `DeleteParagraphBulletsRequest` operam. Como evitar a propagação de bullets vazios quando novas linhas ou quebras de parágrafo são inseridas.
   - Higienização de quebras de linha e espaços extras: Gestão de índices de caracteres (startIndex, endIndex, utf-16 offset rules do Docs), deleção segura de intervalos de espaços/quebras redundantes.
   - Especificação de Formatação ABNT via API:
     * Recuo de primeira linha de 1,25 cm (converter cm para pontos: 1 cm ≈ 28.3465 pt -> 1.25 cm ≈ 35.433 pt no campo `indentFirstLine.magnitude` com unit `PT`).
     * Espaçamento entre linhas 1,5 (`lineSpacing: 150` porcentagem do normal).
     * Alinhamento justificado (`alignment: "JUSTIFIED"`).
     * Margens da página: 3 cm superior e esquerda (85.039 pt), 2 cm inferior e direita (56.693 pt) através de `UpdateDocumentStyleRequest` (`marginTop`, `marginBottom`, `marginLeft`, `marginRight`).
   - Tabelas Nativas no Google Docs:
     * Como usar `InsertTableRequest` (rows, columns, location / endOfSegment).
     * Estrutura de dados de uma Table no documento Google Docs (Table, TableRow, TableCell, Content).
     * Como determinar os índices de início e fim de cada célula após a inserção para preenchimento de texto via `InsertTextRequest`.
     * Formatação e bordas de células (`UpdateTableCellStyleRequest`, `UpdateTableColumnPropertiesRequest`).
   - Imagens Inline:
     * `InsertInlineImageRequest` (uri, location, objectSize).
     * Restrição da Google Docs API sobre a URI da imagem (requer URI acessível pela infraestrutura do Google, ou fluxo via Google Drive com permissão pública/temporária ou webContentLink).

2. Google Drive API (v3):
   - Metadados de integridade de arquivos: endpoint `files.get` com `fields="id, name, md5Checksum, size, modifiedTime, webViewLink, webContentLink"`.
   - Especificação do campo `md5Checksum` (retornado apenas para arquivos binários/não-Google Docs, como PDF).
   - Fluxo de upload/atualização (`files.create` e `files.update`) multipart ou resumable com verificação de hash MD5 local vs `md5Checksum` da nuvem.
   - Fluxo de hospedagem temporária de imagens (diagramas convertidos) no Drive para consumo pelo `InsertInlineImageRequest` do Docs.

3. Práticas de Servidores MCP em Go:
   - Estrutura de Tools no Model Context Protocol, retorno de conteúdo TextContent/ImageContent, tratamento de erros, validação de schema de parâmetros.
   - Padrão arquitetural em Go para compor e executar lotes de `batchUpdate` atômicos com rollback ou tratamento de falha.

Entregáveis:
- Escreva a análise completa e detalhada em: c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\spec_miner_api_1\analysis.md
- Escreva o relatório de handoff em: c:\Projetos\golang\mcp-gdocs-latex\.agents\teamwork\spec_miner_api_1\handoff.md
- Envie uma mensagem concisa ao orquestrador com resumo dos achados e caminhos dos arquivos gerados.
