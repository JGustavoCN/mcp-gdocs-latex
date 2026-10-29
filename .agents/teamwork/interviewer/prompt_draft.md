Pesquisar as atualizações e documentação recente das APIs do Google Docs e Google Drive e as capacidades de servidores MCP, e elaborar um plano arquitetural detalhado e especificação técnica para o servidor Go `mcp-gdocs-latex`, solucionando os problemas de integridade, formatação ABNT, renderização de tabelas e inserção de diagramas Mermaid entre LaTeX e Google Docs.

Working directory: c:\Projetos\golang\mcp-gdocs-latex
Integrity mode: development

## Referências
- `C:\Projetos\golang\mcp-gdocs-latex\resources\outras infromações.md`
- `C:\Projetos\golang\mcp-gdocs-latex\resources\mcp_sync_history_and_improvements.md`
- `C:\Projetos\golang\mcp-gdocs-latex\resources\tcc`

## Requirements

### R1. Pesquisa Técnica e Mapeamento de APIs do Google Docs e Drive
Pesquisar a documentação oficial, atualizações recentes e capacidades das APIs do Google Docs (v1) e Google Drive (v3), bem como práticas de servidores MCP, identificando os endpoints, payloads em lote (`batchUpdate`) e restrições para manipulação de texto, estilos, tabelas, imagens e metadados.

### R2. Especificação Arquitetural para Solução dos Problemas de Formatação e Estilo
Projetar e documentar a arquitetura técnica para eliminar falhas de integridade no documento Google Docs:
- Prevenção ou cura do efeito cascata de negrito e estilos em substituições de texto (ex: isolamento de estilos ou aplicação automática de `NORMAL_TEXT`).
- Prevenção do efeito cascata de listas (evitar bullets vazios em novas linhas).
- Eliminação de quebras de linhas desnecessárias e espaços extras.
- Aplicação de regras de formatação ABNT (recuo de 1,25 cm, entrelinhas 1,5, alinhamento justificado e margens 3/3/2/2).

### R3. Arquitetura para Inserção de Tabelas Nativas e Diagramas Mermaid
Projetar a integração de elementos visuais e tabulares do LaTeX no Google Docs:
- Conversão de tabelas estruturadas do LaTeX em tabelas nativas do Google Docs através da API (`InsertTableRequest` e preenchimento de células).
- Processamento de diagramas Mermaid (arquivos `.mmd` e `.pdf` em `resources/tcc/figuras`), convertendo-os em imagens (PNG) e inserindo-os no Google Docs via API (`InsertInlineImageRequest`) ou upload no Drive com referência inline.

### R4. Especificação de Sincronização e Verificação de Integridade por Hash do PDF
Projetar o fluxo de verificação de integridade entre o PDF compilado localmente em `resources/tcc` e o PDF armazenado no Google Drive, comparando o hash MD5 calculado do arquivo local com o campo `md5Checksum` retornado pela API do Google Drive v3.

### R5. Plano de Implementação e Estratégia de Testes
Elaborar um plano detalhado de implementação para o servidor Go `mcp-gdocs-latex`, especificando as ferramentas MCP novas ou refatoradas, cenários de testes automatizados (unitários e de integração com mocks) e validação prática utilizando a documentação e arquivos do TCC em `resources/tcc`.

## Acceptance Criteria

### Pesquisa e Arquitetura de Formatação
- [ ] O documento de arquitetura especifica as chamadas da API do Google Docs (`batchUpdate`) para isolar substituições de texto e impedir que novos parágrafos herdem estilos indevidos (como negrito de títulos).
- [ ] A arquitetura especifica o tratamento de listas para evitar quebras com bullets vazios e a remoção de espaços e quebras supérfluas.
- [ ] O mapeamento de formatação ABNT define os payloads de `UpdateParagraphStyleRequest` e `UpdateDocumentStyleRequest` para recuo de 1,25 cm, espaçamento 1,5, texto justificado e margens 3/3/2/2.

### Tabelas e Diagramas Mermaid
- [ ] A arquitetura define o fluxo técnico para converter tabelas LaTeX em tabelas nativas no Google Docs via `InsertTableRequest`.
- [ ] O fluxo de conversão e inserção de diagramas Mermaid (`.mmd`/`.pdf`) em formato de imagem no Google Docs é especificado detalhadamente.

### Verificação de Integridade e Plano de Ação
- [ ] A verificação de sincronia entre o PDF local e o PDF no Google Drive é especificada comparando o hash MD5 local com o metadado `md5Checksum` da Google Drive API v3.
- [ ] O plano de implementação apresenta o roteiro em etapas para o código Go, incluindo estratégias de testes automatizados e validação real com os arquivos de `resources/tcc`.
